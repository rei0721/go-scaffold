# RBAC 示例

这个目录包含 `pkg/rbac` 包的使用示例。

## 运行示例

### 基础示例 (basic)

演示如何设置完整的 RBAC 系统，包括创建角色、权限、分配权限和权限检查。

```bash
cd examples/basic
go run main.go
```

**示例输出**：

```
=== 创建角色 ===
创建角色: admin (ID: 1)
创建角色: editor (ID: 2)

=== 创建权限 ===
创建权限: users:read (读取用户)
创建权限: users:write (创建/更新用户)
创建权限: posts:read (读取文章)
创建权限: posts:write (创建/更新文章)
创建权限: posts:delete (删除文章)

=== 分配权限给角色 ===
为 admin 角色分配所有权限
为 editor 角色分配文章读写权限

=== 分配角色给用户 ===
用户 1 被分配为 admin
用户 2 被分配为 editor

=== 权限检查 ===
admin - users:write ✅ 有权限
admin - posts:delete ✅ 有权限
editor - posts:write ✅ 有权限
editor - posts:delete ❌ 无权限
editor - users:write ❌ 无权限

✅ 示例完成！
```

## 示例说明

### basic 示例

**演示内容**：

1. 初始化数据库（SQLite 内存数据库）
2. 创建 Repository 和 Service
3. 创建角色（admin、editor）
4. 创建权限（users:read、users:write、posts:read、posts:write、posts:delete）
5. 为角色分配权限
6. 为用户分配角色
7. 检查用户权限

**适用场景**：

- 快速了解 RBAC 的基本用法
- 学习如何设置角色和权限
- 理解权限检查的工作原理

## 更多示例

未来可能添加的示例：

- **with_cache** - 集成缓存的示例
- **with_middleware** - HTTP 中间件集成示例
- **dynamic_permissions** - 动态权限管理示例
- **hierarchical_roles** - 分层角色示例

## 注意事项

1. 示例使用内存数据库，重启后数据会丢失
2. 生产环境请使用 MySQL、PostgreSQL 等持久化数据库
3. 建议启用缓存以提升性能
