# RBAC 权限管理 API

本文档描述 RBAC（基于角色的访问控制）权限管理相关的所有 API 接口。

## 认证要求

所有 RBAC API 接口都需要：

- ✅ **JWT 认证**：需要在请求头中携带有效的 JWT Token
- ✅ **Admin 权限**：需要用户具有 `admin` 角色

## 通用说明

### 多租户（Domain）支持

RBAC 支持多租户隔离，通过 `domain` 参数实现不同租户间的权限隔离。

- 不指定 `domain`：权限在全局范围生效
- 指定 `domain`：权限仅在指定域内生效

---

## 角色管理

### POST /api/v1/rbac/users/:id/roles

为指定用户分配单个角色。

#### 认证

- 是否需要认证: **是**
- 需要的角色/权限: `admin`

#### 请求

**路径参数:**

| 参数名 | 类型  | 必填 | 说明   |
| ------ | ----- | ---- | ------ |
| id     | int64 | 是   | 用户ID |

**请求体:**

```json
{
  "role": "string", // 必填，角色名称（如: admin, editor, viewer）
  "domain": "string" // 可选，域名（多租户场景使用）
}
```

#### 响应

**成功响应 (200 OK):**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "Role assigned successfully"
  },
  "serverTime": 1640000000
}
```

**错误响应:**

| 状态码 | code | message              | 说明       |
| ------ | ---- | -------------------- | ---------- |
| 400    | 1001 | Invalid user ID      | 用户ID无效 |
| 400    | 1001 | Invalid request body | 请求体无效 |
| 401    | 1002 | Unauthorized         | 未认证     |
| 403    | 1003 | Forbidden            | 权限不足   |
| 500    | 5000 | Internal error       | 内部错误   |

#### 示例

**请求示例:**

```bash
curl -X POST http://localhost:9999/api/v1/rbac/users/123/roles \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "editor",
    "domain": "tenant1"
  }'
```

**响应示例:**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "Role assigned successfully"
  },
  "serverTime": 1705743600
}
```

---

### POST /api/v1/rbac/users/:id/roles/batch

批量为指定用户分配多个角色。

#### 认证

- 是否需要认证: **是**
- 需要的角色/权限: `admin`

#### 请求

**路径参数:**

| 参数名 | 类型  | 必填 | 说明   |
| ------ | ----- | ---- | ------ |
| id     | int64 | 是   | 用户ID |

**请求体:**

```json
{
  "roles": ["string"], // 必填，角色名称数组
  "domain": "string" // 可选，域名
}
```

#### 响应

**成功响应 (200 OK):**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "Roles assigned successfully"
  },
  "serverTime": 1640000000
}
```

#### 示例

**请求示例:**

```bash
curl -X POST http://localhost:9999/api/v1/rbac/users/123/roles/batch \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "roles": ["editor", "viewer"]
  }'
```

---

### DELETE /api/v1/rbac/users/:id/roles/:role

撤销用户的指定角色。

#### 认证

- 是否需要认证: **是**
- 需要的角色/权限: `admin`

#### 请求

**路径参数:**

| 参数名 | 类型   | 必填 | 说明     |
| ------ | ------ | ---- | -------- |
| id     | int64  | 是   | 用户ID   |
| role   | string | 是   | 角色名称 |

**查询参数:**

| 参数名 | 类型   | 必填 | 说明 |
| ------ | ------ | ---- | ---- |
| domain | string | 否   | 域名 |

#### 响应

**成功响应 (200 OK):**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "Role revoked successfully"
  },
  "serverTime": 1640000000
}
```

#### 示例

**请求示例:**

```bash
# 撤销全局角色
curl -X DELETE http://localhost:9999/api/v1/rbac/users/123/roles/editor \
  -H "Authorization: Bearer YOUR_TOKEN"

# 撤销特定域的角色
curl -X DELETE "http://localhost:9999/api/v1/rbac/users/123/roles/editor?domain=tenant1" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

### GET /api/v1/rbac/users/:id/roles

获取用户拥有的所有角色。

#### 认证

- 是否需要认证: **是**
- 需要的角色/权限: `admin`

#### 请求

**路径参数:**

| 参数名 | 类型  | 必填 | 说明   |
| ------ | ----- | ---- | ------ |
| id     | int64 | 是   | 用户ID |

**查询参数:**

| 参数名 | 类型   | 必填 | 说明 |
| ------ | ------ | ---- | ---- |
| domain | string | 否   | 域名 |

#### 响应

**成功响应 (200 OK):**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 123,
    "roles": ["admin", "editor"]
  },
  "serverTime": 1640000000
}
```

