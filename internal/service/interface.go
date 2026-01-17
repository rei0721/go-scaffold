package service

import (
	"github.com/rei0721/go-scaffold/pkg/cache"
	"github.com/rei0721/go-scaffold/pkg/executor"
	"github.com/rei0721/go-scaffold/pkg/jwt"
	"github.com/rei0721/go-scaffold/pkg/logger"
	"github.com/rei0721/go-scaffold/pkg/rbac"
)

type Service interface {
	// SetExecutor 设置协程池管理器（延迟注入）
	// 用于异步任务处理
	// 参数:
	//   exec: 协程池管理器实例，为nil时禁用executor功能
	// 线程安全:
	//   使用原子操作保证并发安全
	SetExecutor(exec executor.Manager)

	// SetCache 设置缓存实例（延迟注入）
	// 用于用户数据缓存
	// 参数:
	//   c: 缓存实例，为nil时禁用缓存功能
	// 线程安全:
	//   使用原子操作保证并发安全
	SetCache(c cache.Cache)

	// SetLogger 设置日志记录器（延迟注入）
	// 用于记录业务操作日志
	// 参数:
	//   l: 日志实例，为nil时禁用日志功能
	// 线程安全:
	//   使用原子操作保证并发安全
	SetLogger(l logger.Logger)

	// SetJWT 设置JWT管理器（延迟注入）
	// 用于生成访问令牌
	// 参数:
	//   j: JWT实例，为nil时使用占位符token
	// 线程安全:
	//   使用原子操作保证并发安全
	SetJWT(j jwt.JWT)

	// SetRBAC 设置RBAC管理器（延迟注入）
	// 用于用户权限管理
	// 参数:
	//   r: RBAC实例，为nil时禁用RBAC功能
	// 线程安全:
	//   使用原子操作保证并发安全
	SetRBAC(r rbac.RBAC)
}
