// Package repository 提供数据访问层的抽象
// 实现仓库模式进行数据库操作
// 设计目标:
// - 抽象数据访问:隐藏数据库实现细节
// - 统一接口:所有实体类型使用相同的 CRUD 接口
// - 便于测试:可以 mock Repository 进行单元测试
// - 支持泛型:使用 Go 1.18+ 泛型提高代码复用
package repository

import (
	"context"
)

// Repository 定义通用的 CRUD 操作接口
// 使用泛型 T 表示仓库管理的实体类型
// 这是一个通用接口,任何实体都可以使用
// 例如:
//
//	Repository[models.User]
//	Repository[models.Product]
//	Repository[models.Order]
//
// 设计考虑:
// - 泛型避免为每个实体写重复代码
// - 提供基础的 CRUD 操作
// - 特定实体可以在此基础上扩展(如 UserRepository)
type Repository[T any] interface {
	// Create 插入一个新实体到数据库
	// 参数:
	//   ctx: 上下文,用于超时控制和取消操作
	//   entity: 要创建的实体
	//     - ID 字段会被自动设置(Snowflake)
	//     - CreatedAt 和 UpdatedAt 由 GORM 自动设置
	// 返回:
	//   error: 创建失败时的错误
	//     - 唯一约束冲突
	//     - 数据库连接错误
	//     - 验证失败
	// 使用场景:
	//   - 用户注册
	//   - 创建订单
	//   - 添加产品
	Create(ctx context.Context, entity *T) error

	// FindByID 根据 ID 检索实体
	// 参数:
	//   ctx: 上下文
	//   id: 实体的 ID(Snowflake 生成的 int64)
	// 返回:
	//   *T: 找到的实体,如果不存在返回 nil
	//   error: 数据库错误(不包括"记录未找到")
	// 注意:
	//   - 返回 nil, nil 表示实体不存在(不是错误)
	//   - 只有数据库错误才返回 error
	//   - 自动排除软删除的记录
	// 使用场景:
	//   - 获取用户详情
	//   - 查看订单信息
	//   - 验证实体是否存在
	FindByID(ctx context.Context, id int64) (*T, error)

	// FindAll 检索所有实体,支持分页
	// 参数:
	//   ctx: 上下文
	//   page: 页码,从 1 开始
	//   pageSize: 每页大小
	// 返回:
	//   []T: 当前页的实体列表
	//   int64: 总记录数(用于计算总页数)
	//   error: 查询错误
	// 查询逻辑:
	//   - 自动计算 offset: (page-1) * pageSize
	//   - 使用 LIMIT 和 OFFSET 实现分页
	//   - 同时返回总记录数(用于前端显示总页数)
	//   - 自动排除软删除的记录
	// 使用场景:
	//   - 用户列表
	//   - 产品目录
	//   - 订单历史
	FindAll(ctx context.Context, page, pageSize int) ([]T, int64, error)

	// Update 修改数据库中的现有实体
	// 参数:
	//   ctx: 上下文
	//   entity: 要更新的实体
	//     - 必须包含有效的 ID
	//     - UpdatedAt 会被 GORM 自动更新
	// 返回:
	//   error: 更新失败时的错误
	//     - 实体不存在
	//     - 唯一约束冲突
	//     - 数据库错误
	// 注意:
	//   - 默认只更新非零值字段
	//   - 如果需要更新为零值,使用 Select 指定字段
	//   - 不会更新软删除的记录
	// 使用场景:
	//   - 修改用户信息
	//   - 更新订单状态
	//   - 编辑产品详情
	Update(ctx context.Context, entity *T) error

	// Delete 根据 ID 删除实体(如果支持则软删除)
	// 参数:
	//   ctx: 上下文
	//   id: 要删除的实体 ID
	// 返回:
	//   error: 删除失败时的错误
	// 删除行为:
	//   - 如果实体有 DeletedAt 字段(如 BaseModel):软删除
	//     * 设置 DeletedAt 为当前时间
	//     * 记录仍在数据库中,但查询时自动排除
	//     * 可以恢复
	//   - 如果实体没有 DeletedAt 字段:硬删除
	//     * 永久从数据库删除
	//     * 无法恢复
	// 为什么使用软删除:
	//   - 数据安全:误删除可以恢复
	//   - 审计友好:保留历史记录
	//   - 关联数据:避免外键约束问题
	// 使用场景:
	//   - 删除用户账户
	//   - 取消订单
	//   - 下架产品
	Delete(ctx context.Context, id int64) error
}
