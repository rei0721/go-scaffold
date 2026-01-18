package yaml2go

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dave/jennifer/jen"
	"gopkg.in/yaml.v3"
)

// converter YAML 转 Go 结构体转换器的实现
type converter struct {
	config *Config
	mu     sync.RWMutex
}

// newConverter 创建转换器实例（内部使用）
func newConverter(config *Config) *converter {
	c := &converter{
		config: normalizeConfig(config),
	}
	return c
}

// normalizeConfig 规范化配置，填充默认值
func normalizeConfig(config *Config) *Config {
	if config == nil {
		config = &Config{}
	}

	// 设置默认包名
	if config.PackageName == "" {
		config.PackageName = DefaultPackageName
	}

	// 设置默认结构体名
	if config.StructName == "" {
		config.StructName = DefaultStructName
	}

	// 设置默认标签
	if len(config.Tags) == 0 {
		config.Tags = copyStringSlice(DefaultTags)
	}

	// 设置默认缩进风格
	if config.IndentStyle == "" {
		config.IndentStyle = DefaultIndentStyle
	}

	return config
}

// Convert 实现 Converter.Convert
func (c *converter) Convert(yamlStr string) (string, error) {
	// 1. 验证输入
	if strings.TrimSpace(yamlStr) == "" {
		return "", ErrEmptyInput
	}

	// 2. 解析 YAML
	var data interface{}
	if err := yaml.Unmarshal([]byte(yamlStr), &data); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidYAML, err)
	}

	// 3. 获取配置
	c.mu.RLock()
	cfg := c.config
	c.mu.RUnlock()

	// 4. 构建结构体信息
	structInfo, err := c.buildStructInfo(data, cfg.StructName)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrTypeInference, err)
	}
	structInfo.PackageName = cfg.PackageName

	// 5. 生成代码
	code, err := c.generateCode(structInfo)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrCodeGeneration, err)
	}

	return code, nil
}

// ConvertToFile 实现 Converter.ConvertToFile
func (c *converter) ConvertToFile(yamlStr string, outputPath string) error {
	// 1. 转换为代码
	code, err := c.Convert(yamlStr)
	if err != nil {
		return err
	}

	// 2. 创建父目录
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("%w: %v", ErrFileWrite, err)
	}

	// 3. 写入文件
	if err := os.WriteFile(outputPath, []byte(code), 0644); err != nil {
		return fmt.Errorf("%w: %v", ErrFileWrite, err)
	}

	return nil
}

// SetConfig 实现 Converter.SetConfig
func (c *converter) SetConfig(config *Config) error {
	if config == nil {
		return ErrInvalidConfig
	}

	// 验证配置
	if config.IndentStyle != "" && config.IndentStyle != IndentStyleTab && config.IndentStyle != IndentStyleSpace {
		return fmt.Errorf("%w: invalid indent style: %s", ErrInvalidConfig, config.IndentStyle)
	}

	c.mu.Lock()
	c.config = normalizeConfig(config)
	c.mu.Unlock()

	return nil
}

// buildStructInfo 从 YAML 数据构建结构体信息
func (c *converter) buildStructInfo(data interface{}, structName string) (*StructInfo, error) {
	structInfo := &StructInfo{
		Name:    structName,
		Comment: structName + " 配置结构",
	}

	// 确保根数据是 map
	rootMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("root element must be an object, got %T", data)
	}

	// 构建字段
	for key, value := range rootMap {
		field, err := c.buildFieldInfo(key, value)
		if err != nil {
			return nil, err
		}
		structInfo.Fields = append(structInfo.Fields, field)
	}

	return structInfo, nil
}

// buildFieldInfo 从键值对构建字段信息
func (c *converter) buildFieldInfo(key string, value interface{}) (*FieldInfo, error) {
	c.mu.RLock()
	cfg := c.config
	c.mu.RUnlock()

	field := &FieldInfo{
		Name:         sanitizeFieldName(key),
		OriginalName: key,
		Tags:         make(map[string]string),
		IsPointer:    cfg.UsePointer,
	}

	// 生成标签
	for _, tagName := range cfg.Tags {
		field.Tags[tagName] = key
	}

	// 推断类型
	fieldType, elementType, children, err := c.inferType(value)
	if err != nil {
		return nil, err
	}

	field.Type = fieldType
	field.ElementType = elementType
	field.Children = children

	// 添加注释
	if cfg.AddComments {
		field.Comment = key + " 字段"
	}

	return field, nil
}

