# 变量命名索引 (Variable Index)

> **命名宪法 (Variable Naming Constitution)**  
> 先检索，后复用。禁止造同义词。

---

## 📋 索引说明

本文档记录项目中所有全局常量的命名规范，用于：

- ✅ 复用已有常量，避免重复定义
- ✅ 统一命名风格，提高代码可读性
- ✅ 快速检索，确定常量所在位置

**使用原则**：

1. 新增常量前，必须先检索本文档
2. 如果已存在语义相同的常量，必须复用
3. 新增全局常量后，必须立即更新本文档

---

## 1. 包级常量 (Package Constants)

### 1.1 CLI 框架 (`pkg/cli`)

**文件**: [`pkg/cli/constants.go`](/pkg/cli/constants.go)

#### 退出码 (Exit Codes)

| 常量名            | 值  | 说明                  |
| ----------------- | --- | --------------------- |
| `ExitSuccess`     | 0   | 成功退出              |
| `ExitError`       | 1   | 通用错误              |
| `ExitUsage`       | 2   | 参数错误（Unix 约定） |
| `ExitConfig`      | 3   | 配置错误              |
| `ExitInterrupted` | 130 | 用户中断 (Ctrl+C)     |

#### 错误消息

| 常量名                   | 值                       |
| ------------------------ | ------------------------ |
| `ErrMsgCommandNotFound`  | "command not found"      |
| `ErrMsgInvalidArgs`      | "invalid arguments"      |
| `ErrMsgMissingRequired`  | "missing required flag"  |
| `ErrMsgDuplicateCommand` | "duplicate command name" |
| `ErrMsgCancelled`        | "operation cancelled"    |
| `ErrMsgInvalidFlagValue` | "invalid flag value"     |

#### 默认选项

| 常量名               | 值        |
| -------------------- | --------- |
| `DefaultHelpFlag`    | "help"    |
| `DefaultVersionFlag` | "version" |

---

### 1.2 缓存 (`pkg/cache`)

**文件**: [`pkg/cache/constants.go`](/pkg/cache/constants.go)

#### 默认配置

| 常量名                | 值          | 说明            |
| --------------------- | ----------- | --------------- |
| `DefaultHost`         | "localhost" | 默认 Redis 主机 |
| `DefaultPort`         | 6379        | Redis 标准端口  |
| `DefaultDB`           | 0           | 默认数据库索引  |
| `DefaultPoolSize`     | 10          | 连接池大小      |
| `DefaultMinIdleConns` | 5           | 最小空闲连接数  |
| `DefaultMaxRetries`   | 3           | 最大重试次数    |
| `DefaultDialTimeout`  | 5           | 连接超时(秒)    |
| `DefaultReadTimeout`  | 3           | 读取超时(秒)    |
| `DefaultWriteTimeout` | 3           | 写入超时(秒)    |

#### 键前缀 (Key Prefixes)

| 常量名             | 前缀       | 用途         |
| ------------------ | ---------- | ------------ |
| `KeyPrefixUser`    | "user:"    | 用户相关数据 |
| `KeyPrefixSession` | "session:" | 会话数据     |
| `KeyPrefixCache`   | "cache:"   | 通用缓存数据 |
| `KeyPrefixLock`    | "lock:"    | 分布式锁     |
| `KeyPrefixCounter` | "counter:" | 计数器       |

#### 过期时间常量

| 常量名             | 值(秒) | 说明                   |
| ------------------ | ------ | ---------------------- |
| `ExpirationShort`  | 300    | 5 分钟，频繁变化的数据 |
| `ExpirationMedium` | 3600   | 1 小时，一般缓存数据   |
| `ExpirationLong`   | 86400  | 24 小时，稳定数据      |
| `ExpirationNever`  | 0      | 永不过期               |

#### 日志/错误消息

