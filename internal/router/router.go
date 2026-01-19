package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rei0721/go-scaffold/internal/handler"
	"github.com/rei0721/go-scaffold/internal/middleware"
	rbacService "github.com/rei0721/go-scaffold/internal/service/rbac"
	"github.com/rei0721/go-scaffold/pkg/jwt"
	"github.com/rei0721/go-scaffold/pkg/logger"
	"github.com/rei0721/go-scaffold/types/constants"
	"github.com/rei0721/go-scaffold/types/result"
)

// Router 管理 HTTP 路由配置和注册
// 它持有必要的依赖(处理器、日志器)用于路由设置
// 设计考虑:
// - 集中管理所有路由,便于维护
// - 通过依赖注入获取处理器,遵循 DI 原则
// - 分离路由定义和处理逻辑,保持清晰的分层
type Router struct {
	// engine Gin 引擎实例
	// 这是 Gin 框架的核心,负责 HTTP 请求路由和处理
	engine *gin.Engine

	// userHandler 用户相关的请求处理器
	// 包含注册、登录、查询等处理方法
	// 通过接口注入,便于测试和替换实现
	authHandler *handler.AuthHandler

	// rbacHandler RBAC权限管理处理器
	rbacHandler *handler.RBACHandler

	// logger 日志记录器
	// 用于记录路由相关的日志
	// 也会传递给中间件使用
	logger logger.Logger

	// jwt JWT管理器
	// 用于认证中间件验证token
	// 如果为nil,则不启用认证保护
	jwt jwt.JWT

	// rbacService RBAC服务
	// 用于中间件权限检查
	rbacService rbacService.RBACService
}

// New 创建一个新的 Router 实例
// 这是工厂函数,遵循依赖注入模式
// 参数:
//
//	authHandler: 认证处理器
//	rbacHandler: RBAC处理器
//	log: 日志记录器,用于记录日志
//	jwtManager: JWT管理器,用于认证中间件(可选,为nil时不启用认证保护)
//	rbacSvc: RBAC服务,用于中间件权限检查(可选)
//
// 返回:
//
//	*Router: 路由器实例
//
// 使用场景:
//
//	在应用初始化时创建,然后调用 Setup() 配置路由
func New(authHandler *handler.AuthHandler, rbacHandler *handler.RBACHandler, log logger.Logger, jwtManager jwt.JWT, rbacSvc rbacService.RBACService) *Router {
	return &Router{
		authHandler: authHandler,
		rbacHandler: rbacHandler,
		logger:      log,
		jwt:         jwtManager,
		rbacService: rbacSvc,
	}
}

// Setup 初始化 Gin 引擎并配置中间件和路由
// 这个方法完成路由器的完整设置
// 参数:
//
//	cfg: 中间件配置,包含各个中间件的启用状态和参数
//
// 返回:
//
//	*gin.Engine: 配置好的 Gin 引擎,可以直接用于启动 HTTP 服务器
//
// 中间件顺序说明:
//
//	必须按 TraceID -> Logger -> Recovery 的顺序,原因:
//	1. TraceID 必须最先,因为后续中间件都需要使用 TraceID
//	2. Logger 在中间,可以记录包含 TraceID 的请求日志
//	3. Recovery 必须最后,才能捕获所有中间件和处理器的 panic
func (r *Router) Setup(cfg middleware.MiddlewareConfig) *gin.Engine {
	// 创建一个新的 Gin 引擎
	// 使用 gin.New() 而不是 gin.Default()
	// gin.Default() 会自动添加 Logger 和 Recovery 中间件
	// 我们使用 gin.New() 从零开始,完全控制中间件
	r.engine = gin.New()

	// 1. 应用 TraceID 中间件(必须第一个)
	// 为每个请求生成或提取 TraceID
	// 后续所有中间件和处理器都可以使用这个 TraceID
	// 这是分布式追踪的基础
	r.engine.Use(middleware.TraceID(cfg.TraceID))

	// 2. 应用 Logger 中间件
	// 记录每个请求的详细信息:方法、路径、状态码、耗时、TraceID 等
	// 这对于监控、调试和问题排查至关重要
	// 可以在配置中指定跳过某些路径(如健康检查)
	r.engine.Use(middleware.Logger(cfg.Logger, r.logger))

	// 3. 应用 Recovery 中间件(必须最后)
	// 捕获所有 panic,防止服务崩溃
	// 必须在所有其他中间件之后,才能捕获它们的 panic
	// 发生 panic 时会记录日志并返回 500 错误
	r.engine.Use(middleware.Recovery(cfg.Recovery, r.logger))

	// 注册所有应用路由
	// 包括健康检查、API 路由等
	r.registerRoutes()

	// 返回配置好的引擎
	return r.engine
}

