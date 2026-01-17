# DBTx - 数据库事务管理器

## 概述

DBTx 是一个优雅的数据库事务管理工具库，基于 GORM 封装，提供函数式 API 简化事务处理代码。

### 特性

- ✅ **函数式 API** - 使用闭包封装事务逻辑，自动管理提交/回滚
- ✅ **自动回滚** - panic 或 error 时自动回滚，成功时自动提交
- ✅ **嵌套事务** - 基于 GORM SavePoint 机制支持嵌套事务
- ✅ **超时控制** - 集成 Context 超时机制
- ✅ **错误包装** - 统一错误处理和日志记录
- ✅ **可选 Logger** - 记录事务开始/提交/回滚事件
- ✅ **线程安全** - 所有操作并发安全

## 安装

DBTx 是项目内部包，无需单独安装。

依赖:

- `gorm.io/gorm` v1.20+
- `github.com/rei0721/go-scaffold/pkg/logger` (可选)

## 快速开始

### 1. 创建事务管理器

```go
package main

import (
    "context"
    "log"

    "github.com/rei0721/go-scaffold/pkg/dbtx"
    "github.com/rei0721/go-scaffold/pkg/logger"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func main() {
    // 创建数据库连接
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // 创建日志器（可选）
    logger, _ := logger.New(&logger.Config{...})

    // 创建事务管理器
    txManager, err := dbtx.NewManager(db, logger)
    if err != nil {
        log.Fatal(err)
    }

    // 使用事务管理器
    ctx := context.Background()
    err = txManager.WithTx(ctx, func(tx *gorm.DB) error {
        // 在这里执行数据库操作
        return tx.Create(&User{Name: "Alice"}).Error
    })
}
```

### 2. 基本事务

```go
// 自动提交
err := txManager.WithTx(ctx, func(tx *gorm.DB) error {
    // 创建用户
    user := &User{Name: "Alice", Email: "alice@example.com"}
    if err := tx.Create(user).Error; err != nil {
        return err // 自动回滚
    }

    // 创建用户配置
    config := &UserConfig{UserID: user.ID, Theme: "dark"}
    if err := tx.Create(config).Error; err != nil {
        return err // 自动回滚
    }

    return nil // 自动提交
})

if err != nil {
    log.Error("transaction failed", "error", err)
}
```

### 3. 使用自定义选项

```go
import (
    "database/sql"
    "time"
)

// 配置事务选项
opts := &dbtx.TxOptions{
    Isolation: sql.LevelSerializable, // 串行化隔离级别
    ReadOnly:  false,                  // 读写事务
    Timeout:   10 * time.Second,       // 10秒超时
}

err := txManager.WithTxOptions(ctx, opts, func(tx *gorm.DB) error {
    return tx.Create(&order).Error
})
```

### 4. 嵌套事务

```go
// 外层事务
err := txManager.WithTx(ctx, func(tx *gorm.DB) error {
    // 创建用户
    user := &User{Name: "Bob"}
    if err := tx.Create(user).Error; err != nil {
        return err
    }

    // 内层事务（使用 SavePoint）
    return txManager.WithTx(ctx, func(tx2 *gorm.DB) error {
        // 分配角色
        role := &UserRole{UserID: user.ID, RoleID: 1}
        return tx2.Create(role).Error
    })
})
```

## 在 Service 层集成

### 示例: 用户注册服务

**重构前** - 手动管理事务:

```go
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) error {
    // 开启事务
    tx := s.db.Begin()
    if tx.Error != nil {
        return tx.Error
    }

    // defer 回滚
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r)
        }
    }()

    // 创建用户
    if err := tx.Create(&user).Error; err != nil {
        tx.Rollback()
        return err
    }

    // 分配角色
    if err := tx.Create(&userRole).Error; err != nil {
        tx.Rollback()
        return err
    }

    // 提交事务
    return tx.Commit().Error
}
```

**重构后** - 使用 dbtx:

```go
type AuthService struct {
    repo      repository.UserRepository
    txManager dbtx.Manager
}

func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) error {
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

        return nil // 自动提交
    })
}
```

**优势对比**:

