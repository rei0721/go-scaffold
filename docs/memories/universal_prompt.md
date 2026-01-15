# 通用记忆提示词 (Universal Memory Prompt)

> **可复用的规则与模式**  
> 本文档沉淀可跨任务复用的协作规范和项目特有经验。

---

## 1. 脚手架吸收规则 (Scaffold Assimilation)

### 1.1 本项目工程化能力盘点

本项目具有以下**成熟工程化特性**：

| 类别            | 能力          | 实现方式                        |
| --------------- | ------------- | ------------------------------- |
| **架构模式**    | 分层架构 + DI | `internal/app.App` 容器         |
| **依赖管理**    | Go Modules    | `go.mod`                        |
| **配置管理**    | 多来源配置    | Viper + 环境变量 + .env         |
| **热重载**      | 配置热更新    | `Reloader` 接口                 |
| **日志系统**    | 结构化日志    | Zap + Lumberjack                |
| **包文档**      | 规范化文档    | 每个包有 `doc.go`               |
| **错误处理**    | 统一错误码    | 常量化错误消息                  |
| **并发控制**    | 协程池管理    | ants + 原子热重载               |
| **缓存抽象**    | Redis 接口化  | `cache.Cache` 接口              |
| **数据库抽象**  | ORM 接口化    | `database.Database` 接口        |
| **国际化**      | 多语言支持    | go-i18n v2                      |
| **HTTP 服务器** | 封装 Gin      | `pkg/httpserver` 热重载         |
| **SQL 生成**    | 双向代码生成  | `pkg/sqlgen` Model ↔ SQL        |
| **延迟注入**    | Executor 注入 | `types.ExecutorInjectable` 接口 |
| **延迟注入**    | Cache 注入    | `types.CacheInjectable` 接口    |
| **延迟注入**    | Logger 注入   | Service 层日志记录              |
| **缓存集成**    | 服务层缓存    | Cache-Aside 模式                |
| **启动模式**    | 多模式启动    | `server` / `initdb` 模式        |

### 1.2 与协议冲突的解法

**冲突点**: IEP v7.2 协议示例代码树与项目现有结构不同

**解决方案**:

- ✅ 保留项目现有目录结构（`cmd/`, `internal/`, `pkg/`, `types/`, `configs/`）
- ✅ 仅补充协议要求的 `docs/` 和 `specs/` 目录
- ✅ 协议示例仅作参考，实际以项目为准

### 1.3 可复用优点固化

**优点 1**: 依赖注入容器模式  
**规则**: 所有基础设施组件通过 `App` 统一管理，按依赖顺序初始化

**优点 2**: 接口抽象 + 实现分离  
**规则**: `pkg/` 包定义接口，便于测试和替换实现

**优点 3**: 原子热重载设计  
**规则**: 使用 `sync.RWMutex` + 原子替换模式，确保线程安全

**优点 4**: 包文档规范 (`doc.go`)  
**规则**: 每个包必须有详细的使用文档、设计目标、示例代码

---

## 2. 依赖注入模式固化

### 2.1 初始化顺序宪法

**强制顺序**：

```
Config → Logger → I18n → Cache → Database → Executor → Business
```

**理由**：

1. `Config` 最先：其他组件需要配置
2. `Logger` 第二：后续初始化需要日志记录
3. `I18n` 第三：HTTP 响应需要国际化
4. `Cache/Database` 第四：业务层依赖数据存储
5. `Executor` 第五：业务层可能需要异步任务
6. `Business` 最后：完整的请求处理链路

### 2.2 优雅关闭顺序

**强制顺序**：

```
HTTP Server → Executor → Cache → Database → Logger
```

**理由**：

1. 先停止接收新请求
2. 等待正在处理的任务完成
3. 关闭外部连接（Redis、DB）
4. 最后刷新日志缓冲区

---

## 3. 热重载设计模式固化

### 3.1 Reloader 接口契约

```go
type Reloader interface {
    Reload(ctx context.Context, newConfig interface{}) error
}
```

**实现要点**：

1. 创建新实例（在锁外）
2. 加写锁
3. 原子替换
4. 释放写锁
5. 关闭旧实例

### 3.2 并发访问模式

**读操作**（高频）：

```go
func (c *Component) Operation() {
    c.mu.RLock()
    defer c.mu.RUnlock()
    // 使用 c.resource
}
```

