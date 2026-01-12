package sqlgen

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
)

// MySQLDialect MySQL 方言实现
type MySQLDialect struct{}

// NewMySQLDialect 创建 MySQL 方言
func NewMySQLDialect() *MySQLDialect {
	return &MySQLDialect{}
}

// GetDataType 获取 MySQL 数据类型
func (d *MySQLDialect) GetDataType(fieldType reflect.Type, gormTag string) string {
	// 处理指针类型
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	// 检查 GORM 标签中是否指定了类型
	tags := parseGormTag(gormTag)
	if typeTag := tags["type"]; typeTag != "" {
		return typeTag
	}

	// 根据 Go 类型映射到 MySQL 类型
	switch fieldType.Kind() {
	case reflect.Bool:
		return "TINYINT(1)"
	case reflect.Int, reflect.Int32:
		return "INT"
	case reflect.Int64:
		return "BIGINT"
	case reflect.Uint, reflect.Uint32:
		return "INT UNSIGNED"
	case reflect.Uint64:
		return "BIGINT UNSIGNED"
	case reflect.Float32:
		return "FLOAT"
	case reflect.Float64:
		return "DOUBLE"
	case reflect.String:
		if size := tags["size"]; size != "" {
			sizeInt := 0
			fmt.Sscanf(size, "%d", &sizeInt)
			if sizeInt <= 255 {
				return fmt.Sprintf("VARCHAR(%s)", size)
			}
		}
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
func (d *MySQLDialect) GetCreateTableSQL(tableName string, fields []Field) string {
	var lines []string
	var indexes []string

	// 生成字段定义
	for _, field := range fields {
		line := fmt.Sprintf("    `%s` %s", field.Name, field.Type)

		// 主键自增
		if field.IsPrimaryKey {
			if field.Type == "BIGINT" {
				line += " AUTO_INCREMENT PRIMARY KEY"
			} else {
				line += " PRIMARY KEY"
			}
		}

		// 非空约束
		if field.IsNotNull && !field.IsPrimaryKey {
			line += " NOT NULL"
		}

		// 默认值
		if field.DefaultValue != "" {
			line += fmt.Sprintf(" DEFAULT %s", field.DefaultValue)
		}

		// 注释
		if field.Comment != "" {
			line += fmt.Sprintf(" COMMENT '%s'", field.Comment)
		}

		lines = append(lines, line)

		// 收集索引信息
		if field.IsUnique && !field.IsPrimaryKey {
			indexes = append(indexes, fmt.Sprintf("    UNIQUE KEY `idx_%s_%s` (`%s`)",
				tableName, field.Name, field.Name))
		} else if field.IsIndex {
			indexes = append(indexes, fmt.Sprintf("    KEY `idx_%s_%s` (`%s`)",
				tableName, field.Name, field.Name))
		}
	}

	// 添加索引到字段定义中
	if len(indexes) > 0 {
		lines = append(lines, indexes...)
	}

	// 构建建表语句
	sql := fmt.Sprintf("CREATE TABLE `%s` (\n%s\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;",
		tableName, strings.Join(lines, ",\n"))

	return sql
}

// GetInsertSQL 生成插入 SQL
func (d *MySQLDialect) GetInsertSQL(tableName string, fields []Field) string {
	var columns []string
	var placeholders []string

	for _, field := range fields {
		// 跳过自动生成的字段（如自增主键）
		if field.IsPrimaryKey && strings.Contains(field.Type, "AUTO_INCREMENT") {
			continue
		}
		columns = append(columns, fmt.Sprintf("`%s`", field.Name))
		placeholders = append(placeholders, "?")
	}

	return fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s);",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))
}

// GetSelectSQL 生成查询 SQL
func (d *MySQLDialect) GetSelectSQL(tableName string, fields []Field) string {
	var columns []string
	for _, field := range fields {
		columns = append(columns, fmt.Sprintf("`%s`", field.Name))
	}

	sql := fmt.Sprintf("-- 查询所有记录\nSELECT %s FROM `%s` WHERE `deleted_at` IS NULL;\n\n",
		strings.Join(columns, ", "), tableName)

	sql += fmt.Sprintf("-- 根据 ID 查询\nSELECT %s FROM `%s` WHERE `id` = ? AND `deleted_at` IS NULL;\n\n",
		strings.Join(columns, ", "), tableName)

	sql += fmt.Sprintf("-- 分页查询\nSELECT %s FROM `%s` WHERE `deleted_at` IS NULL ORDER BY `id` LIMIT ? OFFSET ?;",
		strings.Join(columns, ", "), tableName)

	return sql
}

// GetUpdateSQL 生成更新 SQL
func (d *MySQLDialect) GetUpdateSQL(tableName string, fields []Field) string {
	var setParts []string

	for _, field := range fields {
		// 跳过主键和时间戳字段
		if field.IsPrimaryKey || field.Name == "created_at" || field.Name == "deleted_at" {
			continue
		}
		setParts = append(setParts, fmt.Sprintf("`%s` = ?", field.Name))
	}

	sql := fmt.Sprintf("UPDATE `%s` SET %s, `updated_at` = CURRENT_TIMESTAMP WHERE `id` = ? AND `deleted_at` IS NULL;",
		tableName, strings.Join(setParts, ", "))

	return sql
}

// GetDeleteSQL 生成删除 SQL
func (d *MySQLDialect) GetDeleteSQL(tableName string) string {
	sql := fmt.Sprintf("-- 软删除\nUPDATE `%s` SET `deleted_at` = CURRENT_TIMESTAMP WHERE `id` = ? AND `deleted_at` IS NULL;\n\n", tableName)
	sql += fmt.Sprintf("-- 硬删除 (谨慎使用)\nDELETE FROM `%s` WHERE `id` = ?;", tableName)
	return sql
}
