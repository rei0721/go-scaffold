package repository

import (
	"context"
	"errors"

	"github.com/rei0721/go-scaffold/internal/models"

	"gorm.io/gorm"
)

// UserRepository 扩展了通用的 Repository 接口
// 添加了用户特定的查询方法
// 这种设计模式称为"接口组合",符合 Go 的编程哲学
type UserRepository interface {
	// 嵌入泛型 Repository 接口,继承基本的 CRUD 方法
	// Repository[models.User] 使用 Go 泛型,指定操作的是 User 模型
	Repository[models.User]

	// FindByUsername 根据用户名检索用户
	// 返回 nil 表示用户不存在(不是错误)
	// 这种设计使得调用者可以区分"用户不存在"和"数据库错误"
	FindByUsername(ctx context.Context, username string) (*models.User, error)

	// FindByEmail 根据邮箱地址检索用户
	// 返回 nil 表示用户不存在
	// 邮箱用于登录和找回密码场景
	FindByEmail(ctx context.Context, email string) (*models.User, error)

	// CreateWithTx 在事务中创建用户
	// 参数:
	//   ctx: 上下文
	//   tx: GORM 事务对象
	//   user: 要创建的用户
	// 返回:
	//   error: 创建失败的错误
	// 使用场景:
	//   - 注册时需要同时创建用户和分配角色
	//   - 需要保证数据一致性的操作
	CreateWithTx(ctx context.Context, tx *gorm.DB, user *models.User) error

	// UpdateWithTx 在事务中更新用户
	// 参数:
	//   ctx: 上下文
	//   tx: GORM 事务对象
	//   user: 要更新的用户
	// 返回:
	//   error: 更新失败的错误
	// 使用场景:
	//   - 修改密码时需要同时更新其他相关数据
	//   - 需要保证数据一致性的操作
	UpdateWithTx(ctx context.Context, tx *gorm.DB, user *models.User) error
}

// userRepository 使用 GORM 实现 UserRepository 接口
// 这是一个私有结构体,外部只能通过接口访问
// 这样做的好处是可以轻松替换实现(例如切换到其他 ORM 或数据库)
type userRepository struct {
	// db GORM 数据库连接
	// 使用指针以避免复制大对象
	db *gorm.DB
}

// NewUserRepository 创建一个新的 UserRepository 实例
// 这是一个工厂函数,遵循依赖注入模式
// 参数:
//
//	db: GORM 数据库连接,由调用者传入
//
// 返回:
//
//	UserRepository 接口,而不是具体的实现类型
//	这样可以保持实现细节的隐藏性
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 向数据库插入一个新用户
// 参数:
//
//	ctx: 上下文,用于超时控制和请求追踪
//	user: 要创建的用户对象,ID 字段会在插入后被 GORM 自动填充
//
// 返回:
//
//	error: 如果插入失败返回错误(例如违反唯一约束)
//
// 注意:
//   - GORM 会自动设置 CreatedAt 和 UpdatedAt 字段
//   - 如果 Username 或 Email 重复,会返回数据库约束错误
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	// WithContext 将 Go 的 context 传递给 GORM
	// 这使得查询可以被取消或超时,提高系统的健壮性
	// Create 执行 INSERT 操作
	// .Error 获取执行结果中的错误
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID 根据 ID 检索用户
// 参数:
//
//	ctx: 上下文
//	id: 用户 ID(Snowflake 算法生成的 int64)
//
// 返回:
//
//	*models.User: 找到的用户,如果不存在返回 nil(不是错误)
//	error: 数据库错误(不包括"记录未找到")
//
// 注意:
//
//	这种返回模式允许调用者区分"用户不存在"和"数据库错误"
func (r *userRepository) FindByID(ctx context.Context, id int64) (*models.User, error) {
	// 声明用户变量用于接收查询结果
	var user models.User

	// First 查询第一条匹配的记录
	// 传入 id 作为主键值,GORM 会自动构造 WHERE id = ? 条件
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		// 使用 errors.Is 判断是否为"记录未找到"错误
		// 这是 Go 1.13+ 推荐的错误处理方式
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户不存在不视为错误,返回 (nil, nil)
			// 这样调用者可以通过检查返回的用户是否为 nil 来判断
			return nil, nil
		}
		// 其他数据库错误(连接失败、查询超时等)需要返回
		return nil, err
	}
	// 返回用户指针
	return &user, nil
}

