package rbac

import (
	"embed"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

//go:embed model.conf
var modelFS embed.FS

// rbacImpl Casbin RBAC实现
type rbacImpl struct {
	enforcer *casbin.Enforcer
	config   *Config
	cache    sync.Map // 权限检查结果缓存 map[string]cacheEntry
	mu       sync.RWMutex
}

// cacheEntry 缓存条目
type cacheEntry struct {
	result    bool
	expiresAt time.Time
}

// New 创建新的RBAC实例
//
// 参数:
//
//	cfg: RBAC配置
//
// 返回:
//
//	RBAC: RBAC实例
//	error: 创建失败时的错误
//
// 示例:
//
//	rbac, err := rbac.New(&rbac.Config{
//	    DB: db,
//	    EnableCache: true,
//	})
func New(cfg *Config) (RBAC, error) {
	if cfg == nil {
		cfg = &Config{}
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// 创建Gorm Adapter
	adapter, err := gormadapter.NewAdapterByDBWithCustomTable(cfg.DB, cfg.TablePrefix, DefaultTableName)
	if err != nil {
		return nil, fmt.Errorf("failed to create gorm adapter: %w", err)
	}

	// 加载模型
	var m model.Model
	if cfg.ModelPath != "" {
		// 使用外部模型文件
		m, err = model.NewModelFromFile(cfg.ModelPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load model from file: %w", err)
		}
	} else {
		// 使用内置模型
		modelContent, err := modelFS.ReadFile("model.conf")
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded model: %w", err)
		}
		m, err = model.NewModelFromString(string(modelContent))
		if err != nil {
			return nil, fmt.Errorf("failed to load embedded model: %w", err)
		}
	}

	// 创建Enforcer
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create enforcer: %w", err)
	}

	// 设置自动保存
	enforcer.EnableAutoSave(cfg.AutoSave)

	// 加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrLoadPolicy, err)
	}

	return &rbacImpl{
		enforcer: enforcer,
		config:   cfg,
	}, nil
}

// ========== 权限检查 ==========

// Enforce 检查权限（无域）
func (r *rbacImpl) Enforce(sub, obj, act string) (bool, error) {
	return r.EnforceWithDomain(sub, "", obj, act)
}

// EnforceWithDomain 带域的权限检查
func (r *rbacImpl) EnforceWithDomain(sub, dom, obj, act string) (bool, error) {
	if r.enforcer == nil {
		return false, ErrEnforcerNotInitialized
	}

	// 检查缓存
	if r.config.EnableCache {
		if result, ok := r.getCached(sub, dom, obj, act); ok {
			return result, nil
		}
	}

	// 执行权限检查
	result, err := r.enforcer.Enforce(sub, dom, obj, act)
	if err != nil {
		return false, fmt.Errorf(ErrMsgEnforceFailed, err)
	}

	// 缓存结果
	if r.config.EnableCache {
		r.setCache(sub, dom, obj, act, result)
	}

	return result, nil
}

// ========== 角色管理 ==========

// AddRoleForUser 为用户分配角色（无域）
func (r *rbacImpl) AddRoleForUser(user, role string) error {
	return r.AddRoleForUserInDomain(user, role, "")
}

// AddRoleForUserInDomain 在指定域中为用户分配角色
func (r *rbacImpl) AddRoleForUserInDomain(user, role, domain string) error {
	if r.enforcer == nil {
		return ErrEnforcerNotInitialized
	}

	_, err := r.enforcer.AddRoleForUser(user, role, domain)
	if err != nil {
		return fmt.Errorf(ErrMsgAddRoleFailed, err)
	}

	// 清除缓存
	if r.config.EnableCache {
		r.clearUserCache(user)
	}

	return nil
}

// DeleteRoleForUser 撤销用户的角色（无域）
func (r *rbacImpl) DeleteRoleForUser(user, role string) error {
	return r.DeleteRoleForUserInDomain(user, role, "")
}

// DeleteRoleForUserInDomain 在指定域中撤销用户的角色
func (r *rbacImpl) DeleteRoleForUserInDomain(user, role, domain string) error {
	if r.enforcer == nil {
		return ErrEnforcerNotInitialized
	}

	_, err := r.enforcer.DeleteRoleForUser(user, role, domain)
	if err != nil {
		return fmt.Errorf(ErrMsgRemoveRoleFailed, err)
	}

	// 清除缓存
	if r.config.EnableCache {
		r.clearUserCache(user)
	}

	return nil
}

// GetRolesForUser 获取用户的所有角色（无域）
func (r *rbacImpl) GetRolesForUser(user string) ([]string, error) {
	return r.GetRolesForUserInDomain(user, "")
}

