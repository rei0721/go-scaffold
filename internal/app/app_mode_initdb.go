package app

// runModeInitDB initdb 模式
func (a *App) runModeInitDB() (*App, error) {
	// initdb 模式：仅初始化到数据库，然后执行初始化
	initApps := []func(app *App) error{
		initSqlGenerator, // 初始化 SQL 生成器
		initI18n,         // 初始化 i18n（用于日志消息）
		initDatabase,     // 初始化数据库连接
	}

	for _, initApp := range initApps {
		if err := initApp(a); err != nil {
			return nil, err
		}
	}

	// 初始化 Executor
	if err := initExecutor(a); err != nil {
		return nil, err
	}

	// ⭐ Executor初始化完成后，注入到Logger
	if a.Executor != nil && a.Logger != nil {
		a.Logger.SetExecutor(a.Executor)
		a.Logger.Debug(a.UI18n("internal.app.logger_debug_executor_injected"))
	}

	// 初始化业务逻辑(包括 Router)
	if err := initBusiness(a); err != nil {
		return nil, err
	}

	// 执行数据库初始化
	if err := runInitDB(a); err != nil {
		return nil, err
	}

	a.Logger.Info(a.UI18n("internal.app.logger_info_initdb_mode_completed"))
	return a, nil
}
