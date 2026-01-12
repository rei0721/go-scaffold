package app

import (
	"github.com/gin-gonic/gin"
	"github.com/rei0721/rei0721/internal/handler"
	"github.com/rei0721/rei0721/internal/middleware"
	"github.com/rei0721/rei0721/internal/repository"
	"github.com/rei0721/rei0721/internal/router"
	"github.com/rei0721/rei0721/internal/service"
)

func initBusiness(app *App) error {
	// 初始化 repository layer
	userRepo := repository.NewUserRepository(app.DB.DB())

	// 初始化 service layer (with dependency injection)
	userService := service.NewUserService(userRepo, app.Scheduler)

	// 初始化 handler layer
	userHandler := handler.NewUserHandler(userService)

	// 初始化 router
	r := router.New(userHandler, app.Logger)

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
