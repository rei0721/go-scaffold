package config

import "fmt"

// 环境变量名称常量
// 定义所有支持的环境变量名称,避免魔法字符串
// 命名规范: <模块>_<字段名>,全大写,单词间用下划线分隔

// EnvPrefixJoin 拼接环境变量名称
func EnvPrefixJoin(field string) string {
	return fmt.Sprintf("%s_%s", EnvPrefix, field)
}

// 数据库相关环境变量
const (
	//EnvPrefix 数据库前缀
	EnvPrefix = "REI_APP"
	// EnvDBDriver 数据库驱动类型
	// 可选值: postgres, mysql, sqlite
	// 示例: export DB_DRIVER=postgres
	EnvDBDriver = "DB_DRIVER"

	// EnvDBHost 数据库主机地址
	// 示例: export DB_HOST=localhost
	EnvDBHost = "DB_HOST"

	// EnvDBPort 数据库端口
	// 示例: export DB_PORT=5432
	EnvDBPort = "DB_PORT"

	// EnvDBUser 数据库用户名
	// 示例: export DB_USER=postgres
	EnvDBUser = "DB_USER"

	// EnvDBPassword 数据库密码
	// 重要: 生产环境应该使用环境变量而不是配置文件
	// 示例: export DB_PASSWORD=secret
	EnvDBPassword = "DB_PASSWORD"

	// EnvDBName 数据库名称
	// 示例: export DB_NAME=myapp
	EnvDBName = "DB_NAME"

	// EnvDBMaxOpenConns 最大打开连接数
	// 示例: export DB_MAX_OPEN_CONNS=100
	EnvDBMaxOpenConns = "DB_MAX_OPEN_CONNS"

	// EnvDBMaxIdleConns 最大空闲连接数
	// 示例: export DB_MAX_IDLE_CONNS=10
	EnvDBMaxIdleConns = "DB_MAX_IDLE_CONNS"
)

// Redis 相关环境变量
const (
	// EnvRedisEnabled Redis 是否启用
	// 可选值: true, false
	// 示例: export REDIS_ENABLED=true
	EnvRedisEnabled = "REDIS_ENABLED"

	// EnvRedisHost Redis 主机地址
	// 示例: export REDIS_HOST=localhost
	EnvRedisHost = "REDIS_HOST"

	// EnvRedisPort Redis 端口
	// 示例: export REDIS_PORT=6379
	EnvRedisPort = "REDIS_PORT"

	// EnvRedisPassword Redis 密码
	// 重要: 生产环境应该使用环境变量
	// 示例: export REDIS_PASSWORD=secret
	EnvRedisPassword = "REDIS_PASSWORD"

	// EnvRedisDB Redis 数据库索引
	// 示例: export REDIS_DB=0
	EnvRedisDB = "REDIS_DB"

	// EnvRedisPoolSize Redis 连接池大小
	// 示例: export REDIS_POOL_SIZE=20
	EnvRedisPoolSize = "REDIS_POOL_SIZE"

	// EnvRedisMinIdleConns Redis 最小空闲连接数
	// 示例: export REDIS_MIN_IDLE_CONNS=10
	EnvRedisMinIdleConns = "REDIS_MIN_IDLE_CONNS"

	// EnvRedisMaxRetries Redis 最大重试次数
	// 示例: export REDIS_MAX_RETRIES=3
	EnvRedisMaxRetries = "REDIS_MAX_RETRIES"

	// EnvRedisDialTimeout Redis 连接超时(秒)
	// 示例: export REDIS_DIAL_TIMEOUT=5
	EnvRedisDialTimeout = "REDIS_DIAL_TIMEOUT"

	// EnvRedisReadTimeout Redis 读取超时(秒)
	// 示例: export REDIS_READ_TIMEOUT=3
	EnvRedisReadTimeout = "REDIS_READ_TIMEOUT"

	// EnvRedisWriteTimeout Redis 写入超时(秒)
	// 示例: export REDIS_WRITE_TIMEOUT=3
	EnvRedisWriteTimeout = "REDIS_WRITE_TIMEOUT"
)

// 服务器相关环境变量
const (
	// EnvServerPort HTTP 服务器端口
	// 示例: export SERVER_PORT=8080
	EnvServerPort = "SERVER_PORT"

	// EnvServerMode 服务器运行模式
	// 可选值: debug, release, test
	// 示例: export SERVER_MODE=release
	EnvServerMode = "SERVER_MODE"

	// EnvServerReadTimeout 读取超时(秒)
	// 示例: export SERVER_READ_TIMEOUT=30
	EnvServerReadTimeout = "SERVER_READ_TIMEOUT"

	// EnvServerWriteTimeout 写入超时(秒)
	// 示例: export SERVER_WRITE_TIMEOUT=30
	EnvServerWriteTimeout = "SERVER_WRITE_TIMEOUT"
)

// 日志相关环境变量
const (
	// EnvLogLevel 日志级别
	// 可选值: debug, info, warn, error
	// 示例: export LOG_LEVEL=info
	EnvLogLevel = "LOG_LEVEL"

	// EnvLogFormat 日志格式
	// 可选值: json, console
	// 示例: export LOG_FORMAT=json
	EnvLogFormat = "LOG_FORMAT"

	// EnvLogOutput 日志输出目标
	// 可选值: stdout, file
	// 示例: export LOG_OUTPUT=stdout
	EnvLogOutput = "LOG_OUTPUT"
)

// 国际化相关环境变量
const (
	// EnvI18nDefault 默认语言
	// 示例: export I18N_DEFAULT=zh-CN
	EnvI18nDefault = "I18N_DEFAULT"

	// EnvI18nSupported 支持的语言列表(逗号分隔)
	// 示例: export I18N_SUPPORTED=zh-CN,en-US,ja-JP
	EnvI18nSupported = "I18N_SUPPORTED"
)

// JWT 相关环境变量
const (
	// EnvJWTSecret JWT 签名密钥
	// 重要: 生产环境必须设置,至少32个字符
	// 示例: export JWT_SECRET=your-secret-key-at-least-32-characters-long
	EnvJWTSecret = "JWT_SECRET"

	// EnvJWTExpiresIn JWT 令牌有效期（秒）
	// 示例: export JWT_EXPIRES_IN=3600
	EnvJWTExpiresIn = "JWT_EXPIRES_IN"

	// EnvJWTIssuer JWT 签发者
	// 示例: export JWT_ISSUER=go-scaffold
	EnvJWTIssuer = "JWT_ISSUER"
)

// 其他常量
const (
	// EnvFilePath .env 文件路径
	// 默认在项目根目录
	EnvFilePath = ".env"

	// EnvFilePathExample .env 示例文件路径
	EnvFilePathExample = ".env.example"
)

// 环境变量解析相关常量
const (
	// DefaultSeparator 列表类型环境变量的分隔符
	// 用于解析逗号分隔的值,如语言列表
	// 示例: "zh-CN,en-US,ja-JP" -> ["zh-CN", "en-US", "ja-JP"]
	DefaultSeparator = ","
)

// 应用配置名称常量
const (
	AppServerName   = "server"
	AppDatabaseName = "database"
	AppRedisName    = "redis"
	AppLoggerName   = "logger"
	AppI18nName     = "i18n"
	AppExecutorName = "executor"
	AppJWTName      = "jwt"
	AppInitDBName   = "initdb"
	AppRBACName     = "rbac"
	AppStorageName  = "storage"
)
