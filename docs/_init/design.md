# 架构设计文档

```
项目: Rei0721 | 版本: v1.0 | 更新: 2025-12-29
```

## 设计原则

### 核心原则
1. **资源受控** - 消除隐式并发与资源泄漏
2. **协议至上** - 泛型定义契约，实现依赖注入
3. **防御性封装** - 不信任外部输入和第三方库默认行为

### SOLID
- SRP: 单一职责
- OCP: 开闭原则 (接口+泛型扩展)
- LSP: 里氏替换
- ISP: 接口隔离
- DIP: 依赖倒置

---

## 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                     cmd/server/main.go                      │
└─────────────────────────┬───────────────────────────────────┘
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                    internal/app (IoC)                       │
│              依赖注入 · 组件装配 · 生命周期管理               │
└───────┬─────────────────┬─────────────────┬─────────────────┘
        ▼                 ▼                 ▼
┌───────────────┐ ┌───────────────┐ ┌───────────────┐
│    router     │ │    config     │ │      pkg      │
│   路由定义    │ │   配置定义    │ │   通用工具    │
└───────┬───────┘ └───────────────┘ └───────────────┘
        ▼
┌───────────────┐
│  middleware   │ ← 认证/日志/限流/错误恢复
└───────┬───────┘
        ▼
┌───────────────┐
│    handler    │ ← 参数解析/验证/响应构造
└───────┬───────┘
        ▼
┌───────────────┐
│    service    │ ← 业务逻辑/事务控制
└───────┬───────┘
        ▼
┌───────────────┐
│  repository   │ ← 数据访问抽象
└───────┬───────┘
        ▼
┌───────────────┐     ┌───────────────┐
│    models     │────▶│   Database    │
└───────────────┘     └───────────────┘
```

---

## 分层职责

| 层 | 目录 | 职责 | 限制 |
|---|------|------|------|
| Entry | `cmd/server/` | 启动入口 | 仅调用 app |
| App | `internal/app/` | IoC 容器 | 可引用所有 internal |
| Router | `internal/router/` | 路由定义 | - |
| Middleware | `internal/middleware/` | 请求预/后处理 | 支持配置开关+Hook |
| Handler | `internal/handler/` | HTTP 处理 | 无业务逻辑 |
| Service | `internal/service/` | 业务逻辑 | 不直接操作 DB |
| Repository | `internal/repository/` | 数据访问 | 无业务逻辑 |
| Models | `internal/models/` | 数据模型 | GORM 映射 |
| Config | `internal/config/` | 配置结构 | 禁止引用 Service/Repository |
| Types | `types/` | 类型定义 | **无依赖** |
| Pkg | `pkg/` | 通用工具 | **严禁引用 internal** |

---

## 核心组件设计

### 1. Config 层

```go
type Configurable interface {
    Validate() error
}

type ManagerInterface[T Configurable] interface {
    Load(handles ...HandlerFunc) error
    Get() *T                           // 线程安全快照
    Update(fn func(*T)) error          // 原子更新
    RegisterLogger(h LoggerHandler) Logger
    RegisterHook(h HookHandler)
}
```

**热重载流程:**
```
监听器 → 影子加载 → 重载中间件 → 创建新资源 → 原子切换
```

**设计要点:**
- 原子指针切换实现无锁读取
- 拦截 Viper 异步更新，确保时序可控
- 配置对象不可变

### 2. Scheduler 调度器

```go
type Scheduler interface {
    Submit(ctx context.Context, task func(context.Context)) error
    SubmitWithTimeout(ctx context.Context, timeout time.Duration, task func(context.Context)) error
    Shutdown(ctx context.Context) error
}
```

**设计要点:**
- 基于 `ants` 协程池
- 强制携带 `context.Context`
- 支持优雅关闭

### 3. Database 层

```
Application → Database Manager → [PostgreSQL|MySQL|SQLite] Driver → GORM → Hooks
```

**设计要点:**
- 统一管理入口
- 支持 Hook 注入 (监控耗时)
- 连接池管理

### 4. Logger 日志

```
启动阶段 (Go原生log) → 配置加载 (注入config层) → 重载阶段 (Zap/Logrus)
```

**设计原因:** 配置层需要提前注入日志库

---

## 数据流向

```
Client → Router → Middleware → Handler → Service → Repository → Database
                                                              ↘ Cache
```

**请求处理:**
1. Router: 路由匹配
2. Middleware: 认证/日志/限流
3. Handler: 参数验证
4. Service: 业务逻辑
5. Repository: 数据访问
6. 返回响应

---

## 设计模式

| 模式 | 应用 |
|------|------|
| 依赖注入 | `internal/app` IoC 容器 |
| 仓储模式 | 数据访问抽象 |
| 工厂模式 | 数据库驱动创建 |
| 策略模式 | 多数据库支持 |
| 观察者模式 | 配置热重载 Hook |
| 中间件模式 | 请求处理链 |
| 单例模式 | 配置管理器/调度器 |
| 泛型编程 | 工具库解耦 |

---

## 扩展指南

### 添加新功能步骤

1. `internal/models/` - 定义数据模型
2. `types/` - 定义请求/响应结构
3. `internal/repository/` - 实现数据访问
4. `internal/service/` - 实现业务逻辑
5. `internal/handler/` - 实现 HTTP 处理
6. `internal/router/` - 注册路由
7. `internal/app/` - 完成组件装配

### 最佳实践

- ✅ 常量优先，禁止魔法值
- ✅ 统一错误码
- ✅ 异步任务使用 Scheduler
- ✅ 配置变更注册 Hook
- ✅ 第三方库必须封装
- ✅ 提交前完成自检清单

---

[← README](./README.md) | [protocol.md](./protocol.md) | [api.md →](./api.md)
