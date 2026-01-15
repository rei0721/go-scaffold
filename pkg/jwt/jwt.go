package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

// JWT 定义JWT操作接口
// 提供token的生成、验证和刷新功能
// 为什么使用接口:
// - 定义契约: 明确JWT组件提供的能力
// - 依赖倒置: 使用方依赖接口而非实现
// - 便于测试: 可以创建mock实现进行单元测试
// - 解耦: 可以轻松替换不同的JWT实现
type JWT interface {
	// GenerateToken 生成访问令牌
	// 参数:
	//   userID: 用户ID
	//   username: 用户名
	// 返回:
	//   string: JWT token字符串
	//   error: 生成失败时的错误
	// 业务流程:
	//   1. 创建claims载荷
	//   2. 使用HMAC-SHA256算法签名
	//   3. 生成完整的JWT token
	GenerateToken(userID int64, username string) (string, error)

	// ValidateToken 验证并解析令牌
	// 参数:
	//   tokenString: JWT token字符串
	// 返回:
	//   *Claims: 解析后的载荷信息
	//   error: 验证失败时的错误,如:
	//     - ErrInvalidToken: token格式无效
	//     - ErrExpiredToken: token已过期
	//     - ErrInvalidSignature: 签名验证失败
	// 业务流程:
	//   1. 解析token字符串
	//   2. 验证签名
	//   3. 检查过期时间
	//   4. 返回claims
	ValidateToken(tokenString string) (*Claims, error)

	// RefreshToken 刷新令牌（可选实现）
	// 参数:
	//   tokenString: 旧的JWT token字符串
	// 返回:
	//   string: 新的JWT token字符串
	//   error: 刷新失败时的错误
	// 注意:
	//   当前实现可能暂不支持此功能,返回 ErrNotImplemented
	RefreshToken(tokenString string) (string, error)
}

// Claims JWT载荷
// 包含用户身份信息和JWT标准声明
// 为什么需要自定义Claims:
// - 扩展性: 可以添加自定义字段
// - 类型安全: 明确字段类型
// - 标准遵循: 继承RegisteredClaims获得标准字段
type Claims struct {
	// UserID 用户ID
	// 用于标识token所属的用户
	UserID int64 `json:"user_id"`

	// Username 用户名
	// 用于显示或日志记录
	Username string `json:"username"`

	// jwt.RegisteredClaims 包含标准JWT字段:
	// - Issuer: 签发者
	// - Subject: 主题
	// - Audience: 受众
	// - ExpiresAt: 过期时间
	// - NotBefore: 生效时间
	// - IssuedAt: 签发时间
	// - ID: JWT唯一标识
	jwt.RegisteredClaims
}

// Config JWT配置
// 用于初始化JWT管理器
type Config struct {
	// Secret 签名密钥
	// 生产环境必须从环境变量设置
	// 建议使用至少32个字符的随机字符串
	Secret string

	// ExpiresIn 令牌有效期（秒）
	// 默认: 3600（1小时）
	// 考虑因素:
	// - 安全性: 过期时间越短越安全
	// - 用户体验: 过期时间太短需频繁登录
	// - 业务场景: 根据业务敏感度调整
	ExpiresIn int

	// Issuer 签发者
	// 标识令牌由哪个系统签发
	// 用于多系统环境下区分token来源
	Issuer string
}
