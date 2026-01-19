---
name: config-integration
description: 配置集成到应用的完整流程
---

# 配置集成规范

## 概述

本 skill 指导如何将新组件的配置集成到应用中，包括配置定义、环境变量覆盖和 App 容器集成。

## 配置集成流程

```
1. internal/config/app_{component}.go    # 定义配置结构
2. configs/config.yaml                    # 添加配置项
3. internal/app/app_{component}.go        # 集成到 App
```

## 开发步骤

### 1. 定义配置结构

在 `internal/config/` 创建 `app_{component}.go`：

```go
package config

// {Component}Config 保存 {组件} 配置
type {Component}Config struct {
    // Enabled 是否启用
    Enabled bool `mapstructure:"enabled" json:"enabled" yaml:"enabled" toml:"enabled"`

    // Host 主机地址
    Host string `mapstructure:"host" json:"host" yaml:"host" toml:"host"`

    // Port 端口
    Port int `mapstructure:"port" json:"port" yaml:"port" toml:"port"`
}

// ValidateName 返回配置名称
func (c *{Component}Config) ValidateName() string {
    return "{component}"
}

// Validate 验证配置有效性
func (c *{Component}Config) Validate() error {
    if c.Enabled && c.Host == "" {
        return fmt.Errorf("{component}: host is required when enabled")
    }
    if c.Port < 0 || c.Port > 65535 {
        return fmt.Errorf("{component}: invalid port %d", c.Port)
    }
    return nil
}

// DefaultConfig 设置默认配置
func (c *{Component}Config) DefaultConfig() {
    if c.Host == "" {
        c.Host = "localhost"
    }
    if c.Port == 0 {
        c.Port = 8080
    }
}

// OverrideConfig 从环境变量覆盖配置
func (c *{Component}Config) OverrideConfig() {
    if host := os.Getenv("{COMPONENT}_HOST"); host != "" {
        c.Host = host
    }
    if port := os.Getenv("{COMPONENT}_PORT"); port != "" {
        if p, err := strconv.Atoi(port); err == nil {
            c.Port = p
        }
    }
}
```

### 2. 集成到主配置

在 `internal/config/config.go` 添加字段：

```go
type Config struct {
    // 现有配置...
    App      AppConfig      `mapstructure:"app"`
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`

    // 新增配置
    {Component} {Component}Config `mapstructure:"{component}"`
}
```

### 3. 添加配置文件

在 `configs/config.yaml` 添加：

```yaml
{ component }:
  enabled: true
  host: "localhost"
  port: 8080
```

### 4. 创建 App 集成

在 `internal/app/` 创建 `app_{component}.go`：

```go
package app

import (
    "github.com/rei0721/go-scaffold/pkg/{component}"
)

// init{Component} 初始化 {组件}
func (a *App) init{Component}() error {
    cfg := a.Config.{Component}

    // 应用默认值
    cfg.DefaultConfig()

    // 从环境变量覆盖
    cfg.OverrideConfig()

    // 验证配置
    if err := cfg.Validate(); err != nil {
        return fmt.Errorf("invalid {component} config: %w", err)
    }

    // 如果未启用，跳过初始化
    if !cfg.Enabled {
        a.Logger.Info("{component} is disabled")
        return nil
    }

    // 创建 pkg 配置
    pkgCfg := &{component}.Config{
        Host: cfg.Host,
        Port: cfg.Port,
    }

    // 初始化组件
    instance, err := {component}.New(pkgCfg)
    if err != nil {
        return fmt.Errorf("failed to create {component}: %w", err)
    }

    // 注入依赖（如需要）
    if injector, ok := instance.(types.ExecutorInjectable); ok {
        injector.SetExecutor(a.Executor)
    }

    a.{Component} = instance
    a.Logger.Info("{component} initialized", "host", cfg.Host, "port", cfg.Port)

    return nil
}
```

### 5. 在 App.New 中调用

在 `internal/app/app.go` 的初始化序列中添加：

```go
func New(opts Options) (*App, error) {
    // 现有初始化...

    // 初始化新组件
    if err := app.init{Component}(); err != nil {
        return nil, err
    }

    return app, nil
}
```

### 6. 添加热重载支持（可选）

在 `internal/app/reload.go`：

```go
func (a *App) onConfigReload() {
    // 重载 {Component}
    if a.{Component} != nil {
        if reloader, ok := a.{Component}.({component}.Reloader); ok {
            cfg := a.Config.{Component}
            cfg.DefaultConfig()
            cfg.OverrideConfig()

            pkgCfg := &{component}.Config{...}
            if err := reloader.Reload(pkgCfg); err != nil {
                a.Logger.Error("failed to reload {component}", "error", err)
            }
        }
    }
}
```

## 环境变量命名规范

| 配置项              | 环境变量            |
| ------------------- | ------------------- |
| `component.host`    | `COMPONENT_HOST`    |
| `component.port`    | `COMPONENT_PORT`    |
| `database.password` | `DATABASE_PASSWORD` |

## 检查清单

- [ ] 配置结构定义在 `internal/config/app_{component}.go`
- [ ] 实现 `ValidateName()`、`Validate()`、`DefaultConfig()`、`OverrideConfig()`
- [ ] 在 `Config` 主结构中添加字段
- [ ] 在 `configs/config.yaml` 添加配置项
- [ ] 创建 `internal/app/app_{component}.go` 初始化函数
- [ ] 在 `App.New` 中调用初始化函数
- [ ] 添加热重载支持（如需要）
- [ ] 敏感配置支持环境变量覆盖
