# 认证 API

本文档描述了用户认证相关的 API 接口。

## 概述

认证 API 提供用户注册、登录和登出功能。

- **基础 URL**: `/api/v1/auth`
- **认证方式**: 部分接口需要 JWT Bearer Token

---

## POST /api/v1/auth/register

用户注册接口。

### 请求

**URL**: `POST /api/v1/auth/register`

**认证**: 不需要

**请求体**:

```json
{
  "username": "user123",
  "email": "user@example.com",
  "password": "password123"
}
```

| 字段     | 类型   | 必填 | 说明                           |
| -------- | ------ | ---- | ------------------------------ |
| username | string | 是   | 用户名，3-50个字符             |
| email    | string | 是   | 邮箱地址，必须是有效的邮箱格式 |
| password | string | 是   | 密码，最少8个字符              |

### 响应

**成功响应** (200 OK):

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "userId": 123,
    "username": "user123",
    "email": "user@example.com",
    "status": 1,
    "createdAt": "2024-01-01T00:00:00Z"
  },
  "serverTime": 1640000000
}
```

**错误响应**:

- **400 Bad Request**: 请求参数错误

  ```json
  {
    "code": 1001,
    "message": "无效的请求参数: ...",
    "serverTime": 1640000000
  }
  ```

- **500 Internal Server Error**: 服务器内部错误
  ```json
  {
    "code": 5000,
    "message": "注册失败: ...",
    "serverTime": 1640000000
  }
  ```

### 示例

```bash
curl -X POST http://localhost:9999/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user123",
    "email": "user@example.com",
    "password": "password123"
  }'
```

---

## POST /api/v1/auth/login

用户登录接口。

### 请求

**URL**: `POST /api/v1/auth/login`

**认证**: 不需要

**请求体**:

```json
{
  "username": "user123",
  "password": "password123"
}
```

| 字段     | 类型   | 必填 | 说明   |
| -------- | ------ | ---- | ------ |
| username | string | 是   | 用户名 |
| password | string | 是   | 密码   |

### 响应

**成功响应** (200 OK):

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresIn": 3600,
    "user": {
      "userId": 123,
      "username": "user123",
      "email": "user@example.com",
      "status": 1,
      "createdAt": "2024-01-01T00:00:00Z"
    }
  },
  "serverTime": 1640000000
}
```

| 字段      | 类型   | 说明             |
| --------- | ------ | ---------------- |
| token     | string | JWT 访问令牌     |
| expiresIn | number | 令牌有效期（秒） |
| user      | object | 用户信息         |

**错误响应**:

- **400 Bad Request**: 请求参数错误
- **401 Unauthorized**: 用户名或密码错误
  ```json
  {
    "code": 4010,
    "message": "用户名或密码错误",
    "serverTime": 1640000000
  }
  ```

### 示例

```bash
curl -X POST http://localhost:9999/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user123",
    "password": "password123"
  }'
```

### 使用返回的 Token

登录成功后，将返回的 `token` 放在后续请求的 `Authorization` 头中：

```bash
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  http://localhost:9999/api/v1/auth/logout
```

---

## POST /api/v1/auth/logout

用户登出接口。

### 请求

**URL**: `POST /api/v1/auth/logout`

**认证**: 需要（JWT Bearer Token）

**Headers**:

```
Authorization: Bearer <token>
```

**请求体**: 无

### 响应

**成功响应** (200 OK):

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "登出成功"
  },
  "serverTime": 1640000000
}
```

**错误响应**:

- **401 Unauthorized**: 未认证或 Token 无效

  ```json
  {
    "code": 4010,
    "message": "未认证",
    "serverTime": 1640000000
  }
  ```

- **500 Internal Server Error**: 服务器内部错误
  ```json
  {
    "code": 5000,
    "message": "登出失败: ...",
    "serverTime": 1640000000
  }
  ```

### 示例

```bash
curl -X POST http://localhost:9999/api/v1/auth/logout \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## POST /api/v1/auth/change-password

修改密码接口。

### 请求

**URL**: `POST /api/v1/auth/change-password`

**认证**: 需要（JWT Bearer Token）

