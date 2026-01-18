package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rei0721/go-scaffold/internal/models"
	"github.com/rei0721/go-scaffold/internal/repository"
	"github.com/rei0721/go-scaffold/internal/service"
	"github.com/rei0721/go-scaffold/types"
	"github.com/rei0721/go-scaffold/types/constants"
	"github.com/rei0721/go-scaffold/types/errors"
	"gorm.io/gorm"
)

// authService 实现 AuthService 接口
// 提供完整的认证服务功能
type authService struct {
	service.BaseService[repository.AuthRepository]
}

// NewAuthService 创建一个新的 AuthService 实例
// 参数:
//
//	repo: 认证仓库实例
//
// 返回:
//
//	AuthService: 认证服务接口
//
// 注意:
//
//	其他依赖通过 SetXxx 等方法延迟注入
func NewAuthService(repo repository.AuthRepository) AuthService {
	s := &authService{}
	s.SetRepository(repo)
	return s
}

// Register 用户注册
// 支持事务：同时创建用户和分配默认角色
func (s *authService) Register(ctx context.Context, req *types.RegisterRequest) (*types.UserResponse, error) {
	// 1. 检查用户名是否已存在
	existingUser, err := s.Repo.FindUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.NewBizError(errors.ErrDatabaseError, "failed to check username").WithCause(err)
	}
	if existingUser != nil {
		return nil, errors.NewBizError(errors.ErrDuplicateUsername, "username already exists")
	}

	// 2. 检查邮箱是否已存在
	existingUser, err = s.Repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.NewBizError(errors.ErrDatabaseError, "failed to check email").WithCause(err)
	}
	if existingUser != nil {
		return nil, errors.NewBizError(errors.ErrDuplicateEmail, "email already exists")
	}

	// 3. 加密密码
	hashedPassword, err := s.Crypto.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewBizError(errors.ErrInternalServer, "failed to hash password").WithCause(err)
	}

	// 4. 创建用户对象
	user := &models.DBUser{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Status:   1, // 默认激活
	}

	// 5. 使用事务管理器执行事务
	txManager := s.GetTxManager()
	if txManager == nil {
		// 降级处理：如果未注入txManager，使用传统方式
		if log := s.GetLogger(); log != nil {
			log.Warn("TxManager not injected, falling back to traditional transaction handling")
		}
		return s.registerWithoutTxManager(ctx, user)
	}

	// 使用dbtx执行事务
	err = txManager.WithTx(ctx, func(tx *gorm.DB) error {
		// 6. 在事务中创建用户
		if err := s.Repo.CreateUser(ctx, tx, user); err != nil {
			return errors.NewBizError(errors.ErrDatabaseError, "failed to create user").WithCause(err)
		}

		// 7. 分配默认角色（如果启用了 RBAC）
		if rbacManager := s.GetRBAC(); rbacManager != nil {
			// 注意：这里假设存在一个默认角色，实际应该根据业务需求配置
			// 例如：分配 "user" 角色
			// 这部分需要 RBAC 服务支持通过角色名查找角色ID的方法
			// 此处留作示例，实际使用时需要完善
			if log := s.GetLogger(); log != nil {
				log.Info("RBAC is enabled, but default role assignment is not implemented yet", "userId", user.ID)
			}
		}

		return nil // 成功，自动提交
	})

	if err != nil {
		return nil, err
	}

	// 8. 记录注册成功
	if log := s.GetLogger(); log != nil {
		log.Info("user registered successfully", "userId", user.ID, "username", user.Username)
	}

	// 9. 异步预热缓存
	if c := s.GetCache(); c != nil {
		if exec := s.GetExecutor(); exec != nil {
			userCopy := *user
			_ = exec.Execute(constants.AppPoolCache, func() {
				key := fmt.Sprintf("user:%d", userCopy.ID)
				if data, err := json.Marshal(userCopy); err == nil {
					_ = c.Set(context.Background(), key, string(data), 1*time.Hour)
				}
			})
		}
	}

	// 10. 返回用户信息
	return toUserResponse(user), nil
}

