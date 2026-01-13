package sqlgen

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
)

// generator SQL 代码生成器实现
type generator struct {
	config   *Config
	parser   Parser
	template TemplateEngine
	writer   FileWriter
}

// ParserFactory 解析器工厂函数类型
type ParserFactory func(dbType DatabaseType) (Parser, error)

// TemplateFactory 模板引擎工厂函数类型
type TemplateFactory func() TemplateEngine

// DefaultParserFactory 默认解析器工厂 (需要在外部设置)
var DefaultParserFactory ParserFactory

// DefaultTemplateFactory 默认模板引擎工厂 (需要在外部设置)
var DefaultTemplateFactory TemplateFactory

// NewGenerator 创建代码生成器
func NewGenerator(config *Config) (Generator, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	// 使用工厂创建解析器
	if DefaultParserFactory == nil {
		return nil, fmt.Errorf("parser factory not initialized")
	}
	p, err := DefaultParserFactory(config.DatabaseType)
	if err != nil {
		return nil, err
	}

	// 使用工厂创建模板引擎
	var tmpl TemplateEngine
	if DefaultTemplateFactory != nil {
		tmpl = DefaultTemplateFactory()
	}

	return &generator{
		config:   config,
		parser:   p,
		template: tmpl,
		writer:   NewFileWriter(),
	}, nil
}

// NewGeneratorWithDeps 创建代码生成器 (带依赖注入)
func NewGeneratorWithDeps(config *Config, p Parser, t TemplateEngine, w FileWriter) Generator {
	return &generator{
		config:   config,
		parser:   p,
		template: t,
		writer:   w,
	}
}

// Parse 解析数据库 Schema
func (g *generator) Parse(ctx context.Context, db *sql.DB) (*Schema, error) {
	schema, err := g.parser.ParseDatabase(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseSchema, err)
	}

	// 应用表过滤
	schema.Tables = g.filterTables(schema.Tables)

	return schema, nil
}

// Generate 生成代码到指定目录
func (g *generator) Generate(ctx context.Context, schema *Schema, outputDir string) error {
	if outputDir == "" {
		outputDir = g.config.OutputDir
	}

	for _, table := range schema.Tables {
		if err := g.GenerateTable(ctx, table, outputDir); err != nil {
			return err
		}
	}

	// 生成 DDL/Migration 脚本
	if g.config.Target.Migration {
		ddlGen := NewDDLGenerator(g.config)
		if err := ddlGen.GenerateDDLToFile(schema, ""); err != nil {
			return err
		}
	}

	return nil
}

// GenerateTable 生成单个表的代码
func (g *generator) GenerateTable(ctx context.Context, table *Table, outputDir string) error {
	// 准备模板数据
	data := g.prepareTemplateData(table)

	// 生成模型
	if g.config.Target.Model {
		if err := g.generateModel(data, outputDir); err != nil {
			return err
		}
	}

	// 生成 DAO
	if g.config.Target.DAO {
		if err := g.generateDAO(data, outputDir); err != nil {
			return err
		}
	}

	return nil
}

// filterTables 过滤表
func (g *generator) filterTables(tables []*Table) []*Table {
	if len(g.config.Tables.Include) == 0 && len(g.config.Tables.Exclude) == 0 {
		return tables
	}

	var filtered []*Table
	for _, t := range tables {
		// 检查白名单
		if len(g.config.Tables.Include) > 0 {
			if !matchAny(t.Name, g.config.Tables.Include) {
				continue
			}
		}

		// 检查黑名单
		if matchAny(t.Name, g.config.Tables.Exclude) {
			continue
		}

		filtered = append(filtered, t)
	}

	return filtered
}

// matchAny 检查名称是否匹配任意模式
func matchAny(name string, patterns []string) bool {
	for _, pattern := range patterns {
		if matchPattern(name, pattern) {
			return true
		}
	}
	return false
}

// matchPattern 简单的通配符匹配
func matchPattern(name, pattern string) bool {
	if pattern == "*" {
		return true
	}
	if strings.HasPrefix(pattern, "*") && strings.HasSuffix(pattern, "*") {
		return strings.Contains(name, pattern[1:len(pattern)-1])
	}
	if strings.HasPrefix(pattern, "*") {
		return strings.HasSuffix(name, pattern[1:])
	}
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(name, pattern[:len(pattern)-1])
	}
	return name == pattern
}

