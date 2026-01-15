# pkg/httpserver

`pkg/httpserver` 提供统一的 HTTP 服务器接口，基于标准库 `net/http` 和 Gin 框架，支持配置热更新和优雅关闭。

## 特性

- **统一接口**: 抽象 HTTP 服务器实现，易于测试和替换
- **配置热更新**: 支持运行时更新服务器配置（端口、超时等）
- **优雅关闭**: 等待现有请求完成后关闭服务器
- **线程安全**: 所有操作都是并发安全的
- **自动端口分配**: 支持自动分配可用端口
- **地址验证**: 自动验证和修正监听地址

## 安装

```bash
go get github.com/rei0721/rei0721/pkg/httpserver
```

## 快速开始

### 基本使用

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/rei0721/rei0721/pkg/httpserver"
    "github.com/rei0721/rei0721/pkg/logger"
)

func main() {
    // 创建 Gin Router
    router := gin.Default()
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // 创建日志器
    log, _ := logger.New(&logger.Config{
        Level:  "info",
        Format: "console",
    })

    // 创建 HTTP Server 配置
    config := &httpserver.Config{
        Host:         "localhost",
        Port:         8080,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // 创建 HTTP Server
    server, err := httpserver.New(router, config, log)
    if err != nil {
        log.Fatal("failed to create server", "error", err)
    }

    // 启动服务器
    if err := server.Start(context.Background()); err != nil {
        log.Fatal("failed to start server", "error", err)
    }

    // 等待信号...
    // 优雅关闭
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Error("shutdown error", "error", err)
    }
}
```

## 配置说明

### Config 结构

```go
type Config struct {
    Host         string        // 监听地址，例如 "localhost", "0.0.0.0"
    Port         int           // 监听端口，范围 1-65535，0 表示自动分配
    ReadTimeout  time.Duration // 读取超时
    WriteTimeout time.Duration // 写入超时
    IdleTimeout  time.Duration // 空闲连接超时
}
```

### 默认值

```go
DefaultHost         = "localhost"
DefaultPort         = 8080
DefaultReadTimeout  = 15 * time.Second
DefaultWriteTimeout = 15 * time.Second
DefaultIdleTimeout  = 60 * time.Second
```

### 配置验证

配置会自动验证：

- 端口范围：0-65535
- 超时时间：非负

如果未设置，会自动应用默认值。

## 高级用法

### 配置热重载

支持运行时动态更新配置，不中断服务：

```go
// 创建新配置
newConfig := &httpserver.Config{
    Host:         "0.0.0.0",
    Port:         8081,
    ReadTimeout:  20 * time.Second,
    WriteTimeout: 20 * time.Second,
    IdleTimeout:  90 * time.Second,
}

// 热重载配置
if err := server.Reload(context.Background(), newConfig); err != nil {
    log.Error("failed to reload config", "error", err)
}
```

**热重载行为**：

- **端口未变化**: 启动新服务器 → 关闭旧服务器（无缝切换）
- **端口变化**: 关闭旧服务器 → 启动新服务器（短暂中断）

### 自动端口分配

```go
config := &httpserver.Config{
    Port: 0, // 0 表示自动分配
}

server, _ := httpserver.New(router, config, log)
server.Start(context.Background())
// 服务器会自动分配 9000-30000 范围内的可用端口
```

### 与 DI 容器集成

在应用初始化时集成到 DI 容器：

```go
// internal/app/app_httpserver.go
func initHTTPServer(app *App) error {
    cfg := &httpserver.Config{
        Host:         app.Config.Server.Host,
        Port:         app.Config.Server.Port,
        ReadTimeout:  time.Duration(app.Config.Server.ReadTimeout) * time.Second,
        WriteTimeout: time.Duration(app.Config.Server.WriteTimeout) * time.Second,
        IdleTimeout:  time.Duration(app.Config.Server.IdleTimeout) * time.Second,
    }

    server, err := httpserver.New(app.Router, cfg, app.Logger)
    if err != nil {
        return fmt.Errorf("failed to create HTTP server: %w", err)
    }

    app.HTTPServer = server
    return nil
}
```

## 故障排查

### 端口已被占用

**症状**: 启动失败，错误信息包含 "address already in use"

**解决方案**:

1. 使用不同的端口
2. 使用 `Port: 0` 自动分配端口
3. 检查并关闭占用端口的进程

### 热重载失败

**症状**: 调用 `Reload()` 返回错误

**可能原因**:

- 新配置无效（端口超出范围、超时为负数等）
- 新端口已被占用
- 服务器正在关闭中

**解决方案**:

- 检查新配置的有效性
- 确保新端口可用
- 等待服务器完全启动后再重载

### 优雅关闭超时

**症状**: 关闭时超过 context 超时时间

**可能原因**:

- 有长时间运行的请求
- IdleTimeout 设置过大

**解决方案**:

- 增加 Shutdown context 超时时间
- 优化长时间运行的请求
- 调整 IdleTimeout 配置

## 最佳实践

### 1. 超时配置

```go
// 开发环境：宽松的超时
config := &httpserver.Config{
    ReadTimeout:  30 * time.Second,
    WriteTimeout: 30 * time.Second,
    IdleTimeout:  2 * time.Minute,
}

// 生产环境：严格的超时
config := &httpserver.Config{
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  60 * time.Second,
}
```

### 2. 错误处理

```go
// 启动失败应该终止程序
if err := server.Start(ctx); err != nil {
    log.Fatal("critical: server failed to start", "error", err)
}

// 关闭失败记录日志但继续清理
if err := server.Shutdown(ctx); err != nil {
    log.Error("shutdown error", "error", err)
    // 继续清理其他资源
}

// 热更新失败保持原配置运行
if err := server.Reload(ctx, newConfig); err != nil {
    log.Error("reload failed, keeping old config", "error", err)
    // 服务器继续使用旧配置运行
}
```

### 3. 并发使用

```go
// HTTPServer 的所有方法都是线程安全的
go server.Start(ctx)
go server.Reload(ctx, newConfig)
// 不会发生竞态条件
```

## 许可证

本项目遵循现有项目许可证。
