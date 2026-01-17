package repository

import (
	"context"

	"github.com/rei0721/go-scaffold/internal/models"
	"gorm.io/gorm"
)

// AuthRepository 认证数据访问接口
// 提供用户注册、登录等认证相关的数据库操作
type AuthRepository interface {
	// FindUserByUsername 根据用户名查找用户
	// 用于登录验证
	// 参数:
	//   ctx: 上下文
	//   username: 用户名
	// 返回:
	//   *models.User: 用户对象，不存在时返回nil
	//   error: 数据库错误
	FindUserByUsername(ctx context.Context, username string) (*models.User, error)

	// FindUserByEmail 根据邮箱查找用户
	// 用于邮箱登录和注册唯一性验证
	// 参数:
	//   ctx: 上下文
	//   email: 邮箱地址
	// 返回:
	//   *models.User: 用户对象，不存在时返回nil
	//   error: 数据库错误
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)

	// FindUserByID 根据ID查找用户
	// 用于Token刷新等场景
	// 参数:
	//   ctx: 上下文
	//   userID: 用户ID
	// 返回:
	//   *models.User: 用户对象，不存在时返回nil
	//   error: 数据库错误
	FindUserByID(ctx context.Context, userID int64) (*models.User, error)

	// CreateUser 创建新用户（在事务中）
	// 用于用户注册
	// 参数:
	//   ctx: 上下文
	//   tx: GORM事务对象
	//   user: 要创建的用户
	// 返回:
	//   error: 创建失败的错误
	CreateUser(ctx context.Context, tx *gorm.DB, user *models.User) error

	// UpdateUserPassword 更新用户密码（在事务中）
	// 用于修改密码功能
	// 参数:
	//   ctx: 上下文
	//   tx: GORM事务对象
	//   userID: 用户ID
	//   hashedPassword: 新的加密密码
	// 返回:
	//   error: 更新失败的错误
	UpdateUserPassword(ctx context.Context, tx *gorm.DB, userID int64, hashedPassword string) error

	// UpdateUser 更新用户信息（在事务中）
	// 用于更新用户相关信息
	// 参数:
	//   ctx: 上下文
	//   tx: GORM事务对象
	//   user: 要更新的用户对象
	// 返回:
	//   error: 更新失败的错误
	UpdateUser(ctx context.Context, tx *gorm.DB, user *models.User) error
}
