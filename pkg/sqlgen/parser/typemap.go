package parser

import (
	"strings"
)

// TypeMapping SQL 类型到 Go 类型的映射
type TypeMapping struct {
	GoType     string // Go 类型
	NullType   string // 可空时的类型
	ImportPath string // 需要导入的包路径
}

// 通用类型映射表
var commonTypeMap = map[string]TypeMapping{
	// 整数类型
	"int":       {GoType: "int", NullType: "*int"},
	"integer":   {GoType: "int", NullType: "*int"},
	"tinyint":   {GoType: "int8", NullType: "*int8"},
	"smallint":  {GoType: "int16", NullType: "*int16"},
	"mediumint": {GoType: "int32", NullType: "*int32"},
	"bigint":    {GoType: "int64", NullType: "*int64"},

	// 无符号整数
	"tinyint unsigned":   {GoType: "uint8", NullType: "*uint8"},
	"smallint unsigned":  {GoType: "uint16", NullType: "*uint16"},
	"mediumint unsigned": {GoType: "uint32", NullType: "*uint32"},
	"int unsigned":       {GoType: "uint32", NullType: "*uint32"},
	"integer unsigned":   {GoType: "uint32", NullType: "*uint32"},
	"bigint unsigned":    {GoType: "uint64", NullType: "*uint64"},

	// 浮点类型
	"float":   {GoType: "float32", NullType: "*float32"},
	"double":  {GoType: "float64", NullType: "*float64"},
	"real":    {GoType: "float64", NullType: "*float64"},
	"decimal": {GoType: "string", NullType: "*string"}, // 精度考虑使用 string
	"numeric": {GoType: "string", NullType: "*string"},

	// 字符串类型
	"char":       {GoType: "string", NullType: "*string"},
	"varchar":    {GoType: "string", NullType: "*string"},
	"text":       {GoType: "string", NullType: "*string"},
	"tinytext":   {GoType: "string", NullType: "*string"},
	"mediumtext": {GoType: "string", NullType: "*string"},
	"longtext":   {GoType: "string", NullType: "*string"},

	// 二进制类型
	"binary":     {GoType: "[]byte", NullType: "[]byte"},
	"varbinary":  {GoType: "[]byte", NullType: "[]byte"},
	"blob":       {GoType: "[]byte", NullType: "[]byte"},
	"tinyblob":   {GoType: "[]byte", NullType: "[]byte"},
	"mediumblob": {GoType: "[]byte", NullType: "[]byte"},
	"longblob":   {GoType: "[]byte", NullType: "[]byte"},

	// 时间类型
	"date":      {GoType: "time.Time", NullType: "*time.Time", ImportPath: "time"},
	"time":      {GoType: "time.Time", NullType: "*time.Time", ImportPath: "time"},
	"datetime":  {GoType: "time.Time", NullType: "*time.Time", ImportPath: "time"},
	"timestamp": {GoType: "time.Time", NullType: "*time.Time", ImportPath: "time"},
	"year":      {GoType: "int", NullType: "*int"},

	// 布尔类型
	"bool":    {GoType: "bool", NullType: "*bool"},
	"boolean": {GoType: "bool", NullType: "*bool"},

	// JSON 类型
	"json":  {GoType: "json.RawMessage", NullType: "json.RawMessage", ImportPath: "encoding/json"},
	"jsonb": {GoType: "json.RawMessage", NullType: "json.RawMessage", ImportPath: "encoding/json"},

	// UUID 类型 (PostgreSQL)
	"uuid": {GoType: "string", NullType: "*string"},
}

// MapSQLTypeToGo 将 SQL 类型映射到 Go 类型
// 参数:
//
//	sqlType: SQL 类型 (如 VARCHAR(255), INT, BIGINT UNSIGNED)
//	nullable: 是否可空
//
// 返回:
//
//	goType: Go 类型
//	importPath: 需要导入的包路径 (空字符串表示不需要导入)
func MapSQLTypeToGo(sqlType string, nullable bool) (goType string, importPath string) {
	// 转小写并移除长度信息
	normalized := normalizeSQLType(sqlType)

	// 查找映射
	if mapping, ok := commonTypeMap[normalized]; ok {
		if nullable {
			return mapping.NullType, mapping.ImportPath
		}
		return mapping.GoType, mapping.ImportPath
	}

	// 检查是否是 TINYINT(1) (布尔类型)
	if strings.Contains(strings.ToLower(sqlType), "tinyint(1)") {
		if nullable {
			return "*bool", ""
		}
		return "bool", ""
	}

	// 默认返回 interface{}
	return "interface{}", ""
}

// normalizeSQLType 标准化 SQL 类型
// 移除长度、精度等信息，转为小写
func normalizeSQLType(sqlType string) string {
	// 转小写
	normalized := strings.ToLower(strings.TrimSpace(sqlType))

	// 移除括号内的内容 (如 VARCHAR(255) -> varchar)
	if idx := strings.Index(normalized, "("); idx > 0 {
		normalized = normalized[:idx]
	}

	// 处理 UNSIGNED
	normalized = strings.TrimSpace(normalized)

	return normalized
}

// GetImportsForType 根据 Go 类型返回需要的 import 路径
func GetImportsForType(goType string) string {
	switch {
	case strings.Contains(goType, "time.Time"):
		return "time"
	case strings.Contains(goType, "json.RawMessage"):
		return "encoding/json"
	default:
		return ""
	}
}
