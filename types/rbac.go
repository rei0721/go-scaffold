package types

// RBACPolicy RBAC策略
type RBACPolicy struct {
	// Role 角色名称
	Role string `json:"role" binding:"required"`

	// Domain 域名（租户ID），可选
	Domain string `json:"domain,omitempty"`

	// Resource 资源名称
	Resource string `json:"resource" binding:"required"`

	// Action 操作名称
	Action string `json:"action" binding:"required"`
}

// AssignRoleRequest 分配角色请求
type AssignRoleRequest struct {
	// Role 角色名称
	Role string `json:"role" binding:"required"`

	// Domain 域名（租户ID），可选
	Domain string `json:"domain,omitempty"`
}

// AssignRolesRequest 批量分配角色请求
type AssignRolesRequest struct {
	// Roles 角色列表
	Roles []string `json:"roles" binding:"required,min=1"`

	// Domain 域名（租户ID），可选
	Domain string `json:"domain,omitempty"`
}

// CheckPermissionRequest 权限检查请求
type CheckPermissionRequest struct {
	// UserID 用户ID
	UserID int64 `json:"user_id" binding:"required"`

	// Domain 域名（租户ID），可选
	Domain string `json:"domain,omitempty"`

	// Resource 资源名称
	Resource string `json:"resource" binding:"required"`

	// Action 操作名称
	Action string `json:"action" binding:"required"`
}

// CheckPermissionResponse 权限检查响应
type CheckPermissionResponse struct {
	// Allowed 是否允许
	Allowed bool `json:"allowed"`
}

// AddPolicyRequest 添加策略请求
type AddPolicyRequest struct {
	RBACPolicy
}

// AddPoliciesRequest 批量添加策略请求
type AddPoliciesRequest struct {
	// Policies 策略列表
	Policies []RBACPolicy `json:"policies" binding:"required,min=1"`
}

// RemovePolicyRequest 删除策略请求
type RemovePolicyRequest struct {
	RBACPolicy
}

// UserRolesResponse 用户角色响应
type UserRolesResponse struct {
	// UserID 用户ID
	UserID int64 `json:"user_id"`

	// Roles 角色列表
	Roles []string `json:"roles"`
}

// RoleUsersResponse 角色用户响应
type RoleUsersResponse struct {
	// Role 角色名称
	Role string `json:"role"`

	// UserIDs 用户ID列表
	UserIDs []int64 `json:"user_ids"`
}

// PoliciesResponse 策略列表响应
type PoliciesResponse struct {
	// Policies 策略列表
	Policies []RBACPolicy `json:"policies"`

	// Total 总数
	Total int `json:"total"`
}