// registerWithoutTxManager 降级处理：不使用txManager的传统事务方式
func (s *authService) registerWithoutTxManager(ctx context.Context, user *models.DBUser) (*types.UserResponse, error) {
	// 开启事务
	tx := s.DB.DB().Begin()
	if tx.Error != nil {
		return nil, errors.NewBizError(errors.ErrDatabaseError, "failed to begin transaction").WithCause(tx.Error)
	}

	// 确保事务会被回滚或提交
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // 重新抛出 panic
		}
	}()

	// 在事务中创建用户
	if err := s.Repo.CreateUser(ctx, tx, user); err != nil {
		tx.Rollback()
		return nil, errors.NewBizError(errors.ErrDatabaseError, "failed to create user").WithCause(err)
	}

	// 分配默认角色
	if rbacManager := s.GetRBAC(); rbacManager != nil {
		if log := s.GetLogger(); log != nil {
			log.Info("RBAC is enabled, but default role assignment is not implemented yet", "userId", user.ID)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, errors.NewBizError(errors.ErrDatabaseError, "failed to commit transaction").WithCause(err)
	}

	return toUserResponse(user), nil
}

// Login 用户登录
func (s *authService) Login(ctx context.Context, req *types.LoginRequest) (*types.LoginResponse, error) {
	// 1. 根据用户名查找用户
	user, err := s.Repo.FindUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.NewBizError(errors.ErrDatabaseError, "failed to find user").WithCause(err)
	}
	if user == nil {
		return nil, errors.NewBizError(errors.ErrUserNotFound, "user not found")
	}

	// 2. 验证密码
	if err := s.Crypto.VerifyPassword(user.Password, req.Password); err != nil {
		if log := s.GetLogger(); log != nil {
			log.Warn("login failed: invalid password", "username", req.Username)
		}
		return nil, errors.NewBizError(errors.ErrUnauthorized, "invalid password")
	}

	// 3. 检查用户状态
	if user.Status != 1 {
		if log := s.GetLogger(); log != nil {
			log.Warn("login failed: user inactive", "userId", user.ID, "username", user.Username, "status", user.Status)
		}
		return nil, errors.NewBizError(errors.ErrUnauthorized, "user is inactive")
	}

	// 4. 记录登录成功
	if log := s.GetLogger(); log != nil {
		log.Info("user logged in successfully", "userId", user.ID, "username", user.Username)
	}

	// 5. 异步记录登录事件
	if exec := s.GetExecutor(); exec != nil {
		userID := user.ID
		username := user.Username
		_ = exec.Execute(constants.AppPoolBackground, func() {
			// 这里可以实现：
			// - 记录登录日志（时间、IP、设备等）
			// - 更新最后登录时间
			// - 发送登录通知
			// - 检测异常登录行为
			if log := s.GetLogger(); log != nil {
				log.Debug("login event recorded", "userId", userID, "username", username)
			}
		})
	}

	// 6. 异步预热缓存
	if c := s.GetCache(); c != nil {
		if exec := s.GetExecutor(); exec != nil {
			userCopy := *user
			_ = exec.Execute(constants.AppPoolCache, func() {
				key := fmt.Sprintf("user:%d", userCopy.ID)
				if data, err := json.Marshal(userCopy); err == nil {
					_ = c.Set(context.Background(), key, string(data), 1*time.Hour)
				}
			})
		}
	}

	// 7. 生成访问令牌
	var token string
	var expiresIn int

	if jwtManager := s.GetJWT(); jwtManager != nil {
		var err error
		token, err = jwtManager.GenerateToken(user.ID, user.Username)
		if err != nil {
			if log := s.GetLogger(); log != nil {
				log.Error("failed to generate JWT token", "error", err, "userId", user.ID)
			}
			return nil, errors.NewBizError(errors.ErrInternalServer, "failed to generate token").WithCause(err)
		}
		expiresIn = 3600 // 默认 1 小时，应该从配置读取
	} else {
		// 降级处理
		if log := s.GetLogger(); log != nil {
			log.Warn("JWT manager not injected, using placeholder token")
		}
		token = "placeholder-jwt-token"
		expiresIn = 3600
	}

	// 8. 缓存 token（可选，用于 token 黑名单等功能）
	if c := s.GetCache(); c != nil {
		if exec := s.GetExecutor(); exec != nil {
			userID := user.ID
			tokenCopy := token
			_ = exec.Execute(constants.AppPoolCache, func() {
				key := fmt.Sprintf("%s%d", CacheKeyPrefixAuthToken, userID)
				_ = c.Set(context.Background(), key, tokenCopy, time.Duration(expiresIn)*time.Second)
			})
		}
	}

	// 9. 返回登录响应
	return &types.LoginResponse{
		Token:     token,
		ExpiresIn: expiresIn,
		User:      *toUserResponse(user),
	}, nil
}

