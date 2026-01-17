package dbtx

import "time"

// 默认配置常量
const (
	// DefaultTimeout 默认事务超时时间
	// 推荐值: 30秒，避免长事务阻塞
	DefaultTimeout = 30 * time.Second

	// DefaultMaxRetries 默认最大重试次数
	// 用于可重试的事务（如死锁、序列化失败）
	DefaultMaxRetries = 3

	// DefaultRetryDelay 默认重试延迟
	// 使用指数退避策略时的基准延迟
	DefaultRetryDelay = 100 * time.Millisecond
)

// 日志事件类型
const (
	// LogEventBegin 事务开始
	LogEventBegin = "tx_begin"

	// LogEventCommit 事务提交
	LogEventCommit = "tx_commit"

	// LogEventRollback 事务回滚
	LogEventRollback = "tx_rollback"

	// LogEventPanic panic发生
	LogEventPanic = "tx_panic"

	// LogEventError 错误发生
	LogEventError = "tx_error"

	// LogEventTimeout 超时发生
	LogEventTimeout = "tx_timeout"

	// LogEventNested 嵌套事务
	LogEventNested = "tx_nested"
)

// SavePoint 保存点名称前缀
// GORM 会自动管理 SavePoint，此处仅作文档说明
const (
	// SavePointPrefix SavePoint 名称前缀
	// GORM 格式: sp + 递增数字，如 sp1, sp2, sp3
	SavePointPrefix = "sp"
)
