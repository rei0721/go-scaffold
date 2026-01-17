# RBAC 使用示例

本目录包含pkg/rbac的使用示例。

## 运行示例

### 基础示例

```bash
cd examples/basic
go run main.go
```

**示例输出**：

```
=== RBAC基本示例 ===

1. 添加策略：
  - admin 可以 write users
  - admin 可以 write posts
  - editor 可以 write posts
  - viewer 可以 read posts

2. 分配角色：
  - alice -> admin
  - bob -> editor
  - charlie -> viewer

3. 权限检查：
  alice write users: ✅ 允许
  alice write posts: ✅ 允许
  bob write users: ❌ 拒绝
  bob write posts: ✅ 允许
  charlie read posts: ✅ 允许
  charlie write posts: ❌ 拒绝

=== 角色继承示例 ===

alice被分配为super_admin（继承admin）
alice write users: ✅ 允许

=== 多租户（域）示例 ===

alice在tenant1中是admin
bob在tenant2中是admin
alice能访问tenant1的data吗？true
alice能访问tenant2的data吗？false

=== 批量操作示例 ===

批量添加了3条策略给manager角色

=== 查询示例 ===

alice的角色：[super_admin admin]
拥有admin角色的用户：[alice super_admin bob]
所有策略数量：7
```

## 示例说明

### basic 示例

演示了以下功能：

- 创建RBAC实例
- 添加策略（角色权限）
- 分配角色给用户
- 权限检查
- 角色继承
- 多租户（域）
- 批量操作
- 查询功能

### 关键代码片段

1. **创建实例**：

```go
r, err := rbac.New(&rbac.Config{
    DB:          db,
    EnableCache: true,
})
```

2. **添加策略**：

```go
r.AddPolicy("admin", "users", "write")
```

3. **分配角色**：

```go
r.AddRoleForUser("alice", "admin")
```

4. **检查权限**：

```go
ok, err := r.Enforce("alice", "users", "write")
```

5. **多租户**：

```go
r.AddRoleForUserInDomain("alice", "admin", "tenant1")
ok, _ := r.EnforceWithDomain("alice", "tenant1", "data", "read")
```
