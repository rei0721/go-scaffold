---
name: model-development
description: 在 internal/models/ 目录下创建数据模型
---

# 数据模型开发规范

## 概述

本 skill 指导在 `internal/models/` 目录下创建符合项目规范的 GORM 数据模型。

## 文件结构

```
internal/models/
├── db_base.go       # 基础模型
├── db_{entity}.go   # 实体模型
└── constants.go     # 常量定义
```

## BaseModel 基础模型

项目提供 `BaseModel` 作为所有模型的基类：

```go
// BaseModel 是所有模型的基类
// 提供通用字段：ID、创建时间、更新时间、软删除
type BaseModel struct {
    ID        int64          `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

## 开发步骤

### 1. 创建模型文件 `db_{entity}.go`

```go
package models

import (
    "time"
)

// {Entity} 表示 {实体描述}
type {Entity} struct {
    BaseModel

    // 基本字段
    Name        string `gorm:"size:100;not null" json:"name"`
    Description string `gorm:"size:500" json:"description"`

    // 唯一索引字段
    Code string `gorm:"uniqueIndex;size:50;not null" json:"code"`

    // 外键字段
    UserID int64 `gorm:"index;not null" json:"user_id"`

    // 状态字段
    Status int `gorm:"default:1" json:"status"`

    // 关联关系（使用指针避免循环依赖）
    User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`

    // 多对多关系
    Tags []Tag `gorm:"many2many:{entity}_tags;" json:"tags,omitempty"`
}

// TableName 返回表名
func ({Entity}) TableName() string {
    return "{entities}" // 复数形式
}
```

### 2. GORM 标签规范

| 标签                 | 说明     | 示例             |
| -------------------- | -------- | ---------------- |
| `gorm:"primaryKey"`  | 主键     | `ID int64`       |
| `gorm:"uniqueIndex"` | 唯一索引 | `Email string`   |
| `gorm:"index"`       | 普通索引 | `UserID int64`   |
| `gorm:"size:N"`      | 字段大小 | `Name string`    |
| `gorm:"not null"`    | 非空约束 | 必填字段         |
| `gorm:"default:X"`   | 默认值   | `Status int`     |
| `gorm:"type:X"`      | 指定类型 | `Content string` |

### 3. 关联关系

#### 一对多

```go
// User 拥有多个 Post
type User struct {
    BaseModel
    Posts []Post `gorm:"foreignKey:UserID"`
}

type Post struct {
    BaseModel
    UserID int64 `gorm:"index"`
    User   *User `gorm:"foreignKey:UserID"`
}
```

#### 多对多

```go
type User struct {
    BaseModel
    Roles []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
    BaseModel
    Users []User `gorm:"many2many:user_roles;"`
}
```

### 4. JSON 标签规范

```go
type {Entity} struct {
    // 正常输出
    ID   int64  `json:"id"`
    Name string `json:"name"`

    // 忽略输出（敏感字段）
    Password string `json:"-"`

    // 条件输出（空值不输出）
    Profile *Profile `json:"profile,omitempty"`

    // 字符串形式输出大整数
    BigID int64 `json:"big_id,string"`
}
```

### 5. 状态常量定义

在 `constants.go` 中：

```go
package models

// {Entity} 状态常量
const (
    {Entity}StatusInactive = 0
    {Entity}StatusActive   = 1
    {Entity}StatusDeleted  = 2
)

// {Entity}StatusText 状态文本映射
var {Entity}StatusText = map[int]string{
    {Entity}StatusInactive: "inactive",
    {Entity}StatusActive:   "active",
    {Entity}StatusDeleted:  "deleted",
}
```

### 6. 模型方法

```go
// IsActive 检查是否激活
func (e *{Entity}) IsActive() bool {
    return e.Status == {Entity}StatusActive
}

// BeforeCreate GORM 钩子：创建前
func (e *{Entity}) BeforeCreate(tx *gorm.DB) error {
    // ID 由 Snowflake 生成，在 Repository 层处理
    return nil
}

// AfterFind GORM 钩子：查询后
func (e *{Entity}) AfterFind(tx *gorm.DB) error {
    // 数据后处理
    return nil
}
```

## 数据库迁移

在 `internal/app/app_initdb.go` 中注册模型：

```go
func (a *App) autoMigrate() error {
    return a.DB.DB().AutoMigrate(
        &models.User{},
        &models.{Entity}{},
        // 添加新模型
    )
}
```

## 检查清单

- [ ] 继承 `BaseModel`
- [ ] 实现 `TableName()` 方法
- [ ] 字段有合适的 GORM 标签
- [ ] 字段有合适的 JSON 标签
- [ ] 敏感字段使用 `json:"-"`
- [ ] 外键字段有 `index` 标签
- [ ] 唯一字段有 `uniqueIndex` 标签
- [ ] 状态常量定义在 `constants.go`
- [ ] 在 `autoMigrate` 中注册
