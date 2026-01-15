// Package models 定义 RBAC 相关的数据模型
package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 包含所有 RBAC 模型的公共字段
// 内嵌此结构体确保所有表都有统一的基础字段
type BaseModel struct {
	// ID 主键，使用 Snowflake 算法生成的分布式唯一 ID
	ID int64 `gorm:"primaryKey" json:"id"`

	// CreatedAt 记录创建时间，GORM 自动设置
	CreatedAt time.Time `json:"createdAt"`

	// UpdatedAt 记录最后更新时间，GORM 自动更新
	UpdatedAt time.Time `json:"updatedAt"`

	// DeletedAt 软删除时间，实现软删除功能
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// Role 角色模型
// 角色是权限的集合，如 "管理员"、"编辑者"、"访客"
type Role struct {
	BaseModel
	Name        string       `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Description string       `gorm:"size:255" json:"description"`
	Status      int          `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions,omitempty"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// Permission 权限模型
// 权限定义对特定资源的操作权限，格式为 "resource:action"
type Permission struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;size:100;not null" json:"name"` // 权限名称，例如 "user:read"
	Resource    string `gorm:"size:100;not null" json:"resource"`         // 资源标识，如 "users", "roles"
	Action      string `gorm:"size:50;not null" json:"action"`            // 操作类型，如 "read", "write", "delete"
	Description string `gorm:"size:255" json:"description"`
	Status      int    `gorm:"default:1" json:"status"` // 1: 启用, 0: 禁用
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// UserRole 用户-角色关联表
// 多对多关系：一个用户可以拥有多个角色
type UserRole struct {
	UserID int64 `gorm:"primaryKey"`
	RoleID int64 `gorm:"primaryKey"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}

// RolePermission 角色-权限关联表
// 多对多关系：一个角色包含多个权限
type RolePermission struct {
	RoleID       int64 `gorm:"primaryKey"`
	PermissionID int64 `gorm:"primaryKey"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}
