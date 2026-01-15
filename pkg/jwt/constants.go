package jwt

import "errors"

// 默认配置
const (
	// DefaultExpiresIn 默认过期时间（1小时）
	DefaultExpiresIn = 3600

	// DefaultIssuer 默认签发者
	DefaultIssuer = "go-scaffold"
)

// 预定义错误
var (
	// ErrInvalidToken token 无效
	ErrInvalidToken = errors.New("invalid token")

	// ErrExpiredToken token 已过期
	ErrExpiredToken = errors.New("token has expired")

	// ErrTokenNotYetValid token 尚未生效
	ErrTokenNotYetValid = errors.New("token not yet valid")

	// ErrInvalidSignature 签名无效
	ErrInvalidSignature = errors.New("invalid signature")

	// ErrMissingSecret 缺少签名密钥
	ErrMissingSecret = errors.New("jwt secret is required")
)

// 错误消息常量
const (
	// ErrMsgInvalidToken token 无效错误消息
	ErrMsgInvalidToken = "invalid token"

	// ErrMsgExpiredToken token 过期错误消息
	ErrMsgExpiredToken = "token has expired"

	// ErrMsgTokenNotYetValid token 未生效错误消息
	ErrMsgTokenNotYetValid = "token not yet valid"

	// ErrMsgInvalidSignature 签名无效错误消息
	ErrMsgInvalidSignature = "invalid signature"

	// ErrMsgMissingSecret 缺少密钥错误消息
	ErrMsgMissingSecret = "jwt secret is required"

	// ErrMsgSecretTooShort 密钥太短错误消息
	ErrMsgSecretTooShort = "jwt secret must be at least 32 characters"
)