- ✅ 代码行数减少 50%+
- ✅ 不需要手写 defer/recover
- ✅ 自动错误处理
- ✅ 支持日志记录

## API 参考

### Manager 接口

| 方法                           | 说明                   |
| ------------------------------ | ---------------------- |
| `WithTx(ctx, fn)`              | 使用默认选项执行事务   |
| `WithTxOptions(ctx, opts, fn)` | 使用自定义选项执行事务 |
| `GetDB()`                      | 返回底层数据库实例     |

### TxOptions 配置

| 字段                       | 类型                 | 说明         | 默认值         |
| -------------------------- | -------------------- | ------------ | -------------- |
| `Isolation`                | `sql.IsolationLevel` | 事务隔离级别 | `LevelDefault` |
| `ReadOnly`                 | `bool`               | 是否只读事务 | `false`        |
| `Timeout`                  | `time.Duration`      | 超时时间     | `30s`          |
| `DisableNestedTransaction` | `bool`               | 禁用嵌套事务 | `false`        |

**隔离级别选项**:

- `sql.LevelDefault` - 数据库默认
- `sql.LevelReadUncommitted` - 读未提交
- `sql.LevelReadCommitted` - 读已提交
- `sql.LevelRepeatableRead` - 可重复读
- `sql.LevelSerializable` - 串行化

### 错误定义

| 错误                | 说明           |
| ------------------- | -------------- |
| `ErrNilDB`          | 数据库实例为空 |
| `ErrNilTxFunc`      | 事务函数为空   |
| `ErrInvalidOptions` | 无效的事务选项 |
| `ErrTxNotStarted`   | 事务未开启     |
| `ErrTxRolledBack`   | 事务已回滚     |

## 使用场景

### 场景1: 用户注册（多表操作）

```go
func (s *UserService) Register(ctx context.Context, req *RegisterRequest) error {
    return s.txManager.WithTx(ctx, func(tx *gorm.DB) error {
        // 1. 创建用户
        user := &User{Name: req.Name, Email: req.Email}
        if err := tx.Create(user).Error; err != nil {
            return err
        }

        // 2. 创建用户配置
        config := &UserConfig{UserID: user.ID}
        if err := tx.Create(config).Error; err != nil {
            return err
        }

        // 3. 分配默认角色
        role := &UserRole{UserID: user.ID, RoleID: 1}
        return tx.Create(role).Error
    })
}
```

### 场景2: 订单创建（库存扣减）

```go
func (s *OrderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) error {
    opts := &dbtx.TxOptions{
        Isolation: sql.LevelSerializable, // 防止并发问题
        Timeout:   5 * time.Second,
    }

    return s.txManager.WithTxOptions(ctx, opts, func(tx *gorm.DB) error {
        // 1. 检查库存
        var product Product
        if err := tx.First(&product, req.ProductID).Error; err != nil {
            return err
        }
        if product.Stock < req.Quantity {
            return errors.New("insufficient stock")
        }

        // 2. 扣减库存
        if err := tx.Model(&product).Update("stock", product.Stock-req.Quantity).Error; err != nil {
            return err
        }

        // 3. 创建订单
        order := &Order{ProductID: req.ProductID, Quantity: req.Quantity}
        return tx.Create(order).Error
    })
}
```

### 场景3: 只读事务（一致性查询）

```go
func (s *ReportService) GetMonthlyReport(ctx context.Context) (*Report, error) {
    opts := &dbtx.TxOptions{
        ReadOnly: true, // 只读事务，性能更好
        Timeout:  10 * time.Second,
    }

    var report Report
    err := s.txManager.WithTxOptions(ctx, opts, func(tx *gorm.DB) error {
        // 查询多个表，保证数据一致性
        tx.Find(&report.Orders)
        tx.Find(&report.Users)
        return nil
    })

    return &report, err
}
```

## 最佳实践

### 1. 始终使用 Context 超时

```go
// ✅ 好的做法
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
txManager.WithTx(ctx, func(tx *gorm.DB) error {...})

// ❌ 不好的做法 - 可能导致无限阻塞
txManager.WithTx(context.Background(), func(tx *gorm.DB) error {...})
```

### 2. 保持事务简短

