// Package config 提供配置管理和热重载支持
// 设计目标:
// - 统一配置管理:所有配置集中在一个结构体
// - 类型安全:使用结构体而不是 map
// - 验证支持:提供配置验证机制
// - 环境变量支持:可以从环境变量覆盖配置
package config

import (
	"errors"
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

	// Executor 执行器配置
	// 管理异步任务的协程池
	Executor ExecutorConfig `mapstructure:"executor"`

	// JWT JWT认证配置
	// 管理token的生成和验证
	JWT JWTConfig `mapstructure:"jwt"`
}

// ServerConfig HTTP 服务器配置
// 控制 HTTP 服务的行为
type ServerConfig struct {
	// Host HTTP服务器地址
	// 例如: localhost, 127.0.0.1, db.example.com
	// SQLite 不需要此字段
	Host string `mapstructure:"host"`

	// Port 监听端口
	// 有效范围: 1-65535
	// 常用端口: 8080, 3000, 80(需要 root)
	Port int `mapstructure:"port"`

	// Mode 运行模式
	// 可选值:
	// - debug: 开发模式,详细日志,panic 堆栈
	// - release: 生产模式,性能优化,简化日志
	// - test: 测试模式
	// 影响:
	// - Gin 的日志详细程度
	// - 性能优化级别
	// - panic 恢复行为
	Mode string `mapstructure:"mode"`

	// ReadTimeout 读取请求的超时时间(秒)
	// 从连接建立到读取完整请求体的最大时间
	// 防止慢速客户端占用连接
	// 推荐: 5-60 秒
	ReadTimeout int `mapstructure:"readTimeout"`

	// WriteTimeout 写入响应的超时时间(秒)
	// 从请求处理完成到写入完整响应的最大时间
	// 防止慢速客户端占用连接
	// 推荐: 10-120 秒(取决于响应大小)
	WriteTimeout int `mapstructure:"writeTimeout"`

	// IdleTimeout 空闲连接的超时时间(秒)
	// 从连接建立到空闲的最大时间
	// 防止慢速客户端占用连接
	// 推荐: 60-300 秒
	IdleTimeout int `mapstructure:"idleTimeout"`
}

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
	MaxOpenConns int `mapstructure:"maxOpenConns"`

	// MaxIdleConns 最大空闲连接数
	// 建议设置为 MaxOpenConns 的 50%-100%
	// 保持空闲连接可以提高响应速度
	MaxIdleConns int `mapstructure:"maxIdleConns"`
}

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
	PoolSize int `mapstructure:"poolSize"`

	// MinIdleConns 最小空闲连接数
	// 保持一定数量的空闲连接可以提高响应速度
	// 推荐: PoolSize 的 30-50%
	MinIdleConns int `mapstructure:"minIdleConns"`

	// MaxRetries 最大重试次数
	// 当命令执行失败时自动重试的次数
	// 0 表示不重试
	// 推荐: 2-3 次
	MaxRetries int `mapstructure:"maxRetries"`

	// DialTimeout 连接超时时间(秒)
	// 建立 TCP 连接的最大等待时间
	// 推荐: 5 秒
	DialTimeout int `mapstructure:"dialTimeout"`

	// ReadTimeout 读取超时时间(秒)
	// 从 Redis 读取响应的最大等待时间
	// 推荐: 3 秒
	ReadTimeout int `mapstructure:"readTimeout"`

	// WriteTimeout 写入超时时间(秒)
	// 向 Redis 写入命令的最大等待时间
	// 推荐: 3 秒
	WriteTimeout int `mapstructure:"writeTimeout"`
}

