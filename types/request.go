// Package types 定义了应用程序的请求和响应类型
// 这个包除了 Go 标准库外没有外部依赖
// 将类型定义从处理器中分离出来,提高可重用性
package types

// RegisterRequest 表示用户注册请求
// 使用 Gin 的 binding tag 进行数据验证
type RegisterRequest struct {
	// Username 用户名
	// binding:"required" - 必填字段
	// min=3 - 最小长度 3 个字符(防止过短的用户名)
	// max=50 - 最大长度 50 个字符(与数据库字段长度一致)
	Username string `json:"username" binding:"required,min=3,max=50"`

	// Email 邮箱地址
	// binding:"required" - 必填字段
	// email - 使用内置的邮箱格式验证器
	// 确保邮箱格式正确,包含 @ 和域名
	Email string `json:"email" binding:"required,email"`

	// Password 密码(明文)
	// binding:"required" - 必填字段
	// min=8 - 最小长度 8 位(安全性要求)
	// 注意:密码在服务端会立即使用 bcrypt 加密,不会存储明文
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest 表示用户登录请求
// 简化的验证规则,只要求字段非空
type LoginRequest struct {
	// Username 用户名
	// 登录时不需要验证长度,只验证是否非空
	Username string `json:"username" binding:"required"`

	// Password 密码
	// 登录时不需要验证长度,只验证是否非空
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest 表示用户更新请求
// 所有字段均为可选，只更新传入的字段
// 使用指针类型区分"未传入"和"传入零值"
type UpdateUserRequest struct {
	// Username 新用户名（可选）
	// 如果传入，需要验证唯一性
	// binding:"omitempty" - 允许不传入此字段
	// min=3,max=50 - 如果传入，需要满足长度要求
	Username *string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`

	// Email 新邮箱（可选）
	// 如果传入，需要验证唯一性和格式
	// email - 如果传入，必须是有效的邮箱格式
	Email *string `json:"email,omitempty" binding:"omitempty,email"`

	// Status 新状态（可选）
	// 1: 激活, 0: 禁用
	// oneof=0 1 - 如果传入，只能是 0 或 1
	Status *int `json:"status,omitempty" binding:"omitempty,oneof=0 1"`
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=50"`
	Description string `json:"description" binding:"max=255"`
	Status      int    `json:"status" binding:"oneof=0 1"` // 默认1
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=2,max=50"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=255"`
	Status      *int    `json:"status,omitempty" binding:"omitempty,oneof=0 1"`
}

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Resource    string `json:"resource" binding:"required,min=2,max=100"`
	Action      string `json:"action" binding:"required,min=2,max=50"`
	Description string `json:"description" binding:"max=255"`
	Status      int    `json:"status" binding:"oneof=0 1"`
}

// UpdatePermissionRequest 更新权限请求
type UpdatePermissionRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Resource    *string `json:"resource,omitempty" binding:"omitempty,min=2,max=100"`
	Action      *string `json:"action,omitempty" binding:"omitempty,min=2,max=50"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=255"`
	Status      *int    `json:"status,omitempty" binding:"omitempty,oneof=0 1"`
}

// AssignRoleRequest 分配角色请求
type AssignRoleRequest struct {
	RoleID int64 `json:"role_id" binding:"required"`
}