// FindAll 检索所有用户,支持分页
// 参数:
//
//	ctx: 上下文
//	page: 页码,从 1 开始
//	pageSize: 每页大小
//
// 返回:
//
//	[]models.User: 当前页的用户列表
//	int64: 总记录数,用于计算总页数
//	error: 数据库错误
//
// 注意:
//
//	软删除的用户会被自动过滤(DeletedAt IS NULL)
func (r *userRepository) FindAll(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	// 声明用户切片用于接收查询结果
	var users []models.User
	// 声明总记录数变量
	var total int64

	// 先统计总记录数
	// Model(&models.User{}) 指定要查询的表
	// Count(&total) 执行 COUNT(*) 查询并将结果赋值给 total
	// 这样可以让前端计算总页数和显示分页导航
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 计算偏移量(OFFSET)
	// 例如: page=1 时 offset=0, page=2 时 offset=pageSize
	// 这是标准的分页计算公式
	offset := (page - 1) * pageSize

	// 执行分页查询
	// Offset(offset) 设置跳过的记录数(SQL: OFFSET offset)
	// Limit(pageSize) 设置返回的最大记录数(SQL: LIMIT pageSize)
	// Find(&users) 查询所有匹配的记录
	// GORM 会自动添加 WHERE deleted_at IS NULL 条件(因为 User 使用了软删除)
	if err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// 返回用户列表、总数和错误(如果有)
	return users, total, nil
}

// Update 修改数据库中的现有用户
// 参数:
//
//	ctx: 上下文
//	user: 要更新的用户对象,必须包含有效的 ID
//
// 返回:
//
//	error: 更新失败时返回错误
//
// 注意:
//   - Save 会更新所有字段,包括零值字段
//   - GORM 会自动更新 UpdatedAt 字段为当前时间
//   - 如果要只更新部分字段,应使用 Updates 方法
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	// Save 执行 UPDATE 操作,根据主键 ID 更新记录
	// 它会保存所有字段,即使是零值
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 根据 ID 删除用户(软删除)
// 参数:
//
//	ctx: 上下文
//	id: 要删除的用户 ID
//
// 返回:
//
//	error: 删除失败时返回错误
//
// 注意:
//   - 这是软删除,实际上是将 DeletedAt 字段设置为当前时间
//   - 数据仍保留在数据库中,可以恢复
//   - 软删除的记录在正常查询时会被自动过滤
//   - 如果需要硬删除(物理删除),应使用 Unscoped().Delete()
func (r *userRepository) Delete(ctx context.Context, id int64) error {
	// Delete 执行软删除操作
	// &models.User{} 提供模型类型信息
	// id 是要删除的记录的主键值
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// FindByUsername 根据用户名检索用户
// 参数:
//
//	ctx: 上下文
//	username: 用户名(已建立唯一索引)
//
// 返回:
//
//	*models.User: 找到的用户,不存在时返回 nil
//	error: 数据库错误
//
// 使用场景:
//   - 用户登录
//   - 注册时检查用户名是否已存在
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	// Where 添加查询条件
	// 使用 "?" 占位符防止 SQL 注入
	// GORM 会自动转义参数
	// 由于 username 字段有唯一索引,查询性能很高
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户不存在,返回 nil 而不是错误
			return nil, nil
		}
		// 其他数据库错误
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱地址检索用户
// 参数:
//
//	ctx: 上下文
//	email: 邮箱地址(已建立唯一索引)
//
// 返回:
//
//	*models.User: 找到的用户,不存在时返回 nil
//	error: 数据库错误
//
// 使用场景:
//   - 邮箱登录
//   - 注册时检查邮箱是否已存在
//   - 找回密码功能
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	// Where 添加 email 查询条件
	// 使用参数化查询防止 SQL 注入攻击
	// email 字段有唯一索引,查询会很快
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 邮箱不存在,返回 nil
			return nil, nil
		}
		// 其他数据库错误
		return nil, err
	}
	return &user, nil
}

// CreateWithTx 在事务中创建用户
// 参数:
//
//	ctx: 上下文
//	tx: GORM 事务对象
//	user: 要创建的用户
//
// 返回:
//
//	error: 创建失败的错误
//
// 注意:
//   - 使用传入的事务对象而不是 r.db
//   - 调用方负责事务的开启、提交和回滚
func (r *userRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, user *models.User) error {
	return tx.WithContext(ctx).Create(user).Error
}

// UpdateWithTx 在事务中更新用户
// 参数:
//
//	ctx: 上下文
//	tx: GORM 事务对象
//	user: 要更新的用户
//
// 返回:
//
//	error: 更新失败的错误
//
// 注意:
//   - 使用传入的事务对象而不是 r.db
//   - 调用方负责事务的开启、提交和回滚
func (r *userRepository) UpdateWithTx(ctx context.Context, tx *gorm.DB, user *models.User) error {
	return tx.WithContext(ctx).Save(user).Error
}
