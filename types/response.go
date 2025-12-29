package types

import "time"

// UserResponse 表示 API 响应中的用户信息
// 这是一个 DTO(Data Transfer Object),用于将数据库模型转换为 API 响应
// 注意:不包含 Password 字段,确保密码永远不会被返回
type UserResponse struct {
	// UserID 用户 ID
	// 使用 userId 而不是 id,遵循前端命名规范
	UserID int64 `json:"userId"`

	// Username 用户名
	Username string `json:"username"`

	// Email 邮箱地址
	Email string `json:"email"`

	// Status 用户状态
	// 1: 激活, 0: 未激活/禁用
	Status int `json:"status"`

	// CreatedAt 创建时间
	// 使用 RFC3339 格式进行 JSON 序列化
	CreatedAt time.Time `json:"createdAt"`
}

// LoginResponse 表示成功登录的响应
// 包含访问令牌和用户信息
type LoginResponse struct {
	// Token 访问令牌(JWT)
	// 前端需要将它存储在 localStorage 或 cookie 中
	// 并在后续请求中放在 Authorization header 中
	Token string `json:"token"`

	// ExpiresIn 令牌有效期(秒)
	// 前端可以用来计算令牌过期时间
	ExpiresIn int `json:"expiresIn"`

	// User 用户信息
	// 嵌入用户数据,避免前端需要再次请求获取用户信息
	User UserResponse `json:"user"`
}
