package rbac

import "errors"

// 预定义错误（Sentinel Errors）
// 可使用 errors.Is() 判断
var (
	// ErrEnforcerNotInitialized Enforcer未初始化
	ErrEnforcerNotInitialized = errors.New("enforcer not initialized")

	// ErrInvalidPolicy 无效的策略
	ErrInvalidPolicy = errors.New("invalid policy")

	// ErrPolicyNotFound 策略不存在
	ErrPolicyNotFound = errors.New("policy not found")

	// ErrRoleNotFound 角色不存在
	ErrRoleNotFound = errors.New("role not found")

	// ErrLoadPolicy 加载策略失败
	ErrLoadPolicy = errors.New("failed to load policy")

	// ErrSavePolicy 保存策略失败
	ErrSavePolicy = errors.New("failed to save policy")
)

// 错误消息模板常量
// 用于 fmt.Errorf() 包装错误
const (
	ErrMsgEnforceFailed      = "enforce check failed: %w"
	ErrMsgAddPolicyFailed    = "add policy failed: %w"
	ErrMsgRemovePolicyFailed = "remove policy failed: %w"
	ErrMsgAddRoleFailed      = "add role failed: %w"
	ErrMsgRemoveRoleFailed   = "remove role failed: %w"
)
