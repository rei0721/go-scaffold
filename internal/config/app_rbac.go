package config

import (
	"time"
)

// RBACConfig RBAC配置
type RBACConfig struct {
	// 模型文件路径（可选）
	// 如果为空，将使用pkg/rbac目录下的内置model.conf
	ModelPath string `mapstructure:"model_path"`

	// 是否启用缓存（默认true）
	// 缓存可以显著提升权限检查性能
	EnableCache bool `mapstructure:"enable_cache"`

	// 缓存过期时间（默认30分钟）
	// 仅在EnableCache=true时生效
	CacheTTL time.Duration `mapstructure:"cache_ttl"`

	// 是否自动保存策略（默认true）
	// 设置为true时，每次策略变更都会立即持久化到数据库
	// 设置为false时，需要手动调用SavePolicy()
	AutoSave bool `mapstructure:"auto_save"`

	// 表名前缀（可选）
	// 用于Casbin策略表的前缀，默认为空
	TablePrefix string `mapstructure:"table_prefix"`
}

func (c *RBACConfig) ValidateName() string {
	return AppRBACName
}

func (c *RBACConfig) ValidateRequired() bool {
	return false
}

// Validate 验证 RBAC 配置
func (c *RBACConfig) Validate() error {
	return nil
}
