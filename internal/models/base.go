// Package models 定义了应用程序的 GORM 兼容数据模型
// 这些模型用于与数据库交互,提供了统一的数据结构定义
package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 包含所有模型的公共字段
// 通过嵌入此结构体,可以确保所有数据表都有统一的基础字段
// 这遵循了 DRY(Don't Repeat Yourself)原则
type BaseModel struct {
	// ID 使用 Snowflake 算法生成的分布式唯一 ID
	// 优点:
	// - 分布式环境下也能保证唯一性
	// - 按时间递增,有利于数据库索引性能
	// - 不依赖数据库自增,可以在应用层生成
	// gorm:"primaryKey" 标记此字段为主键
	// json:"id" 指定 JSON 序列化时的字段名
	ID int64 `gorm:"primaryKey" json:"id"`

	// CreatedAt 记录创建时间
	// GORM 会在插入记录时自动设置此字段
	// json:"createdAt" 使用驼峰命名,符合前端 JavaScript 习惯
	CreatedAt time.Time `json:"createdAt"`

	// UpdatedAt 记录最后更新时间
	// GORM 会在每次更新记录时自动更新此字段
	// 这对于追踪数据变更历史非常有用
	UpdatedAt time.Time `json:"updatedAt"`

	// DeletedAt 实现软删除功能
	// gorm.DeletedAt 是 GORM v2 的软删除类型
	// 优点:
	// - 删除的数据不会真正从数据库移除,可以恢复
	// - 查询时默认会排除已删除的记录
	// - 有利于数据审计和问题排查
	// gorm:"index" 为此字段创建索引,提高查询性能
	// json:"deletedAt,omitempty" 如果为空则不包含在 JSON 中
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}
