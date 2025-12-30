// Package app 提供依赖注入容器和应用程序生命周期管理
// 它按照正确的依赖顺序初始化所有组件,并提供优雅关闭功能
// 这是应用程序的核心,负责协调各个组件的创建和销毁
package app

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rei0721/rei0721/pkg/cache"
	"github.com/rei0721/rei0721/pkg/daemon"
	"github.com/rei0721/rei0721/pkg/i18n"

	"github.com/rei0721/rei0721/internal/config"
	"github.com/rei0721/rei0721/internal/daemons"
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

	// I18n 国际化
	I18n i18n.I18n

	// Cache Redis 缓存
	// 用于提高性能,减轻数据库压力
	// 如果 Redis 未启用,此字段为 nil
	Cache cache.Cache

	// Scheduler 任务调度器,用于执行异步任务
	// 基于 ants 协程池实现,提高并发性能
	Scheduler scheduler.Scheduler

	// Logger 结构化日志记录器
	// 支持多种输出格式(JSON/控制台)和日志级别
	Logger logger.Logger

	// Router Gin HTTP 路由引擎
	// 包含所有HTTP路由和中间件配置
	Router *gin.Engine

	// DaemonManager 守护进程管理器
	// 统一管理所有长期运行的服务(HTTP、gRPC、Kafka 等)
	// 使用 daemon.Manager 实现优雅启动和关闭
	DaemonManager *daemon.Manager
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
// config → logger → i18n → database → scheduler → repository → service → handler → router
// 为什么这个顺序:
// - config 最先:其他组件需要配置信息
// - logger 第二:后续初始化过程需要记录日志
// - i18n 第三:响应HTTP需用到
// - database 第四:repository 依赖数据库
// - scheduler 第五:service 需要调度器执行异步任务
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

	// 记录配置加载信息
	// 展示环境变量支持功能已生效
	app.Logger.Info("configuration loaded successfully",
		"config_file", opts.ConfigPath,
		"env_support", "enabled")

	// 记录关键配置信息（不记录敏感信息）
	app.Logger.Info("server configuration",
		"port", app.Config.Server.Port,
		"mode", app.Config.Server.Mode)

	app.Logger.Info("database configuration",
		"driver", app.Config.Database.Driver,
		"host", app.Config.Database.Host,
		"db", app.Config.Database.DBName)

	if app.Config.Redis.Enabled {
		app.Logger.Info("redis configuration",
			"enabled", true,
			"host", app.Config.Redis.Host,
			"db", app.Config.Redis.DB)
	} else {
		app.Logger.Info("redis configuration", "enabled", false)
	}

	// 3. 初始化i18n
	i18nCfg := &i18n.Config{
		DefaultLanguage:    ConstantsI18nDefaultLanguage,
		SupportedLanguages: ConstantsI18nSupportedLanguages,
		MessagesDir:        ConstantsI18nMessagesDir,
	}
	i18nApp, i18nErr := i18n.New(i18nCfg)
	if i18nErr != nil {
		return nil, fmt.Errorf("failed to create i18n: %w", i18nErr)
	}
	app.I18n = i18nApp

	// 4. 初始化 Redis 缓存(可选)
	// 如果配置中启用了 Redis,则创建缓存实例
	if app.Config.Redis.Enabled {
		cacheCfg := &cache.Config{
			Host:         app.Config.Redis.Host,
			Port:         app.Config.Redis.Port,
			Password:     app.Config.Redis.Password,
			DB:           app.Config.Redis.DB,
			PoolSize:     app.Config.Redis.PoolSize,
			MinIdleConns: app.Config.Redis.MinIdleConns,
			MaxRetries:   app.Config.Redis.MaxRetries,
			DialTimeout:  time.Duration(app.Config.Redis.DialTimeout) * time.Second,
			ReadTimeout:  time.Duration(app.Config.Redis.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(app.Config.Redis.WriteTimeout) * time.Second,
		}

		cacheClient, err := cache.NewRedis(cacheCfg, log)
		if err != nil {
			// Redis 连接失败
			// 可以选择:
			// 1. 返回错误,强制要求 Redis 可用
			// 2. 警告但继续,允许无缓存运行
			// 这里选择方案 2,提高可用性
			app.Logger.Warn("failed to connect to redis, running without cache", "error", err)
			app.Cache = nil
		} else {
			app.Cache = cacheClient
			app.Logger.Info("redis cache connected successfully")
		}
	} else {
		app.Logger.Info("redis cache disabled")
		app.Cache = nil
	}

	// 5. 初始化数据库连接
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

	// 5. Initialize scheduler
	sched, err := scheduler.New(&scheduler.Config{
		PoolSize:       10000,
		ExpiryDuration: time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}
	app.Scheduler = sched
	app.Logger.Info("scheduler initialized", "poolSize", 10000)

	// 6. Initialize repository layer
	userRepo := repository.NewUserRepository(db.DB())

	// 7. Initialize service layer (with dependency injection)
	userService := service.NewUserService(userRepo, sched)

	// 8. Initialize handler layer
	userHandler := handler.NewUserHandler(userService)

	// 9. Initialize router
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

	// 10. Start config file watching for hot-reload
	if err := configManager.Watch(); err != nil {
		app.Logger.Warn("failed to start config watcher", "error", err)
	}

	// Register config change hook
	// 当配置文件变化时自动调用
	configManager.RegisterHook(func(old, new *config.Config) {
		app.Logger.Info("configuration file changed, processing updates...")

		// 1. 检查 Redis 配置是否变化
		// 使用工具方法比较新旧配置
		if isRedisConfigChanged(old, new) {
			app.Logger.Info("redis configuration changed, reloading cache...")

			// 只有在 Cache 不为 nil 且新配置启用了 Redis 时才重载
			if app.Cache != nil && new.Redis.Enabled {
				// 创建新的缓存配置
				newCacheCfg := &cache.Config{
					Host:         new.Redis.Host,
					Port:         new.Redis.Port,
					Password:     new.Redis.Password,
					DB:           new.Redis.DB,
					PoolSize:     new.Redis.PoolSize,
					MinIdleConns: new.Redis.MinIdleConns,
					MaxRetries:   new.Redis.MaxRetries,
					DialTimeout:  time.Duration(new.Redis.DialTimeout) * time.Second,
					ReadTimeout:  time.Duration(new.Redis.ReadTimeout) * time.Second,
					WriteTimeout: time.Duration(new.Redis.WriteTimeout) * time.Second,
				}

				// 使用超时上下文进行重载
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				// 原子化重载缓存配置
				err := app.Cache.Reload(ctx, newCacheCfg)
				if err != nil {
					app.Logger.Error("failed to reload redis cache", "error", err)
				} else {
					app.Logger.Info("redis cache reloaded successfully")
				}
			} else if !new.Redis.Enabled {
				app.Logger.Info("redis disabled in new config")
			} else {
				app.Logger.Warn("cache is nil, cannot reload redis configuration")
			}
		}

		// 2. 检查服务器配置是否变化
		if isServerConfigChanged(old, new) {
			app.Logger.Info("server configuration changed",
				"oldPort", old.Server.Port,
				"newPort", new.Server.Port)
		}

		// 3. 更新应用配置引用
		app.Config = new
		app.Logger.Info("configuration update completed")
	})

	// 11. 初始化守护进程管理器
	// 用于统一管理所有长期运行的服务
	app.DaemonManager = daemon.NewManager(log)

	// 注册 HTTP 守护进程
	// 将 HTTP 服务器封装为 daemon 并注册到管理器
	addr := fmt.Sprintf(":%d", app.Config.Server.Port)
	httpDaemon := daemons.NewHTTPDaemon(addr, app.Config.Server.ReadTimeout, app.Config.Server.WriteTimeout, app.Router, log)
	app.DaemonManager.Register(httpDaemon)

	app.Logger.Info("application initialized successfully")
	return app, nil
}

// Start 启动所有守护进程
// 这个方法会启动所有注册的守护进程(包括 HTTP 服务器)
// 并且不会阻塞,服务在后台运行
// 参数:
//
//	ctx: 上下文,用于控制启动过程
//
// 返回:
//
//	error: 启动失败时的错误
func (a *App) Start(ctx context.Context) error {
	// 启动所有守护进程
	// Manager.Start() 会并发启动所有注册的守护进程
	if err := a.DaemonManager.Start(ctx); err != nil {
		return fmt.Errorf("failed to start daemons: %w", err)
	}

	a.Logger.Info("all daemons started successfully")
	return nil
}

// Run 启动应用并阻塞直到接收到停止信号
// 这个方法是为了保持向后兼容性
// 实际上它只是调用 Start() 然后阻塞
// 返回:
//
//	error: 启动失败时的错误
func (a *App) Run() error {
	// 启动所有守护进程
	if err := a.Start(context.Background()); err != nil {
		return err
	}

	// 阻塞,直到 Shutdown() 被调用
	// 这里使用一个简单的 select {}
	// 实际应用中通常会监听信号来控制关闭
	select {}
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

	// 1. 关闭所有守护进程(包括 HTTP 服务器)
	// DaemonManager.Stop() 会并发关闭所有注册的守护进程
	// 步骤:
	// - 停止接收新连接/任务
	// - 等待现有连接/任务完成
	// - 或者直到 context 超时
	if a.DaemonManager != nil {
		if err := a.DaemonManager.Stop(ctx); err != nil {
			// 关闭失败,记录错误但继续关闭其他组件
			a.Logger.Error("failed to stop daemons", "error", err)
			errs = append(errs, fmt.Errorf("daemons shutdown: %w", err))
		} else {
			a.Logger.Info("all daemons stopped")
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

	// 3. 关闭缓存连接
	// 步骤:
	// - 关闭 Redis 连接
	// - 释放连接池资源
	if a.Cache != nil {
		if err := a.Cache.Close(); err != nil {
			a.Logger.Error("failed to close cache", "error", err)
			errs = append(errs, fmt.Errorf("cache close: %w", err))
		} else {
			a.Logger.Info("cache closed")
		}
	}

	// 4. 关闭数据库连接
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
