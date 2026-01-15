package middleware

import (
	"github.com/gin-gonic/gin"
	rbacservice "github.com/rei0721/rei0721/pkg/rbac/service"
	"github.com/rei0721/rei0721/types/result"
)

// RequirePermission 权限校验中间件
// 验证当前用户是否拥有指定资源的操作权限
// 使用方式:
//
//	router.Use(middleware.RequirePermission(rbacService, "users", "read"))
func RequirePermission(rbacService rbacservice.RBACService, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := GetUserID(c)
		if !ok {
			result.Unauthorized(c, "User not authenticated")
			c.Abort()
			return
		}

		// 调用服务层检查权限（带缓存）
		hasPermission, err := rbacService.CheckPermission(c.Request.Context(), userID, resource, action)
		if err != nil {
			// 发生错误（如数据库连接失败），安全起见由于无法确认权限，拒绝访问并记录错误
			// 这里记录日志最好
			result.Fail(c, 500, "Internal Server Error during permission check")
			c.Abort()
			return
		}

		if !hasPermission {
			result.Forbidden(c, "Permission denied")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole 角色校验中间件
// 验证当前用户是否拥有指定角色
// 注意：相比 RequirePermission，角色校验粒度较粗，建议优先使用 RequirePermission
// 使用方式:
//
//	router.Use(middleware.RequireRole(rbacService, "admin"))
func RequireRole(rbacService rbacservice.RBACService, roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := GetUserID(c)
		if !ok {
			result.Unauthorized(c, "User not authenticated")
			c.Abort()
			return
		}

		// 获取用户的所有角色
		roles, err := rbacService.GetUserRoles(c.Request.Context(), userID)
		if err != nil {
			result.Fail(c, 500, "Internal Server Error during role check")
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range roles {
			if role.Name == roleName {
				hasRole = true
				break
			}
		}

		if !hasRole {
			result.Forbidden(c, "Role denied")
			c.Abort()
			return
		}

		c.Next()
	}
}
