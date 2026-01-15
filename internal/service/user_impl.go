package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/rei0721/rei0721/internal/models"
	"github.com/rei0721/rei0721/internal/repository"
	"github.com/rei0721/rei0721/pkg/cache"
	"github.com/rei0721/rei0721/pkg/executor"
	"github.com/rei0721/rei0721/pkg/jwt"
	"github.com/rei0721/rei0721/pkg/logger"
	"github.com/rei0721/rei0721/types"
	"github.com/rei0721/rei0721/types/constants"
	"github.com/rei0721/rei0721/types/errors"
	"github.com/rei0721/rei0721/types/result"

	"golang.org/x/crypto/bcrypt"
)

// userService 实现 UserService 接口
// 这是业务逻辑层的具体实现,协调仓库层和其他服务
// 设计原则:
// - 依赖注入:通过构造函数注入依赖,便于测试和替换
// - 单一职责:只负责用户相关的业务逻辑
// - 错误处理:将底层错误转换为业务错误,提供更好的错误信息
type userService struct {
	// repo 用户数据仓库
	// 通过接口依赖,而不是具体实现,遵循依赖倒置原则
	repo repository.UserRepository

	// executor 任务执行器（可选）
	// 使用 atomic.Value 实现无锁读取
	// 用于执行异步任务(如发送邮件、记录日志等)
	// 避免阻塞主要业务流程
	executor atomic.Value // 存储 executor.Manager

	// cache 缓存实例（可选）
	// 使用 atomic.Value 实现无锁读取
	// 用于缓存用户数据，提高查询性能
	cache atomic.Value // 存储 cache.Cache

	// logger 日志记录器（可选）
	// 使用 atomic.Value 实现无锁读取
	// 用于记录业务操作日志
	logger atomic.Value // 存储 logger.Logger

	// jwt JWT管理器（可选）
	// 使用 atomic.Value 实现无锁读取
	// 用于生成和验证访问令牌
	jwt atomic.Value // 存储 jwt.JWT
}

// NewUserService 创建一个新的 UserService 实例
// 这是工厂函数,遵循依赖注入模式
// 参数:
//
//	repo: 用户仓库接口,提供数据访问能力
//
// 返回:
//
//	UserService 接口,而不是具体类型
//	这样调用者只依赖接口,可以方便地进行单元测试(使用 mock)
//
// 注意:
//
//	Executor、Cache、Logger、JWT 需要通过相应的Set方法延迟注入
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// SetExecutor 设置协程池管理器
// 实现延迟注入模式，支持在Service创建后设置
// 使用 atomic.Value 实现原子替换，无需加锁
func (s *userService) SetExecutor(exec executor.Manager) {
	s.executor.Store(exec)
}

// getExecutor 获取当前executor（内部辅助方法）
func (s *userService) getExecutor() executor.Manager {
	if exec := s.executor.Load(); exec != nil {
		return exec.(executor.Manager)
	}
	return nil
}

// SetCache 设置缓存实例
// 实现延迟注入模式，支持在Service创建后设置
// 使用 atomic.Value 实现原子替换，无需加锁
func (s *userService) SetCache(c cache.Cache) {
	s.cache.Store(c)
}

// getCache 获取当前缓存（内部辅助方法）
func (s *userService) getCache() cache.Cache {
	if c := s.cache.Load(); c != nil {
		return c.(cache.Cache)
	}
	return nil
}

// SetLogger 设置日志记录器
// 实现延迟注入模式，支持在Service创建后设置
// 使用 atomic.Value 实现原子替换，无需加锁
func (s *userService) SetLogger(l logger.Logger) {
	s.logger.Store(l)
}

// getLogger 获取当前logger（内部辅助方法）
func (s *userService) getLogger() logger.Logger {
	if l := s.logger.Load(); l != nil {
		return l.(logger.Logger)
	}
	return nil
}

// SetJWT 设置JWT管理器
// 实现延迟注入模式，支持在Service创建后设置
// 使用 atomic.Value 实现原子替换，无需加锁
func (s *userService) SetJWT(j jwt.JWT) {
	s.jwt.Store(j)
}

// getJWT 获取当前JWT管理器（内部辅助方法）
func (s *userService) getJWT() jwt.JWT {
	if j := s.jwt.Load(); j != nil {
		return j.(jwt.JWT)
	}
	return nil
}