// Config 保存日志配置
// 包含日志库初始化所需的所有参数
// 这些配置通常从配置文件加载
type LoggerConfig struct {
	// Level 最低日志级别
	// 可选值: debug, info, warn, error
	// 只有 >= 此级别的日志会被记录
	// 例如:如果设置为 info,debug 日志不会输出
	// 开发环境推荐: debug
	// 生产环境推荐: info 或 warn
	Level string `mapstructure:"level"`

	// Format 默认输出格式(用于所有输出)
	// 可选值:
	// - json: 结构化 JSON 格式,便于日志系统解析
	// - console: 人类可读的控制台格式,便于开发调试
	// 如果设置了 ConsoleFormat 或 FileFormat,则此字段作为后备默认值
	// 生产环境推荐: json(便于 ELK、Splunk 等系统分析)
	// 开发环境推荐: console(易读)
	Format string `mapstructure:"format"`

	// ConsoleFormat 控制台输出专用格式(可选)
	// 可选值: json, console
	// 如果为空,则使用 Format 的值
	// 使用场景: 希望控制台用易读的 console 格式,文件用 json 格式
	ConsoleFormat string `mapstructure:"console_format"`

	// FileFormat 文件输出专用格式(可选)
	// 可选值: json, console
	// 如果为空,则使用 Format 的值
	// 使用场景: 希望控制台用 console 格式,文件用 json 格式
	FileFormat string `mapstructure:"file_format"`

	// Output 输出目标
	// 可选值:
	// - stdout: 仅标准输出,适合容器环境(日志收集器会捕获)
	// - file: 仅文件输出,需要配合 FilePath,适合传统部署
	// - both: 同时输出到文件和标准输出,适合开发环境
	// 推荐:
	// - 容器/K8s 环境: stdout
	// - 传统部署: file
	// - 开发环境: both
	Output string `mapstructure:"output"`

	// FilePath 日志文件路径
	// 仅当 Output="file" 或 Output="both" 时有效
	// 例如: /var/log/app/app.log
	// 注意:
	// - 确保目录存在且有写权限
	// - 建议使用绝对路径
	FilePath string `mapstructure:"file_path"`

	// MaxSize 单个日志文件的最大大小(MB)
	// 超过此大小会触发日志轮转
	// 推荐值: 100-500 MB
	// 设置过大:单个文件难以处理
	// 设置过小:文件过多
	MaxSize int `mapstructure:"max_size"`

	// MaxBackups 保留的旧日志文件最大数量
	// 超过此数量的旧文件会被删除
	// 推荐值: 3-10
	// 用途:
	// - 防止日志占满磁盘
	// - 保留足够的历史日志用于问题排查
	MaxBackups int `mapstructure:"max_backups"`

	// MaxAge 保留旧日志文件的最大天数
	// 超过此天数的日志文件会被删除
	// 推荐值: 7-30 天
	// 考虑因素:
	// - 法规要求(某些行业要求保留审计日志)
	// - 磁盘空间
	// - 问题排查需求
	MaxAge int `mapstructure:"max_age"`
}

// I18nConfig 国际化配置
// 支持多语言
type I18nConfig struct {
	// Default 默认语言
	// 当请求的语言不支持时使用
	// 例如: en, zh-CN, ja
	Default string `mapstructure:"default"`

	// Supported 支持的语言列表
	// 必须包含 Default 语言
	// 例如: ["en", "zh-CN", "ja"]
	Supported []string `mapstructure:"supported"`
}

// ExecutorPoolConfig 单个执行器池的配置
// 定义了一个协程池的行为参数
type ExecutorPoolConfig struct {
	// Name 池的唯一标识符
	// 例如: "http", "database", "background"
	// 业务层应该定义常量来引用这些名称
	Name string `mapstructure:"name"`

	// Size 池的容量,即最大并发 worker 数量
	// 推荐值:
	// - CPU 密集型: runtime.NumCPU() * 2
	// - IO 密集型: 100-500
	Size int `mapstructure:"size"`

	// Expiry worker 的过期时间(秒)
	// 闲置超过此时间的 worker 会被回收
	// 推荐: 10-60 秒
	Expiry int `mapstructure:"expiry"`

	// NonBlocking 是否使用非阻塞模式
	// true:  池满时立即返回错误
	// false: 池满时阻塞等待
	// 推荐使用 true
	NonBlocking bool `mapstructure:"nonBlocking"`
}

// ExecutorConfig 执行器配置
// 管理多个异步任务池
type ExecutorConfig struct {
	// Enabled 是否启用执行器
	// false 时,应用不会创建执行器
	Enabled bool `mapstructure:"enabled"`

	// Pools 池配置列表
	// 每个池可以有不同的参数
	Pools []ExecutorPoolConfig `mapstructure:"pools"`
}

// JWTConfig JWT认证配置
// 用于token的生成和验证
type JWTConfig struct {
	// Secret 签名密钥
	// 生产环境必须从环境变量设置
	// 建议使用至少32个字符的随机字符串
	// 注意: 此字段非常敏感,必须保密
	Secret string `mapstructure:"secret"`

	// ExpiresIn 令牌有效期（秒）
	// 默认: 3600（1小时）
	// 考虑因素:
	// - 安全性: 过期时间越短越安全
	// - 用户体验: 过期时间太短需频繁登录
	// - 业务场景: 根据业务敏感度调整
	ExpiresIn int `mapstructure:"expiresIn"`

	// Issuer 签发者
	// 标识令牌由哪个系统签发
	// 用于多系统环境下区分token来源
	// 默认: "go-scaffold"
	Issuer string `mapstructure:"issuer"`
}

