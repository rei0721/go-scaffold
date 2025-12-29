# API 规范文档

```
项目: Rei0721 | 版本: v1.0 | 更新: 2025-12-29
```

## 设计原则

| 原则 | 说明 |
|------|------|
| 一致性 | 统一格式和规范 |
| 可预测性 | 相同输入 = 相同输出 |
| 简洁性 | 避免冗余 |
| 安全性 | 验证所有输入 |
| 国际化 | 多语言错误消息 |
| 可观测性 | TraceID 链路追踪 |

---

## 统一响应格式

```go
type Response[T any] struct {
    Code       int    `json:"code"`
    Message    string `json:"message"`
    Data       T      `json:"data,omitempty"`
    TraceID    string `json:"traceId,omitempty"`
    ServerTime int64  `json:"serverTime"`
}
```

### 响应示例

**成功:**
```json
{"code": 0, "message": "success", "data": {"id": 1001, "username": "alice"}, "traceId": "abc123", "serverTime": 1735475400}
```

**错误:**
```json
{"code": 1001, "message": "用户名不能为空", "traceId": "abc123", "serverTime": 1735475400}
```

**分页:**
```json
{
  "code": 0, "message": "success",
  "data": {
    "list": [{"id": 1, "name": "Item 1"}],
    "pagination": {"page": 1, "pageSize": 10, "total": 100, "totalPages": 10}
  }
}
```

---

## 错误码定义

| 范围 | 类型 | 说明 |
|------|------|------|
| 0 | 成功 | 请求成功 |
| 1000-1999 | 参数错误 | 请求参数验证失败 |
| 2000-2999 | 业务错误 | 业务逻辑错误 |
| 3000-3999 | 认证错误 | 认证/授权相关 |
| 4000-4999 | 资源错误 | 资源不存在/已占用 |
| 5000-5999 | 系统错误 | 服务器内部错误 |

### 常用错误码

```go
// types/errors/codes.go
const (
    CodeSuccess = 0

    // 参数错误 1000-1999
    ErrInvalidParams   = 1000
    ErrInvalidUsername = 1001
    ErrInvalidEmail    = 1002
    ErrInvalidPassword = 1003

    // 业务错误 2000-2999
    ErrBusinessLogic     = 2000
    ErrDuplicateUsername = 2001
    ErrDuplicateEmail    = 2002

    // 认证错误 3000-3999
    ErrUnauthorized      = 3000
    ErrInvalidToken      = 3001
    ErrTokenExpired      = 3002
    ErrPermissionDenied  = 3003

    // 资源错误 4000-4999
    ErrResourceNotFound = 4000
    ErrUserNotFound     = 4001

    // 系统错误 5000-5999
    ErrInternalServer = 5000
    ErrDatabaseError  = 5001
    ErrCacheError     = 5002
)
```

---

## 请求规范

### HTTP 方法

| 方法 | 用途 | 幂等 |
|------|------|------|
| GET | 查询 | ✅ |
| POST | 创建 | ❌ |
| PUT | 完整更新 | ✅ |
| PATCH | 部分更新 | ❌ |
| DELETE | 删除 | ✅ |

### 请求头

```http
Content-Type: application/json
Accept: application/json
Accept-Language: zh-CN,en-US
X-Request-ID: <uuid>
Authorization: Bearer <token>
```

### 参数位置

| 位置 | 场景 | 示例 |
|------|------|------|
| Path | 资源标识 | `/api/v1/users/{userId}` |
| Query | 过滤/排序/分页 | `?page=1&pageSize=10` |
| Body | 创建/更新数据 | `{"username": "alice"}` |
| Header | 认证/元数据 | `Authorization: Bearer token` |

---

## RESTful 设计

### 命名规范

- ✅ 复数名词: `/api/v1/users`
- ✅ 小写字母: `/api/v1/users`
- ✅ 连字符: `/api/v1/user-profiles`
- ❌ 避免动词 (用 HTTP 方法表达)

### 标准 CRUD

| 操作 | 方法 | URI |
|------|------|-----|
| 创建 | POST | `/api/v1/users` |
| 列表 | GET | `/api/v1/users` |
| 详情 | GET | `/api/v1/users/{id}` |
| 更新 | PUT/PATCH | `/api/v1/users/{id}` |
| 删除 | DELETE | `/api/v1/users/{id}` |

### 嵌套资源 (≤2层)

```
GET  /api/v1/users/{userId}/orders
GET  /api/v1/users/{userId}/orders/{orderId}
POST /api/v1/users/{userId}/orders
```

---

## 版本管理

```
/api/v1/users    # 版本 1
/api/v2/users    # 版本 2
```

- 主版本号: 重大变更时递增
- 向后兼容改动不增加版本
- 弃用 API 保留至少一个版本周期

---

## API 示例

### 用户注册

```http
POST /api/v1/users/register
Content-Type: application/json

{"username": "alice", "email": "alice@example.com", "password": "SecurePass123!"}
```

```json
{
  "code": 0, "message": "success",
  "data": {"userId": 1001, "username": "alice", "createdAt": "2025-12-29T21:30:00Z"}
}
```

### 用户登录

```http
POST /api/v1/users/login
Content-Type: application/json

{"username": "alice", "password": "SecurePass123!"}
```

```json
{
  "code": 0, "message": "success",
  "data": {"token": "eyJhbGci...", "expiresIn": 7200, "user": {"userId": 1001}}
}
```

### 获取用户

```http
GET /api/v1/users/1001
Authorization: Bearer eyJhbGci...
```

```json
{
  "code": 0, "message": "success",
  "data": {"userId": 1001, "username": "alice", "email": "alice@example.com"}
}
```

### 用户列表 (分页)

```http
GET /api/v1/users?page=1&pageSize=10&sort=createdAt&order=desc
Authorization: Bearer eyJhbGci...
```

```json
{
  "code": 0, "message": "success",
  "data": {
    "list": [{"userId": 1001, "username": "alice"}],
    "pagination": {"page": 1, "pageSize": 10, "total": 100}
  }
}
```

### 错误响应

```json
{"code": 2001, "message": "用户名已存在", "traceId": "pqr678", "serverTime": 1735475900}
```

---

## 测试清单

- [ ] 正常流程
- [ ] 参数验证 (缺失/格式错误)
- [ ] 认证测试 (无Token/过期Token)
- [ ] 权限测试
- [ ] 并发测试
- [ ] 边界值测试
- [ ] 国际化测试

---

[← README](./README.md) | [protocol.md](./protocol.md) | [deployment.md →](./deployment.md)
