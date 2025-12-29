# 架构设计文档

## 概述

Rei0721 是一个遵循分层架构设计的 Go 高性能脚手架项目。本文档详细描述了系统架构、组件设计和依赖关系。

## 架构原则

### 1. 资源受控

所有异步任务必须通过 `Scheduler` 统一管理，而不是直接使用 `go` 关键字。这确保了：

- 协程数量可控
- 资源使用可预测
- 优雅关闭机制完善

```go
// ❌ 禁止
go func() {
    // 异步任务
}()

// ✅ 推荐
scheduler.Submit(ctx, func(ctx context.Context) {
    // 异步任务
})
```

### 2. 协议至上

通过接口驱动设计实现：

- 依赖注入
- 易于测试
- 灵活扩展

```go
// 定义接口
type UserService interface {
    Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error)
    Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
}

// 依赖注入
type UserHandler struct {
    service UserService
}

func NewUserHandler(svc UserService) *UserHandler {
    return &UserHandler{service: svc}
}
```

### 3. 防御性封装

分层隔离，单向依赖：

```
cmd → app → handler → service → repository → models
         ↘ config
         ↘ pkg ✓

pkg → internal ✗ (严禁)
types → * ✗ (无依赖)
```

## 分层架构

### 1. 入口层 (cmd/server)

**职责**: 应用启动和信号处理

```go
// cmd/server/main.go
func main() {
    // 1. 初始化应用容器
    app, err := app.New(app.Options{
        ConfigPath: "configs/config.yaml",
    })
    
    // 2. 启动服务
    go app.Run()
    
    // 3. 监听信号
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    
    // 4. 优雅关闭
    app.Shutdown(ctx)
}
```

**特点**:
- 最小化业务逻辑
- 专注于生命周期管理
- 清晰的错误处理

### 2. IoC 容器层 (internal/app)

**职责**: 依赖注入和组件装配

```go
// internal/app/app.go
type App struct {
    Config        *config.Config
    ConfigManager config.Manager
    DB            database.Database
    Scheduler     scheduler.Scheduler
    Logger        logger.Logger
    Router        *gin.Engine
    server        *http.Server
}

func New(opts Options) (*App, error) {
    // 按依赖顺序初始化
    // 1. config → logger → database → scheduler
    // 2. repository → service → handler → router
    // 3. 启动配置监听
}
```

**初始化顺序**:

```
1. ConfigManager (配置管理)
   ↓
2. Logger (日志)
   ↓
3. Database (数据库)
   ↓
4. Scheduler (协程调度)
   ↓
5. Repository (数据访问)
   ↓
6. Service (业务逻辑)
   ↓
7. Handler (HTTP 处理)
   ↓
8. Router (路由)
```

### 3. 路由层 (internal/router)

**职责**: 路由定义和中间件配置

```go
// internal/router/router.go
func (r *Router) Setup(cfg middleware.MiddlewareConfig) *gin.Engine {
    // 中间件顺序很重要
    // 1. TraceID - 必须第一个
    // 2. Logger - 记录请求
    // 3. Recovery - 必须最后
    
    r.engine.Use(middleware.TraceID(cfg.TraceID))
    r.engine.Use(middleware.Logger(cfg.Logger, r.logger))
    r.engine.Use(middleware.Recovery(cfg.Recovery, r.logger))
    
    // 注册路由
    r.registerRoutes()
}
```

**路由组织**:

```
GET  /health                    - 健康检查
POST /api/v1/users/register     - 用户注册
POST /api/v1/users/login        - 用户登录
GET  /api/v1/users/:id          - 获取用户
GET  /api/v1/users              - 列表用户
```

### 4. 中间件层 (internal/middleware)

**职责**: 横切关注点处理

#### TraceID 中间件

```go
// 生成或提取 X-Request-ID
// 注入到 context
// 传递给后续处理
```

#### Logger 中间件

```go
// 记录请求信息
// - 方法、路径、状态码
// - 响应时间
// - TraceID
```

#### Recovery 中间件

```go
// 捕获 panic
// 返回 500 错误
// 包含 TraceID
```

### 5. 处理层 (internal/handler)

**职责**: HTTP 请求处理

