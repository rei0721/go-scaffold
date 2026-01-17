package auth

import "github.com/rei0721/go-scaffold/pkg/crypto"

// SetCrypto 设置密码加密器（延迟注入）
func (s *authService) SetCrypto(c crypto.Crypto) {
	s.Crypto = c
}
