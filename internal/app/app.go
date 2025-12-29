// Package app 提供依赖注入容器和应用程序生命周期管理
// 它按照正确的依赖顺序初始化所有组件,并提供优雅关闭功能
// 这是应用程序的核心,负责协调各个组件的创建和销毁
package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rei0721/rei0721/internal/config"
	"github.com/rei0721/rei0721/internal/handler"
	"github.com/rei0721/rei0721/internal/middleware"
	"github.com/rei0721/rei0721/internal/repository"
	"github.com/rei0721/rei0721/internal/router"
	"github.com/rei0721/rei0721/internal/service"
	"github.com/rei0721/rei0721/pkg/database"
	"github.com/rei0721/rei0721/pkg/logger"
	"github.com/rei0721/rei0721/pkg/scheduler"
)

// App 是主应用程序容器,持有所有组件并管理它们的生命周期
// 这是一个依赖注入(DI)容器模式的实现
// 优点:
// - 集中管理所有组件的创建和销毁
// - 明确的依赖关系,便于测试和维护
// - 支持优雅关闭,确保资源正确释放
type App struct {
	// Config 应用配置,从配置文件加载
	Config *config.Config

	// ConfigManager 配置管理器,支持配置热更新
	// 当配置文件变化时,可以动态重新加载
	ConfigManager config.Manager

	// DB 数据库连接抽象层
	// 使用接口而非具体实现,便于切换数据库
	DB database.Database

	// Scheduler 任务调度器,用于执行异步任务
	// 基于 ants 协程池实现,提高并发性能
	Scheduler scheduler.Scheduler

	// Logger 结构化日志记录器
	// 支持多种输出格式(JSON/控制台)和日志级别
	Logger logger.Logger

	// Router Gin HTTP 路由引擎
	// 包含所有HTTP路由和中间件配置
	Router *gin.Engine

	// server HTTP 服务器实例
	// 私有字段,仅供内部使用
	server *http.Server
}

// Options 创建新 App 时的配置选项
// 使用选项模式(Options Pattern)提高API的可扩展性
type Options struct {
	// ConfigPath 配置文件的路径
	// 支持相对路径和绝对路径
	ConfigPath string
}