**写操作**（热重载）：

```go
func (c *Component) Reload(newConfig) error {
    newResource := create(newConfig)

    c.mu.Lock()
    oldResource := c.resource
    c.resource = newResource  // 原子
    c.mu.Unlock()

    oldResource.Close()
}
```

---

## 4. 包文档规范固化

### 4.1 doc.go 模板

```go
/*
Package <name> <一句话描述>

# 设计目标

- 目标1
- 目标2

# 核心概念

概念解释...

# 使用示例

基本用法:

    example code

高级用法:

    example code

# 最佳实践

1. 实践1
2. 实践2

# 线程安全

所有方法都是线程安全的...

# 与其他包的区别

pkg/a: 用途A
pkg/b: 用途B
*/
package name
```

### 4.2 注释规范

**导出符号必须注释**：

```go
// NewClient 创建一个新的客户端实例
// 参数:
//   cfg: 配置对象
//   logger: 日志器
// 返回:
//   Client: 客户端实例
//   error: 创建失败时的错误
func NewClient(cfg *Config, logger logger.Logger) (Client, error) {
    // ...
}
```

---

## 5. 常量管理规范固化

### 5.1 常量定义位置约定

| 常量类型          | 定义位置                       | 示例                    |
| ----------------- | ------------------------------ | ----------------------- |
| **包级基础常量**  | `pkg/*/constants.go`           | 默认值、错误消息        |
| **池名称/资源名** | `types/constants/executor.go`  | `PoolHTTP`, `PoolCache` |
| **应用常量**      | `types/constants/app.go`       | 应用级别常量            |
| **环境变量名**    | `internal/config/constants.go` | `EnvDBPassword`         |
| **应用层常量**    | `internal/app/constants.go`    | 应用特定常量            |
| **公共接口定义**  | `types/interfaces.go`          | `ExecutorInjectable`    |

### 5.2 命名前缀约定

