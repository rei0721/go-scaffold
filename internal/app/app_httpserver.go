package app

import (
	"fmt"
	"time"

	"github.com/rei0721/go-scaffold/pkg/httpserver"
)

// initHTTPServer 初始化 HTTP 服务器
// 使用 pkg/httpserver 封装，替代原来的直接创建 http.Server
// 这个函数应该在 Router 初始化之后调用
func (app *App) initHTTPServer() error {
	// 创建 HTTP 服务器配置
	cfg := &httpserver.Config{
		Host:         app.Config.Server.Host,
		Port:         app.Config.Server.Port,
		ReadTimeout:  time.Duration(app.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(app.Config.Server.WriteTimeout) * time.Second,
	}

	// 创建 HTTP 服务器实例（不直接注入executor）
	server, err := httpserver.New(app.Router, cfg, app.Logger)
	if err != nil {
		return fmt.Errorf("failed to create http server: %w", err)
	}

	app.HTTPServer = server

	// ⭐ 延迟注入executor到HTTPServer
	if app.Executor != nil {
		server.SetExecutor(app.Executor)
		app.Logger.Debug("executor injected into httpserver")
	}

	return nil
}
