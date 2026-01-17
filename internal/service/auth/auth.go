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
	"github.com/rei0721/go-scaffold/pkg/executor"
	"github.com/rei0721/go-scaffold/pkg/jwt"
	"github.com/rei0721/go-scaffold/pkg/logger"
	"github.com/rei0721/go-scaffold/pkg/rbac"
	"github.com/rei0721/go-scaffold/pkg/utils"
	"github.com/rei0721/go-scaffold/types"
)

// AuthService 定义认证服务的接口
// 提供用户注册、登录、登出、密码管理等认证相关功能
type AuthService interface {
	// Register 用户注册
	// 流程：
	//   1. 验证用户名和邮箱唯一性
	//   2. 加密密码（bcrypt）
	//   3. 开启事务
	//   4. 创建用户
	//   5. 分配默认角色（如果启用 RBAC）
	//   6. 提交事务
	//   7. 缓存用户信息
	// 参数：
	//   ctx: 上下文
	//   req: 注册请求（用户名、邮箱、密码）
	// 返回：
	//   *types.UserResponse: 创建的用户信息
	//   error: 注册失败的错误
	Register(ctx context.Context, req *types.RegisterRequest) (*types.UserResponse, error)

	// Login 用户登录
	// 流程：
	//   1. 根据用户名查找用户
	//   2. 验证密码（bcrypt.CompareHashAndPassword）
	//   3. 检查用户状态（是否被禁用）
	//   4. 生成 JWT token
	//   5. 缓存用户信息和 token
	//   6. 记录登录日志（异步）
	// 参数：
	//   ctx: 上下文
	//   req: 登录请求（用户名、密码）
	// 返回：
	//   *types.LoginResponse: 登录响应（token、用户信息）
	//   error: 登录失败的错误
	Login(ctx context.Context, req *types.LoginRequest) (*types.LoginResponse, error)

	// Logout 用户登出
	// 流程：
	//   1. 清除缓存的用户信息
	//   2. 清除缓存的 token
	//   3. 记录登出日志（异步）
	// 参数：
	//   ctx: 上下文
	//   userID: 用户 ID
	// 返回：
	//   error: 登出失败的错误
	Logout(ctx context.Context, userID int64) error

	// ChangePassword 修改密码
	// 流程：
	//   1. 查找用户
	//   2. 验证旧密码
	//   3. 加密新密码
	//   4. 更新用户密码
	//   5. 清除缓存
	//   6. 记录密码修改日志（异步）
	// 参数：
	//   ctx: 上下文
	//   userID: 用户 ID
	//   req: 修改密码请求（旧密码、新密码）
	// 返回：
	//   error: 修改失败的错误
	ChangePassword(ctx context.Context, userID int64, req *types.ChangePasswordRequest) error

	// RefreshToken 刷新访问令牌
	// 流程：
	//   1. 验证 refresh token
	//   2. 提取用户信息
	//   3. 生成新的 access token
	//   4. 可选：生成新的 refresh token（refresh token rotation）
	//   5. 更新缓存
	// 参数：
	//   ctx: 上下文
	//   req: 刷新 token 请求
	// 返回：
	//   *types.TokenResponse: 新的 token 响应
	//   error: 刷新失败的错误
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
}
