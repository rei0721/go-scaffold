# 认证说明

## 认证方式

本API使用 **JWT (JSON Web Token)** 进行认证。

## 获取 Token

### 1. 用户注册

```bash
POST /api/v1/auth/register
```

### 2. 用户登录

```bash
POST /api/v1/auth/login
```

成功登录后会返回 access token:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresIn": 3600
  }
}
```

## 使用 Token

在需要认证的接口请求头中添加 `Authorization` 字段:

```
Authorization: Bearer {token}
```

**示例:**

```bash
curl -X GET http://localhost:9999/api/v1/protected-endpoint \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## Token 有效期

- **默认有效期:** 1小时 (3600秒)
- **配置位置:** `configs/config.yaml` 中的 `jwt.expiresIn`

## Token 刷新

当 token 即将过期时,可以使用 refresh token 获取新的 access token:

```bash
POST /api/v1/auth/refresh
```

## 权限控制

某些接口除了需要认证外,还需要特定的角色或权限。详见各接口文档的"认证"部分。

### 权限等级

- **公开接口** - 无需认证
- **认证接口** - 需要有效的 JWT token
- **角色接口** - 需要特定角色(如 `admin`, `user`)
- **权限接口** - 需要特定权限(如对 `users` 资源的 `write` 权限)

## 常见认证错误

| 状态码 | code | message                  | 说明             |
| ------ | ---- | ------------------------ | ---------------- |
| 401    | 1002 | Missing authorization    | 缺少认证头       |
| 401    | 1002 | Invalid authorization    | 认证格式错误     |
| 401    | 1002 | Invalid or expired token | Token 无效或过期 |
| 403    | 1003 | Forbidden                | 权限不足         |
