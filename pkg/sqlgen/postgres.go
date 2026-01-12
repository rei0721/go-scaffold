package sqlgen

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
)

// PostgresDialect PostgreSQL 方言实现
type PostgresDialect struct{}

// NewPostgresDialect 创建 PostgreSQL 方言
func NewPostgresDialect() *PostgresDialect {
	return &PostgresDialect{}
}

// GetDataType 获取 PostgreSQL 数据类型
func (d *PostgresDialect) GetDataType(fieldType reflect.Type, gormTag string) string {
	// 处理指针类型
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	// 检查 GORM 标签中是否指定了类型
	tags := parseGormTag(gormTag)
	if typeTag := tags["type"]; typeTag != "" {
		return typeTag
	}

	// 根据 Go 类型映射到 PostgreSQL 类型
	switch fieldType.Kind() {
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Int, reflect.Int32:
		return "INTEGER"
	case reflect.Int64:
		return "BIGINT"
	case reflect.Uint, reflect.Uint32:
		return "INTEGER"
	case reflect.Uint64:
		return "BIGINT"
	case reflect.Float32:
		return "REAL"
	case reflect.Float64:
		return "DOUBLE PRECISION"
	case reflect.String:
		if size := tags["size"]; size != "" {
			return fmt.Sprintf("VARCHAR(%s)", size)
		}
		return "TEXT"
	}

	// 处理特殊类型
	switch fieldType {
	case reflect.TypeOf(time.Time{}):
		return "TIMESTAMP WITH TIME ZONE"
	case reflect.TypeOf(gorm.DeletedAt{}):
		return "TIMESTAMP WITH TIME ZONE"
	}

	// 默认类型
	return "TEXT"
}

// GetCreateTableSQL 生成建表 SQL
func (d *PostgresDialect) GetCreateTableSQL(tableName string, fields []Field) string {
	var lines []string
	var indexes []string

	// 生成字段定义
	for _, field := range fields {
		line := fmt.Sprintf("    %s %s", field.Name, field.Type)

		// 主键
		if field.IsPrimaryKey {
			line += " PRIMARY KEY"
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
func (d *PostgresDialect) GetInsertSQL(tableName string, fields []Field) string {
	var columns []string
	var placeholders []string

	for i, field := range fields {
		// 跳过自动生成的字段（如自增主键）
		if field.IsPrimaryKey && field.Type == "BIGSERIAL" {
			continue
		}
		columns = append(columns, field.Name)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id;",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))
}

// GetSelectSQL 生成查询 SQL
func (d *PostgresDialect) GetSelectSQL(tableName string, fields []Field) string {
	var columns []string
	for _, field := range fields {
		columns = append(columns, field.Name)
	}

	sql := fmt.Sprintf("-- 查询所有记录\nSELECT %s FROM %s WHERE deleted_at IS NULL;\n\n",
		strings.Join(columns, ", "), tableName)

	sql += fmt.Sprintf("-- 根据 ID 查询\nSELECT %s FROM %s WHERE id = $1 AND deleted_at IS NULL;\n\n",
		strings.Join(columns, ", "), tableName)

	sql += fmt.Sprintf("-- 分页查询\nSELECT %s FROM %s WHERE deleted_at IS NULL ORDER BY id LIMIT $1 OFFSET $2;",
		strings.Join(columns, ", "), tableName)

	return sql
}

// GetUpdateSQL 生成更新 SQL
func (d *PostgresDialect) GetUpdateSQL(tableName string, fields []Field) string {
	var setParts []string
	paramIndex := 1

	for _, field := range fields {
		// 跳过主键和时间戳字段
		if field.IsPrimaryKey || field.Name == "created_at" || field.Name == "deleted_at" {
			continue
		}
		setParts = append(setParts, fmt.Sprintf("%s = $%d", field.Name, paramIndex))
		paramIndex++
	}

	sql := fmt.Sprintf("UPDATE %s SET %s, updated_at = CURRENT_TIMESTAMP WHERE id = $%d AND deleted_at IS NULL;",
		tableName, strings.Join(setParts, ", "), paramIndex)

	return sql
}

// GetDeleteSQL 生成删除 SQL
func (d *PostgresDialect) GetDeleteSQL(tableName string) string {
	sql := fmt.Sprintf("-- 软删除\nUPDATE %s SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL;\n\n", tableName)
	sql += fmt.Sprintf("-- 硬删除 (谨慎使用)\nDELETE FROM %s WHERE id = $1;", tableName)
	return sql
}

// parseGormTag 解析 GORM 标签的辅助函数
func parseGormTag(tag string) map[string]string {
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
