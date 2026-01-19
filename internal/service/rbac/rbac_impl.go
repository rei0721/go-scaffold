package rbac

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/rei0721/go-scaffold/pkg/logger"
	"github.com/rei0721/go-scaffold/pkg/rbac"
	"github.com/rei0721/go-scaffold/types"
)

// rbacServiceImpl 是 RBACService 的具体实现
type rbacServiceImpl struct {
	// 延迟注入的依赖（使用 atomic.Value）
	rbac   atomic.Value // rbac.RBAC
	logger atomic.Value // logger.Logger
}

// NewRBACService 创建新的RBAC服务实例
func NewRBACService() RBACService {
	return &rbacServiceImpl{}
}

// ========== 延迟注入方法 ==========

// SetRBAC 设置RBAC管理器（延迟注入）
func (s *rbacServiceImpl) SetRBAC(r rbac.RBAC) {
	s.rbac.Store(r)
}

// SetLogger 设置日志记录器（延迟注入）
func (s *rbacServiceImpl) SetLogger(l logger.Logger) {
	s.logger.Store(l)
}

// ========== 辅助方法 ==========

// getRBAC 获取RBAC实例
func (s *rbacServiceImpl) getRBAC() rbac.RBAC {
	if r := s.rbac.Load(); r != nil {
		return r.(rbac.RBAC)
	}
	return nil
}

// getLogger 获取日志实例
func (s *rbacServiceImpl) getLogger() logger.Logger {
	if l := s.logger.Load(); l != nil {
		return l.(logger.Logger)
	}
	return nil
}

// userIDToString 将用户ID转换为字符串
// Casbin使用string作为subject
func userIDToString(userID int64) string {
	return strconv.FormatInt(userID, 10)
}

