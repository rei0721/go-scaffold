# pkg/rbac - RBAC 权限控制工具库

基于角色的访问控制（Role-Based Access Control）通用库，提供完整的角色、权限管理功能。

## 特性

- ✅ 完整的 RBAC 模型（角色、权限、用户-角色、角色-权限）
- ✅ GORM 数据库抽象，支持 MySQL、PostgreSQL、SQLite
- ✅ 缓存集成，减少数据库查询
- ✅ 权限通配符支持（`*:*`, `resource:*`）
- ✅ 异步缓存操作，不阻塞主请求
- ✅ 延迟依赖注入（Cache、Executor）
- ✅ 线程安全

## 快速开始

### 1. 安装

```bash
go get github.com/rei0721/rei0721/pkg/rbac
```

### 2. 基本使用

```go
import (
    "github.com/rei0721/rei0721/pkg/rbac/repository"
    "github.com/rei0721/rei0721/pkg/rbac/service"
)

// 创建 Repository
repo := repository.NewGormRBACRepository(db)

// 创建 Service
rbacService := service.NewRBACService(repo)

// 可选：注入缓存和协程池提升性能
rbacService.SetCache(cache)
rbacService.SetExecutor(executor)
```

### 3. 创建角色和权限

```go
// 创建角色
role, err := rbacService.CreateRole(ctx, &types.CreateRoleRequest{
    Name:        "editor",
    Description: "内容编辑者",
    Status:      1,
})

// 创建权限
perm, err := rbacService.CreatePermission(ctx, &types.CreatePermissionRequest{
    Name:        "posts:write",
    Resource:    "posts",
    Action:      "write",
    Description: "文章写权限",
    Status:      1,
})

// 为角色分配权限
err = rbacService.AssignPermission(ctx, role.ID, perm.ID)
```

### 4. 权限检查

```go
// 检查用户权限
hasPermission, err := rbacService.CheckPermission(ctx, userID, "posts", "write")
if err != nil {
    // 处理错误
}
if !hasPermission {
    // 无权限
}
```

### 5. 在 HTTP 中间件中使用

```go
import "github.com/rei0721/rei0721/internal/middleware"

// 权限中间件
router.Use(middleware.RequirePermission(rbacService, "posts", "write"))

// 角色中间件
router.Use(middleware.RequireRole(rbacService, "admin"))
```

## 数据模型

### Role - 角色

```go
type Role struct {
    ID          int64
    Name        string       // 角色名称，如 "admin", "editor"
    Description string       // 描述
    Status      int          // 1: 启用, 0: 禁用
    Permissions []Permission // 角色拥有的权限
}
```

### Permission - 权限

```go
type Permission struct {
    ID          int64
    Name        string // 权限名称，如 "posts:write"
    Resource    string // 资源标识，如 "posts", "users"
    Action      string // 操作类型，如 "read", "write", "delete"
    Description string
    Status      int    // 1: 启用, 0: 禁用
}
```

## 权限格式

权限使用 `resource:action` 格式：

| 格式              | 说明           | 示例                         |
| ----------------- | -------------- | ---------------------------- |
| `resource:action` | 基本格式       | `users:read`, `posts:write`  |
| `*:*`             | 超级管理员权限 | 匹配所有资源和操作           |
| `resource:*`      | 资源通配符     | `posts:*` 匹配文章的所有操作 |

## 缓存策略

- **缓存键**: `user:perms:{userID}`
- **缓存内容**: 用户的权限集合 `["users:read", "posts:write", ...]`
- **TTL**: 60 分钟
- **失效时机**: 用户角色变更时自动清除
- **写入方式**: 异步，不阻塞主请求

## API 参考

完整 API 文档请查看：

- [Repository 接口](repository/interface.go) - 数据访问层
- [Service 接口](service/interface.go) - 业务逻辑层
- [包文档](doc.go) - 详细使用说明

## 最佳实践

### 1. 角色设计

遵循最小权限原则，按职责划分角色：

```
superadmin  -> *:*
admin       -> users:*, roles:*, posts:*
editor      -> posts:write, posts:read
viewer      -> posts:read
```

### 2. 权限粒度

根据业务需求平衡粗细粒度：

- ✅ 推荐：`posts:write`（创建和更新文章）
- ❌ 过细：`posts:create`, `posts:update`
- ❌ 过粗：`admin:*`（权限范围太大）

### 3. 性能优化

生产环境务必注入 Cache 和 Executor：

```go
rbacService.SetCache(cache)      // 启用权限缓存
rbacService.SetExecutor(executor) // 异步缓存操作
```

### 4. 数据库索引

确保以下表有正确的索引：

```sql
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);
```

## 依赖项

### 必须依赖

- `gorm.io/gorm` - ORM 框架

### 可选依赖（通过接口注入）

- `github.com/rei0721/rei0721/pkg/cache` - 缓存抽象
- `github.com/rei0721/rei0721/pkg/executor` - 协程池管理

## 与 pkg/jwt 的配合

RBAC 通常与 JWT 认证配合使用：

1. **JWT** - 身份认证（Authentication），验证"你是谁"
2. **RBAC** - 授权（Authorization），验证"你能做什么"

典型流程：

```
用户登录 -> JWT 验证身份 -> RBAC 检查权限 -> 允许/拒绝访问
```

## 迁移指南

如果你的项目中已有 RBAC 实现，迁移步骤：

1. 导入 `pkg/rbac`
2. 将现有 Repository 替换为 `repository.NewGormRBACRepository(db)`
3. 将现有 Service 替换为 `service.NewRBACService(repo)`
4. 注入可选依赖（Cache、Executor）
5. 更新中间件引用

## 许可证

本项目遵循项目根目录的许可证。

## 相关链接

- [项目文档](../../docs/)
- [系统架构](../../docs/architecture/system_map.md)
