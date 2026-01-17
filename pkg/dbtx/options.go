package dbtx

import (
	"database/sql"
	"time"
)

// TxOptions 事务选项配置
// 用于自定义事务的隔离级别、超时等参数
type TxOptions struct {
	// Isolation 事务隔离级别
	// 可选值:
	//   - sql.LevelDefault: 使用数据库默认隔离级别
	//   - sql.LevelReadUncommitted: 读未提交
	//   - sql.LevelReadCommitted: 读已提交（PostgreSQL 默认）
	//   - sql.LevelRepeatableRead: 可重复读（MySQL 默认）
	//   - sql.LevelSerializable: 串行化
	// 注意: 不同数据库支持的隔离级别不同
	Isolation sql.IsolationLevel

	// ReadOnly 是否为只读事务
	// 只读事务的优势:
	//   - 性能更好（无需加写锁）
	//   - 某些数据库可以使用快照读
	//   - 避免误操作修改数据
	// 注意: 只读事务中执行写操作会返回错误
	ReadOnly bool

	// Timeout 事务超时时间
	// 如果设置为 0，使用 DefaultTimeout
	// 超时后 Context 会被取消，事务会回滚
	// 推荐: 根据业务复杂度设置，一般 5-30 秒
	Timeout time.Duration

	// DisableNestedTransaction 禁用嵌套事务
	// 如果为 true，嵌套调用 WithTx 时会返回错误
	// 如果为 false（默认），使用 GORM SavePoint 机制
	DisableNestedTransaction bool
}

// DefaultOptions 返回默认的事务选项
// 隔离级别: 使用数据库默认
// 只读: false
// 超时: 30秒
// 嵌套事务: 启用
func DefaultOptions() *TxOptions {
	return &TxOptions{
		Isolation:                sql.LevelDefault,
		ReadOnly:                 false,
		Timeout:                  DefaultTimeout,
		DisableNestedTransaction: false,
	}
}

// Validate 验证选项的有效性
// 返回:
//
//	error: 如果选项无效，返回错误
func (opts *TxOptions) Validate() error {
	if opts == nil {
		return ErrInvalidOptions
	}

	// 验证超时时间
	if opts.Timeout < 0 {
		return ErrInvalidOptions
	}

	// 如果超时为 0，使用默认值
	if opts.Timeout == 0 {
		opts.Timeout = DefaultTimeout
	}

	return nil
}

// WithIsolation 设置隔离级别
func (opts *TxOptions) WithIsolation(level sql.IsolationLevel) *TxOptions {
	opts.Isolation = level
	return opts
}

// WithReadOnly 设置只读模式
func (opts *TxOptions) WithReadOnly(readOnly bool) *TxOptions {
	opts.ReadOnly = readOnly
	return opts
}

// WithTimeout 设置超时时间
func (opts *TxOptions) WithTimeout(timeout time.Duration) *TxOptions {
	opts.Timeout = timeout
	return opts
}

// WithDisableNested 设置是否禁用嵌套事务
func (opts *TxOptions) WithDisableNested(disable bool) *TxOptions {
	opts.DisableNestedTransaction = disable
	return opts
}