// stringToUserID 将字符串转换为用户ID
func stringToUserID(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// ========== 权限检查 ==========

// CheckPermission 检查用户权限
func (s *rbacServiceImpl) CheckPermission(ctx context.Context, userID int64, resource, action string) (bool, error) {
	r := s.getRBAC()
	if r == nil {
		return false, fmt.Errorf("RBAC not initialized")
	}

	log := s.getLogger()
	user := userIDToString(userID)

	allowed, err := r.Enforce(user, resource, action)
	if err != nil {
		if log != nil {
			log.Error("failed to check permission", "user_id", userID, "resource", resource, "action", action, "error", err)
		}
		return false, fmt.Errorf("failed to check permission: %w", err)
	}

	if log != nil {
		log.Debug("permission checked", "user_id", userID, "resource", resource, "action", action, "allowed", allowed)
	}

	return allowed, nil
}

// CheckPermissionWithDomain 检查用户在指定域中的权限
func (s *rbacServiceImpl) CheckPermissionWithDomain(ctx context.Context, userID int64, domain, resource, action string) (bool, error) {
	r := s.getRBAC()
	if r == nil {
		return false, fmt.Errorf("RBAC not initialized")
	}

	log := s.getLogger()
	user := userIDToString(userID)

	allowed, err := r.EnforceWithDomain(user, domain, resource, action)
	if err != nil {
		if log != nil {
			log.Error("failed to check permission with domain", "user_id", userID, "domain", domain, "resource", resource, "action", action, "error", err)
		}
		return false, fmt.Errorf("failed to check permission with domain: %w", err)
	}

	if log != nil {
		log.Debug("permission checked with domain", "user_id", userID, "domain", domain, "resource", resource, "action", action, "allowed", allowed)
	}

	return allowed, nil
}

// ========== 角色管理 ==========

// AssignRole 为用户分配角色
func (s *rbacServiceImpl) AssignRole(ctx context.Context, userID int64, role string) error {
	r := s.getRBAC()
	if r == nil {
		return fmt.Errorf("RBAC not initialized")
	}

	log := s.getLogger()
	user := userIDToString(userID)

	if err := r.AddRoleForUser(user, role); err != nil {
		if log != nil {
			log.Error("failed to assign role", "user_id", userID, "role", role, "error", err)
		}
		return fmt.Errorf("failed to assign role: %w", err)
	}

	if log != nil {
		log.Info("role assigned", "user_id", userID, "role", role)
	}

	return nil
}

// AssignRoleInDomain 在指定域中为用户分配角色
func (s *rbacServiceImpl) AssignRoleInDomain(ctx context.Context, userID int64, role, domain string) error {
	r := s.getRBAC()
	if r == nil {
		return fmt.Errorf("RBAC not initialized")
	}

	log := s.getLogger()
	user := userIDToString(userID)

	if err := r.AddRoleForUserInDomain(user, role, domain); err != nil {
		if log != nil {
			log.Error("failed to assign role in domain", "user_id", userID, "role", role, "domain", domain, "error", err)
		}
		return fmt.Errorf("failed to assign role in domain: %w", err)
	}

	if log != nil {
		log.Info("role assigned in domain", "user_id", userID, "role", role, "domain", domain)
	}

	return nil
}

// RevokeRole 撤销用户的角色
func (s *rbacServiceImpl) RevokeRole(ctx context.Context, userID int64, role string) error {
	r := s.getRBAC()
	if r == nil {
		return fmt.Errorf("RBAC not initialized")
	}

	log := s.getLogger()
	user := userIDToString(userID)

	if err := r.DeleteRoleForUser(user, role); err != nil {
		if log != nil {
			log.Error("failed to revoke role", "user_id", userID, "role", role, "error", err)
		}
		return fmt.Errorf("failed to revoke role: %w", err)
	}

	if log != nil {
		log.Info("role revoked", "user_id", userID, "role", role)
	}

	return nil
}

// RevokeRoleInDomain 在指定域中撤销用户的角色
func (s *rbacServiceImpl) RevokeRoleInDomain(ctx context.Context, userID int64, role, domain string) error {
	r := s.getRBAC()
	if r == nil {
		return fmt.Errorf("RBAC not initialized")
	}

	log := s.getLogger()
	user := userIDToString(userID)

	if err := r.DeleteRoleForUserInDomain(user, role, domain); err != nil {
		if log != nil {
			log.Error("failed to revoke role in domain", "user_id", userID, "role", role, "domain", domain, "error", err)
		}
		return fmt.Errorf("failed to revoke role in domain: %w", err)
	}

	if log != nil {
		log.Info("role revoked in domain", "user_id", userID, "role", role, "domain", domain)
	}

	return nil
}

// GetUserRoles 获取用户的所有角色
func (s *rbacServiceImpl) GetUserRoles(ctx context.Context, userID int64) ([]string, error) {
	r := s.getRBAC()
	if r == nil {
		return nil, fmt.Errorf("RBAC not initialized")
	}

	user := userIDToString(userID)
	roles, err := r.GetRolesForUser(user)
	if err != nil {
		log := s.getLogger()
		if log != nil {
			log.Error("failed to get user roles", "user_id", userID, "error", err)
		}
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	return roles, nil
}

// GetUserRolesInDomain 获取用户在指定域中的角色
func (s *rbacServiceImpl) GetUserRolesInDomain(ctx context.Context, userID int64, domain string) ([]string, error) {
	r := s.getRBAC()
	if r == nil {
		return nil, fmt.Errorf("RBAC not initialized")
	}

	user := userIDToString(userID)
	roles, err := r.GetRolesForUserInDomain(user, domain)
	if err != nil {
		log := s.getLogger()
		if log != nil {
			log.Error("failed to get user roles in domain", "user_id", userID, "domain", domain, "error", err)
		}
		return nil, fmt.Errorf("failed to get user roles in domain: %w", err)
	}

	return roles, nil
}

// GetRoleUsers 获取拥有指定角色的所有用户
func (s *rbacServiceImpl) GetRoleUsers(ctx context.Context, role string) ([]int64, error) {
	r := s.getRBAC()
	if r == nil {
		return nil, fmt.Errorf("RBAC not initialized")
	}

	users, err := r.GetUsersForRole(role)
	if err != nil {
		log := s.getLogger()
		if log != nil {
			log.Error("failed to get role users", "role", role, "error", err)
		}
		return nil, fmt.Errorf("failed to get role users: %w", err)
	}

	// 转换字符串数组为int64数组
	userIDs := make([]int64, 0, len(users))
	for _, u := range users {
		id, err := stringToUserID(u)
		if err != nil {
			// 跳过无效的用户ID
			log := s.getLogger()
			if log != nil {
				log.Warn("invalid user ID in role", "role", role, "user", u, "error", err)
			}
			continue
		}
		userIDs = append(userIDs, id)
	}

	return userIDs, nil
}

// ========== 策略管理 ==========

// AddPolicy 添加策略
func (s *rbacServiceImpl) AddPolicy(ctx context.Context, role, resource, action string) error {
	r := s.getRBAC()
	if r == nil {
		return fmt.Errorf("RBAC not initialized")
	}

	log := s.getLogger()

	if err := r.AddPolicy(role, resource, action); err != nil {
		if log != nil {
			log.Error("failed to add policy", "role", role, "resource", resource, "action", action, "error", err)
		}
		return fmt.Errorf("failed to add policy: %w", err)
	}

	if log != nil {
		log.Info("policy added", "role", role, "resource", resource, "action", action)
	}

	return nil
}

// AddPolicyWithDomain 添加带域的策略
func (s *rbacServiceImpl) AddPolicyWithDomain(ctx context.Context, role, domain, resource, action string) error {
	r := s.getRBAC()
	if r == nil {
		return fmt.Errorf("RBAC not initialized")
	}

	log := s.getLogger()

	if err := r.AddPolicyWithDomain(role, domain, resource, action); err != nil {
		if log != nil {
			log.Error("failed to add policy with domain", "role", role, "domain", domain, "resource", resource, "action", action, "error", err)
		}
		return fmt.Errorf("failed to add policy with domain: %w", err)
	}

	if log != nil {
		log.Info("policy added with domain", "role", role, "domain", domain, "resource", resource, "action", action)
	}

	return nil
}

// RemovePolicy 删除策略
func (s *rbacServiceImpl) RemovePolicy(ctx context.Context, role, resource, action string) error {
	r := s.getRBAC()
	if r == nil {
		return fmt.Errorf("RBAC not initialized")
	}

	log := s.getLogger()

	if err := r.RemovePolicy(role, resource, action); err != nil {
		if log != nil {
			log.Error("failed to remove policy", "role", role, "resource", resource, "action", action, "error", err)
		}
		return fmt.Errorf("failed to remove policy: %w", err)
	}

	if log != nil {
		log.Info("policy removed", "role", role, "resource", resource, "action", action)
	}

	return nil
}

// RemovePolicyWithDomain 删除带域的策略
func (s *rbacServiceImpl) RemovePolicyWithDomain(ctx context.Context, role, domain, resource, action string) error {
	r := s.getRBAC()
	if r == nil {
		return fmt.Errorf("RBAC not initialized")
	}

	log := s.getLogger()

	if err := r.RemovePolicyWithDomain(role, domain, resource, action); err != nil {
		if log != nil {
			log.Error("failed to remove policy with domain", "role", role, "domain", domain, "resource", resource, "action", action, "error", err)
		}
		return fmt.Errorf("failed to remove policy with domain: %w", err)
	}

	if log != nil {
		log.Info("policy removed with domain", "role", role, "domain", domain, "resource", resource, "action", action)
	}

	return nil
}

// GetPolicies 获取所有策略
func (s *rbacServiceImpl) GetPolicies(ctx context.Context) ([]types.RBACPolicy, error) {
	r := s.getRBAC()
	if r == nil {
		return nil, fmt.Errorf("RBAC not initialized")
	}

	policies := r.GetPolicy()
	return convertCasbinPoliciesToTypes(policies), nil
}

// GetPoliciesByRole 获取指定角色的所有策略
func (s *rbacServiceImpl) GetPoliciesByRole(ctx context.Context, role string) ([]types.RBACPolicy, error) {
	r := s.getRBAC()
	if r == nil {
		return nil, fmt.Errorf("RBAC not initialized")
	}

	// 使用GetFilteredPolicy根据角色过滤
	// fieldIndex=0 表示过滤第一个字段(subject/role)
	policies := r.GetFilteredPolicy(0, role)
	return convertCasbinPoliciesToTypes(policies), nil
}

// ========== 批量操作 ==========

// AssignRoles 批量为用户分配角色
func (s *rbacServiceImpl) AssignRoles(ctx context.Context, userID int64, roles []string) error {
	for _, role := range roles {
		if err := s.AssignRole(ctx, userID, role); err != nil {
			return err
		}
	}
	return nil
}

// AddPolicies 批量添加策略
func (s *rbacServiceImpl) AddPolicies(ctx context.Context, policies []types.RBACPolicy) error {
	r := s.getRBAC()
	if r == nil {
		return fmt.Errorf("RBAC not initialized")
	}

	log := s.getLogger()

	// 转换为Casbin策略格式
	rules := make([][]string, 0, len(policies))
	for _, p := range policies {
		if p.Domain != "" {
			// 带域的策略: [role, domain, resource, action]
			rules = append(rules, []string{p.Role, p.Domain, p.Resource, p.Action})
		} else {
			// 不带域的策略: [role, resource, action]
			rules = append(rules, []string{p.Role, p.Resource, p.Action})
		}
	}

	if err := r.AddPolicies(rules); err != nil {
		if log != nil {
			log.Error("failed to add policies", "count", len(policies), "error", err)
		}
		return fmt.Errorf("failed to add policies: %w", err)
	}

	if log != nil {
		log.Info("policies added", "count", len(policies))
	}

	return nil
}

// ========== 辅助函数 ==========

// convertCasbinPoliciesToTypes 将Casbin策略格式转换为types.RBACPolicy
func convertCasbinPoliciesToTypes(casbinPolicies [][]string) []types.RBACPolicy {
	policies := make([]types.RBACPolicy, 0, len(casbinPolicies))
	for _, p := range casbinPolicies {
		if len(p) == 3 {
			// 不带域: [role, resource, action]
			policies = append(policies, types.RBACPolicy{
				Role:     p[0],
				Resource: p[1],
				Action:   p[2],
			})
		} else if len(p) == 4 {
			// 带域: [role, domain, resource, action]
			policies = append(policies, types.RBACPolicy{
				Role:     p[0],
				Domain:   p[1],
				Resource: p[2],
				Action:   p[3],
			})
		}
	}
	return policies
}
