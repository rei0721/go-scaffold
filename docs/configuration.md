# 配置管理指南

本文档详细说明了 Rei0721 项目的配置管理系统，包括配置文件、环境变量、热重载和最佳实践。

## 配置系统概述

Rei0721 采用分层配置系统，支持多种配置源和热重载机制。

### 配置优先级

配置值按以下优先级应用（从低到高）：

```
1. 代码默认值
   ↓
2. 配置文件 (configs/config.yaml)
   ↓
3. 环境变量 (${VAR_NAME})
   ↓
4. 运行时更新 (ConfigManager.Update)
```

## 配置文件

### 主配置文件 (configs/config.yaml)

```yaml
server:
  port: ${SERVER_PORT:8080}
  mode: ${SERVER_MODE:debug}
  readTimeout: 10
  writeTimeout: 10

database:
  driver: ${DB_DRIVER:postgres}
  host: ${DB_HOST:localhost}
  port: ${DB_PORT:5432}
  user: ${DB_USER:postgres}
  password: ${DB_PASSWORD:password}
  dbname: ${DB_NAME:rei0721}
  maxOpenConns: 100
  maxIdleConns: 10

redis:
  enabled: true
  host: ${REDIS_HOST:localhost}
  port: ${REDIS_PORT:6379}
  password: ${REDIS_PASSWORD:}
  db: ${REDIS_DB:0}
  poolSize: 10

logger:
  level: ${LOG_LEVEL:info}
  format: ${LOG_FORMAT:json}
  output: stdout

i18n:
  default: zh-CN
  supported:
    - zh-CN
    - en-US
```

### 环境变量替换语法

配置文件支持 `${VAR_NAME:default_value}` 语法：

```yaml
# 使用环境变量，如果未设置则使用默认值
port: ${SERVER_PORT:8080}

# 只使用环境变量，如果未设置则为空
password: ${DB_PASSWORD}

# 直接使用值，不进行替换
timeout: 30
```

### 环境变量模板 (.env.example)

```bash
# 服务器配置
SERVER_PORT=8080
SERVER_MODE=debug

# 数据库配置
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=rei0721

# Redis 配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# 日志配置
LOG_LEVEL=info
LOG_FORMAT=json

# 其他配置
CONFIG_PATH=configs/config.yaml
```

## 配置结构

### Config 结构体

```go
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    Logger   LoggerConfig   `mapstructure:"logger"`
    I18n     I18nConfig     `mapstructure:"i18n"`
}

func (c *Config) Validate() error {
    // 验证配置有效性
}
```

### ServerConfig

```go
type ServerConfig struct {
    Port         int    `mapstructure:"port"`
    Mode         string `mapstructure:"mode"`        // debug, release, test
    ReadTimeout  int    `mapstructure:"readTimeout"`
    WriteTimeout int    `mapstructure:"writeTimeout"`
}
```

**说明**:
- `port` - 服务器监听端口 (1-65535)
- `mode` - 运行模式
  - `debug` - 调试模式，详细日志
  - `release` - 生产模式，性能优化
  - `test` - 测试模式
- `readTimeout` - 读超时 (秒)
- `writeTimeout` - 写超时 (秒)

### DatabaseConfig

```go
type DatabaseConfig struct {
    Driver       string `mapstructure:"driver"`
    Host         string `mapstructure:"host"`
    Port         int    `mapstructure:"port"`
    User         string `mapstructure:"user"`
    Password     string `mapstructure:"password"`
    DBName       string `mapstructure:"dbname"`
    MaxOpenConns int    `mapstructure:"maxOpenConns"`
    MaxIdleConns int    `mapstructure:"maxIdleConns"`
}
```

**说明**:
- `driver` - 数据库驱动 (postgres, mysql, sqlite)
- `host` - 数据库主机
- `port` - 数据库端口
- `user` - 数据库用户
- `password` - 数据库密码
- `dbname` - 数据库名称
- `maxOpenConns` - 最大连接数
- `maxIdleConns` - 最大空闲连接数

### RedisConfig

```go
type RedisConfig struct {
    Enabled  bool   `mapstructure:"enabled"`
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
    PoolSize int    `mapstructure:"poolSize"`
}
```

### LoggerConfig

```go
type LoggerConfig struct {
    Level  string `mapstructure:"level"`
    Format string `mapstructure:"format"`
    Output string `mapstructure:"output"`
}
```

