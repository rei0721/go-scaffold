package sqlgen

// Config 生成器配置
type Config struct {
	// 数据库配置
	DatabaseType DatabaseType `yaml:"database_type" json:"database_type"`
	DSN          string       `yaml:"dsn" json:"dsn"`

	// 输出配置
	OutputDir   string `yaml:"output_dir" json:"output_dir"`
	PackageName string `yaml:"package_name" json:"package_name"`

	// 命名规范
	TableNameRule  NamingRule `yaml:"table_name_rule" json:"table_name_rule"`
	ColumnNameRule NamingRule `yaml:"column_name_rule" json:"column_name_rule"`

	// 生成目标
	Target GenerateTarget `yaml:"target" json:"target"`

	// 表过滤
	Tables TableFilter `yaml:"tables" json:"tables"`

	// 企业特性
	SoftDelete SoftDeleteOptions `yaml:"soft_delete" json:"soft_delete"`
	Timestamp  TimestampOptions  `yaml:"timestamp" json:"timestamp"`
	Version    VersionOptions    `yaml:"version" json:"version"`

	// 标签生成
	Tags TagOptions `yaml:"tags" json:"tags"`

	// 模板配置
	TemplatePath string `yaml:"template_path" json:"template_path"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		DatabaseType:   DatabaseSQLite,
		OutputDir:      DefaultOutputDir,
		PackageName:    DefaultPackageName,
		TableNameRule:  NamingSnakeCase,
		ColumnNameRule: NamingCamelCase,
		Target: GenerateTarget{
			Model:     true,
			DAO:       false,
			Query:     false,
			Migration: false,
		},
		SoftDelete: SoftDeleteOptions{
			Enabled: false,
			Field:   DefaultSoftDeleteField,
		},
		Timestamp: TimestampOptions{
			Enabled:      true,
			CreatedField: DefaultCreatedAtField,
			UpdatedField: DefaultUpdatedAtField,
		},
		Version: VersionOptions{
			Enabled: false,
			Field:   DefaultVersionField,
		},
		Tags: TagOptions{
			JSON: true,
			GORM: true,
		},
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.OutputDir == "" {
		return &ConfigError{Field: "output_dir", Message: "output directory is required"}
	}
	if c.PackageName == "" {
		return &ConfigError{Field: "package_name", Message: "package name is required"}
	}
	if c.DatabaseType == "" {
		return &ConfigError{Field: "database_type", Message: "database type is required"}
	}

	// 验证数据库类型
	switch c.DatabaseType {
	case DatabaseMySQL, DatabasePostgres, DatabaseSQLite:
		// 有效类型
	default:
		return &ConfigError{
			Field:   "database_type",
			Message: "unsupported database type: " + string(c.DatabaseType),
		}
	}

	return nil
}

// Merge 合并配置 (other 覆盖 c 中的非零值)
func (c *Config) Merge(other *Config) *Config {
	if other == nil {
		return c
	}

	merged := *c

	if other.DatabaseType != "" {
		merged.DatabaseType = other.DatabaseType
	}
	if other.DSN != "" {
		merged.DSN = other.DSN
	}
	if other.OutputDir != "" {
		merged.OutputDir = other.OutputDir
	}
	if other.PackageName != "" {
		merged.PackageName = other.PackageName
	}
	if other.TemplatePath != "" {
		merged.TemplatePath = other.TemplatePath
	}

	return &merged
}
