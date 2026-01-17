// Package config 提供配置管理和热重载支持
// 设计目标:
// - 统一配置管理:所有配置集中在一个结构体
// - 类型安全:使用结构体而不是 map
// - 验证支持:提供配置验证机制
// - 环境变量支持:可以从环境变量覆盖配置
package config

import (
	"fmt"
)

// Configurable 定义可验证配置的接口
// 所有配置结构体都应该实现这个接口
// 为什么需要接口:
// - 统一验证方式
// - 便于组合验证
// - 支持递归验证
type Configurable interface {
	// Validate 验证配置是否有效
	// 返回:
	//   error: 验证失败时的错误
	Validate() error
}

// Config 包含所有应用程序配置
// 这是顶层配置结构,聚合了所有子配置
// mapstructure tag 用于从配置文件(YAML/JSON)加载
// 设计考虑:
// - 分组管理:按功能分为 Server、Database 等
// - 清晰层次:便于理解和维护
// - 完整性:包含应用所需的所有配置
type Config struct {
	// Server HTTP 服务器配置
	// 包含端口、超时等
	Server ServerConfig `mapstructure:"server"`

	// Database 数据库连接配置
	// 支持 PostgreSQL、MySQL、SQLite
	Database DatabaseConfig `mapstructure:"database"`

	// Redis 缓存配置
	// 可选,通过 Enabled 控制是否启用
	Redis RedisConfig `mapstructure:"redis"`

	// Logger 日志配置
	// 控制日志级别、格式、输出等
	Logger LoggerConfig `mapstructure:"logger"`

	// I18n 国际化配置
	// 支持多语言
	I18n I18nConfig `mapstructure:"i18n"`

	// InitDB 数据库初始化配置
	InitDB InitDBConfig `mapstructure:"initdb"`

	// Executor 执行器配置
	// 管理异步任务的协程池
	Executor ExecutorConfig `mapstructure:"executor"`

	// JWT JWT认证配置
	// 管理token的生成和验证
	JWT JWTConfig `mapstructure:"jwt"`
}

// Validator 定义可验证配置的接口
type Validator interface {
	Validate() error
	ValidateName() string
	ValidateRequired() bool
}

// Validate 验证整个配置
// 实现 Configurable 接口
// 会递归验证所有子配置
// 返回:
//
//	error: 第一个验证失败的错误,包含错误路径
func (c *Config) Validate() error {
	validators := []Validator{
		&c.Server,
		&c.Database,
		&c.Redis,
		&c.Logger,
		&c.I18n,
		&c.InitDB,
		&c.Executor,
		&c.JWT,
	}
	for _, validator := range validators {
		if validator == nil {
			continue
		}
		if err := validator.Validate(); err != nil {
			return fmt.Errorf("%s config: %w", validator.ValidateName(), err)
		}
	}
	return nil
}

// ValidateOld 旧的验证整个配置
// 实现 Configurable 接口
// 会递归验证所有子配置
// 返回:
//
//	error: 第一个验证失败的错误,包含错误路径
func (c *Config) ValidateOld() error {
	// 验证服务器配置
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("server config: %w", err)
	}

	// 验证数据库配置
	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("database config: %w", err)
	}

	// 验证 Redis 配置
	if err := c.Redis.Validate(); err != nil {
		return fmt.Errorf("redis config: %w", err)
	}

	// 验证日志配置
	if err := c.Logger.Validate(); err != nil {
		return fmt.Errorf("logger config: %w", err)
	}

	// 验证国际化配置
	if err := c.I18n.Validate(); err != nil {
		return fmt.Errorf("i18n config: %w", err)
	}

	// 验证执行器配置
	if err := c.Executor.Validate(); err != nil {
		return fmt.Errorf("executor config: %w", err)
	}

	// 验证 JWT 配置
	if err := c.JWT.Validate(); err != nil {
		return fmt.Errorf("jwt config: %w", err)
	}

	return nil
}
