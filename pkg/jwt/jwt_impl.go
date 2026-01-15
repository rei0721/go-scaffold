package jwt

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtManager 实现 JWT 接口
// 这是JWT管理器的具体实现,负责token的生成和验证
// 设计原则:
// - 线程安全: 使用 RWMutex 保护配置
// - 配置驱动: 通过Config初始化
// - 错误明确: 提供清晰的错误信息
type jwtManager struct {
	// secret 签名密钥
	// 使用HMAC-SHA256算法时的密钥
	// 必须保密,不能泄露
	secret []byte

	// expiresIn token有效期
	// 从签发时间开始计算
	expiresIn time.Duration

	// issuer 签发者标识
	// 用于标识token的来源
	issuer string

	// mu 读写锁
	// 保护配置字段的并发访问
	// 读多写少的场景使用RWMutex性能更好
	mu sync.RWMutex
}

// New 创建一个新的 JWT 管理器实例
// 这是工厂函数,遵循依赖注入模式
// 参数:
//
//	cfg: JWT配置对象
//
// 返回:
//
//	JWT: JWT接口实例
//	error: 创建失败时的错误
//
// 验证规则:
//
//  1. secret不能为空
//  2. secret长度至少32个字符（安全性考虑）
//  3. expiresIn必须大于0
func New(cfg *Config) (JWT, error) {
	// 1. 验证配置
	if cfg.Secret == "" {
		return nil, ErrMissingSecret
	}

	// 2. 验证密钥长度（安全性要求）
	if len(cfg.Secret) < 32 {
		return nil, errors.New(ErrMsgSecretTooShort)
	}

	// 3. 设置默认值
	expiresIn := cfg.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = DefaultExpiresIn
	}

	issuer := cfg.Issuer
	if issuer == "" {
		issuer = DefaultIssuer
	}

	// 4. 创建实例
	return &jwtManager{
		secret:    []byte(cfg.Secret),
		expiresIn: time.Duration(expiresIn) * time.Second,
		issuer:    issuer,
	}, nil
}

// GenerateToken 生成访问令牌
// 实现JWT接口的GenerateToken方法
// 参数:
//
//	userID: 用户ID
//	username: 用户名
//
// 返回:
//
//	string: JWT token字符串
//	error: 生成失败时的错误
//
// 业务流程:
//  1. 创建claims载荷
//  2. 创建JWT token对象
//  3. 使用HMAC-SHA256算法签名
//  4. 生成完整的token字符串
func (m *jwtManager) GenerateToken(userID int64, username string) (string, error) {
	// 使用读锁保护配置读取
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 1. 创建claims
	now := time.Now()
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// 签发者
			Issuer: m.issuer,

			// 签发时间
			IssuedAt: jwt.NewNumericDate(now),

			// 过期时间 = 当前时间 + 有效期
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expiresIn)),

			// 生效时间 = 当前时间
			// token立即生效,不设置延迟
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// 2. 创建token对象
	// SigningMethodHS256 使用HMAC-SHA256算法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. 签名并生成token字符串
	// SignedString会:
	// - 将header和claims编码为base64
	// - 使用secret对它们进行HMAC-SHA256签名
	// - 拼接成完整的JWT: header.claims.signature
	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken 验证并解析令牌
// 实现JWT接口的ValidateToken方法
// 参数:
//
//	tokenString: JWT token字符串
//
// 返回:
//
//	*Claims: 解析后的载荷信息
//	error: 验证失败时的错误
//
// 验证步骤:
//  1. 解析token字符串
//  2. 验证签名
//  3. 检查过期时间
//  4. 检查生效时间
//  5. 提取claims
func (m *jwtManager) ValidateToken(tokenString string) (*Claims, error) {
	// 使用读锁保护配置读取
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 1. 解析token
	// ParseWithClaims会:
	// - 解析token字符串
	// - 使用keyFunc验证签名
	// - 检查标准声明（过期时间、生效时间等）
	// - 将载荷解析到Claims结构
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		// 防止攻击者使用其他算法（如none）绕过签名验证
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// 返回密钥用于验证签名
		return m.secret, nil
	})

	// 2. 处理解析错误
	if err != nil {
		// 根据错误类型返回更具体的错误
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotYetValid
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, ErrInvalidSignature
		}
		// 其他错误统一返回无效token
		return nil, ErrInvalidToken
	}

	// 3. 提取claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		// token解析成功但claims类型不匹配或token无效
		return nil, ErrInvalidToken
	}

	// 4. 返回claims
	return claims, nil
}

// RefreshToken 刷新令牌
// 实现JWT接口的RefreshToken方法
// 注意: 当前实现为占位符,可根据需求实现
// 典型的实现方式:
//  1. 验证旧token（允许一定程度的过期）
//  2. 提取用户信息
//  3. 生成新token
//  4. 可选：记录刷新事件
func (m *jwtManager) RefreshToken(tokenString string) (string, error) {
	// 1. 验证旧token
	// 注意: 某些实现允许在一定时间窗口内刷新过期的token
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		// 如果token已过期但在允许的刷新窗口内,可以继续
		// 这里简化处理,不允许过期token刷新
		return "", err
	}

	// 2. 生成新token
	// 使用相同的用户信息,但更新时间戳
	return m.GenerateToken(claims.UserID, claims.Username)
}
