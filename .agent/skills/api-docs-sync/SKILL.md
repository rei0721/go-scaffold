---
name: api-docs-sync
description: HTTP API实现后必须同步更新项目API文档
---

# API 文档同步规范

## 概述

本 skill 确保每次新增或修改 HTTP API 接口后,必须同步更新项目的 API 文档(`docs/api/**`)。这保证了文档与代码的一致性,提升团队协作效率。

## 适用场景

- ✅ 新增 HTTP API 接口
- ✅ 修改现有 API 的请求/响应格式
- ✅ 修改 API 路由路径
- ✅ 修改 API 认证/权限要求
- ✅ 废弃或删除 API 接口

## API 文档目录结构

```
docs/api/
├── README.md                    # API 文档索引
├── authentication.md            # 认证说明
├── error-codes.md              # 错误码说明
├── common/
│   ├── request-format.md       # 通用请求格式
│   ├── response-format.md      # 通用响应格式
│   └── pagination.md           # 分页规范
└── endpoints/
    ├── auth.md                 # 认证相关接口
    ├── user.md                 # 用户管理接口
    ├── rbac.md                 # RBAC权限管理接口
    └── ...
```

## API 文档模板

### 单个API接口文档格式

```markdown
## {方法} {路径}

### 描述

简短描述这个接口的用途

### 认证

- 是否需要认证: **是/否**
- 需要的角色/权限: `admin` / `user` / 无

### 请求

**路径参数:**

| 参数名 | 类型  | 必填 | 说明   |
| ------ | ----- | ---- | ------ |
| id     | int64 | 是   | 用户ID |

**查询参数:**

| 参数名   | 类型 | 必填 | 默认值 | 说明     |
| -------- | ---- | ---- | ------ | -------- |
| page     | int  | 否   | 1      | 页码     |
| pageSize | int  | 否   | 10     | 每页数量 |

**请求体:**

\`\`\`json
{
"username": "string", // 必填,用户名
"email": "string", // 必填,邮箱
"password": "string" // 必填,密码(最少8位)
}
\`\`\`

### 响应

**成功响应 (200 OK):**

\`\`\`json
{
"code": 0,
"message": "success",
"data": {
"id": 123,
"username": "john_doe",
"email": "john@example.com",
"createdAt": "2024-01-01T00:00:00Z"
},
"serverTime": 1640000000
}
\`\`\`

**错误响应:**

| 状态码 | code | message           | 说明           |
| ------ | ---- | ----------------- | -------------- |
| 400    | 1001 | Invalid parameter | 参数错误       |
| 401    | 1002 | Unauthorized      | 未认证         |
| 403    | 1003 | Forbidden         | 权限不足       |
| 500    | 5000 | Internal error    | 服务器内部错误 |

### 示例

**请求示例:**

\`\`\`bash
curl -X POST http://localhost:9999/api/v1/users \\
-H "Authorization: Bearer YOUR_TOKEN" \\
-H "Content-Type: application/json" \\
-d '{
"username": "john_doe",
"email": "john@example.com",
"password": "password123"
}'
\`\`\`

**响应示例:**

\`\`\`json
{
"code": 0,
"message": "success",
"data": {
"id": 123,
"username": "john_doe",
"email": "john@example.com",
"createdAt": "2024-01-01T00:00:00Z"
},
"serverTime": 1640000000
}
\`\`\`

### 注意事项

- 特殊说明和注意事项
- 业务规则限制
- 性能相关提示
```

## 同步更新流程

### 步骤 1: 创建/更新 Handler 后立即更新文档

在完成 Handler 开发后(**在提交代码前**),立即更新对应的 API 文档:

```bash
# 1. 确定接口所属的功能模块
# 例如: RBAC 相关接口 -> docs/api/endpoints/rbac.md

# 2. 如果是新模块,创建新文档
touch docs/api/endpoints/{module}.md

# 3. 编辑文档,添加或更新接口说明
code docs/api/endpoints/{module}.md
```

### 步骤 2: 使用模板填写接口信息

按照上述模板,依次填写:

1. **方法和路径** - 从 router.go 中确认
2. **描述** - 从 Handler 注释中提取
3. **认证要求** - 检查是否使用 AuthMiddleware 和 RBAC 中间件
4. **请求参数** - 从 Handler 中的 request 结构体提取
5. **响应格式** - 从 Handler 中的 response 结构体提取
6. **错误码** - 列出可能的错误响应
7. **示例** - 提供可运行的 curl 示例

### 步骤 3: 更新 API 索引

在 `docs/api/README.md` 中更新 API 列表:

```markdown
# API 文档索引

## RBAC 权限管理

| 方法   | 路径                               | 说明         | 文档链接                                 |
| ------ | ---------------------------------- | ------------ | ---------------------------------------- |
| POST   | /api/v1/rbac/users/:id/roles       | 分配角色     | [详情](./endpoints/rbac.md#分配角色)     |
| DELETE | /api/v1/rbac/users/:id/roles/:role | 撤销角色     | [详情](./endpoints/rbac.md#撤销角色)     |
| GET    | /api/v1/rbac/users/:id/roles       | 获取用户角色 | [详情](./endpoints/rbac.md#获取用户角色) |
```

