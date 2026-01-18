package yaml2go

// FieldInfo 字段信息
// 描述结构体中的一个字段
type FieldInfo struct {
	// Name 字段名（Go 格式，驼峰命名）
	Name string

	// OriginalName 原始名称（YAML 中的名称）
	OriginalName string

	// Type 字段类型
	Type FieldType

	// IsPointer 是否为指针类型
	IsPointer bool

	// Comment 字段注释
	Comment string

	// Tags 标签映射 {tagName: tagValue}
	// 例如: {"json": "field_name", "yaml": "field_name"}
	Tags map[string]string

	// Children 子字段（用于嵌套结构体）
	Children []*FieldInfo

	// ElementType 数组元素类型（用于数组类型）
	ElementType *FieldInfo
}

// FieldType 字段类型枚举
type FieldType int

const (
	// TypeUnknown 未知类型
	TypeUnknown FieldType = iota

	// TypeString 字符串类型
	TypeString

	// TypeInt 整数类型
	TypeInt

	// TypeFloat 浮点数类型
	TypeFloat

	// TypeBool 布尔类型
	TypeBool

	// TypeStruct 结构体类型
	TypeStruct

	// TypeSlice 数组类型
	TypeSlice

	// TypeMap Map 类型
	TypeMap

	// TypeInterface 接口类型（用于无法推断的类型）
	TypeInterface
)

// String 返回类型的字符串表示
func (t FieldType) String() string {
	switch t {
	case TypeString:
		return "string"
	case TypeInt:
		return "int64"
	case TypeFloat:
		return "float64"
	case TypeBool:
		return "bool"
	case TypeStruct:
		return "struct"
	case TypeSlice:
		return "slice"
	case TypeMap:
		return "map[string]interface{}"
	case TypeInterface:
		return "interface{}"
	default:
		return "unknown"
	}
}

// StructInfo 结构体信息
// 描述整个结构体的元数据
type StructInfo struct {
	// Name 结构体名称
	Name string

	// PackageName 包名
	PackageName string

	// Fields 字段列表
	Fields []*FieldInfo

	// Comment 结构体注释
	Comment string
}
