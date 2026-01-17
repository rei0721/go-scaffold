// Package app 提供依赖注入容器和应用程序生命周期管理
// 它按照正确的依赖顺序初始化所有组件,并提供优雅关闭功能
// 这是应用程序的核心,负责协调各个组件的创建和销毁
package app

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/rei0721/go-scaffold/pkg/cache"
	"github.com/rei0721/go-scaffold/pkg/executor"
	"github.com/rei0721/go-scaffold/pkg/httpserver"
	"github.com/rei0721/go-scaffold/pkg/i18n"
	"github.com/rei0721/go-scaffold/pkg/jwt"
	"github.com/rei0721/go-scaffold/pkg/rbac"
	"github.com/rei0721/go-scaffold/pkg/sqlgen"
	"github.com/rei0721/go-scaffold/pkg/utils"

	"github.com/rei0721/go-scaffold/internal/config"
	"github.com/rei0721/go-scaffold/pkg/database"
	"github.com/rei0721/go-scaffold/pkg/logger"
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

	// Sqlgen SQL 生成器
	// 用于生成数据库建表语句
	Sqlgen *sqlgen.Generator

	// I18n 国际化
	I18n      i18n.I18n
	I18nUtils *utils.I18nUtils

	// Cache Redis 缓存
	// 用于提高性能,减轻数据库压力
	// 如果 Redis 未启用,此字段为 nil
	Cache cache.Cache

	// Executor 异步任务执行器
	// 管理多个协程池,支持动态热重载
	// 如果 Executor 未启用,此字段为 nil
	Executor executor.Manager

	// Logger 结构化日志记录器
	// 支持多种输出格式(JSON/控制台)和日志级别
	Logger logger.Logger

	// Router Gin HTTP 路由引擎
	// 包含所有HTTP路由和中间件配置
	Router *gin.Engine

	// HTTPServer HTTP 服务器实例
	// 使用 pkg/httpserver 接口，支持配置热更新
	HTTPServer httpserver.HTTPServer

	// JWT JWT认证管理器
	// 用于生成和验证访问令牌
	JWT jwt.JWT

	// RBAC 角色访问控制
	// 管理用户权限和角色
	RBAC rbac.RBAC

	// Options 应用选项
	Options Options
}

// Options 创建新 App 时的配置选项
// 使用选项模式(Options Pattern)提高API的可扩展性
type Options struct {
	// ConfigPath 配置文件的路径
	// 支持相对路径和绝对路径
	ConfigPath string

	// Mode 启动模式
	// 支持 ModeServer（默认）和 ModeInitDB 两种模式
	Mode AppMode
}

// New 创建一个新的 App 实例
//
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

	// 设置默认模式
	if opts.Mode == "" {
		opts.Mode = ModeServer
	}

	// 备份选项
	app.Options = opts

	// 初始化配置管理器并加载配置
	// 配置是整个应用的基础,必须最先加载
	if err := app.initConfig(opts); err != nil {
		// 配置加载失败,应用无法启动
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 初始化日志记录器
	// 日志系统应该尽早初始化,便于记录后续的初始化过程
	if err := app.initLogger(); err != nil {
		return nil, err
	}

	// 初始化 i18n
	if err := app.initI18n(); err != nil {
		return nil, err
	}

	// debug
	app.Logger.Debug(app.UI18n("internal.app.logger_debug_init_config_path"), "config_path", opts.ConfigPath)
	app.Logger.Debug(app.UI18n("internal.app.logger_debug_init_mode"), "mode", opts.Mode)

	// 将日志器注册到配置管理器
	// 这样配置变更时可以记录日志
	app.ConfigManager.RegisterLogger(func() logger.Logger {
		return app.Logger
	})

	// Config 调试配置
	debugConfig(app, opts)

	// 初始化模式
	return app.initMode()
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
	// 启动 HTTP 服务器（非阻塞）
	// 使用新的 httpserver 包
	if err := a.HTTPServer.Start(ctx); err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

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
// 1. 守护进程 - 停止接收新请求
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

	// 关闭 HTTP 服务器
	// 使用新的 httpserver 接口
	// 步骤:
	// - 停止接收新连接
	// - 等待现有请求处理完成
	// - 或者直到 context 超时
	if a.HTTPServer != nil {
		if err := a.HTTPServer.Shutdown(ctx); err != nil {
			// 关闭失败,记录错误但继续关闭其他组件
			a.Logger.Error("failed to shutdown HTTP server", "error", err)
			errs = append(errs, fmt.Errorf("http server shutdown: %w", err))
		} else {
			a.Logger.Info("HTTP server stopped")
		}
	}

	// 关闭 RBAC
	if a.RBAC != nil {
		a.RBAC.Close()
		a.Logger.Info("rbac stopped")
	}

	// 关闭执行器(等待运行中的任务)
	// 步骤:
	// - 停止接收新任务
	// - 等待运行中的任务完成
	// - 释放协程池资源
	if a.Executor != nil {
		a.Executor.Shutdown()
		a.Logger.Info("executor stopped")
	}

	// 关闭缓存连接
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

	// 关闭数据库连接
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

	// 同步日志器
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
