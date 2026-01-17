package config

import "errors"

// JWTConfig JWT认证配置
// 用于token的生成和验证
type JWTConfig struct {
	// Secret 签名密钥
	// 生产环境必须从环境变量设置
	// 建议使用至少32个字符的随机字符串
	// 注意: 此字段非常敏感,必须保密
	Secret string `mapstructure:"secret"`

	// ExpiresIn 令牌有效期（秒）
	// 默认: 3600（1小时）
	// 考虑因素:
	// - 安全性: 过期时间越短越安全
	// - 用户体验: 过期时间太短需频繁登录
	// - 业务场景: 根据业务敏感度调整
	ExpiresIn int `mapstructure:"expiresIn"`

	// Issuer 签发者
	// 标识令牌由哪个系统签发
	// 用于多系统环境下区分token来源
	// 默认: "go-scaffold"
	Issuer string `mapstructure:"issuer"`
}

func (c *JWTConfig) ValidateName() string {
	return AppJWTName
}

func (c *JWTConfig) ValidateRequired() bool {
	return true
}

// Validate 验证 JWT 配置
// 实现 Configurable 接口
func (c *JWTConfig) Validate() error {
	// 验证密钥
	if c.Secret == "" {
		return errors.New("jwt secret is required")
	}

	// 验证密钥长度（安全性要求）
	if len(c.Secret) < 32 {
		return errors.New("jwt secret must be at least 32 characters")
	}

	// 验证过期时间
	if c.ExpiresIn <= 0 {
		return errors.New("jwt expiresIn must be positive")
	}

	return nil
}
