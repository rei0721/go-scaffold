# 项目结构

本文档详细说明了 Rei0721 项目的目录结构和文件组织。

## 目录树

```
rei0721/
├── cmd/                          # 应用入口
│   └── server/
│       └── main.go              # 服务器启动入口
│
├── internal/                     # 内部包 (不对外暴露)
│   ├── app/
│   │   └── app.go               # IoC 容器和生命周期管理
│   │
│   ├── config/
│   │   ├── config.go            # 配置结构定义
│   │   └── manager.go           # 配置管理器 (热重载)
│   │
│   ├── handler/
│   │   └── user.go              # HTTP 请求处理器
│   │
│   ├── middleware/
│   │   ├── config.go            # 中间件配置
│   │   ├── logger.go            # 日志中间件
│   │   ├── recovery.go          # 恢复中间件 (panic 处理)
│   │   └── traceid.go           # 追踪 ID 中间件
│   │
│   ├── models/
│   │   ├── base.go              # 基础模型 (ID, 时间戳)
│   │   └── user.go              # 用户数据模型
│   │
│   ├── repository/
│   │   ├── repository.go        # 泛型仓储接口
│   │   └── user.go              # 用户仓储实现
│   │
│   ├── router/
│   │   ├── router.go            # 路由定义
│   │   └── router_test.go       # 路由测试
│   │
│   └── service/
│       ├── user.go              # 用户服务接口
│       └── user_impl.go         # 用户服务实现
│
├── pkg/                          # 可复用工具包 (无 internal 依赖)
│   ├── database/
│   │   ├── database.go          # 数据库接口定义
│   │   └── factory.go           # 数据库工厂函数
│   │
│   ├── id/
│   │   └── snowflake.go         # Snowflake ID 生成器
│   │
│   ├── logger/
│   │   ├── logger.go            # 日志接口定义
│   │   └── zap.go               # Zap 日志实现
│   │
│   └── scheduler/
│       ├── scheduler.go         # 调度器接口定义
│       └── ants.go              # ants 协程池实现
│
├── types/                        # 类型定义 (零依赖)
│   ├── errors/
│   │   ├── codes.go             # 错误码常量
│   │   └── error.go             # 业务错误类型
│   │
│   ├── result/
│   │   ├── result.go            # 统一响应类型
│   │   └── pagination.go        # 分页类型
│   │
│   ├── request.go               # 请求类型
│   └── response.go              # 响应类型
│
├── configs/                      # 配置文件
│   ├── config.yaml              # 主配置文件
│   ├── .env.example             # 环境变量模板
│   └── i18n/
│       ├── zh-CN.yaml           # 中文消息
│       └── en-US.yaml           # 英文消息
│
├── docs/                         # 项目文档
│   ├── README.md                # 文档首页
│   ├── getting-started.md       # 快速开始
│   ├── architecture.md          # 架构设计
│   ├── project-structure.md     # 项目结构 (本文件)
│   ├── configuration.md         # 配置管理
│   ├── database.md              # 数据库
│   ├── middleware.md            # 中间件
│   ├── api.md                   # API 规范
│   ├── protocol.md              # 开发规范
│   ├── deployment.md            # 部署指南
│   └── faq.md                   # 常见问题
│
├── logs/                         # 日志目录
│   └── .gitkeep
│
├── .gitignore                    # Git 忽略文件
├── go.mod                        # Go 模块定义
├── go.sum                        # Go 依赖校验
├── Dockerfile                    # Docker 镜像定义
├── docker-compose.yml            # Docker Compose 配置
├── Makefile                      # 构建脚本
└── README.md                     # 项目说明
```

## 目录说明

### cmd/ - 应用入口

存放应用的启动入口代码。

```
cmd/
└── server/
    └── main.go          # 服务器启动入口
```

**职责**:
- 初始化应用容器
- 处理 OS 信号
- 优雅关闭

**特点**:
- 最小化业务逻辑
- 清晰的错误处理
- 完整的生命周期管理

### internal/ - 内部包

存放项目内部的业务代码，不对外暴露。

#### internal/app/

**职责**: 依赖注入容器和生命周期管理

```go
type App struct {
    Config        *config.Config
    ConfigManager config.Manager
    DB            database.Database
    Scheduler     scheduler.Scheduler
    Logger        logger.Logger
    Router        *gin.Engine
}
```

#### internal/config/

**职责**: 配置管理

- `config.go` - 配置结构定义
- `manager.go` - 配置管理器 (支持热重载)

#### internal/handler/

**职责**: HTTP 请求处理

- 参数解析和验证
- 调用服务
- 错误处理
- 响应构造

#### internal/middleware/

**职责**: 横切关注点处理

- `traceid.go` - 生成/提取 TraceID
- `logger.go` - 记录请求信息
- `recovery.go` - 捕获 panic

#### internal/models/

**职责**: 数据模型定义

- `base.go` - 基础模型 (ID, 时间戳)
- `user.go` - 用户模型

#### internal/repository/

**职责**: 数据访问抽象

- `repository.go` - 泛型仓储接口
- `user.go` - 用户仓储实现

#### internal/router/

**职责**: 路由定义

- 路由分组
- 中间件配置
- 端点注册

#### internal/service/

**职责**: 业务逻辑处理

- `user.go` - 用户服务接口
- `user_impl.go` - 用户服务实现

### pkg/ - 可复用工具包

存放可复用的工具库，不依赖 internal 包。

#### pkg/database/

**职责**: 数据库抽象

