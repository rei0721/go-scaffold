package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/rei0721/rei0721/pkg/cache"
	"github.com/rei0721/rei0721/pkg/executor"
	"github.com/rei0721/rei0721/pkg/rbac"
	"github.com/rei0721/rei0721/pkg/rbac/models"
	"github.com/rei0721/rei0721/pkg/rbac/repository"
	"github.com/rei0721/rei0721/types"
	"github.com/rei0721/rei0721/types/constants"
)

// rbacServiceImpl RBAC Service 实现
type rbacServiceImpl struct {
	repo     repository.RBACRepository
	cache    atomic.Value // 存储 cache.Cache
	executor atomic.Value // 存储 executor.Manager
}

// NewRBACService 创建一个新的 RBAC Service 实例
func NewRBACService(repo repository.RBACRepository) RBACService {
	return &rbacServiceImpl{repo: repo}
}

// SetExecutor 实现 ExecutorInjectable
func (s *rbacServiceImpl) SetExecutor(exec executor.Manager) {
	s.executor.Store(exec)
}

func (s *rbacServiceImpl) getExecutor() executor.Manager {
	if exec := s.executor.Load(); exec != nil {
		return exec.(executor.Manager)
	}
	return nil
}

// SetCache 实现 CacheInjectable
func (s *rbacServiceImpl) SetCache(c cache.Cache) {
	s.cache.Store(c)
}

func (s *rbacServiceImpl) getCache() cache.Cache {
	if c := s.cache.Load(); c != nil {
		return c.(cache.Cache)
	}
	return nil
}

// --- 角色管理 ---

func (s *rbacServiceImpl) CreateRole(ctx context.Context, req *types.CreateRoleRequest) (*models.Role, error) {
	// 检查角色名是否存在
	exist, _ := s.repo.GetRoleByName(ctx, req.Name)
	if exist != nil {
		return nil, fmt.Errorf("role %s already exists", req.Name)
	}

	role := &models.Role{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}
	if req.Status == 0 {
		role.Status = 1 // 默认启用
	}

	if err := s.repo.CreateRole(ctx, role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *rbacServiceImpl) GetRole(ctx context.Context, id int64) (*models.Role, error) {
	return s.repo.GetRoleByID(ctx, id)
}

func (s *rbacServiceImpl) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	return s.repo.GetRoleByName(ctx, name)
}

func (s *rbacServiceImpl) ListRoles(ctx context.Context, page, pageSize int) ([]*models.Role, int64, error) {
	return s.repo.ListRoles(ctx, page, pageSize)
}

func (s *rbacServiceImpl) UpdateRole(ctx context.Context, id int64, req *types.UpdateRoleRequest) (*models.Role, error) {
	role, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("role not found")
	}

	if req.Name != nil {
		// 检查名称唯一性
		if *req.Name != role.Name {
			exist, _ := s.repo.GetRoleByName(ctx, *req.Name)
			if exist != nil {
				return nil, fmt.Errorf("role name %s already exists", *req.Name)
			}
			role.Name = *req.Name
		}
	}
	if req.Description != nil {
		role.Description = *req.Description
	}
	if req.Status != nil {
		role.Status = *req.Status
	}

	if err := s.repo.UpdateRole(ctx, role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *rbacServiceImpl) DeleteRole(ctx context.Context, id int64) error {
	role, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return err
	}
	if role == nil {
		return fmt.Errorf("role not found")
	}
	return s.repo.DeleteRole(ctx, id)
}

// --- 权限管理 ---

func (s *rbacServiceImpl) CreatePermission(ctx context.Context, req *types.CreatePermissionRequest) (*models.Permission, error) {
	exist, _ := s.repo.GetPermissionByName(ctx, req.Name)
	if exist != nil {
		return nil, fmt.Errorf("permission %s already exists", req.Name)
	}

	perm := &models.Permission{
		Name:        req.Name,
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
		Status:      req.Status,
	}
	if req.Status == 0 {
		perm.Status = 1
	}

	if err := s.repo.CreatePermission(ctx, perm); err != nil {
		return nil, err
	}
	return perm, nil
}

func (s *rbacServiceImpl) ListPermissions(ctx context.Context, page, pageSize int) ([]*models.Permission, int64, error) {
	return s.repo.ListPermissions(ctx, page, pageSize)
}

