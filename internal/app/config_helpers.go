package app

import (
	"time"

	"github.com/rei0721/go-scaffold/internal/config"
	"github.com/rei0721/go-scaffold/pkg/executor"
)

// isRedisConfigChanged 检查 Redis 配置是否发生变化
// 比较新旧配置的所有 Redis 相关字段
// 参数:
//
//	oldCfg: 旧配置
//	newCfg: 新配置
//
// 返回:
//
//	bool: 如果配置有任何差异返回 true,否则返回 false
//
// 使用示例:
//
//	if isRedisConfigChanged(oldConfig, newConfig) {
//	    // 配置变化了，需要重载
//	    app.Cache.Reload(ctx, newCacheCfg)
//	}
func isRedisConfigChanged(oldCfg, newCfg *config.Config) bool {
	// 快速检查：如果两个配置指针相同，肯定没变化
	if oldCfg == newCfg {
		return false
	}

	// 比较 Enabled 状态
	// 如果启用状态改变，肯定需要重载
	if oldCfg.Redis.Enabled != newCfg.Redis.Enabled {
		return true
	}

	// 如果都未启用，不需要关心其他字段
	if !newCfg.Redis.Enabled {
		return false
	}

	// 比较所有 Redis 配置字段
	// 任何一个字段不同都认为配置改变了

	// 连接相关
	if oldCfg.Redis.Host != newCfg.Redis.Host {
		return true
	}
	if oldCfg.Redis.Port != newCfg.Redis.Port {
		return true
	}
	if oldCfg.Redis.Password != newCfg.Redis.Password {
		return true
	}
	if oldCfg.Redis.DB != newCfg.Redis.DB {
		return true
	}

	// 连接池相关
	if oldCfg.Redis.PoolSize != newCfg.Redis.PoolSize {
		return true
	}
	if oldCfg.Redis.MinIdleConns != newCfg.Redis.MinIdleConns {
		return true
	}

	// 重试和超时相关
	if oldCfg.Redis.MaxRetries != newCfg.Redis.MaxRetries {
		return true
	}
	if oldCfg.Redis.DialTimeout != newCfg.Redis.DialTimeout {
		return true
	}
	if oldCfg.Redis.ReadTimeout != newCfg.Redis.ReadTimeout {
		return true
	}
	if oldCfg.Redis.WriteTimeout != newCfg.Redis.WriteTimeout {
		return true
	}

	// 所有字段都相同
	return false
}

// isDatabaseConfigChanged 检查数据库配置是否发生变化
// 比较新旧配置的所有数据库相关字段
// 参数:
//
//	oldCfg: 旧配置
//	newCfg: 新配置
//
// 返回:
//
//	bool: 如果配置有任何差异返回 true,否则返回 false
func isDatabaseConfigChanged(oldCfg, newCfg *config.Config) bool {
	if oldCfg == newCfg {
		return false
	}

	// 比较数据库配置字段
	if oldCfg.Database.Driver != newCfg.Database.Driver {
		return true
	}
	if oldCfg.Database.Host != newCfg.Database.Host {
		return true
	}
	if oldCfg.Database.Port != newCfg.Database.Port {
		return true
	}
	if oldCfg.Database.User != newCfg.Database.User {
		return true
	}
	if oldCfg.Database.Password != newCfg.Database.Password {
		return true
	}
	if oldCfg.Database.DBName != newCfg.Database.DBName {
		return true
	}
	if oldCfg.Database.MaxOpenConns != newCfg.Database.MaxOpenConns {
		return true
	}
	if oldCfg.Database.MaxIdleConns != newCfg.Database.MaxIdleConns {
		return true
	}

	return false
}

// isServerConfigChanged 检查服务器配置是否发生变化
// 参数:
//
//	oldCfg: 旧配置
//	newCfg: 新配置
//
// 返回:
//
//	bool: 如果配置有任何差异返回 true,否则返回 false
func isServerConfigChanged(oldCfg, newCfg *config.Config) bool {
	if oldCfg == newCfg {
		return false
	}

	if oldCfg.Server.Port != newCfg.Server.Port {
		return true
	}
	if oldCfg.Server.Mode != newCfg.Server.Mode {
		return true
	}
	if oldCfg.Server.ReadTimeout != newCfg.Server.ReadTimeout {
		return true
	}
	if oldCfg.Server.WriteTimeout != newCfg.Server.WriteTimeout {
		return true
	}

	return false
}

