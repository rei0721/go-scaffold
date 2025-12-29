# 数据库指南

本文档详细说明了 Rei0721 项目的数据库设计、模型定义和最佳实践。

## 数据库支持

Rei0721 支持多种数据库，通过 GORM ORM 框架进行统一管理。

### 支持的数据库

| 数据库 | 驱动 | 推荐用途 |
|--------|------|----------|
| PostgreSQL | postgres | 生产环境 |
| MySQL | mysql | 生产环境 |
| SQLite | sqlite | 开发/测试 |

## 数据库配置

### PostgreSQL

```yaml
database:
  driver: postgres
  host: localhost
  port: 5432
  user: postgres
  password: password
  dbname: rei0721
  maxOpenConns: 100
  maxIdleConns: 10
```

### MySQL

```yaml
database:
  driver: mysql
  host: localhost
  port: 3306
  user: root
  password: password
  dbname: rei0721
  maxOpenConns: 100
  maxIdleConns: 10
```

### SQLite

```yaml
database:
  driver: sqlite
  dbname: ./rei0721.db
```

## 数据模型

### BaseModel

所有数据模型都继承自 BaseModel，包含通用字段：

```go
type BaseModel struct {
    ID        int64          `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time      `json:"createdAt"`
    UpdatedAt time.Time      `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}
```

字段说明:
- ID - 主键，使用 Snowflake 算法生成的分布式 ID
- CreatedAt - 创建时间，自动设置
- UpdatedAt - 更新时间，自动更新
- DeletedAt - 软删除时间，用于逻辑删除

### User 模型

```go
type User struct {
    BaseModel
    Username string `gorm:"uniqueIndex;size:50;not null" json:"username"`
    Email    string `gorm:"uniqueIndex;size:100;not null" json:"email"`
    Password string `gorm:"size:255;not null" json:"-"`
    Status   int    `gorm:"default:1" json:"status"`
}

func (User) TableName() string {
    return "users"
}
```

字段说明:
- Username - 用户名，唯一索引，最大 50 字符
- Email - 邮箱，唯一索引，最大 100 字符
- Password - 密码哈希，最大 255 字符，不返回给客户端
- Status - 状态，1 表示活跃，0 表示禁用

## GORM 标签

### 常用标签

| 标签 | 说明 |
|------|------|
| gorm:"primaryKey" | 主键 |
| gorm:"uniqueIndex" | 唯一索引 |
| gorm:"index" | 普通索引 |
| gorm:"not null" | 非空约束 |
| gorm:"default:value" | 默认值 |
| gorm:"size:n" | 字段大小 |
| gorm:"column:name" | 列名 |
| gorm:"type:type" | 数据类型 |

### JSON 标签

| 标签 | 说明 |
|------|------|
| json:"fieldName" | JSON 字段名 |
| json:"-" | 不序列化 |
| json:"fieldName,omitempty" | 空值时省略 |

## 数据库操作

### Repository 接口

```go
type Repository[T any] interface {
    Create(ctx context.Context, entity *T) error
    FindByID(ctx context.Context, id int64) (*T, error)
    FindAll(ctx context.Context, page, pageSize int) ([]T, int64, error)
    Update(ctx context.Context, entity *T) error
    Delete(ctx context.Context, id int64) error
}
```

### 创建记录

```go
user := &models.User{
    Username: "john_doe",
    Email:    "john@example.com",
    Password: hashPassword("password123"),
    Status:   1,
}

if err := repo.Create(ctx, user); err != nil {
    return nil, err
}
```

### 查询记录

```go
// 按 ID 查询
user, err := repo.FindByID(ctx, 123)

// 按用户名查询
user, err := repo.FindByUsername(ctx, "john_doe")

// 列表查询 (分页)
users, total, err := repo.FindAll(ctx, 1, 10)
```

### 更新记录

```go
user.Status = 0
if err := repo.Update(ctx, user); err != nil {
    return err
}
```

### 删除记录

```go
// 软删除
if err := repo.Delete(ctx, 123); err != nil {
    return err
}

// 硬删除 (需要自定义)
if err := db.Unscoped().Delete(&User{}, 123).Error; err != nil {
    return err
}
```

## 高级查询

### 条件查询

```go
var users []models.User

// 单个条件
db.Where("status = ?", 1).Find(&users)

// 多个条件
db.Where("status = ? AND email LIKE ?", 1, "%@example.com").Find(&users)

// IN 查询
db.Where("id IN ?", []int64{1, 2, 3}).Find(&users)
```

### 排序

```go
// 升序
db.Order("created_at ASC").Find(&users)

// 降序
db.Order("created_at DESC").Find(&users)

// 多字段排序
db.Order("status DESC, created_at ASC").Find(&users)
```

### 分页

```go
page := 1
pageSize := 10
offset := (page - 1) * pageSize

var users []models.User
var total int64

db.Model(&models.User{}).Count(&total)
db.Offset(offset).Limit(pageSize).Find(&users)
```

## 事务处理

### 基本事务

```go
tx := db.BeginTx(ctx, nil)

// 创建用户
user := &models.User{...}
if err := tx.Create(user).Error; err != nil {
    tx.Rollback()
    return err
}

// 更新用户
if err := tx.Model(user).Update("status", 1).Error; err != nil {
    tx.Rollback()
    return err
}

// 提交事务
if err := tx.Commit().Error; err != nil {
    return err
}
```

## 性能优化

### 连接池配置

```yaml
database:
  maxOpenConns: 100    # 最大连接数
  maxIdleConns: 10     # 最大空闲连接数
```

建议:
- 开发环境: maxOpenConns=10, maxIdleConns=5
- 生产环境: maxOpenConns=100, maxIdleConns=20

### 查询优化

```go
// 不好 - N+1 查询
var users []models.User
db.Find(&users)
for _, user := range users {
    db.Model(&user).Association("Orders").Find(&user.Orders)
}

// 好 - 预加载
var users []models.User
db.Preload("Orders").Find(&users)

// 好 - 只查询需要的字段
var users []models.User
db.Select("id", "username", "email").Find(&users)
```

## 最佳实践

### 1. 使用 Context

所有数据库操作都应该接收 context：

```go
// 好
func (r *userRepository) FindByID(ctx context.Context, id int64) (*User, error) {
    return r.db.WithContext(ctx).First(&User{}, id).Error
}

// 不好
func (r *userRepository) FindByID(id int64) (*User, error) {
    return r.db.First(&User{}, id).Error
}
```

### 2. 错误处理

```go
// 好
if err := db.Create(user).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, &errors.BizError{Code: errors.ErrUserNotFound}
    }
    return nil, &errors.BizError{Code: errors.ErrDatabaseError}
}

// 不好
if err := db.Create(user).Error; err != nil {
    return nil, err
}
```

### 3. 使用事务

```go
// 好 - 使用事务确保数据一致性
tx := db.BeginTx(ctx, nil)
if err := tx.Create(user).Error; err != nil {
    tx.Rollback()
    return err
}
if err := tx.Commit().Error; err != nil {
    return err
}

// 不好 - 没有事务保护
db.Create(user)
db.Update(user)
```

### 4. 软删除

```go
// 好 - 使用软删除
db.Delete(&user)  // 只设置 DeletedAt

// 查询时自动排除已删除的记录
db.Find(&users)   // 不包括已删除的

// 包括已删除的记录
db.Unscoped().Find(&users)

// 不好 - 硬删除
db.Unscoped().Delete(&user)  // 永久删除
```

---

**最后更新**: 2025-12-30
