package dbtx

import "errors"

// 预定义错误（Sentinel Errors）
// 可使用 errors.Is() 判断
var (
	// ErrTxNotStarted 事务未开启
	ErrTxNotStarted = errors.New("transaction not started")

	// ErrTxAlreadyStarted 事务已经开启
	ErrTxAlreadyStarted = errors.New("transaction already started")

	// ErrTxRolledBack 事务已回滚
	ErrTxRolledBack = errors.New("transaction rolled back")

	// ErrTxCommitted 事务已提交
	ErrTxCommitted = errors.New("transaction already committed")

	// ErrInvalidOptions 无效的事务选项
	ErrInvalidOptions = errors.New("invalid transaction options")

	// ErrNilDB 数据库实例为空
	ErrNilDB = errors.New("database instance is nil")

	// ErrNilTxFunc 事务函数为空
	ErrNilTxFunc = errors.New("transaction function is nil")
)

// 错误消息模板常量
// 用于 fmt.Errorf() 包装错误
const (
	// ErrMsgBeginFailed 开启事务失败
	ErrMsgBeginFailed = "failed to begin transaction: %w"

	// ErrMsgCommitFailed 提交事务失败
	ErrMsgCommitFailed = "failed to commit transaction: %w"

	// ErrMsgRollbackFailed 回滚事务失败
	ErrMsgRollbackFailed = "failed to rollback transaction: %w"

	// ErrMsgTxFuncPanic 事务函数发生panic
	ErrMsgTxFuncPanic = "transaction function panicked: %v"

	// ErrMsgTxFuncError 事务函数返回错误
	ErrMsgTxFuncError = "transaction function returned error: %w"

	// ErrMsgContextCanceled Context已取消
	ErrMsgContextCanceled = "transaction context canceled: %w"
)
