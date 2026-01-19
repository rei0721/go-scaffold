# pkg/rbac - 基于Casbin的RBAC权限控制

基于Casbin实现的通用RBAC（基于角色的访问控制）工具包，提供简单易用的权限管理接口。

## 特性

- ✅ **简单易用** - 封装Casbin复杂性，提供清晰的API
- ✅ **高性能** - 内置缓存机制，减少数据库查询
- ✅ **持久化** - 使用Gorm Adapter将策略存储到数据库
- ✅ **多租户** - 支持域（Domain），实现租户间权限隔离
- ✅ **灵活** - 支持角色继承、动态策略管理

## 快速开始

### 1. 安装

依赖已包含在项目中：

- `github.com/casbin/casbin/v3`
- `github.com/casbin/gorm-adapter/v3`

### 2. 基本使用

```go
package main

import (
    "fmt"
    "log"

    "github.com/rei0721/go-scaffold/pkg/rbac"
    "gorm.io/gorm"
)

func main() {
    // 假设已有GORM数据库连接
    var db *gorm.DB

    // 1. 创建RBAC实例
    r, err := rbac.New(&rbac.Config{
        DB:          db,
        EnableCache: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer r.Close()

    // 2. 添加策略（定义角色权限）
    r.AddPolicy("admin", "users", "write")
    r.AddPolicy("admin", "posts", "write")
    r.AddPolicy("editor", "posts", "write")
    r.AddPolicy("viewer", "posts", "read")

    // 3. 分配角色
    r.AddRoleForUser("alice", "admin")
    r.AddRoleForUser("bob", "editor")
    r.AddRoleForUser("charlie", "viewer")

    // 4. 检查权限
    ok, _ := r.Enforce("alice", "users", "write")
    fmt.Println("alice can write users:", ok) // true

    ok, _ = r.Enforce("bob", "users", "write")
    fmt.Println("bob can write users:", ok)   // false

    ok, _ = r.Enforce("charlie", "posts", "read")
    fmt.Println("charlie can read posts:", ok) // true
}
```

## API 参考

### 配置

```go
type Config struct {
    DB          *gorm.DB      // 必须：GORM数据库连接
    ModelPath   string        // 可选：自定义模型文件路径
    EnableCache bool          // 可选：是否启用缓存（默认true）
    CacheTTL    time.Duration // 可选：缓存过期时间（默认30分钟）
    AutoSave    bool          // 可选：是否自动保存（默认true）
    TablePrefix string        // 可选：表名前缀
}
```

### 权限检查

```go
// 基本权限检查
ok, err := rbac.Enforce("alice", "data", "read")

// 带域的权限检查（多租户）
ok, err := rbac.EnforceWithDomain("alice", "tenant1", "data", "read")
```

### 角色管理

```go
// 分配角色
rbac.AddRoleForUser("alice", "admin")

// 在指定域中分配角色
rbac.AddRoleForUserInDomain("alice", "admin", "tenant1")

// 撤销角色
rbac.DeleteRoleForUser("alice", "admin")

// 获取用户的角色
roles, err := rbac.GetRolesForUser("alice")

// 获取拥有某角色的所有用户
users, err := rbac.GetUsersForRole("admin")
```

### 策略管理

```go
// 添加策略
rbac.AddPolicy("admin", "users", "write")

// 带域的策略
rbac.AddPolicyWithDomain("admin", "tenant1", "users", "write")

// 删除策略
rbac.RemovePolicy("admin", "users", "write")

// 获取所有策略
policies := rbac.GetPolicy()

// 批量添加策略
rules := [][]string{
    {"admin", "users", "read"},
    {"admin", "users", "write"},
}
rbac.AddPolicies(rules)
```

## 高级用法

### 多租户（域）

```go
// tenant1的admin
rbac.AddRoleForUserInDomain("alice", "admin", "tenant1")
rbac.AddPolicyWithDomain("admin", "tenant1", "data", "read")

// tenant2的admin
rbac.AddRoleForUserInDomain("bob", "admin", "tenant2")
rbac.AddPolicyWithDomain("admin", "tenant2", "data", "read")

// alice只能访问tenant1的数据
ok, _ := rbac.EnforceWithDomain("alice", "tenant1", "data", "read") // true
ok, _ := rbac.EnforceWithDomain("alice", "tenant2", "data", "read") // false
```