// prepareTemplateData 准备模板数据
func (g *generator) prepareTemplateData(table *Table) *TemplateData {
	structName := ToPascalCase(table.Name)
	structName = ToSingular(structName)

	// 收集需要的 imports
	imports := make(map[string]bool)

	// 准备字段信息
	var columnInfos []*ColumnInfo
	var insertColumns []*ColumnInfo
	var updateColumns []*ColumnInfo

	for _, col := range table.Columns {
		fieldName := ToPascalCase(col.Name)

		info := &ColumnInfo{
			Column:    col,
			FieldName: fieldName,
			JSONTag:   g.buildJSONTag(col),
			GORMTag:   g.buildGORMTag(col),
		}

		// 收集 import
		if imp := getImportsForType(col.GoType); imp != "" {
			imports[imp] = true
		}

		columnInfos = append(columnInfos, info)

		// 排除自增主键的插入列
		if !col.IsAutoIncrement {
			insertColumns = append(insertColumns, info)
		}

		// 排除主键的更新列
		if !col.IsPrimaryKey {
			updateColumns = append(updateColumns, info)
		}
	}

	// 转换 imports map 为 slice
	var importList []string
	for imp := range imports {
		importList = append(importList, imp)
	}

	return &TemplateData{
		Header:        GeneratedFileHeader,
		PackageName:   g.config.PackageName,
		StructName:    structName,
		Table:         table,
		ColumnInfos:   columnInfos,
		InsertColumns: insertColumns,
		UpdateColumns: updateColumns,
		Imports:       importList,
		Tags:          g.config.Tags,
		SoftDelete:    g.config.SoftDelete,
		Timestamp:     g.config.Timestamp,
		Version:       g.config.Version,
	}
}

// TemplateData 模板数据
type TemplateData struct {
	Header        string
	PackageName   string
	StructName    string
	Table         *Table
	ColumnInfos   []*ColumnInfo
	InsertColumns []*ColumnInfo
	UpdateColumns []*ColumnInfo
	Imports       []string
	Tags          TagOptions
	SoftDelete    SoftDeleteOptions
	Timestamp     TimestampOptions
	Version       VersionOptions
}

// buildJSONTag 构建 JSON tag
func (g *generator) buildJSONTag(col *Column) string {
	name := ToSnakeCase(col.Name)
	if strings.EqualFold(col.Name, "password") {
		return "-"
	}
	if col.Nullable {
		return name + ",omitempty"
	}
	return name
}

// buildGORMTag 构建 GORM tag
func (g *generator) buildGORMTag(col *Column) string {
	var parts []string

	if col.IsPrimaryKey {
		parts = append(parts, "primaryKey")
	}
	if col.IsAutoIncrement {
		parts = append(parts, "autoIncrement")
	}
	if !col.Nullable && !col.IsPrimaryKey {
		parts = append(parts, "not null")
	}

	// 添加类型信息
	if col.DataType != "" {
		parts = append(parts, "type:"+col.DataType)
	}

	return strings.Join(parts, ";")
}

// generateModel 生成模型文件
func (g *generator) generateModel(data *TemplateData, outputDir string) error {
	content, err := g.template.Render("model", data)
	if err != nil {
		return &GenerateError{Table: data.Table.Name, Message: "failed to render model", Cause: err}
	}

	// 输出路径: outputDir/models/table_name.go
	fileName := ToSnakeCase(data.Table.Name) + ".go"
	filePath := filepath.Join(outputDir, "models", fileName)

	if err := g.writer.WriteAtomic(filePath, []byte(content)); err != nil {
		return err
	}

	return nil
}

// generateDAO 生成 DAO 文件
func (g *generator) generateDAO(data *TemplateData, outputDir string) error {
	content, err := g.template.Render("dao", data)
	if err != nil {
		return &GenerateError{Table: data.Table.Name, Message: "failed to render dao", Cause: err}
	}

	// 输出路径: outputDir/dao/table_name_dao.go
	fileName := ToSnakeCase(data.Table.Name) + "_dao.go"
	filePath := filepath.Join(outputDir, "dao", fileName)

	if err := g.writer.WriteAtomic(filePath, []byte(content)); err != nil {
		return err
	}

	return nil
}

// getImportsForType 根据 Go 类型返回需要的 import 路径
func getImportsForType(goType string) string {
	switch {
	case strings.Contains(goType, "time.Time"):
		return "time"
	case strings.Contains(goType, "json.RawMessage"):
		return "encoding/json"
	default:
		return ""
	}
}
