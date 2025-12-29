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
