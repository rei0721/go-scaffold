package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rei0721/go-scaffold/internal/service/rbac"
	"github.com/rei0721/go-scaffold/types/result"
)

// RequirePermission 权限检查中间件
// 检查当前用户是否有访问指定资源的权限
// 用法:
//
//	router.Use(middleware.RequirePermission(rbacSvc, "users", "write"))
//
// 参数:
//
//	rbacSvc: RBAC服务实例
//	resource: 资源名称
//	action: 操作名称
//
// 返回:
//
//	gin.HandlerFunc: Gin中间件处理函数
//
// 注意:
//
//	此中间件依赖于AuthMiddleware,必须在认证中间件之后使用
func RequirePermission(rbacSvc rbac.RBACService, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户ID
		userID, ok := GetUserID(c)
		if !ok {
			// 未认证或用户ID获取失败
			result.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		// 检查权限
		allowed, err := rbacSvc.CheckPermission(c.Request.Context(), userID, resource, action)
		if err != nil {
			// 权限检查失败（内部错误）
			result.InternalError(c, "Failed to check permission")
			c.Abort()
			return
		}

		if !allowed {
			// 没有权限
			result.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		// 有权限，继续处理
		c.Next()
	}
}

// RequirePermissionWithDomain 带域的权限检查中间件
// 检查当前用户在指定域中是否有访问指定资源的权限
// 用于多租户场景
// 用法:
//
//	router.Use(middleware.RequirePermissionWithDomain(rbacSvc, "tenant1", "users", "write"))
//
// 参数:
//
//	rbacSvc: RBAC服务实例
//	domain: 域名（租户ID）
//	resource: 资源名称
//	action: 操作名称
func RequirePermissionWithDomain(rbacSvc rbac.RBACService, domain, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户ID
		userID, ok := GetUserID(c)
		if !ok {
			result.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		// 检查权限
		allowed, err := rbacSvc.CheckPermissionWithDomain(c.Request.Context(), userID, domain, resource, action)
		if err != nil {
			result.InternalError(c, "Failed to check permission")
			c.Abort()
			return
		}

		if !allowed {
			result.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole 角色检查中间件
// 检查当前用户是否拥有指定角色
// 用法:
//
//	router.Use(middleware.RequireRole(rbacSvc, "admin"))
//
// 参数:
//
//	rbacSvc: RBAC服务实例
//	role: 角色名称
//
// 返回:
//
//	gin.HandlerFunc: Gin中间件处理函数
func RequireRole(rbacSvc rbac.RBACService, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户ID
		userID, ok := GetUserID(c)
		if !ok {
			result.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		// 获取用户的所有角色
		roles, err := rbacSvc.GetUserRoles(c.Request.Context(), userID)
		if err != nil {
			result.InternalError(c, "Failed to get user roles")
			c.Abort()
			return
		}

		// 检查是否拥有指定角色
		hasRole := false
		for _, r := range roles {
			if r == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			result.Forbidden(c, "Required role not found")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRoleInDomain 带域的角色检查中间件
// 检查当前用户在指定域中是否拥有指定角色
// 用于多租户场景
// 用法:
//
//	router.Use(middleware.RequireRoleInDomain(rbacSvc, "admin", "tenant1"))
//
// 参数:
//
//	rbacSvc: RBAC服务实例
//	role: 角色名称
//	domain: 域名（租户ID）
func RequireRoleInDomain(rbacSvc rbac.RBACService, role, domain string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户ID
		userID, ok := GetUserID(c)
		if !ok {
			result.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		// 获取用户在指定域中的角色
		roles, err := rbacSvc.GetUserRolesInDomain(c.Request.Context(), userID, domain)
		if err != nil {
			result.InternalError(c, "Failed to get user roles in domain")
			c.Abort()
			return
		}

		// 检查是否拥有指定角色
		hasRole := false
		for _, r := range roles {
			if r == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			result.Forbidden(c, "Required role not found in domain")
			c.Abort()
			return
		}

		c.Next()
	}
}
