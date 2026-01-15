package models

import (
	rbacmodels "github.com/rei0721/rei0721/pkg/rbac/models"
)

// User 表示系统中的用户实体
// 它嵌入了 BaseModel 以继承公共字段,并添加了用户特定的字段
type User struct {
	// 嵌入 BaseModel,继承 ID、CreatedAt、UpdatedAt、DeletedAt 字段
	// 这是 Go 语言的组合模式,比继承更灵活
	BaseModel

	// Username 用户名
	// gorm:"uniqueIndex" 创建唯一索引,确保用户名不重复
	// size:50 限制字段长度为50个字符,防止过长的输入
	// not null 设置为必填字段,确保数据完整性
	// json:"username" JSON序列化时的字段名
	Username string `gorm:"uniqueIndex;size:50;not null" json:"username"`

	// Email 邮箱地址
	// gorm:"uniqueIndex" 确保邮箱唯一,用于登录和找回密码
	// size:100 邮箱地址的合理最大长度
	// not null 邮箱为必填项
	Email string `gorm:"uniqueIndex;size:100;not null" json:"email"`

	// Password 加密后的密码
	// size:255 bcrypt 等加密算法生成的哈希值长度
	// not null 密码为必填项
	// json:"-" 非常重要的安全措施:在 JSON 序列化时忽略此字段
	// 这确保密码哈希永远不会被返回给客户端
	Password string `gorm:"size:255;not null" json:"-"`

	// Status 用户状态
	// 1: 激活(active) - 用户可以正常登录使用
	// 0: 未激活(inactive) - 用户被禁用或未完成激活流程
	// gorm:"default:1" 默认为激活状态,新用户注册后即可使用
	Status int `gorm:"default:1" json:"status"`

	// Roles 用户拥有的角色列表
	// many2many:user_roles 指定多对多关联的中间表表名
	// 使用 pkg/rbac/models 中的 Role 类型
	Roles []rbacmodels.Role `gorm:"many2many:user_roles" json:"roles,omitempty"`
}

// TableName 返回 User 模型对应的数据库表名
// 显式指定表名 "users",而不是使用 GORM 的默认命名规则
// 这样做的好处:
// - 表名更清晰,符合数据库命名规范
// - 避免不同数据库方言的命名差异
// - 便于与现有数据库集成
func (User) TableName() string {
	return "users"
}
