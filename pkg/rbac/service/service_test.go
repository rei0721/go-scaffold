package service_test

import (
	"context"
	"testing"

	"github.com/rei0721/go-scaffold/pkg/rbac/models"
	"github.com/rei0721/go-scaffold/pkg/rbac/service"
	"github.com/rei0721/go-scaffold/types"
)

// mockRBACRepository 是用于测试的 mock Repository
type mockRBACRepository struct {
	roles       map[int64]*models.Role
	permissions map[int64]*models.Permission
	userRoles   map[int64][]int64 // userID -> roleIDs
	rolePerms   map[int64][]int64 // roleID -> permIDs
}

func newMockRepository() *mockRBACRepository {
	return &mockRBACRepository{
		roles:       make(map[int64]*models.Role),
		permissions: make(map[int64]*models.Permission),
		userRoles:   make(map[int64][]int64),
		rolePerms:   make(map[int64][]int64),
	}
}

func (m *mockRBACRepository) CreateRole(ctx context.Context, role *models.Role) error {
	role.ID = int64(len(m.roles) + 1)
	m.roles[role.ID] = role
	return nil
}

func (m *mockRBACRepository) GetRoleByID(ctx context.Context, id int64) (*models.Role, error) {
	return m.roles[id], nil
}

func (m *mockRBACRepository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	for _, role := range m.roles {
		if role.Name == name {
			return role, nil
		}
	}
	return nil, nil
}

func (m *mockRBACRepository) ListRoles(ctx context.Context, page, pageSize int) ([]*models.Role, int64, error) {
	roles := make([]*models.Role, 0, len(m.roles))
	for _, role := range m.roles {
		roles = append(roles, role)
	}
	return roles, int64(len(roles)), nil
}

func (m *mockRBACRepository) UpdateRole(ctx context.Context, role *models.Role) error {
	m.roles[role.ID] = role
	return nil
}

func (m *mockRBACRepository) DeleteRole(ctx context.Context, id int64) error {
	delete(m.roles, id)
	return nil
}

func (m *mockRBACRepository) CreatePermission(ctx context.Context, perm *models.Permission) error {
	perm.ID = int64(len(m.permissions) + 1)
	m.permissions[perm.ID] = perm
	return nil
}

func (m *mockRBACRepository) GetPermissionByID(ctx context.Context, id int64) (*models.Permission, error) {
	return m.permissions[id], nil
}

func (m *mockRBACRepository) GetPermissionByName(ctx context.Context, name string) (*models.Permission, error) {
	for _, perm := range m.permissions {
		if perm.Name == name {
			return perm, nil
		}
	}
	return nil, nil
}

func (m *mockRBACRepository) ListPermissions(ctx context.Context, page, pageSize int) ([]*models.Permission, int64, error) {
	perms := make([]*models.Permission, 0, len(m.permissions))
	for _, perm := range m.permissions {
		perms = append(perms, perm)
	}
	return perms, int64(len(perms)), nil
}

func (m *mockRBACRepository) UpdatePermission(ctx context.Context, perm *models.Permission) error {
	m.permissions[perm.ID] = perm
	return nil
}

func (m *mockRBACRepository) DeletePermission(ctx context.Context, id int64) error {
	delete(m.permissions, id)
	return nil
}

func (m *mockRBACRepository) AssignRoleToUser(ctx context.Context, userID, roleID int64) error {
	m.userRoles[userID] = append(m.userRoles[userID], roleID)
	return nil
}

func (m *mockRBACRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID int64) error {
	roles := m.userRoles[userID]
	for i, rid := range roles {
		if rid == roleID {
			m.userRoles[userID] = append(roles[:i], roles[i+1:]...)
			break
		}
	}
	return nil
}

func (m *mockRBACRepository) GetUserRoles(ctx context.Context, userID int64) ([]*models.Role, error) {
	roleIDs := m.userRoles[userID]
	roles := make([]*models.Role, 0, len(roleIDs))
	for _, rid := range roleIDs {
		if role, ok := m.roles[rid]; ok {
			roles = append(roles, role)
		}
	}
	return roles, nil
}

