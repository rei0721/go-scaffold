package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// hashPassword 使用 bcrypt 加密密码
// 参数:
//
//	password: 明文密码
//
// 返回:
//
//	string: 加密后的密码哈希
//	error: 加密失败的错误
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// verifyPassword 验证密码是否匹配
// 参数:
//
//	hashedPassword: 数据库中存储的密码哈希
//	password: 用户输入的明文密码
//
// 返回:
//
//	error: 密码不匹配或验证失败的错误
func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