```go
// internal/handler/user.go
type UserHandler struct {
    service service.UserService
}

func (h *UserHandler) Register(c *gin.Context) {
    // 1. 解析请求
    var req types.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, result.Error(
            errors.ErrInvalidParams,
            "invalid request",
        ))
        return
    }
    
    // 2. 调用服务
    resp, err := h.service.Register(c.Request.Context(), &req)
    if err != nil {
        // 3. 错误处理
        c.JSON(http.StatusInternalServerError, result.Error(
            errors.ErrInternalServer,
            err.Error(),
        ))
        return
    }
    
    // 4. 返回响应
    c.JSON(http.StatusOK, result.Success(resp))
}
```

**职责**:
- ✅ 参数解析和验证
- ✅ 调用服务
- ✅ 错误处理
- ✅ 响应构造
- ❌ 业务逻辑

### 6. 服务层 (internal/service)

**职责**: 业务逻辑处理

```go
// internal/service/user_impl.go
type userService struct {
    repo      repository.UserRepository
    scheduler scheduler.Scheduler
}

func (s *userService) Register(ctx context.Context, req *types.RegisterRequest) (*types.UserResponse, error) {
    // 1. 参数验证
    if err := validateRegisterRequest(req); err != nil {
        return nil, err
    }
    
    // 2. 业务逻辑
    user := &models.User{
        Username: req.Username,
        Email:    req.Email,
        Password: hashPassword(req.Password),
        Status:   1,
    }
    
    // 3. 数据持久化
    if err := s.repo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    // 4. 异步任务 (使用 Scheduler)
    s.scheduler.Submit(ctx, func(ctx context.Context) {
        // 发送欢迎邮件
        sendWelcomeEmail(user.Email)
    })
    
    return &types.UserResponse{
        UserID:    user.ID,
        Username:  user.Username,
        Email:     user.Email,
        Status:    user.Status,
        CreatedAt: user.CreatedAt,
    }, nil
}
```

**职责**:
- ✅ 业务逻辑
- ✅ 事务控制
- ✅ 异步任务调度
- ❌ HTTP 处理
- ❌ 数据库操作

### 7. 数据访问层 (internal/repository)

**职责**: 数据库操作抽象

```go
// internal/repository/user.go
type userRepository struct {
    db *gorm.DB
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*models.User, error) {
    var user models.User
    if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) FindAll(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
    var users []models.User
    var total int64
    
    offset := (page - 1) * pageSize
    if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
        return nil, 0, err
    }
    
    return users, total, nil
}
```

**职责**:
- ✅ CRUD 操作
- ✅ 查询构建
- ✅ 分页处理
- ❌ 业务逻辑
- ❌ 事务控制

### 8. 数据模型层 (internal/models)

**职责**: 数据结构定义

```go
// internal/models/user.go
type User struct {
    BaseModel
    Username string `gorm:"uniqueIndex;size:50;not null"`
    Email    string `gorm:"uniqueIndex;size:100;not null"`
    Password string `gorm:"size:255;not null"`
    Status   int    `gorm:"default:1"`
}

func (User) TableName() string {
    return "users"
}
```

### 9. 类型定义层 (types/)

**职责**: 跨层类型定义

```go
// types/errors/codes.go - 错误码
const (
    CodeSuccess = 0
    ErrInvalidParams = 1000
    ErrDuplicateUsername = 2001
    ErrUserNotFound = 4001
    ErrInternalServer = 5000
)

// types/result/result.go - 统一响应
type Result[T any] struct {
    Code       int    `json:"code"`
    Message    string `json:"message"`
    Data       T      `json:"data,omitempty"`
    TraceID    string `json:"traceId,omitempty"`
    ServerTime int64  `json:"serverTime"`
}

// types/request.go - 请求类型
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

// types/response.go - 响应类型
type UserResponse struct {
    UserID    int64     `json:"userId"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Status    int       `json:"status"`
    CreatedAt time.Time `json:"createdAt"`
}
```

### 10. 通用工具层 (pkg/)

**职责**: 可复用的工具库

#### Logger (pkg/logger)

```go
type Logger interface {
    Debug(msg string, keysAndValues ...interface{})
    Info(msg string, keysAndValues ...interface{})
    Warn(msg string, keysAndValues ...interface{})
    Error(msg string, keysAndValues ...interface{})
    Fatal(msg string, keysAndValues ...interface{})
    With(keysAndValues ...interface{}) Logger
    Sync() error
}
```

#### Scheduler (pkg/scheduler)

```go
type Scheduler interface {
    Submit(ctx context.Context, task func(context.Context)) error
    SubmitWithTimeout(ctx context.Context, timeout time.Duration, task func(context.Context)) error
    Shutdown(ctx context.Context) error
    Running() int
    Cap() int
}
```

#### Database (pkg/database)

```go
type Database interface {
    DB() *gorm.DB
    Close() error
    Ping() error
}
```

#### ID Generator (pkg/id)

```go
type Generator interface {
    NextID() int64
    NextString() string
}
```

## 配置管理

### 配置层次

```
1. 默认值 (代码中硬编码)
   ↓
