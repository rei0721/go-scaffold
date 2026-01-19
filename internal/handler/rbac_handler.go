package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/rei0721/go-scaffold/internal/middleware"
	"github.com/rei0721/go-scaffold/internal/service/rbac"
	"github.com/rei0721/go-scaffold/pkg/logger"
	"github.com/rei0721/go-scaffold/types"
	"github.com/rei0721/go-scaffold/types/result"
)

// RBACHandler RBAC管理处理器
type RBACHandler struct {
	rbacService rbac.RBACService
	logger      logger.Logger
}

// NewRBACHandler 创建新的RBAC处理器
func NewRBACHandler(rbacService rbac.RBACService, logger logger.Logger) *RBACHandler {
	return &RBACHandler{
		rbacService: rbacService,
		logger:      logger,
	}
}

// ========== 角色管理接口 ==========

// AssignRole 为用户分配角色
// POST /rbac/users/:id/roles
// Body: {"role": "admin", "domain": "tenant1"}
func (h *RBACHandler) AssignRole(c *gin.Context) {
	// 获取用户ID参数
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		result.BadRequest(c, "Invalid user ID")
		return
	}

	// 解析请求体
	var req types.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, "Invalid request body")
		return
	}

	// 分配角色
	if req.Domain != "" {
		err = h.rbacService.AssignRoleInDomain(c.Request.Context(), userID, req.Role, req.Domain)
	} else {
		err = h.rbacService.AssignRole(c.Request.Context(), userID, req.Role)
	}

	if err != nil {
		h.logger.Error("failed to assign role", "user_id", userID, "role", req.Role, "error", err)
		result.InternalError(c, "Failed to assign role")
		return
	}

	c.JSON(http.StatusOK, result.Success(gin.H{
		"message": "Role assigned successfully",
	}))
}

// RevokeRole 撤销用户的角色
// DELETE /rbac/users/:id/roles/:role
// Query: ?domain=tenant1
func (h *RBACHandler) RevokeRole(c *gin.Context) {
	// 获取用户ID参数
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		result.BadRequest(c, "Invalid user ID")
		return
	}

	// 获取角色参数
	role := c.Param("role")
	if role == "" {
		result.BadRequest(c, "Role is required")
		return
	}

	// 获取域参数（可选）
	domain := c.Query("domain")

	// 撤销角色
	if domain != "" {
		err = h.rbacService.RevokeRoleInDomain(c.Request.Context(), userID, role, domain)
	} else {
		err = h.rbacService.RevokeRole(c.Request.Context(), userID, role)
	}

	if err != nil {
		h.logger.Error("failed to revoke role", "user_id", userID, "role", role, "error", err)
		result.InternalError(c, "Failed to revoke role")
		return
	}

	c.JSON(http.StatusOK, result.Success(gin.H{
		"message": "Role revoked successfully",
	}))
}

// GetUserRoles 获取用户的所有角色
// GET /rbac/users/:id/roles
// Query: ?domain=tenant1
func (h *RBACHandler) GetUserRoles(c *gin.Context) {
	// 获取用户ID参数
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		result.BadRequest(c, "Invalid user ID")
		return
	}

	// 获取域参数（可选）
	domain := c.Query("domain")

	// 获取角色
	var roles []string
	if domain != "" {
		roles, err = h.rbacService.GetUserRolesInDomain(c.Request.Context(), userID, domain)
	} else {
		roles, err = h.rbacService.GetUserRoles(c.Request.Context(), userID)
	}

	if err != nil {
		h.logger.Error("failed to get user roles", "user_id", userID, "error", err)
		result.InternalError(c, "Failed to get user roles")
		return
	}

	c.JSON(http.StatusOK, result.Success(types.UserRolesResponse{
		UserID: userID,
		Roles:  roles,
	}))
}

// GetRoleUsers 获取拥有指定角色的所有用户
// GET /rbac/roles/:role/users
func (h *RBACHandler) GetRoleUsers(c *gin.Context) {
	// 获取角色参数
	role := c.Param("role")
	if role == "" {
		result.BadRequest(c, "Role is required")
		return
	}

	// 获取用户列表
	userIDs, err := h.rbacService.GetRoleUsers(c.Request.Context(), role)
	if err != nil {
		h.logger.Error("failed to get role users", "role", role, "error", err)
		result.InternalError(c, "Failed to get role users")
		return
	}

	c.JSON(http.StatusOK, result.Success(types.RoleUsersResponse{
		Role:    role,
		UserIDs: userIDs,
	}))
}

// ========== 策略管理接口 ==========

// AddPolicy 添加策略
// POST /rbac/policies
// Body: {"role": "admin", "resource": "users", "action": "write", "domain": "tenant1"}
func (h *RBACHandler) AddPolicy(c *gin.Context) {
	// 解析请求体
	var req types.AddPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, "Invalid request body")
		return
	}

	// 添加策略
	var err error
	if req.Domain != "" {
		err = h.rbacService.AddPolicyWithDomain(c.Request.Context(), req.Role, req.Domain, req.Resource, req.Action)
	} else {
		err = h.rbacService.AddPolicy(c.Request.Context(), req.Role, req.Resource, req.Action)
	}

	if err != nil {
		h.logger.Error("failed to add policy", "policy", req, "error", err)
		result.InternalError(c, "Failed to add policy")
		return
	}

	c.JSON(http.StatusOK, result.Success(gin.H{
		"message": "Policy added successfully",
	}))
}

