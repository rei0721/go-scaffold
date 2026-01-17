package app

// initMode 根据启动模式执行不同的初始化流程
func (a *App) initMode() (*App, error) {
	switch a.Options.Mode {
	case ModeServer:
		// 服务器模式
		return a.runModeServer()
	case ModeInitDB:
		// 数据库初始化模式
		return a.runModeInitDB()
	default:
		// 暂时默认模式 server
		return a.runModeServer()
	}
}