2. 配置文件 (configs/config.yaml)
   ↓
3. 环境变量 (${VAR_NAME})
   ↓
4. 运行时更新 (ConfigManager.Update)
```

### 配置热重载

```go
// 监听文件变更
configManager.Watch()

// 注册变更 Hook
configManager.RegisterHook(func(old, new *config.Config) {
    log.Info("configuration updated",
        "oldPort", old.Server.Port,
        "newPort", new.Server.Port,
    )
})

// 原子更新
configManager.Update(func(cfg *config.Config) {
    cfg.Server.Port = 8081
})
```

## 错误处理

### 错误流

```
Handler → Service → Repository → Database
   ↓         ↓          ↓           ↓
BizError  BizError   DBError    GORMError
   ↓         ↓          ↓           ↓
   └─────────┴──────────┴───────────┘
                    ↓
              Result[T] with error code
```

### 错误码范围

| 范围 | 类型 | 示例 |
|------|------|------|
| 0 | 成功 | CodeSuccess |
| 1000-1999 | 参数错误 | ErrInvalidParams |
| 2000-2999 | 业务错误 | ErrDuplicateUsername |
| 3000-3999 | 认证错误 | ErrUnauthorized |
| 4000-4999 | 资源错误 | ErrUserNotFound |
| 5000-5999 | 系统错误 | ErrInternalServer |

## 依赖关系

### 允许的依赖

```
✅ cmd → app
✅ app → internal/*
✅ app → pkg/*
✅ internal/handler → internal/service
✅ internal/service → internal/repository
✅ internal/repository → internal/models
✅ internal/* → types
✅ internal/* → pkg
```

### 禁止的依赖

```
❌ pkg → internal
❌ types → 任何包
❌ internal/models → internal/service
❌ internal/repository → internal/handler
```

## 生命周期

### 启动流程

```
1. 加载配置
2. 初始化日志
3. 连接数据库
4. 初始化调度器
5. 初始化仓储
6. 初始化服务
7. 初始化处理器
8. 初始化路由
9. 启动配置监听
10. 启动 HTTP 服务
```

### 关闭流程

```
1. 停止接收新请求
2. 等待现有请求完成
3. 关闭 HTTP 服务
4. 关闭调度器
5. 关闭数据库连接
6. 同步日志
```

## 扩展点

### 中间件扩展

```go
// 自定义中间件
func CustomMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 前置处理
        c.Next()
        // 后置处理
    }
}

// 注册中间件
r.engine.Use(CustomMiddleware())
```

### 数据库 Hook

```go
// 实现 Hook 接口
type CustomHook struct{}

func (h *CustomHook) BeforeCreate(tx *gorm.DB) {
    // 创建前处理
}

func (h *CustomHook) AfterCreate(tx *gorm.DB) {
    // 创建后处理
}

// 注册 Hook
db, err := database.NewWithHooks(cfg, &CustomHook{})
```

### 配置变更通知

```go
// 注册 Hook
configManager.RegisterHook(func(old, new *config.Config) {
    // 处理配置变更
    if old.Logger.Level != new.Logger.Level {
        // 更新日志级别
    }
})
```

## 性能考虑

### 连接池

```yaml
database:
  maxOpenConns: 100    # 最大连接数
  maxIdleConns: 10     # 最大空闲连接数
```

### 协程池

```go
scheduler.Config{
    PoolSize:       10000,
    ExpiryDuration: time.Second,
}
```

### 缓存策略

- 配置缓存 (原子更新)
- 日志缓存 (批量写入)
- 数据库连接缓存 (连接池)

---

**最后更新**: 2025-12-30
