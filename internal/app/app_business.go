package app

import (
	"github.com/gin-gonic/gin"
	"github.com/rei0721/rei0721/internal/handler"
	"github.com/rei0721/rei0721/internal/middleware"
	"github.com/rei0721/rei0721/internal/repository"
	"github.com/rei0721/rei0721/internal/router"
	"github.com/rei0721/rei0721/internal/service"
	rbacrepo "github.com/rei0721/rei0721/pkg/rbac/repository"
	rbacservice "github.com/rei0721/rei0721/pkg/rbac/service"
)

func initBusiness(app *App) error {
	// 初始化 repository layer
	userRepo := repository.NewUserRepository(app.DB.DB())

	// 初始化 service layer (不直接注入executor)
	userService := service.NewUserService(userRepo)

	// ⭐ 延迟注入 executor 到 Service 层
	if app.Executor != nil {
		userService.SetExecutor(app.Executor)
		app.Logger.Debug("executor injected into user service")
	}

	// ⭐ 延迟注入 cache 到 Service 层
	if app.Cache != nil {
		userService.SetCache(app.Cache)
		app.Logger.Debug("cache injected into user service")
	}

	// ⭐ 延迟注入 logger 到 Service 层
	if app.Logger != nil {
		userService.SetLogger(app.Logger)
		app.Logger.Debug("logger injected into user service")
	}

	// ⭐ 延迟注入 JWT 到 Service 层
	if app.JWT != nil {
		userService.SetJWT(app.JWT)
		app.Logger.Debug("JWT injected into user service")
	}

	// 初始化 handler layer
	userHandler := handler.NewUserHandler(userService)

	// 初始化 RBAC repository（使用 pkg/rbac）
	rbacRepo := rbacrepo.NewGormRBACRepository(app.DB.DB())

	// 初始化 RBAC service（使用 pkg/rbac）
	rbacService := rbacservice.NewRBACService(rbacRepo)

	// 注入 RBAC 依赖
	if app.Executor != nil {
		rbacService.SetExecutor(app.Executor)
	}
	if app.Cache != nil {
		rbacService.SetCache(app.Cache)
	}

	// 初始化 RBAC handler
	rbacHandler := handler.NewRBACHandler(rbacService, app.Logger)

	// 初始化 router（传入JWT用于认证中间件）
	r := router.New(userHandler, rbacHandler, app.Logger, app.JWT, rbacService)

	// Set Gin mode based on config
	if app.Config.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else if app.Config.Server.Mode == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Setup router with middleware
	middlewareCfg := middleware.DefaultMiddlewareConfig()
	app.Router = r.Setup(middlewareCfg)

	return nil
}
