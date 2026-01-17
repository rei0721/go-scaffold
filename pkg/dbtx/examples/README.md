# DBTx 示例

本目录包含 `pkg/dbtx` 的可运行示例。

## 运行示例

### 基础示例

```bash
cd examples/basic
go run main.go
```

**示例输出**：

```
=== DBTx Basic Example ===
✓ Transaction Manager created
✓ Database migrated
✓ Basic transaction committed successfully
  Created user: Alice (ID: 1)
✓ Rollback on error works correctly
  Error: duplicate email
✓ Nested transaction committed successfully
  Outer user: Bob (ID: 2)
  Inner role: UserRole{UserID: 2, RoleID: 1}
✓ All tests passed!
```

## 示例说明

### basic 示例

演示 dbtx 的核心功能:

- 创建事务管理器
- 基本事务提交
- 错误自动回滚
- 嵌套事务（SavePoint）
- Logger 集成

### 代码结构

```go
// 1. 初始化
db := setupDatabase()
logger := setupLogger()
txManager := dbtx.NewManager(db, logger)

// 2. 基本事务
txManager.WithTx(ctx, func(tx *gorm.DB) error {
    return tx.Create(&user).Error
})

// 3. 错误回滚
txManager.WithTx(ctx, func(tx *gorm.DB) error {
    return errors.New("simulate error") // 自动回滚
})

// 4. 嵌套事务
txManager.WithTx(ctx, func(tx *gorm.DB) error {
    tx.Create(&user)
    return txManager.WithTx(ctx, func(tx2 *gorm.DB) error {
        return tx2.Create(&role).Error
    })
})
```

## 依赖

- `gorm.io/gorm`
- `gorm.io/driver/sqlite`
- `github.com/rei0721/go-scaffold/pkg/logger`
