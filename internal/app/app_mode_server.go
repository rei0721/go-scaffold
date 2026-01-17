package app

import "github.com/rei0721/go-scaffold/internal/config"

func (app *App) runModeServer() (*App, error) {
	// server 模式：完整初始化流程
	// 阶段1：核心基础设施
	if err := initCache(app); err != nil {
		return nil, err
	}
	if err := initDatabase(app); err != nil {
		return nil, err
	}

	// 阶段2：初始化Executor
	if err := initExecutor(app); err != nil {
		return nil, err
	}

	// ⭐ Executor初始化完成后，立即注入到Logger
	if app.Executor != nil && app.Logger != nil {
		app.Logger.SetExecutor(app.Executor)
		app.Logger.Debug("executor injected into logger")
	}

	// 阶段2.5：初始化JWT认证
	if err := initJWT(app); err != nil {
		return nil, err
	}

	// 阶段3：业务层和HTTP服务器
	// 注意：initBusiness和initHTTPServer内部会自动注入executor
	if err := initBusiness(app); err != nil {
		return nil, err
	}
	if err := initHTTPServer(app); err != nil {
		return nil, err
	}

	// Start config file watching for hot-reload
	if err := app.ConfigManager.Watch(); err != nil {
		app.Logger.Warn("failed to start config watcher", "error", err)
	}
	app.Logger.Debug("config watcher started")

	// Register config change hook
	// 当配置文件变化时自动调用
	app.ConfigManager.RegisterHook(func(old, new *config.Config) {
		app.Logger.Info("configuration file changed, processing updates...")

		// 重载 app
		app.reload(old, new)

		// 更新应用配置引用
		app.Config = new
		app.Logger.Info("configuration update completed")
	})

	app.Logger.Info("application initialized successfully")
	return app, nil
}
