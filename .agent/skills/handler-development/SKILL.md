---
name: handler-development
description: 在 internal/handler/ 目录下创建 HTTP 处理器
---

# Handler 开发规范

## 概述

本 skill 指导在 `internal/handler/` 目录下创建符合项目规范的 HTTP 处理器。

## 文件结构

```
internal/handler/
├── {feature}_handler.go  # Handler 实现
└── constants.go          # 常量定义
```

## 开发步骤

### 1. 创建 Handler 结构体

```go
package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/rei0721/go-scaffold/internal/service/{service}"
    "github.com/rei0721/go-scaffold/pkg/logger"
    "github.com/rei0721/go-scaffold/types"
    "github.com/rei0721/go-scaffold/types/result"
)

// {Feature}Handler 处理 {功能} 相关的 HTTP 请求
type {Feature}Handler struct {
    service {service}.{Service}Service
    logger  logger.Logger
}

// New{Feature}Handler 创建新的 Handler 实例
func New{Feature}Handler(svc {service}.{Service}Service, log logger.Logger) *{Feature}Handler {
    return &{Feature}Handler{
        service: svc,
        logger:  log,
    }
}
```

### 2. 实现处理方法

```go
// Create 创建资源
// @Summary 创建 {资源}
// @Tags {Feature}
// @Accept json
// @Produce json
// @Param request body types.CreateRequest true "创建请求"
// @Success 200 {object} result.Result{data=types.Response}
// @Router /api/v1/{feature} [post]
func (h *{Feature}Handler) Create(c *gin.Context) {
    // 1. 绑定请求
    var req types.CreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, result.Fail(400, "invalid request: "+err.Error()))
        return
    }

    // 2. 调用服务
    resp, err := h.service.Create(c.Request.Context(), &req)
    if err != nil {
        h.logger.Error("failed to create", "error", err)
        c.JSON(500, result.Fail(500, err.Error()))
        return
    }

    // 3. 返回成功响应
    c.JSON(200, result.Success(resp))
}

// GetByID 获取单个资源
func (h *{Feature}Handler) GetByID(c *gin.Context) {
    // 1. 获取路径参数
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(400, result.Fail(400, "invalid id"))
        return
    }

    // 2. 调用服务
    resp, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(500, result.Fail(500, err.Error()))
        return
    }
    if resp == nil {
        c.JSON(404, result.Fail(404, "not found"))
        return
    }

    c.JSON(200, result.Success(resp))
}

// List 获取列表（分页）
func (h *{Feature}Handler) List(c *gin.Context) {
    // 1. 获取分页参数
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

    // 2. 调用服务
    items, total, err := h.service.List(c.Request.Context(), page, pageSize)
    if err != nil {
        c.JSON(500, result.Fail(500, err.Error()))
        return
    }

    // 3. 返回分页结果
    c.JSON(200, result.SuccessWithPagination(items, page, pageSize, total))
}

// Update 更新资源
func (h *{Feature}Handler) Update(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(400, result.Fail(400, "invalid id"))
        return
    }

    var req types.UpdateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, result.Fail(400, "invalid request: "+err.Error()))
        return
    }

    if err := h.service.Update(c.Request.Context(), id, &req); err != nil {
        c.JSON(500, result.Fail(500, err.Error()))
        return
    }

    c.JSON(200, result.Success(nil))
}

// Delete 删除资源
func (h *{Feature}Handler) Delete(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(400, result.Fail(400, "invalid id"))
        return
    }

    if err := h.service.Delete(c.Request.Context(), id); err != nil {
        c.JSON(500, result.Fail(500, err.Error()))
        return
    }

    c.JSON(200, result.Success(nil))
}
```

### 3. 注册路由

在 `internal/router/router.go` 中：

```go
func SetupRouter(
    router *gin.Engine,
    {feature}Handler *handler.{Feature}Handler,
    authMiddleware gin.HandlerFunc,
) {
    api := router.Group("/api/v1")

    // {Feature} 路由组
    {feature}Group := api.Group("/{feature}")
    {feature}Group.Use(authMiddleware) // 需要认证
    {
        {feature}Group.POST("", {feature}Handler.Create)
        {feature}Group.GET("", {feature}Handler.List)
        {feature}Group.GET("/:id", {feature}Handler.GetByID)
        {feature}Group.PUT("/:id", {feature}Handler.Update)
        {feature}Group.DELETE("/:id", {feature}Handler.Delete)
    }
}
```

## 响应格式

使用 `types/result` 统一响应格式：

```go
// 成功响应
result.Success(data)

// 失败响应
result.Fail(code, message)

// 分页响应
result.SuccessWithPagination(items, page, pageSize, total)
```

响应示例：

```json
{
    "code": 200,
    "message": "success",
    "data": {...}
}
```

## 从 Context 获取用户信息

```go
func (h *{Feature}Handler) Create(c *gin.Context) {
    // 从中间件设置的 context 获取用户 ID
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(401, result.Fail(401, "unauthorized"))
        return
    }

    // 类型断言
    uid := userID.(int64)
    // ...
}
```

## 检查清单

- [ ] Handler 结构体包含 service 和 logger
- [ ] 使用 `New{Feature}Handler` 构造函数
- [ ] 使用 `c.ShouldBindJSON` 绑定请求
- [ ] 使用 `c.Request.Context()` 传递上下文
- [ ] 使用 `result.Success/Fail` 统一响应
- [ ] 路径参数使用 `c.Param`
- [ ] 查询参数使用 `c.Query` 或 `c.DefaultQuery`
- [ ] 在 `internal/router/` 中注册路由
- [ ] 添加必要的中间件保护
