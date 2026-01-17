package rbac

// RBAC 权限控制接口
//
// 基于Casbin实现的RBAC（基于角色的访问控制）
// 支持角色继承、域（多租户）等高级特性
type RBAC interface {
	// ========== 权限检查 ==========

	// Enforce 检查权限
	// 参数:
	//   sub: 主体（用户ID或角色）
	//   obj: 对象（资源）
	//   act: 操作（read/write/delete等）
	// 返回:
	//   bool: 是否有权限
	//   error: 检查过程中的错误
	// 示例:
	//   ok, err := rbac.Enforce("alice", "data1", "read")
	Enforce(sub, obj, act string) (bool, error)

	// EnforceWithDomain 带域的权限检查
	// 用于多租户场景，不同域的权限相互隔离
	// 参数:
	//   sub: 主体（用户ID或角色）
	//   dom: 域（租户ID）
	//   obj: 对象（资源）
	//   act: 操作
	// 示例:
	//   ok, err := rbac.EnforceWithDomain("alice", "tenant1", "data1", "read")
	EnforceWithDomain(sub, dom, obj, act string) (bool, error)

	// ========== 角色管理 ==========

	// AddRoleForUser 为用户分配角色
	// 参数:
	//   user: 用户ID
	//   role: 角色名称
	// 示例:
	//   rbac.AddRoleForUser("alice", "admin")
	AddRoleForUser(user, role string) error

	// AddRoleForUserInDomain 在指定域中为用户分配角色
	// 参数:
	//   user: 用户ID
	//   role: 角色名称
	//   domain: 域名
	AddRoleForUserInDomain(user, role, domain string) error

	// DeleteRoleForUser 撤销用户的角色
	// 参数:
	//   user: 用户ID
	//   role: 角色名称
	DeleteRoleForUser(user, role string) error

	// DeleteRoleForUserInDomain 在指定域中撤销用户的角色
	DeleteRoleForUserInDomain(user, role, domain string) error

	// GetRolesForUser 获取用户的所有角色
	// 参数:
	//   user: 用户ID
	// 返回:
	//   []string: 角色列表
	GetRolesForUser(user string) ([]string, error)

	// GetRolesForUserInDomain 获取用户在指定域中的角色
	GetRolesForUserInDomain(user, domain string) ([]string, error)

	// GetUsersForRole 获取拥有指定角色的所有用户
	// 参数:
	//   role: 角色名称
	// 返回:
	//   []string: 用户ID列表
	GetUsersForRole(role string) ([]string, error)

	// ========== 策略管理 ==========

	// AddPolicy 添加策略
	// 参数:
	//   sub: 主体（通常是角色）
	//   obj: 对象
	//   act: 操作
	// 示例:
	//   rbac.AddPolicy("admin", "users", "write")
	AddPolicy(sub, obj, act string) error

	// AddPolicyWithDomain 添加带域的策略
	AddPolicyWithDomain(sub, domain, obj, act string) error

	// RemovePolicy 删除策略
	RemovePolicy(sub, obj, act string) error

	// RemovePolicyWithDomain 删除带域的策略
	RemovePolicyWithDomain(sub, domain, obj, act string) error

	// GetPolicy 获取所有策略
	// 返回:
	//   [][]string: 策略列表，每个策略是[sub, obj, act]
	GetPolicy() [][]string

	// GetFilteredPolicy 获取过滤后的策略
	// 参数:
	//   fieldIndex: 字段索引（0=sub, 1=obj, 2=act）
	//   fieldValues: 过滤值
	GetFilteredPolicy(fieldIndex int, fieldValues ...string) [][]string

	// ========== 批量操作 ==========

	// AddPolicies 批量添加策略
	// 参数:
	//   rules: 策略列表，每个策略是[sub, obj, act]或[sub, dom, obj, act]
	// 示例:
	//   rules := [][]string{
	//       {"admin", "users", "read"},
	//       {"admin", "users", "write"},
	//   }
	//   rbac.AddPolicies(rules)
	AddPolicies(rules [][]string) error

	// RemovePolicies 批量删除策略
	RemovePolicies(rules [][]string) error

	// ========== 工具方法 ==========

	// LoadPolicy 从存储加载策略
	// 通常在Enforcer初始化时自动调用
	LoadPolicy() error

	// SavePolicy 保存策略到存储
	// 如果AutoSave=true，策略变更会自动保存，无需手动调用
	SavePolicy() error

	// ClearCache 清除缓存
	// 当策略变更时，应该清除缓存以确保一致性
	ClearCache() error

	// Close 关闭RBAC实例
	// 释放资源
	Close() error
}
