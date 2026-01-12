package cli

// 退出码常量
// 标准化程序退出状态,便于脚本和监控系统识别
const (
	// ExitSuccess 成功退出码
	// 程序正常完成所有任务
	ExitSuccess = 0

	// ExitError 一般错误退出码
	// 程序执行过程中遇到错误
	ExitError = 1

	// ExitInvalidUsage 无效使用退出码
	// 用户提供了无效的命令行参数或选项
	ExitInvalidUsage = 2

	// ExitInvalidConfig 无效配置退出码
	// 配置文件格式错误或内容无效
	ExitInvalidConfig = 3

	// ExitIOError IO 错误退出码
	// 文件读写或网络操作失败
	ExitIOError = 4

	// ExitInterrupted 中断退出码
	// 程序被用户或系统中断(Ctrl+C)
	ExitInterrupted = 130
)

// 默认标志常量
// 常用命令行标志的标准名称
const (
	// DefaultHelpFlag 帮助标志
	// 用于显示命令使用说明
	DefaultHelpFlag = "help"

	// DefaultHelpShortFlag 帮助标志简写
	DefaultHelpShortFlag = "h"

	// DefaultVersionFlag 版本标志
	// 用于显示程序版本信息
	DefaultVersionFlag = "version"

	// DefaultVersionShortFlag 版本标志简写
	DefaultVersionShortFlag = "v"

	// DefaultVerboseFlag 详细输出标志
	// 启用详细日志输出
	DefaultVerboseFlag = "verbose"

	// DefaultQuietFlag 静默模式标志
	// 禁用非必要的输出
	DefaultQuietFlag = "quiet"

	// DefaultConfigFlag 配置文件标志
	// 指定配置文件路径
	DefaultConfigFlag = "config"

	// DefaultOutputFlag 输出路径标志
	// 指定输出文件或目录
	DefaultOutputFlag = "output"
)

// 默认值常量
// 命令行参数的常见默认值
const (
	// DefaultConfigPath 默认配置文件路径
	DefaultConfigPath = "./config.yaml"

	// DefaultOutputDir 默认输出目录
	DefaultOutputDir = "./output"

	// DefaultTimeout 默认超时时间(秒)
	DefaultTimeout = 30

	// DefaultMaxRetries 默认最大重试次数
	DefaultMaxRetries = 3
)

// 日志消息常量
// 避免在代码中使用魔法字符串,便于统一管理和修改
const (
	// MsgCommandStarting 命令开始执行消息
	MsgCommandStarting = "command starting"

	// MsgCommandCompleted 命令完成消息
	MsgCommandCompleted = "command completed successfully"

	// MsgCommandFailed 命令失败消息
	MsgCommandFailed = "command execution failed"

	// MsgParsingFlags 解析标志消息
	MsgParsingFlags = "parsing command line flags"

	// MsgFlagsParsed 标志解析完成消息
	MsgFlagsParsed = "flags parsed successfully"

	// MsgValidatingConfig 验证配置消息
	MsgValidatingConfig = "validating configuration"

	// MsgConfigValid 配置有效消息
	MsgConfigValid = "configuration is valid"

	// MsgLoadingConfig 加载配置消息
	MsgLoadingConfig = "loading configuration"

	// MsgConfigLoaded 配置加载完成消息
	MsgConfigLoaded = "configuration loaded successfully"

	// MsgShowingHelp 显示帮助消息
	MsgShowingHelp = "showing help information"

	// MsgShowingVersion 显示版本消息
	MsgShowingVersion = "showing version information"
)

// 错误消息常量
// 用于创建错误时的统一消息格式
const (
	// ErrMsgInvalidFlag 无效标志错误消息格式
	// 使用 fmt.Sprintf(ErrMsgInvalidFlag, flagName)
	ErrMsgInvalidFlag = "invalid flag: %s"

	// ErrMsgMissingRequired 缺少必需参数错误消息格式
	// 使用 fmt.Sprintf(ErrMsgMissingRequired, paramName)
	ErrMsgMissingRequired = "missing required parameter: %s"

	// ErrMsgInvalidValue 无效值错误消息格式
	// 使用 fmt.Sprintf(ErrMsgInvalidValue, paramName, value)
	ErrMsgInvalidValue = "invalid value for %s: %v"

	// ErrMsgCommandFailed 命令失败错误消息格式
	// 使用 fmt.Sprintf(ErrMsgCommandFailed, commandName, err)
	ErrMsgCommandFailedFmt = "command '%s' failed: %w"

	// ErrMsgConfigLoadFailed 配置加载失败错误消息格式
	// 使用 fmt.Sprintf(ErrMsgConfigLoadFailed, path, err)
	ErrMsgConfigLoadFailed = "failed to load config from %s: %w"

	// ErrMsgConfigInvalid 配置无效错误消息格式
	// 使用 fmt.Sprintf(ErrMsgConfigInvalid, err)
	ErrMsgConfigInvalid = "invalid configuration: %w"

	// ErrMsgOutputFailed 输出失败错误消息格式
	// 使用 fmt.Sprintf(ErrMsgOutputFailed, path, err)
	ErrMsgOutputFailed = "failed to write output to %s: %w"

	// ErrMsgNoCommand 未指定命令错误消息
	ErrMsgNoCommand = "no command specified"

	// ErrMsgUnknownCommand 未知命令错误消息格式
	// 使用 fmt.Sprintf(ErrMsgUnknownCommand, commandName)
	ErrMsgUnknownCommand = "unknown command: %s"
)

// 帮助信息模板常量
// 用于格式化帮助文档的模板字符串
const (
	// HelpTemplateHeader 帮助信息头部模板
	// 使用格式: fmt.Sprintf(HelpTemplateHeader, appName, description)
	HelpTemplateHeader = `%s - %s

Usage:
  %s [options] [arguments]
`

	// HelpTemplateOptions 选项部分模板
	HelpTemplateOptions = `
Options:
`

	// HelpTemplateExamples 示例部分模板
	HelpTemplateExamples = `
Examples:
`

	// HelpTemplateFooter 帮助信息底部模板
	HelpTemplateFooter = `
For more information, visit: %s
`
)

// 版本信息模板常量
const (
	// VersionTemplate 版本信息模板
	// 使用格式: fmt.Sprintf(VersionTemplate, appName, version, buildTime, commit)
	VersionTemplate = `%s version %s
Build Time: %s
Git Commit: %s
`
)

// 进度指示器常量
// 用于控制台输出的视觉元素
const (
	// IndicatorSuccess 成功指示器
	IndicatorSuccess = "✅"

	// IndicatorError 错误指示器
	IndicatorError = "❌"

	// IndicatorWarning 警告指示器
	IndicatorWarning = "⚠️"

	// IndicatorInfo 信息指示器
	IndicatorInfo = "ℹ️"

	// IndicatorProgress 进度指示器
	IndicatorProgress = "⏳"

	// IndicatorFolder 文件夹指示器
	IndicatorFolder = "📁"

	// IndicatorFile 文件指示器
	IndicatorFile = "📄"
)