### 步骤 4: 验证文档完整性

使用以下检查清单验证文档:

- [ ] 接口路径与 router.go 中的注册一致
- [ ] 请求参数与 types 包中的结构体一致
- [ ] 响应格式与 Handler 返回的结构一致
- [ ] 认证要求准确(检查中间件)
- [ ] 错误码与 types/errors 包一致
- [ ] curl 示例可以实际运行
- [ ] 在 API 索引中已添加链接

## 实战示例：同步 RBAC API 文档

### 场景

刚刚实现了 RBAC 相关的 HTTP 接口,需要同步更新文档。

### 执行步骤

**1. 创建 RBAC API 文档:**

```bash
touch docs/api/endpoints/rbac.md
```

**2. 填写文档内容:**

从 `internal/handler/rbac_handler.go` 和 `internal/router/router.go` 提取信息:

```markdown
# RBAC 权限管理 API

## POST /api/v1/rbac/users/:id/roles

### 描述

为指定用户分配角色

### 认证

- 是否需要认证: **是**
- 需要的角色/权限: `admin`

### 请求

**路径参数:**

| 参数名 | 类型  | 必填 | 说明   |
| ------ | ----- | ---- | ------ |
| id     | int64 | 是   | 用户ID |

**请求体:**

\`\`\`json
{
"role": "string", // 必填,角色名称
"domain": "string" // 可选,域名(多租户场景)
}
\`\`\`

### 响应

**成功响应 (200 OK):**

\`\`\`json
{
"code": 0,
"message": "success",
"data": {
"message": "Role assigned successfully"
},
"serverTime": 1640000000
}
\`\`\`

...
```

**3. 更新 API 索引:**

编辑 `docs/api/README.md`,添加 RBAC 相关接口列表。

**4. 提交代码:**

```bash
git add internal/handler/rbac_handler.go
git add docs/api/endpoints/rbac.md
git add docs/api/README.md
git commit -m "feat: add RBAC management API and documentation"
```

## 文档维护最佳实践

### 1. 文档与代码同步提交

```bash
# ✅ 好的做法：同时提交代码和文档
git add internal/handler/user_handler.go
git add docs/api/endpoints/user.md
git commit -m "feat: add user update API"

# ❌ 不好的做法：只提交代码,稍后再补文档
git add internal/handler/user_handler.go
git commit -m "feat: add user update API"
# (文档被遗忘...)
```

### 2. 使用清晰的 commit message

```bash
# API 新增
git commit -m "feat(api): add RBAC role assignment endpoint"

# API 修改
git commit -m "fix(api): update user creation response format"

# API 废弃
git commit -m "docs(api): mark /api/v1/old-endpoint as deprecated"
```

### 3. 定期审查文档一致性

建议每个迭代结束时,审查文档与代码的一致性:

```bash
# 检查所有 API Handler
find internal/handler -name "*_handler.go"

# 检查对应的 API 文档
find docs/api/endpoints -name "*.md"
```

### 4. 使用工具辅助

推荐使用以下工具:

- **Swagger/OpenAPI** - 从注释生成文档
- **Postman Collection** - 导出 API 测试集合
- **API Blueprint** - 使用结构化格式定义 API

## 文档规范补充

### 命名规范

**文件命名:**

- 使用小写字母和连字符
- 模块名.md,如 `user.md`, `rbac.md`

**标题命名:**

- 使用 HTTP 方法 + 路径,如 `## POST /api/v1/users`

### 版本管理

当 API 有多个版本时:

```
docs/api/
├── v1/
│   └── endpoints/
│       └── user.md
└── v2/
    └── endpoints/
        └── user.md
```

### 废弃标记

废弃的 API 需要明确标记:

```markdown
## ~~POST /api/v1/old-endpoint~~ (已废弃)

> ⚠️ **已废弃:** 请使用 [POST /api/v2/new-endpoint](#post-apiv2new-endpoint) 替代
>
> **废弃时间:** 2024-01-01
>
> **移除计划:** 2024-06-01
```

## 检查清单

### 代码实现完成后

- [ ] 已在 `docs/api/endpoints/` 中创建或更新对应文档
- [ ] 接口路径与 router.go 一致
- [ ] 请求参数文档与 types 包一致
- [ ] 响应格式文档与 Handler 返回一致
- [ ] 认证要求准确反映中间件配置
- [ ] 错误码与 types/errors 包一致
- [ ] 提供了可运行的 curl 示例
- [ ] 在 `docs/api/README.md` 更新了索引
- [ ] 文档与代码在同一个 commit 中提交

### Code Review 时

- [ ] Reviewer 检查文档是否与代码一致
- [ ] 检查 curl 示例是否可以实际运行
- [ ] 检查错误响应是否完整
- [ ] 检查认证和权限说明是否准确

## 相关 Skills

- [handler-development](../handler-development/SKILL.md) - Handler 开发规范
- [docs-sync](../docs-sync/SKILL.md) - 通用文档同步规范
- [error-handling](../error-handling/SKILL.md) - 错误码定义规范
