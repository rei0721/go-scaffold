package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rei0721/go-scaffold/pkg/logger"
	rbacservice "github.com/rei0721/go-scaffold/pkg/rbac/service"
	"github.com/rei0721/go-scaffold/types"
	"github.com/rei0721/go-scaffold/types/result"
)

type RBACHandler struct {
	service rbacservice.RBACService
	logger  logger.Logger
}

func NewRBACHandler(svc rbacservice.RBACService, log logger.Logger) *RBACHandler {
	return &RBACHandler{
		service: svc,
		logger:  log,
	}
}

// --- 角色管理 API ---

// CreateRole 创建角色
// POST /api/v1/roles
func (h *RBACHandler) CreateRole(c *gin.Context) {
	var req types.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	role, err := h.service.CreateRole(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("failed to create role", "error", err)
		result.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.OK(c, role)
}

// ListRoles 获取角色列表
// GET /api/v1/roles
func (h *RBACHandler) ListRoles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	roles, total, err := h.service.ListRoles(c.Request.Context(), page, pageSize)
	if err != nil {
		h.logger.Error("failed to list roles", "error", err)
		result.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Page(c, roles, total, page, pageSize)
}

// GetRole 获取角色详情
// GET /api/v1/roles/:id
func (h *RBACHandler) GetRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		result.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}

	role, err := h.service.GetRole(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to get role", "error", err)
		result.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	if role == nil {
		result.Fail(c, http.StatusNotFound, "role not found")
		return
	}

	result.OK(c, role)
}

// UpdateRole 更新角色
// PUT /api/v1/roles/:id
func (h *RBACHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		result.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req types.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	role, err := h.service.UpdateRole(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.Error("failed to update role", "error", err)
		result.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.OK(c, role)
}

// DeleteRole 删除角色
// DELETE /api/v1/roles/:id
func (h *RBACHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		result.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.service.DeleteRole(c.Request.Context(), id); err != nil {
		h.logger.Error("failed to delete role", "error", err)
		result.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.OK[any](c, nil)
}

// --- 权限管理 API ---

// CreatePermission 创建权限
// POST /api/v1/permissions
func (h *RBACHandler) CreatePermission(c *gin.Context) {
	var req types.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	perm, err := h.service.CreatePermission(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("failed to create permission", "error", err)
		result.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.OK(c, perm)
}

// ListPermissions 获取权限列表
// GET /api/v1/permissions
func (h *RBACHandler) ListPermissions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	perms, total, err := h.service.ListPermissions(c.Request.Context(), page, pageSize)
	if err != nil {
		h.logger.Error("failed to list permissions", "error", err)
		result.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.Page(c, perms, total, page, pageSize)
}

// --- 用户角色关联 API ---

// AssignRoleToUser 为用户分配角色
// POST /api/v1/users/:id/roles
func (h *RBACHandler) AssignRoleToUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		result.Fail(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req types.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.AssignRole(c.Request.Context(), userID, req.RoleID); err != nil {
		h.logger.Error("failed to assign role", "error", err)
		result.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.OK[any](c, nil)
}

// RevokeRoleFromUser 撤销用户角色
// DELETE /api/v1/users/:id/roles/:roleID
func (h *RBACHandler) RevokeRoleFromUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		result.Fail(c, http.StatusBadRequest, "invalid user id")
		return
	}
	roleID, err := strconv.ParseInt(c.Param("roleID"), 10, 64)
	if err != nil {
		result.Fail(c, http.StatusBadRequest, "invalid role id")
		return
	}

	if err := h.service.RevokeRole(c.Request.Context(), userID, roleID); err != nil {
		h.logger.Error("failed to revoke role", "error", err)
		result.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.OK[any](c, nil)
}

// GetUserRoles 获取用户角色
// GET /api/v1/users/:id/roles
func (h *RBACHandler) GetUserRoles(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		result.Fail(c, http.StatusBadRequest, "invalid user id")
		return
	}

	roles, err := h.service.GetUserRoles(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get user roles", "error", err)
		result.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.OK(c, roles)
}

// --- 角色权限关联 API (可选实现) ---
// POST /api/v1/roles/:id/permissions
func (h *RBACHandler) AssignPermissionToRole(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		result.Fail(c, http.StatusBadRequest, "invalid role id")
		return
	}

	// 使用 map 接收 { "permission_id": 123 }
	var req map[string]int64
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	permID, ok := req["permission_id"]
	if !ok {
		result.Fail(c, http.StatusBadRequest, "permission_id required")
		return
	}

	if err := h.service.AssignPermission(c.Request.Context(), roleID, permID); err != nil {
		h.logger.Error("failed to assign permission", "error", err)
		result.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.OK[any](c, nil)
}

func (h *RBACHandler) GetRolePermissions(c *gin.Context) {
	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		result.Fail(c, http.StatusBadRequest, "invalid role id")
		return
	}

	perms, err := h.service.GetRolePermissions(c.Request.Context(), roleID)
	if err != nil {
		h.logger.Error("failed to get role permissions", "error", err)
		result.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	result.OK(c, perms)
}
