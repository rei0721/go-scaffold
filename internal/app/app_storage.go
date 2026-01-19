package app

import (
	"fmt"

	"github.com/rei0721/go-scaffold/pkg/storage"
)

// initStorage 初始化 Storage 文件服务
// 用于提供统一的文件操作API,支持文件监听、复制、Excel和图片处理
// 参数:
//
//	app: 应用实例
//
// 返回:
//
//	error: 初始化失败时的错误
func initStorage(app *App) error {
	app.Logger.Info("Initializing Storage...")

	cfg := app.Config.Storage

	// 应用默认值
	cfg.DefaultConfig()

	// 从环境变量覆盖
	cfg.OverrideConfig()

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid storage config: %w", err)
	}

	// 如果未启用,跳过初始化
	if !cfg.Enabled {
		app.Logger.Info("Storage is disabled")
		return nil
	}

	// 转换为 pkg/storage.Config
	pkgCfg := cfg.ToPkgConfig()

	// 创建 Storage 实例
	storageService, err := storage.New(pkgCfg)
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}

	// 注入依赖 - Storage 不需要 Executor 依赖
	// 如果 Storage 实现了 ExecutorInjectable 接口,可以在这里注入
	// if injector, ok := storageService.(types.ExecutorInjectable); ok {
	// 	injector.SetExecutor(app.Executor)
	// }

	app.Storage = storageService
	app.Logger.Info("Storage initialized successfully",
		"fs_type", cfg.FSType,
		"base_path", cfg.BasePath,
		"enable_watch", cfg.EnableWatch)

	return nil
}
