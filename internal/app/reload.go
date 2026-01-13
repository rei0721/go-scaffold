package app

import (
	"context"
	"time"

	"github.com/rei0721/rei0721/internal/config"
	"github.com/rei0721/rei0721/pkg/cache"
	"github.com/rei0721/rei0721/pkg/database"
	"github.com/rei0721/rei0721/pkg/logger"
)

// reload
func (a *App) reload(old, new *config.Config) {
	// 重新加载配置
	// a.Logger.Debug("reloading configuration...")

	// cache
	// 检查 Redis 配置是否变化
	// 使用工具方法比较新旧配置
	if isRedisConfigChanged(old, new) {
		a.Logger.Info("redis configuration changed, reloading cache...")

		// 只有在 Cache 不为 nil 且新配置启用了 Redis 时才重载
		if a.Cache != nil && new.Redis.Enabled {
			// 创建新的缓存配置
			newCacheCfg := &cache.Config{
				Host:         new.Redis.Host,
				Port:         new.Redis.Port,
				Password:     new.Redis.Password,
				DB:           new.Redis.DB,
				PoolSize:     new.Redis.PoolSize,
				MinIdleConns: new.Redis.MinIdleConns,
				MaxRetries:   new.Redis.MaxRetries,
				DialTimeout:  time.Duration(new.Redis.DialTimeout) * time.Second,
				ReadTimeout:  time.Duration(new.Redis.ReadTimeout) * time.Second,
				WriteTimeout: time.Duration(new.Redis.WriteTimeout) * time.Second,
			}

			// 使用超时上下文进行重载
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// 原子化重载缓存配置
			err := a.Cache.Reload(ctx, newCacheCfg)
			if err != nil {
				a.Logger.Error("failed to reload redis cache", "error", err)
			} else {
				a.Logger.Info("redis cache reloaded successfully")
			}
		} else if !new.Redis.Enabled {
			a.Logger.Info("redis disabled in new config")
		} else {
			a.Logger.Warn("cache is nil, cannot reload redis configuration")
		}
	}

	if isDatabaseConfigChanged(old, new) {
		a.Logger.Info("database configuration changed, reloading database...")

		// 重新加载数据库配置
		newDBCfg := &database.Config{
			Driver:       database.Driver(new.Database.Driver),
			Host:         new.Database.Host,
			Port:         new.Database.Port,
			User:         new.Database.User,
			Password:     new.Database.Password,
			DBName:       new.Database.DBName,
			MaxOpenConns: new.Database.MaxOpenConns,
			MaxIdleConns: new.Database.MaxIdleConns,
		}

		if err := a.DB.Reload(newDBCfg); err != nil {
			a.Logger.Error("failed to reload database", "error", err)
		} else {
			a.Logger.Info("database reloaded successfully")
		}
	}

	// logger
	// 检查日志配置是否变化
	if isLoggerConfigChanged(old, new) {
		a.Logger.Info("logger configuration changed, reloading logger...")

		// 创建新的日志配置
		newLoggerCfg := &logger.Config{
			Level:         new.Logger.Level,
			Format:        new.Logger.Format,
			ConsoleFormat: new.Logger.ConsoleFormat,
			FileFormat:    new.Logger.FileFormat,
			Output:        new.Logger.Output,
			FilePath:      new.Logger.FilePath,
			MaxSize:       new.Logger.MaxSize,
			MaxBackups:    new.Logger.MaxBackups,
			MaxAge:        new.Logger.MaxAge,
		}

		// 原子化重载日志配置
		if err := a.Logger.Reload(newLoggerCfg); err != nil {
			a.Logger.Error("failed to reload logger", "error", err)
		} else {
			a.Logger.Info("logger reloaded successfully")
		}
	}

	// executor
	// 检查执行器配置是否变化
	if isExecutorConfigChanged(old, new) {
		a.Logger.Info("executor configuration changed, reloading executor...")

		// 只有在 Executor 不为 nil 且新配置启用了执行器时才重载
		if a.Executor != nil && new.Executor.Enabled {
			// 转换配置格式
			newExecutorConfigs := makeExecutorConfigs(new)

			// 原子化重载执行器配置
			if err := a.Executor.Reload(newExecutorConfigs); err != nil {
				a.Logger.Error("failed to reload executor", "error", err)
			} else {
				a.Logger.Info("executor reloaded successfully", "pools", len(newExecutorConfigs))
			}
		} else if !new.Executor.Enabled {
			a.Logger.Info("executor disabled in new config")
		} else {
			a.Logger.Warn("executor is nil, cannot reload configuration")
		}
	}
}
