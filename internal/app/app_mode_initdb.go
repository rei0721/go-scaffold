package app

// runModeInitDB initdb 模式
func (app *App) runModeInitDB() (*App, error) {
	// initdb 模式：仅初始化到数据库，然后执行初始化

	// 初始化 SQL 生成器
	if err := app.initSqlGenerator(); err != nil {
		return nil, err
	}

	// 初始化 i18n
	if err := app.initI18n(); err != nil {
		return nil, err
	}

	// 初始化数据库连接
	if err := app.initDatabase(); err != nil {
		return nil, err
	}

	// 初始化 Executor
	if err := app.initExecutor(); err != nil {
		return nil, err
	}

	// ⭐ Executor初始化完成后，注入到Logger
	if app.Executor != nil && app.Logger != nil {
		app.Logger.SetExecutor(app.Executor)
		app.Logger.Debug(app.UI18n("internal.app.logger_debug_executor_injected"))
	}

	// 初始化业务逻辑(包括 Router)
	if err := app.initBusiness(); err != nil {
		return nil, err
	}

	// 执行数据库初始化
	if err := runInitDB(app); err != nil {
		return nil, err
	}

	app.Logger.Info(app.UI18n("internal.app.logger_info_initdb_mode_completed"))
	return app, nil
}
