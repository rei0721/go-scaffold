/*
Package rbac 提供基于Casbin的RBAC（基于角色的访问控制）功能

# 设计目标

- 简单易用：提供清晰的API，隐藏Casbin的复杂性
- 高性能：内置缓存机制，减少数据库查询
- 持久化：使用Gorm Adapter将策略持久化到数据库
- 灵活性：支持域（多租户）、角色继承等高级特性

# 核心概念

RBAC模型包含以下核心概念：

1. 主体（Subject）：用户或角色
2. 对象（Object）：被访问的资源
3. 操作（Action）：对资源的操作（read/write/delete等）
4. 域（Domain）：用于多租户场景，实现权限隔离

# 使用示例

基本用法：

	// 1. 创建RBAC实例
	rbac, err := rbac.New(&rbac.Config{
	    DB: db,  // GORM数据库连接
	    EnableCache: true,
	    CacheTTL: 30 * time.Minute,
	})
	if err != nil {
	    log.Fatal(err)
	}
	defer rbac.Close()

	// 2. 添加策略：admin角色可以对users资源进行write操作
	rbac.AddPolicy("admin", "users", "write")

	// 3. 为用户分配角色：alice是admin
	rbac.AddRoleForUser("alice", "admin")

	// 4. 检查权限：alice能否对users进行write操作？
	ok, err := rbac.Enforce("alice", "users", "write")
	if ok {
	    fmt.Println("允许访问")
	} else {
	    fmt.Println("拒绝访问")
	}

带域（多租户）的用法：

	// 在tenant1域中，alice是admin
	rbac.AddRoleForUserInDomain("alice", "admin", "tenant1")

	// 在tenant1域中，admin可以访问data
	rbac.AddPolicyWithDomain("admin", "tenant1", "data", "read")

	// 检查：alice在tenant1域中能否读取data？
	ok, err := rbac.EnforceWithDomain("alice", "tenant1", "data", "read")

批量操作：

	// 批量添加策略
	rules := [][]string{
	    {"admin", "users", "read"},
	    {"admin", "users", "write"},
	    {"admin", "posts", "read"},
	}
	rbac.AddPolicies(rules)

# 最佳实践

1. 角色命名：使用小写和下划线，如 "super_admin", "content_editor"
2. 资源命名：使用复数形式，如 "users", "posts"
3. 操作命名：使用标准HTTP动词，如 "read", "write", "delete"
4. 缓存管理：策略变更后会自动清除缓存，无需手动处理

# 性能优化

1. 启用缓存：默认启用，可显著提升权限检查性能
2. 批量操作：使用 AddPolicies 而不是多次调用 AddPolicy
3. 合理的缓存TTL：根据实际需求调整，默认30分钟

# 数据库表结构

Casbin会自动创建 casbin_rule 表，结构如下：

	CREATE TABLE casbin_rule (
	    id BIGINT PRIMARY KEY,
	    ptype VARCHAR(100),  -- p(策略) 或 g(角色)
	    v0 VARCHAR(100),     -- subject/user
	    v1 VARCHAR(100),     -- domain/role
	    v2 VARCHAR(100),     -- object
	    v3 VARCHAR(100),     -- action
	    v4 VARCHAR(100),
	    v5 VARCHAR(100)
	);

# 与其他包的区别

- pkg/jwt：处理身份认证（用户是谁）
- pkg/rbac：处理授权（用户能做什么）
- pkg/cache：通用缓存，rbac内部使用sync.Map做权限结果缓存

# 线程安全

所有方法都是线程安全的，可以在并发环境下使用。
*/
package rbac
