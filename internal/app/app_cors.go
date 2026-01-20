package app

import "github.com/rei0721/go-scaffold/internal/middleware"

// initCORS 初始化 CORS 中间件配置
// 从配置文件加载 CORS 配置，应用默认值并验证有效性
// 返回:
//
//	error: 初始化失败时的错误
//
// 执行步骤:
//  1. 获取 CORS 配置
//  2. 应用默认配置
//  3. 从环境变量覆盖
//  4. 验证配置有效性
//  5. 转换为中间件配置格式并存储
//
// 使用场景:
//
//	在应用初始化时调用，为路由器准备 CORS 配置
func (a *App) initCORS() error {
	// 获取 CORS 配置
	cfg := &a.Config.CORS

	// 应用默认配置
	// 为未配置的字段设置合理的默认值
	cfg.DefaultConfig()

	// 从环境变量覆盖
	// 生产环境可以通过环境变量覆盖配置文件中的值
	cfg.OverrideConfig()

	// 验证配置
	// 确保配置有效，否则提前失败
	if err := cfg.Validate(); err != nil {
		return err
	}

	// 记录配置状态
	if cfg.Enabled {
		a.Logger.Info("CORS middleware enabled",
			"allow_origins", cfg.AllowOrigins,
			"allow_credentials", cfg.AllowCredentials,
			"max_age", cfg.MaxAge,
		)
	} else {
		a.Logger.Info("CORS middleware disabled")
	}

	// CORS 配置会在路由器初始化时使用
	// 这里只需要验证和记录日志

	return nil
}

// getCORSMiddlewareConfig 获取 CORS 中间件配置
// 将应用配置转换为中间件配置格式
// 返回:
//
//	middleware.CORSConfig: CORS 中间件配置
//
// 使用场景:
//
//	在路由器初始化时调用，获取 CORS 配置
func (a *App) getCORSMiddlewareConfig() middleware.CORSConfig {
	cfg := a.Config.CORS

	return middleware.CORSConfig{
		Enabled:          cfg.Enabled,
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     cfg.AllowMethods,
		AllowHeaders:     cfg.AllowHeaders,
		ExposeHeaders:    cfg.ExposeHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           cfg.MaxAge,
	}
}
