package cli

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// Parser 提供命令行参数解析功能
// 封装 flag.FlagSet,提供更友好的 API
// 设计目标:
// - 简化参数定义
// - 提供类型安全的访问
// - 友好的错误消息
type Parser struct {
	// flagSet 底层的 flag.FlagSet
	flagSet *flag.FlagSet

	// flags 存储所有定义的标志
	flags map[string]*Flag

	// validators 存储标志的验证函数
	validators map[string]ValidatorFunc

	// parsed 标记是否已解析
	parsed bool

	// args 解析后的位置参数
	args []string
}

// NewParser 创建一个新的 Parser
// 参数:
//
//	name: 程序名或命令名,用于错误消息
//
// 返回:
//
//	*Parser: 新创建的解析器
//
// 使用示例:
//
//	parser := NewParser("myapp")
//	parser.String("output", "o", "./output", "输出目录")
//	parser.Parse(os.Args[1:])
func NewParser(name string) *Parser {
	return &Parser{
		flagSet:    flag.NewFlagSet(name, flag.ContinueOnError),
		flags:      make(map[string]*Flag),
		validators: make(map[string]ValidatorFunc),
		parsed:     false,
	}
}

// String 定义一个字符串标志
// 参数:
//
//	name: 标志名称
//	shortName: 短标志名称(可以为空)
//	defaultValue: 默认值
//	usage: 使用说明
//
// 返回:
//
//	*string: 指向标志值的指针
//
// 使用示例:
//
//	output := parser.String("output", "o", "./output", "输出目录")
func (p *Parser) String(name, shortName, defaultValue, usage string) *string {
	value := p.flagSet.String(name, defaultValue, usage)
	if shortName != "" {
		p.flagSet.StringVar(value, shortName, defaultValue, usage)
	}

	p.flags[name] = &Flag{
		Name:         name,
		ShortName:    shortName,
		Usage:        usage,
		DefaultValue: defaultValue,
		Type:         "string",
	}

	return value
}

// Int 定义一个整数标志
// 参数:
//
//	name: 标志名称
//	shortName: 短标志名称(可以为空)
//	defaultValue: 默认值
//	usage: 使用说明
//
// 返回:
//
//	*int: 指向标志值的指针
//
// 使用示例:
//
//	port := parser.Int("port", "p", 8080, "监听端口")
func (p *Parser) Int(name, shortName string, defaultValue int, usage string) *int {
	value := p.flagSet.Int(name, defaultValue, usage)
	if shortName != "" {
		p.flagSet.IntVar(value, shortName, defaultValue, usage)
	}

	p.flags[name] = &Flag{
		Name:         name,
		ShortName:    shortName,
		Usage:        usage,
		DefaultValue: defaultValue,
		Type:         "int",
	}

	return value
}

// Bool 定义一个布尔标志
// 参数:
//
//	name: 标志名称
//	shortName: 短标志名称(可以为空)
//	defaultValue: 默认值
//	usage: 使用说明
//
// 返回:
//
//	*bool: 指向标志值的指针
//
// 使用示例:
//
//	verbose := parser.Bool("verbose", "v", false, "详细输出")
func (p *Parser) Bool(name, shortName string, defaultValue bool, usage string) *bool {
	value := p.flagSet.Bool(name, defaultValue, usage)
	if shortName != "" {
		p.flagSet.BoolVar(value, shortName, defaultValue, usage)
	}

	p.flags[name] = &Flag{
		Name:         name,
		ShortName:    shortName,
		Usage:        usage,
		DefaultValue: defaultValue,
		Type:         "bool",
	}

	return value
}

// StringSlice 定义一个字符串切片标志
// 允许多次指定同一标志来收集多个值
// 参数:
//
//	name: 标志名称
//	shortName: 短标志名称(可以为空)
//	defaultValue: 默认值(逗号分隔的字符串)
//	usage: 使用说明
//
// 返回:
//
//	*[]string: 指向标志值的指针
//
// 使用示例:
//
//	tags := parser.StringSlice("tag", "t", "", "标签(可多次指定)")
//	// 使用: -tag foo -tag bar 或 -tag foo,bar
func (p *Parser) StringSlice(name, shortName, defaultValue, usage string) *[]string {
	var result []string
	if defaultValue != "" {
		result = strings.Split(defaultValue, ",")
	}

	valueStr := p.flagSet.String(name, defaultValue, usage)
	if shortName != "" {
		p.flagSet.StringVar(valueStr, shortName, defaultValue, usage)
	}

	p.flags[name] = &Flag{
		Name:         name,
		ShortName:    shortName,
		Usage:        usage,
		DefaultValue: defaultValue,
		Type:         "[]string",
	}

	// 返回一个指针,指向解析后会被填充的切片
	resultPtr := &result
	// 注册一个钩子,在解析后转换字符串为切片
	p.validators[name] = func(value interface{}) error {
		if str, ok := value.(string); ok && str != "" {
			*resultPtr = strings.Split(str, ",")
		}
		return nil
	}

	return resultPtr
}