func (m *mockRBACRepository) GetUserPermissions(ctx context.Context, userID int64) ([]*models.Permission, error) {
	permMap := make(map[int64]bool)
	roleIDs := m.userRoles[userID]

	for _, rid := range roleIDs {
		permIDs := m.rolePerms[rid]
		for _, pid := range permIDs {
			permMap[pid] = true
		}
	}

	perms := make([]*models.Permission, 0, len(permMap))
	for pid := range permMap {
		if perm, ok := m.permissions[pid]; ok {
			perms = append(perms, perm)
		}
	}
	return perms, nil
}

func (m *mockRBACRepository) AssignPermissionToRole(ctx context.Context, roleID, permID int64) error {
	m.rolePerms[roleID] = append(m.rolePerms[roleID], permID)
	return nil
}

func (m *mockRBACRepository) RemovePermissionFromRole(ctx context.Context, roleID, permID int64) error {
	perms := m.rolePerms[roleID]
	for i, pid := range perms {
		if pid == permID {
			m.rolePerms[roleID] = append(perms[:i], perms[i+1:]...)
			break
		}
	}
	return nil
}

func (m *mockRBACRepository) GetRolePermissions(ctx context.Context, roleID int64) ([]*models.Permission, error) {
	permIDs := m.rolePerms[roleID]
	perms := make([]*models.Permission, 0, len(permIDs))
	for _, pid := range permIDs {
		if perm, ok := m.permissions[pid]; ok {
			perms = append(perms, perm)
		}
	}
	return perms, nil
}

func (m *mockRBACRepository) UserHasPermission(ctx context.Context, userID int64, resource, action string) (bool, error) {
	perms, err := m.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, perm := range perms {
		if perm.Resource == resource && perm.Action == action {
			return true, nil
		}
		if perm.Resource == "*" && perm.Action == "*" {
			return true, nil
		}
		if perm.Resource == resource && perm.Action == "*" {
			return true, nil
		}
	}
	return false, nil
}

// TestCreateRole 测试创建角色
func TestCreateRole(t *testing.T) {
	repo := newMockRepository()
	svc := service.NewRBACService(repo)

	ctx := context.Background()
	req := &types.CreateRoleRequest{
		Name:        "admin",
		Description: "Administrator",
		Status:      1,
	}

	role, err := svc.CreateRole(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if role.Name != req.Name {
		t.Errorf("expected role name %s, got %s", req.Name, role.Name)
	}
}

// TestCreatePermission 测试创建权限
func TestCreatePermission(t *testing.T) {
	repo := newMockRepository()
	svc := service.NewRBACService(repo)

	ctx := context.Background()
	req := &types.CreatePermissionRequest{
		Name:        "users:read",
		Resource:    "users",
		Action:      "read",
		Description: "Read users",
		Status:      1,
	}

	perm, err := svc.CreatePermission(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if perm.Resource != req.Resource {
		t.Errorf("expected resource %s, got %s", req.Resource, perm.Resource)
	}
}

// TestCheckPermission 测试权限检查
func TestCheckPermission(t *testing.T) {
	repo := newMockRepository()
	svc := service.NewRBACService(repo)
	ctx := context.Background()

	// 创建角色
	role, _ := svc.CreateRole(ctx, &types.CreateRoleRequest{
		Name:   "editor",
		Status: 1,
	})

	// 创建权限
	perm, _ := svc.CreatePermission(ctx, &types.CreatePermissionRequest{
		Name:     "posts:write",
		Resource: "posts",
		Action:   "write",
		Status:   1,
	})

	// 分配权限给角色
	_ = svc.AssignPermission(ctx, role.ID, perm.ID)

	// 分配角色给用户
	userID := int64(100)
	_ = svc.AssignRole(ctx, userID, role.ID)

	// 检查权限
	hasPermission, err := svc.CheckPermission(ctx, userID, "posts", "write")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !hasPermission {
		t.Error("expected user to have permission")
	}

	// 检查不存在的权限
	hasPermission, err = svc.CheckPermission(ctx, userID, "posts", "delete")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hasPermission {
		t.Error("expected user to not have permission")
	}
}
