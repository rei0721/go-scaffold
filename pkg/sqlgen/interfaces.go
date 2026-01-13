package sqlgen

import (
	"context"
	"database/sql"
)

// Generator SQL 代码生成器接口
// 提供解析数据库 Schema 和生成代码的能力
type Generator interface {
	// Parse 解析数据库 Schema
	// 参数:
	//   ctx: 上下文
	//   db: 数据库连接
	// 返回:
	//   *Schema: 解析后的 Schema
	//   error: 解析失败时的错误
	Parse(ctx context.Context, db *sql.DB) (*Schema, error)

	// Generate 生成代码到指定目录
	// 参数:
	//   ctx: 上下文
	//   schema: 解析后的 Schema
	//   outputDir: 输出目录
	// 返回:
	//   error: 生成失败时的错误
	Generate(ctx context.Context, schema *Schema, outputDir string) error

	// GenerateTable 生成单个表的代码
	// 参数:
	//   ctx: 上下文
	//   table: 表结构
	//   outputDir: 输出目录
	// 返回:
	//   error: 生成失败时的错误
	GenerateTable(ctx context.Context, table *Table, outputDir string) error
}

// Parser 数据库 Schema 解析器接口
// 用于从数据库读取表结构信息
type Parser interface {
	// ParseDatabase 解析整个数据库
	// 参数:
	//   ctx: 上下文
	//   db: 数据库连接
	// 返回:
	//   *Schema: 解析后的 Schema
	//   error: 解析失败时的错误
	ParseDatabase(ctx context.Context, db *sql.DB) (*Schema, error)

	// ParseTable 解析单个表
	// 参数:
	//   ctx: 上下文
	//   db: 数据库连接
	//   tableName: 表名
	// 返回:
	//   *Table: 解析后的表结构
	//   error: 解析失败时的错误
	ParseTable(ctx context.Context, db *sql.DB, tableName string) (*Table, error)

	// GetDatabaseType 获取数据库类型
	// 返回:
	//   DatabaseType: 数据库类型
	GetDatabaseType() DatabaseType
}

// TemplateEngine 模板引擎接口
// 用于渲染代码模板
type TemplateEngine interface {
	// Render 渲染模板
	// 参数:
	//   tmplName: 模板名称
	//   data: 模板数据
	// 返回:
	//   string: 渲染后的内容
	//   error: 渲染失败时的错误
	Render(tmplName string, data interface{}) (string, error)

	// RegisterFunc 注册自定义模板函数
	// 参数:
	//   name: 函数名
	//   fn: 函数实现
	// 返回:
	//   error: 注册失败时的错误
	RegisterFunc(name string, fn interface{}) error

	// LoadTemplate 加载模板
	// 参数:
	//   name: 模板名称
	//   content: 模板内容
	// 返回:
	//   error: 加载失败时的错误
	LoadTemplate(name string, content string) error

	// LoadTemplateFile 从文件加载模板
	// 参数:
	//   name: 模板名称
	//   path: 文件路径
	// 返回:
	//   error: 加载失败时的错误
	LoadTemplateFile(name string, path string) error
}

// FileWriter 文件写入器接口
// 提供原子性文件写入能力
type FileWriter interface {
	// Write 写入文件
	// 参数:
	//   path: 文件路径
	//   content: 文件内容
	// 返回:
	//   error: 写入失败时的错误
	Write(path string, content []byte) error

	// WriteAtomic 原子性写入文件 (临时文件 + 重命名)
	// 参数:
	//   path: 文件路径
	//   content: 文件内容
	// 返回:
	//   error: 写入失败时的错误
	WriteAtomic(path string, content []byte) error
}