**说明**:
- `level` - 日志级别 (debug, info, warn, error)
- `format` - 日志格式 (json, console)
- `output` - 输出目标 (stdout, file)

### I18nConfig

```go
type I18nConfig struct {
    Default   string   `mapstructure:"default"`
    Supported []string `mapstructure:"supported"`
}
```

## 配置管理器

### ConfigManager 接口

```go
type Manager interface {
    // 从文件加载配置
    Load(configPath string) error
    
    // 获取当前配置 (只读)
    Get() *Config
    
    // 原子更新配置
    Update(fn func(*Config)) error
    
    // 注册配置变更 Hook
    RegisterHook(h HookHandler)
    
    // 注册日志处理器
    RegisterLogger(h LoggerHandler) logger.Logger
    
    // 启动文件监听 (热重载)
    Watch() error
}
```

### 使用示例

#### 加载配置

```go
configManager := config.NewManager()
if err := configManager.Load("configs/config.yaml"); err != nil {
    log.Fatal("failed to load config:", err)
}

cfg := configManager.Get()
fmt.Println("Server port:", cfg.Server.Port)
```

#### 获取配置

```go
// 获取当前配置 (只读)
cfg := configManager.Get()

// 配置是原子的，可以安全地并发读取
go func() {
    cfg := configManager.Get()
    // 使用配置
}()
```

#### 更新配置

```go
// 原子更新配置
err := configManager.Update(func(cfg *config.Config) {
    cfg.Server.Port = 8081
    cfg.Logger.Level = "debug"
})

if err != nil {
    log.Error("failed to update config:", err)
}
```

#### 注册变更 Hook

```go
// 监听配置变更
configManager.RegisterHook(func(old, new *config.Config) {
    if old.Server.Port != new.Server.Port {
        log.Info("server port changed",
            "old", old.Server.Port,
            "new", new.Server.Port,
        )
    }
    
    if old.Logger.Level != new.Logger.Level {
        log.Info("log level changed",
            "old", old.Logger.Level,
            "new", new.Logger.Level,
        )
    }
})
```

#### 启动热重载

```go
// 启动文件监听
if err := configManager.Watch(); err != nil {
    log.Warn("failed to start config watcher:", err)
}

// 修改 configs/config.yaml 文件时，会自动重新加载
```

## 环境变量配置

### 设置环境变量

#### Linux/macOS

```bash
# 临时设置
export SERVER_PORT=8081
export DB_DRIVER=postgres

# 运行命令
go run ./cmd/server/main.go

# 或在命令前设置
SERVER_PORT=8081 go run ./cmd/server/main.go
```

#### Windows (CMD)

```cmd
set SERVER_PORT=8081
set DB_DRIVER=postgres
go run ./cmd/server/main.go
```

#### Windows (PowerShell)

```powershell
$env:SERVER_PORT = "8081"
$env:DB_DRIVER = "postgres"
go run ./cmd/server/main.go
```

### 从 .env 文件加载

```bash
# 创建 .env 文件
cp configs/.env.example .env

# 编辑 .env 文件
# SERVER_PORT=8081
# DB_DRIVER=postgres

# 加载环境变量 (Linux/macOS)
source .env

# 加载环境变量 (Windows PowerShell)
Get-Content .env | ForEach-Object {
    if ($_ -match '^\s*([^=]+)=(.*)$') {
        [Environment]::SetEnvironmentVariable($matches[1], $matches[2])
    }
}
```

## 配置验证

### 验证配置

```go
// 配置在加载时自动验证
if err := configManager.Load("configs/config.yaml"); err != nil {
    log.Fatal("invalid config:", err)
}

// 手动验证
cfg := configManager.Get()
if err := cfg.Validate(); err != nil {
    log.Fatal("config validation failed:", err)
}
```

### 验证规则

```go
func (c *Config) Validate() error {
    // 验证服务器配置
    if c.Server.Port < 1 || c.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", c.Server.Port)
    }
    
    // 验证数据库配置
    if c.Database.Driver == "" {
        return fmt.Errorf("database driver is required")
    }
    
    // 验证日志配置
    validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
    if !validLevels[c.Logger.Level] {
        return fmt.Errorf("invalid log level: %s", c.Logger.Level)
    }
    
    return nil
}
```

