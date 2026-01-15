// Package service 定义 RBAC 业务逻辑层接口
package service

import (
	"context"

	"github.com/rei0721/rei0721/pkg/rbac/models"
	"github.com/rei0721/rei0721/types"
)

// RBACService RBAC 业务逻辑接口
// 提供高层的 RBAC 操作，包括缓存集成和权限检查
type RBACService interface {
	// ExecutorInjectable 支持延迟注入 Executor
	types.ExecutorInjectable

	// CacheInjectable 支持延迟注入 Cache
	types.CacheInjectable

	// 角色管理
	CreateRole(ctx context.Context, req *types.CreateRoleRequest) (*models.Role, error)
	GetRole(ctx context.Context, id int64) (*models.Role, error)
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
	ListRoles(ctx context.Context, page, pageSize int) ([]*models.Role, int64, error)
	UpdateRole(ctx context.Context, id int64, req *types.UpdateRoleRequest) (*models.Role, error)
	DeleteRole(ctx context.Context, id int64) error

	// 权限管理
	CreatePermission(ctx context.Context, req *types.CreatePermissionRequest) (*models.Permission, error)
	ListPermissions(ctx context.Context, page, pageSize int) ([]*models.Permission, int64, error)

	// 用户角色管理
	AssignRole(ctx context.Context, userID, roleID int64) error
	RevokeRole(ctx context.Context, userID, roleID int64) error
	GetUserRoles(ctx context.Context, userID int64) ([]*models.Role, error)

	// 角色权限管理
	AssignPermission(ctx context.Context, roleID, permID int64) error
	RevokePermission(ctx context.Context, roleID, permID int64) error
	GetRolePermissions(ctx context.Context, roleID int64) ([]*models.Permission, error)

	// 权限检查
	CheckPermission(ctx context.Context, userID int64, resource, action string) (bool, error)
}