**Headers**:

```
Authorization: Bearer <token>
```

**请求体**:

```json
{
  "old_password": "oldpass123",
  "new_password": "newpass123"
}
```

| 字段         | 类型   | 必填 | 说明                |
| ------------ | ------ | ---- | ------------------- |
| old_password | string | 是   | 原密码              |
| new_password | string | 是   | 新密码，最少8个字符 |

### 响应

**成功响应** (200 OK):

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "密码修改成功"
  },
  "serverTime": 1640000000
}
```

**错误响应**:

- **400 Bad Request**: 请求参数错误
- **401 Unauthorized**: 未认证或原密码错误
  ```json
  {
    "code": 4010,
    "message": "未认证",
    "serverTime": 1640000000
  }
  ```
- **500 Internal Server Error**: 服务器内部错误

### 示例

```bash
curl -X POST http://localhost:9999/api/v1/auth/change-password \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "oldpass123",
    "new_password": "newpass123"
  }'
```

### 注意事项

- 修改密码后，原有的 Token 仍然有效，直到过期
- 建议修改密码后重新登录
- 服务器会清除缓存中的用户信息

---

## POST /api/v1/auth/refresh

刷新访问令牌接口。

### 请求

**URL**: `POST /api/v1/auth/refresh`

**认证**: 不需要（但需要有效的 refresh token）

**请求体**:

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

| 字段          | 类型   | 必填 | 说明     |
| ------------- | ------ | ---- | -------- |
| refresh_token | string | 是   | 刷新令牌 |

### 响应

**成功响应** (200 OK):

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600,
    "token_type": "Bearer"
  },
  "serverTime": 1640000000
}
```

| 字段          | 类型   | 说明                       |
| ------------- | ------ | -------------------------- |
| access_token  | string | 新的访问令牌               |
| refresh_token | string | 原 refresh token（或新的） |
| expires_in    | number | 访问令牌有效期（秒）       |
| token_type    | string | 令牌类型，通常为 "Bearer"  |

**错误响应**:

- **400 Bad Request**: 请求参数错误
- **401 Unauthorized**: refresh token 无效或过期
  ```json
  {
    "code": 4010,
    "message": "refresh token 无效或过期",
    "serverTime": 1640000000
  }
  ```
- **500 Internal Server Error**: 服务器内部错误

### 示例

```bash
curl -X POST http://localhost:9999/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

### 使用场景

当访问令牌（access token）过期时，前端可以使用 refresh token 获取新的访问令牌，而无需用户重新登录。

### 注意事项

- Refresh token 的有效期通常较长（例如 7 天）
- 当前实现保持原 refresh token 不变
- 可以实现 refresh token rotation 提高安全性

---

## 错误码

| 错误码 | 说明            |
| ------ | --------------- |
| 0      | 成功            |
| 1001   | 请求参数错误    |
| 4010   | 未授权/认证失败 |
| 5000   | 服务器内部错误  |

更多错误码请参考 [错误码说明](../error-codes.md)。

---

## API 列表

| 方法 | 路径                         | 认证 | 说明     |
| ---- | ---------------------------- | ---- | -------- |
| POST | /api/v1/auth/register        | ❌   | 用户注册 |
| POST | /api/v1/auth/login           | ❌   | 用户登录 |
| POST | /api/v1/auth/logout          | ✅   | 用户登出 |
| POST | /api/v1/auth/change-password | ✅   | 修改密码 |
| POST | /api/v1/auth/refresh         | ❌   | 刷新令牌 |

**图例**: ✅ = 需要认证 | ❌ = 公开接口

---

## 注意事项

1. **密码安全**: 密码在传输和存储时都会被加密处理
2. **Token 有效期**: JWT Token 默认有效期为 1 小时
3. **并发登录**: 同一账号可以在多个设备同时登录
4. **登出行为**: 登出会清除服务器端的用户缓存，但不会使 Token 立即失效
5. **密码修改**: 修改密码后建议重新登录以确保安全
6. **Token 刷新**: 使用 refresh token 可以在 access token 过期时获取新的令牌，无需重新登录
