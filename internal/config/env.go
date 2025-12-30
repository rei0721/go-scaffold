package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// LoadEnv 加载 .env 文件
// .env 文件是可选的,如果不存在不会报错
// 这个函数应该在加载 config.yaml 之前调用
//
// 工作流程:
//  1. 尝试加载项目根目录的 .env 文件
//  2. 如果文件不存在,静默跳过
//  3. 如果文件存在但格式错误,记录错误但不中断
//
// 使用场景:
//   - 本地开发: 创建 .env 文件存放敏感配置
//   - 生产环境: 不使用 .env 文件,直接使用系统环境变量
//
// 注意事项:
//   - .env 文件不应该提交到 Git
//   - .env 文件中的变量会被加载到进程环境变量中
//   - 如果系统环境变量已存在同名变量,.env 文件的值会被忽略
func LoadEnv() {
	// 尝试加载 .env 文件
	// godotenv.Load() 会:
	// 1. 读取 .env 文件
	// 2. 解析 KEY=VALUE 格式
	// 3. 将变量设置到进程环境变量中
	// 4. 不会覆盖已存在的环境变量
	err := godotenv.Load(EnvFilePath)
	if err != nil {
		// .env 文件不存在或读取失败
		// 这是正常情况,不需要报错
		// 生产环境通常不使用 .env 文件

		// 调试: 打印到 stderr 以便诊断
		fmt.Fprintf(os.Stderr, "[DEBUG] .env file not loaded: %v\n", err)
		return
	}

	// .env 文件加载成功
	// 调试: 打印成功信息
	fmt.Fprintf(os.Stderr, "[DEBUG] .env file loaded successfully\n")

	// 调试: 打印一些关键环境变量
	fmt.Fprintf(os.Stderr, "[DEBUG] DB_DRIVER=%s\n", os.Getenv("DB_DRIVER"))
	fmt.Fprintf(os.Stderr, "[DEBUG] REDIS_HOST=%s\n", os.Getenv("REDIS_HOST"))
	fmt.Fprintf(os.Stderr, "[DEBUG] DB_HOST=%s\n", os.Getenv("DB_HOST"))
	fmt.Fprintf(os.Stderr, "[DEBUG] REDIS_ENABLED=%s\n", os.Getenv("REDIS_ENABLED"))
}

// OverrideWithEnv 使用环境变量覆盖配置
// 优先级: 环境变量 > config.yaml
//
// 参数:
//
//	cfg: 从 config.yaml 加载的配置
//
// 工作流程:
//  1. 检查每个支持的环境变量
//  2. 如果环境变量存在,使用其值覆盖配置
//  3. 如果环境变量不存在,保持 config.yaml 的值
//
// 使用示例:
//
//	config := loadFromYaml()
//	OverrideWithEnv(config)
//	// 此时 config 中的值可能已被环境变量覆盖
func OverrideWithEnv(cfg *Config) {
	// 调试: 显示开始覆盖配置
	fmt.Fprintf(os.Stderr, "[DEBUG] OverrideWithEnv: starting environment variable override\n")

	// 数据库配置
	overrideDatabaseConfig(&cfg.Database)

	// Redis 配置
	overrideRedisConfig(&cfg.Redis)

	// 服务器配置
	overrideServerConfig(&cfg.Server)

	// 日志配置
	overrideLoggerConfig(&cfg.Logger)

	// 国际化配置
	overrideI18nConfig(&cfg.I18n)

	// 调试: 显示覆盖后的值
	fmt.Fprintf(os.Stderr, "[DEBUG] After override - DB_DRIVER=%s, DB_HOST=%s, REDIS_ENABLED=%v\n",
		cfg.Database.Driver, cfg.Database.Host, cfg.Redis.Enabled)
}