// SetValidator 为标志设置验证函数
// 在解析后自动调用验证函数
// 参数:
//
//	name: 标志名称
//	validator: 验证函数
//
// 使用示例:
//
//	port := parser.Int("port", "p", 8080, "端口")
//	parser.SetValidator("port", func(v interface{}) error {
//	  if port := v.(int); port < 1 || port > 65535 {
//	    return fmt.Errorf("端口必须在 1-65535 之间")
//	  }
//	  return nil
//	})
func (p *Parser) SetValidator(name string, validator ValidatorFunc) {
	p.validators[name] = validator
}

// SetRequired 标记标志为必需
// 如果用户未提供该标志,Parse 会返回错误
// 参数:
//
//	name: 标志名称
//
// 使用示例:
//
//	input := parser.String("input", "i", "", "输入文件")
//	parser.SetRequired("input")
func (p *Parser) SetRequired(name string) {
	if f, ok := p.flags[name]; ok {
		f.Required = true
	}
}

// Parse 解析命令行参数
// 参数:
//
//	args: 命令行参数(通常是 os.Args[1:])
//
// 返回:
//
//	error: 解析失败或验证失败时的错误
//
// 使用示例:
//
//	if err := parser.Parse(os.Args[1:]); err != nil {
//	  fmt.Fprintf(os.Stderr, "Error: %v\n", err)
//	  os.Exit(ExitInvalidUsage)
//	}
func (p *Parser) Parse(args []string) error {
	// 解析标志
	if err := p.flagSet.Parse(args); err != nil {
		return fmt.Errorf(ErrMsgInvalidFlag, err)
	}

	// 保存位置参数
	p.args = p.flagSet.Args()
	p.parsed = true

	// 验证必需标志
	for name, f := range p.flags {
		if f.Required {
			// 检查标志是否被设置(不是默认值)
			if !p.isFlagSet(name) {
				return fmt.Errorf(ErrMsgMissingRequired, name)
			}
		}
	}

	// 运行验证器
	for name, validator := range p.validators {
		if value := p.getFlagValue(name); value != nil {
			if err := validator(value); err != nil {
				return fmt.Errorf(ErrMsgInvalidValue, name, err)
			}
		}
	}

	return nil
}

// Args 返回解析后的位置参数
// 必须在 Parse 之后调用
// 返回:
//
//	[]string: 位置参数列表
//
// 使用示例:
//
//	parser.Parse(os.Args[1:])
//	files := parser.Args() // 获取所有位置参数
func (p *Parser) Args() []string {
	if !p.parsed {
		return nil
	}
	return p.args
}

// NArg 返回位置参数的数量
// 返回:
//
//	int: 位置参数数量
func (p *Parser) NArg() int {
	return len(p.args)
}

// Arg 返回第 i 个位置参数
// 参数:
//
//	i: 参数索引(从0开始)
//
// 返回:
//
//	string: 位置参数,如果索引越界返回空字符串
func (p *Parser) Arg(i int) string {
	if i < 0 || i >= len(p.args) {
		return ""
	}
	return p.args[i]
}

// GetFlags 返回所有定义的标志
// 用于生成帮助文档
// 返回:
//
//	map[string]*Flag: 标志映射
func (p *Parser) GetFlags() map[string]*Flag {
	return p.flags
}

// isFlagSet 检查标志是否被用户设置
// 内部方法,用于验证必需标志
func (p *Parser) isFlagSet(name string) bool {
	found := false
	p.flagSet.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// getFlagValue 获取标志的当前值
// 内部方法,用于验证器
func (p *Parser) getFlagValue(name string) interface{} {
	f, ok := p.flags[name]
	if !ok {
		return nil
	}

	var value interface{}
	p.flagSet.Visit(func(fl *flag.Flag) {
		if fl.Name == name {
			switch f.Type {
			case "string":
				value = fl.Value.String()
			case "int":
				if v, err := strconv.Atoi(fl.Value.String()); err == nil {
					value = v
				}
			case "bool":
				value = fl.Value.String() == "true"
			case "[]string":
				value = fl.Value.String()
			}
		}
	})

	return value
}

// PrintDefaults 打印所有标志的默认值和用法
// 用于生成帮助信息
// 使用示例:
//
//	parser.PrintDefaults()
func (p *Parser) PrintDefaults() {
	p.flagSet.PrintDefaults()
}
