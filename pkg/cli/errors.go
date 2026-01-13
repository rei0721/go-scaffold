package cli

import "fmt"

// UsageError 表示参数使用错误
type UsageError struct {
	Command string
	Message string
}

// Error 实现 error 接口
func (e *UsageError) Error() string {
	if e.Command != "" {
		return fmt.Sprintf("%s: %s\nRun '%s --help' for usage", e.Command, e.Message, e.Command)
	}
	return fmt.Sprintf("%s\nRun '--help' for usage", e.Message)
}

// ExitCode 返回退出码
func (e *UsageError) ExitCode() int {
	return ExitUsage
}

// CommandError 表示命令执行错误
type CommandError struct {
	Command string
	Message string
	Cause   error
}

// Error 实现 error 接口
func (e *CommandError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Command, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Command, e.Message)
}

// Unwrap 实现 errors.Unwrap 接口
func (e *CommandError) Unwrap() error {
	return e.Cause
}

// ExitCode 返回退出码
func (e *CommandError) ExitCode() int {
	return ExitError
}

// CancelledError 表示用户取消操作
type CancelledError struct {
	Message string
}

// Error 实现 error 接口
func (e *CancelledError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return ErrMsgCancelled
}

// ExitCode 返回退出码
func (e *CancelledError) ExitCode() int {
	return ExitInterrupted
}

// ExitCoder 定义可返回退出码的错误接口
type ExitCoder interface {
	error
	ExitCode() int
}

// GetExitCode 从错误中提取退出码
// 如果错误实现了 ExitCoder 接口,返回其退出码
// 否则返回通用错误码
func GetExitCode(err error) int {
	if err == nil {
		return ExitSuccess
	}
	if ec, ok := err.(ExitCoder); ok {
		return ec.ExitCode()
	}
	return ExitError
}