// --- 用户角色管理 ---

func (s *rbacServiceImpl) AssignRole(ctx context.Context, userID, roleID int64) error {
	if err := s.repo.AssignRoleToUser(ctx, userID, roleID); err != nil {
		return err
	}
	s.clearUserPermissionCache(ctx, userID)
	return nil
}

func (s *rbacServiceImpl) RevokeRole(ctx context.Context, userID, roleID int64) error {
	if err := s.repo.RemoveRoleFromUser(ctx, userID, roleID); err != nil {
		return err
	}
	s.clearUserPermissionCache(ctx, userID)
	return nil
}

func (s *rbacServiceImpl) GetUserRoles(ctx context.Context, userID int64) ([]*models.Role, error) {
	return s.repo.GetUserRoles(ctx, userID)
}

// --- 角色权限管理 ---

func (s *rbacServiceImpl) AssignPermission(ctx context.Context, roleID, permID int64) error {
	if err := s.repo.AssignPermissionToRole(ctx, roleID, permID); err != nil {
		return err
	}
	// 角色权限变更，所有拥有该角色的用户的权限缓存都失效
	// 这是一个昂贵的操作，如果用户量大，需要优化（例如使用版本号）
	// 这里暂不处理，依赖 TTL 过期
	return nil
}

func (s *rbacServiceImpl) RevokePermission(ctx context.Context, roleID, permID int64) error {
	if err := s.repo.RemovePermissionFromRole(ctx, roleID, permID); err != nil {
		return err
	}
	return nil
}

func (s *rbacServiceImpl) GetRolePermissions(ctx context.Context, roleID int64) ([]*models.Permission, error) {
	return s.repo.GetRolePermissions(ctx, roleID)
}

// --- 权限检查 ---

func (s *rbacServiceImpl) CheckPermission(ctx context.Context, userID int64, resource, action string) (bool, error) {
	// 1. 尝试从缓存获取
	cacheKey := fmt.Sprintf("%s%d", rbac.CacheKeyPrefixUserPermissions, userID)
	cache := s.getCache()
	if cache != nil {
		val, err := cache.Get(ctx, cacheKey)
		if err == nil && val != "" {
			// 缓存命中
			var perms []string
			if err := json.Unmarshal([]byte(val), &perms); err == nil {
				target := fmt.Sprintf("%s:%s", resource, action)
				for _, p := range perms {
					if p == target || p == "*:*" { // 支持超级管理员权限
						return true, nil
					}
					// 也可以支持 resource:* 的通配符
					if p == fmt.Sprintf("%s:*", resource) {
						return true, nil
					}
				}
				// 缓存中没有找到，说明无权限
				return false, nil
			}
		}
	}

	// 2. 缓存未命中，查询数据库
	perms, err := s.repo.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	// 3. 构建权限集合并存入缓存
	permStrings := make([]string, 0, len(perms))
	hasPermission := false
	target := fmt.Sprintf("%s:%s", resource, action)
	wildcardResource := fmt.Sprintf("%s:*", resource)

	for _, p := range perms {
		// 格式: resource:action
		pStr := fmt.Sprintf("%s:%s", p.Resource, p.Action)
		permStrings = append(permStrings, pStr)

		if pStr == target || pStr == "*:*" || pStr == wildcardResource {
			hasPermission = true
		}
	}

	// 异步写入缓存
	if cache != nil {
		go func() {
			// 使用独立的 context 防止原 request context 取消
			exec := s.getExecutor()
			task := func() {
				data, _ := json.Marshal(permStrings)
				_ = cache.Set(context.Background(), cacheKey, string(data), rbac.CacheTTLPermission)
			}

			if exec != nil {
				_ = exec.Execute(constants.PoolCache, task)
			} else {
				task()
			}
		}()
	}

	return hasPermission, nil
}

func (s *rbacServiceImpl) clearUserPermissionCache(ctx context.Context, userID int64) {
	cache := s.getCache()
	if cache != nil {
		cacheKey := fmt.Sprintf("%s%d", rbac.CacheKeyPrefixUserPermissions, userID)
		// 异步清除
		exec := s.getExecutor()
		task := func() {
			_ = cache.Delete(context.Background(), cacheKey)
		}
		if exec != nil {
			_ = exec.Execute(constants.PoolCache, task)
		} else {
			task()
		}
	}
}