// inferType 推断值的类型
func (c *converter) inferType(value interface{}) (FieldType, *FieldInfo, []*FieldInfo, error) {
	if value == nil {
		return TypeInterface, nil, nil, nil
	}

	switch v := value.(type) {
	case string:
		return TypeString, nil, nil, nil

	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return TypeInt, nil, nil, nil

	case float32, float64:
		return TypeFloat, nil, nil, nil

	case bool:
		return TypeBool, nil, nil, nil

	case []interface{}:
		if len(v) == 0 {
			// 空数组，无法推断元素类型
			return TypeSlice, &FieldInfo{Type: TypeInterface}, nil, nil
		}

		// 推断第一个元素的类型
		elemType, elemElementType, elemChildren, err := c.inferType(v[0])
		if err != nil {
			return TypeUnknown, nil, nil, err
		}

		elementInfo := &FieldInfo{
			Type:        elemType,
			ElementType: elemElementType,
			Children:    elemChildren,
		}

		return TypeSlice, elementInfo, nil, nil

	case map[string]interface{}:
		// 嵌套对象
		var children []*FieldInfo
		for key, val := range v {
			child, err := c.buildFieldInfo(key, val)
			if err != nil {
				return TypeUnknown, nil, nil, err
			}
			children = append(children, child)
		}
		return TypeStruct, nil, children, nil

	default:
		// 未知类型，使用 interface{}
		return TypeInterface, nil, nil, nil
	}
}

// generateCode 生成 Go 代码
func (c *converter) generateCode(structInfo *StructInfo) (string, error) {
	f := jen.NewFile(structInfo.PackageName)

	// 生成根结构体
	c.generateStruct(f, structInfo.Name, structInfo.Fields, structInfo.Comment)

	// 生成嵌套结构体
	c.generateNestedStructs(f, structInfo.Fields)

	// 渲染代码
	buf := &bytes.Buffer{}
	if err := f.Render(buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// generateStruct 生成单个结构体定义
func (c *converter) generateStruct(f *jen.File, name string, fields []*FieldInfo, comment string) {
	c.mu.RLock()
	cfg := c.config
	c.mu.RUnlock()

	// 添加注释
	if comment != "" {
		f.Comment(comment)
	}

	// 构建结构体
	structFields := []jen.Code{}

	for _, field := range fields {
		// 字段注释
		var fieldCode *jen.Statement
		if cfg.AddComments && field.Comment != "" {
			fieldCode = jen.Comment(field.Comment).Line()
		} else {
			fieldCode = jen.Null()
		}

		// 字段定义
		fieldType := c.buildFieldType(field)
		tagStr := buildTags(field.Tags, cfg.OmitEmpty)

		fieldCode = fieldCode.Id(field.Name).Add(fieldType)
		if tagStr != "" {
			fieldCode = fieldCode.Tag(map[string]string{"": strings.Trim(tagStr, "`")})
		}

		structFields = append(structFields, fieldCode)
	}

	f.Type().Id(name).Struct(structFields...)
}

// buildFieldType 构建字段类型
func (c *converter) buildFieldType(field *FieldInfo) jen.Code {
	c.mu.RLock()
	cfg := c.config
	c.mu.RUnlock()

	var typeCode jen.Code

	switch field.Type {
	case TypeString:
		typeCode = jen.String()
	case TypeInt:
		typeCode = jen.Int64()
	case TypeFloat:
		typeCode = jen.Float64()
	case TypeBool:
		typeCode = jen.Bool()
	case TypeInterface:
		typeCode = jen.Interface()
	case TypeSlice:
		elemType := c.buildFieldType(field.ElementType)
		typeCode = jen.Index().Add(elemType)
	case TypeStruct:
		// 内联结构体
		structFields := []jen.Code{}
		for _, child := range field.Children {
			childType := c.buildFieldType(child)
			tagStr := buildTags(child.Tags, cfg.OmitEmpty)

			childCode := jen.Id(child.Name).Add(childType)
			if tagStr != "" {
				childCode = childCode.Tag(map[string]string{"": strings.Trim(tagStr, "`")})
			}
			structFields = append(structFields, childCode)
		}
		typeCode = jen.Struct(structFields...)
	default:
		typeCode = jen.Interface()
	}

	// 添加指针
	if cfg.UsePointer && field.Type != TypeInterface {
		typeCode = jen.Op("*").Add(typeCode)
	}

	return typeCode
}

// generateNestedStructs 生成嵌套结构体（当前实现使用内联结构体，此方法预留）
func (c *converter) generateNestedStructs(f *jen.File, fields []*FieldInfo) {
	// 当前使用内联结构体，不需要额外生成
	// 如果未来需要提取嵌套结构体为独立类型，在此实现
}