// overrideDatabaseConfig 使用环境变量覆盖数据库配置
func overrideDatabaseConfig(cfg *DatabaseConfig) {
	// Driver
	if val := os.Getenv(EnvDBDriver); val != "" {
		cfg.Driver = val
	}

	// Host
	if val := os.Getenv(EnvDBHost); val != "" {
		cfg.Host = val
	}

	// Port
	if val := os.Getenv(EnvDBPort); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			cfg.Port = port
		}
	}

	// User
	if val := os.Getenv(EnvDBUser); val != "" {
		cfg.User = val
	}

	// Password
	// 密码应该优先使用环境变量
	if val := os.Getenv(EnvDBPassword); val != "" {
		cfg.Password = val
	}

	// DBName
	if val := os.Getenv(EnvDBName); val != "" {
		cfg.DBName = val
	}

	// MaxOpenConns
	if val := os.Getenv(EnvDBMaxOpenConns); val != "" {
		if conns, err := strconv.Atoi(val); err == nil {
			cfg.MaxOpenConns = conns
		}
	}

	// MaxIdleConns
	if val := os.Getenv(EnvDBMaxIdleConns); val != "" {
		if conns, err := strconv.Atoi(val); err == nil {
			cfg.MaxIdleConns = conns
		}
	}
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

// overrideServerConfig 使用环境变量覆盖服务器配置
func overrideServerConfig(cfg *ServerConfig) {
	// Port
	if val := os.Getenv(EnvServerPort); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			cfg.Port = port
		}
	}

	// Mode
	if val := os.Getenv(EnvServerMode); val != "" {
		cfg.Mode = val
	}

	// ReadTimeout
	if val := os.Getenv(EnvServerReadTimeout); val != "" {
		if timeout, err := strconv.Atoi(val); err == nil {
			cfg.ReadTimeout = timeout
		}
	}

	// WriteTimeout
	if val := os.Getenv(EnvServerWriteTimeout); val != "" {
		if timeout, err := strconv.Atoi(val); err == nil {
			cfg.WriteTimeout = timeout
		}
	}
}

// overrideLoggerConfig 使用环境变量覆盖日志配置
func overrideLoggerConfig(cfg *LoggerConfig) {
	// Level
	if val := os.Getenv(EnvLogLevel); val != "" {
		cfg.Level = val
	}

	// Format
	if val := os.Getenv(EnvLogFormat); val != "" {
		cfg.Format = val
	}

	// Output
	if val := os.Getenv(EnvLogOutput); val != "" {
		cfg.Output = val
	}
}

// overrideI18nConfig 使用环境变量覆盖国际化配置
func overrideI18nConfig(cfg *I18nConfig) {
	// Default
	if val := os.Getenv(EnvI18nDefault); val != "" {
		cfg.Default = val
	}

	// Supported
	// 环境变量格式: "zh-CN,en-US,ja-JP"
	// 解析为: ["zh-CN", "en-US", "ja-JP"]
	if val := os.Getenv(EnvI18nSupported); val != "" {
		langs := strings.Split(val, DefaultSeparator)
		// 去除空白
		var supported []string
		for _, lang := range langs {
			trimmed := strings.TrimSpace(lang)
			if trimmed != "" {
				supported = append(supported, trimmed)
			}
		}
		if len(supported) > 0 {
			cfg.Supported = supported
		}
	}
}

// getEnvOrDefault 获取环境变量,如果不存在则返回默认值
// 这是一个辅助函数,用于简化环境变量读取
//
// 参数:
//
//	key: 环境变量名称
//	defaultValue: 默认值
//
// 返回:
//
//	string: 环境变量的值,或默认值
//
// 使用示例:
//
//	host := getEnvOrDefault("DB_HOST", "localhost")
func getEnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为整数
// 如果环境变量不存在或转换失败,返回默认值
//
// 参数:
//
//	key: 环境变量名称
//	defaultValue: 默认值
//
// 返回:
//
//	int: 环境变量转换后的整数值,或默认值
//
// 使用示例:
//
//	port := getEnvAsInt("SERVER_PORT", 8080)
func getEnvAsInt(key string, defaultValue int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getEnvAsBool 获取环境变量并转换为布尔值
// 如果环境变量不存在或转换失败,返回默认值
//
// 参数:
//
//	key: 环境变量名称
//	defaultValue: 默认值
//
// 返回:
//
//	bool: 环境变量转换后的布尔值,或默认值
//
// 使用示例:
//
//	enabled := getEnvAsBool("REDIS_ENABLED", true)
func getEnvAsBool(key string, defaultValue bool) bool {
	if val := os.Getenv(key); val != "" {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