// Validate 验证整个配置
// 实现 Configurable 接口
// 会递归验证所有子配置
// 返回:
//
//	error: 第一个验证失败的错误,包含错误路径
func (c *Config) Validate() error {
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

// Validate 验证服务器配置
// 实现 Configurable 接口
func (c *ServerConfig) Validate() error {
	// 验证端口范围
	if c.Port <= 0 || c.Port > 65535 {
		// 端口必须在 1-65535 之间
		// 1-1023 是保留端口,需要 root 权限
		// 1024-65535 是用户端口
		return errors.New("port must be between 1 and 65535")
	}

	// 验证运行模式
	if c.Mode != "debug" && c.Mode != "release" && c.Mode != "test" {
		// 只允许这三种模式
		return errors.New("mode must be debug, release, or test")
	}

	// 验证读取超时
	if c.ReadTimeout <= 0 {
		// 必须大于 0
		// 0 表示无超时,可能导致连接泄漏
		return errors.New("readTimeout must be positive")
	}

	// 验证写入超时
	if c.WriteTimeout <= 0 {
		// 必须大于 0
		return errors.New("writeTimeout must be positive")
	}

	return nil
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

// Validate 验证日志配置
// 实现 Configurable 接口
func (c *LoggerConfig) Validate() error {
	// 验证日志级别
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[c.Level] {
		return errors.New("level must be debug, info, warn, or error")
	}

	// 验证日志格式
	validFormats := map[string]bool{"json": true, "console": true}
	if !validFormats[c.Format] {
		return errors.New("format must be json or console")
	}

	// 验证输出目标
	validOutputs := map[string]bool{"stdout": true, "file": true, "both": true}
	if !validOutputs[c.Output] {
		return errors.New("output must be stdout, file, or both")
	}

	return nil
}

// Validate 验证国际化配置
// 实现 Configurable 接口
func (c *I18nConfig) Validate() error {
	// 验证默认语言
	if c.Default == "" {
		return errors.New("default locale is required")
	}

	// 验证支持的语言列表
	if len(c.Supported) == 0 {
		return errors.New("at least one supported locale is required")
	}

	// 确保默认语言在支持列表中
	found := false
	for _, s := range c.Supported {
		if s == c.Default {
			found = true
			break
		}
	}
	if !found {
		// 默认语言必须是支持的语言之一
		// 否则会导致运行时错误
		return errors.New("default locale must be in supported list")
	}

	return nil
}

// Validate 验证执行器配置
// 实现 Configurable 接口
func (c *ExecutorConfig) Validate() error {
	// 如果未启用,跳过验证
	if !c.Enabled {
		return nil
	}

	// 启用时必须至少有一个池
	if len(c.Pools) == 0 {
		return errors.New("at least one pool is required when executor is enabled")
	}

	// 验证每个池的配置
	poolNames := make(map[string]bool)
	for i, pool := range c.Pools {
		// 验证池名称
		if pool.Name == "" {
			return fmt.Errorf("pool %d: name is required", i)
		}

		// 检查重复名称
		if poolNames[pool.Name] {
			return fmt.Errorf("duplicate pool name: %s", pool.Name)
		}
		poolNames[pool.Name] = true

		// 验证池大小
		if pool.Size <= 0 {
			return fmt.Errorf("pool %s: size must be positive", pool.Name)
		}
		if pool.Size > 10000 {
			return fmt.Errorf("pool %s: size must not exceed 10000", pool.Name)
		}

		// 验证过期时间
		if pool.Expiry < 0 {
			return fmt.Errorf("pool %s: expiry must be non-negative", pool.Name)
		}
	}

	return nil
}

// Validate 验证 JWT 配置
// 实现 Configurable 接口
func (c *JWTConfig) Validate() error {
	// 验证密钥
	if c.Secret == "" {
		return errors.New("jwt secret is required")
	}

	// 验证密钥长度（安全性要求）
	if len(c.Secret) < 32 {
		return errors.New("jwt secret must be at least 32 characters")
	}

	// 验证过期时间
	if c.ExpiresIn <= 0 {
		return errors.New("jwt expiresIn must be positive")
	}

	return nil
}
