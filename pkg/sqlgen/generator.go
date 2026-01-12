// Package sqlgen 提供了基于 GORM 模型生成 SQL 语句的工具
// 支持生成建表语句和基础的 CRUD 操作 SQL
package sqlgen

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Generator SQL 生成器
type Generator struct {
	dialect Dialect // 数据库方言
}

// Dialect 数据库方言接口
type Dialect interface {
	// GetDataType 根据 Go 类型和 GORM 标签获取数据库字段类型
	GetDataType(fieldType reflect.Type, gormTag string) string
	// GetCreateTableSQL 生成建表 SQL
	GetCreateTableSQL(tableName string, fields []Field) string
	// GetInsertSQL 生成插入 SQL
	GetInsertSQL(tableName string, fields []Field) string
	// GetSelectSQL 生成查询 SQL
	GetSelectSQL(tableName string, fields []Field) string
	// GetUpdateSQL 生成更新 SQL
	GetUpdateSQL(tableName string, fields []Field) string
	// GetDeleteSQL 生成删除 SQL
	GetDeleteSQL(tableName string) string
}

// Field 字段信息
type Field struct {
	Name         string // 字段名
	Type         string // 数据库类型
	IsPrimaryKey bool   // 是否主键
	IsUnique     bool   // 是否唯一
	IsIndex      bool   // 是否索引
	IsNotNull    bool   // 是否非空
	DefaultValue string // 默认值
	Size         int    // 字段大小
	Comment      string // 注释
}

// ModelInfo 模型信息
type ModelInfo struct {
	TableName string  // 表名
	Fields    []Field // 字段列表
}

// New 创建新的 SQL 生成器
func New(dialect Dialect) *Generator {
	return &Generator{
		dialect: dialect,
	}
}

// GenerateSQL 生成模型的所有 SQL 语句
func (g *Generator) GenerateSQL(model interface{}) (*SQLResult, error) {
	modelInfo, err := g.parseModel(model)
	if err != nil {
		return nil, fmt.Errorf("failed to parse model: %w", err)
	}

	result := &SQLResult{
		TableName:   modelInfo.TableName,
		CreateTable: g.dialect.GetCreateTableSQL(modelInfo.TableName, modelInfo.Fields),
		Insert:      g.dialect.GetInsertSQL(modelInfo.TableName, modelInfo.Fields),
		Select:      g.dialect.GetSelectSQL(modelInfo.TableName, modelInfo.Fields),
		Update:      g.dialect.GetUpdateSQL(modelInfo.TableName, modelInfo.Fields),
		Delete:      g.dialect.GetDeleteSQL(modelInfo.TableName),
	}

	return result, nil
}

// SQLResult SQL 生成结果
type SQLResult struct {
	TableName   string // 表名
	CreateTable string // 建表 SQL
	Insert      string // 插入 SQL
	Select      string // 查询 SQL
	Update      string // 更新 SQL
	Delete      string // 删除 SQL
}

// parseModel 解析模型结构
func (g *Generator) parseModel(model interface{}) (*ModelInfo, error) {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	if modelType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("model must be a struct, got %s", modelType.Kind())
	}

	// 获取表名
	tableName := g.getTableName(model, modelType)

	// 解析字段
	fields, err := g.parseFields(modelType)
	if err != nil {
		return nil, fmt.Errorf("failed to parse fields: %w", err)
	}

	return &ModelInfo{
		TableName: tableName,
		Fields:    fields,
	}, nil
}

// getTableName 获取表名
func (g *Generator) getTableName(model interface{}, modelType reflect.Type) string {
	// 尝试调用 TableName 方法
	if tableNamer, ok := model.(interface{ TableName() string }); ok {
		return tableNamer.TableName()
	}

	// 使用结构体名称的复数形式作为表名
	return strings.ToLower(modelType.Name()) + "s"
}

