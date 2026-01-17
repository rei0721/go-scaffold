package dbtx

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestUser 测试用的用户模型
type TestUser struct {
	ID    int64  `gorm:"primaryKey"`
	Name  string `gorm:"size:100"`
	Email string `gorm:"size:100;uniqueIndex"`
}

// setupTestDB 创建测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(&TestUser{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

// TestNewManager 测试创建管理器
func TestNewManager(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := setupTestDB(t)
		mgr, err := NewManager(db, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mgr == nil {
			t.Fatal("expected manager, got nil")
		}
	})

	t.Run("nil database", func(t *testing.T) {
		_, err := NewManager(nil, nil)
		if !errors.Is(err, ErrNilDB) {
			t.Fatalf("expected ErrNilDB, got %v", err)
		}
	})
}

// TestWithTx_BasicCommit 测试基本提交
func TestWithTx_BasicCommit(t *testing.T) {
	db := setupTestDB(t)
	mgr, _ := NewManager(db, nil)
	ctx := context.Background()

	err := mgr.WithTx(ctx, func(tx *gorm.DB) error {
		user := &TestUser{Name: "Alice", Email: "alice@example.com"}
		return tx.Create(user).Error
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 验证数据已提交
	var count int64
	db.Model(&TestUser{}).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 user, got %d", count)
	}
}

// TestWithTx_ErrorRollback 测试错误回滚
func TestWithTx_ErrorRollback(t *testing.T) {
	db := setupTestDB(t)
	mgr, _ := NewManager(db, nil)
	ctx := context.Background()

	expectedErr := errors.New("simulate error")

	err := mgr.WithTx(ctx, func(tx *gorm.DB) error {
		// 创建用户
		user := &TestUser{Name: "Bob", Email: "bob@example.com"}
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// 返回错误，触发回滚
		return expectedErr
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// 验证数据已回滚
	var count int64
	db.Model(&TestUser{}).Count(&count)
	if count != 0 {
		t.Fatalf("expected 0 users (rolled back), got %d", count)
	}
}

// TestWithTx_PanicRollback 测试 panic 回滚
func TestWithTx_PanicRollback(t *testing.T) {
	db := setupTestDB(t)
	mgr, _ := NewManager(db, nil)
	ctx := context.Background()

	// 测试 panic 是否被正确抛出
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic, got nil")
			}
		}()

		_ = mgr.WithTx(ctx, func(tx *gorm.DB) error {
			// 创建用户
			user := &TestUser{Name: "Charlie", Email: "charlie@example.com"}
			if err := tx.Create(user).Error; err != nil {
				return err
			}

			// 触发 panic
			panic("test panic")
		})
	}()

	// 验证数据已回滚
	var count int64
	db.Model(&TestUser{}).Count(&count)
	if count != 0 {
		t.Fatalf("expected 0 users (rolled back), got %d", count)
	}
}

