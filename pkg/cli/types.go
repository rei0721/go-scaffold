package cli

import (
	"context"
	"io"
)

// Command 定义 CLI 命令的接口
// 每个命令都应该实现这个接口,提供统一的执行模式
// 设计考虑:
// - 名称和描述用于帮助文档生成
// - Run 方法接收 Context 以支持超时和取消
// 使用示例:
//
//	type MyCommand struct {}
//	func (c *MyCommand) Name() string { return "mycommand" }
//	func (c *MyCommand) Description() string { return "My command description" }
//	func (c *MyCommand) Run(ctx context.Context, args []string) error { ... }
type Command interface {
	// Name 返回命令的名称
	// 用于命令行调用和帮助文档
	// 应该使用小写字母和连字符,如 "migrate-db"
	Name() string

	// Description 返回命令的简短描述
	// 用于帮助文档的命令列表
	// 应该是一句话的概括,不超过80个字符
	Description() string

	// Run 执行命令
	// 参数:
	//   ctx: 上下文,用于超时控制和取消
	//   args: 命令行参数(不包括程序名和命令名)
	// 返回:
	//   error: 执行失败时的错误
	Run(ctx context.Context, args []string) error
}

// FlagSet 表示命令行标志集合
// 封装了 flag.FlagSet 的功能,提供类型安全的参数访问
// 设计目标:
// - 简化标志定义和解析
// - 提供友好的错误消息
// - 支持常用的参数类型
type FlagSet struct {
	// name 标志集的名称,通常是命令名或程序名
	name string

	// flags 存储所有定义的标志
	flags map[string]*Flag

	// parsed 标记是否已解析
	parsed bool

	// args 存储解析后的位置参数
	args []string
}

// Flag 表示单个命令行标志
// 存储标志的元数据和值
type Flag struct {
	// Name 标志的名称
	Name string

	// ShortName 标志的简写形式(可选)
	ShortName string

	// Usage 标志的使用说明
	Usage string

	// Value 标志的值
	Value interface{}

	// DefaultValue 标志的默认值
	DefaultValue interface{}

	// Required 是否为必需标志
	Required bool

	// Type 标志的类型(string, int, bool 等)
	Type string
}

// Context 表示命令执行上下文
// 提供命令执行所需的所有信息和资源
// 设计考虑:
// - 包含标准输入输出,便于测试
// - 提供配置访问
// - 支持取消和超时
type Context struct {
	// context.Context 嵌入标准上下文,支持取消和超时
	context.Context

	// Command 当前执行的命令
	Command Command

	// FlagSet 解析后的标志集
	FlagSet *FlagSet

	// Args 位置参数
	Args []string

	// Stdin 标准输入
	// 默认为 os.Stdin,测试时可以替换
	Stdin io.Reader

	// Stdout 标准输出
	// 默认为 os.Stdout,测试时可以替换
	Stdout io.Writer

	// Stderr 标准错误输出
	// 默认为 os.Stderr,测试时可以替换
	Stderr io.Writer

	// Config 配置对象(可选)
	// 可以是任何类型,由命令自行解释
	Config interface{}
}

// HelpFormatter 定义帮助信息格式化器接口
// 允许自定义帮助文档的格式和样式
// 使用示例:
//
//	formatter := NewDefaultHelpFormatter()
//	help := formatter.Format(command, flagSet)
type HelpFormatter interface {
	// Format 格式化帮助信息
	// 参数:
	//   command: 要显示帮助的命令
	//   flagSet: 命令的标志集
	// 返回:
	//   string: 格式化后的帮助文本
	Format(command Command, flagSet *FlagSet) string

	// FormatUsage 格式化用法信息
	// 生成简短的用法说明
	// 返回:
	//   string: 格式化后的用法文本
	FormatUsage(command Command) string
}

// Runner 定义命令运行器接口
// 负责解析参数、执行命令和处理错误
// 设计目标:
// - 统一的命令执行流程
// - 标准化的错误处理
// - 可测试性
type Runner interface {
	// Run 运行命令
	// 这是 CLI 应用的主入口点
	// 参数:
	//   ctx: 上下文
	//   args: 命令行参数(通常是 os.Args)
	// 返回:
	//   int: 退出码
	Run(ctx context.Context, args []string) int

	// RegisterCommand 注册命令
	// 允许应用注册多个子命令
	// 参数:
	//   command: 要注册的命令
	RegisterCommand(command Command)

	// SetHelpFormatter 设置帮助格式化器
	// 允许自定义帮助文档格式
	// 参数:
	//   formatter: 自定义的格式化器
	SetHelpFormatter(formatter HelpFormatter)
}

// VersionInfo 表示版本信息
// 用于 --version 标志输出
type VersionInfo struct {
	// Name 应用程序名称
	Name string

	// Version 版本号
	// 建议使用语义化版本,如 "1.2.3"
	Version string

	// BuildTime 构建时间
	// 通常通过 ldflags 在编译时注入
	BuildTime string

	// GitCommit Git 提交哈希
	// 通常通过 ldflags 在编译时注入
	GitCommit string

	// GoVersion Go 版本
	GoVersion string
}

// String 返回版本信息的字符串表示
func (v *VersionInfo) String() string {
	if v.BuildTime != "" && v.GitCommit != "" {
		return VersionTemplate + "\n"
	}
	return v.Name + " version " + v.Version + "\n"
}

// Option 定义配置选项函数类型
// 使用函数选项模式(Functional Options Pattern)
// 提供灵活的配置方式
// 使用示例:
//
//	runner := NewRunner(
//	  WithStdout(os.Stdout),
//	  WithStderr(os.Stderr),
//	)
type Option func(interface{})

// ValidatorFunc 定义标志值验证函数类型
// 用于自定义参数验证逻辑
// 使用示例:
//
//	func validatePort(value interface{}) error {
//	  port := value.(int)
//	  if port < 1 || port > 65535 {
//	    return fmt.Errorf("port must be between 1 and 65535")
//	  }
//	  return nil
//	}
type ValidatorFunc func(value interface{}) error
