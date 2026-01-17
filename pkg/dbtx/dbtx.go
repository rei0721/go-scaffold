package dbtx

import (
	"context"

	"gorm.io/gorm"
)

// TxFunc 事务函数签名
// 在事务函数中执行业务逻辑
// 参数:
//
//	tx: GORM 事务实例，使用此实例执行数据库操作
//
// 返回:
//
//	error: 如果返回非 nil 错误，事务会自动回滚
//	       如果返回 nil，事务会自动提交
//
// 示例:
//
//	err := txManager.WithTx(ctx, func(tx *gorm.DB) error {
//	    // 使用 tx 执行数据库操作
//	    if err := tx.Create(&user).Error; err != nil {
//	        return err // 自动回滚
//	    }
//	    return nil // 自动提交
//	})
type TxFunc func(tx *gorm.DB) error

// Manager 事务管理器接口
// 提供事务管理的核心功能，支持自动提交/回滚、嵌套事务等
type Manager interface {
	// WithTx 在事务中执行函数
	// 自动处理事务的开启、提交、回滚
	// 使用默认选项（隔离级别: 数据库默认，超时: 30秒）
	//
	// 参数:
	//   ctx: 上下文，用于超时控制和取消
	//   fn: 事务函数，在此函数中执行业务逻辑
	//
	// 返回:
	//   error: 事务执行失败的错误
	//
	// 行为:
	//   - 如果 fn 返回 nil，事务会自动提交
	//   - 如果 fn 返回 error，事务会自动回滚
	//   - 如果 fn 发生 panic，事务会自动回滚并重新抛出 panic
	//   - 如果 ctx 被取消，事务会回滚
	//
	// 示例:
	//   err := manager.WithTx(ctx, func(tx *gorm.DB) error {
	//       return tx.Create(&user).Error
	//   })
	WithTx(ctx context.Context, fn TxFunc) error

	// WithTxOptions 使用自定义选项执行事务
	// 允许指定隔离级别、只读模式、超时等
	//
	// 参数:
	//   ctx: 上下文
	//   opts: 事务选项，如果为 nil 则使用默认选项
	//   fn: 事务函数
	//
	// 返回:
	//   error: 事务执行失败的错误
	//
	// 示例:
	//   opts := &dbtx.TxOptions{
	//       Isolation: sql.LevelSerializable,
	//       ReadOnly: false,
	//       Timeout: 10 * time.Second,
	//   }
	//   err := manager.WithTxOptions(ctx, opts, func(tx *gorm.DB) error {
	//       return tx.Create(&user).Error
	//   })
	WithTxOptions(ctx context.Context, opts *TxOptions, fn TxFunc) error

	// GetDB 返回底层的数据库实例
	// 用于需要直接访问数据库的场景（非事务操作）
	//
	// 返回:
	//   *gorm.DB: GORM 数据库实例
	GetDB() *gorm.DB
}