// RemovePolicy 删除策略
// DELETE /rbac/policies
// Body: {"role": "admin", "resource": "users", "action": "write", "domain": "tenant1"}
func (h *RBACHandler) RemovePolicy(c *gin.Context) {
	// 解析请求体
	var req types.RemovePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, "Invalid request body")
		return
	}

	// 删除策略
	var err error
	if req.Domain != "" {
		err = h.rbacService.RemovePolicyWithDomain(c.Request.Context(), req.Role, req.Domain, req.Resource, req.Action)
	} else {
		err = h.rbacService.RemovePolicy(c.Request.Context(), req.Role, req.Resource, req.Action)
	}

	if err != nil {
		h.logger.Error("failed to remove policy", "policy", req, "error", err)
		result.InternalError(c, "Failed to remove policy")
		return
	}

	c.JSON(http.StatusOK, result.Success(gin.H{
		"message": "Policy removed successfully",
	}))
}

// GetPolicies 获取所有策略
// GET /rbac/policies
func (h *RBACHandler) GetPolicies(c *gin.Context) {
	// 获取策略列表
	policies, err := h.rbacService.GetPolicies(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to get policies", "error", err)
		result.InternalError(c, "Failed to get policies")
		return
	}

	c.JSON(http.StatusOK, result.Success(types.PoliciesResponse{
		Policies: policies,
		Total:    len(policies),
	}))
}

// GetPoliciesByRole 获取指定角色的所有策略
// GET /rbac/roles/:role/policies
func (h *RBACHandler) GetPoliciesByRole(c *gin.Context) {
	// 获取角色参数
	role := c.Param("role")
	if role == "" {
		result.BadRequest(c, "Role is required")
		return
	}

	// 获取策略列表
	policies, err := h.rbacService.GetPoliciesByRole(c.Request.Context(), role)
	if err != nil {
		h.logger.Error("failed to get policies by role", "role", role, "error", err)
		result.InternalError(c, "Failed to get policies by role")
		return
	}

	c.JSON(http.StatusOK, result.Success(types.PoliciesResponse{
		Policies: policies,
		Total:    len(policies),
	}))
}

// ========== 权限检查接口 ==========

// CheckPermission 检查权限
// POST /rbac/check
// Body: {"user_id": 123, "resource": "users", "action": "write", "domain": "tenant1"}
func (h *RBACHandler) CheckPermission(c *gin.Context) {
	// 解析请求体
	var req types.CheckPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, "Invalid request body")
		return
	}

	// 检查权限
	var allowed bool
	var err error
	if req.Domain != "" {
		allowed, err = h.rbacService.CheckPermissionWithDomain(c.Request.Context(), req.UserID, req.Domain, req.Resource, req.Action)
	} else {
		allowed, err = h.rbacService.CheckPermission(c.Request.Context(), req.UserID, req.Resource, req.Action)
	}

	if err != nil {
		h.logger.Error("failed to check permission", "request", req, "error", err)
		result.InternalError(c, "Failed to check permission")
		return
	}

	c.JSON(http.StatusOK, result.Success(types.CheckPermissionResponse{
		Allowed: allowed,
	}))
}

// AssignRoles 批量为用户分配角色
// POST /rbac/users/:id/roles/batch
// Body: {"roles": ["admin", "editor"], "domain": "tenant1"}
func (h *RBACHandler) AssignRoles(c *gin.Context) {
	// 获取用户ID参数
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		result.BadRequest(c, "Invalid user ID")
		return
	}

	// 解析请求体
	var req types.AssignRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, "Invalid request body")
		return
	}

	// 批量分配角色
	err = h.rbacService.AssignRoles(c.Request.Context(), userID, req.Roles)
	if err != nil {
		h.logger.Error("failed to assign roles", "user_id", userID, "roles", req.Roles, "error", err)
		result.InternalError(c, "Failed to assign roles")
		return
	}

	c.JSON(http.StatusOK, result.Success(gin.H{
		"message": "Roles assigned successfully",
	}))
}

// AddPolicies 批量添加策略
// POST /rbac/policies/batch
// Body: {"policies": [{"role": "admin", "resource": "users", "action": "write"}]}
func (h *RBACHandler) AddPolicies(c *gin.Context) {
	// 解析请求体
	var req types.AddPoliciesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.BadRequest(c, "Invalid request body")
		return
	}

	// 批量添加策略
	err := h.rbacService.AddPolicies(c.Request.Context(), req.Policies)
	if err != nil {
		h.logger.Error("failed to add policies", "count", len(req.Policies), "error", err)
		result.InternalError(c, "Failed to add policies")
		return
	}

	c.JSON(http.StatusOK, result.Success(gin.H{
		"message": "Policies added successfully",
	}))
}

// GetCurrentUserID 从上下文获取当前用户ID
// 用于需要操作当前用户权限的场景
func (h *RBACHandler) GetCurrentUserID(c *gin.Context) (int64, bool) {
	return middleware.GetUserID(c)
}
