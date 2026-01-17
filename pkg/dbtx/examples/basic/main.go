package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/rei0721/go-scaffold/pkg/dbtx"
	"github.com/rei0721/go-scaffold/pkg/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        int64     `gorm:"primaryKey"`
	Name      string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:100;uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// UserRole 用户角色关联模型
type UserRole struct {
	ID        int64     `gorm:"primaryKey"`
	UserID    int64     `gorm:"not null;index"`
	RoleID    int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func main() {
	fmt.Println("=== DBTx Basic Example ===")

	// 1. 创建数据库连接（使用内存数据库）
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// 2. 自动迁移
	if err := db.AutoMigrate(&User{}, &UserRole{}); err != nil {
		log.Fatal("failed to migrate:", err)
	}
	fmt.Println("✓ Database migrated")

	// 3. 创建日志器（可选）
	loggerInstance, err := logger.New(&logger.Config{
		Level:  "debug",
		Output: "stdout",
		Format: "console",
	})
	if err != nil {
		log.Fatal("failed to create logger:", err)
	}

	// 4. 创建事务管理器
	txManager, err := dbtx.NewManager(db, loggerInstance)
	if err != nil {
		log.Fatal("failed to create tx manager:", err)
	}
	fmt.Println("✓ Transaction Manager created")

	ctx := context.Background()

	// ========== 示例 1: 基本事务提交 ==========
	fmt.Println("\n--- Example 1: Basic Transaction ---")
	err = basicTransaction(ctx, txManager)
	if err != nil {
		log.Fatal("basic transaction failed:", err)
	}
	fmt.Println("✓ Basic transaction committed successfully")

	// ========== 示例 2: 错误自动回滚 ==========
	fmt.Println("\n--- Example 2: Rollback on Error ---")
	err = errorRollback(ctx, txManager)
	if err == nil {
		log.Fatal("expected error but got nil")
	}
	fmt.Println("✓ Rollback on error works correctly")
	fmt.Println("  Error:", err)

	// ========== 示例 3: 嵌套事务 ==========
	fmt.Println("\n--- Example 3: Nested Transaction ---")
	err = nestedTransaction(ctx, txManager)
	if err != nil {
		log.Fatal("nested transaction failed:", err)
	}
	fmt.Println("✓ Nested transaction committed successfully")

	// ========== 示例 4: 自定义选项 ==========
	fmt.Println("\n--- Example 4: Custom Options ---")
	err = customOptions(ctx, txManager)
	if err != nil {
		log.Fatal("custom options transaction failed:", err)
	}
	fmt.Println("✓ Custom options transaction committed successfully")

	fmt.Println("\n=== All Examples Completed ===")
}

// basicTransaction 示例 1: 基本事务
func basicTransaction(ctx context.Context, txManager dbtx.Manager) error {
	return txManager.WithTx(ctx, func(tx *gorm.DB) error {
		// 创建用户
		user := &User{
			Name:  "Alice",
			Email: "alice@example.com",
		}

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		fmt.Printf("  Created user: %s (ID: %d)\n", user.Name, user.ID)
		return nil // 自动提交
	})
}

// errorRollback 示例 2: 错误自动回滚
func errorRollback(ctx context.Context, txManager dbtx.Manager) error {
	return txManager.WithTx(ctx, func(tx *gorm.DB) error {
		// 创建用户
		user := &User{
			Name:  "Bob",
			Email: "alice@example.com", // 重复的邮箱
		}

		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("duplicate email: %w", err)
		}

		return nil
	})
}

// nestedTransaction 示例 3: 嵌套事务
func nestedTransaction(ctx context.Context, txManager dbtx.Manager) error {
	return txManager.WithTx(ctx, func(tx *gorm.DB) error {
		// 外层事务: 创建用户
		user := &User{
			Name:  "Bob",
			Email: "bob@example.com",
		}

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		fmt.Printf("  Outer user: %s (ID: %d)\n", user.Name, user.ID)

		// 内层事务: 分配角色（使用 SavePoint）
		return txManager.WithTx(ctx, func(tx2 *gorm.DB) error {
			role := &UserRole{
				UserID: user.ID,
				RoleID: 1, // 假设 roleID=1 是默认角色
			}

			if err := tx2.Create(role).Error; err != nil {
				return err
			}

			fmt.Printf("  Inner role: UserRole{UserID: %d, RoleID: %d}\n",
				role.UserID, role.RoleID)
			return nil
		})
	})
}

// customOptions 示例 4: 自定义选项
func customOptions(ctx context.Context, txManager dbtx.Manager) error {
	// 配置只读事务，超时 5 秒
	opts := &dbtx.TxOptions{
		ReadOnly: true,
		Timeout:  5 * time.Second,
	}

	return txManager.WithTxOptions(ctx, opts, func(tx *gorm.DB) error {
		// 查询用户（只读）
		var users []User
		if err := tx.Find(&users).Error; err != nil {
			return err
		}

		fmt.Printf("  Found %d users\n", len(users))
		return nil
	})
}

// panicRecovery 示例 5: Panic 自动回滚
func panicRecovery(ctx context.Context, txManager dbtx.Manager) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("panic recovered")
		}
	}()

	return txManager.WithTx(ctx, func(tx *gorm.DB) error {
		// 创建用户
		user := &User{
			Name:  "Charlie",
			Email: "charlie@example.com",
		}

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// 模拟 panic
		panic("something went wrong!")
	})
}
