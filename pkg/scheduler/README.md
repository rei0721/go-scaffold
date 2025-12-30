# Scheduler - 协程池任务调度器

[English](#english) | [中文](#chinese)

---

## <a name="chinese"></a>中文文档

### 概述

Scheduler 是一个基于 `ants` 的高性能协程池管理器，用于统一管理应用中的所有异步任务。它提供了任务提交、超时控制、优雅关闭等功能，帮助你避免 goroutine 泄漏和资源耗尽问题。

### 特性

- ✅ **协程池管理** - 复用 goroutine，减少创建/销毁开销
- ✅ **并发限制** - 防止无限制创建 goroutine 导致资源耗尽
- ✅ **Context 支持** - 完整支持上下文取消和超时
- ✅ **优雅关闭** - 等待所有任务完成后再关闭
- ✅ **资源监控** - 提供运行中任务数和池容量查询
- ✅ **自动回收** - 空闲 goroutine 自动过期回收
- ✅ **高性能** - 基于 ants 库，经过大规模生产验证

### 安装

```bash
go get github.com/panjf2000/ants/v2
```

### 快速开始

#### 1. 创建 Scheduler 实例

```go
package main

import (
    "context"
    "log"
    "time"
    "github.com/rei0721/rei0721/pkg/scheduler"
)

func main() {
    // 方式1: 使用默认配置
    sched := scheduler.Default()

    // 方式2: 使用自定义配置
    cfg := &scheduler.Config{
        PoolSize:         5000,           // 最大 goroutine 数
        MaxBlockingTasks: 0,              // 无限制等待队列
        ExpiryDuration:   time.Second,    // 1秒后回收空闲 goroutine
    }
    sched, err := scheduler.New(cfg)
    if err != nil {
        log.Fatal(err)
    }

    defer func() {
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        sched.Shutdown(ctx)
    }()
}
```

#### 2. 提交异步任务

```go
// 提交普通任务
err := sched.Submit(ctx, func(taskCtx context.Context) {
    // 执行异步操作
    fmt.Println("Task is running")

    // 检查上下文是否被取消
    select {
    case <-taskCtx.Done():
        fmt.Println("Task cancelled")
        return
    default:
        // 继续执行
    }
})
```

#### 3. 提交带超时的任务

```go
// 任务最多执行 5 秒
err := sched.SubmitWithTimeout(ctx, 5*time.Second, func(taskCtx context.Context) {
    // 这个任务有 5 秒超时限制
    // 超时后 taskCtx.Done() 会被关闭

    select {
    case <-taskCtx.Done():
        fmt.Println("Task timeout or cancelled")
        return
    case <-time.After(3 * time.Second):
        fmt.Println("Task completed within timeout")
    }
})
```

### API 文档

#### Scheduler 接口

```go
type Scheduler interface {
    // Submit 提交一个任务到协程池
    Submit(ctx context.Context, task func(context.Context)) error

    // SubmitWithTimeout 提交一个带超时的任务
    SubmitWithTimeout(ctx context.Context, timeout time.Duration, task func(context.Context)) error

    // Shutdown 优雅关闭调度器
    Shutdown(ctx context.Context) error

    // Running 返回当前正在运行的 goroutine 数量
    Running() int

    // Cap 返回协程池的容量
    Cap() int
}
```

#### Config 配置

```go
type Config struct {
    // PoolSize 协程池中的最大 goroutine 数量
    // 推荐值:
    //   - CPU 密集型: CPU 核心数的 1-2 倍
    //   - IO 密集型: CPU 核心数的 10-100 倍
    // 默认: 10000
    PoolSize int

    // MaxBlockingTasks 最大阻塞任务数
    // 0 表示无限制队列
    // >0 表示限制等待队列大小
    // 默认: 0
    MaxBlockingTasks int

    // ExpiryDuration 空闲 goroutine 的过期时间
    // 超过此时间未使用的 goroutine 会被回收
    // 默认: 1 秒
    ExpiryDuration time.Duration
}
```

### 使用场景

#### 场景 1: 异步发送邮件

```go
// 在用户注册后发送欢迎邮件
func (s *UserService) RegisterUser(ctx context.Context, user *User) error {
    // 保存用户到数据库
    if err := s.repo.Create(ctx, user); err != nil {
        return err
    }

    // 异步发送欢迎邮件(不阻塞主流程)
    s.scheduler.Submit(ctx, func(taskCtx context.Context) {
        if err := s.emailService.SendWelcome(taskCtx, user.Email); err != nil {
            log.Error("failed to send welcome email", "error", err)
        }
    })

    return nil
}
```

#### 场景 2: 记录审计日志

```go
// 记录用户操作日志
func (h *Handler) UpdateProfile(c *gin.Context) {
    var req UpdateProfileRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 更新用户资料
    user, err := h.service.UpdateProfile(c.Request.Context(), req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    // 异步记录审计日志
    h.scheduler.Submit(c.Request.Context(), func(taskCtx context.Context) {
        auditLog := &AuditLog{
            UserID:    user.ID,
            Action:    "update_profile",
            Timestamp: time.Now(),
        }
        h.auditService.Log(taskCtx, auditLog)
    })

    c.JSON(200, user)
}
```

#### 场景 3: 调用外部 API

```go
// 调用第三方服务，设置超时
func (s *PaymentService) ProcessPayment(ctx context.Context, order *Order) error {
    // 调用支付网关 API，最多等待 10 秒
    var paymentResult *PaymentResult
    var apiErr error

    err := s.scheduler.SubmitWithTimeout(ctx, 10*time.Second, func(taskCtx context.Context) {
        paymentResult, apiErr = s.paymentGateway.Charge(taskCtx, order.Amount)
    })

    if err != nil {
        return fmt.Errorf("failed to submit payment task: %w", err)
    }

    if apiErr != nil {
        return fmt.Errorf("payment API error: %w", apiErr)
    }

    return s.repo.UpdatePaymentStatus(ctx, order.ID, paymentResult.Status)
}
```

#### 场景 4: 批量处理任务

```go
// 批量导入用户数据
func (s *UserService) BatchImport(ctx context.Context, users []*User) error {
    for _, user := range users {
        // 为每个用户提交一个导入任务
        u := user // 避免闭包问题

        if err := s.scheduler.Submit(ctx, func(taskCtx context.Context) {
            if err := s.repo.Create(taskCtx, u); err != nil {
                log.Error("failed to import user", "userId", u.ID, "error", err)
            }
        }); err != nil {
            log.Error("failed to submit import task", "error", err)
        }
    }

    return nil
}
```

### 监控和调试

#### 检查调度器状态

```go
// 获取当前运行的任务数
running := sched.Running()
fmt.Printf("Running tasks: %d\n", running)

// 获取协程池容量
capacity := sched.Cap()
fmt.Printf("Pool capacity: %d\n", capacity)

// 计算资源利用率
utilization := float64(running) / float64(capacity) * 100
fmt.Printf("Pool utilization: %.2f%%\n", utilization)
```

#### 集成到监控系统

```go
// Prometheus metrics 示例
import "github.com/prometheus/client_golang/prometheus"

var (
    poolUtilization = prometheus.NewGaugeFunc(
        prometheus.GaugeOpts{
            Name: "scheduler_pool_utilization",
            Help: "Current utilization of the scheduler pool",
        },
        func() float64 {
            return float64(sched.Running()) / float64(sched.Cap())
        },
    )
)

func init() {
    prometheus.MustRegister(poolUtilization)
}
```

### 最佳实践

#### 1. 始终通过 Context 控制生命周期

```go
// ✅ 正确: 任务尊重 context 的取消信号
sched.Submit(ctx, func(taskCtx context.Context) {
    select {
    case <-taskCtx.Done():
        log.Info("task cancelled")
        return
    default:
        // 执行任务
    }
})

// ❌ 错误: 任务忽略 context，可能导致资源泄漏
sched.Submit(ctx, func(taskCtx context.Context) {
    // 没有检查 taskCtx.Done()，任务可能永远运行
    time.Sleep(10 * time.Minute)
})
```

#### 2. 长时间运行的任务应定期检查取消信号

```go
sched.Submit(ctx, func(taskCtx context.Context) {
    for i := 0; i < 1000; i++ {
        // 每次迭代都检查是否被取消
        select {
        case <-taskCtx.Done():
            log.Info("task cancelled at iteration", "i", i)
            return
        default:
        }

        // 执行工作
        processItem(i)
    }
})
```

#### 3. 为关键任务设置超时

```go
// API 调用应该有明确的超时限制
sched.SubmitWithTimeout(ctx, 5*time.Second, func(taskCtx context.Context) {
    resp, err := http.Get("https://api.example.com/data")
    if err != nil {
        log.Error("API call failed", "error", err)
        return
    }
    defer resp.Body.Close()
    // 处理响应...
})
```

#### 4. 优雅关闭时等待任务完成

```go
func main() {
    sched := scheduler.Default()

    // 注册关闭信号处理
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    go func() {
        <-sigChan
        log.Info("shutting down scheduler...")

        // 最多等待 30 秒
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()

        if err := sched.Shutdown(ctx); err != nil {
            log.Error("scheduler shutdown error", "error", err)
        } else {
            log.Info("scheduler shutdown successfully")
        }

        os.Exit(0)
    }()

    // 应用逻辑...
}
```

#### 5. 根据任务类型调整池大小

```go
// CPU 密集型任务 (计算、数据处理等)
cpuIntensiveConfig := &scheduler.Config{
    PoolSize: runtime.NumCPU() * 2, // CPU 核心数的 2 倍
}

// IO 密集型任务 (数据库查询、API 调用等)
ioIntensiveConfig := &scheduler.Config{
    PoolSize: runtime.NumCPU() * 50, // CPU 核心数的 50 倍
}
```

### 性能考虑

#### 协程池大小

- **过小**: 任务会在队列中等待，降低吞吐量
- **过大**: 消耗过多内存，上下文切换开销增加
- **推荐**: 根据任务类型和系统资源调整

#### 过期时间

- **短过期时间** (1-5 秒): 快速回收空闲 goroutine，节省内存
- **长过期时间** (30-60 秒): 减少 goroutine 创建/销毁次数，适合持续高负载
- **默认值** (1 秒): 平衡内存使用和性能

#### 阻塞队列大小

- **无限制** (`MaxBlockingTasks=0`): 确保所有任务最终都能执行，但可能占用大量内存
- **有限制** (`MaxBlockingTasks>0`): 防止内存耗尽，但任务可能被拒绝

### 错误处理

```go
// 提交任务时处理错误
if err := sched.Submit(ctx, task); err != nil {
    if errors.Is(err, ants.ErrPoolClosed) {
        log.Error("scheduler is closed")
    } else if errors.Is(err, ants.ErrPoolOverload) {
        log.Error("scheduler pool is overloaded")
    } else {
        log.Error("failed to submit task", "error", err)
    }
}

// 在任务内部处理错误
sched.Submit(ctx, func(taskCtx context.Context) {
    defer func() {
        if r := recover(); r != nil {
            log.Error("task panic", "panic", r)
        }
    }()

    // 任务逻辑...
})
```

### 注意事项

1. **Context 取消不会强制终止 goroutine** - 任务必须主动检查 `ctx.Done()` 并退出
2. **避免闭包陷阱** - 在循环中提交任务时，要注意变量捕获问题
3. **不要阻塞任务** - 长时间阻塞会降低协程池效率
4. **监控资源使用** - 定期检查 `Running()` 和 `Cap()`，避免过载

### 项目结构示例

```
pkg/scheduler/
├── scheduler.go    # 接口定义和配置
├── ants.go         # ants 实现
├── README.md       # 本文档
└── doc.go          # 包文档
```

### 参考链接

- [ants 协程池](https://github.com/panjf2000/ants)
- [Go Context 包](https://pkg.go.dev/context)
- [Go 并发模式](https://go.dev/blog/pipelines)

---

## <a name="english"></a>English Documentation

### Overview

Scheduler is a high-performance goroutine pool manager based on `ants`, designed to uniformly manage all asynchronous tasks in your application. It provides task submission, timeout control, graceful shutdown, and other features to help you avoid goroutine leaks and resource exhaustion.

### Features

- ✅ **Goroutine Pool Management** - Reuse goroutines, reduce creation/destruction overhead
- ✅ **Concurrency Limiting** - Prevent unlimited goroutine creation from exhausting resources
- ✅ **Context Support** - Full support for context cancellation and timeout
- ✅ **Graceful Shutdown** - Wait for all tasks to complete before closing
- ✅ **Resource Monitoring** - Query running task count and pool capacity
- ✅ **Auto Recycling** - Idle goroutines automatically expire and get recycled
- ✅ **High Performance** - Based on ants library, proven in large-scale production

### Installation

```bash
go get github.com/panjf2000/ants/v2
```

### Quick Start

#### 1. Create Scheduler Instance

```go
package main

import (
    "context"
    "log"
    "time"
    "github.com/rei0721/rei0721/pkg/scheduler"
)

func main() {
    // Method 1: Use default config
    sched := scheduler.Default()

    // Method 2: Use custom config
    cfg := &scheduler.Config{
        PoolSize:         5000,           // Max goroutines
        MaxBlockingTasks: 0,              // Unlimited queue
        ExpiryDuration:   time.Second,    // Recycle idle goroutines after 1s
    }
    sched, err := scheduler.New(cfg)
    if err != nil {
        log.Fatal(err)
    }

    defer func() {
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        sched.Shutdown(ctx)
    }()
}
```

#### 2. Submit Async Tasks

```go
// Submit a regular task
err := sched.Submit(ctx, func(taskCtx context.Context) {
    // Perform async operation
    fmt.Println("Task is running")

    // Check if context is cancelled
    select {
    case <-taskCtx.Done():
        fmt.Println("Task cancelled")
        return
    default:
        // Continue execution
    }
})
```

#### 3. Submit Tasks with Timeout

```go
// Task runs for maximum 5 seconds
err := sched.SubmitWithTimeout(ctx, 5*time.Second, func(taskCtx context.Context) {
    // This task has a 5-second timeout
    // taskCtx.Done() will close after timeout

    select {
    case <-taskCtx.Done():
        fmt.Println("Task timeout or cancelled")
        return
    case <-time.After(3 * time.Second):
        fmt.Println("Task completed within timeout")
    }
})
```

### License

MIT License
