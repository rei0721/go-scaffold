package daemons

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rei0721/rei0721/pkg/daemon"
	"github.com/rei0721/rei0721/pkg/scheduler"
)

// HTTPDaemon HTTP 服务器守护进程
// 将 HTTP 服务器封装为 Daemon,便于统一管理
//
// 职责:
//   - 启动 HTTP 服务器
//   - 优雅关闭 HTTP 服务器
//   - 记录日志
//
// 使用示例:
//
//	router := gin.Default()
//	httpDaemon := daemons.NewHTTPDaemon(":8080", router, logger)
//	manager.Register(httpDaemon)
type HTTPDaemon struct {
	// server HTTP 服务器实例
	// 使用标准库的 http.Server
	server *http.Server

	// logger 日志记录器
	// 用于记录服务器的启动、停止和错误信息
	logger daemon.Logger

	scheduler scheduler.Scheduler

	// name 守护进程名称
	// 用于日志记录和错误报告
	name string
}

// NewHTTPDaemon 创建一个新的 HTTP 守护进程
//
// 参数:
//
//	addr: 服务器监听地址
//	  - 格式: "host:port"
//	  - 例如: ":8080" (监听所有接口的 8080 端口)
//	  - 例如: "localhost:8080" (只监听本地的 8080 端口)
//	handler: HTTP 请求处理器
//	  - 通常是 Gin 的 Engine: gin.Default()
//	  - 或者标准库的 http.Handler
//	logger: 日志记录器
//	  - 必须实现 daemon.Logger 接口
//
// 返回:
//
//	*HTTPDaemon: HTTP 守护进程实例
//
// 使用示例:
//
//	// 使用 Gin
//	router := gin.Default()
//	router.GET("/ping", func(c *gin.Context) {
//	    c.JSON(200, gin.H{"message": "pong"})
//	})
//	httpDaemon := NewHTTPDaemon(":8080", router, logger)
//
//	// 使用标准库
//	mux := http.NewServeMux()
//	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
//	    w.Write([]byte("pong"))
//	})
//	httpDaemon := NewHTTPDaemon(":8080", mux, logger)
func NewHTTPDaemon(addr string, ReadTimeout, WriteTimeout int, sched scheduler.Scheduler, handler http.Handler, logger daemon.Logger) *HTTPDaemon {

	// schedHandler := &SchedServeHTTP{
	// 	handler: handler,
	// 	sched:   sched,
	// 	logger:  logger,
	// }

	return &HTTPDaemon{
		// 创建 HTTP 服务器实例
		server: &http.Server{
			Addr:    addr,
			Handler: handler,
			// 可以在这里添加其他配置:
			// ReadTimeout:  15 * time.Second,
			// WriteTimeout: 15 * time.Second,
			// IdleTimeout:  60 * time.Second,

			// ReadTimeout 读取请求的最大时间
			// 包括读取请求头和请求体
			// 防止慢速客户端长时间占用连接
			ReadTimeout: time.Duration(ReadTimeout) * time.Second,

			// WriteTimeout 写入响应的最大时间
			// 从请求处理完成到写入完整响应
			// 防止慢速客户端长时间占用连接
			WriteTimeout: time.Duration(WriteTimeout) * time.Second,
		},

		logger: logger,
		// 使用 "http-server" 作为守护进程名称
		// 如果需要运行多个 HTTP 服务器,可以添加后缀区分
		// 例如: "http-server-api", "http-server-admin"
		name: "http-server",
	}
}

// Name 返回守护进程的名称
// 实现 daemon.Daemon 接口
//
// 返回:
//
//	string: 守护进程名称 "http-server"
func (d *HTTPDaemon) Name() string {
	return d.name
}

