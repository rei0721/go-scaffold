package config

// InitDBConfig 数据库初始化配置
type InitDBConfig struct {
	// ScriptDir 初始化脚本目录
	ScriptDir string `mapstructure:"script_dir"`
	// LockFile 初始化锁文件
	LockFile string `mapstructure:"lock_file"`
	// ScriptFilePrefix 初始化脚本文件名前缀
	ScriptFilePrefix string `mapstructure:"script_file_prefix"`
}

func (c *InitDBConfig) ValidateName() string {
	return AppInitDBName
}

func (c *InitDBConfig) ValidateRequired() bool {
	return true
}

// Validate 验证数据库初始化配置
// 实现 Configurable 接口
func (c *InitDBConfig) Validate() error {
	return nil
}
