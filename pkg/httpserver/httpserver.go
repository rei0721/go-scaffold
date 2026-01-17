package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/rei0721/go-scaffold/pkg/logger"
	"github.com/rei0721/go-scaffold/pkg/utils"
)

// httpServer HTTP 服务器实现
type httpServer struct {
	// server 标准库 http.Server 实例
	server *http.Server

	// handler HTTP 请求处理器 (Gin Router)
	handler Handler

	// config 当前配置
	config *Config

	// logger 日志记录器
	logger logger.Logger

	// executor 协程池管理器（可选）
	// 使用 atomic.Value 实现无锁读取
	// 用于异步处理HTTP请求相关任务
	executor atomic.Value // 存储 executor.Manager

	// mu 保护并发访问（热重载使用）
	mu sync.RWMutex

	// state 服务器运行状态
	state atomic.Int32

	// errChan 服务器错误通道
	errChan chan error
}

// New 创建新的 HTTP Server 实例
// 参数:
//
//	handler: HTTP 请求处理器，通常是 Gin Router
//	cfg: 服务器配置
//	log: 日志记录器
//
// 返回:
//
//	HTTPServer: 服务器实例
//	error: 创建失败时的错误
//
// 注意:
//
//	Executor需要通过SetExecutor方法延迟注入
func New(handler Handler, cfg *Config, log logger.Logger) (HTTPServer, error) {
	if handler == nil {
		return nil, &ServerError{
			Op:      "new",
			Message: "handler cannot be nil",
		}
	}

	if cfg == nil {
		cfg = &Config{}
	}

	// 应用默认值
	cfg.ApplyDefaults()

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, &ServerError{
			Op:      "new",
			Message: ErrMsgInvalidConfig,
			Err:     err,
		}
	}

	s := &httpServer{
		handler: handler,
		config:  cfg,
		logger:  log,
		errChan: make(chan error, 1),
	}

	// 设置初始状态为已停止
	s.state.Store(int32(stateStopped))

	return s, nil
}

// Start 启动 HTTP 服务器（非阻塞）
func (s *httpServer) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查服务器是否已在运行
	currentState := serverState(s.state.Load())
	if currentState == stateRunning || currentState == stateStarting {
		return &ServerError{
			Op:      "start",
			Message: ErrMsgServerAlreadyRunning,
		}
	}

	// 设置状态为启动中
	s.state.Store(int32(stateStarting))

	// 如果端口为 0，自动分配可用端口
	if s.config.Port == 0 {
		port, err := utils.GetAvailablePort(9000, 30000)
		if err != nil {
			s.state.Store(int32(stateStopped))
			return &ServerError{
				Op:      "start",
				Message: ErrMsgPortUnavailable,
				Err:     err,
			}
		}
		s.config.Port = port
	}

	// 构造监听地址
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	// 校验地址是否合法
	if err := utils.IsValidHTTPListenAddr(addr); err != nil {
		// 如果地址不合法，使用默认 host
		s.logger.Warn("invalid listen address, using default host", "addr", addr, "error", err)
		addr = fmt.Sprintf("%s:%d", DefaultHost, s.config.Port)
	}

	// 创建 HTTP 服务器实例
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.handler,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	// 记录启动信息
	s.logger.Info(fmt.Sprintf("starting HTTP server on http://%s", addr), "addr", addr)

	// 在新的 goroutine 中启动服务器
	go func() {
		// 设置状态为运行中
		s.state.Store(int32(stateRunning))

		// 启动服务器并开始监听
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// ErrServerClosed 是正常的关闭，不是错误
			s.logger.Error("HTTP server error", "error", err)
			s.errChan <- &ServerError{
				Op:      "start",
				Message: ErrMsgServerStartFailed,
				Err:     err,
			}
			s.state.Store(int32(stateStopped))
		}
	}()

	return nil
}

// Shutdown 优雅关闭服务器
func (s *httpServer) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查服务器是否在运行
	currentState := serverState(s.state.Load())
	if currentState != stateRunning {
		s.logger.Warn("attempting to shutdown a non-running server", "state", currentState.String())
		return nil
	}

	// 设置状态为停止中
	s.state.Store(int32(stateStopping))

	s.logger.Info("shutting down HTTP server...")

	// 执行优雅关闭
	if s.server != nil {
		if err := s.server.Shutdown(ctx); err != nil {
			s.state.Store(int32(stateStopped))
			return &ServerError{
				Op:      "shutdown",
				Message: ErrMsgServerShutdownFailed,
				Err:     err,
			}
		}
	}

	// 设置状态为已停止
	s.state.Store(int32(stateStopped))
	s.logger.Info("HTTP server stopped gracefully")

	return nil
}

// Reload 热重载配置（原子操作）
func (s *httpServer) Reload(ctx context.Context, cfg *Config) error {
	if cfg == nil {
		return &ServerError{
			Op:      "reload",
			Message: ErrMsgInvalidConfig,
		}
	}

	// 应用默认值
	cfg.ApplyDefaults()

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return &ServerError{
			Op:      "reload",
			Message: ErrMsgInvalidConfig,
			Err:     err,
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	currentState := serverState(s.state.Load())

	// 如果服务器未运行，直接更新配置
	if currentState != stateRunning {
		s.config = cfg
		s.logger.Info("HTTP server config updated (server not running)")
		return nil
	}

	s.logger.Info("reloading HTTP server with new config...")

	// 保存旧服务器实例
	oldServer := s.server

	// 构造新的监听地址
	newAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	// 检查端口是否变化
	oldAddr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	portChanged := newAddr != oldAddr

	// 如果端口变化，需要先关闭旧服务器
	if portChanged {
		s.logger.Info("port changed, shutting down old server first", "old", oldAddr, "new", newAddr)

		// 创建一个临时上下文用于关闭旧服务器
		shutdownCtx, cancel := context.WithTimeout(context.Background(), DefaultWriteTimeout)
		defer cancel()

		if err := oldServer.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("failed to shutdown old server during reload", "error", err)
			return &ServerError{
				Op:      "reload",
				Message: ErrMsgReloadFailed,
				Err:     err,
			}
		}
	}

	// 更新配置
	s.config = cfg

	// 校验地址是否合法
	if err := utils.IsValidHTTPListenAddr(newAddr); err != nil {
		s.logger.Warn("invalid listen address, using default host", "addr", newAddr, "error", err)
		newAddr = fmt.Sprintf("%s:%d", DefaultHost, cfg.Port)
	}

	// 创建新的服务器实例
	s.server = &http.Server{
		Addr:         newAddr,
		Handler:      s.handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// 启动新服务器
	go func() {
		s.logger.Info(fmt.Sprintf("restarting HTTP server on http://%s", newAddr), "addr", newAddr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("reloaded HTTP server error", "error", err)
			s.errChan <- &ServerError{
				Op:      "reload",
				Message: ErrMsgReloadFailed,
				Err:     err,
			}
		}
	}()

	// 如果端口未变化，现在关闭旧服务器
	if !portChanged && oldServer != nil {
		go func() {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), DefaultWriteTimeout)
			defer cancel()

			if err := oldServer.Shutdown(shutdownCtx); err != nil {
				s.logger.Error("failed to shutdown old server after reload", "error", err)
			}
		}()
	}

	s.logger.Info("HTTP server reloaded successfully")
	return nil
}
