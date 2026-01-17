package repository

import (
	"context"
	"errors"

	"github.com/rei0721/go-scaffold/pkg/rbac/models"
	"gorm.io/gorm"
)

// gormRBACRepository GORM 实现的 RBAC Repository
type gormRBACRepository struct {
	db *gorm.DB
}

// NewGormRBACRepository 创建一个新的 GORM RBAC Repository 实例
func NewGormRBACRepository(db *gorm.DB) RBACRepository {
	return &gormRBACRepository{db: db}
}

// --- 角色操作 ---

func (r *gormRBACRepository) CreateRole(ctx context.Context, role *models.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *gormRBACRepository) GetRoleByID(ctx context.Context, id int64) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).Preload("Permissions").First(&role, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *gormRBACRepository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).Preload("Permissions").Where("name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *gormRBACRepository) ListRoles(ctx context.Context, page, pageSize int) ([]*models.Role, int64, error) {
	var roles []*models.Role
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *gormRBACRepository) UpdateRole(ctx context.Context, role *models.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *gormRBACRepository) DeleteRole(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.Role{}, id).Error
}

// --- 权限操作 ---

func (r *gormRBACRepository) CreatePermission(ctx context.Context, perm *models.Permission) error {
	return r.db.WithContext(ctx).Create(perm).Error
}

func (r *gormRBACRepository) GetPermissionByID(ctx context.Context, id int64) (*models.Permission, error) {
	var perm models.Permission
	err := r.db.WithContext(ctx).First(&perm, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &perm, nil
}

func (r *gormRBACRepository) GetPermissionByName(ctx context.Context, name string) (*models.Permission, error) {
	var perm models.Permission
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&perm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &perm, nil
}

func (r *gormRBACRepository) ListPermissions(ctx context.Context, page, pageSize int) ([]*models.Permission, int64, error) {
	var perms []*models.Permission
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.Permission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&perms).Error; err != nil {
		return nil, 0, err
	}

	return perms, total, nil
}

func (r *gormRBACRepository) UpdatePermission(ctx context.Context, perm *models.Permission) error {
	return r.db.WithContext(ctx).Save(perm).Error
}

func (r *gormRBACRepository) DeletePermission(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.Permission{}, id).Error
}

// --- 关联操作 ---

func (r *gormRBACRepository) AssignRoleToUser(ctx context.Context, userID, roleID int64) error {
	// 使用关联表 user_roles
	// 通过 GORM Association API 添加关联
	// 注意：这里需要一个临时的 User 结构，只包含 ID 和 Roles 字段
	type User struct {
		ID    int64
		Roles []models.Role `gorm:"many2many:user_roles"`
	}
	var user User
	user.ID = userID
	var role models.Role
	role.ID = roleID
	return r.db.WithContext(ctx).Model(&user).Association("Roles").Append(&role)
}

func (r *gormRBACRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID int64) error {
	type User struct {
		ID    int64
		Roles []models.Role `gorm:"many2many:user_roles"`
	}
	var user User
	user.ID = userID
	var role models.Role
	role.ID = roleID
	return r.db.WithContext(ctx).Model(&user).Association("Roles").Delete(&role)
}

func (r *gormRBACRepository) GetUserRoles(ctx context.Context, userID int64) ([]*models.Role, error) {
	type User struct {
		ID    int64
		Roles []models.Role `gorm:"many2many:user_roles"`
	}
	var user User
	if err := r.db.WithContext(ctx).Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}
	// 转换 []models.Role 到 []*models.Role
	roles := make([]*models.Role, len(user.Roles))
	for i := range user.Roles {
		roles[i] = &user.Roles[i]
	}
	return roles, nil
}

func (r *gormRBACRepository) GetUserPermissions(ctx context.Context, userID int64) ([]*models.Permission, error) {
	// 复杂的连接查询: User -> Roles -> Permissions
	// SELECT DISTINCT p.* FROM permissions p
	// JOIN role_permissions rp ON p.id = rp.permission_id
	// JOIN user_roles ur ON rp.role_id = ur.role_id
	// WHERE ur.user_id = ? AND p.status = 1

	var perms []*models.Permission
	err := r.db.WithContext(ctx).Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND permissions.status = 1", userID).
		Distinct().
		Find(&perms).Error

	return perms, err
}

func (r *gormRBACRepository) AssignPermissionToRole(ctx context.Context, roleID, permID int64) error {
	var role models.Role
	role.ID = roleID
	var perm models.Permission
	perm.ID = permID
	return r.db.WithContext(ctx).Model(&role).Association("Permissions").Append(&perm)
}

func (r *gormRBACRepository) RemovePermissionFromRole(ctx context.Context, roleID, permID int64) error {
	var role models.Role
	role.ID = roleID
	var perm models.Permission
	perm.ID = permID
	return r.db.WithContext(ctx).Model(&role).Association("Permissions").Delete(&perm)
}

func (r *gormRBACRepository) GetRolePermissions(ctx context.Context, roleID int64) ([]*models.Permission, error) {
	var role models.Role
	if err := r.db.WithContext(ctx).Preload("Permissions").First(&role, roleID).Error; err != nil {
		return nil, err
	}
	perms := make([]*models.Permission, len(role.Permissions))
	for i := range role.Permissions {
		perms[i] = &role.Permissions[i]
	}
	return perms, nil
}

// --- 权限检查 ---

func (r *gormRBACRepository) UserHasPermission(ctx context.Context, userID int64, resource, action string) (bool, error) {
	// 检查用户是否拥有特定资源的特定权限
	// 优化：直接使用 SQL 查询是否存在匹配记录
	// SELECT count(1) FROM permissions p
	// JOIN role_permissions rp ON p.id = rp.permission_id
	// JOIN user_roles ur ON rp.role_id = ur.role_id
	// WHERE ur.user_id = ? AND p.resource = ? AND p.action = ? AND p.status = 1

	var count int64
	err := r.db.WithContext(ctx).Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND permissions.resource = ? AND permissions.action = ? AND permissions.status = 1",
			userID, resource, action).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
