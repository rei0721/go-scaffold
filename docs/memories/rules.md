# 项目编码规则 (Coding Rules)

> **项目特定的 AI 行为准则**  
> 本文档定义项目的编码规范和协作约定。

---

## 1. 代码风格约定

### 1.1 Go 语言规范

遵循官方 Go 编码规范：

- **格式化**: 所有代码必须通过 `gofmt` 格式化
- **Lint**: 必须通过 `golint` 检查（或 `golangci-lint`）
- **导入分组**: 标准库 → 第三方库 → 本项目包

**示例**:

```go
import (
    "context"
    "fmt"

    "github.com/gin-gonic/gin"

    "github.com/rei0721/rei0721/internal/config"
    "github.com/rei0721/rei0721/pkg/logger"
)
```

### 1.2 命名规范

| 类型     | 规范                     | 示例                             |
| -------- | ------------------------ | -------------------------------- |
| 包名     | 小写单词，简短           | `cache`, `logger`, `executor`    |
| 导出     | PascalCase（首字母大写） | `NewApp`, `Config`, `Execute`    |
| 私有     | camelCase（首字母小写）  | `initLogger`, `parseConfig`      |
| 常量     | PascalCase               | `DefaultPoolSize`, `ExitSuccess` |
| 接口     | 单一方法用 `-er` 后缀    | `Logger`, `Reloader`, `Database` |
| 环境变量 | UPPER_SNAKE_CASE         | `DB_PASSWORD`, `REDIS_HOST`      |

### 1.3 注释规范

**必须注释**：

- 所有导出符号（变量、函数、类型）
- 包级文档（`doc.go`）
- 复杂逻辑（解释"为什么"而非"是什么"）

**注释格式**：

```go
// NewApp 创建一个新的 App 实例
// 按照正确的依赖顺序初始化所有组件
// 参数:
//   opts: 应用选项
// 返回:
//   *App: 初始化完成的应用实例
//   error: 初始化失败时的错误
func NewApp(opts Options) (*App, error) {
    // ...
}
```

---

## 2. 架构约束

### 2.1 分层依赖规则

```
cmd/      → internal/, pkg/     [可依赖业务层和基础层]
internal/ → pkg/, types/        [可依赖基础层]
pkg/      → 第三方库 ONLY       [禁止依赖 internal/]
types/    → 无依赖              [纯数据类型]
```

**违规示例**（禁止）：

```go
// pkg/cache/redis.go
import "github.com/rei0721/rei0721/internal/config"  // ❌ pkg 不能依赖 internal
```

**正确做法**：

```go
// pkg/cache/redis.go
// 通过参数传递配置，而非直接导入 internal/config
func NewRedis(cfg *Config, logger logger.Logger) (Cache, error) {
    // ...
}
```

### 2.2 Context 传递规则

**强制要求**：

- 所有 I/O 操作（HTTP、数据库、Redis、RPC）必须透传 `context.Context`
- 第一个参数必须是 `ctx context.Context`
- 禁止使用 `context.Background()`，除非在顶层入口（如 `main.go`）

**示例**：

```go
// ✅ 正确
func (r *userRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
    return r.db.GetDB().WithContext(ctx).First(&user, id).Error
}

// ❌ 错误
func (r *userRepository) GetUserByID(id int64) (*models.User, error) {
    return r.db.GetDB().First(&user, id).Error  // 缺少 context
}
```

### 2.3 错误处理规则

**禁止忽略错误**：

```go
// ❌ 绝对禁止
_ = db.Close()
_ = logger.Sync()

// ✅ 正确
if err := db.Close(); err != nil {
    logger.Error("failed to close database", "error", err)
}
```

**错误包装**：
使用 `fmt.Errorf` 和 `%w` 包装错误

```go
if err := someOperation(); err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

---

## 3. 依赖注入 (DI) 规范

### 3.1 容器模式

使用 `internal/app.App` 作为 DI 容器：

```go
type App struct {
    Config   *config.Config
    Logger   logger.Logger    // 接口，非具体实现
    DB       database.Database
    Cache    cache.Cache
    Executor executor.Manager
    // ...
}
```

### 3.2 接口优先

**基础设施组件必须定义接口**：

```go
// pkg/logger/logger.go
type Logger interface {
    Debug(msg string, fields ...any)
    Info(msg string, fields ...any)
    // ...
}

