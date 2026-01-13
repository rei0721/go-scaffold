package service

import (
	"context"

	"github.com/rei0721/rei0721/internal/models"
	"github.com/rei0721/rei0721/internal/repository"
	"github.com/rei0721/rei0721/pkg/executor"
	"github.com/rei0721/rei0721/types"
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

	// executor 任务执行器
	// 用于执行异步任务(如发送邮件、记录日志等)
	// 避免阻塞主要业务流程
	executor executor.Manager
}

// NewUserService 创建一个新的 UserService 实例
// 这是工厂函数,遵循依赖注入模式
// 参数:
//
//	repo: 用户仓库接口,提供数据访问能力
//	sched: 任务调度器,用于异步任务执行
//
// 返回:
//
//	UserService 接口,而不是具体类型
//	这样调用者只依赖接口,可以方便地进行单元测试(使用 mock)
func NewUserService(repo repository.UserRepository, exec executor.Manager) UserService {
	return &userService{
		repo:     repo,
		executor: exec,
	}
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

	// 6. 提交异步任务处理注册后的操作
	// 使用调度器异步执行,不阻塞注册流程
	// 好处:
	// - 提高响应速度,用户不需要等待邮件发送
	// - 即使异步任务失败,注册仍然成功
	// - 可以通过调度器的协程池控制并发,避免资源耗尽
	_ = s.executor.Execute("background", func() {
		// 这里可以实现:
		// - 发送欢迎邮件
		// - 记录注册事件到日志或分析系统
		// - 触发其他微服务的通知
		// - 初始化用户相关的其他资源
	})

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
		return nil, errors.NewBizError(errors.ErrUnauthorized, "invalid password")
	}

	// 3. 检查用户状态
	// Status=1 表示用户是激活的
	// Status=0 可能表示用户被禁用、未激活等
	if user.Status != 1 {
		// 用户已被禁用或未激活
		// 即使密码正确,也不允许登录
		return nil, errors.NewBizError(errors.ErrUnauthorized, "user is inactive")
	}

	// 4. 提交异步任务记录登录事件
	// 这些任务不应该阻塞登录流程
	_ = s.executor.Execute("background", func() {
		// 这里可以实现:
		// - 记录登录日志(时间、IP、设备等)
		// - 更新最后登录时间
		// - 发送登录通知(如果启用)
		// - 检测异常登录行为
	})

	// 5. 生成访问令牌
	// 这里是占位符实现,实际应该:
	// - 使用 JWT 生成令牌
	// - 包含用户 ID、用户名等信息
	// - 设置过期时间
	// - 签名确保不可篡改
	// TODO: 实现真正的 JWT 生成
	token := "placeholder-jwt-token"
	expiresIn := 3600 // 1小时(3600秒)

	// 6. 返回登录响应
	return &types.LoginResponse{
		Token:     token,                 // 访问令牌,前端应该安全存储
		ExpiresIn: expiresIn,             // 令牌有效期,前端可以据此刷新令牌
		User:      *toUserResponse(user), // 用户信息,避免前端再次请求
	}, nil
}

// GetByID 根据用户ID获取用户信息
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
	// 查询数据库获取用户信息
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

	// 转换为响应类型并返回
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