// Logout 用户登出
func (s *authService) Logout(ctx context.Context, userID int64) error {
	// 1. 清除缓存的用户信息
	if c := s.GetCache(); c != nil {
		userKey := fmt.Sprintf("user:%d", userID)
		tokenKey := fmt.Sprintf("%s%d", CacheKeyPrefixAuthToken, userID)

		if exec := s.GetExecutor(); exec != nil {
			_ = exec.Execute(constants.AppPoolCache, func() {
				_ = c.Delete(context.Background(), userKey, tokenKey)
			})
		} else {
			// 如果没有 executor，同步删除
			_ = c.Delete(ctx, userKey, tokenKey)
		}
	}

	// 2. 记录登出日志
	if log := s.GetLogger(); log != nil {
		log.Info("user logged out", "userId", userID)
	}

	return nil
}

// ChangePassword 修改密码
func (s *authService) ChangePassword(ctx context.Context, userID int64, req *types.ChangePasswordRequest) error {
	// 1. 查找用户
	user, err := s.Repo.FindUserByID(ctx, userID)
	if err != nil {
		return errors.NewBizError(errors.ErrDatabaseError, "failed to find user").WithCause(err)
	}
	if user == nil {
		return errors.NewBizError(errors.ErrUserNotFound, "user not found")
	}

	// 2. 验证旧密码
	if err := s.Crypto.VerifyPassword(user.Password, req.OldPassword); err != nil {
		if log := s.GetLogger(); log != nil {
			log.Warn("change password failed: invalid old password", "userId", userID)
		}
		return errors.NewBizError(errors.ErrUnauthorized, "invalid old password")
	}

	// 3. 加密新密码
	hashedPassword, err := s.Crypto.HashPassword(req.NewPassword)
	if err != nil {
		return errors.NewBizError(errors.ErrInternalServer, "failed to hash password").WithCause(err)
	}

	// 4. 更新密码（注意：这里应该在事务中处理，暂时使用DB直接更新）
	// TODO: 重构为使用事务
	tx := s.DB.DB()
	if err := s.Repo.UpdateUserPassword(ctx, tx, userID, hashedPassword); err != nil {
		return errors.NewBizError(errors.ErrDatabaseError, "failed to update password").WithCause(err)
	}

	// 5. 清除缓存
	if c := s.GetCache(); c != nil {
		userKey := fmt.Sprintf("user:%d", userID)
		tokenKey := fmt.Sprintf("%s%d", CacheKeyPrefixAuthToken, userID)

		if exec := s.GetExecutor(); exec != nil {
			_ = exec.Execute(constants.AppPoolCache, func() {
				_ = c.Delete(context.Background(), userKey, tokenKey)
			})
		}
	}

	// 6. 记录密码修改日志
	if log := s.GetLogger(); log != nil {
		log.Info("user password changed", "userId", userID)
	}

	return nil
}

// RefreshToken 刷新访问令牌
func (s *authService) RefreshToken(ctx context.Context, req *types.RefreshTokenRequest) (*types.TokenResponse, error) {
	// 1. 验证 refresh token
	jwtManager := s.GetJWT()
	if jwtManager == nil {
		return nil, errors.NewBizError(errors.ErrInternalServer, "JWT manager not available")
	}

	// 2. 验证并提取 token 信息
	claims, err := jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		if log := s.GetLogger(); log != nil {
			log.Warn("refresh token validation failed", "error", err)
		}
		return nil, errors.NewBizError(errors.ErrUnauthorized, "invalid refresh token").WithCause(err)
	}

	// 3. 生成新的 access token
	accessToken, err := jwtManager.GenerateToken(claims.UserID, claims.Username)
	if err != nil {
		if log := s.GetLogger(); log != nil {
			log.Error("failed to generate new access token", "error", err, "userId", claims.UserID)
		}
		return nil, errors.NewBizError(errors.ErrInternalServer, "failed to generate token").WithCause(err)
	}

	// 4. 可选：生成新的 refresh token（refresh token rotation）
	// 这里暂时不实现，使用原 refresh token
	// newRefreshToken, err := jwtManager.GenerateRefreshToken(claims.UserID, claims.Username)

	// 5. 返回新的 token 响应
	return &types.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken, // 保持原 refresh token
		ExpiresIn:    3600,             // 应该从配置读取
		TokenType:    "Bearer",
	}, nil
}

// toUserResponse 将 User 模型转换为 UserResponse
func toUserResponse(user *models.DBUser) *types.UserResponse {
	return &types.UserResponse{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}
}
