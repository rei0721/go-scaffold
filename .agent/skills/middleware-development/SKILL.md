---
name: middleware-development
description: 在 internal/middleware/ 目录下创建 Gin 中间件
---

# 中间件开发规范

## 概述

本 skill 指导在 `internal/middleware/` 目录下创建符合项目规范的 Gin 中间件。

## 文件结构

```
internal/middleware/
├── {middleware}.go    # 中间件实现
└── constants.go       # 常量定义
```

## 中间件模式

### 基础模式

```go
package middleware

import "github.com/gin-gonic/gin"

// {Name}Middleware 返回 {功能} 中间件
func {Name}Middleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 前置处理
        // ...

        // 2. 调用下一个处理器
        c.Next()

        // 3. 后置处理
        // ...
    }
}
```

### 带依赖的中间件

```go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/rei0721/go-scaffold/pkg/jwt"
    "github.com/rei0721/go-scaffold/pkg/logger"
)

// AuthMiddleware 返回认证中间件
func AuthMiddleware(jwtManager jwt.JWT, log logger.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从 Header 获取 Token
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }

        // 2. 验证 Token
        claims, err := jwtManager.ValidateToken(token)
        if err != nil {
            log.Warn("invalid token", "error", err)
            c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
            return
        }

        // 3. 将用户信息存入 Context
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)

        c.Next()
    }
}
```

## 常用中间件示例

### 日志中间件

```go
func LoggerMiddleware(log logger.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path

        c.Next()

        latency := time.Since(start)
        log.Info("request completed",
            "method", c.Request.Method,
            "path", path,
            "status", c.Writer.Status(),
            "latency", latency,
        )
    }
}
```

### 恢复中间件

```go
func RecoveryMiddleware(log logger.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                log.Error("panic recovered",
                    "error", err,
                    "stack", string(debug.Stack()),
                )
                c.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
            }
        }()
        c.Next()
    }
}
```

### TraceID 中间件

```go
func TraceIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        traceID := c.GetHeader("X-Trace-ID")
        if traceID == "" {
            traceID = uuid.New().String()
        }

        c.Set("trace_id", traceID)
        c.Header("X-Trace-ID", traceID)

        c.Next()
    }
}
```

### RBAC 权限中间件

```go
func RBACMiddleware(rbacManager rbac.RBAC, log logger.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("user_id")
        if !exists {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }

        // 检查权限
        subject := fmt.Sprintf("user:%d", userID)
        object := c.Request.URL.Path
        action := c.Request.Method

        allowed, err := rbacManager.Enforce(subject, object, action)
        if err != nil || !allowed {
            log.Warn("access denied",
                "user_id", userID,
                "path", object,
                "action", action,
            )
            c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
            return
        }

        c.Next()
    }
}
```

## 注册中间件

### 全局中间件

```go
func SetupRouter(router *gin.Engine, log logger.Logger) {
    // 全局中间件
    router.Use(RecoveryMiddleware(log))
    router.Use(TraceIDMiddleware())
    router.Use(LoggerMiddleware(log))
}
```

### 路由组中间件

```go
api := router.Group("/api/v1")
api.Use(AuthMiddleware(jwtManager, log))
api.Use(RBACMiddleware(rbacManager, log))
```

## 检查清单

- [ ] 返回 `gin.HandlerFunc`
- [ ] 使用 `c.Next()` 调用后续处理器
- [ ] 使用 `c.Abort*` 中断请求链
- [ ] 使用 `c.Set` 传递数据到下游
- [ ] 依赖通过函数参数注入
- [ ] 错误使用 `c.AbortWithStatusJSON`
- [ ] 在 `internal/router/` 中注册
