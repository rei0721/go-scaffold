package config

import (
	"errors"
	"os"
	"strconv"
)

// RedisConfig Redis 连接配置
// Redis 用于缓存、会话存储等
type RedisConfig struct {
	// Enabled 是否启用 Redis
	// false 时,应用不会连接 Redis
	// 可以在开发环境中禁用
	Enabled bool `mapstructure:"enabled"`

	// Host Redis 服务器地址
	// 例如: localhost, 127.0.0.1, redis.example.com
	Host string `mapstructure:"host"`

	// Port Redis 端口
	// 默认: 6379
	Port int `mapstructure:"port"`

	// Password Redis 密码
	// 如果 Redis 未设置密码,留空
	Password string `mapstructure:"password"`

	// DB Redis 数据库编号
	// Redis 支持 0-15 共 16 个数据库
	// 默认: 0
	// 可以用不同的 DB 隔离不同环境的数据
	DB int `mapstructure:"db"`

	// PoolSize 连接池大小
	// 0 表示使用默认值(通常是 CPU 核心数 * 10)
	// 推荐: 10-100
	PoolSize int `mapstructure:"pool_size"`

	// MinIdleConns 最小空闲连接数
	// 保持一定数量的空闲连接可以提高响应速度
	// 推荐: PoolSize 的 30-50%
	MinIdleConns int `mapstructure:"min_idle_conns"`

	// MaxRetries 最大重试次数
	// 当命令执行失败时自动重试的次数
	// 0 表示不重试
	// 推荐: 2-3 次
	MaxRetries int `mapstructure:"max_retries"`

	// DialTimeout 连接超时时间(秒)
	// 建立 TCP 连接的最大等待时间
	// 推荐: 5 秒
	DialTimeout int `mapstructure:"dial_timeout"`

	// ReadTimeout 读取超时时间(秒)
	// 从 Redis 读取响应的最大等待时间
	// 推荐: 3 秒
	ReadTimeout int `mapstructure:"read_timeout"`

	// WriteTimeout 写入超时时间(秒)
	// 向 Redis 写入命令的最大等待时间
	// 推荐: 3 秒
	WriteTimeout int `mapstructure:"write_timeout"`
}

func (c *RedisConfig) ValidateName() string {
	return AppRedisName
}

func (c *RedisConfig) ValidateRequired() bool {
	return false
}

// Validate 验证 Redis 配置
// 实现 Configurable 接口
func (c *RedisConfig) Validate() error {
	// 如果未启用,跳过验证
	if !c.Enabled {
		return nil
	}

	// 启用时必须提供配置
	if c.Host == "" {
		return errors.New("host is required when redis is enabled")
	}

	// 验证端口
	if c.Port <= 0 || c.Port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}

	// 验证数据库编号
	if c.DB < 0 || c.DB > 15 {
		// Redis 只支持 0-15
		return errors.New("db must be between 0 and 15")
	}

	// 验证连接池大小
	if c.PoolSize < 0 {
		// 必须 >= 0
		// 0 表示使用默认值
		return errors.New("poolSize must be non-negative")
	}

	return nil
}

// overrideRedisConfig 使用环境变量覆盖 Redis 配置
func overrideRedisConfig(cfg *RedisConfig) {
	// Enabled
	if val := os.Getenv(EnvRedisEnabled); val != "" {
		if enabled, err := strconv.ParseBool(val); err == nil {
			cfg.Enabled = enabled
		}
	}

	// Host
	if val := os.Getenv(EnvRedisHost); val != "" {
		cfg.Host = val
	}

	// Port
	if val := os.Getenv(EnvRedisPort); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			cfg.Port = port
		}
	}

	// Password
	// 密码应该优先使用环境变量
	if val := os.Getenv(EnvRedisPassword); val != "" {
		cfg.Password = val
	}

	// DB
	if val := os.Getenv(EnvRedisDB); val != "" {
		if db, err := strconv.Atoi(val); err == nil {
			cfg.DB = db
		}
	}

	// PoolSize
	if val := os.Getenv(EnvRedisPoolSize); val != "" {
		if size, err := strconv.Atoi(val); err == nil {
			cfg.PoolSize = size
		}
	}

	// MinIdleConns
	if val := os.Getenv(EnvRedisMinIdleConns); val != "" {
		if conns, err := strconv.Atoi(val); err == nil {
			cfg.MinIdleConns = conns
		}
	}

	// MaxRetries
	if val := os.Getenv(EnvRedisMaxRetries); val != "" {
		if retries, err := strconv.Atoi(val); err == nil {
			cfg.MaxRetries = retries
		}
	}

	// DialTimeout
	if val := os.Getenv(EnvRedisDialTimeout); val != "" {
		if timeout, err := strconv.Atoi(val); err == nil {
			cfg.DialTimeout = timeout
		}
	}

	// ReadTimeout
	if val := os.Getenv(EnvRedisReadTimeout); val != "" {
		if timeout, err := strconv.Atoi(val); err == nil {
			cfg.ReadTimeout = timeout
		}
	}

	// WriteTimeout
	if val := os.Getenv(EnvRedisWriteTimeout); val != "" {
		if timeout, err := strconv.Atoi(val); err == nil {
			cfg.WriteTimeout = timeout
		}
	}
}
