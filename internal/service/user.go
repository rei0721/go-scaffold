// Package service 提供业务逻辑层的实现
// 服务层依赖 Repository 接口,使用 Scheduler 执行异步任务
// 设计原则:
// - 业务逻辑与数据访问分离
// - 通过接口定义规范,便于测试和扩展
// - 使用上下文传递请求信息
package service

import (
	"context"

	"github.com/rei0721/rei0721/types"
	"github.com/rei0721/rei0721/types/result"
)

// UserService 定义用户相关业务操作的接口
// 这是一个接口定义,具体实现在 user_impl.go 中
// 为什么使用接口:
// - 定义契约:明确业务层提供哪些功能
// - 依赖倒置:handler 层依赖接口而非实现
// - 便于测试:可以创建 mock 实现进行单元测试
// - 解耦:可以轻松替换不同的实现
type UserService interface {
	// Register 创建一个新的用户账户
	// 这是用户注册的业务接口
	// 参数:
	//   ctx: 上下文,用于:
	//     - 传递请求超时
	//     - 取消操作
	//     - 传递 TraceID 等元数据
	//   req: 注册请求,包含用户名、邮箱和密码
	// 返回:
	//   *types.UserResponse: 创建成功的用户信息(不含密码)
	//   error: 业务错误,如:
	//     - ErrDuplicateUsername: 用户名已存在
	//     - ErrDuplicateEmail: 邮箱已存在
	//     - ErrDatabaseError: 数据库错误
	// 业务流程:
	//   1. 验证用户名是否已存在
	//   2. 验证邮箱是否已存在
	//   3. 对密码进行加密(bcrypt)
	//   4. 创建用户记录
	//   5. 触发异步后续任务(如发送欢迎邮件)
	Register(ctx context.Context, req *types.RegisterRequest) (*types.UserResponse, error)

	// Login 使用用户名和密码进行身份验证
	// 这是用户登录的业务接口
	// 参数:
	//   ctx: 上下文
	//   req: 登录请求,包含用户名和密码
	// 返回:
	//   *types.LoginResponse: 登录响应,包含:
	//     - Token: 访问令牌(JWT)
	//     - ExpiresIn: 令牌有效期(秒)
	//     - User: 用户信息
	//   error: 认证错误,如:
	//     - ErrUserNotFound: 用户不存在
	//     - ErrUnauthorized: 密码错误
	//     - ErrDatabaseError: 数据库错误
	// 业务流程:
	//   1. 根据用户名查找用户
	//   2. 验证密码(bcrypt.CompareHashAndPassword)
	//   3. 检查用户状态(是否被禁用)
	//   4. 生成访问令牌
	//   5. 记录登录事件(异步)
	Login(ctx context.Context, req *types.LoginRequest) (*types.LoginResponse, error)

	// GetByID 根据用户 ID 获取用户信息
	// 这是用户查询的业务接口
	// 参数:
	//   ctx: 上下文
	//   id: 用户 ID(Snowflake 算法生成的 int64)
	// 返回:
	//   *types.UserResponse: 用户信息,如果用户不存在返回 nil
	//   error: 查询错误,如:
	//     - ErrUserNotFound: 用户不存在
	//     - ErrDatabaseError: 数据库错误
	// 使用场景:
	//   - 获取用户详情页信息
	//   - 验证用户是否存在
	//   - 获取用户基本信息用于显示
	GetByID(ctx context.Context, id int64) (*types.UserResponse, error)

	// List 获取用户列表(分页)
	// 这是用户列表查询的业务接口
	// 参数:
	//   ctx: 上下文
	//   page: 页码,从 1 开始
	//   pageSize: 每页大小,建议范围 1-100
	// 返回:
	//   *result.PageResult[types.UserResponse]: 分页结果,包含:
	//     - List: 当前页的用户列表
	//     - Pagination: 分页信息(总数、总页数等)
	//   error: 查询错误,如:
	//     - ErrDatabaseError: 数据库错误
	// 使用场景:
	//   - 用户管理后台
	//   - 用户列表展示
	// 注意:
	//   - 返回的用户不包含密码等敏感信息
	//   - 只返回未被软删除的用户
	//   - 服务层会对 page 和 pageSize 进行验证
	List(ctx context.Context, page, pageSize int) (*result.PageResult[types.UserResponse], error)
}
