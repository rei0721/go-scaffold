// Package auth 提供认证服务的实现
// 职责：
// - 用户注册（创建用户 + 分配默认角色）
// - 用户登录（验证凭证 + 生成 Token）
// - 用户登出（清除缓存/会话）
// - 密码修改（验证旧密码 + 更新新密码）
// - Token 刷新（验证 refresh token + 生成新 access token）
//
// 设计原则：
// - 与 UserService 职责分离：Auth 负责认证，User 负责用户资料管理
// - 支持事务：注册等操作使用事务保证数据一致性
// - 集成现有组件：JWT、RBAC、Cache、Logger、Executor
package auth

import (
	"context"

	"github.com/rei0721/go-scaffold/pkg/cache"
	"github.com/rei0721/go-scaffold/pkg/crypto"
	"github.com/rei0721/go-scaffold/pkg/database"
	"github.com/rei0721/go-scaffold/pkg/dbtx"
	"github.com/rei0721/go-scaffold/pkg/executor"
	"github.com/rei0721/go-scaffold/pkg/jwt"
	"github.com/rei0721/go-scaffold/pkg/logger"
	"github.com/rei0721/go-scaffold/pkg/rbac"
	"github.com/rei0721/go-scaffold/pkg/utils"
	"github.com/rei0721/go-scaffold/types"
)

// AuthService 定义认证服务的接口
type AuthService interface {
	// Register 用户注册
	Register(ctx context.Context, req *types.RegisterRequest) (*types.UserResponse, error)

	// Login 用户登录
	Login(ctx context.Context, req *types.LoginRequest) (*types.LoginResponse, error)

	// Logout 用户登出
	Logout(ctx context.Context, userID int64) error

	// ChangePassword 修改密码
	ChangePassword(ctx context.Context, userID int64, req *types.ChangePasswordRequest) error

	// RefreshToken 刷新访问令牌
	RefreshToken(ctx context.Context, req *types.RefreshTokenRequest) (*types.TokenResponse, error)

	// SetDB 设置DB依赖（延迟注入）
	SetDB(db database.Database)

	// SetExecutor 设置协程池管理器（延迟注入）
	SetExecutor(exec executor.Manager)

	// SetCache 设置缓存实例（延迟注入）
	SetCache(c cache.Cache)

	// SetLogger 设置日志记录器（延迟注入）
	SetLogger(l logger.Logger)

	// SetJWT 设置JWT管理器（延迟注入）
	SetJWT(j jwt.JWT)

	// SetRBAC 设置RBAC管理器（延迟注入）
	SetRBAC(r rbac.RBAC)

	// SetIDGenerator 设置ID生成器（延迟注入）
	SetIDGenerator(idGenerator utils.IDGenerator)

	// SetCrypto 设置密码加密器（延迟注入）
	SetCrypto(c crypto.Crypto)

	// SetTxManager 设置事务管理器（延迟注入）
	SetTxManager(txMgr dbtx.Manager)
}