// Start 启动 HTTP 服务器
// 实现 daemon.Daemon 接口
//
// 参数:
//
//	ctx: 上下文,用于控制启动过程
//	  - 如果 ctx 被取消,应该停止启动
//	  - 传递给服务器以支持优雅关闭
//
// 返回:
//
//	error: 启动失败时返回错误,成功返回 nil
//
// 工作流程:
//  1. 记录启动日志
//  2. 在后台 goroutine 中启动服务器
//  3. 快速返回(不阻塞)
//
// 注意事项:
//   - 这个方法不会阻塞,服务器在后台运行
//   - ListenAndServe 会一直运行直到出错或被关闭
//   - 正常关闭时,ListenAndServe 返回 http.ErrServerClosed
//   - 如果是其他错误,会记录错误日志
//
// 为什么在后台 goroutine 中运行?
//   - Start 方法需要快速返回,让 Manager 继续启动其他守护进程
//   - 如果在 Start 中直接调用 ListenAndServe,会阻塞
//   - 其他守护进程将无法启动
func (d *HTTPDaemon) Start(ctx context.Context) error {
	// 记录启动日志
	// 包含监听地址,方便调试
	d.logger.Info("HTTP server starting", "addr", d.server.Addr)

	// 在后台 goroutine 中启动服务器
	// 这样 Start 方法可以快速返回
	go func() {
		// 启动 HTTP 服务器
		// ListenAndServe 会阻塞直到服务器关闭或出错
		if err := d.server.ListenAndServe(); err != nil {
			// 检查是否是正常关闭
			// http.ErrServerClosed 表示服务器被 Shutdown 方法正常关闭
			if err != http.ErrServerClosed {
				// 不是正常关闭,记录错误日志
				// 这可能是端口被占用、权限不足等问题
				d.logger.Error("HTTP server error", "error", err)
			}
			// 如果是 ErrServerClosed,不记录错误
			// 因为这是正常的关闭流程
		}
	}()

	// 快速返回
	// 服务器已经在后台启动,不需要等待
	return nil
}

// Stop 停止 HTTP 服务器
// 实现 daemon.Daemon 接口
//
// 参数:
//
//	ctx: 上下文,用于控制关闭超时
//	  - 如果 ctx 超时,Shutdown 会返回错误
//	  - 建议设置 30 秒的超时时间
//
// 返回:
//
//	error: 关闭失败或超时时返回错误
//
// 工作流程:
//  1. 记录停止日志
//  2. 调用 server.Shutdown 进行优雅关闭
//  3. 等待所有连接关闭或超时
//
// 优雅关闭流程:
//  1. 停止接收新的 HTTP 请求
//  2. 等待所有正在处理的请求完成
//  3. 关闭所有空闲连接
//  4. 返回
//
// 注意事项:
//   - Shutdown 是阻塞的,会等待所有请求处理完成
//   - 如果 ctx 超时,Shutdown 会立即返回错误
//   - 超时后可能有请求未完成,但服务器会立即停止
//   - 建议设置合理的超时时间(如 30 秒)
//
// 为什么使用 Shutdown 而不是 Close?
//   - Close 会立即关闭所有连接,可能导致请求中断
//   - Shutdown 会等待请求完成,实现优雅关闭
//   - 优雅关闭对用户体验更友好
func (d *HTTPDaemon) Stop(ctx context.Context) error {
	// 记录停止日志
	d.logger.Info("HTTP server shutting down")

	// 调用 Shutdown 进行优雅关闭
	// Shutdown 会:
	//   1. 停止监听新连接
	//   2. 等待活跃连接处理完成
	//   3. 关闭空闲连接
	// 如果 ctx 超时,Shutdown 会立即返回 ctx.Err()
	if err := d.server.Shutdown(ctx); err != nil {
		// 关闭失败(通常是超时)
		// 返回详细的错误信息
		return fmt.Errorf("HTTP server shutdown failed: %w", err)
	}

	// 关闭成功
	d.logger.Info("HTTP server stopped successfully")
	return nil
}

type SchedServeHTTP struct {
	sched   scheduler.Scheduler
	logger  daemon.Logger
	handler http.Handler
}

func (s *SchedServeHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 提交普通任务
	err := s.sched.Submit(context.Background(), func(taskCtx context.Context) {
		// 执行异步操作
		fmt.Println("Task is running")

		// 检查上下文是否被取消
		select {
		case <-taskCtx.Done():
			s.handler.ServeHTTP(w, r)
			return
		default:
			// 继续执行
		}
	})
	if err != nil {
		s.logger.Error("Submit task failed:", "error", err)
	}
}
