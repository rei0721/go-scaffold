package dbtx

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/rei0721/go-scaffold/pkg/logger"
	"gorm.io/gorm"
)

// manager 实现 Manager 接口
// 提供基于 GORM 的事务管理功能
type manager struct {
	// db GORM 数据库实例
	db *gorm.DB

	// logger 日志记录器（可选）
	// 使用 atomic.Value 支持延迟注入
	logger atomic.Value
}

// NewManager 创建一个新的事务管理器
// 参数:
//
//	db: GORM 数据库实例，必须非空
//	log: 日志记录器，可以为 nil（不记录日志）
//
// 返回:
//
//	Manager: 事务管理器实例
//	error: 如果 db 为 nil，返回 ErrNilDB
//
// 示例:
//
//	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
//	txManager, err := dbtx.NewManager(db, logger)
//	if err != nil {
//	    log.Fatal(err)
//	}
func NewManager(db *gorm.DB, log logger.Logger) (Manager, error) {
	if db == nil {
		return nil, ErrNilDB
	}

	m := &manager{
		db: db,
	}

	// 注入 logger（可选）
	if log != nil {
		m.logger.Store(log)
	}

	return m, nil
}

// SetLogger 设置日志记录器（延迟注入）
// 参数:
//
//	log: 日志记录器实例
func (m *manager) SetLogger(log logger.Logger) {
	m.logger.Store(log)
}

// getLogger 获取日志记录器
// 返回:
//
//	logger.Logger: 如果已注入返回实例，否则返回 nil
func (m *manager) getLogger() logger.Logger {
	if log := m.logger.Load(); log != nil {
		return log.(logger.Logger)
	}
	return nil
}

// logEvent 记录事务事件
// 参数:
//
//	event: 事件类型（如 tx_begin, tx_commit, tx_rollback）
//	message: 事件消息
//	fields: 附加字段
func (m *manager) logEvent(event string, message string, fields ...interface{}) {
	if log := m.getLogger(); log != nil {
		allFields := []interface{}{"event", event, "message", message}
		allFields = append(allFields, fields...)
		log.Debug("transaction event", allFields...)
	}
}

// GetDB 返回底层的数据库实例
func (m *manager) GetDB() *gorm.DB {
	return m.db
}

// WithTx 在事务中执行函数
// 使用默认选项
func (m *manager) WithTx(ctx context.Context, fn TxFunc) error {
	return m.WithTxOptions(ctx, nil, fn)
}

// WithTxOptions 使用自定义选项执行事务
// 这是核心实现，负责事务的完整生命周期管理
func (m *manager) WithTxOptions(ctx context.Context, opts *TxOptions, fn TxFunc) error {
	// 1. 验证参数
	if fn == nil {
		return ErrNilTxFunc
	}

	// 2. 处理选项
	if opts == nil {
		opts = DefaultOptions()
	} else {
		if err := opts.Validate(); err != nil {
			return err
		}
	}

	// 3. 创建带超时的 Context
	txCtx := ctx
	var cancel context.CancelFunc
	if opts.Timeout > 0 {
		txCtx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	// 4. 检查是否在事务中（嵌套事务检测）
	// GORM 会自动处理嵌套事务（使用 SavePoint）
	// 如果禁用嵌套事务，需要检查
	tx := m.db.WithContext(txCtx)

	// 5. 配置事务选项
	if opts.ReadOnly {
		// 设置只读模式
		tx = tx.Set("gorm:query_option", "FOR SHARE")
	}

	// 6. 开启事务
	m.logEvent(LogEventBegin, "starting transaction",
		"isolation", opts.Isolation,
		"readonly", opts.ReadOnly,
		"timeout", opts.Timeout,
	)

	tx = tx.Begin()
	if tx.Error != nil {
		m.logEvent(LogEventError, "failed to begin transaction", "error", tx.Error)
		return fmt.Errorf(ErrMsgBeginFailed, tx.Error)
	}

	// 7. 设置隔离级别（如果需要）
	// 注意: 必须在 Begin 之后设置
	if opts.Isolation != 0 {
		// GORM 不直接支持设置隔离级别，需要使用原始 SQL
		// 不同数据库的 SQL 语法不同，这里仅作示例
		// 实际使用时建议在连接配置中设置默认隔离级别
	}

	// 8. 确保事务会被提交或回滚
	var committed bool
	defer func() {
		if r := recover(); r != nil {
			// panic 发生，回滚事务
			m.logEvent(LogEventPanic, "panic occurred, rolling back",
				"panic", r,
			)
			tx.Rollback()
			panic(r) // 重新抛出 panic
		} else if !committed {
			// 如果没有提交（说明发生了错误），回滚
			m.logEvent(LogEventRollback, "rolling back transaction")
			tx.Rollback()
		}
	}()

	// 9. 执行业务逻辑
	err := fn(tx)

	// 10. 检查 Context 是否被取消
	select {
	case <-txCtx.Done():
		m.logEvent(LogEventTimeout, "transaction context done",
			"error", txCtx.Err(),
		)
		return fmt.Errorf(ErrMsgContextCanceled, txCtx.Err())
	default:
	}

	// 11. 根据错误决定提交或回滚
	if err != nil {
		m.logEvent(LogEventError, "transaction function returned error",
			"error", err,
		)
		return fmt.Errorf(ErrMsgTxFuncError, err)
	}

	// 12. 提交事务
	if commitErr := tx.Commit().Error; commitErr != nil {
		m.logEvent(LogEventError, "failed to commit transaction",
			"error", commitErr,
		)
		return fmt.Errorf(ErrMsgCommitFailed, commitErr)
	}

	committed = true
	m.logEvent(LogEventCommit, "transaction committed successfully")
	return nil
}
