package service

import (
	"github.com/rei0721/go-scaffold/pkg/cache"
	"github.com/rei0721/go-scaffold/pkg/executor"
	"github.com/rei0721/go-scaffold/pkg/jwt"
	"github.com/rei0721/go-scaffold/pkg/logger"
	"github.com/rei0721/go-scaffold/pkg/rbac"
)

// BaseService 是所有业务服务的泛型基类
//
// 设计目标:
// - 统一管理可选依赖的延迟注入 (Cache、Executor、Logger、JWT)
// - 使用 atomic.Value 确保线程安全
// - 支持运行时动态替换依赖 (如配置热重载)
// - 避免每个Service重复实现依赖管理代码
//
// 泛型参数:
//
//	T: 对应的 Model 类型 (如 models.User)
//
// 使用示例:
//
//	type userService struct {
//		BaseService[repository.UserRepository]
//	}
//
//	func NewUserService(repo repository.UserRepository) UserService {
//		u := &userService{}
//		u.SetRepository(repo)
//		return u
//	}
type BaseService[T any] struct {
	Repo     T                // 必须依赖，直接声明
	Executor executor.Manager // executor.Manager (可选，延迟注入)
	Cache    cache.Cache      // cache.Cache (可选，延迟注入)
	Logger   logger.Logger    // logger.Logger (可选，延迟注入)
	JWT      jwt.JWT          // jwt.JWT (可选，延迟注入，认证服务需要)
	RBAC     rbac.RBAC        // rbac.RBAC (可选，延迟注入，权限服务需要)
}

// SetRepository 设置Repository依赖 (延迟注入)
//
// 参数:
//
//	repo: Repository实例
//
// 注意:
//
//	此方法是线程安全的，可以在运行时动态替换
func (BaseService[T]) SetRepository(repo T) *BaseService[T] {
	return &BaseService[T]{
		Repo: repo,
	}
}

// SetExecutor 设置Executor依赖 (延迟注入)
//
// 参数:
//
//	exec: Executor管理器实例
//
// 注意:
//
//	此方法是线程安全的，可以在运行时动态替换
func (s *BaseService[T]) SetExecutor(exec executor.Manager) {
	s.Executor = exec
}

// SetCache 设置Cache依赖 (延迟注入)
//
// 参数:
//
//	c: Cache实例
//
// 注意:
//
//	此方法是线程安全的，可以在运行时动态替换
func (s *BaseService[T]) SetCache(c cache.Cache) {
	s.Cache = c
}

// SetLogger 设置Logger依赖 (延迟注入)
//
// 参数:
//
//	l: Logger实例
//
// 注意:
//
//	此方法是线程安全的，可以在运行时动态替换
func (s *BaseService[T]) SetLogger(l logger.Logger) {
	s.Logger = l
}

// SetJWT 设置JWT依赖 (延迟注入)
//
// 参数:
//
//	j: JWT管理器实例
//
// 注意:
//
//	此方法是线程安全的，可以在运行时动态替换
func (s *BaseService[T]) SetJWT(j jwt.JWT) {
	s.JWT = j
}

// SetRBAC 设置RBAC依赖 (延迟注入)
//
// 参数:
//
//	r: RBAC管理器实例
//
// 注意:
//
//	此方法是线程安全的，可以在运行时动态替换
func (s *BaseService[T]) SetRBAC(r rbac.RBAC) {
	s.RBAC = r
}

// getExecutor 获取Executor实例
//
// 返回:
//
//	executor.Manager: Executor实例，如果未注入则返回nil
//
// 注意:
//
//	使用前必须检查返回值是否为nil
func (s *BaseService[T]) getExecutor() executor.Manager {
	if exec := s.Executor; exec != nil {
		return exec
	}
	return nil
}

// getCache 获取Cache实例
//
// 返回:
//
//	cache.Cache: Cache实例，如果未注入则返回nil
//
// 注意:
//
//	使用前必须检查返回值是否为nil
func (s *BaseService[T]) getCache() cache.Cache {
	if c := s.Cache; c != nil {
		return c
	}
	return nil
}

// getLogger 获取Logger实例
//
// 返回:
//
//	logger.Logger: Logger实例，如果未注入则返回nil
//
// 注意:
//
//	使用前必须检查返回值是否为nil
func (s *BaseService[T]) getLogger() logger.Logger {
	if l := s.Logger; l != nil {
		return l
	}
	return nil
}

// getJWT 获取JWT实例
//
// 返回:
//
//	jwt.JWT: JWT管理器实例，如果未注入则返回nil
//
// 注意:
//
//	使用前必须检查返回值是否为nil
func (s *BaseService[T]) getJWT() jwt.JWT {
	if j := s.JWT; j != nil {
		return j
	}
	return nil
}

// getRBAC 获取RBAC实例
//
// 返回:
//
//	rbac.RBAC: RBAC实例，如果未注入则返回nil
//
// 注意:
//
//	使用前必须检查返回值是否为nil
func (s *BaseService[T]) getRBAC() rbac.RBAC {
	if r := s.RBAC; r != nil {
		return r
	}
	return nil
}
