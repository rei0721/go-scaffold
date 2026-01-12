package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

// SimpleRunner 简单的命令运行器实现
// 提供基本的命令执行功能
// 适合单命令 CLI 应用
type SimpleRunner struct {
	// config 运行器配置
	config *Config

	// command 要执行的命令
	command Command

	// parser 参数解析器
	parser *Parser

	// helpFormatter 帮助格式化器
	helpFormatter HelpFormatter

	// stdout 标准输出
	stdout io.Writer

	// stderr 标准错误输出
	stderr io.Writer

	// versionInfo 版本信息
	versionInfo *VersionInfo
}

// NewSimpleRunner 创建一个新的简单运行器
// 参数:
//
//	command: 要执行的命令
//	options: 配置选项(可选)
//
// 返回:
//
//	*SimpleRunner: 新创建的运行器
//
// 使用示例:
//
//	runner := NewSimpleRunner(myCommand, WithConfig(config))
//	exitCode := runner.Run(context.Background(), os.Args)
func NewSimpleRunner(command Command, options ...Option) *SimpleRunner {
	runner := &SimpleRunner{
		config:        DefaultConfig(),
		command:       command,
		parser:        NewParser(command.Name()),
		helpFormatter: &DefaultHelpFormatter{},
		stdout:        os.Stdout,
		stderr:        os.Stderr,
	}

	// 应用选项
	for _, option := range options {
		option(runner)
	}

	return runner
}

// Run 运行命令
// 实现 Runner 接口
// 参数:
//
//	ctx: 上下文
//	args: 命令行参数
//
// 返回:
//
//	int: 退出码
func (r *SimpleRunner) Run(ctx context.Context, args []string) int {
	// 1. 解析参数
	if len(args) > 0 {
		args = args[1:] // 去掉程序名
	}

	// 2. 检查内置标志
	for _, arg := range args {
		// 检查 help 标志
		if r.config.EnableHelp && (arg == "-h" || arg == "--help") {
			r.showHelp()
			return ExitSuccess
		}

		// 检查 version 标志
		if r.config.EnableVersion && (arg == "-v" || arg == "--version") {
			r.showVersion()
			return ExitSuccess
		}
	}

	// 3. 解析参数
	if err := r.parser.Parse(args); err != nil {
		fmt.Fprintf(r.stderr, "%s Error: %v\n", IndicatorError, err)
		if r.config.ExitOnError {
			r.showUsage()
		}
		return ExitInvalidUsage
	}

	// 4. 创建命令上下文
	cmdCtx := &Context{
		Context: ctx,
		Command: r.command,
		FlagSet: &FlagSet{
			name:   r.parser.flagSet.Name(),
			flags:  r.parser.flags,
			parsed: r.parser.parsed,
			args:   r.parser.args,
		},
		Args:   r.parser.Args(),
		Stdin:  os.Stdin,
		Stdout: r.stdout,
		Stderr: r.stderr,
	}

	// 5. 执行命令
	if err := r.command.Run(cmdCtx, cmdCtx.Args); err != nil {
		fmt.Fprintf(r.stderr, "%s %s\n", IndicatorError, err)
		return ExitError
	}

	return ExitSuccess
}

// RegisterCommand 注册命令
// 对于 SimpleRunner,这个方法会替换当前命令
// 参数:
//
//	command: 要注册的命令
func (r *SimpleRunner) RegisterCommand(command Command) {
	r.command = command
	r.parser = NewParser(command.Name())
}

// SetHelpFormatter 设置帮助格式化器
// 参数:
//
//	formatter: 自定义的格式化器
func (r *SimpleRunner) SetHelpFormatter(formatter HelpFormatter) {
	r.helpFormatter = formatter
}

// GetParser 获取参数解析器
// 用于定义命令行标志
// 返回:
//
//	*Parser: 参数解析器
//
// 使用示例:
//
//	parser := runner.GetParser()
//	parser.String("output", "o", "./output", "输出目录")
func (r *SimpleRunner) GetParser() *Parser {
	return r.parser
}

// showHelp 显示帮助信息
func (r *SimpleRunner) showHelp() {
	help := r.helpFormatter.Format(r.command, &FlagSet{
		name:  r.parser.flagSet.Name(),
		flags: r.parser.flags,
	})
	fmt.Fprint(r.stdout, help)
}

// showVersion 显示版本信息
func (r *SimpleRunner) showVersion() {
	if r.versionInfo != nil {
		fmt.Fprint(r.stdout, r.versionInfo.String())
	} else {
		fmt.Fprintf(r.stdout, "%s version %s\n", r.config.AppName, r.config.Version)
	}
}

// showUsage 显示用法信息
func (r *SimpleRunner) showUsage() {
	usage := r.helpFormatter.FormatUsage(r.command)
	fmt.Fprint(r.stderr, usage)
	fmt.Fprintf(r.stderr, "\nUse -h or --help for more information.\n")
}

// DefaultHelpFormatter 默认的帮助格式化器
type DefaultHelpFormatter struct{}

// Format 格式化帮助信息
func (f *DefaultHelpFormatter) Format(command Command, flagSet *FlagSet) string {
	var help strings.Builder

	// 头部
	help.WriteString(fmt.Sprintf("%s - %s\n\n", command.Name(), command.Description()))

	// 用法
	help.WriteString("Usage:\n")
	help.WriteString(fmt.Sprintf("  %s [options] [arguments]\n\n", command.Name()))

	// 选项
	if len(flagSet.flags) > 0 {
		help.WriteString("Options:\n")
		for name, flag := range flagSet.flags {
			if flag.ShortName != "" {
				help.WriteString(fmt.Sprintf("  -%s, --%s", flag.ShortName, name))
			} else {
				help.WriteString(fmt.Sprintf("  --%s", name))
			}

			if flag.Type != "bool" {
				help.WriteString(fmt.Sprintf(" (%s)", flag.Type))
			}

			help.WriteString("\n")
			if flag.Usage != "" {
				help.WriteString(fmt.Sprintf("        %s", flag.Usage))
			}
			if flag.DefaultValue != nil && flag.DefaultValue != "" && flag.DefaultValue != false && flag.DefaultValue != 0 {
				help.WriteString(fmt.Sprintf(" (default: %v)", flag.DefaultValue))
			}
			if flag.Required {
				help.WriteString(" [required]")
			}
			help.WriteString("\n")
		}
	}

	return help.String()
}

// FormatUsage 格式化用法信息
func (f *DefaultHelpFormatter) FormatUsage(command Command) string {
	return fmt.Sprintf("Usage: %s [options] [arguments]", command.Name())
}

// 配置选项函数

// WithConfig 设置配置
func WithConfig(config *Config) Option {
	return func(v interface{}) {
		if runner, ok := v.(*SimpleRunner); ok {
			runner.config = config
		}
	}
}

// WithStdout 设置标准输出
func WithStdout(stdout io.Writer) Option {
	return func(v interface{}) {
		if runner, ok := v.(*SimpleRunner); ok {
			runner.stdout = stdout
		}
	}
}

// WithStderr 设置标准错误输出
func WithStderr(stderr io.Writer) Option {
	return func(v interface{}) {
		if runner, ok := v.(*SimpleRunner); ok {
			runner.stderr = stderr
		}
	}
}

// WithVersionInfo 设置版本信息
func WithVersionInfo(versionInfo *VersionInfo) Option {
	return func(v interface{}) {
		if runner, ok := v.(*SimpleRunner); ok {
			runner.versionInfo = versionInfo
		}
	}
}
