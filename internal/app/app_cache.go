package app

import (
	"time"

	"github.com/rei0721/go-scaffold/pkg/cache"
)

// 初始化 Redis 缓存(可选)
func initCache(app *App) error {
	// 如果配置中启用了 Redis,则创建缓存实例
	if app.Config.Redis.Enabled {
		cacheCfg := &cache.Config{
			Host:         app.Config.Redis.Host,
			Port:         app.Config.Redis.Port,
			Password:     app.Config.Redis.Password,
			DB:           app.Config.Redis.DB,
			PoolSize:     app.Config.Redis.PoolSize,
			MinIdleConns: app.Config.Redis.MinIdleConns,
			MaxRetries:   app.Config.Redis.MaxRetries,
			DialTimeout:  time.Duration(app.Config.Redis.DialTimeout) * time.Second,
			ReadTimeout:  time.Duration(app.Config.Redis.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(app.Config.Redis.WriteTimeout) * time.Second,
		}

		cacheClient, err := cache.NewRedis(cacheCfg, app.Logger)
		if err != nil {
			// Redis 连接失败
			// 可以选择:
			// 1. 返回错误,强制要求 Redis 可用
			// 2. 警告但继续,允许无缓存运行
			// 这里选择方案 2,提高可用性
			app.Logger.Warn("failed to connect to redis, running without cache", "error", err)
			app.Cache = nil
		} else {
			app.Cache = cacheClient
			app.Logger.Info("redis cache connected successfully")
		}
	} else {
		app.Logger.Info("redis cache disabled")
		app.Cache = nil
	}

	return nil
}
