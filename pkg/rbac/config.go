package rbac

import (
	"time"

	"gorm.io/gorm"
)

// Config RBAC配置
type Config struct {
	// 数据库连接（使用Gorm Adapter）
	// 必须已经初始化的GORM DB实例
	DB *gorm.DB

	// 模型文件路径（可选）
	// 如果为空，将使用pkg/rbac目录下的内置model.conf
	ModelPath string

	// 是否启用缓存（默认true）
	// 缓存可以显著提升权限检查性能
	EnableCache bool

	// 缓存过期时间（默认30分钟）
	// 仅在EnableCache=true时生效
	CacheTTL time.Duration

	// 是否自动保存策略（默认true）
	// 设置为true时，每次策略变更都会立即持久化到数据库
	// 设置为false时，需要手动调用SavePolicy()
	AutoSave bool

	// 表名前缀（可选）
	// 用于Casbin策略表的前缀，默认为空
	TablePrefix string
}

// DefaultConfig 返回默认配置
func DefaultConfig(db *gorm.DB) *Config {
	return &Config{
		DB:          db,
		ModelPath:   "",
		EnableCache: true,
		CacheTTL:    30 * time.Minute,
		AutoSave:    true,
		TablePrefix: "",
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.DB == nil {
		return ErrInvalidPolicy
	}
	if c.CacheTTL <= 0 {
		c.CacheTTL = 30 * time.Minute
	}
	return nil
}