见 [`docs/architecture/variable_index.md`](/docs/architecture/variable_index.md#62-前缀约定)

---

## 6. 错误处理模式固化

### 6.1 错误定义模式

**消息模板**（在 `constants.go`）：

```go
const (
    ErrMsgOperationFailed = "operation failed: %w"
    ErrMsgInvalidConfig = "invalid config: %w"
```

**预定义错误**（sentinel errors）：

```go
var (
    ErrNotFound = errors.New("not found")
    ErrInvalidInput = errors.New("invalid input")
)
```

### 6.2 错误包装规范

**必须使用 `%w`**：

```go
if err := doSomething(); err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

**错误判断**：

```go
if errors.Is(err, ErrNotFound) {
    // 处理 NotFound
}

var configErr *ConfigError
if errors.As(err, &configErr) {
    // 处理配置错误
}
```

---

## 7. Context 传递模式固化

### 7.1 强制规则

**所有 I/O 操作必须透传 Context**：

- HTTP 请求
- 数据库查询
- Redis 操作
- gRPC 调用

**第一个参数必须是 `ctx`**：

```go
func (r *Repository) GetUser(ctx context.Context, id int64) (*User, error) {
    return r.db.GetDB().WithContext(ctx).First(&user, id).Error
}
```

### 7.2 Context 传递链路

```
HTTP Request → Handler → Service → Repository → Database
     ctx ----------> ctx -----> ctx --------> ctx
```

---

## 8. 测试模式固化

### 8.1 表驱动测试 (Table-Driven Tests)

所有需要多 case 测试的函数都应使用表驱动：

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name string
        input Type
        want Type
        wantErr bool
    }{
        {"case 1", input1, want1, false},
        {"case 2", input2, want2, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("unexpected error: %v", err)
            }
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 8.2 Mock 策略

**优先手写 Mock**（简单场景）：

```go
type mockLogger struct{}
func (m *mockLogger) Info(msg string, fields ...any) {}
// ...
```

**复杂场景使用第三方库**：

- `go.uber.org/mock`

---

## 9. 配置管理模式固化

### 9.1 三层配置优先级

```
1. 环境变量（最高优先级）
2. .env 文件
3. config.yaml（默认值）
```

### 9.2 敏感信息处理

**禁止硬编码**：

- ❌ `password: "secret"` in `config.yaml`
- ✅ `password: ""` in `config.yaml` + `DB_PASSWORD=xxx` in `.env`

**Git 忽略**：

```gitignore
.env
config.yaml  # 如果包含敏感信息
```

**提供示例文件**：

- `config.example.yaml`
- `.env.example`

---

## 10. 项目特有经验固化

### 10.1 Gin 模式切换

| 模式      | 用途 | 特点                 |
| --------- | ---- | -------------------- |
| `debug`   | 开发 | 详细日志、panic 堆栈 |
| `release` | 生产 | 性能优化、简化日志   |
| `test`    | 测试 | 测试专用             |

### 10.2 数据库驱动选择策略

| 驱动       | 场景       | 理由                |
| ---------- | ---------- | ------------------- |
| SQLite     | 开发/测试  | 无需安装、快速启动  |
| PostgreSQL | 生产       | 功能强大、ACID 保证 |
| MySQL      | 高性能场景 | 读写分离友好        |

### 10.3 日志输出策略

| 环境     | 推荐输出             | 推荐格式  |
| -------- | -------------------- | --------- |
| 开发     | `both` (控制台+文件) | `console` |
| 容器/K8s | `stdout`             | `json`    |
| 传统部署 | `file`               | `json`    |

### 10.4 启动模式

| 模式     | 用途         | 说明                       |
| -------- | ------------ | -------------------------- |
| `server` | 生产/开发    | 完整启动，包含 HTTP 服务器 |
| `initdb` | 数据库初始化 | 仅初始化数据库表结构后退出 |

---

## 11. 延迟注入模式固化

### 11.1 ExecutorInjectable 接口

**目的**: 解除组件初始化顺序依赖，支持运行时动态替换 executor

**定义位置**: `types/interfaces.go`

```go
type ExecutorInjectable interface {
    SetExecutor(exec executor.Manager)
}
```

**实现模式**:

```go
type Component struct {
    executor atomic.Value // 存储 executor.Manager
}

func (c *Component) SetExecutor(exec executor.Manager) {
    c.executor.Store(exec)
}

func (c *Component) getExecutor() executor.Manager {
    if exec := c.executor.Load(); exec != nil {
        return exec.(executor.Manager)
    }
    return nil
}
```

**使用场景**:

- 需要异步处理但初始化时 executor 尚未创建的组件
- 测试时可选不注入 executor
- 支持运行时动态替换

---

## 13. pkg 通用工具库标准规范

### 13.1 强制文件清单

每个 `pkg/*` 工具库**必须**包含以下标准文件：

| 文件           | 必须性  | 说明                                           |
| -------------- | ------- | ---------------------------------------------- |
| `doc.go`       | ✅ 必须 | 包级文档，详细说明设计目标、核心概念、使用示例 |
| `README.md`    | ✅ 必须 | 快速开始指南，概览和基本用法                   |
| `errors.go`    | ✅ 必须 | 错误定义（Sentinel Errors + 错误消息模板）     |
| `constants.go` | 可选    | 包内常量，如默认值、缓存键、TTL 等             |
| `*_test.go`    | ✅ 必须 | 单元测试，覆盖核心功能                         |
| `examples/`    | ✅ 必须 | 可运行的示例代码                               |

### 13.2 标准目录结构

```
pkg/yourlib/
├── doc.go              # 包文档（GoDoc）
├── README.md           # 快速开始
├── errors.go           # 错误定义
├── constants.go        # 常量定义（可选）
├── yourlib.go          # 接口定义
├── yourlib_impl.go     # 实现（如有必要）
└── examples/
    ├── README.md       # 示例说明
    └── basic/
        └── main.go     # 基础示例
```

**如果包较复杂，可以分层**：

```
pkg/yourlib/
├── doc.go
├── README.md
├── errors.go
├── constants.go
├── models/             # 数据模型
│   └── models.go
├── repository/         # 数据访问层（如适用）
│   ├── interface.go
│   └── impl.go
├── service/            # 业务逻辑层（如适用）
│   ├── interface.go
│   ├── service.go
│   └── service_test.go
└── examples/
    └── basic/
        └── main.go
```

### 13.3 doc.go 标准模板

```go
/*
Package yourlib 一句话描述包的功能

# 设计目标

- 目标1：xxx
- 目标2：xxx
- 目标3：xxx

# 核心概念

概念解释...

# 使用示例

基本用法:

    // 示例代码
    lib := yourlib.New(config)
    result, err := lib.DoSomething()

高级用法:

    // 高级示例代码

# 最佳实践

1. 实践1
2. 实践2

# 线程安全

所有方法都是线程安全的...

# 与其他包的区别

pkg/a: 用途A
pkg/b: 用途B
*/
package yourlib
```

### 13.4 errors.go 标准模板

```go
package yourlib

import "errors"

// 预定义错误（Sentinel Errors）
// 可使用 errors.Is() 判断
var (
    ErrNotFound = errors.New("resource not found")
    ErrInvalidInput = errors.New("invalid input")
    // ...
)

// 错误消息模板常量
// 用于 fmt.Errorf() 包装错误
const (
    ErrMsgOperationFailed = "operation failed: %w"
    ErrMsgInvalidConfig = "invalid config: %w"
    // ...
)
```

**使用方式**：

```go
// 返回预定义错误
if resource == nil {
    return ErrNotFound
}

// 包装错误保留上下文
if err := doSomething(); err != nil {
    return fmt.Errorf(ErrMsgOperationFailed, err)
}

// 判断错误
if errors.Is(err, ErrNotFound) {
    // 处理未找到
}
```

### 13.5 README.md 标准模板

```markdown
# pkg/yourlib - 一句话说明

简短描述，一两句话。

## 特性

- ✅ 特性 1
- ✅ 特性 2
- ✅ 特性 3

## 快速开始

### 1. 安装

\`\`\`bash
go get github.com/yourorg/yourlib
\`\`\`

### 2. 基本使用

\`\`\`go
// 示例代码
\`\`\`

## API 参考

- [接口文档](interface.go)
- [包文档](doc.go)

## 最佳实践

1. 实践 1
2. 实践 2

## 依赖项

### 必须依赖

- xxx

### 可选依赖

- xxx
```

### 13.6 单元测试规范

**测试文件位置**：与被测试文件同目录，命名为 `*_test.go`

**测试函数命名**：`TestXxx`

**表驱动测试**（推荐）：

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   Type
        want    Type
        wantErr bool
    }{
        {"case1", input1, want1, false},
        {"case2", input2, want2, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("unexpected error: %v", err)
            }
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

**Mock 策略**：

- 简单场景：手写 Mock
- 复杂场景：使用 `go.uber.org/mock`

### 13.7 示例代码规范

**目录结构**：

```
examples/
├── README.md           # 示例说明
├── basic/              # 基础示例
│   └── main.go
├── advanced/           # 高级示例（可选）
│   └── main.go
└── integration/        # 集成示例（可选）
    └── main.go
```

**示例代码要求**：

1. ✅ 可以直接运行（`go run main.go`）
2. ✅ 有详细注释
3. ✅ 展示常见用法
4. ✅ 包含错误处理

**examples/README.md 模板**：

```markdown
# 示例

## 运行示例

### 基础示例

\`\`\`bash
cd examples/basic
go run main.go
\`\`\`

**示例输出**：
\`\`\`
预期的输出...
\`\`\`

## 示例说明

### basic 示例

- 演示内容 1
- 演示内容 2
```

### 13.8 接口设计原则

**接口命名**：

```go
// 接口名通常与包名相同或相近
type Cache interface { ... }        // pkg/cache
type Logger interface { ... }       // pkg/logger
type Database interface { ... }     // pkg/database
```

**接口职责**：

- ✅ 单一职责原则
- ✅ 方法数量适中（5-10 个）
- ✅ 避免过度抽象

**实现命名**：

```go
// 实现类型通常带具体技术名
type redisCache struct { ... }      // Redis 实现
type zapLogger struct { ... }       // Zap 实现
type gormDatabase struct { ... }    // GORM 实现
```

### 13.9 配置模式

**使用 Option 模式**（推荐）：

```go
type Config struct {
    Host string
    Port int
    // ...
}

type Option func(*Config)

func WithHost(host string) Option {
    return func(c *Config) {
        c.Host = host
    }
}

func New(opts ...Option) Library {
    cfg := &Config{
        Host: "localhost", // 默认值
        Port: 8080,
    }
    for _, opt := range opts {
        opt(cfg)
    }
    return &libraryImpl{config: cfg}
}

// 使用
lib := New(
    WithHost("example.com"),
    WithPort(9000),
)
```

### 13.10 依赖注入模式

**延迟注入**（用于可选依赖）：

```go
type Library interface {
    types.ExecutorInjectable  // 嵌入接口
    types.CacheInjectable
    // ...
}

type libraryImpl struct {
    executor atomic.Value  // 存储 executor.Manager
    cache    atomic.Value  // 存储 cache.Cache
}

func (l *libraryImpl) SetExecutor(exec executor.Manager) {
    l.executor.Store(exec)
}

func (l *libraryImpl) getExecutor() executor.Manager {
    if exec := l.executor.Load(); exec != nil {
        return exec.(executor.Manager)
    }
    return nil
}
```

### 13.11 创建检查清单

创建新的 `pkg/*` 工具库时，使用此检查清单：

- [ ] 创建目录结构
- [ ] 编写 `doc.go`（包文档）
- [ ] 编写 `README.md`（快速开始）
- [ ] 定义接口（`yourlib.go`）
- [ ] 实现功能（`yourlib_impl.go`）
- [ ] 定义错误（`errors.go`）
- [ ] 定义常量（`constants.go`，如需要）
- [ ] 编写单元测试（`*_test.go`）
- [ ] 创建示例代码（`examples/basic/main.go`）
- [ ] 编写示例说明（`examples/README.md`）
- [ ] 运行测试确保通过（`go test ./...`）
- [ ] 确保编译通过（`go build ./...`）
- [ ] 更新 `docs/architecture/system_map.md`
- [ ] 更新 `docs/memories/universal_prompt.md` 组件清单

### 13.12 参考示例包

项目中已有的标准示例：

| 包          | 特点         | 可参考内容                            |
| ----------- | ------------ | ------------------------------------- |
| `pkg/rbac`  | 分层结构     | 完整的 models/repository/service 分层 |
| `pkg/cache` | 接口抽象     | 接口定义 + Redis 实现                 |
| `pkg/jwt`   | 简单实用     | 接口 + 实现在同一个包内               |
| `pkg/cli`   | 复杂错误处理 | 自定义错误类型 + ExitCode             |

---

## 14. pkg 组件清单

| 包名             | 职责            | 接口                       |
| ---------------- | --------------- | -------------------------- |
| `pkg/cache`      | Redis 缓存抽象  | `cache.Cache`              |
| `pkg/cli`        | CLI 框架        | `cli.App`, `cli.Command`   |
| `pkg/database`   | 数据库抽象      | `database.Database`        |
| `pkg/executor`   | 协程池管理      | `executor.Manager`         |
| `pkg/httpserver` | HTTP 服务器封装 | `httpserver.HTTPServer`    |
| `pkg/i18n`       | 国际化          | `i18n.I18n`                |
| `pkg/jwt`        | JWT 认证        | `jwt.JWT`                  |
| `pkg/logger`     | 结构化日志      | `logger.Logger`            |
| `pkg/rbac`       | RBAC 权限控制   | `rbac/service.RBACService` |
| `pkg/sqlgen`     | SQL 双向生成    | Model ↔ SQL Script         |
| `pkg/utils`      | 工具函数        | 独立工具函数               |

---

## 更新记录

| 日期       | 内容                                                            |
| ---------- | --------------------------------------------------------------- |
| 2026-01-15 | 固化 pkg 通用工具库标准规范：强制文件、目录结构、模板、检查清单 |
| 2026-01-15 | RBAC 重构：将 RBAC 功能从 internal 抽离到 pkg/rbac 作为通用库   |
| 2026-01-15 | 新增 JWT 认证系统：pkg/jwt 包、认证中间件、路由保护             |
| 2026-01-15 | 用户模块完善：Logger 注入、Update/Delete 功能、缓存失效策略     |
| 2026-01-15 | 新增缓存集成模式：CacheInjectable 接口和 Cache-Aside 模式       |
| 2026-01-15 | 补充 HTTPServer、SQLGen、ExecutorInjectable 模式；更新常量位置  |
| 2026-01-15 | 初始创建，固化项目工程化能力和设计模式                          |

---

> **提醒**: 每次发现可复用的模式后，立即更新本文档！