// registerRoutes 注册所有应用路由
// 这个方法定义了应用的 URL 结构
// 设计考虑:
// - 使用路由分组,保持 URL 层次清晰
// - 版本化 API(/api/v1),便于未来升级
// - RESTful 风格,语义清晰
func (r *Router) registerRoutes() {
	// 健康检查端点
	// GET /health
	// 用途:
	// - K8s/Docker 的健康探针
	// - 负载均衡器的健康检查
	// - 监控系统检查服务是否存活
	// 这个端点应该:
	// - 响应快速(不访问数据库)
	// - 始终返回 200(除非服务真的挂了)
	// - 不需要认证
	r.engine.GET("/health", r.healthCheck)

	// API v1 路由组
	// 所有 v1 API 都在 /api/v1 路径下
	// 好处:
	// - 明确的版本标识
	// - 可以同时运行多个版本(v1, v2)
	// - URL 清晰,易于理解
	v1 := r.engine.Group("/api/v1")
	{
		// ==================== 公开路由 ====================
		// 这些路由不需要认证即可访问

		// 认证相关路由组
		// 公开接口：注册和登录不需要认证
		authGroup := v1.Group("/auth")
		{
			// POST /api/v1/auth/register - 用户注册
			// 任何人都可以注册账号
			authGroup.POST("/register", r.authHandler.Register)

			// POST /api/v1/auth/login - 用户登录
			// 使用用户名和密码登录获取 token
			authGroup.POST("/login", r.authHandler.Login)
		}

		// 需要认证的认证路由
		// 这些接口需要提供有效的 JWT token
		authProtected := v1.Group("/auth")
		authProtected.Use(middleware.AuthMiddleware(r.jwt))
		{
			// POST /api/v1/auth/logout - 用户登出
			// 需要认证后才能登出
			authProtected.POST("/logout", r.authHandler.Logout)

			// POST /api/v1/auth/change-password - 修改密码
			// 需要认证，用户修改自己的密码
			authProtected.POST("/change-password", r.authHandler.ChangePassword)
		}

		// Token 刷新路由（公开，但需要有效的 refresh token）
		refreshGroup := v1.Group("/auth")
		{
			// POST /api/v1/auth/refresh - 刷新访问令牌
			// 使用 refresh token 获取新的 access token
			refreshGroup.POST("/refresh", r.authHandler.RefreshToken)
		}

		// ==================== 受保护路由 ====================
		// RBAC管理路由组(需要认证+admin权限)
		if r.rbacHandler != nil && r.jwt != nil && r.rbacService != nil {
			rbacGroup := v1.Group("/rbac")
			// 认证中间件
			rbacGroup.Use(middleware.AuthMiddleware(r.jwt))
			// 需要admin角色
			rbacGroup.Use(middleware.RequireRole(r.rbacService, "admin"))
			{
				// 角色管理
				rbacGroup.POST("/users/:id/roles", r.rbacHandler.AssignRole)
				rbacGroup.POST("/users/:id/roles/batch", r.rbacHandler.AssignRoles)
				rbacGroup.DELETE("/users/:id/roles/:role", r.rbacHandler.RevokeRole)
				rbacGroup.GET("/users/:id/roles", r.rbacHandler.GetUserRoles)
				rbacGroup.GET("/roles/:role/users", r.rbacHandler.GetRoleUsers)

				// 策略管理
				rbacGroup.POST("/policies", r.rbacHandler.AddPolicy)
				rbacGroup.POST("/policies/batch", r.rbacHandler.AddPolicies)
				rbacGroup.DELETE("/policies", r.rbacHandler.RemovePolicy)
				rbacGroup.GET("/policies", r.rbacHandler.GetPolicies)
				rbacGroup.GET("/roles/:role/policies", r.rbacHandler.GetPoliciesByRole)

				// 权限检查
				rbacGroup.POST("/check", r.rbacHandler.CheckPermission)
			}
		}

	}
}

// healthCheck 处理健康检查请求
// GET /health
// 用途:
//
//	用于容器编排平台(K8s)、负载均衡器等检查服务是否健康
//
// 响应:
//
//	总是返回 200 OK 和简单的状态信息
//
// 设计考虑:
//   - 不访问数据库或外部服务,保证快速响应
//   - 即使数据库连接失败,健康检查也应该返回 200
//     (这样可以区分"服务进程挂了"和"依赖服务挂了")
//   - 如果需要深度健康检查(包括数据库),应该另外提供 /health/deep 端点
func (r *Router) healthCheck(c *gin.Context) {
	// 返回 200 OK 和简单的状态信息
	// 使用 result.Success 保持响应格式一致
	// gin.H 是 map[string]interface{} 的简写
	c.JSON(http.StatusOK, result.Success(gin.H{
		"status":  "ok",
		"version": constants.AppVersion,
	}))
}

// Engine 返回底层的 Gin 引擎
// 这是一个访问器方法,用于特殊场景
// 使用场景:
//   - 单元测试中需要直接访问引擎
//   - 需要在外部进行高级配置
//   - 集成其他中间件或插件
//
// 注意:
//
//	大多数情况下不需要直接访问引擎
//	应该通过 Router 的方法进行操作
func (r *Router) Engine() *gin.Engine {
	return r.engine
}
