package app

import (
	"fmt"

	"github.com/rei0721/go-scaffold/internal/config"
)

func initConfig(app *App, opts Options) error {
	// 1. 初始化配置管理器并加载配置
	// 配置是整个应用的基础,必须最先加载
	configManager := config.NewManager()
	if err := configManager.Load(opts.ConfigPath); err != nil {
		// 配置加载失败,应用无法启动
		return fmt.Errorf("failed to load config: %w", err)
	}
	app.ConfigManager = configManager
	app.Config = configManager.Get()
	return nil
}
