# Daemon Manager - 守护进程管理器

---

## <a name="chinese"></a>中文文档

### 概述

Daemon Manager 是一个用于管理长期运行服务的管理器。它提供了统一的接口来启动、停止和监控各种类型的守护进程(HTTP 服务器、gRPC 服务器、消息队列消费者等)。

### 与其他包的区别

| 包名                   | 用途         | 适用场景                               |
| ---------------------- | ------------ | -------------------------------------- |
| **`pkg/daemon`**       | 管理长期服务 | HTTP 服务器、gRPC 服务器、Kafka 消费者 |
| **`pkg/scheduler`**    | 管理短期任务 | 发送邮件、记录日志、API 调用           |
| **`internal/service`** | 业务逻辑层   | 用户管理、订单处理等业务功能           |

**简单来说**：

- `daemon` = 管理**服务**的生命周期
- `scheduler` = 管理**任务**的执行
- `service` = 实现**业务**逻辑

### 特性

- ✅ **统一管理** - 所有长期服务在一个地方注册和管理
- ✅ **优雅启动** - 并发启动所有服务，提高启动速度
- ✅ **优雅关闭** - 并发停止所有服务，支持超时控制
- ✅ **错误处理** - 完整的错误日志和错误传播
- ✅ **易于扩展** - 实现 Daemon 接口即可添加新服务
- ✅ **线程安全** - 可以在多个 goroutine 中安全调用

### 快速开始

#### 1. 实现 Daemon 接口

```go
package daemons

import (
    "context"
    "net/http"
    "github.com/rei0721/rei0721/pkg/daemon"
)

// HTTPDaemon HTTP 服务器守护进程
type HTTPDaemon struct {
    server *http.Server
    logger daemon.Logger
}

// NewHTTPDaemon 创建 HTTP 守护进程
func NewHTTPDaemon(addr string, handler http.Handler, logger daemon.Logger) *HTTPDaemon {
    return &HTTPDaemon{
        server: &http.Server{
            Addr:    addr,
            Handler: handler,
        },
        logger: logger,
    }
}

// Name 返回守护进程名称
func (d *HTTPDaemon) Name() string {
    return "http-server"
}

// Start 启动 HTTP 服务器
func (d *HTTPDaemon) Start(ctx context.Context) error {
    d.logger.Info("HTTP server starting", "addr", d.server.Addr)

    // 在后台启动服务器
    go func() {
        if err := d.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            d.logger.Error("HTTP server error", "error", err)
        }
    }()

    return nil
}

// Stop 停止 HTTP 服务器
func (d *HTTPDaemon) Stop(ctx context.Context) error {
    d.logger.Info("HTTP server shutting down")
    return d.server.Shutdown(ctx)
}
```

#### 2. 使用 Manager 管理守护进程

```go
package main

import (
    "context"
    "time"

    "github.com/rei0721/rei0721/pkg/daemon"
    "github.com/rei0721/rei0721/internal/daemons"
)

func main() {
    // 创建管理器
    manager := daemon.NewManager(logger)

    // 注册 HTTP 守护进程
    httpDaemon := daemons.NewHTTPDaemon(":8080", router, logger)
    manager.Register(httpDaemon)

    // 启动所有守护进程
    ctx := context.Background()
    if err := manager.Start(ctx); err != nil {
        panic(err)
    }

    // 等待终止信号...

    // 优雅关闭(最多等待 30 秒)
    stopCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := manager.Stop(stopCtx); err != nil {
        log.Error("shutdown error", "error", err)
    }
}
```

### API 文档

#### Daemon 接口

```go
type Daemon interface {
    // Name 返回守护进程的名称
    Name() string

    // Start 启动守护进程
    Start(ctx context.Context) error

    // Stop 停止守护进程
    Stop(ctx context.Context) error
}
```

**核心方法说明**：

##### Name()

- 返回守护进程的唯一名称
- 用于日志记录和错误报告
- 建议使用小写字母和连字符，如 `"http-server"`

##### Start(ctx)

- 启动守护进程
- 应该快速返回，长期运行的逻辑在后台 goroutine 中
- 必须检查 `ctx.Done()` 以支持取消
- 启动失败返回有意义的错误

##### Stop(ctx)

- 优雅关闭守护进程
- 停止接收新请求
- 等待现有请求完成
- 清理所有资源
- 必须在 `ctx` 超时前完成