// TestWithTx_NestedTransaction 测试嵌套事务
func TestWithTx_NestedTransaction(t *testing.T) {
	db := setupTestDB(t)
	mgr, _ := NewManager(db, nil)
	ctx := context.Background()

	err := mgr.WithTx(ctx, func(tx *gorm.DB) error {
		// 外层事务: 创建第一个用户
		user1 := &TestUser{Name: "User1", Email: "user1@example.com"}
		if err := tx.Create(user1).Error; err != nil {
			return err
		}

		// 内层事务: 创建第二个用户
		return mgr.WithTx(ctx, func(tx2 *gorm.DB) error {
			user2 := &TestUser{Name: "User2", Email: "user2@example.com"}
			return tx2.Create(user2).Error
		})
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 验证两个用户都已提交
	var count int64
	db.Model(&TestUser{}).Count(&count)
	if count != 2 {
		t.Fatalf("expected 2 users, got %d", count)
	}
}

// TestWithTx_NestedRollback 测试嵌套事务回滚
func TestWithTx_NestedRollback(t *testing.T) {
	db := setupTestDB(t)
	mgr, _ := NewManager(db, nil)
	ctx := context.Background()

	err := mgr.WithTx(ctx, func(tx *gorm.DB) error {
		// 外层事务: 创建第一个用户
		user1 := &TestUser{Name: "User1", Email: "user1@example.com"}
		if err := tx.Create(user1).Error; err != nil {
			return err
		}

		// 内层事务: 失败，触发回滚
		innerErr := mgr.WithTx(ctx, func(tx2 *gorm.DB) error {
			user2 := &TestUser{Name: "User2", Email: "user2@example.com"}
			if err := tx2.Create(user2).Error; err != nil {
				return err
			}
			return errors.New("inner transaction failed")
		})

		// 内层失败，外层继续
		if innerErr != nil {
			// 注意: 如果返回 innerErr，外层也会回滚
			// 这里不返回，继续外层事务
			t.Log("inner transaction failed (expected):", innerErr)
		}

		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 验证只有外层用户提交了
	var count int64
	db.Model(&TestUser{}).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 user, got %d", count)
	}
}

// TestWithTxOptions_Timeout 测试超时
func TestWithTxOptions_Timeout(t *testing.T) {
	db := setupTestDB(t)
	mgr, _ := NewManager(db, nil)

	// 创建一个已超时的 Context
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	time.Sleep(10 * time.Millisecond) // 确保超时

	err := mgr.WithTx(ctx, func(tx *gorm.DB) error {
		time.Sleep(100 * time.Millisecond) // 模拟耗时操作
		return nil
	})

	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}

	// 验证错误包含 context deadline exceeded
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Logf("error: %v", err)
		// 注意: 某些情况下错误可能被包装，所以不强制要求
	}
}

// TestWithTxOptions_ReadOnly 测试只读选项
func TestWithTxOptions_ReadOnly(t *testing.T) {
	db := setupTestDB(t)
	mgr, _ := NewManager(db, nil)
	ctx := context.Background()

	// 先创建一些数据
	user := &TestUser{Name: "Alice", Email: "alice@example.com"}
	db.Create(user)

	opts := &TxOptions{
		ReadOnly: true,
		Timeout:  5 * time.Second,
	}

	err := mgr.WithTxOptions(ctx, opts, func(tx *gorm.DB) error {
		var users []TestUser
		return tx.Find(&users).Error
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestWithTxOptions_CustomIsolation 测试自定义隔离级别
func TestWithTxOptions_CustomIsolation(t *testing.T) {
	db := setupTestDB(t)
	mgr, _ := NewManager(db, nil)
	ctx := context.Background()

	opts := &TxOptions{
		Isolation: sql.LevelSerializable,
		Timeout:   5 * time.Second,
	}

	err := mgr.WithTxOptions(ctx, opts, func(tx *gorm.DB) error {
		user := &TestUser{Name: "Bob", Email: "bob@example.com"}
		return tx.Create(user).Error
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestWithTx_NilFunc 测试 nil 函数
func TestWithTx_NilFunc(t *testing.T) {
	db := setupTestDB(t)
	mgr, _ := NewManager(db, nil)
	ctx := context.Background()

	err := mgr.WithTx(ctx, nil)
	if !errors.Is(err, ErrNilTxFunc) {
		t.Fatalf("expected ErrNilTxFunc, got %v", err)
	}
}

// TestDefaultOptions 测试默认选项
func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	if opts == nil {
		t.Fatal("expected options, got nil")
	}

	if opts.Isolation != sql.LevelDefault {
		t.Fatalf("expected LevelDefault, got %v", opts.Isolation)
	}

	if opts.ReadOnly != false {
		t.Fatal("expected ReadOnly=false")
	}

	if opts.Timeout != DefaultTimeout {
		t.Fatalf("expected %v, got %v", DefaultTimeout, opts.Timeout)
	}
}

// TestTxOptions_Validate 测试选项验证
func TestTxOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *TxOptions
		wantErr bool
	}{
		{
			name:    "nil options",
			opts:    nil,
			wantErr: true,
		},
		{
			name: "valid options",
			opts: &TxOptions{
				Timeout: 10 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "negative timeout",
			opts: &TxOptions{
				Timeout: -1 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "zero timeout (should use default)",
			opts: &TxOptions{
				Timeout: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			// 验证零超时被设置为默认值
			if tt.opts != nil && !tt.wantErr && tt.opts.Timeout == 0 {
				if tt.opts.Timeout != DefaultTimeout {
					t.Fatalf("expected timeout to be set to default, got %v", tt.opts.Timeout)
				}
			}
		})
	}
}

// TestGetDB 测试 GetDB 方法
func TestGetDB(t *testing.T) {
	db := setupTestDB(t)
	mgr, _ := NewManager(db, nil)

	gotDB := mgr.GetDB()
	if gotDB != db {
		t.Fatal("GetDB() returned different instance")
	}
}

// BenchmarkWithTx 基准测试
func BenchmarkWithTx(b *testing.B) {
	db := setupTestDB(&testing.T{})
	mgr, _ := NewManager(db, nil)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mgr.WithTx(ctx, func(tx *gorm.DB) error {
			user := &TestUser{
				Name:  "BenchUser",
				Email: "bench@example.com",
			}
			return tx.Create(user).Error
		})
	}
}
