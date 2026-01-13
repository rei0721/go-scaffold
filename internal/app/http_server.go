package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rei0721/rei0721/pkg/utils"
)

func listenHttpServer(app *App) error {
	if app.Config.Server.Port == 0 {
		port, err := utils.GetAvailablePort(9000, 30000)
		if err != nil {
			return err
		}
		app.Config.Server.Port = port
	}
	addr := fmt.Sprintf("%s:%d", app.Config.Server.Host, app.Config.Server.Port)
	// 校验地址是否合法，如果合法则使用，否则使用默认地址
	addrErr := utils.IsValidHTTPListenAddr(addr)
	if addrErr != nil {
		addr = fmt.Sprintf("%s:%d", ConstantsDefaultHost, app.Config.Server.Port)
	}

	// 创建 HTTP 服务器实例
	// 配置超时参数防止慢速客户端占用连接
	app.HTTPServer = &http.Server{
		// Addr 监听地址
		Addr: addr,

		// Handler HTTP 请求处理器(Gin Router)
		Handler: app.Router,
		// 可以在这里添加其他配置:
		// ReadTimeout:  15 * time.Second,
		// WriteTimeout: 15 * time.Second,
		// IdleTimeout:  60 * time.Second,

		// ReadTimeout 读取请求的最大时间
		// 包括读取请求头和请求体
		// 防止慢速客户端长时间占用连接
		ReadTimeout: time.Duration(app.Config.Server.ReadTimeout) * time.Second,

		// WriteTimeout 写入响应的最大时间
		// 从请求处理完成到写入完整响应
		// 防止慢速客户端长时间占用连接
		WriteTimeout: time.Duration(app.Config.Server.WriteTimeout) * time.Second,

		// IdleTimeout 空闲连接的超时时间
		// 从连接建立到空闲的最大时间
		// 防止慢速客户端长时间占用连接
		IdleTimeout: time.Duration(app.Config.Server.IdleTimeout) * time.Second,
	}

	app.Logger.Info(fmt.Sprintf("starting HTTP server on http://%s", addr), "addr", addr)

	// 启动服务器并开始监听
	// ListenAndServe 会阻塞,直到:
	// - 发生错误
	// - Shutdown() 被调用
	if err := app.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		// ErrServerClosed 是正常的关闭,不是错误
		// 其他错误才需要返回
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}