#### Manager 管理器

```go
type Manager struct {
    // 私有字段
}

// NewManager 创建管理器
func NewManager(logger Logger) *Manager

// Register 注册守护进程
func (m *Manager) Register(daemon Daemon)

// Start 启动所有守护进程
func (m *Manager) Start(ctx context.Context) error

// Stop 停止所有守护进程
func (m *Manager) Stop(ctx context.Context) error
```

**方法说明**：

##### Register(daemon)

- 注册一个守护进程到管理器
- 可以注册多个守护进程
- 应该在 Start 之前调用
- 线程安全

##### Start(ctx)

- 并发启动所有已注册的守护进程
- 任何一个失败都会立即返回错误
- 建议使用 `context.Background()`

##### Stop(ctx)

- 并发停止所有守护进程
- 建议设置超时，如 30 秒
- 超时后强制返回

### 使用场景

#### 场景 1: 管理多个服务器

```go
// 创建管理器
manager := daemon.NewManager(logger)

// 注册 HTTP 服务器
httpDaemon := daemons.NewHTTPDaemon(":8080", httpRouter, logger)
manager.Register(httpDaemon)

// 注册 gRPC 服务器
grpcDaemon := daemons.NewGRPCDaemon(":9090", logger)
manager.Register(grpcDaemon)

// 注册指标服务器
metricsDaemon := daemons.NewMetricsDaemon(":9091", logger)
manager.Register(metricsDaemon)

// 一次性启动所有服务
ctx := context.Background()
if err := manager.Start(ctx); err != nil {
    log.Fatal("failed to start services", "error", err)
}
```

#### 场景 2: Kafka 消费者

```go
// KafkaDaemon Kafka 消费者守护进程
type KafkaDaemon struct {
    consumer *kafka.Consumer
    logger   daemon.Logger
    ctx      context.Context
    cancel   context.CancelFunc
}

func (d *KafkaDaemon) Name() string {
    return "kafka-consumer"
}

func (d *KafkaDaemon) Start(ctx context.Context) error {
    d.ctx, d.cancel = context.WithCancel(ctx)

    go func() {
        for {
            select {
            case <-d.ctx.Done():
                return
            default:
                msg, err := d.consumer.ReadMessage(d.ctx)
                if err != nil {
                    d.logger.Error("kafka read error", "error", err)
                    continue
                }
                d.handleMessage(msg)
            }
        }
    }()

    return nil
}

func (d *KafkaDaemon) Stop(ctx context.Context) error {
    d.cancel() // 停止消费循环
    return d.consumer.Close()
}
```

#### 场景 3: 定时任务调度器

```go
// CronDaemon 定时任务守护进程
type CronDaemon struct {
    cron   *cron.Cron
    logger daemon.Logger
}

func (d *CronDaemon) Name() string {
    return "cron-scheduler"
}

func (d *CronDaemon) Start(ctx context.Context) error {
    d.logger.Info("cron scheduler starting")

    // 添加定时任务
    d.cron.AddFunc("@every 1h", func() {
        d.logger.Info("running hourly task")
        // 执行任务...
    })

    d.cron.Start()
    return nil
}

func (d *CronDaemon) Stop(ctx context.Context) error {
    d.logger.Info("cron scheduler stopping")
    stopCtx := d.cron.Stop()

    select {
    case <-stopCtx.Done():
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

### 最佳实践

#### 1. 守护进程命名规范

```go
// ✅ 好的命名
"http-server"
"grpc-server"
"kafka-consumer"
"cron-scheduler"
"metrics-server"

// ❌ 不好的命名
"Server"        // 太泛化
"http"          // 不够具体
"myService123"  // 含义不明
```

#### 2. 正确处理 Context

```go
func (d *MyDaemon) Start(ctx context.Context) error {
    // ✅ 正确：在后台 goroutine 中检查 ctx.Done()
    go func() {
        for {
            select {
            case <-ctx.Done():
                d.logger.Info("daemon stopped by context")
                return
            default:
                // 执行工作...
            }
        }
    }()

    return nil
}

