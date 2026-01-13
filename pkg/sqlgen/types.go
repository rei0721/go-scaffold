package sqlgen

import "time"

// Schema 表示数据库 Schema
// 包含数据库的所有表结构信息
type Schema struct {
	// Name 数据库名称
	Name string
	// DatabaseType 数据库类型
	DatabaseType DatabaseType
	// Tables 表列表
	Tables []*Table
	// ParsedAt 解析时间
	ParsedAt time.Time
}

// Table 表示数据库表结构
type Table struct {
	// Name 表名
	Name string
	// Comment 表注释
	Comment string
	// Columns 字段列表
	Columns []*Column
	// PrimaryKey 主键字段名列表 (支持联合主键)
	PrimaryKey []string
	// Indexes 索引列表
	Indexes []*Index
	// ForeignKeys 外键列表
	ForeignKeys []*ForeignKey
}

// Column 表示表字段
type Column struct {
	// Name 字段名
	Name string
	// DataType 数据库原生类型 (如 VARCHAR(255), INT)
	DataType string
	// GoType Go 类型 (如 string, int64)
	GoType string
	// Nullable 是否可空
	Nullable bool
	// DefaultValue 默认值
	DefaultValue *string
	// Comment 字段注释
	Comment string
	// IsPrimaryKey 是否为主键
	IsPrimaryKey bool
	// IsAutoIncrement 是否自增
	IsAutoIncrement bool
	// Length 字段长度 (用于 VARCHAR 等)
	Length int
	// Precision 精度 (用于 DECIMAL 等)
	Precision int
	// Scale 小数位数 (用于 DECIMAL 等)
	Scale int
}

// Index 表示索引
type Index struct {
	// Name 索引名称
	Name string
	// Columns 索引包含的字段名
	Columns []string
	// IsUnique 是否唯一索引
	IsUnique bool
	// IsPrimary 是否主键索引
	IsPrimary bool
}

// ForeignKey 表示外键约束
type ForeignKey struct {
	// Name 外键名称
	Name string
	// Column 本表字段名
	Column string
	// RefTable 引用表名
	RefTable string
	// RefColumn 引用字段名
	RefColumn string
	// OnDelete 删除时的行为 (CASCADE, RESTRICT, SET NULL 等)
	OnDelete string
	// OnUpdate 更新时的行为
	OnUpdate string
}

// TableFilter 表过滤配置
type TableFilter struct {
	// Include 白名单 (为空则包含所有)
	Include []string
	// Exclude 黑名单 (支持通配符 * )
	Exclude []string
}

// GenerateTarget 生成目标选项
type GenerateTarget struct {
	// Model 生成模型
	Model bool
	// DAO 生成 DAO
	DAO bool
	// Query 生成查询构建器
	Query bool
	// Migration 生成迁移脚本
	Migration bool
}

// TagOptions 标签生成选项
type TagOptions struct {
	// JSON 生成 JSON tag
	JSON bool
	// GORM 生成 GORM tag
	GORM bool
	// Validate 生成 validate tag
	Validate bool
}

// SoftDeleteOptions 软删除选项
type SoftDeleteOptions struct {
	// Enabled 是否启用
	Enabled bool
	// Field 软删除字段名
	Field string
}

// TimestampOptions 时间戳选项
type TimestampOptions struct {
	// Enabled 是否启用
	Enabled bool
	// CreatedField 创建时间字段名
	CreatedField string
	// UpdatedField 更新时间字段名
	UpdatedField string
}

// VersionOptions 乐观锁版本选项
type VersionOptions struct {
	// Enabled 是否启用
	Enabled bool
	// Field 版本字段名
	Field string
}

// ColumnInfo 用于模板渲染的字段信息
type ColumnInfo struct {
	*Column
	// FieldName Go 字段名 (驼峰)
	FieldName string
	// JSONTag JSON 标签
	JSONTag string
	// GORMTag GORM 标签
	GORMTag string
	// ValidateTag validate 标签
	ValidateTag string
}

// TableInfo 用于模板渲染的表信息
type TableInfo struct {
	*Table
	// StructName Go 结构体名
	StructName string
	// PackageName 包名
	PackageName string
	// Columns 字段信息列表 (带渲染信息)
	ColumnInfos []*ColumnInfo
	// Imports 需要导入的包
	Imports []string
}