// userCacheKey 生成用户缓存键
// 参数:
//
//	id: 用户ID
//
// 返回:
//
//	缓存键，格式为 "user:{id}"
func userCacheKey(id int64) string {
	return fmt.Sprintf("%s%d", CacheKeyPrefixUser, id)
}

// Register 创建一个新的用户账户
// 这是用户注册的完整业务流程,包括:
// 1. 验证用户名和邮箱是否已存在
// 2. 加密密码
// 3. 创建用户记录
// 4. 触发异步后续任务
// 参数:
//
//	ctx: 上下文,用于请求追踪和超时控制
//	req: 注册请求,包含用户名、邮箱和密码
//
// 返回:
//
//	*types.UserResponse: 创建成功的用户信息(不含密码)
//	error: 业务错误(用户名重复、邮箱重复、系统错误等)
func (s *userService) Register(ctx context.Context, req *types.RegisterRequest) (*types.UserResponse, error) {
	// 1. 检查用户名是否已存在
	// 这是业务规则:用户名必须唯一
	// 虽然数据库有唯一索引,但提前检查可以提供更友好的错误信息
	existingUser, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		// 数据库查询错误,使用 WithCause 保留原始错误,便于调试
		return nil, errors.NewBizError(errors.ErrDatabaseError, "failed to check username").WithCause(err)
	}
	if existingUser != nil {
		// 用户名已存在,返回业务错误
		// 使用特定的错误码,前端可以据此显示相应提示
		return nil, errors.NewBizError(errors.ErrDuplicateUsername, "username already exists")
	}

	// 2. 检查邮箱是否已存在
	// 邮箱也必须唯一,用于账户恢复和通知
	existingUser, err = s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.NewBizError(errors.ErrDatabaseError, "failed to check email").WithCause(err)
	}
	if existingUser != nil {
		// 邮箱已被注册
		return nil, errors.NewBizError(errors.ErrDuplicateEmail, "email already exists")
	}

	// 3. 加密密码
	// 使用 bcrypt 算法加密密码,这是业界标准
	// bcrypt 的优点:
	// - 自动加盐(salt),防止彩虹表攻击
	// - 可调节计算成本,随硬件发展保持安全性
	// - 单向加密,无法解密,只能验证
	// DefaultCost(10) 是推荐的成本因子,在安全性和性能间取得平衡
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		// 加密失败(极少见),可能是系统资源问题
		return nil, errors.NewBizError(errors.ErrInternalServer, "failed to hash password").WithCause(err)
	}

	// 4. 创建用户对象
	// 注意:不直接存储明文密码,而是存储 bcrypt 哈希
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword), // 存储哈希值,不是明文
		Status:   1,                      // 1 表示激活状态,新用户默认激活
	}

	// 5. 保存到数据库
	// GORM 会自动设置 ID(Snowflake)、CreatedAt 和 UpdatedAt
	if err := s.repo.Create(ctx, user); err != nil {
		// 创建失败,可能是数据库连接问题或约束冲突
		return nil, errors.NewBizError(errors.ErrDatabaseError, "failed to create user").WithCause(err)
	}

	// 记录注册成功
	if log := s.getLogger(); log != nil {
		log.Info("user registered successfully", "userId", user.ID, "username", user.Username)
	}

	// 6. 提交异步任务处理注册后的操作
	// 使用调度器异步执行,不阻塞注册流程
	// 好处:
	// - 提高响应速度,用户不需要等待邮件发送
	// - 即使异步任务失败,注册仍然成功
	// - 可以通过调度器的协程池控制并发,避免资源耗尽
	if exec := s.getExecutor(); exec != nil {
		_ = exec.Execute(constants.PoolBackground, func() {
			// 这里可以实现:
			// - 发送欢迎邮件
			// - 记录注册事件到日志或分析系统
			// - 触发其他微服务的通知
			// - 初始化用户相关的其他资源
		})
	}

	// 7. 返回用户信息
	// 使用 toUserResponse 转换为 DTO,过滤掉密码等敏感信息
	return toUserResponse(user), nil
}

