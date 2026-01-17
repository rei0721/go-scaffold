package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rei0721/go-scaffold/pkg/jwt"
	"github.com/rei0721/go-scaffold/types/result"
)

// AuthMiddleware JWT认证中间件
// 验证请求头中的JWT token,并将用户信息存入上下文
// 使用方式:
//
//	router.Use(middleware.AuthMiddleware(jwtManager))
//
// 或在特定路由组使用:
//
//	protected := router.Group("/api/v1")
//	protected.Use(middleware.AuthMiddleware(jwtManager))
//
// 参数:
//
//	jwtManager: JWT管理器实例
//
// 返回:
//
//	gin.HandlerFunc: Gin中间件处理函数
//
// 工作流程:
//  1. 从请求头获取 Authorization 字段
//  2. 验证 Bearer token 格式
//  3. 验证 token 有效性
//  4. 将用户信息存入上下文
//  5. 调用下一个处理器
func AuthMiddleware(jwtManager jwt.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头获取 token
		// 标准HTTP认证头格式: Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 缺少认证头,返回401未授权
			result.Unauthorized(c, "Missing authorization header")
			c.Abort()
			return
		}

		// 2. 验证 Bearer 格式
		// JWT标准要求使用 "Bearer " 前缀
		// 使用 SplitN 限制分割次数为2,防止token中包含空格导致解析错误
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			// 格式错误,返回401未授权
			result.Unauthorized(c, "Invalid authorization format")
			c.Abort()
			return
		}

		// 3. 验证 token
		// 提取token字符串（去除"Bearer "前缀）
		tokenString := parts[1]
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			// Token验证失败（无效、过期、签名错误等）
			// 统一返回401,不泄露具体原因给客户端
			result.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// 4. 将用户信息存入上下文
		// 后续的处理器可以通过GetUserID和GetUsername获取
		// 使用常量键避免拼写错误
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUsername, claims.Username)

		// 5. 调用下一个处理器
		c.Next()
	}
}

// 上下文键常量
// 定义为常量避免魔法字符串,提高可维护性
const (
	// ContextKeyUserID 用户ID在上下文中的键
	ContextKeyUserID = "user_id"

	// ContextKeyUsername 用户名在上下文中的键
	ContextKeyUsername = "username"
)

// GetUserID 从上下文获取用户ID
// 在需要获取当前用户ID的处理器中使用
// 参数:
//
//	c: Gin上下文
//
// 返回:
//
//	int64: 用户ID
//	bool: 是否成功获取
//
// 使用示例:
//
//	func MyHandler(c *gin.Context) {
//	    userID, ok := middleware.GetUserID(c)
//	    if !ok {
//	        c.JSON(500, gin.H{"error": "Failed to get user ID"})
//	        return
//	    }
//	    // 使用 userID...
//	}
func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get(ContextKeyUserID)
	if !exists {
		return 0, false
	}
	id, ok := userID.(int64)
	return id, ok
}

// GetUsername 从上下文获取用户名
// 在需要获取当前用户名的处理器中使用
// 参数:
//
//	c: Gin上下文
//
// 返回:
//
//	string: 用户名
//	bool: 是否成功获取
//
// 使用示例:
//
//	func MyHandler(c *gin.Context) {
//	    username, ok := middleware.GetUsername(c)
//	    if !ok {
//	        c.JSON(500, gin.H{"error": "Failed to get username"})
//	        return
//	    }
//	    // 使用 username...
//	}
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get(ContextKeyUsername)
	if !exists {
		return "", false
	}
	name, ok := username.(string)
	return name, ok
}