// parseFields 解析字段
func (g *Generator) parseFields(modelType reflect.Type) ([]Field, error) {
	var fields []Field

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		// 跳过非导出字段
		if !field.IsExported() {
			continue
		}

		// 处理嵌入字段
		if field.Anonymous {
			embeddedFields, err := g.parseFields(field.Type)
			if err != nil {
				return nil, fmt.Errorf("failed to parse embedded field %s: %w", field.Name, err)
			}
			fields = append(fields, embeddedFields...)
			continue
		}

		// 跳过 JSON 标签为 "-" 的字段
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		// 解析 GORM 标签
		gormTag := field.Tag.Get("gorm")
		fieldInfo := g.parseFieldInfo(field, gormTag)

		fields = append(fields, fieldInfo)
	}

	return fields, nil
}

// parseFieldInfo 解析字段信息
func (g *Generator) parseFieldInfo(field reflect.StructField, gormTag string) Field {
	fieldInfo := Field{
		Name: g.getColumnName(field),
		Type: g.dialect.GetDataType(field.Type, gormTag),
	}

	// 解析 GORM 标签
	tags := g.parseGormTag(gormTag)

	// 设置字段属性
	fieldInfo.IsPrimaryKey = g.hasTag(tags, "primaryKey")
	fieldInfo.IsUnique = g.hasTag(tags, "uniqueIndex") || g.hasTag(tags, "unique")
	fieldInfo.IsIndex = g.hasTag(tags, "index")
	fieldInfo.IsNotNull = g.hasTag(tags, "not null")

	// 获取字段大小
	if sizeTag := g.getTagValue(tags, "size"); sizeTag != "" {
		fieldInfo.Size = g.parseSize(sizeTag)
	}

	// 获取默认值
	if defaultTag := g.getTagValue(tags, "default"); defaultTag != "" {
		fieldInfo.DefaultValue = defaultTag
	}

	// 获取注释
	if commentTag := g.getTagValue(tags, "comment"); commentTag != "" {
		fieldInfo.Comment = commentTag
	}

	return fieldInfo
}

// getColumnName 获取列名
func (g *Generator) getColumnName(field reflect.StructField) string {
	// 检查 GORM 标签中的 column 设置
	gormTag := field.Tag.Get("gorm")
	tags := g.parseGormTag(gormTag)

	if columnName := g.getTagValue(tags, "column"); columnName != "" {
		return columnName
	}

	// 使用字段名的蛇形命名
	return g.toSnakeCase(field.Name)
}

// parseGormTag 解析 GORM 标签
func (g *Generator) parseGormTag(tag string) map[string]string {
	tags := make(map[string]string)

	if tag == "" {
		return tags
	}

	parts := strings.Split(tag, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, ":") {
			kv := strings.SplitN(part, ":", 2)
			tags[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		} else {
			tags[part] = ""
		}
	}

	return tags
}

// hasTag 检查是否有指定标签
func (g *Generator) hasTag(tags map[string]string, tag string) bool {
	_, exists := tags[tag]
	return exists
}

// getTagValue 获取标签值
func (g *Generator) getTagValue(tags map[string]string, tag string) string {
	return tags[tag]
}

// parseSize 解析大小
func (g *Generator) parseSize(sizeStr string) int {
	var size int
	fmt.Sscanf(sizeStr, "%d", &size)
	return size
}

// toSnakeCase 转换为蛇形命名
func (g *Generator) toSnakeCase(str string) string {
	// 处理常见的缩写
	switch str {
	case "ID":
		return "id"
	case "URL":
		return "url"
	case "HTTP":
		return "http"
	case "API":
		return "api"
	case "JSON":
		return "json"
	case "XML":
		return "xml"
	}

	// 处理以 ID 结尾的情况
	if strings.HasSuffix(str, "ID") && len(str) > 2 {
		prefix := str[:len(str)-2]
		return g.toSnakeCase(prefix) + "_id"
	}

	var result strings.Builder

	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}

	return strings.ToLower(result.String())
}

// IsGormDeletedAt 检查是否是 GORM 的 DeletedAt 字段
func (g *Generator) IsGormDeletedAt(fieldType reflect.Type) bool {
	return fieldType == reflect.TypeOf(gorm.DeletedAt{}) ||
		fieldType == reflect.TypeOf(&gorm.DeletedAt{}) ||
		(fieldType.Kind() == reflect.Ptr && fieldType.Elem() == reflect.TypeOf(time.Time{}))
}
