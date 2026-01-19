---
name: pkg-development
description: 在 pkg/ 目录下创建新的可复用工具包
---

# pkg 工具包开发规范

## 概述

本 skill 指导在 `pkg/` 目录下创建符合项目规范的可复用工具包。

## 目录结构

```
pkg/{package-name}/
├── {package}.go      # 接口定义（必需）
├── {package}_impl.go # 具体实现（必需）
├── config.go         # 配置结构（可选）
├── constants.go      # 常量定义（可选）
├── errors.go         # 错误定义（可选）
├── doc.go            # 包文档（推荐）
└── README.md         # 使用说明（推荐）
```

## 开发步骤

### 1. 创建接口文件 `{package}.go`

```go
// Package {name} 提供 {功能描述}
// 设计目标:
// - 提供统一接口,屏蔽具体实现
// - 支持依赖注入和测试 mock
// - 便于切换实现
package {name}

// {Name} 定义 {功能} 的接口
type {Name} interface {
    // Method 方法说明
    // 参数:
    //   ctx: 上下文
    // 返回:
    //   error: 错误描述
    Method(ctx context.Context) error
    
    // Reloader 嵌入热重载接口（如需支持）
    Reloader
}

// Reloader 定义配置重载接口
type Reloader interface {
    Reload(cfg *Config) error
}
```

### 2. 创建配置结构 `config.go`

```go
package {name}

// Config 保存 {功能} 配置
type Config struct {
    // Field 字段说明
    Field string `mapstructure:"field"`
}

// ValidateName 返回配置名称
func (c *Config) ValidateName() string {
    return "{name}"
}

// Validate 验证配置有效性
func (c *Config) Validate() error {
    // 配置验证逻辑
    return nil
}

// DefaultConfig 返回默认配置
func (c *Config) DefaultConfig() {
    // 设置默认值
}

// OverrideConfig 从环境变量覆盖配置
func (c *Config) OverrideConfig() {
    // 环境变量覆盖逻辑
}
```

### 3. 创建实现文件 `{package}_impl.go`

```go
package {name}

import (
    "sync"
    "sync/atomic"
    
    "github.com/rei0721/go-scaffold/pkg/executor"
)

// impl 是 {Name} 接口的具体实现
type impl struct {
    config   *Config
    mu       sync.RWMutex      // 保护并发访问
    executor atomic.Value      // 延迟注入的 executor
}

// New 创建新的 {Name} 实例
func New(cfg *Config) ({Name}, error) {
    if err := cfg.Validate(); err != nil {
        return nil, err
    }
    
    i := &impl{
        config: cfg,
    }
    
    return i, nil
}

// Method 实现接口方法
func (i *impl) Method(ctx context.Context) error {
    i.mu.RLock()
    defer i.mu.RUnlock()
    // 实现逻辑
    return nil
}

// Reload 热重载配置
func (i *impl) Reload(cfg *Config) error {
    if err := cfg.Validate(); err != nil {
        return err
    }
    
    i.mu.Lock()
    defer i.mu.Unlock()
    
    // 原子替换配置
    i.config = cfg
    return nil
}

// SetExecutor 设置协程池管理器（延迟注入）
func (i *impl) SetExecutor(exec executor.Manager) {
    i.executor.Store(exec)
}

// getExecutor 获取协程池管理器
func (i *impl) getExecutor() executor.Manager {
    if exec := i.executor.Load(); exec != nil {
        return exec.(executor.Manager)
    }
    return nil
}
```

### 4. 创建常量文件 `constants.go`

```go
package {name}

const (
    // DefaultValue 默认值说明
    DefaultValue = "default"
)
```

### 5. 创建错误文件 `errors.go`（如需要）

```go
package {name}

import "errors"

var (
    // ErrInvalidConfig 无效配置错误
    ErrInvalidConfig = errors.New("{name}: invalid configuration")
)
```

### 6. 创建包文档 `doc.go`

```go
// Package {name} 提供 {功能描述}
//
// 设计目标:
// - 目标1
// - 目标2
//
// 使用示例:
//
//     cfg := &{name}.Config{...}
//     instance, err := {name}.New(cfg)
//     if err != nil {
//         log.Fatal(err)
//     }
//     defer instance.Close()
//
// 接口设计:
// - {Name}: 主接口
// - Reloader: 热重载接口
//
// 线程安全:
// - 所有方法都是并发安全的
// - 使用 sync.RWMutex 保护内部状态
package {name}
```

### 7. 创建 README.md

```markdown
# {Package Name}

## 功能

{功能描述}

## 安装

```go
import "github.com/rei0721/go-scaffold/pkg/{name}"
```

## 使用

```go
cfg := &{name}.Config{...}
instance, err := {name}.New(cfg)
```

## 配置

| 字段 | 类型 | 默认值 | 说明 |
|-----|------|-------|------|
| Field | string | "" | 字段说明 |

## 接口

- `{Name}`: 主接口
- `Reloader`: 热重载接口
```

## 集成到 App 容器

在 `internal/app/` 下创建 `app_{name}.go`:

```go
package app

func (a *App) init{Name}() error {
    cfg := &{name}.Config{
        // 从 a.Config 读取配置
    }
    
    instance, err := {name}.New(cfg)
    if err != nil {
        return err
    }
    
    a.{Name} = instance
    return nil
}
```

## 检查清单

- [ ] 接口定义在 `{package}.go`
- [ ] 实现在 `{package}_impl.go`
- [ ] 支持 `Reloader` 接口（如需热重载）
- [ ] 使用 `sync.RWMutex` 保护并发访问
- [ ] 延迟注入使用 `atomic.Value`
- [ ] 配置有 `Validate()`、`DefaultConfig()`、`OverrideConfig()` 方法
- [ ] 有 `doc.go` 包文档
- [ ] 有 `README.md` 使用说明
- [ ] 在 `internal/app/` 中集成