### 角色继承

默认的RBAC模型支持角色继承：

```go
// 定义角色层级
rbac.AddRoleForUser("super_admin", "admin")      // super_admin继承admin
rbac.AddRoleForUser("admin", "editor")           // admin继承editor
rbac.AddRoleForUser("editor", "viewer")          // editor继承viewer

// 策略定义
rbac.AddPolicy("viewer", "posts", "read")
rbac.AddPolicy("editor", "posts", "write")
rbac.AddPolicy("admin", "users", "write")

// alice是super_admin，继承了所有下级角色的权限
rbac.AddRoleForUser("alice", "super_admin")
ok, _ := rbac.Enforce("alice", "posts", "read")  // true（继承自viewer）
ok, _ := rbac.Enforce("alice", "users", "write") // true（继承自admin）
```

## 性能优化

### 缓存策略

```go
// 启用缓存（默认）
rbac, _ := rbac.New(&rbac.Config{
    DB:          db,
    EnableCache: true,
    CacheTTL:    30 * time.Minute,
})

// 手动清除缓存（策略变更时自动清除）
rbac.ClearCache()
```

### 批量操作

```go
// ❌ 不推荐：多次单个操作
for _, rule := range rules {
    rbac.AddPolicy(rule[0], rule[1], rule[2])
}

// ✅ 推荐：批量操作
rbac.AddPolicies(rules)
```

## 最佳实践

### 1. 命名规范

```go
// ✅ 角色：小写+下划线
"super_admin", "content_editor", "guest"

// ✅ 资源：复数形式
"users", "posts", "comments"

// ✅ 操作：标准HTTP动词
"read", "write", "delete", "update"
```

### 2. 权限粒度

```go
// ✅ 推荐：细粒度权限
rbac.AddPolicy("editor", "posts", "create")
rbac.AddPolicy("editor", "posts", "update")
rbac.AddPolicy("editor", "posts", "delete")

// ❌ 不推荐：粗粒度权限
rbac.AddPolicy("editor", "posts", "manage")
```

### 3. 使用域隔离租户

```go
// 为每个租户使用独立的域
rbac.AddPolicyWithDomain("admin", "company_a", "data", "write")
rbac.AddPolicyWithDomain("admin", "company_b", "data", "write")

// 检查时必须指定域
rbac.EnforceWithDomain("alice", "company_a", "data", "write")
```

## 数据库表结构

Casbin会自动创建 `casbin_rule` 表：

```sql
CREATE TABLE casbin_rule (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    ptype VARCHAR(100),   -- 类型：p(策略) 或 g(角色)
    v0 VARCHAR(100),      -- 主体/用户
    v1 VARCHAR(100),      -- 域/角色
    v2 VARCHAR(100),      -- 对象
    v3 VARCHAR(100),      -- 操作
    v4 VARCHAR(100),
    v5 VARCHAR(100)
);
```

## 故障排除

### 1. 策略不生效

```go
// 确保策略已保存
rbac.SavePolicy()

// 重新加载策略
rbac.LoadPolicy()
```

### 2. 权限检查总是false

```go
// 检查策略是否存在
policies := rbac.GetPolicy()
fmt.Println(policies)

// 检查用户角色
roles, _ := rbac.GetRolesForUser("alice")
fmt.Println(roles)
```

### 3. 缓存不一致

```go
// 清除缓存
rbac.ClearCache()

// 或禁用缓存
rbac, _ := rbac.New(&rbac.Config{
    DB:          db,
    EnableCache: false,
})
```

## 依赖项

### 必须依赖

- `github.com/casbin/casbin/v3` - Casbin核心库
- `github.com/casbin/gorm-adapter/v3` - GORM适配器
- `gorm.io/gorm` - GORM ORM

### 可选依赖

无

## 参考资料

- [Casbin官方文档](https://casbin.org/)
- [RBAC模型说明](https://casbin.org/docs/rbac)
- [Gorm Adapter](https://github.com/casbin/gorm-adapter)