// New 创建一个新的 App 实例
// 按照正确的依赖顺序初始化所有组件:
// config → logger → database → scheduler → repository → service → handler → router
// 为什么这个顺序:
// - config 最先:其他组件需要配置信息
// - logger 第二:后续初始化过程需要记录日志
// - database 第三:repository 依赖数据库
// - scheduler 第四:service 需要调度器执行异步任务
// - repository、service、handler、router 依次初始化,形成完整的请求链路
// 参数:
//
//	opts: 应用选项,包含配置文件路径等
//
// 返回:
//
//	*App: 初始化完成的应用实例
//	error: 初始化失败时的错误
func New(opts Options) (*App, error) {
	app := &App{}

	// 1. 初始化配置管理器并加载配置
	// 配置是整个应用的基础,必须最先加载
	configManager := config.NewManager()
	if err := configManager.Load(opts.ConfigPath); err != nil {
		// 配置加载失败,应用无法启动
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	app.ConfigManager = configManager
	app.Config = configManager.Get()

	// 2. 初始化日志记录器
	// 日志系统应该尽早初始化,便于记录后续的初始化过程
	log, err := logger.New(&logger.Config{
		Level:  app.Config.Logger.Level,  // 从配置读取日志级别
		Format: app.Config.Logger.Format, // 从配置读取日志格式
		Output: app.Config.Logger.Output, // 从配置读取输出目标
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}
	app.Logger = log

	// 将日志器注册到配置管理器
	// 这样配置变更时可以记录日志
	configManager.RegisterLogger(func() logger.Logger {
		return app.Logger
	})

	// 3. 初始化数据库连接
	db, err := database.New(&database.Config{
		Driver:       database.Driver(app.Config.Database.Driver),
		Host:         app.Config.Database.Host,
		Port:         app.Config.Database.Port,
		User:         app.Config.Database.User,
		Password:     app.Config.Database.Password,
		DBName:       app.Config.Database.DBName,
		MaxOpenConns: app.Config.Database.MaxOpenConns,
		MaxIdleConns: app.Config.Database.MaxIdleConns,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	app.DB = db
	app.Logger.Info("database connected successfully")

	// 4. Initialize scheduler
	sched, err := scheduler.New(&scheduler.Config{
		PoolSize:       10000,
		ExpiryDuration: time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}
	app.Scheduler = sched
	app.Logger.Info("scheduler initialized", "poolSize", 10000)

	// 5. Initialize repository layer
	userRepo := repository.NewUserRepository(db.DB())

	// 6. Initialize service layer (with dependency injection)
	userService := service.NewUserService(userRepo, sched)

	// 7. Initialize handler layer
	userHandler := handler.NewUserHandler(userService)

	// 8. Initialize router
	r := router.New(userHandler, log)

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

	// 9. Start config file watching for hot-reload
	if err := configManager.Watch(); err != nil {
		app.Logger.Warn("failed to start config watcher", "error", err)
	}

	// Register config change hook
	configManager.RegisterHook(func(old, new *config.Config) {
		app.Logger.Info("configuration updated",
			"oldPort", old.Server.Port,
			"newPort", new.Server.Port,
		)
		app.Config = new
	})

	app.Logger.Info("application initialized successfully")
	return app, nil
}

// Run 启动 HTTP 服务器并阻塞直到服务器停止
// 这个方法会一直运行直到:
// - 服务器发生错误
// - Shutdown() 被调用
// 返回:
//
//	error: 服务器启动失败或运行时错误
func (a *App) Run() error {
	// 构造监听地址
	// 格式: ":8080" 表示监听所有网络接口的 8080 端口
	addr := fmt.Sprintf(":%d", a.Config.Server.Port)

	// 创建 HTTP 服务器实例
	// 配置超时参数防止慢速客户端占用连接
	a.server = &http.Server{
		// Addr 监听地址
		Addr: addr,

		// Handler HTTP 请求处理器(Gin Router)
		Handler: a.Router,

		// ReadTimeout 读取请求的最大时间
		// 包括读取请求头和请求体
		// 防止慢速客户端长时间占用连接
		ReadTimeout: time.Duration(a.Config.Server.ReadTimeout) * time.Second,

		// WriteTimeout 写入响应的最大时间
		// 从请求处理完成到写入完整响应
		// 防止慢速客户端长时间占用连接
		WriteTimeout: time.Duration(a.Config.Server.WriteTimeout) * time.Second,
	}

	a.Logger.Info("starting HTTP server", "addr", addr)

	// 启动服务器并开始监听
	// ListenAndServe 会阻塞,直到:
	// - 发生错误
	// - Shutdown() 被调用
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		// ErrServerClosed 是正常的关闭,不是错误
		// 其他错误才需要返回
		return fmt.Errorf("failed to start server: %w", err)
	}

	// 正常关闭,返回 nil
	return nil
}

// Shutdown 优雅地关闭应用程序
// 关闭顺序很重要:
// 1. HTTP 服务器 - 停止接收新请求
// 2. 调度器 - 等待异步任务完成
// 3. 数据库 - 关闭连接
// 4. 日志器 - 刷新缓冲区
// 参数:
//
//	ctx: 上下文,用于控制关闭超时
//
// 返回:
//
//	error: 关闭过程中的错误(可能包含多个)
//
// 设计考虑:
//   - 即使某个组件关闭失败,也继续关闭其他组件
//   - 收集所有错误并返回
//   - 使用 context 控制整体超时
func (a *App) Shutdown(ctx context.Context) error {
	a.Logger.Info("shutting down application...")

	// 收集所有关闭过程中的错误
	var errs []error

	// 1. 关闭 HTTP 服务器
	// 步骤:
	// - 停止接收新连接
	// - 等待现有请求处理完成
	// - 或者直到 context 超时
	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil {
			// 关闭失败,记录错误但继续关闭其他组件
			a.Logger.Error("failed to shutdown HTTP server", "error", err)
			errs = append(errs, fmt.Errorf("http server shutdown: %w", err))
		} else {
			a.Logger.Info("HTTP server stopped")
		}
	}

	// 2. 关闭调度器(等待运行中的任务)
	// 步骤:
	// - 停止接收新任务
	// - 等待运行中的任务完成
	// - 释放协程池资源
	if a.Scheduler != nil {
		if err := a.Scheduler.Shutdown(ctx); err != nil {
			a.Logger.Error("failed to shutdown scheduler", "error", err)
			errs = append(errs, fmt.Errorf("scheduler shutdown: %w", err))
		} else {
			a.Logger.Info("scheduler stopped")
		}
	}

	// 3. 关闭数据库连接
	// 步骤:
	// - 关闭所有连接池中的连接
	// - 释放相关资源
	// 注意: 此时不应该有活跃的数据库操作
	if a.DB != nil {
		if err := a.DB.Close(); err != nil {
			a.Logger.Error("failed to close database", "error", err)
			errs = append(errs, fmt.Errorf("database close: %w", err))
		} else {
			a.Logger.Info("database connection closed")
		}
	}

	// 4. 同步日志器
	// 确保所有缓冲的日志都写入磁盘
	// 这应该是最后一步,确保所有关闭日志都被记录
	if a.Logger != nil {
		// 忽略 Sync 的错误
		// 某些平台(如 Linux /dev/stdout)可能返回无害的错误
		_ = a.Logger.Sync()
	}

	// 检查是否有错误发生
	if len(errs) > 0 {
		// 有错误但已尽力关闭所有组件
		return fmt.Errorf("shutdown completed with %d errors", len(errs))
	}

	a.Logger.Info("application shutdown complete")
	return nil
}
