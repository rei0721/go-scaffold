package app

func debugConfig(app *App, opts Options) {
	// 记录配置加载信息
	// 展示环境变量支持功能已生效
	app.Logger.Debug("configuration loaded successfully",
		"config_file", opts.ConfigPath,
		"env_support", "enabled")

	// 记录关键配置信息（不记录敏感信息）
	app.Logger.Debug("server configuration",
		"port", app.Config.Server.Port,
		"mode", app.Config.Server.Mode)

	app.Logger.Debug("database configuration",
		"driver", app.Config.Database.Driver,
		"host", app.Config.Database.Host,
		"db", app.Config.Database.DBName)

	if app.Config.Redis.Enabled {
		app.Logger.Debug("redis configuration",
			"enabled", true,
			"host", app.Config.Redis.Host,
			"db", app.Config.Redis.DB)
	} else {
		app.Logger.Debug("redis configuration", "enabled", false)
	}
}
