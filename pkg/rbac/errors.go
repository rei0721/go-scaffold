package rbac

import "errors"

// 预定义错误（Sentinel Errors）
// 这些错误可以使用 errors.Is() 进行判断
var (
	// ErrRoleNotFound 角色不存在
	ErrRoleNotFound = errors.New("role not found")

	// ErrRoleAlreadyExists 角色已存在（名称重复）
	ErrRoleAlreadyExists = errors.New("role already exists")

	// ErrPermissionNotFound 权限不存在
	ErrPermissionNotFound = errors.New("permission not found")

	// ErrPermissionAlreadyExists 权限已存在（名称重复）
	ErrPermissionAlreadyExists = errors.New("permission already exists")

	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidPermissionFormat 权限格式无效
	// 权限格式应为 "resource:action"
	ErrInvalidPermissionFormat = errors.New("invalid permission format, expected 'resource:action'")

	// ErrEmptyRoleName 角色名称为空
	ErrEmptyRoleName = errors.New("role name cannot be empty")

	// ErrEmptyPermissionName 权限名称为空
	ErrEmptyPermissionName = errors.New("permission name cannot be empty")

	// ErrRoleDisabled 角色已禁用
	ErrRoleDisabled = errors.New("role is disabled")

	// ErrPermissionDisabled 权限已禁用
	ErrPermissionDisabled = errors.New("permission is disabled")

	// ErrPermissionDenied 权限被拒绝
	ErrPermissionDenied = errors.New("permission denied")
)

// 错误消息模板常量
// 用于 fmt.Errorf() 包装错误
const (
	// ErrMsgCreateRoleFailed 创建角色失败
	ErrMsgCreateRoleFailed = "failed to create role: %w"

	// ErrMsgUpdateRoleFailed 更新角色失败
	ErrMsgUpdateRoleFailed = "failed to update role: %w"

	// ErrMsgDeleteRoleFailed 删除角色失败
	ErrMsgDeleteRoleFailed = "failed to delete role: %w"

	// ErrMsgGetRoleFailed 获取角色失败
	ErrMsgGetRoleFailed = "failed to get role: %w"

	// ErrMsgCreatePermissionFailed 创建权限失败
	ErrMsgCreatePermissionFailed = "failed to create permission: %w"

	// ErrMsgAssignRoleFailed 分配角色失败
	ErrMsgAssignRoleFailed = "failed to assign role to user: %w"

	// ErrMsgRevokeRoleFailed 撤销角色失败
	ErrMsgRevokeRoleFailed = "failed to revoke role from user: %w"

	// ErrMsgAssignPermissionFailed 分配权限失败
	ErrMsgAssignPermissionFailed = "failed to assign permission to role: %w"

	// ErrMsgRevokePermissionFailed 撤销权限失败
	ErrMsgRevokePermissionFailed = "failed to revoke permission from role: %w"

	// ErrMsgCheckPermissionFailed 检查权限失败
	ErrMsgCheckPermissionFailed = "failed to check permission: %w"

	// ErrMsgCacheOperationFailed 缓存操作失败
	ErrMsgCacheOperationFailed = "cache operation failed: %w"
)