// isLoggerConfigChanged 检查日志配置是否发生变化
// 参数:
//
//	oldCfg: 旧配置
//	newCfg: 新配置
//
// 返回:
//
//	bool: 如果配置有任何差异返回 true,否则返回 false
//
// 使用示例:
//
//	if isLoggerConfigChanged(oldConfig, newConfig) {
//	    // 配置变化了，需要重载
//	    logger.Reload(newLoggerCfg)
//	}
func isLoggerConfigChanged(oldCfg, newCfg *config.Config) bool {
	if oldCfg == newCfg {
		return false
	}

	// 比较日志级别
	if oldCfg.Logger.Level != newCfg.Logger.Level {
		return true
	}
	// 比较日志格式
	if oldCfg.Logger.Format != newCfg.Logger.Format {
		return true
	}
	// 比较控制台格式
	if oldCfg.Logger.ConsoleFormat != newCfg.Logger.ConsoleFormat {
		return true
	}
	// 比较文件格式
	if oldCfg.Logger.FileFormat != newCfg.Logger.FileFormat {
		return true
	}
	// 比较输出目标
	if oldCfg.Logger.Output != newCfg.Logger.Output {
		return true
	}

	// 比较日志文件路径
	if oldCfg.Logger.FilePath != newCfg.Logger.FilePath {
		return true
	}

	// 比较日志文件最大大小
	if oldCfg.Logger.MaxSize != newCfg.Logger.MaxSize {
		return true
	}

	// 比较日志文件最大备份数
	if oldCfg.Logger.MaxBackups != newCfg.Logger.MaxBackups {
		return true
	}

	// 比较日志文件最大年龄
	if oldCfg.Logger.MaxAge != newCfg.Logger.MaxAge {
		return true
	}

	return false
}

// isExecutorConfigChanged 检查执行器配置是否发生变化
// 参数:
//
//	oldCfg: 旧配置
//	newCfg: 新配置
//
// 返回:
//
//	bool: 如果配置有任何差异返回 true,否则返回 false
func isExecutorConfigChanged(oldCfg, newCfg *config.Config) bool {
	if oldCfg == newCfg {
		return false
	}

	// 比较启用状态
	if oldCfg.Executor.Enabled != newCfg.Executor.Enabled {
		return true
	}

	// 如果都未启用,不需要关心其他字段
	if !newCfg.Executor.Enabled {
		return false
	}

	// 比较池数量
	if len(oldCfg.Executor.Pools) != len(newCfg.Executor.Pools) {
		return true
	}

	// 比较每个池的配置
	// 使用 map 快速查找
	oldPools := make(map[string]config.ExecutorPoolConfig)
	for _, pool := range oldCfg.Executor.Pools {
		oldPools[pool.Name] = pool
	}

	for _, newPool := range newCfg.Executor.Pools {
		oldPool, exists := oldPools[newPool.Name]
		if !exists {
			// 新池,配置改变了
			return true
		}

		// 比较池配置
		if oldPool.Size != newPool.Size {
			return true
		}
		if oldPool.Expiry != newPool.Expiry {
			return true
		}
		if oldPool.NonBlocking != newPool.NonBlocking {
			return true
		}
	}

	return false
}

// makeExecutorConfigs 从应用配置创建执行器配置
// 转换 internal/config.ExecutorPoolConfig 到 pkg/executor.Config
// 参数:
//
//	cfg: 应用配置
//
// 返回:
//
//	[]executor.Config: 执行器配置列表
func makeExecutorConfigs(cfg *config.Config) []executor.Config {
	configs := make([]executor.Config, 0, len(cfg.Executor.Pools))
	for _, poolCfg := range cfg.Executor.Pools {
		configs = append(configs, executor.Config{
			Name:        executor.PoolName(poolCfg.Name),
			Size:        poolCfg.Size,
			Expiry:      time.Duration(poolCfg.Expiry) * time.Second,
			NonBlocking: poolCfg.NonBlocking,
		})
	}
	return configs
}

// isStorageConfigChanged 检查 Storage 配置是否发生变化
// 比较新旧配置的所有 Storage 相关字段
// 参数:
//
//	oldCfg: 旧配置
//	newCfg: 新配置
//
// 返回:
//
//	bool: 如果配置有任何差异返回 true,否则返回 false
func isStorageConfigChanged(oldCfg, newCfg *config.Config) bool {
	if oldCfg == newCfg {
		return false
	}

	// 比较启用状态
	if oldCfg.Storage.Enabled != newCfg.Storage.Enabled {
		return true
	}

	// 如果都未启用,不需要关心其他字段
	if !newCfg.Storage.Enabled {
		return false
	}

	// 比较文件系统类型
	if oldCfg.Storage.FSType != newCfg.Storage.FSType {
		return true
	}

	// 比较基础路径
	if oldCfg.Storage.BasePath != newCfg.Storage.BasePath {
		return true
	}

	// 比较是否启用监听
	if oldCfg.Storage.EnableWatch != newCfg.Storage.EnableWatch {
		return true
	}

	// 比较监听缓冲区大小
	if oldCfg.Storage.WatchBufferSize != newCfg.Storage.WatchBufferSize {
		return true
	}

	return false
}
