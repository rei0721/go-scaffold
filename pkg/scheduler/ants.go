package scheduler

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/panjf2000/ants/v2"
)

// antsScheduler 使用 ants 协程池实现 Scheduler 接口
// ants 是一个高性能的 goroutine 池,优点:
// - 复用 goroutine,减少创建/销毁开销
// - 限制并发数量,防止资源耗尽
// - 自动管理协程生命周期
type antsScheduler struct {
	// pool ants 协程池实例
	// 所有提交的任务都会在这个池中的 goroutine 上执行
	pool *ants.Pool
}

// New 创建一个新的 Scheduler 实例
// 参数:
//
//	cfg: 调度器配置,包含池大小、过期时间等
//
// 返回:
//
//	Scheduler 接口
//	error: 创建失败时的错误
//
// 使用场景:
//   - 异步发送邮件
//   - 后台日志记录
//   - 定时任务触发
func New(cfg *Config) (Scheduler, error) {
	// 如果未提供配置,使用默认配置
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// 配置 ants 池的选项
	opts := []ants.Option{
		// WithExpiryDuration 设置空闲 goroutine 的过期时间
		// 超过此时间未使用的 goroutine 会被回收,释放内存
		ants.WithExpiryDuration(cfg.ExpiryDuration),

		// WithPreAlloc(false) 不预分配 goroutine
		// false 表示按需创建,节省初始内存
		// true 会预先创建所有 goroutine,适合高并发场景
		ants.WithPreAlloc(false),

		// WithNonblocking(false) 使用阻塞模式
		// false: 当池满时,Submit 会阻塞等待
		// true: 当池满时,Submit 立即返回错误
		// 阻塞模式更适合后台任务,确保所有任务最终都能执行
		ants.WithNonblocking(false),
	}

	// 如果配置了最大阻塞任务数
	if cfg.MaxBlockingTasks > 0 {
		// 限制等待队列的大小,防止无限制等待
		opts = append(opts, ants.WithMaxBlockingTasks(cfg.MaxBlockingTasks))
	}

	// 创建协程池
	// PoolSize 决定了最多同时运行的 goroutine 数量
	pool, err := ants.NewPool(cfg.PoolSize, opts...)
	if err != nil {
		// 创建失败,可能是参数无效(如 PoolSize <= 0)
		return nil, err
	}

	// 返回包装后的调度器
	return &antsScheduler{
		pool: pool,
	}, nil
}

// Default 使用默认配置创建 Scheduler
// 这是一个便捷函数,适合快速开始使用
// 注意:如果创建失败会 panic,因为默认配置不应该失败
func Default() Scheduler {
	s, err := New(DefaultConfig())
	if err != nil {
		// 使用默认配置不应该失败
		// 如果失败了,说明代码有 bug,应该立即暴露
		panic("failed to create default scheduler: " + err.Error())
	}
	return s
}

// Submit submits a task to the goroutine pool.
// The task function receives the provided context for cancellation support.
func (s *antsScheduler) Submit(ctx context.Context, task func(context.Context)) error {
	if ctx == nil {
		ctx = context.Background()
	}

	return s.pool.Submit(func() {
		// Check if context is already cancelled before executing
		select {
		case <-ctx.Done():
			return
		default:
		}
		task(ctx)
	})
}

// http
func (s *antsScheduler) SubmitWithHTTP(ctx context.Context, task func(context.Context, *gin.Context)) error {
	if ctx == nil {
		ctx = context.Background()
	}

	return s.pool.Submit(func() {
		// Check if context is already cancelled before executing
		select {
		case <-ctx.Done():
			return
		default:
		}
		task(ctx, nil)
	})
}

// SubmitWithTimeout submits a task with a specific timeout.
// A new context with timeout is derived from the provided context.
func (s *antsScheduler) SubmitWithTimeout(ctx context.Context, timeout time.Duration, task func(context.Context)) error {
	if ctx == nil {
		ctx = context.Background()
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)

	return s.pool.Submit(func() {
		defer cancel()
		// Check if context is already cancelled before executing
		select {
		case <-timeoutCtx.Done():
			return
		default:
		}
		task(timeoutCtx)
	})
}

// Shutdown gracefully shuts down the scheduler.
// It waits for all running tasks to complete or until context is cancelled.
func (s *antsScheduler) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		s.pool.Release()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Running returns the number of currently running goroutines.
func (s *antsScheduler) Running() int {
	return s.pool.Running()
}

// Cap returns the capacity of the goroutine pool.
func (s *antsScheduler) Cap() int {
	return s.pool.Cap()
}