// 实现
type zapLogger struct { /* ... */ }
func (z *zapLogger) Debug(msg string, fields ...any) { /* ... */ }
```

**好处**：

- 便于测试（Mock）
- 便于替换实现
- 降低耦合

---

## 4. 配置管理规范

### 4.1 配置来源优先级

```
环境变量 > .env 文件 > config.yaml
```

### 4.2 敏感信息处理

**禁止硬编码**：

- 数据库密码
- Redis 密码
- API 密钥
- Token

**正确做法**：

```yaml
# config.yaml
database:
  password: "" # 留空，通过环境变量覆盖
```

```bash
# .env (不提交到 Git)
DB_PASSWORD=your-secure-password
```

### 4.3 配置验证

所有配置结构体必须实现 `Validate()` 方法：

```go
type Configurable interface {
    Validate() error
}

func (c *DatabaseConfig) Validate() error {
    if c.Host == "" {
        return errors.New("host is required")
    }
    // ...
}
```

---

## 5. 热重载规范

### 5.1 Reloader 接口

支持热重载的组件必须实现：

```go
type Reloader interface {
    Reload(ctx context.Context, newConfig interface{}) error
}
```

### 5.2 原子替换模式

使用 `sync.RWMutex` 保护并发访问：

```go
type redisCache struct {
    client *redis.Client
    mu     sync.RWMutex  // 保护 client 并发访问
}

func (c *redisCache) Reload(ctx context.Context, cfg *Config) error {
    newClient := createClient(cfg)

    c.mu.Lock()
    oldClient := c.client
    c.client = newClient  // 原子替换
    c.mu.Unlock()

    oldClient.Close()
    return nil
}

func (c *redisCache) Get(ctx context.Context, key string) (string, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.client.Get(ctx, key).Result()
}
```

---

## 6. 测试规范

### 6.1 单元测试

- 所有 `pkg/` 包必须有单元测试
- 测试文件命名：`*_test.go`
- 测试函数命名：`TestXxx`

### 6.2 表驱动测试

推荐使用表驱动测试（Table-Driven Tests）：

```go
func TestParseLevel(t *testing.T) {
    tests := []struct {
        name string
        input string
        want Level
    }{
        {"debug level", "debug", LevelDebug},
        {"info level", "info", LevelInfo},
        {"unknown defaults to info", "unknown", LevelInfo},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := parseLevel(tt.input)
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 6.3 Mock 策略

使用接口 Mock：

- 手写 Mock（简单场景）
- 使用 `go.uber.org/mock`（复杂场景）

---

## 7. 包文档规范

### 7.1 每个包必须有 `doc.go`

```go
/*
Package cache 提供统一的缓存操作接口

# 设计目标

- 统一接口
- 易于使用
- 线程安全

# 使用示例

    cache, err := cache.NewRedis(cfg, logger)
    cache.Set(ctx, "key", "value", time.Hour)

# 最佳实践

1. 始终设置过期时间
2. 使用命名空间前缀
3. 区分键不存在和其他错误
*/
package cache
```

### 7.2 文档结构

推荐章节：

1. 概述
2. 设计目标/核心概念
3. 使用示例
4. 最佳实践
5. 线程安全说明

---

## 8. Git 提交规范

### 8.1 Commit Message 格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Type**:

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建/工具相关

**示例**:

```
feat(executor): add atomic hot-reload support

- Implement Reloader interface
- Add RWMutex for concurrent safety
- Update documentation

Closes #123
```

---

## 9. 常见反模式 (Anti-Patterns)

### ❌ 禁止事项

1. **Magic Numbers/Strings**: 使用常量代替魔法值

```go
// ❌ 错误
if mode == "release" { }

// ✅ 正确
const ModeRelease = "release"
if mode == ModeRelease { }
```

2. **全局可变状态**: 避免全局变量（常量除外）

```go
// ❌ 错误
var db *gorm.DB  // 全局可变

// ✅ 正确
type App struct {
    DB database.Database  // 通过 DI 容器管理
}
```

3. **Init 副作用**: `init()` 函数不应有副作用

```go
// ❌ 错误
func init() {
    db, _ = connectDatabase()  // 副作用
}

// ✅ 正确
func New(cfg *Config) (*App, error) {
    db, err := connectDatabase(cfg)
    // ...
}
```

---

## 10. IDE 配置建议

### VS Code 推荐设置

```json
{
  "go.formatTool": "gofmt",
  "go.lintOnSave": "package",
  "editor.formatOnSave": true,
  "[go]": {
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  }
}
```

---

> **提醒**: 违反本文档规范的代码将在 Code Review 中被要求修改。
