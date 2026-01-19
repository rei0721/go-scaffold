// Package rbac 提供RBAC权限管理服务的实现
// 职责：
// - 权限检查（用户权限验证）
// - 角色管理（分配/撤销角色）
// - 策略管理（添加/删除/查询策略）
//
// 设计原则：
// - 封装 pkg/rbac 复杂性，提供业务友好的API
// - 支持延迟注入依赖
// - UserID 使用 int64，内部转换为 string
package rbac

import (
	"context"

	"github.com/rei0721/go-scaffold/pkg/logger"
	"github.com/rei0721/go-scaffold/pkg/rbac"
	"github.com/rei0721/go-scaffold/types"
)

// RBACService 定义RBAC服务的接口
type RBACService interface {
	// ========== 权限检查 ==========

	// CheckPermission 检查用户权限
	// 参数:
	//   ctx: 上下文
	//   userID: 用户ID
	//   resource: 资源名称
	//   action: 操作名称
	// 返回:
	//   bool: 是否有权限
	//   error: 检查过程中的错误
	CheckPermission(ctx context.Context, userID int64, resource, action string) (bool, error)

	// CheckPermissionWithDomain 检查用户在指定域中的权限
	// 用于多租户场景
	// 参数:
	//   ctx: 上下文
	//   userID: 用户ID
	//   domain: 域名（租户ID）
	//   resource: 资源名称
	//   action: 操作名称
	CheckPermissionWithDomain(ctx context.Context, userID int64, domain, resource, action string) (bool, error)

	// ========== 角色管理 ==========

	// AssignRole 为用户分配角色
	// 参数:
	//   ctx: 上下文
	//   userID: 用户ID
	//   role: 角色名称
	AssignRole(ctx context.Context, userID int64, role string) error

	// AssignRoleInDomain 在指定域中为用户分配角色
	// 参数:
	//   ctx: 上下文
	//   userID: 用户ID
	//   role: 角色名称
	//   domain: 域名
	AssignRoleInDomain(ctx context.Context, userID int64, role, domain string) error

	// RevokeRole 撤销用户的角色
	// 参数:
	//   ctx: 上下文
	//   userID: 用户ID
	//   role: 角色名称
	RevokeRole(ctx context.Context, userID int64, role string) error

	// RevokeRoleInDomain 在指定域中撤销用户的角色
	RevokeRoleInDomain(ctx context.Context, userID int64, role, domain string) error

	// GetUserRoles 获取用户的所有角色
	// 参数:
	//   ctx: 上下文
	//   userID: 用户ID
	// 返回:
	//   []string: 角色列表
	GetUserRoles(ctx context.Context, userID int64) ([]string, error)

	// GetUserRolesInDomain 获取用户在指定域中的角色
	GetUserRolesInDomain(ctx context.Context, userID int64, domain string) ([]string, error)

	// GetRoleUsers 获取拥有指定角色的所有用户
	// 参数:
	//   ctx: 上下文
	//   role: 角色名称
	// 返回:
	//   []int64: 用户ID列表
	GetRoleUsers(ctx context.Context, role string) ([]int64, error)

	// ========== 策略管理 ==========

	// AddPolicy 添加策略
	// 参数:
	//   ctx: 上下文
	//   role: 角色名称
	//   resource: 资源名称
	//   action: 操作名称
	AddPolicy(ctx context.Context, role, resource, action string) error

	// AddPolicyWithDomain 添加带域的策略
	AddPolicyWithDomain(ctx context.Context, role, domain, resource, action string) error

	// RemovePolicy 删除策略
	// 参数:
	//   ctx: 上下文
	//   role: 角色名称
	//   resource: 资源名称
	//   action: 操作名称
	RemovePolicy(ctx context.Context, role, resource, action string) error

	// RemovePolicyWithDomain 删除带域的策略
	RemovePolicyWithDomain(ctx context.Context, role, domain, resource, action string) error

	// GetPolicies 获取所有策略
	// 返回:
	//   []types.RBACPolicy: 策略列表
	GetPolicies(ctx context.Context) ([]types.RBACPolicy, error)

	// GetPoliciesByRole 获取指定角色的所有策略
	// 参数:
	//   ctx: 上下文
	//   role: 角色名称
	GetPoliciesByRole(ctx context.Context, role string) ([]types.RBACPolicy, error)

	// ========== 批量操作 ==========

	// AssignRoles 批量为用户分配角色
	// 参数:
	//   ctx: 上下文
	//   userID: 用户ID
	//   roles: 角色列表
	AssignRoles(ctx context.Context, userID int64, roles []string) error

	// AddPolicies 批量添加策略
	// 参数:
	//   ctx: 上下文
	//   policies: 策略列表
	AddPolicies(ctx context.Context, policies []types.RBACPolicy) error

	// ========== 延迟注入方法 ==========

	// SetRBAC 设置RBAC管理器（延迟注入）
	SetRBAC(r rbac.RBAC)

	// SetLogger 设置日志记录器（延迟注入）
	SetLogger(l logger.Logger)
}
