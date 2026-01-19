---
name: repository-development
description: 在 internal/repository/ 目录下创建数据访问层
---

# Repository 开发规范

## 概述

本 skill 指导在 `internal/repository/` 目录下创建符合项目规范的数据访问层。

## 文件结构

```
internal/repository/
├── repository.go      # 通用 Repository 接口
├── {entity}.go        # 实体特定接口
├── {entity}_impl.go   # 实体特定实现
└── constants.go       # 常量定义
```

## 通用 Repository 接口

项目已提供泛型 `Repository[T any]` 接口：

```go
type Repository[T any] interface {
    Create(ctx context.Context, entity *T) error
    FindByID(ctx context.Context, id int64) (*T, error)
    FindAll(ctx context.Context, page, pageSize int) ([]T, int64, error)
    Update(ctx context.Context, entity *T) error
    Delete(ctx context.Context, id int64) error
}
```

## 开发步骤

### 1. 定义特定接口 `{entity}.go`

```go
package repository

import (
    "context"

    "github.com/rei0721/go-scaffold/internal/models"
    "gorm.io/gorm"
)

// {Entity}Repository 定义 {实体} 特定的数据访问接口
// 扩展通用 Repository 接口
type {Entity}Repository interface {
    Repository[models.{Entity}]

    // 特定查询方法
    FindByUsername(ctx context.Context, username string) (*models.{Entity}, error)
    FindByEmail(ctx context.Context, email string) (*models.{Entity}, error)
    ExistsByUsername(ctx context.Context, username string) (bool, error)

    // 事务支持
    WithTx(tx *gorm.DB) {Entity}Repository
}
```

### 2. 实现接口 `{entity}_impl.go`

```go
package repository

import (
    "context"
    "errors"

    "github.com/rei0721/go-scaffold/internal/models"
    "gorm.io/gorm"
)

// {entity}Repository 是 {Entity}Repository 的具体实现
type {entity}Repository struct {
    db *gorm.DB
}

// New{Entity}Repository 创建新的 Repository 实例
func New{Entity}Repository(db *gorm.DB) {Entity}Repository {
    return &{entity}Repository{db: db}
}

// WithTx 返回使用事务的 Repository
func (r *{entity}Repository) WithTx(tx *gorm.DB) {Entity}Repository {
    return &{entity}Repository{db: tx}
}

// Create 创建新实体
func (r *{entity}Repository) Create(ctx context.Context, entity *models.{Entity}) error {
    return r.db.WithContext(ctx).Create(entity).Error
}

// FindByID 根据 ID 查找
func (r *{entity}Repository) FindByID(ctx context.Context, id int64) (*models.{Entity}, error) {
    var entity models.{Entity}
    err := r.db.WithContext(ctx).First(&entity, id).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil // 不存在返回 nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &entity, nil
}

// FindAll 分页查询
func (r *{entity}Repository) FindAll(ctx context.Context, page, pageSize int) ([]models.{Entity}, int64, error) {
    var entities []models.{Entity}
    var total int64

    // 计算总数
    if err := r.db.WithContext(ctx).Model(&models.{Entity}{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // 分页查询
    offset := (page - 1) * pageSize
    if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&entities).Error; err != nil {
        return nil, 0, err
    }

    return entities, total, nil
}

// Update 更新实体
func (r *{entity}Repository) Update(ctx context.Context, entity *models.{Entity}) error {
    return r.db.WithContext(ctx).Save(entity).Error
}

// Delete 删除实体（软删除）
func (r *{entity}Repository) Delete(ctx context.Context, id int64) error {
    return r.db.WithContext(ctx).Delete(&models.{Entity}{}, id).Error
}

// FindByUsername 根据用户名查找
func (r *{entity}Repository) FindByUsername(ctx context.Context, username string) (*models.{Entity}, error) {
    var entity models.{Entity}
    err := r.db.WithContext(ctx).Where("username = ?", username).First(&entity).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &entity, nil
}

// ExistsByUsername 检查用户名是否存在
func (r *{entity}Repository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
    var count int64
    err := r.db.WithContext(ctx).Model(&models.{Entity}{}).Where("username = ?", username).Count(&count).Error
    return count > 0, err
}
```

## 最佳实践

### 软删除

模型继承 `BaseModel` 自动支持软删除：

```go
type {Entity} struct {
    models.BaseModel
    // 字段...
}
```

### 条件查询

```go
func (r *{entity}Repository) FindByCondition(ctx context.Context, cond *Condition) ([]models.{Entity}, error) {
    query := r.db.WithContext(ctx).Model(&models.{Entity}{})

    if cond.Status != 0 {
        query = query.Where("status = ?", cond.Status)
    }
    if cond.Keyword != "" {
        query = query.Where("name LIKE ?", "%"+cond.Keyword+"%")
    }

    var entities []models.{Entity}
    return entities, query.Find(&entities).Error
}
```

### 预加载关联

```go
func (r *{entity}Repository) FindWithRelations(ctx context.Context, id int64) (*models.{Entity}, error) {
    var entity models.{Entity}
    err := r.db.WithContext(ctx).
        Preload("Roles").
        Preload("Profile").
        First(&entity, id).Error
    return &entity, err
}
```

## 检查清单

- [ ] 接口定义在 `{entity}.go`
- [ ] 实现在 `{entity}_impl.go`
- [ ] 实现 `WithTx` 方法支持事务
- [ ] 所有方法第一个参数是 `context.Context`
- [ ] 使用 `WithContext(ctx)` 传递 context
- [ ] 记录不存在返回 `nil, nil` 而非错误
- [ ] 软删除使用 `BaseModel`
- [ ] 分页查询同时返回总数