见 [`cache/constants.go:L47-L102`](/pkg/cache/constants.go#L47-L102)

---

### 1.3 数据库 (`pkg/database`)

**文件**: [`pkg/database/constants.go`](/pkg/database/constants.go)

#### 默认配置

| 常量名                   | 值     | 说明             |
| ------------------------ | ------ | ---------------- |
| `DefaultReloadTimeout`   | 30 秒  | 重载操作超时时间 |
| `DefaultConnMaxLifetime` | 1 小时 | 连接最大生命周期 |

#### 错误消息

| 常量名                           | 值                                         |
| -------------------------------- | ------------------------------------------ |
| `ErrMsgFailedToCreateConnection` | "failed to create new database connection" |
| `ErrMsgConnectionPingFailed`     | "database connection ping failed"          |
| `ErrMsgFailedToCloseConnection`  | "failed to close database connection"      |
| `ErrMsgUnsupportedDriver`        | "unsupported database driver"              |

---

### 1.4 执行器 (`pkg/executor`)

**文件**: [`pkg/executor/constants.go`](/pkg/executor/constants.go)

#### 预定义错误

| 变量名             | 类型  | 说明                 |
| ------------------ | ----- | -------------------- |
| `ErrPoolNotFound`  | error | 池不存在             |
| `ErrPoolOverload`  | error | 池过载（非阻塞模式） |
| `ErrManagerClosed` | error | 管理器已关闭         |
| `ErrInvalidConfig` | error | 无效配置             |

#### 默认配置

| 常量名                | 值    | 说明            |
| --------------------- | ----- | --------------- |
| `DefaultPoolSize`     | 100   | 默认池大小      |
| `DefaultWorkerExpiry` | 10 秒 | Worker 过期时间 |
| `DefaultNonBlocking`  | true  | 默认非阻塞模式  |
| `ShutdownTimeout`     | 5 秒  | 关闭超时时间    |
| `MinPoolSize`         | 1     | 最小池大小      |
| `MaxPoolSize`         | 10000 | 最大池大小      |

---

### 1.5 日志 (`pkg/logger`)

**文件**: [`pkg/logger/constants.go`](/pkg/logger/constants.go)

#### 日志级别

| 常量名       | 枚举值 |
| ------------ | ------ |
| `LevelDebug` | -1     |
| `LevelInfo`  | 0      |
| `LevelWarn`  | 1      |
| `LevelError` | 2      |
| `LevelFatal` | 3      |

#### 输出模式

| 常量名         | 值       | 说明     |
| -------------- | -------- | -------- |
| `OutputStdout` | "stdout" | 仅控制台 |
| `OutputFile`   | "file"   | 仅文件   |
| `OutputBoth`   | "both"   | 同时输出 |

#### 默认值

| 常量名          | 值        |
| --------------- | --------- |
| `DefaultLevel`  | "debug"   |
| `DefaultFormat` | "console" |
| `DefaultOutput` | "stdout"  |

---

### 1.6 国际化 (`pkg/i18n`)

**文件**: [`pkg/i18n/constants.go`](/pkg/i18n/constants.go)

#### 语言代码 (BCP 47)

| 常量名             | 值      | 说明               |
| ------------------ | ------- | ------------------ |
| `DefaultLanguage`  | "zh-CN" | 默认语言(简体中文) |
| `LanguageEnglish`  | "en-US" | 英语(美国)         |
| `LanguageChinese`  | "zh-CN" | 简体中文           |
| `LanguageJapanese` | "ja-JP" | 日语               |

#### HTTP 头部

| 常量名           | 值                |
| ---------------- | ----------------- |
| `LanguageHeader` | "Accept-Language" |

#### 文件格式

| 常量名               | 值     |
| -------------------- | ------ |
| `FilenameFormatJson` | "json" |
| `FilenameFormatYaml` | "yaml" |
| `FilenameFormatYml`  | "yml"  |

---

### 1.7 HTTP Server (`pkg/httpserver`)

**文件**: [`pkg/httpserver/constants.go`](/pkg/httpserver/constants.go)

#### 默认配置

| 常量名                | 值                | 说明         |
| --------------------- | ----------------- | ------------ |
| `DefaultHost`         | "localhost"       | 默认监听地址 |
| `DefaultPort`         | 8080              | 默认端口     |
| `DefaultReadTimeout`  | 15 \* time.Second | 读取超时     |
| `DefaultWriteTimeout` | 15 \* time.Second | 写入超时     |
| `DefaultIdleTimeout`  | 60 \* time.Second | 空闲超时     |

#### 错误消息

| 常量名                       | 值                               |
| ---------------------------- | -------------------------------- |
| `ErrMsgInvalidAddress`       | "invalid listen address"         |
| `ErrMsgServerStartFailed`    | "failed to start server"         |
| `ErrMsgServerShutdownFailed` | "failed to shutdown server"      |
| `ErrMsgPortUnavailable`      | "port is not available"          |
| `ErrMsgServerAlreadyRunning` | "server is already running"      |
| `ErrMsgServerNotRunning`     | "server is not running"          |
| `ErrMsgInvalidConfig`        | "invalid server config"          |
| `ErrMsgReloadFailed`         | "failed to reload server config" |

---

### 1.8 JWT (`pkg/jwt`)

**文件**: [`pkg/jwt/constants.go`](/pkg/jwt/constants.go)

#### 默认配置

| 常量名             | 值            | 说明                   |
| ------------------ | ------------- | ---------------------- |
| `DefaultExpiresIn` | 3600          | 默认过期时间（1 小时） |
| `DefaultIssuer`    | "go-scaffold" | 默认签发者             |

#### 预定义错误

| 变量名                | 类型  | 说明           |
| --------------------- | ----- | -------------- |
| `ErrInvalidToken`     | error | Token 无效     |
| `ErrExpiredToken`     | error | Token 已过期   |
| `ErrTokenNotYetValid` | error | Token 尚未生效 |
| `ErrInvalidSignature` | error | 签名无效       |
| `ErrMissingSecret`    | error | 缺少签名密钥   |

#### 错误消息

| 常量名                   | 值                                          |
| ------------------------ | ------------------------------------------- |
| `ErrMsgInvalidToken`     | "invalid token"                             |
| `ErrMsgExpiredToken`     | "token has expired"                         |
| `ErrMsgTokenNotYetValid` | "token not yet valid"                       |
| `ErrMsgInvalidSignature` | "invalid signature"                         |
| `ErrMsgMissingSecret`    | "jwt secret is required"                    |
| `ErrMsgSecretTooShort`   | "jwt secret must be at least 32 characters" |

---

## 2. 配置键名 (Configuration Keys)

### 2.1 内部配置常量 (`internal/config`)

**文件**: [`internal/config/constants.go`](/internal/config/constants.go)

#### 环境变量名称

**命名规范**: `<模块>_<字段名>`，全大写，单词间下划线分隔

##### 数据库

| 常量名              | 环境变量            | 说明                              |
| ------------------- | ------------------- | --------------------------------- |
| `EnvDBDriver`       | `DB_DRIVER`         | 数据库驱动(postgres/mysql/sqlite) |
| `EnvDBHost`         | `DB_HOST`           | 数据库主机                        |
| `EnvDBPort`         | `DB_PORT`           | 数据库端口                        |
| `EnvDBUser`         | `DB_USER`           | 数据库用户名                      |
| `EnvDBPassword`     | `DB_PASSWORD`       | 数据库密码(**敏感**)              |
| `EnvDBName`         | `DB_NAME`           | 数据库名称                        |
| `EnvDBMaxOpenConns` | `DB_MAX_OPEN_CONNS` | 最大打开连接数                    |
| `EnvDBMaxIdleConns` | `DB_MAX_IDLE_CONNS` | 最大空闲连接数                    |

##### Redis

| 常量名                 | 环境变量               | 说明                 |
| ---------------------- | ---------------------- | -------------------- |
| `EnvRedisEnabled`      | `REDIS_ENABLED`        | 是否启用 Redis       |
| `EnvRedisHost`         | `REDIS_HOST`           | Redis 主机           |
| `EnvRedisPort`         | `REDIS_PORT`           | Redis 端口           |
| `EnvRedisPassword`     | `REDIS_PASSWORD`       | Redis 密码(**敏感**) |
| `EnvRedisDB`           | `REDIS_DB`             | Redis 数据库索引     |
| `EnvRedisPoolSize`     | `REDIS_POOL_SIZE`      | 连接池大小           |
| `EnvRedisMinIdleConns` | `REDIS_MIN_IDLE_CONNS` | 最小空闲连接         |
| `EnvRedisMaxRetries`   | `REDIS_MAX_RETRIES`    | 最大重试次数         |
| `EnvRedisDialTimeout`  | `REDIS_DIAL_TIMEOUT`   | 连接超时(秒)         |
| `EnvRedisReadTimeout`  | `REDIS_READ_TIMEOUT`   | 读取超时(秒)         |
| `EnvRedisWriteTimeout` | `REDIS_WRITE_TIMEOUT`  | 写入超时(秒)         |

##### 服务器

| 常量名                  | 环境变量               | 说明                         |
| ----------------------- | ---------------------- | ---------------------------- |
| `EnvServerPort`         | `SERVER_PORT`          | HTTP 端口                    |
| `EnvServerMode`         | `SERVER_MODE`          | 运行模式(debug/release/test) |
| `EnvServerReadTimeout`  | `SERVER_READ_TIMEOUT`  | 读取超时(秒)                 |
| `EnvServerWriteTimeout` | `SERVER_WRITE_TIMEOUT` | 写入超时(秒)                 |

##### 日志

| 常量名         | 环境变量     | 说明     |
| -------------- | ------------ | -------- |
| `EnvLogLevel`  | `LOG_LEVEL`  | 日志级别 |
| `EnvLogFormat` | `LOG_FORMAT` | 日志格式 |
| `EnvLogOutput` | `LOG_OUTPUT` | 日志输出 |

##### 国际化

| 常量名             | 环境变量         | 说明                     |
| ------------------ | ---------------- | ------------------------ |
| `EnvI18nDefault`   | `I18N_DEFAULT`   | 默认语言                 |
| `EnvI18nSupported` | `I18N_SUPPORTED` | 支持的语言列表(逗号分隔) |

##### 其他

| 常量名               | 值             | 说明          |
| -------------------- | -------------- | ------------- |
| `EnvPrefix`          | "REI_APP"      | 环境变量前缀  |
| `EnvFilePath`        | ".env"         | .env 文件路径 |
| `EnvFilePathExample` | ".env.example" | 示例文件路径  |
| `DefaultSeparator`   | ","            | 列表分隔符    |

### 2.2 YAML 配置键

**文件**: [`internal/config/config.go`](/internal/config/config.go)

**使用 `mapstructure` 标签映射**，结构如下：

```yaml
server:          # ServerConfig
  host: string
  port: int
  mode: string
  readTimeout: int
  writeTimeout: int
  idleTimeout: int

database:        # DatabaseConfig
  driver: string
  host: string
  port: int
  user: string
  password: string
  dbname: string
  maxOpenConns: int
  maxIdleConns: int

redis:           # RedisConfig
  enabled: bool
  host: string
  port: int
  password: string
  db: int
  poolSize: int
  minIdleConns: int
  maxRetries: int
  dialTimeout: int
  readTimeout: int
  writeTimeout: int

logger:          # LoggerConfig
  level: string
  format: string
  console_format: string
  file_format: string
  output: string
  file_path: string
  max_size: int
  max_backups: int
  max_age: int

i18n:            # I18nConfig
  default: string
  supported: []string

executor:        # ExecutorConfig
  enabled: bool
  pools:
    - name: string
      size: int
      expiry: int
      nonBlocking: bool
```

---

## 3. 应用层常量 (`internal/app`)

**文件**: [`internal/app/constants.go`](/internal/app/constants.go)

#### 启动模式

| 常量名       | 值       | 说明                                   |
| ------------ | -------- | -------------------------------------- |
| `ModeServer` | "server" | 服务器模式（默认），完整启动流程       |
| `ModeInitDB` | "initdb" | 数据库初始化模式，执行初始化脚本后退出 |

#### 国际化

| 常量名                            | 值                  | 说明         |
| --------------------------------- | ------------------- | ------------ |
| `ConstantsI18nMessagesDir`        | "./configs/locales" | 翻译文件目录 |
| `ConstantsI18nDefaultLanguage`    | "zh-CN"             | 默认语言     |
| `ConstantsDefaultHost`            | "localhost"         | 默认主机     |
| `ConstantsI18nSupportedLanguages` | ["zh-CN", "en-US"]  | 支持的语言   |

#### 数据库初始化

| 常量名                 | 值                 | 说明                 |
| ---------------------- | ------------------ | -------------------- |
| `InitDBScriptDir`      | "./scripts/initdb" | 初始化 SQL 脚本目录  |
| `InitDBLockFile`       | ".initialized"     | 初始化锁文件名       |
| `InitDBScriptFileName` | "init\_%s.sql"     | 初始化脚本文件名模板 |

---

## 4. 业务层常量 (`internal/service`)

**文件**: [`internal/service/constants.go`](/internal/service/constants.go)

### 4.1 缓存相关

#### 缓存键前缀

| 常量名               | 值      | 说明               |
| -------------------- | ------- | ------------------ |
| `CacheKeyPrefixUser` | "user:" | 用户数据缓存键前缀 |

#### 缓存过期时间

| 常量名         | 值                | 说明                        |
| -------------- | ----------------- | --------------------------- |
| `CacheTTLUser` | 30 \* time.Minute | 用户缓存过期时间（30 分钟） |

---

## 5. 错误类型/错误码

### 5.1 CLI 错误

见 [`pkg/cli/error.go`](/pkg/cli/error.go)（推断）

- `UsageError` - 参数错误（ExitUsage = 2）
- `CommandError` - 命令执行错误（ExitError = 1）
- `CancelledError` - 用户取消（ExitInterrupted = 130）

### 4.2 Executor 错误

**文件**: [`pkg/executor/constants.go:L31-L50`](/pkg/executor/constants.go#L31-L50)

- `ErrPoolNotFound` - 池不存在
- `ErrPoolOverload` - 池过载
- `ErrManagerClosed` - 管理器已关闭
- `ErrInvalidConfig` - 无效配置

---

## 5. 资源名称/池名称

### 5.1 Executor 池名称

**定义位置**: [`types/constants/executor.go`](file:///D:/coder/go/PicHub/main/types/constants/executor.go)

| 常量名           | 值           | 类型                | 非阻塞 | 池大小 | 用途              |
| ---------------- | ------------ | ------------------- | ------ | ------ | ----------------- |
| `PoolHTTP`       | "http"       | `executor.PoolName` | 是     | 200    | HTTP 请求异步处理 |
| `PoolDatabase`   | "database"   | `executor.PoolName` | 否     | 50     | 数据库异步操作    |
| `PoolCache`      | "cache"      | `executor.PoolName` | 是     | 30     | 缓存异步更新      |
| `PoolLogger`     | "logger"     | `executor.PoolName` | 否     | 10     | 日志异步处理      |
| `PoolBackground` | "background" | `executor.PoolName` | 是     | 30     | 通用后台任务      |

**使用示例**:

```go
import (
    "github.com/rei0721/rei0721/types/constants"
    "github.com/rei0721/rei0721/pkg/executor"
)

// 在业务代码中使用常量引用池名
func (s *userService) CreateAsync(ctx context.Context) error {
    return s.executor.Execute(constants.PoolDatabase, func() {
        // 异步数据库任务
    })
}

// 在HTTP Server中使用
func handleRequest() {
    _ = executor.Execute(constants.PoolHTTP, func() {
        // 异步HTTP任务
    })
}
```

**命名规范**:

- 前缀：`Pool` + 功能描述（PascalCase）
- 类型：必须使用 `executor.PoolName` 类型
- 值：与 `config.yaml` 中的池名称严格一致（小写）

---

## 6. 命名规范总结

### 6.1 常量命名风格

| 层级         | 命名风格         | 示例                     |
| ------------ | ---------------- | ------------------------ |
| 导出常量     | PascalCase       | `DefaultPoolSize`        |
| 私有常量     | camelCase        | `defaultTimeout`（少用） |
| 环境变量     | UPPER_SNAKE_CASE | `DB_PASSWORD`            |
| 配置键(YAML) | snake_case       | `max_open_conns`         |

### 6.2 前缀约定

| 前缀               | 用途          | 示例                |
| ------------------ | ------------- | ------------------- |
| `Default*`         | 默认值        | `DefaultHost`       |
| `Err*` / `ErrMsg*` | 错误/错误消息 | `ErrPoolNotFound`   |
| `Msg*`             | 日志消息      | `MsgCacheConnected` |
| `Env*`             | 环境变量名    | `EnvDBPassword`     |
| `KeyPrefix*`       | 缓存键前缀    | `KeyPrefixUser`     |
| `Exit*`            | 退出码        | `ExitSuccess`       |
| `Expiration*`      | 过期时间      | `ExpirationShort`   |

### 6.3 值的约定

| 类型           | 约定                                          |
| -------------- | --------------------------------------------- |
| 超时时间       | 使用 `time.Duration`（如 `30 * time.Second`） |
| 整数超时配置   | 单位为**秒**（配置文件中）                    |
| 布尔值环境变量 | "true" / "false" (字符串)                     |
| 列表环境变量   | 逗号分隔（`,`）                               |

---

## 🔄 更新日志

| 日期       | 变更内容                                                     |
| ---------- | ------------------------------------------------------------ |
| 2026-01-15 | 新增业务层缓存常量（CacheKeyPrefixUser、CacheTTLUser）       |
| 2026-01-15 | 新增启动模式常量（ModeServer、ModeInitDB）和数据库初始化常量 |
| 2026-01-15 | 新增 `pkg/httpserver` 包常量（HTTP Server 封装）             |
| 2026-01-15 | 初始创建，扫描所有现有常量定义                               |

---

> **提醒**: 新增全局常量后，请立即更新本文档对应章节！
