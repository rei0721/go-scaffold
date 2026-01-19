# 错误码说明

## 错误响应格式

当API发生错误时,会返回统一的错误响应格式:

```json
{
  "code": 1001, // 错误码
  "message": "Invalid parameter", // 错误消息
  "traceId": "abc123", // 追踪ID(可用于问题排查)
  "serverTime": 1640000000 // 服务器时间戳
}
```

## 错误码分类

错误码采用4位数字,按照以下规则分类:

- **0** - 成功
- **1xxx** - 客户端错误 (对应 4xx HTTP 状态码)
- **5xxx** - 服务器错误 (对应 5xx HTTP 状态码)

## 通用错误码

| code | HTTP状态码 | message               | 说明                 |
| ---- | ---------- | --------------------- | -------------------- |
| 0    | 200        | success               | 请求成功             |
| 1000 | 400        | Bad request           | 请求格式错误         |
| 1001 | 400        | Invalid parameter     | 参数错误             |
| 1002 | 401        | Unauthorized          | 未认证               |
| 1003 | 403        | Forbidden             | 权限不足             |
| 1004 | 404        | Not found             | 资源不存在           |
| 1005 | 409        | Conflict              | 资源冲突(如重复创建) |
| 1006 | 429        | Too many requests     | 请求过于频繁         |
| 5000 | 500        | Internal server error | 服务器内部错误       |
| 5001 | 503        | Service unavailable   | 服务暂时不可用       |

## 业务错误码

### 认证模块 (10xx)

| code | message             | 说明             |
| ---- | ------------------- | ---------------- |
| 1010 | Invalid credentials | 用户名或密码错误 |
| 1011 | User already exists | 用户已存在       |
| 1012 | User not found      | 用户不存在       |
| 1013 | Invalid token       | Token 无效       |
| 1014 | Token expired       | Token 已过期     |

### RBAC模块 (11xx)

| code | message               | 说明       |
| ---- | --------------------- | ---------- |
| 1110 | Role not found        | 角色不存在 |
| 1111 | Permission denied     | 权限不足   |
| 1112 | Invalid policy        | 策略无效   |
| 1113 | Role already assigned | 角色已分配 |

## 使用建议

### 1. 客户端错误处理

```javascript
if (response.code !== 0) {
  switch (response.code) {
    case 1002: // Unauthorized
      // 跳转到登录页
      redirectToLogin();
      break;
    case 1003: // Forbidden
      // 显示权限不足提示
      showPermissionDenied();
      break;
    default:
      // 显示错误消息
      showError(response.message);
  }
}
```

### 2. 使用 traceId 排查问题

当遇到错误时,记录 `traceId` 并提交给技术支持:

```
错误: Internal server error
TraceId: abc123def456
```

技术人员可以通过 traceId 在日志中快速定位问题。

## 扩展错误码

如需添加新的业务错误码,请:

1. 在 `types/errors` 包中定义常量
2. 更新此文档
3. 确保错误码在整个项目中唯一
