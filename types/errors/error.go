package errors

import "fmt"

// BizError 表示一个业务错误,包含错误码、错误消息和可选的原因错误
// 这是应用程序中所有业务错误的基础类型
// 设计思想:
// - 错误码: 用于区分不同类型的错误,前端可以根据错误码做不同处理
// - 错误消息: 人类可读的错误描述,可以直接展示给用户
// - 原因错误: 保留底层错误,便于调试和日志记录
// 优点:
// - 统一的错误格式,便于 API 响应
// - 支持错误链,可以追踪错误的根本原因
// - 实现了 error 接口,可以像标准错误一样使用
type BizError struct {
	// Code 错误码
	// 每种错误类型都有唯一的数字代码
	// 例如: 1001-参数错误, 2001-业务错误, 5001-系统错误
	Code int

	// Message 错误消息
	// 这是给用户看的友好描述,应该是人类可读的
	// 不应该包含敏感信息(如数据库错误详情、堆栈信息等)
	Message string

	// Cause 原因错误
	// 保存底层的错误对象,用于错误追踪
	// 例如:数据库查询失败的具体原因、网络请求超时等
	// 可选字段,如果没有底层错误可以为 nil
	Cause error
}

// NewBizError 创建一个新的业务错误
// 这是创建 BizError 的推荐方式
// 参数:
//   code: 错误码,应该使用 codes.go 中定义的常量
//   message: 错误消息,应该简洁明确
// 返回:
//   *BizError: 业务错误对象
// 示例:
//   err := NewBizError(ErrUserNotFound, "用户不存在")
//   err := NewBizError(ErrInvalidParams, "用户名长度必须在3-50之间")
func NewBizError(code int, message string) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
	}
}

// Error 实现 error 接口
// 这使得 BizError 可以像标准 error 一样使用
// 返回格式:
//   - 有原因错误: "[错误码] 错误消息: 原因错误"
//   - 无原因错误: "[错误码] 错误消息"
// 这种格式既包含了业务信息,又保留了技术细节
func (e *BizError) Error() string {
	if e.Cause != nil {
		// 如果有原因错误,包含在错误信息中
		// 这样在日志中可以看到完整的错误链
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Cause)
	}
	// 只有业务错误本身
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// WithCause 附加原因错误到 BizError 并返回自身
// 这是一个链式调用方法,使代码更简洁
// 参数:
//   err: 底层错误,通常来自数据库、网络等操作
// 返回:
//   *BizError: 返回自身,支持链式调用
// 使用场景:
//   - 数据库操作失败
//   - 网络请求失败
//   - 文件操作失败
// 示例:
//   return nil, NewBizError(ErrDatabaseError, "查询失败").WithCause(sqlErr)
//   这样既有友好的业务错误消息,又保留了技术错误详情
func (e *BizError) WithCause(err error) *BizError {
	e.Cause = err
	// 返回自身,支持链式调用
	// 例如: NewBizError(code, msg).WithCause(err)
	return e
}

// Unwrap 返回底层的原因错误
// 实现这个方法是为了支持 Go 1.13+ 的错误处理功能
// 好处:
//   - errors.Is() 可以判断错误链中是否包含特定错误
//   - errors.As() 可以从错误链中提取特定类型的错误
//   - 错误可以被正确地包装和解包
// 使用场景:
//   if errors.Is(err, sql.ErrNoRows) {
//       // 可以正确判断底层是否是"记录不存在"错误
//   }
func (e *BizError) Unwrap() error {
	return e.Cause
}