// GetRolesForUserInDomain 获取用户在指定域中的角色
func (r *rbacImpl) GetRolesForUserInDomain(user, domain string) ([]string, error) {
	if r.enforcer == nil {
		return nil, ErrEnforcerNotInitialized
	}

	roles, err := r.enforcer.GetRolesForUser(user, domain)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// GetUsersForRole 获取拥有指定角色的所有用户
func (r *rbacImpl) GetUsersForRole(role string) ([]string, error) {
	if r.enforcer == nil {
		return nil, ErrEnforcerNotInitialized
	}

	users, err := r.enforcer.GetUsersForRole(role, "")
	if err != nil {
		return nil, err
	}

	return users, nil
}

// ========== 策略管理 ==========

// AddPolicy 添加策略（无域）
func (r *rbacImpl) AddPolicy(sub, obj, act string) error {
	return r.AddPolicyWithDomain(sub, "", obj, act)
}

// AddPolicyWithDomain 添加带域的策略
func (r *rbacImpl) AddPolicyWithDomain(sub, domain, obj, act string) error {
	if r.enforcer == nil {
		return ErrEnforcerNotInitialized
	}

	_, err := r.enforcer.AddPolicy(sub, domain, obj, act)
	if err != nil {
		return fmt.Errorf(ErrMsgAddPolicyFailed, err)
	}

	// 清除缓存
	if r.config.EnableCache {
		r.ClearCache()
	}

	return nil
}

// RemovePolicy 删除策略（无域）
func (r *rbacImpl) RemovePolicy(sub, obj, act string) error {
	return r.RemovePolicyWithDomain(sub, "", obj, act)
}

// RemovePolicyWithDomain 删除带域的策略
func (r *rbacImpl) RemovePolicyWithDomain(sub, domain, obj, act string) error {
	if r.enforcer == nil {
		return ErrEnforcerNotInitialized
	}

	_, err := r.enforcer.RemovePolicy(sub, domain, obj, act)
	if err != nil {
		return fmt.Errorf(ErrMsgRemovePolicyFailed, err)
	}

	// 清除缓存
	if r.config.EnableCache {
		r.ClearCache()
	}

	return nil
}

// GetPolicy 获取所有策略
func (r *rbacImpl) GetPolicy() [][]string {
	if r.enforcer == nil {
		return nil
	}
	policies, _ := r.enforcer.GetPolicy()
	return policies
}

// GetFilteredPolicy 获取过滤后的策略
func (r *rbacImpl) GetFilteredPolicy(fieldIndex int, fieldValues ...string) [][]string {
	if r.enforcer == nil {
		return nil
	}
	policies, _ := r.enforcer.GetFilteredPolicy(fieldIndex, fieldValues...)
	return policies
}

// ========== 批量操作 ==========

// AddPolicies 批量添加策略
func (r *rbacImpl) AddPolicies(rules [][]string) error {
	if r.enforcer == nil {
		return ErrEnforcerNotInitialized
	}

	_, err := r.enforcer.AddPolicies(rules)
	if err != nil {
		return fmt.Errorf(ErrMsgAddPolicyFailed, err)
	}

	// 清除缓存
	if r.config.EnableCache {
		r.ClearCache()
	}

	return nil
}

// RemovePolicies 批量删除策略
func (r *rbacImpl) RemovePolicies(rules [][]string) error {
	if r.enforcer == nil {
		return ErrEnforcerNotInitialized
	}

	_, err := r.enforcer.RemovePolicies(rules)
	if err != nil {
		return fmt.Errorf(ErrMsgRemovePolicyFailed, err)
	}

	// 清除缓存
	if r.config.EnableCache {
		r.ClearCache()
	}

	return nil
}

// ========== 工具方法 ==========

// LoadPolicy 从存储加载策略
func (r *rbacImpl) LoadPolicy() error {
	if r.enforcer == nil {
		return ErrEnforcerNotInitialized
	}

	if err := r.enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("%w: %v", ErrLoadPolicy, err)
	}

	// 清除缓存
	if r.config.EnableCache {
		r.ClearCache()
	}

	return nil
}

// SavePolicy 保存策略到存储
func (r *rbacImpl) SavePolicy() error {
	if r.enforcer == nil {
		return ErrEnforcerNotInitialized
	}

	if err := r.enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("%w: %v", ErrSavePolicy, err)
	}

	return nil
}

// ClearCache 清除所有缓存
func (r *rbacImpl) ClearCache() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cache = sync.Map{}
	return nil
}

// Close 关闭RBAC实例
func (r *rbacImpl) Close() error {
	// Casbin enforcer 没有Close方法，只需清理资源
	r.cache = sync.Map{}
	r.enforcer = nil
	return nil
}

// ========== 内部辅助方法 ==========

// getCached 从缓存获取权限检查结果
func (r *rbacImpl) getCached(sub, dom, obj, act string) (bool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := r.cacheKey(sub, dom, obj, act)
	if val, ok := r.cache.Load(key); ok {
		entry := val.(cacheEntry)
		if time.Now().Before(entry.expiresAt) {
			return entry.result, true
		}
		// 缓存已过期，删除
		r.cache.Delete(key)
	}
	return false, false
}

// setCache 设置缓存
func (r *rbacImpl) setCache(sub, dom, obj, act string, result bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := r.cacheKey(sub, dom, obj, act)
	r.cache.Store(key, cacheEntry{
		result:    result,
		expiresAt: time.Now().Add(r.config.CacheTTL),
	})
}

// clearUserCache 清除用户相关的缓存
func (r *rbacImpl) clearUserCache(user string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 遍历删除所有包含该用户的缓存
	r.cache.Range(func(key, value interface{}) bool {
		keyStr := key.(string)
		if len(keyStr) > len(user) && keyStr[:len(user)] == user {
			r.cache.Delete(key)
		}
		return true
	})
}

// cacheKey 生成缓存键
func (r *rbacImpl) cacheKey(sub, dom, obj, act string) string {
	if dom == "" {
		return fmt.Sprintf("%s:%s:%s", sub, obj, act)
	}
	return fmt.Sprintf("%s:%s:%s:%s", sub, dom, obj, act)
}

// GetModelPath 获取模型文件路径（用于测试）
func GetModelPath() string {
	return filepath.Join("model.conf")
}