#### 示例

**请求示例:**

```bash
# 获取全局角色
curl -X GET http://localhost:9999/api/v1/rbac/users/123/roles \
  -H "Authorization: Bearer YOUR_TOKEN"

# 获取特定域的角色
curl -X GET "http://localhost:9999/api/v1/rbac/users/123/roles?domain=tenant1" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

### GET /api/v1/rbac/roles/:role/users

获取拥有指定角色的所有用户。

#### 认证

- 是否需要认证: **是**
- 需要的角色/权限: `admin`

#### 请求

**路径参数:**

| 参数名 | 类型   | 必填 | 说明     |
| ------ | ------ | ---- | -------- |
| role   | string | 是   | 角色名称 |

#### 响应

**成功响应 (200 OK):**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "role": "admin",
    "user_ids": [123, 456, 789]
  },
  "serverTime": 1640000000
}
```

#### 示例

**请求示例:**

```bash
curl -X GET http://localhost:9999/api/v1/rbac/roles/admin/users \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 策略管理

### POST /api/v1/rbac/policies

添加单个策略（定义角色对资源的操作权限）。

#### 认证

- 是否需要认证: **是**
- 需要的角色/权限: `admin`

#### 请求

**请求体:**

```json
{
  "role": "string", // 必填，角色名称
  "resource": "string", // 必填，资源名称（如: users, posts）
  "action": "string", // 必填，操作（如: read, write, delete）
  "domain": "string" // 可选，域名
}
```

#### 响应

**成功响应 (200 OK):**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "Policy added successfully"
  },
  "serverTime": 1640000000
}
```

#### 示例

**请求示例:**

```bash
curl -X POST http://localhost:9999/api/v1/rbac/policies \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "editor",
    "resource": "posts",
    "action": "write",
    "domain": "tenant1"
  }'
```

---

### POST /api/v1/rbac/policies/batch

批量添加多个策略。

#### 认证

- 是否需要认证: **是**
- 需要的角色/权限: `admin`

#### 请求

**请求体:**

```json
{
  "policies": [
    {
      "role": "string",
      "resource": "string",
      "action": "string"
    }
  ]
}
```

#### 响应

**成功响应 (200 OK):**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "Policies added successfully"
  },
  "serverTime": 1640000000
}
```

#### 示例

**请求示例:**

```bash
curl -X POST http://localhost:9999/api/v1/rbac/policies/batch \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "policies": [
      {"role": "editor", "resource": "posts", "action": "read"},
      {"role": "editor", "resource": "posts", "action": "write"},
      {"role": "viewer", "resource": "posts", "action": "read"}
    ]
  }'
```

---

### DELETE /api/v1/rbac/policies

删除指定策略。

#### 认证

- 是否需要认证: **是**
- 需要的角色/权限: `admin`

#### 请求

**请求体:**

```json
{
  "role": "string", // 必填，角色名称
  "resource": "string", // 必填，资源名称
  "action": "string", // 必填，操作
  "domain": "string" // 可选，域名
}
```

#### 响应

**成功响应 (200 OK):**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "Policy removed successfully"
  },
  "serverTime": 1640000000
}
```

#### 示例

**请求示例:**

```bash
curl -X DELETE http://localhost:9999/api/v1/rbac/policies \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "editor",
    "resource": "posts",
    "action": "write"
  }'
```

---

### GET /api/v1/rbac/policies

获取所有策略列表。

#### 认证

- 是否需要认证: **是**
- 需要的角色/权限: `admin`

#### 响应

**成功响应 (200 OK):**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "policies": [
      ["admin", "", "users", "write"],
      ["editor", "tenant1", "posts", "write"]
    ],
    "total": 2
  },
  "serverTime": 1640000000
}
```

**说明:**

- 策略数组格式: `[role, domain, resource, action]`
- domain 为空字符串表示全局策略

#### 示例

**请求示例:**

```bash
curl -X GET http://localhost:9999/api/v1/rbac/policies \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

### GET /api/v1/rbac/roles/:role/policies

获取指定角色的所有策略。

#### 认证

