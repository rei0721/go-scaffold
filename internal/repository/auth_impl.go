package repository

import (
	"context"
	"errors"

	"github.com/rei0721/go-scaffold/internal/models"
	"gorm.io/gorm"
)

// authRepository 认证数据访问实现
// 使用GORM实现AuthRepository接口
type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository 创建AuthRepository实例
// 参数:
//
//	db: GORM数据库连接
//
// 返回:
//
//	AuthRepository: 认证数据访问接口
func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

// FindUserByUsername 根据用户名查找用户
func (r *authRepository) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 用户不存在，返回nil而非错误
		}
		return nil, err
	}
	return &user, nil
}

// FindUserByEmail 根据邮箱查找用户
func (r *authRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 用户不存在，返回nil而非错误
		}
		return nil, err
	}
	return &user, nil
}

// FindUserByID 根据ID查找用户
func (r *authRepository) FindUserByID(ctx context.Context, userID int64) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 用户不存在，返回nil而非错误
		}
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建新用户（在事务中）
func (r *authRepository) CreateUser(ctx context.Context, tx *gorm.DB, user *models.User) error {
	return tx.WithContext(ctx).Create(user).Error
}

// UpdateUserPassword 更新用户密码（在事务中）
func (r *authRepository) UpdateUserPassword(ctx context.Context, tx *gorm.DB, userID int64, hashedPassword string) error {
	return tx.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("password", hashedPassword).
		Error
}

// UpdateUser 更新用户信息（在事务中）
func (r *authRepository) UpdateUser(ctx context.Context, tx *gorm.DB, user *models.User) error {
	return tx.WithContext(ctx).Save(user).Error
}
