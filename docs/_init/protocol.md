# 工程规范协议

```
项目: Rei0721 | 版本: v1.0 | 状态: Enforced | 更新: 2025-12-29
```

## 三大核心原则

1. **资源受控** - 消除隐式并发与资源泄漏
2. **协议至上** - 泛型定义契约，实现依赖注入
3. **防御性封装** - 不信任外部输入，不信任第三方库默认行为

---

## 1. 资源调度协议

### 协程管控

| 规则 | 说明 |
|------|------|
| 🚫 禁止 | `internal/` 中直接使用 `go` 关键字 |
| ✅ 必须 | 异步任务提交至 `pkg/scheduler` |
| ✅ 必须 | 任务携带 `context.Context` |

```go
// ❌ 错误
go func() { processData() }()

// ✅ 正确
scheduler.Submit(ctx, func(taskCtx context.Context) {
    processData(taskCtx)
})
```

### 池化资源

| 规则 | 说明 |
|------|------|
| 🚫 禁止 | 业务层手动创建 `gorm.Open` / `redis.NewClient` |
| ✅ 必须 | 通过 `internal/app` 或 `pkg/database` 统一初始化 |

---

## 2. 配置状态协议

### 接口契约

```go
type Configurable interface {
    Validate() error
}

type ManagerInterface[T Configurable] interface {
    Get() *T                      // 只读快照
    Update(fn func(*T)) error     // 原子更新
    RegisterHook(h HookHandler)   // 变更钩子
}
```

### 变更规范

| 规则 | 说明 |
|------|------|
| 不可变性 | `Get()` 返回对象禁止修改 |
| 热重载 | 影子加载 → 校验 → 中间件拦截 → 原子切换 |
| Hook | 需感知变更必须注册 Hook，禁止轮询 |

```go
// ❌ 错误
cfg := configManager.Get()
cfg.Server.Port = 8080

// ✅ 正确
configManager.Update(func(cfg *Config) {
    cfg.Server.Port = 8080
})

configManager.RegisterHook(func(old, new *Config) {
    if old.Database.DSN != new.Database.DSN {
        reinitDatabase(new.Database)
    }
})
```

---

## 3. 错误与日志协议

### 日志规范

| 规则 | 说明 |
|------|------|
| 🚫 禁止 | `fmt.Println` / `log.Println` |
| ✅ 必须 | 使用 `pkg/logger` |
| ⚠️ 例外 | 配置加载前可用原生 log，加载后立即切换 |

```go
// ❌ 错误
fmt.Println("User login:", userID)

// ✅ 正确
logger.Info("User login", "userID", userID, "ip", clientIP)
```

### 错误处理

| 规则 | 说明 |
|------|------|
| 标准化 | API 错误必须包含 Code, Message, TraceID |
| 国际化 | 用户提示使用 `go-i18n` key，禁止硬编码 |
| 常量优先 | 错误码在 `types/errors` 定义，禁止魔法数字 |

```go
// ❌ 错误
return c.JSON(400, map[string]interface{}{
    "code": 1001, "message": "用户名不能为空",
})

// ✅ 正确
return result.Error(errors.ErrInvalidUsername, i18n.T(ctx, "error.username_required"))
```

---

## 4. 封装与扩展协议

### 泛型驱动

| 规则 | 说明 |
|------|------|
| 规范 | `pkg/` 工具库必须使用泛型 |
| 解耦 | 工具库严禁引用 `internal/` |

```go
// pkg/result/result.go
type Result[T any] struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    T      `json:"data,omitempty"`
}

func Success[T any](data T) *Result[T] {
    return &Result[T]{Code: 0, Message: "success", Data: data}
}
```

### 行为拦截

| 规则 | 说明 |
|------|------|
| 必须 | 第三方库封装预留 Hook 插槽 |
| 示例 | GORM 注入 `BeforeCreate` / `AfterQuery` 钩子 |
| 中间件 | 支持配置开关 + Hook 扩展 |

---

## 5. 目录与分层规范

| 目录 | 职责 | 限制 |
|------|------|------|
| `cmd/server/` | 启动入口 | 无业务逻辑 |
| `internal/app/` | 依赖注入容器 | 可引用所有 internal |
| `internal/config/` | 配置定义 | 禁止引用 Service/Repository |
| `types/` | 契约定义 | **无依赖** |
| `pkg/` | 通用工具 | **严禁引用 internal** |

### 依赖方向

```
cmd → app → handler → service → repository → models
         ↘ config
         ↘ pkg ✓

pkg → internal ✗
types → * ✗
```

---

## 6. 开发自检清单

### 资源检查
- [ ] 无 `go func()` (改用 `scheduler.Submit`)
- [ ] 异步任务传递 `context`

### 配置检查
- [ ] 无直接修改配置对象 (改用 `Manager.Update`)
- [ ] 配置变更注册 Hook

### 代码质量
- [ ] 无硬编码错误文本 (改用 i18n key)
- [ ] 错误码使用常量

### 架构检查
- [ ] `pkg/` 未引用 `internal/`
- [ ] 无跨层级反向依赖

### 测试与文档
- [ ] 添加单元测试
- [ ] 更新相关文档

---

## 违规案例

### 案例 1: 直接使用 go

```go
// ❌
go func() { handleOrder(orderID) }()

// ✅
scheduler.Submit(ctx, func(taskCtx context.Context) {
    handleOrder(taskCtx, orderID)
})
```

### 案例 2: 硬编码错误

```go
// ❌
return errors.New("用户不存在")

// ✅
return errors.NewBizError(errors.ErrUserNotFound, i18n.T(ctx, "error.user_not_found"))
```

### 案例 3: pkg 引用 internal

```go
// ❌ pkg/helper/user.go
import "myproject/internal/service"

// ✅ 使用泛型和接口
type UserGetter interface { GetName(id int) string }
func GetUserName[T UserGetter](svc T, id int) string { return svc.GetName(id) }
```

---

> 违反协议的代码将被 CI/CD 或 Code Review 驳回

[← README](./README.md) | [design.md →](./design.md)