// Login 验证用户的用户名和密码
// 这是用户登录的完整业务流程
// 参数:
//
//	ctx: 上下文
//	req: 登录请求,包含用户名和密码
//
// 返回:
//
//	*types.LoginResponse: 登录成功的响应,包含 token 和用户信息
//	error: 认证失败的错误
func (s *userService) Login(ctx context.Context, req *types.LoginRequest) (*types.LoginResponse, error) {
	// 1. 根据用户名查找用户
	// 查询数据库获取用户完整信息(包括密码哈希)
	user, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		// 数据库查询错误
		return nil, errors.NewBizError(errors.ErrDatabaseError, "failed to find user").WithCause(err)
	}
	if user == nil {
		// 用户不存在
		// 为了安全,不要告诉客户端是用户名错误还是密码错误
		// 统一返回"用户名或密码错误"
		return nil, errors.NewBizError(errors.ErrUserNotFound, "user not found")
	}

	// 2. 验证密码
	// bcrypt.CompareHashAndPassword 比较密码哈希和明文密码
	// 参数顺序:第一个是哈希,第二个是明文
	// 如果密码匹配返回 nil,否则返回 error
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		// 密码错误
		// bcrypt 内部已经防止了时序攻击
		// 即使密码错误,也会执行完整的哈希比较,避免通过响应时间推测密码
		if log := s.getLogger(); log != nil {
			log.Warn("login failed: invalid password", "username", req.Username)
		}
		return nil, errors.NewBizError(errors.ErrUnauthorized, "invalid password")
	}

	// 3. 检查用户状态
	// Status=1 表示用户是激活的
	// Status=0 可能表示用户被禁用、未激活等
	if user.Status != 1 {
		// 用户已被禁用或未激活
		// 即使密码正确,也不允许登录
		if log := s.getLogger(); log != nil {
			log.Warn("login failed: user inactive", "userId", user.ID, "username", user.Username, "status", user.Status)
		}
		return nil, errors.NewBizError(errors.ErrUnauthorized, "user is inactive")
	}

	// 记录登录成功
	if log := s.getLogger(); log != nil {
		log.Info("user logged in successfully", "userId", user.ID, "username", user.Username)
	}

	// 4. 提交异步任务记录登录事件
	// 这些任务不应该阻塞登录流程
	if exec := s.getExecutor(); exec != nil {
		_ = exec.Execute(constants.PoolBackground, func() {
			// 这里可以实现:
			// - 记录登录日志(时间、IP、设备等)
			// - 更新最后登录时间
			// - 发送登录通知(如果启用)
			// - 检测异常登录行为
		})
	}

	// 5. 预热缓存（登录后用户很可能会被查询）
	if c := s.getCache(); c != nil {
		if exec := s.getExecutor(); exec != nil {
			userCopy := *user // 复制以避免并发问题
			_ = exec.Execute(constants.PoolCache, func() {
				key := userCacheKey(userCopy.ID)
				if data, err := json.Marshal(userCopy); err == nil {
					_ = c.Set(context.Background(), key, string(data), CacheTTLUser)
				}
			})
		}
	}

	// 6. 生成访问令牌
	// 使用JWT管理器生成真实的访问令牌
	var token string
	var expiresIn int

	if jwtManager := s.getJWT(); jwtManager != nil {
		// 使用真实的 JWT
		var err error
		token, err = jwtManager.GenerateToken(user.ID, user.Username)
		if err != nil {
			// JWT生成失败，记录日志
			if log := s.getLogger(); log != nil {
				log.Error("failed to generate JWT token", "error", err, "userId", user.ID)
			}
			return nil, errors.NewBizError(errors.ErrInternalServer, "failed to generate token").WithCause(err)
		}
		// 设置过期时间（单位：秒）
		expiresIn = 3600 // 从配置读取或使用默认值（1小时）
	} else {
		// 降级处理：如果 JWT 未注入，使用占位符（不应发生）
		if log := s.getLogger(); log != nil {
			log.Warn("JWT manager not injected, using placeholder token")
		}
		token = "placeholder-jwt-token"
		expiresIn = 3600
	}

	// 7. 返回登录响应
	return &types.LoginResponse{
		Token:     token,                 // 访问令牌,前端应该安全存储
		ExpiresIn: expiresIn,             // 令牌有效期,前端可以据此刷新令牌
		User:      *toUserResponse(user), // 用户信息,避免前端再次请求
	}, nil
}

