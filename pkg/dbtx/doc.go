/*
Package dbtx 提供优雅的数据库事务管理

# 设计目标

- **简化事务代码** - 使用函数式 API 封装事务逻辑，减少样板代码
- **自动提交/回滚** - 根据函数返回值自动决定提交或回滚
- **Panic 安全** - 捕获 panic 并自动回滚事务
- **嵌套事务支持** - 基于 GORM SavePoint 机制
- **超时控制** - 集成 Context 超时机制
- **日志集成** - 可选的日志记录，便于调试

# 核心概念

## 函数式事务

传统的事务代码需要手动管理 Begin/Commit/Rollback，代码冗长且容易出错：

	tx := db.Begin()
	if tx.Error != nil {
	    return tx.Error
	}
	defer func() {
	    if r := recover(); r != nil {
	        tx.Rollback()
	        panic(r)
	    }
	}()

	if err := doSomething(tx); err != nil {
	    tx.Rollback()
	    return err
	}

	return tx.Commit().Error

使用 dbtx，将事务逻辑封装在闭包中，自动处理提交和回滚：

	err := txManager.WithTx(ctx, func(tx *gorm.DB) error {
	    return doSomething(tx) // 自动提交或回滚
	})

## 嵌套事务

GORM 使用 SavePoint 机制支持嵌套事务。当内层事务回滚时，只回滚到 SavePoint，外层事务不受影响：

	txManager.WithTx(ctx, func(tx *gorm.DB) error {
	    // 外层事务：创建用户
	    db.Create(&user)

	    // 内层事务：分配角色（使用 SavePoint）
	    return txManager.WithTx(ctx, func(tx2 *gorm.DB) error {
	        return db.Create(&userRole)
	    })
	})

# 使用示例

## 基本用法

创建事务管理器：

	import (
	    "github.com/rei0721/go-scaffold/pkg/dbtx"
	    "gorm.io/gorm"
	)

	// 创建管理器
	db, _ := gorm.Open(...)
	txManager, err := dbtx.NewManager(db, logger)
	if err != nil {
	    log.Fatal(err)
	}

执行事务：

	err := txManager.WithTx(ctx, func(tx *gorm.DB) error {
	    // 创建用户
	    if err := tx.Create(&user).Error; err != nil {
	        return err // 自动回滚
	    }

	    // 创建用户配置
	    if err := tx.Create(&userConfig).Error; err != nil {
	        return err // 自动回滚
	    }

	    return nil // 自动提交
	})

## 高级用法

使用自定义选项：

	opts := &dbtx.TxOptions{
	    Isolation: sql.LevelSerializable, // 串行化隔离级别
	    ReadOnly:  false,                  // 读写事务
	    Timeout:   10 * time.Second,       // 10秒超时
	}

	err := txManager.WithTxOptions(ctx, opts, func(tx *gorm.DB) error {
	    return tx.Create(&order).Error
	})

在 Service 层集成：

	type UserService struct {
	    repo      repository.UserRepository
	    txManager dbtx.Manager
	}

	func (s *UserService) Register(ctx context.Context, req *RegisterRequest) error {
	    return s.txManager.WithTx(ctx, func(tx *gorm.DB) error {
	        // 1. 创建用户
	        user := &models.User{...}
	        if err := s.repo.CreateWithTx(ctx, tx, user); err != nil {
	            return err
	        }

	        // 2. 分配默认角色
	        role := &models.UserRole{...}
	        if err := s.repo.AssignRoleWithTx(ctx, tx, role); err != nil {
	            return err
	        }

	        return nil
	    })
	}

## 错误处理

捕获特定错误：

	err := txManager.WithTx(ctx, func(tx *gorm.DB) error {
	    return tx.Create(&user).Error
	})

	if err != nil {
	    if errors.Is(err, dbtx.ErrTxRolledBack) {
	        log.Warn("transaction was rolled back")
	    } else {
	        log.Error("transaction failed", "error", err)
	    }
	}

# 最佳实践

1. **始终使用 Context** - 传递带超时的 Context，避免长事务

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	txManager.WithTx(ctx, func(tx *gorm.DB) error {...})

2. **保持事务简短** - 事务中只执行必要的数据库操作，避免长时间I/O

	// ❌ 不好的做法 - 在事务中发送邮件
	txManager.WithTx(ctx, func(tx *gorm.DB) error {
	    tx.Create(&user)
	    sendEmail(user.Email) // 可能很慢
	    return nil
	})

	// ✅ 好的做法 - 事务后发送邮件
	txManager.WithTx(ctx, func(tx *gorm.DB) error {
	    return tx.Create(&user).Error
	})
	sendEmail(user.Email)

3. **显式返回错误** - 不要忽略事务函数中的错误

	// ❌ 不好的做法
	txManager.WithTx(ctx, func(tx *gorm.DB) error {
	    tx.Create(&user) // 忽略错误
	    return nil
	})

	// ✅ 好的做法
	txManager.WithTx(ctx, func(tx *gorm.DB) error {
	    return tx.Create(&user).Error
	})

4. **避免过深的嵌套** - 嵌套事务有性能开销，尽量扁平化

5. **使用只读事务优化** - 只读操作使用 ReadOnly 选项提升性能

	opts := &dbtx.TxOptions{ReadOnly: true}
	txManager.WithTxOptions(ctx, opts, func(tx *gorm.DB) error {
	    return tx.Find(&users).Error
	})

# 线程安全

所有方法都是线程安全的。可以在多个 goroutine 中并发调用 WithTx。

# 与其他包的区别

- pkg/database: 提供数据库连接管理，dbtx 依赖 database 包
- gorm.DB.Transaction: GORM 原生事务方法，dbtx 提供更优雅的封装和扩展功能
- Repository 层: 负责数据访问，dbtx 在 Service 层协调多个 Repository 的事务

# 性能考虑

- **嵌套事务开销**: SavePoint 有额外开销，避免过度嵌套
- **锁竞争**: 长事务会持有锁，影响并发性能
- **超时设置**: 根据业务复杂度设置合理的超时时间

# 限制

- 依赖 GORM v1.20+ 的 SavePoint 功能
- 不同数据库对隔离级别的支持程度不同
- Context 取消不保证原子性（数据库可能已部分执行）
*/
package dbtx
