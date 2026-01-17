package app

import (
	"fmt"

	"github.com/rei0721/go-scaffold/pkg/jwt"
)

// initJWT 初始化 JWT 认证管理器
// 用于生成和验证访问令牌
// 参数:
//
//	app: 应用实例
//
// 返回:
//
//	error: 初始化失败时的错误
func initJWT(app *App) error {
	app.Logger.Info("Initializing JWT manager...")

	// 创建 JWT 配置
	jwtCfg := &jwt.Config{
		Secret:    app.Config.JWT.Secret,
		ExpiresIn: app.Config.JWT.ExpiresIn,
		Issuer:    app.Config.JWT.Issuer,
	}

	// 创建 JWT 管理器
	jwtManager, err := jwt.New(jwtCfg)
	if err != nil {
		return fmt.Errorf("failed to create JWT manager: %w", err)
	}

	app.JWT = jwtManager
	app.Logger.Info("JWT manager initialized successfully",
		"expires_in", app.Config.JWT.ExpiresIn,
		"issuer", app.Config.JWT.Issuer)

	return nil
}