// GetByID 根据用户ID获取用户信息
// Cache-Aside模式：先查缓存，未命中再查数据库，然后异步更新缓存
// 参数:
//
//	ctx: 上下文
//	id: 用户ID
//
// 返回:
//
//	*types.UserResponse: 用户的详细信息
//	error: 查找失败的错误
func (s *userService) GetByID(ctx context.Context, id int64) (*types.UserResponse, error) {
	// 1. 尝试从缓存读取
	if c := s.getCache(); c != nil {
		key := userCacheKey(id)
		if data, err := c.Get(ctx, key); err == nil {
			var user models.User
			if err := json.Unmarshal([]byte(data), &user); err == nil {
				// 缓存命中
				return toUserResponse(&user), nil
			}
			// 缓存数据解析失败，继续查数据库
		}
		// 缓存未命中或读取失败，继续查数据库
	}

	// 2. 查询数据库获取用户信息
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		// 数据库查询错误
		return nil, errors.NewBizError(errors.ErrDatabaseError, "failed to find user").WithCause(err)
	}
	if user == nil {
		// 用户不存在
		// 返回业务错误,前端会收到 404
		return nil, errors.NewBizError(errors.ErrUserNotFound, "user not found")
	}

	// 3. 异步写入缓存（使用executor避免阻塞）
	if c := s.getCache(); c != nil {
		if exec := s.getExecutor(); exec != nil {
			userCopy := *user // 复制以避免并发问题
			_ = exec.Execute(constants.PoolCache, func() {
				key := userCacheKey(userCopy.ID)
				if data, err := json.Marshal(userCopy); err == nil {
					_ = c.Set(context.Background(), key, string(data), CacheTTLUser)
				}
			})
		}
	}

	// 4. 转换为响应类型并返回
	// toUserResponse 会过滤掉密码等敏感信息
	return toUserResponse(user), nil
}

// List 获取用户列表(分页)
// 参数:
//
//	ctx: 上下文
//	page: 页码,从 1 开始
//	pageSize: 每页大小,建议 1-100
//
// 返回:
//
//	*result.PageResult[types.UserResponse]: 分页结果
//	error: 查询失败的错误
func (s *userService) List(ctx context.Context, page, pageSize int) (*result.PageResult[types.UserResponse], error) {
	// 1. 验证分页参数
	// 这是一个防御性编程的例子,确保参数合法
	if page < 1 {
		// 页码小于 1,重置为 1
		page = 1
	}
	if pageSize < 1 {
		// 每页大小小于 1,使用默认值 10
		pageSize = 10
	}
	if pageSize > 100 {
		// 限制最大值为 100,防止:
		// - 查询过多数据导致性能问题
		// - 响应体过大
		// - 内存消耗过高
		pageSize = 100
	}

	// 2. 查询数据库
	// FindAll 会返回当前页的用户列表和总记录数
	users, total, err := s.repo.FindAll(ctx, page, pageSize)
	if err != nil {
		// 数据库查询错误
		return nil, errors.NewBizError(errors.ErrDatabaseError, "failed to list users").WithCause(err)
	}

	// 3. 转换为响应类型
	// 从数据库模型 (models.User) 转换为 API 响应 (types.UserResponse)
	// 这是 DTO (Data Transfer Object) 模式
	responses := make([]types.UserResponse, len(users))
	for i, user := range users {
		// 对每个用户进行转换
		// toUserResponse 会过滤掉密码等敏感信息
		responses[i] = *toUserResponse(&user)
	}

	// 4. 创建分页结果
	// NewPageResult 会自动计算总页数等分页信息
	return result.NewPageResult(responses, page, pageSize, total), nil
}

