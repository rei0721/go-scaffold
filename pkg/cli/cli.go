package cli

import (
	"io"
)

// FlagType 表示选项的类型
type FlagType int

const (
	// FlagTypeString 字符串类型
	FlagTypeString FlagType = iota
	// FlagTypeInt 整数类型
	FlagTypeInt
	// FlagTypeBool 布尔类型
	FlagTypeBool
	// FlagTypeStringSlice 字符串数组类型
	FlagTypeStringSlice
)

// Flag 表示一个命令行选项
type Flag struct {
	// Name 选项名称 (长选项，如 "output")
	Name string
	// ShortName 短选项名称 (如 "o")
	ShortName string
	// Type 选项类型
	Type FlagType
	// Required 是否必填
	Required bool
	// Default 默认值
	Default interface{}
	// Description 选项描述
	Description string
	// EnvVar 环境变量名 (用于回退)
	EnvVar string
}

// App CLI 应用接口
type App interface {
	// Name 返回应用名称
	Name() string
	// Version 返回应用版本
	Version() string
	// Description 返回应用描述
	Description() string
	// SetVersion 设置版本号
	SetVersion(version string)
	// SetDescription 设置描述
	SetDescription(desc string)
	// AddCommand 注册子命令
	AddCommand(cmd Command) error
	// Run 执行 CLI，解析参数并路由到对应命令
	Run(args []string) error
	// RunWithIO 执行 CLI，使用自定义 I/O (用于测试)
	RunWithIO(args []string, stdin io.Reader, stdout, stderr io.Writer) error
}

// Command 命令接口
type Command interface {
	// Name 返回命令名称 (如 "generate", "migrate")
	Name() string
	// Description 返回命令描述 (用于 help 输出)
	Description() string
	// Usage 返回使用说明 (可选，用于详细帮助)
	Usage() string
	// Flags 返回命令支持的选项列表
	Flags() []Flag
	// Execute 执行命令逻辑
	Execute(ctx *Context) error
}

// Context 命令执行上下文
type Context struct {
	// Args 位置参数 (去除命令名和选项后的参数)
	Args []string
	// Flags 解析后的选项值
	Flags map[string]interface{}
	// Stdin 标准输入
	Stdin io.Reader
	// Stdout 标准输出
	Stdout io.Writer
	// Stderr 标准错误输出
	Stderr io.Writer
}

// GetString 获取字符串类型的选项值
func (c *Context) GetString(name string) string {
	if v, ok := c.Flags[name]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetInt 获取整数类型的选项值
func (c *Context) GetInt(name string) int {
	if v, ok := c.Flags[name]; ok {
		if i, ok := v.(int); ok {
			return i
		}
	}
	return 0
}

// GetBool 获取布尔类型的选项值
func (c *Context) GetBool(name string) bool {
	if v, ok := c.Flags[name]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// GetStringSlice 获取字符串数组类型的选项值
func (c *Context) GetStringSlice(name string) []string {
	if v, ok := c.Flags[name]; ok {
		if s, ok := v.([]string); ok {
			return s
		}
	}
	return nil
}
