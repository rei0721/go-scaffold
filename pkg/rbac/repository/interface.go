// Package repository 定义 RBAC 数据访问层接口
package repository

import (
	"context"

	"github.com/rei0721/rei0721/pkg/rbac/models"
)

// RBACRepository RBAC 数据访问接口
// 统一管理角色、权限以及它们之间的关联关系
type RBACRepository interface {
	// 角色操作
	CreateRole(ctx context.Context, role *models.Role) error
	GetRoleByID(ctx context.Context, id int64) (*models.Role, error)
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
	ListRoles(ctx context.Context, page, pageSize int) ([]*models.Role, int64, error)
	UpdateRole(ctx context.Context, role *models.Role) error
	DeleteRole(ctx context.Context, id int64) error

	// 权限操作
	CreatePermission(ctx context.Context, perm *models.Permission) error
	GetPermissionByID(ctx context.Context, id int64) (*models.Permission, error)
	GetPermissionByName(ctx context.Context, name string) (*models.Permission, error)
	ListPermissions(ctx context.Context, page, pageSize int) ([]*models.Permission, int64, error)
	UpdatePermission(ctx context.Context, perm *models.Permission) error
	DeletePermission(ctx context.Context, id int64) error

	// 关联操作
	AssignRoleToUser(ctx context.Context, userID, roleID int64) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID int64) error
	GetUserRoles(ctx context.Context, userID int64) ([]*models.Role, error)
	GetUserPermissions(ctx context.Context, userID int64) ([]*models.Permission, error)

	AssignPermissionToRole(ctx context.Context, roleID, permID int64) error
	RemovePermissionFromRole(ctx context.Context, roleID, permID int64) error
	GetRolePermissions(ctx context.Context, roleID int64) ([]*models.Permission, error)

	// 权限检查
	UserHasPermission(ctx context.Context, userID int64, resource, action string) (bool, error)
}
