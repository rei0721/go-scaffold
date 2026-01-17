package config

import (
	"errors"
	"os"
	"strconv"
)

// DatabaseConfig 数据库连接配置
// 包含连接数据库所需的所有信息
type DatabaseConfig struct {
	// Driver 数据库驱动类型
	// 可选值: postgres, mysql, sqlite
	// 影响连接字符串格式和 SQL 方言
	Driver string `mapstructure:"driver"`

	// Host 数据库服务器地址
	// 例如: localhost, 127.0.0.1, db.example.com
	// SQLite 不需要此字段
	Host string `mapstructure:"host"`

	// Port 数据库端口
	// PostgreSQL 默认: 5432
	// MySQL 默认: 3306
	// SQLite 不需要此字段
	Port int `mapstructure:"port"`

	// User 数据库用户名
	// SQLite 不需要此字段
	User string `mapstructure:"user"`

	// Password 数据库密码
	// 生产环境应该从环境变量或密钥管理服务读取
	// 不要硬编码在配置文件中
	Password string `mapstructure:"password"`

	// DBName 数据库名称
	// PostgreSQL/MySQL: 数据库名
	// SQLite: 文件路径
	DBName string `mapstructure:"dbname"`

	// MaxOpenConns 最大打开连接数
	// 0 表示无限制(不推荐)
	// 推荐: 10-100,根据并发量调整
	MaxOpenConns int `mapstructure:"max_open_conns"`

	// MaxIdleConns 最大空闲连接数
	// 建议设置为 MaxOpenConns 的 50%-100%
	// 保持空闲连接可以提高响应速度
	MaxIdleConns int `mapstructure:"max_idle_conns"`
}

func (c *DatabaseConfig) ValidateName() string {
	return AppDatabaseName
}

func (c *DatabaseConfig) ValidateRequired() bool {
	return true
}

// Validate 验证数据库配置
// 实现 Configurable 接口
func (c *DatabaseConfig) Validate() error {
	// 验证驱动类型
	validDrivers := map[string]bool{"postgres": true, "mysql": true, "sqlite": true}
	if !validDrivers[c.Driver] {
		return errors.New("driver must be postgres, mysql, or sqlite")
	}

	// SQLite 以外的数据库需要网络配置
	if c.Driver != "sqlite" {
		// 验证主机地址
		if c.Host == "" {
			return errors.New("host is required")
		}

		// 验证端口
		if c.Port <= 0 || c.Port > 65535 {
			return errors.New("port must be between 1 and 65535")
		}

		// 验证用户名
		if c.User == "" {
			return errors.New("user is required")
		}

		// 注意:密码可以为空(虽然不推荐)
	}

	// 验证数据库名称
	if c.DBName == "" {
		// 所有数据库都需要数据库名
		// SQLite: 文件路径
		// PostgreSQL/MySQL: 数据库名
		return errors.New("dbname is required")
	}

	// 验证连接池参数
	if c.MaxOpenConns < 0 {
		// 必须 >= 0
		// 0 表示无限制(不推荐)
		return errors.New("maxOpenConns must be non-negative")
	}

	if c.MaxIdleConns < 0 {
		// 必须 >= 0
		return errors.New("maxIdleConns must be non-negative")
	}

	return nil
}

// overrideDatabaseConfig 使用环境变量覆盖数据库配置
func overrideDatabaseConfig(cfg *DatabaseConfig) {
	// Driver
	if val := os.Getenv(EnvPrefixJoin(EnvDBDriver)); val != "" {
		cfg.Driver = val
	}

	// Host
	if val := os.Getenv(EnvPrefixJoin(EnvDBHost)); val != "" {
		cfg.Host = val
	}

	// Port
	if val := os.Getenv(EnvPrefixJoin(EnvDBPort)); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			cfg.Port = port
		}
	}

	// User
	if val := os.Getenv(EnvPrefixJoin(EnvDBUser)); val != "" {
		cfg.User = val
	}

	// Password
	// 密码应该优先使用环境变量
	if val := os.Getenv(EnvPrefixJoin(EnvDBPassword)); val != "" {
		cfg.Password = val
	}

	// DBName
	if val := os.Getenv(EnvPrefixJoin(EnvDBName)); val != "" {
		cfg.DBName = val
	}

	// MaxOpenConns
	if val := os.Getenv(EnvPrefixJoin(EnvDBMaxOpenConns)); val != "" {
		if conns, err := strconv.Atoi(val); err == nil {
			cfg.MaxOpenConns = conns
		}
	}

	// MaxIdleConns
	if val := os.Getenv(EnvPrefixJoin(EnvDBMaxIdleConns)); val != "" {
		if conns, err := strconv.Atoi(val); err == nil {
			cfg.MaxIdleConns = conns
		}
	}
}