```go
// ❌ 不好的做法 - 在事务中执行耗时操作
txManager.WithTx(ctx, func(tx *gorm.DB) error {
    tx.Create(&user)
    sendEmail(user.Email) // 可能很慢！
    uploadToS3(user.Avatar) // 耗时 I/O！
    return nil
})

// ✅ 好的做法 - 事务后执行耗时操作
txManager.WithTx(ctx, func(tx *gorm.DB) error {
    return tx.Create(&user).Error
})
sendEmail(user.Email)
uploadToS3(user.Avatar)
```

### 3. 显式返回错误

```go
// ❌ 不好的做法 - 忽略错误
txManager.WithTx(ctx, func(tx *gorm.DB) error {
    tx.Create(&user) // 忽略错误
    return nil
})

// ✅ 好的做法 - 显式处理错误
txManager.WithTx(ctx, func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil {
        return err
    }
    return nil
})
```

### 4. 使用只读事务优化

```go
// 只读操作使用 ReadOnly 提升性能
opts := &dbtx.TxOptions{ReadOnly: true}
txManager.WithTxOptions(ctx, opts, func(tx *gorm.DB) error {
    return tx.Find(&users).Error
})
```

### 5. 避免过深的嵌套

```go
// ❌ 不好的做法 - 嵌套过深
txManager.WithTx(ctx, func(tx1 *gorm.DB) error {
    return txManager.WithTx(ctx, func(tx2 *gorm.DB) error {
        return txManager.WithTx(ctx, func(tx3 *gorm.DB) error {
            // 3 层嵌套，性能开销大
            return nil
        })
    })
})

// ✅ 好的做法 - 扁平化设计
txManager.WithTx(ctx, func(tx *gorm.DB) error {
    // 所有操作在同一事务中
    tx.Create(&user)
    tx.Create(&role)
    tx.Create(&config)
    return nil
})
```

## 性能考虑

| 因素     | 影响                       | 建议                  |
| -------- | -------------------------- | --------------------- |
| 事务时长 | 持有锁时间越长，并发性越差 | 保持事务简短（< 1秒） |
| 嵌套事务 | SavePoint 有额外开销       | 减少嵌套层级          |
| 隔离级别 | 级别越高，锁竞争越激烈     | 根据需求选择最低级别  |
| 只读事务 | 性能更好                   | 查询操作使用 ReadOnly |

## 错误处理

```go
err := txManager.WithTx(ctx, func(tx *gorm.DB) error {
    return tx.Create(&user).Error
})

if err != nil {
    // 检查具体错误类型
    if errors.Is(err, dbtx.ErrTxRolledBack) {
        log.Warn("transaction was rolled back")
    } else if strings.Contains(err.Error(), "context deadline exceeded") {
        log.Error("transaction timeout")
    } else {
        log.Error("transaction failed", "error", err)
    }
}
```

## 常见问题

### Q: 嵌套事务如何工作？

A: GORM 使用 SavePoint 机制。内层事务回滚时，只回滚到 SavePoint，外层事务不受影响。

### Q: Context 取消会回滚事务吗？

A: 会。但注意数据库可能已执行部分操作，Context 取消不保证完全原子性。

### Q: 可以在事务中调用其他 Service 吗？

A: 可以，但需要确保被调用的 Service 支持接收 tx 参数，避免嵌套过深。

### Q: ReadOnly 事务能提升多少性能？

A: 取决于数据库和场景。通常可以减少锁竞争，提升 10-30% 性能。

## 项目结构

```
pkg/dbtx/
├── doc.go          # GoDoc 包文档
├── README.md       # 本文档
├── dbtx.go         # Manager 接口定义
├── manager.go      # manager 实现
├── options.go      # TxOptions 定义
├── constants.go    # 常量定义
├── errors.go       # 错误定义
└── examples/
    ├── README.md   # 示例说明
    └── basic/
        └── main.go # 基础示例
```

## 参考链接

- [GORM 事务文档](https://gorm.io/docs/transactions.html)
- [Go database/sql 事务](https://pkg.go.dev/database/sql#Tx)
- [项目架构地图](../../docs/architecture/system_map.md)
