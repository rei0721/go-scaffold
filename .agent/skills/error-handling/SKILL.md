---
name: error-handling
description: 项目统一的错误处理和错误码规范
---

# 错误处理规范

## 概述

本 skill 指导在项目中使用统一的错误处理模式，包括错误码定义、错误包装和 HTTP 响应。

## 错误码范围

项目在 `types/errors/codes.go` 中定义了分段的错误码：

| 范围      | 类型     | HTTP 状态码 | 说明           |
| --------- | -------- | ----------- | -------------- |
| 0         | 成功     | 200         | 请求成功       |
| 1000-1999 | 参数错误 | 400         | 客户端输入有误 |
| 2000-2999 | 业务错误 | 422         | 业务规则不满足 |
| 3000-3999 | 认证错误 | 401/403     | 身份验证失败   |
| 4000-4999 | 资源错误 | 404         | 资源不存在     |
| 5000-5999 | 系统错误 | 500         | 服务端内部错误 |

## 添加新错误码

### 1. 在 `types/errors/codes.go` 添加常量

```go
const (
    // ==================== 业务错误 (2000-2999) ====================

    // 新增：订单相关错误码
    // ErrOrderNotPaid 订单未支付
    ErrOrderNotPaid = 2010

    // ErrOrderCancelled 订单已取消
    ErrOrderCancelled = 2011
)
```

### 2. 命名规范

| 规范            | 示例                   |
| --------------- | ---------------------- |
| 使用 `Err` 前缀 | `ErrUserNotFound`      |
| 大驼峰命名      | `ErrInvalidPassword`   |
| 描述性名称      | `ErrDuplicateUsername` |

## 错误包装

### 使用 fmt.Errorf 添加上下文

```go
func (s *userService) GetByID(ctx context.Context, id int64) (*User, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        // 使用 %w 包装错误，保留原始错误链
        return nil, fmt.Errorf("failed to get user %d: %w", id, err)
    }
    if user == nil {
        return nil, ErrUserNotFound
    }
    return user, nil
}
```

### 错误检查

```go
// 使用 errors.Is 检查特定错误
if errors.Is(err, ErrUserNotFound) {
    // 处理用户不存在
}

// 使用 errors.As 获取错误类型
var validationErr *ValidationError
if errors.As(err, &validationErr) {
    // 处理验证错误
}
```

## 自定义错误类型

```go
// types/errors/error.go

// BusinessError 业务错误
type BusinessError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func (e *BusinessError) Error() string {
    return e.Message
}

// NewBusinessError 创建业务错误
func NewBusinessError(code int, message string) *BusinessError {
    return &BusinessError{Code: code, Message: message}
}
```

## Handler 层错误处理

### 使用 result 统一响应

```go
func (h *UserHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(400, result.Fail(errors.ErrInvalidParams, "invalid user id"))
        return
    }

    user, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        // 根据错误类型返回不同状态码
        switch {
        case errors.Is(err, ErrUserNotFound):
            c.JSON(404, result.Fail(errors.ErrUserNotFound, "user not found"))
        default:
            h.logger.Error("failed to get user", "error", err)
            c.JSON(500, result.Fail(errors.ErrInternalServer, "internal server error"))
        }
        return
    }

    c.JSON(200, result.Success(user))
}
```

### 错误响应格式

```json
{
  "code": 4001,
  "message": "user not found",
  "data": null
}
```

## Service 层错误处理

```go
func (s *authService) Login(ctx context.Context, req *LoginRequest) (*TokenResponse, error) {
    // 1. 查找用户
    user, err := s.repo.FindByUsername(ctx, req.Username)
    if err != nil {
        return nil, fmt.Errorf("failed to find user: %w", err)
    }
    if user == nil {
        return nil, NewBusinessError(errors.ErrUnauthorized, "invalid credentials")
    }

    // 2. 验证密码
    if !s.crypto.VerifyPassword(req.Password, user.Password) {
        return nil, NewBusinessError(errors.ErrUnauthorized, "invalid credentials")
    }

    // 3. 生成 Token
    token, err := s.jwt.GenerateToken(user.ID)
    if err != nil {
        return nil, fmt.Errorf("failed to generate token: %w", err)
    }

    return &TokenResponse{Token: token}, nil
}
```

## Repository 层错误处理

```go
func (r *userRepository) FindByID(ctx context.Context, id int64) (*User, error) {
    var user User
    err := r.db.WithContext(ctx).First(&user, id).Error

    // 记录不存在返回 nil, nil（不是错误）
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }

    // 其他数据库错误返回 nil, error
    if err != nil {
        return nil, err
    }

    return &user, nil
}
```

## 错误日志规范

```go
// 记录错误日志（包含上下文信息）
log.Error("failed to process request",
    "error", err,
    "user_id", userID,
    "action", "login",
)

// 不要记录敏感信息
log.Error("login failed",
    "username", username,
    // 错误：不要记录密码
    // "password", password,
)
```

## 检查清单

- [ ] 错误码定义在 `types/errors/codes.go`
- [ ] 错误码在正确的范围内
- [ ] 使用 `fmt.Errorf` + `%w` 包装错误
- [ ] 使用 `errors.Is` 检查特定错误
- [ ] Handler 层返回合适的 HTTP 状态码
- [ ] Repository 层：记录不存在返回 `nil, nil`
- [ ] 不在日志中记录敏感信息
- [ ] 系统错误不暴露技术细节给用户
