# API 文档

这里包含项目所有 HTTP API 的详细文档。

## 文档组织

- [认证说明](./authentication.md) - API 认证和授权机制
- [错误码说明](./error-codes.md) - 统一的错误码定义
- [通用规范](./common/) - 请求/响应格式、分页等通用规范
- [API 接口](./endpoints/) - 各功能模块的具体 API 接口

## 快速开始

### 基础URL

```
开发环境: http://localhost:9999
生产环境: https://api.example.com
```

### 认证方式

所有需要认证的接口都使用 JWT Bearer Token:

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:9999/api/v1/protected-endpoint
```

### 通用响应格式

所有API响应遵循统一格式:

```json
{
  "code": 0, // 0表示成功,其他值表示错误
  "message": "success", // 响应消息
  "data": {}, // 实际数据(可选)
  "traceId": "...", // 追踪ID(可选)
  "serverTime": 1640000000 // 服务器时间戳
}
```

## API 列表

### 认证相关

| 方法 | 路径                         | 说明     | 文档链接                                                  |
| ---- | ---------------------------- | -------- | --------------------------------------------------------- |
| POST | /api/v1/auth/register        | 用户注册 | [详情](./endpoints/auth.md#post-apiv1authregister)        |
| POST | /api/v1/auth/login           | 用户登录 | [详情](./endpoints/auth.md#post-apiv1authlogin)           |
| POST | /api/v1/auth/logout          | 用户登出 | [详情](./endpoints/auth.md#post-apiv1authlogout)          |
| POST | /api/v1/auth/change-password | 修改密码 | [详情](./endpoints/auth.md#post-apiv1authchange-password) |
| POST | /api/v1/auth/refresh         | 刷新令牌 | [详情](./endpoints/auth.md#post-apiv1authrefresh)         |

### RBAC 权限管理

| 方法   | 路径                               | 说明         | 文档链接                                                     |
| ------ | ---------------------------------- | ------------ | ------------------------------------------------------------ |
| POST   | /api/v1/rbac/users/:id/roles       | 分配角色     | [详情](./endpoints/rbac.md#post-apiv1rbacusersid roles)      |
| POST   | /api/v1/rbac/users/:id/roles/batch | 批量分配角色 | [详情](./endpoints/rbac.md#post-apiv1rbacusersidrolesbatch)  |
| DELETE | /api/v1/rbac/users/:id/roles/:role | 撤销角色     | [详情](./endpoints/rbac.md#delete-apiv1rbacusersidrolesrole) |
| GET    | /api/v1/rbac/users/:id/roles       | 获取用户角色 | [详情](./endpoints/rbac.md#get-apiv1rbacusersidroles)        |
| GET    | /api/v1/rbac/roles/:role/users     | 获取角色用户 | [详情](./endpoints/rbac.md#get-apiv1rbacroleseroleusers)     |
| POST   | /api/v1/rbac/policies              | 添加策略     | [详情](./endpoints/rbac.md#post-apiv1rbacpolicies)           |
| POST   | /api/v1/rbac/policies/batch        | 批量添加策略 | [详情](./endpoints/rbac.md#post-apiv1rbacpoliciesbatch)      |
| DELETE | /api/v1/rbac/policies              | 删除策略     | [详情](./endpoints/rbac.md#delete-apiv1rbacpolicies)         |
| GET    | /api/v1/rbac/policies              | 获取所有策略 | [详情](./endpoints/rbac.md#get-apiv1rbacpolicies)            |
| GET    | /api/v1/rbac/roles/:role/policies  | 获取角色策略 | [详情](./endpoints/rbac.md#get-apiv1rbacroles rolepolicies)  |
| POST   | /api/v1/rbac/check                 | 检查权限     | [详情](./endpoints/rbac.md#post-apiv1rbaccheck)              |

## 版本历史

- **v1** (当前版本)
  - 认证模块
  - RBAC 权限管理模块

## 文档维护

本文档应与代码保持同步。每次新增或修改 API 时,请参考 [api-docs-sync skill](../.agent/skills/api-docs-sync/SKILL.md) 更新文档。