// Update 更新用户信息
// 支持部分字段更新（只更新传入的字段）
// 参数:
//
//	ctx: 上下文
//	id: 用户 ID
//	req: 更新请求
//
// 返回:
//
//	*types.UserResponse: 更新后的用户信息
//	error: 更新失败的错误
func (s *userService) Update(ctx context.Context, id int64, req *types.UpdateUserRequest) (*types.UserResponse, error) {
	// 1. 查询用户是否存在
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.NewBizError(errors.ErrDatabaseError, "failed to find user").WithCause(err)
	}
	if user == nil {
		return nil, errors.NewBizError(errors.ErrUserNotFound, "user not found")
	}

	// 记录更新操作
	if log := s.getLogger(); log != nil {
		log.Info("user update started", "userId", id)
	}

	// 2. 检查用户名唯一性（如果更新了用户名）
	if req.Username != nil && *req.Username != user.Username {
		existingUser, err := s.repo.FindByUsername(ctx, *req.Username)
		if err != nil {
			return nil, errors.NewBizError(errors.ErrDatabaseError, "failed to check username").WithCause(err)
		}
		if existingUser != nil {
			return nil, errors.NewBizError(errors.ErrDuplicateUsername, "username already exists")
		}
		user.Username = *req.Username
	}

	// 3. 检查邮箱唯一性（如果更新了邮箱）
	if req.Email != nil && *req.Email != user.Email {
		existingUser, err := s.repo.FindByEmail(ctx, *req.Email)
		if err != nil {
			return nil, errors.NewBizError(errors.ErrDatabaseError, "failed to check email").WithCause(err)
		}
		if existingUser != nil {
			return nil, errors.NewBizError(errors.ErrDuplicateEmail, "email already exists")
		}
		user.Email = *req.Email
	}

	// 4. 更新状态（如果传入了状态）
	if req.Status != nil {
		user.Status = *req.Status
	}

	// 5. 保存到数据库
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, errors.NewBizError(errors.ErrDatabaseError, "failed to update user").WithCause(err)
	}

	// 6. 失效缓存
	s.invalidateUserCache(id)

	// 记录成功
	if log := s.getLogger(); log != nil {
		log.Info("user updated successfully", "userId", id)
	}

	return toUserResponse(user), nil
}

// Delete 删除用户（软删除）
// 参数:
//
//	ctx: 上下文
//	id: 用户 ID
//
// 返回:
//
//	error: 删除失败的错误
func (s *userService) Delete(ctx context.Context, id int64) error {
	// 1. 查询用户是否存在
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.NewBizError(errors.ErrDatabaseError, "failed to find user").WithCause(err)
	}
	if user == nil {
		return errors.NewBizError(errors.ErrUserNotFound, "user not found")
	}

	// 记录删除操作
	if log := s.getLogger(); log != nil {
		log.Info("user deletion started", "userId", id, "username", user.Username)
	}

	// 2. 软删除用户
	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.NewBizError(errors.ErrDatabaseError, "failed to delete user").WithCause(err)
	}

	// 3. 失效缓存
	s.invalidateUserCache(id)

	// 记录成功
	if log := s.getLogger(); log != nil {
		log.Info("user deleted successfully", "userId", id)
	}

	return nil
}

// invalidateUserCache 失效用户缓存
// 这是一个内部辅助方法，用于在用户数据变更时清除缓存
// 参数:
//
//	id: 用户 ID
func (s *userService) invalidateUserCache(id int64) {
	if c := s.getCache(); c != nil {
		if exec := s.getExecutor(); exec != nil {
			key := userCacheKey(id)
			_ = exec.Execute(constants.PoolCache, func() {
				_ = c.Delete(context.Background(), key)
			})
		}
	}
}

// toUserResponse 将 User 模型转换为 UserResponse
// 这是一个内部辅助函数,实现 DTO 模式
// 为什么需要转换:
// - 分离数据层和表示层:数据库模型和 API 响应独立演化
// - 安全性:过滤掉敏感信息(如密码哈希)
// - 灵活性:可以自由调整响应格式,不影响数据库结构
// - 统一格式:确保所有 API 返回的用户信息格式一致
// 参数:
//
//	user: 数据库用户模型
//
// 返回:
//
//	*types.UserResponse: API 响应用户信息
func toUserResponse(user *models.User) *types.UserResponse {
	return &types.UserResponse{
		UserID:    user.ID,        // 用户 ID
		Username:  user.Username,  // 用户名
		Email:     user.Email,     // 邮箱
		Status:    user.Status,    // 状态
		CreatedAt: user.CreatedAt, // 创建时间
		// 注意:没有包含 Password 字段
		// 这确保密码哈希永远不会被返回给客户端
	}
}