- 是否需要认证: **是**
- 需要的角色/权限: `admin`

#### 请求

**路径参数:**

| 参数名 | 类型   | 必填 | 说明     |
| ------ | ------ | ---- | -------- |
| role   | string | 是   | 角色名称 |

#### 响应

**成功响应 (200 OK):**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "policies": [
      ["editor", "", "posts", "read"],
      ["editor", "", "posts", "write"]
    ],
    "total": 2
  },
  "serverTime": 1640000000
}
```

#### 示例

**请求示例:**

```bash
curl -X GET http://localhost:9999/api/v1/rbac/roles/editor/policies \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 权限检查

### POST /api/v1/rbac/check

检查指定用户是否有权限执行某个操作。

#### 认证

- 是否需要认证: **是**
- 需要的角色/权限: `admin`

#### 请求

**请求体:**

```json
{
  "user_id": 0, // 必填，用户ID
  "resource": "string", // 必填，资源名称
  "action": "string", // 必填，操作
  "domain": "string" // 可选，域名
}
```

#### 响应

**成功响应 (200 OK):**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "allowed": true
  },
  "serverTime": 1640000000
}
```

**字段说明:**

- `allowed`: `true` 表示有权限，`false` 表示无权限

#### 示例

**请求示例:**

```bash
curl -X POST http://localhost:9999/api/v1/rbac/check \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 123,
    "resource": "posts",
    "action": "write",
    "domain": "tenant1"
  }'
```

**响应示例:**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "allowed": true
  },
  "serverTime": 1705743600
}
```

---

## 最佳实践

### 角色命名规范

```json
// ✅ 推荐：使用小写和下划线
"super_admin", "content_editor", "guest_viewer"

// ❌ 不推荐：使用大写或空格
"SuperAdmin", "Content Editor"
```

### 资源命名规范

```json
// ✅ 推荐：使用复数形式
"users", "posts", "comments"

// ❌ 不推荐：使用单数
"user", "post"
```

### 操作命名规范

```json
// ✅ 推荐：使用标准动词
"read", "write", "delete", "update"

// ❌ 不推荐：使用自定义动词
"view", "modify", "remove"
```

### 批量操作优化

```bash
# ✅ 推荐：使用批量接口
curl -X POST /api/v1/rbac/users/123/roles/batch \
  -d '{"roles": ["editor", "viewer"]}'

# ❌ 不推荐：多次调用单个接口
curl -X POST /api/v1/rbac/users/123/roles -d '{"role": "editor"}'
curl -X POST /api/v1/rbac/users/123/roles -d '{"role": "viewer"}'
```

## 常见问题

### 1. 如何实现角色继承？

通过为角色分配角色来实现继承：

```bash
# super_admin 继承 admin 的所有权限
curl -X POST /api/v1/rbac/users/super_admin/roles \
  -d '{"role": "admin"}'
```

### 2. 多租户如何隔离权限？

使用 `domain` 参数实现租户间权限隔离：

```bash
# tenant1 的 admin
curl -X POST /api/v1/rbac/users/123/roles \
  -d '{"role": "admin", "domain": "tenant1"}'

# tenant2 的 admin（完全独立）
curl -X POST /api/v1/rbac/users/456/roles \
  -d '{"role": "admin", "domain": "tenant2"}'
```

### 3. 权限检查失败如何排查？

1. 检查用户是否有对应角色
2. 检查角色是否有对应策略
3. 检查域（domain）是否匹配

```bash
# 1. 检查用户角色
GET /api/v1/rbac/users/123/roles?domain=tenant1

# 2. 检查角色策略
GET /api/v1/rbac/roles/editor/policies

# 3. 手动检查权限
POST /api/v1/rbac/check
```

## 注意事项

- ⚠️ **角色分配立即生效**：分配或撤销角色后，权限检查会立即使用新的角色配置
- ⚠️ **策略缓存**：为提升性能，权限检查结果会被缓存，默认TTL为30分钟
- ⚠️ **域名大小写敏感**：domain 参数大小写敏感，请保持一致
- ⚠️ **删除角色不会删除策略**：撤销用户角色不会删除该角色的策略定义

## 相关文档

- [RBAC 包文档](../../pkg/rbac/README.md) - RBAC 工具包使用说明
- [认证说明](../authentication.md) - JWT 认证机制
- [错误码说明](../error-codes.md) - 错误码定义