- 支持 PostgreSQL, MySQL, SQLite
- 连接池管理
- Hook 扩展点

#### pkg/id/

**职责**: ID 生成

- Snowflake 算法
- 分布式 ID 生成

#### pkg/logger/

**职责**: 日志抽象

- 统一日志接口
- Zap 实现
- 结构化日志

#### pkg/scheduler/

**职责**: 协程调度

- ants 协程池
- 统一任务提交接口
- 优雅关闭

### types/ - 类型定义

存放跨层的类型定义，零依赖。

#### types/errors/

**职责**: 错误定义

- `codes.go` - 错误码常量
- `error.go` - 业务错误类型

#### types/result/

**职责**: 响应类型

- `result.go` - 统一响应结构
- `pagination.go` - 分页类型

#### types/request.go

**职责**: 请求类型

```go
type RegisterRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

#### types/response.go

**职责**: 响应类型

```go
type UserResponse struct {
    UserID    int64
    Username  string
    Email     string
    Status    int
    CreatedAt time.Time
}
```

### configs/ - 配置文件

#### config.yaml

主配置文件，支持环境变量替换。

```yaml
server:
  port: ${SERVER_PORT:8080}
  mode: ${SERVER_MODE:debug}

database:
  driver: ${DB_DRIVER:postgres}
  host: ${DB_HOST:localhost}
  port: ${DB_PORT:5432}
```

#### .env.example

环境变量模板，用于快速配置。

```bash
SERVER_PORT=8080
SERVER_MODE=debug
DB_DRIVER=postgres
DB_HOST=localhost
```

#### i18n/

国际化消息文件。

- `zh-CN.yaml` - 中文消息
- `en-US.yaml` - 英文消息

### docs/ - 项目文档

| 文件 | 说明 |
|------|------|
| README.md | 文档首页 |
| getting-started.md | 快速开始指南 |
| architecture.md | 架构设计文档 |
| project-structure.md | 项目结构说明 (本文件) |
| configuration.md | 配置管理指南 |
| database.md | 数据库指南 |
| middleware.md | 中间件说明 |
| api.md | API 规范 |
| protocol.md | 开发规范 |
| deployment.md | 部署指南 |
| faq.md | 常见问题 |

### logs/ - 日志目录

存放应用运行时生成的日志文件。

## 依赖关系

### 允许的导入

```go
// ✅ 允许
import "rei0721/types"
import "rei0721/pkg/logger"
import "rei0721/internal/service"

// ❌ 禁止
import "rei0721/internal" // 在 pkg 中
import "rei0721/pkg"      // 在 types 中
```

### 分层依赖

```
cmd
 ↓
internal/app
 ├→ internal/router
 ├→ internal/handler
 ├→ internal/service
 ├→ internal/repository
 ├→ internal/models
 ├→ internal/config
 ├→ internal/middleware
 ├→ pkg/*
 └→ types/*

pkg/*
 └→ types/*

types/*
 └→ (无依赖)
```

## 文件命名规范

### Go 源文件

- **接口定义**: `{name}.go`
  - 例: `user.go` (UserService 接口)

- **实现文件**: `{name}_impl.go`
  - 例: `user_impl.go` (UserService 实现)

- **测试文件**: `{name}_test.go`
  - 例: `user_test.go`

- **工厂函数**: `factory.go`
  - 例: `pkg/database/factory.go`

### 配置文件

- **主配置**: `config.yaml`
- **环境模板**: `.env.example`
- **国际化**: `{lang}.yaml`
  - 例: `zh-CN.yaml`, `en-US.yaml`

## 包导入规范

### 导入顺序

```go
import (
    // 标准库
    "context"
    "fmt"
    "time"

    // 第三方库
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    // 项目内部
    "rei0721/internal/models"
    "rei0721/pkg/logger"
    "rei0721/types"
)
```

### 导入别名

```go
import (
    "rei0721/types/errors"
    "rei0721/types/result"
    
    bizErr "rei0721/types/errors"  // 避免冲突
)
```

## 最佳实践

### 1. 单一职责

每个包只负责一个功能：

```
✅ pkg/logger - 只处理日志
✅ pkg/database - 只处理数据库
❌ pkg/utils - 混合多个功能
```

### 2. 接口驱动

定义接口，而不是具体实现：

```go
// ✅ 好
type Logger interface {
    Info(msg string, keysAndValues ...interface{})
}

// ❌ 不好
type ZapLogger struct {
    // 具体实现
}
```

### 3. 依赖注入

通过构造函数注入依赖：

```go
// ✅ 好
func NewUserService(repo UserRepository, sched Scheduler) UserService {
    return &userService{repo: repo, sched: sched}
}

// ❌ 不好
func NewUserService() UserService {
    repo := NewUserRepository()
    sched := NewScheduler()
    return &userService{repo: repo, sched: sched}
}
```

### 4. 错误处理

使用统一的错误类型：

```go
// ✅ 好
return nil, &errors.BizError{
    Code:    errors.ErrUserNotFound,
    Message: "user not found",
}

// ❌ 不好
return nil, fmt.Errorf("user not found")
```

### 5. 上下文传递

所有 I/O 操作都应该接收 context：

```go
// ✅ 好
func (r *userRepository) FindByID(ctx context.Context, id int64) (*User, error) {
    return r.db.WithContext(ctx).First(&User{}, id).Error
}

// ❌ 不好
func (r *userRepository) FindByID(id int64) (*User, error) {
    return r.db.First(&User{}, id).Error
}
```

---

**最后更新**: 2025-12-30
