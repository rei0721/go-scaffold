package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware 返回 CORS 中间件
// 基于配置处理跨域资源共享(CORS)
// 参数:
//
//	cfg: CORS 配置,包含允许的源、方法、请求头等
//
// 返回:
//
//	gin.HandlerFunc: CORS 中间件
//
// 功能:
//  1. 处理 OPTIONS 预检请求
//  2. 设置 CORS 相关响应头
//  3. 支持配置化的源匹配(精确匹配和通配符)
//
// 使用场景:
//
//	当前端应用和后端 API 不在同一域名时,需要启用 CORS
//	例如: 前端在 http://localhost:3000, 后端在 http://localhost:8080
//
// 中间件顺序:
//
//	建议在 TraceID 之后,Logger 之前注册
//	这样可以确保预检请求也能被正确追踪和记录
func CORSMiddleware(cfg CORSConfig) gin.HandlerFunc {
	// 如果未启用 CORS,返回空中间件
	// 空中间件直接调用 c.Next(),不做任何处理
	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// 构建 gin-contrib/cors 的配置
	// 这个库提供了完整的 CORS 实现,遵循 W3C CORS 规范
	corsConfig := cors.Config{
		// AllowOrigins 允许的源列表
		// 支持:
		//   - 精确匹配: "http://localhost:3000"
		//   - 通配符: "*" (允许所有源)
		AllowOrigins: cfg.AllowOrigins,

		// AllowMethods 允许的 HTTP 方法
		// 常用方法: GET, POST, PUT, DELETE, PATCH, OPTIONS
		// OPTIONS 用于预检请求,通常需要包含
		AllowMethods: cfg.AllowMethods,

		// AllowHeaders 允许的请求头
		// 客户端可以在跨域请求中携带这些请求头
		AllowHeaders: cfg.AllowHeaders,

		// ExposeHeaders 暴露给浏览器的响应头
		// 默认情况下浏览器只能访问简单响应头
		// 通过此配置可以让浏览器访问自定义响应头
		ExposeHeaders: cfg.ExposeHeaders,

		// AllowCredentials 是否允许携带凭证
		// true: 允许跨域请求携带 Cookie、HTTP Auth 等
		// false: 不允许携带凭证
		// 安全警告: 设置为 true 时,AllowOrigins 不能使用通配符 "*"
		AllowCredentials: cfg.AllowCredentials,

		// MaxAge 预检请求缓存时间
		// 浏览器会缓存 OPTIONS 预检请求的结果
		// 在缓存有效期内,相同的跨域请求不会再发送预检请求
		// 作用: 减少网络开销,提升性能
		MaxAge: time.Duration(cfg.MaxAge) * time.Second,
	}

	// 使用 gin-contrib/cors 创建中间件
	// 这个库会:
	//   1. 检查请求的 Origin 头是否在允许列表中
	//   2. 对于 OPTIONS 预检请求,直接返回 CORS 响应头
	//   3. 对于实际请求,添加必要的 CORS 响应头后继续处理
	return cors.New(corsConfig)
}
