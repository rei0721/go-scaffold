/*
Package rbac 提供基于角色的访问控制（RBAC）功能

# 设计目标

- 通用性：可在任何 Go 项目中复用的 RBAC 实现
- 灵活性：支持多种数据库后端（通过 GORM）
- 高性能：集成缓存机制，减少数据库查询
- 线程安全：所有操作都是并发安全的
- 易集成：提供清晰的接口和依赖注入支持

# 核心概念

RBAC（Role-Based Access Control）包含以下核心实体：

1. **用户（User）**: 系统使用者，由调用方定义
2. **角色（Role）**: 权限的集合，如 "管理员"、"编辑者"、"访客"
3. **权限（Permission）**: 对特定资源的操作权限，格式为 "resource:action"
4. **用户-角色关联**: 多对多关系，一个用户可以拥有多个角色
5. **角色-权限关联**: 多对多关系，一个角色包含多个权限

权限检查流程：

	用户 -> 用户角色 -> 角色权限 -> 权限

# 包结构

	pkg/rbac/
	├── doc.go              # 本文件，包文档
	├── constants.go        # RBAC 相关常量
	├── models/             # 数据模型
	│   └── models.go       # Role, Permission, UserRole, RolePermission
	├── repository/         # 数据访问层
	│   ├── interface.go    # Repository 接口定义
	│   └── gorm.go         # GORM 实现
	└── service/            # 业务逻辑层
	    ├── interface.go    # Service 接口定义
	    └── service.go      # Service 实现

# 使用示例

基本用法：

	import (
		"github.com/rei0721/rei0721/pkg/rbac/models"
		"github.com/rei0721/rei0721/pkg/rbac/repository"
		"github.com/rei0721/rei0721/pkg/rbac/service"
	)

	// 1. 创建 Repository（传入 GORM DB 实例）
	repo := repository.NewGormRBACRepository(db)

	// 2. 创建 Service
	rbacService := service.NewRBACService(repo)

	// 3. 可选：注入缓存和协程池以提升性能
	rbacService.SetCache(cacheInstance)
	rbacService.SetExecutor(executorInstance)

	// 4. 使用 Service 进行权限管理
	// 创建角色
	role, err := rbacService.CreateRole(ctx, &types.CreateRoleRequest{
		Name:        "admin",
		Description: "管理员角色",
		Status:      1,
	})

	// 创建权限
	perm, err := rbacService.CreatePermission(ctx, &types.CreatePermissionRequest{
		Name:        "user:write",
		Resource:    "users",
		Action:      "write",
		Description: "用户写权限",
		Status:      1,
	})

	// 为角色分配权限
	err = rbacService.AssignPermission(ctx, role.ID, perm.ID)

	// 为用户分配角色
	err = rbacService.AssignRole(ctx, userID, role.ID)

	// 检查用户权限
	hasPermission, err := rbacService.CheckPermission(ctx, userID, "users", "write")

HTTP 中间件集成（示例，实际实现在调用方的 middleware 包中）：

	// 权限检查中间件
	func RequirePermission(rbacService service.RBACService, resource, action string) gin.HandlerFunc {
		return func(c *gin.Context) {
			userID := getUserIDFromContext(c)
			hasPermission, err := rbacService.CheckPermission(c.Request.Context(), userID, resource, action)
			if err != nil || !hasPermission {
				c.AbortWithStatusJSON(403, gin.H{"error": "Permission denied"})
				return
			}
			c.Next()
		}
	}

	// 在路由中使用
	api.GET("/admin/users", RequirePermission(rbacService, "users", "read"), listUsersHandler)

# 缓存策略

为了提升性能，Service 层支持可选的缓存集成：

- **权限缓存**: 用户的权限集合缓存 60 分钟
- **缓存键格式**: "user:perms:{userID}"
- **缓存失效**: 当用户角色变更时自动清除

缓存查询流程：
 1. 检查缓存是否命中
 2. 缓存命中：直接返回结果
 3. 缓存未命中：查询数据库 -> 异步写入缓存 -> 返回结果

# 权限格式

权限使用 "resource:action" 格式：

- resource: 资源标识，如 "users", "roles", "posts"
- action: 操作类型，如 "read", "write", "delete"

特殊权限：
- "*:*": 超级管理员权限，匹配所有资源和操作
- "resource:*": 资源通配符，匹配该资源的所有操作

示例：
- "users:read": 读取用户信息
- "users:write": 创建/更新用户
- "posts:delete": 删除文章
- "admin:*": 管理员模块的所有权限

# 最佳实践

1. **角色设计**: 遵循最小权限原则，按职责划分角色
2. **权限粒度**: 根据业务需求平衡权限的粗细粒度
3. **性能优化**: 生产环境务必启用缓存
4. **数据库索引**: 确保 user_roles 和 role_permissions 表有正确的索引
5. **错误处理**: 权限检查失败时记录详细日志，便于审计

# 线程安全

- Repository 实现是并发安全的（GORM 连接池）
- Service 实现使用 atomic.Value 存储可选依赖，并发安全
- 缓存操作使用协程池异步执行，避免阻塞主请求

# 依赖项

必须依赖：
- gorm.io/gorm: ORM 框架

可选依赖（通过接口注入）：
- github.com/rei0721/rei0721/pkg/cache: 缓存抽象
- github.com/rei0721/rei0721/pkg/executor: 协程池管理

# 与其他包的区别

- pkg/jwt: 负责身份认证（Authentication），验证用户是谁
- pkg/rbac: 负责授权（Authorization），验证用户能做什么
- 两者通常配合使用：JWT 验证身份后，RBAC 检查权限

# 数据库迁移

使用本包前，需要在数据库中创建以下表：
- roles: 角色表
- permissions: 权限表
- user_roles: 用户角色关联表
- role_permissions: 角色权限关联表

可以使用 pkg/sqlgen 从模型生成建表 SQL，或使用 GORM 的 AutoMigrate：

	db.AutoMigrate(
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
	)
*/
package rbac
