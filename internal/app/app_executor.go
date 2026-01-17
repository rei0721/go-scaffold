package app

import (
	"fmt"
	"time"

	"github.com/rei0721/go-scaffold/pkg/executor"
)

// initExecutor 初始化执行器管理器
// 从配置创建执行器,管理多个协程池
// 参数:
//
//	app: 应用实例
//
// 返回:
//
//	error: 初始化失败时的错误
func (app *App) initExecutor() error {
	// 检查是否启用执行器
	if !app.Config.Executor.Enabled {
		app.Logger.Info("executor is disabled, skipping initialization")
		return nil
	}

	// 转换配置格式
	// internal/config.ExecutorPoolConfig -> pkg/executor.Config
	configs := make([]executor.Config, 0, len(app.Config.Executor.Pools))
	for _, poolCfg := range app.Config.Executor.Pools {
		configs = append(configs, executor.Config{
			Name:        executor.PoolName(poolCfg.Name),
			Size:        poolCfg.Size,
			Expiry:      time.Duration(poolCfg.Expiry) * time.Second,
			NonBlocking: poolCfg.NonBlocking,
		})
	}

	// 创建执行器管理器
	mgr, err := executor.NewManager(configs)
	if err != nil {
		return fmt.Errorf("failed to create executor manager: %w", err)
	}

	app.Executor = mgr
	app.Logger.Info("executor initialized", "pools", len(configs))

	return nil
}