## 数据库配置

### PostgreSQL

```yaml
database:
  driver: postgres
  host: localhost
  port: 5432
  user: postgres
  password: password
  dbname: rei0721
  maxOpenConns: 100
  maxIdleConns: 10
```

### MySQL

```yaml
database:
  driver: mysql
  host: localhost
  port: 3306
  user: root
  password: password
  dbname: rei0721
  maxOpenConns: 100
  maxIdleConns: 10
```

### SQLite

```yaml
database:
  driver: sqlite
  dbname: ./rei0721.db
```

## 日志配置

### 日志级别

| 级别 | 说明 | 用途 |
|------|------|------|
| debug | 调试信息 | 开发调试 |
| info | 一般信息 | 正常运行 |
| warn | 警告信息 | 潜在问题 |
| error | 错误信息 | 错误处理 |

### 日志格式

#### JSON 格式

```json
{
  "level": "info",
  "ts": 1735475400.123,
  "msg": "user registered",
  "userId": 123,
  "email": "user@example.com"
}
```

#### Console 格式

```
2025-12-30T10:30:00.123+0800    info    user registered    {"userId": 123, "email": "user@example.com"}
```

### 日志输出

```yaml
logger:
  level: info
  format: json
  output: stdout    # stdout 或 file
```

## 国际化配置

### 支持的语言

```yaml
i18n:
  default: zh-CN
  supported:
    - zh-CN
    - en-US
```

### 消息文件

#### zh-CN.yaml

```yaml
errors:
  invalid_params: "参数无效"
  user_not_found: "用户不存在"
  duplicate_username: "用户名已存在"

messages:
  welcome: "欢迎使用 Rei0721"
  goodbye: "再见"
```

#### en-US.yaml

```yaml
errors:
  invalid_params: "Invalid parameters"
  user_not_found: "User not found"
  duplicate_username: "Username already exists"

messages:
  welcome: "Welcome to Rei0721"
  goodbye: "Goodbye"
```

## 最佳实践

### 1. 使用环境变量

敏感信息应该使用环境变量，而不是硬编码在配置文件中：

```yaml
# ✅ 好
database:
  password: ${DB_PASSWORD}

# ❌ 不好
database:
  password: "my_secret_password"
```

### 2. 提供默认值

为环境变量提供合理的默认值：

```yaml
# ✅ 好
port: ${SERVER_PORT:8080}

# ❌ 不好
port: ${SERVER_PORT}
```

### 3. 验证配置

在应用启动时验证配置的有效性：

```go
if err := cfg.Validate(); err != nil {
    log.Fatal("invalid configuration:", err)
}
```

### 4. 使用 Hook 处理变更

注册 Hook 以响应配置变更：

```go
configManager.RegisterHook(func(old, new *config.Config) {
    if old.Logger.Level != new.Logger.Level {
        // 更新日志级别
        updateLogLevel(new.Logger.Level)
    }
})
```

### 5. 原子更新

使用 `Update` 方法进行原子更新，避免部分更新：

```go
// ✅ 好 - 原子更新
configManager.Update(func(cfg *config.Config) {
    cfg.Server.Port = 8081
    cfg.Logger.Level = "debug"
})

// ❌ 不好 - 非原子更新
cfg := configManager.Get()
cfg.Server.Port = 8081
cfg.Logger.Level = "debug"
```

### 6. 只读访问

通过 `Get()` 获取的配置应该视为只读：

```go
// ✅ 好
cfg := configManager.Get()
port := cfg.Server.Port

// ❌ 不好
cfg := configManager.Get()
cfg.Server.Port = 8081  // 不会生效
```

## 故障排查

### 配置文件找不到

```bash
# 检查文件是否存在
ls -la configs/config.yaml

# 指定完整路径
CONFIG_PATH=/path/to/config.yaml go run ./cmd/server/main.go
```

### 环境变量未生效

```bash
# 检查环境变量是否设置
echo $SERVER_PORT

# 重新设置环境变量
export SERVER_PORT=8081

# 验证配置
go run ./cmd/server/main.go
```

### 配置验证失败

```bash
# 检查配置文件格式
cat configs/config.yaml

# 查看详细错误信息
LOG_LEVEL=debug go run ./cmd/server/main.go
```

---

**最后更新**: 2025-12-30
