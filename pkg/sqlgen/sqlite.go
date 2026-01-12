package sqlgen

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
)

// SQLiteDialect SQLite 方言实现
type SQLiteDialect struct{}

// NewSQLiteDialect 创建 SQLite 方言
func NewSQLiteDialect() *SQLiteDialect {
	return &SQLiteDialect{}
}

// GetDataType 获取 SQLite 数据类型
func (d *SQLiteDialect) GetDataType(fieldType reflect.Type, gormTag string) string {
	// 处理指针类型
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	// 检查 GORM 标签中是否指定了类型
	tags := parseGormTag(gormTag)
	if typeTag := tags["type"]; typeTag != "" {
		return typeTag
	}

	// 根据 Go 类型映射到 SQLite 类型
	switch fieldType.Kind() {
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Int, reflect.Int32, reflect.Int64:
		return "INTEGER"
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		return "INTEGER"
	case reflect.Float32, reflect.Float64:
		return "REAL"
	case reflect.String:
		return "TEXT"
	}

	// 处理特殊类型
	switch fieldType {
	case reflect.TypeOf(time.Time{}):
		return "DATETIME"
	case reflect.TypeOf(gorm.DeletedAt{}):
		return "DATETIME"
	}

	// 默认类型
	return "TEXT"
}

// GetCreateTableSQL 生成建表 SQL
func (d *SQLiteDialect) GetCreateTableSQL(tableName string, fields []Field) string {
	var lines []string
	var indexes []string

	// 生成字段定义
	for _, field := range fields {
		line := fmt.Sprintf("    %s %s", field.Name, field.Type)

		// 主键自增
		if field.IsPrimaryKey {
			line += " PRIMARY KEY AUTOINCREMENT"
		}

		// 非空约束
		if field.IsNotNull && !field.IsPrimaryKey {
			line += " NOT NULL"
		}

		// 默认值
		if field.DefaultValue != "" {
			line += fmt.Sprintf(" DEFAULT %s", field.DefaultValue)
		}

		lines = append(lines, line)

		// 收集索引信息
		if field.IsUnique && !field.IsPrimaryKey {
			indexes = append(indexes, fmt.Sprintf("CREATE UNIQUE INDEX idx_%s_%s ON %s (%s);",
				tableName, field.Name, tableName, field.Name))
		} else if field.IsIndex {
			indexes = append(indexes, fmt.Sprintf("CREATE INDEX idx_%s_%s ON %s (%s);",
				tableName, field.Name, tableName, field.Name))
		}
	}

	// 构建建表语句
	sql := fmt.Sprintf("CREATE TABLE %s (\n%s\n);", tableName, strings.Join(lines, ",\n"))

	// 添加索引语句
	if len(indexes) > 0 {
		sql += "\n\n" + strings.Join(indexes, "\n")
	}

	return sql
}

// GetInsertSQL 生成插入 SQL
func (d *SQLiteDialect) GetInsertSQL(tableName string, fields []Field) string {
	var columns []string
	var placeholders []string

	for _, field := range fields {
		// 跳过自动生成的字段（如自增主键）
		if field.IsPrimaryKey && field.Type == "INTEGER" {
			continue
		}
		columns = append(columns, field.Name)
		placeholders = append(placeholders, "?")
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))
}

// GetSelectSQL 生成查询 SQL
func (d *SQLiteDialect) GetSelectSQL(tableName string, fields []Field) string {
	var columns []string
	for _, field := range fields {
		columns = append(columns, field.Name)
	}

	sql := fmt.Sprintf("-- 查询所有记录\nSELECT %s FROM %s WHERE deleted_at IS NULL;\n\n",
		strings.Join(columns, ", "), tableName)

	sql += fmt.Sprintf("-- 根据 ID 查询\nSELECT %s FROM %s WHERE id = ? AND deleted_at IS NULL;\n\n",
		strings.Join(columns, ", "), tableName)

	sql += fmt.Sprintf("-- 分页查询\nSELECT %s FROM %s WHERE deleted_at IS NULL ORDER BY id LIMIT ? OFFSET ?;",
		strings.Join(columns, ", "), tableName)

	return sql
}

// GetUpdateSQL 生成更新 SQL
func (d *SQLiteDialect) GetUpdateSQL(tableName string, fields []Field) string {
	var setParts []string

	for _, field := range fields {
		// 跳过主键和时间戳字段
		if field.IsPrimaryKey || field.Name == "created_at" || field.Name == "deleted_at" {
			continue
		}
		setParts = append(setParts, fmt.Sprintf("%s = ?", field.Name))
	}

	sql := fmt.Sprintf("UPDATE %s SET %s, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL;",
		tableName, strings.Join(setParts, ", "))

	return sql
}

// GetDeleteSQL 生成删除 SQL
func (d *SQLiteDialect) GetDeleteSQL(tableName string) string {
	sql := fmt.Sprintf("-- 软删除\nUPDATE %s SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL;\n\n", tableName)
	sql += fmt.Sprintf("-- 硬删除 (谨慎使用)\nDELETE FROM %s WHERE id = ?;", tableName)
	return sql
}
