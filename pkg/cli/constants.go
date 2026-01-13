package cli

// 错误常量
// 遵循 Unix 退出码约定
const (
	// ExitSuccess 成功退出码
	ExitSuccess = 0
	// ExitError 通用错误
	ExitError = 1
	// ExitUsage 参数错误
	ExitUsage = 2
	// ExitConfig 配置错误
	ExitConfig = 3
	// ExitInterrupted 用户中断 (Ctrl+C)
	ExitInterrupted = 130
)

// 错误消息常量
const (
	// ErrMsgCommandNotFound 命令未找到
	ErrMsgCommandNotFound = "command not found"
	// ErrMsgInvalidArgs 无效参数
	ErrMsgInvalidArgs = "invalid arguments"
	// ErrMsgMissingRequired 缺少必填参数
	ErrMsgMissingRequired = "missing required flag"
	// ErrMsgDuplicateCommand 重复的命令名
	ErrMsgDuplicateCommand = "duplicate command name"
	// ErrMsgCancelled 操作已取消
	ErrMsgCancelled = "operation cancelled"
	// ErrMsgInvalidFlagValue 无效的选项值
	ErrMsgInvalidFlagValue = "invalid flag value"
)

// 默认值
const (
	// DefaultHelpFlag help 选项名
	DefaultHelpFlag = "help"
	// DefaultVersionFlag version 选项名
	DefaultVersionFlag = "version"
)