// ❌ 错误：忽略 context，导致无法停止
func (d *MyDaemon) Start(ctx context.Context) error {
    go func() {
        for {
            // 没有检查 ctx.Done()
            // 这个 goroutine 永远不会停止！
            doWork()
        }
    }()
    return nil
}
```

#### 3. 优雅关闭

```go
func (d *MyDaemon) Stop(ctx context.Context) error {
    // ✅ 正确的优雅关闭流程

    // 1. 停止接收新请求/任务
    d.stopAccepting()

    // 2. 等待现有请求/任务完成
    done := make(chan struct{})
    go func() {
        d.waitForCompletion()
        close(done)
    }()

    // 3. 等待完成或超时
    select {
    case <-done:
        d.logger.Info("all tasks completed")
    case <-ctx.Done():
        d.logger.Warn("shutdown timeout, forcing stop")
    }

    // 4. 清理资源
    d.cleanup()

    return nil
}
```

#### 4. 错误处理

```go
func (d *MyDaemon) Start(ctx context.Context) error {
    // ✅ 返回有意义的错误
    if err := d.initialize(); err != nil {
        return fmt.Errorf("failed to initialize: %w", err)
    }

    if err := d.connect(); err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    return nil
}
```

#### 5. 设置合理的超时

```go
func main() {
    manager := daemon.NewManager(logger)
    manager.Register(httpDaemon)

    // 启动：使用 Background context
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal(err)
    }

    // 关闭：设置 30 秒超时
    stopCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := manager.Stop(stopCtx); err != nil {
        log.Error("shutdown error", "error", err)
    }
}
```

### 常见问题

#### Q: Daemon 和 Service 有什么区别？

**A**:

- **Daemon** (守护进程): 管理**服务的生命周期**（启动、停止）
  - 例如：HTTP 服务器、gRPC 服务器
  - 位置：`pkg/daemon`、`internal/daemons`
- **Service** (业务服务): 实现**业务逻辑**
  - 例如：用户管理、订单处理
  - 位置：`internal/service`

#### Q: 什么时候用 Daemon，什么时候用 Scheduler？

**A**:

- **Daemon**: 长期运行的服务
  - HTTP/gRPC 服务器
  - Kafka 消费者
  - 定时任务调度器
- **Scheduler**: 短期异步任务
  - 发送邮件
  - 记录日志
  - API 调用

#### Q: Start 方法应该阻塞吗？

**A**:
不应该。Start 方法应该快速返回，长期运行的逻辑应该在后台 goroutine 中执行。

```go
// ✅ 正确
func (d *MyDaemon) Start(ctx context.Context) error {
    go func() {
        // 长期运行的逻辑
    }()
    return nil // 快速返回
}

// ❌ 错误
func (d *MyDaemon) Start(ctx context.Context) error {
    // 阻塞在这里，其他守护进程无法启动
    d.server.ListenAndServe()
    return nil
}
```

#### Q: 如何确保所有守护进程都已启动？

**A**:
Manager.Start() 会等待所有守护进程启动完成或出错。如果你需要更精细的控制，可以在守护进程中使用 ready channel。

### 项目结构

```
pkg/daemon/              # 核心包
├── constants.go        # 常量定义
├── daemon.go           # Daemon 接口
├── manager.go          # Manager 实现
├── doc.go              # 包文档
└── README.md           # 本文档

internal/daemons/        # 具体实现
├── http_daemon.go      # HTTP 服务器
├── grpc_daemon.go      # gRPC 服务器（未来）
└── kafka_daemon.go     # Kafka 消费者（未来）

internal/service/        # 业务逻辑（不变）
├── user.go
└── user_impl.go
```

### 参考链接

- [Go Context 包](https://pkg.go.dev/context)
- [优雅关闭](https://go.dev/blog/context)
- [并发模式](https://go.dev/blog/pipelines)

---

## <a name="english"></a>English Documentation

### Overview

Daemon Manager is a manager for long-running services. It provides a unified interface to start, stop, and monitor various types of daemons (HTTP servers, gRPC servers, message queue consumers, etc.).

### Features

- ✅ **Unified Management** - All long-running services registered and managed in one place
- ✅ **Graceful Startup** - Concurrent startup of all services for faster boot time
- ✅ **Graceful Shutdown** - Concurrent shutdown with timeout control
- ✅ **Error Handling** - Complete error logging and propagation
- ✅ **Easy Extension** - Just implement the Daemon interface
- ✅ **Thread-Safe** - Safe for concurrent use

### Quick Start

See the Chinese documentation above for detailed examples.

### License

MIT License
