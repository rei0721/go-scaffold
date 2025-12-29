// Package scheduler 提供协程池管理的统一接口
// internal/ 中的所有异步任务都必须通过这个调度器提交
// 设计目标:
// - 统一管理异步任务
// - 限制并发数量,防止资源耗尽
// - 支持优雅关闭,确保任务完成
// 为什么需要调度器:
// - 直接创建 goroutine 会导致资源不可控
// - 协程池可以复用 goroutine,减少创建/销毁开销
// - 便于监控和限流
package scheduler

import (
	"context"
	"time"
)

// Scheduler 定义协程池管理的接口
// 所有异步任务都必须通过这个接口提交,并传递 context.Context
// 为什么使用接口:
// - 抽象实现:可以使用不同的协程池实现(ants、自定义等)
// - 便于测试:可以 mock 调度器进行测试
// - 统一管理:强制所有异步任务使用统一的调度器
type Scheduler interface {
	// Submit 提交一个任务到协程池
	// 任务函数会接收一个 context,应该尊重上下文的取消信号
	// 参数:
	//   ctx: 上下文,用于:
	//     - 传递给任务函数
	//     - 任务应该检查 ctx.Done() 以支持取消
	//   task: 任务函数,接收 context 参数
	//     - 应该定期检查 ctx.Done()
	//     - 长时间运行的任务应该支持中断
	// 返回:
	//   error: 提交失败时的错误,如:
	//     - 协程池已满且队列已满(非阻塞模式)
	//     - 调度器已关闭
	// 使用示例:
	//   scheduler.Submit(ctx, func(taskCtx context.Context) {
	//       // 执行异步任务
	//       select {
	//       case <-taskCtx.Done():
	//           // 任务被取消,清理并退出
	//           return
	//       default:
	//           // 执行实际工作
	//       }
	//   })
	// 使用场景:
	//   - 发送邮件通知
	//   - 记录审计日志
	//   - 更新缓存
	//   - 触发其他服务
	Submit(ctx context.Context, task func(context.Context)) error

	// SubmitWithTimeout 提交一个带超时的任务
	// 如果任务执行时间超过 timeout,会被自动取消
	// 参数:
	//   ctx: 父上下文
	//   timeout: 超时时间
	//     - 从任务开始执行计时
	//     - 超时后 taskCtx.Done() 会被关闭
	//   task: 任务函数
	// 返回:
	//   error: 提交失败时的错误
	// 使用示例:
	//   scheduler.SubmitWithTimeout(ctx, 5*time.Second, func(taskCtx context.Context) {
	//       // 这个任务最多执行 5 秒
	//       // 超时后 taskCtx.Done() 会被关闭
	//   })
	// 使用场景:
	//   - 调用外部 API(设置超时防止长时间等待)
	//   - 执行可能耗时的操作
	//   - 需要保证任务在限定时间内完成
	// 注意:
	//   - 任务函数必须尊重 context 的取消信号
	//   - 超时不会强制终止 goroutine,只是关闭 context
	//   - 如果任务不检查 ctx.Done(),仍然会继续执行
	SubmitWithTimeout(ctx context.Context, timeout time.Duration, task func(context.Context)) error

	// Shutdown 优雅关闭调度器
	// 等待所有正在运行的任务完成,或直到上下文被取消
	// 参数:
	//   ctx: 上下文,控制关闭超时
	//     - ctx.Done() 被关闭时,停止等待并返回
	//     - 建议设置超时(如 30 秒)
	// 返回:
	//   error: 关闭失败或超时时的错误
	// 关闭流程:
	//   1. 停止接收新任务
	//   2. 等待所有运行中的任务完成
	//   3. 释放协程池资源
	// 使用示例:
	//   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//   defer cancel()
	//   if err := scheduler.Shutdown(ctx); err != nil {
	//       log.Error("scheduler shutdown error", "error", err)
	//   }
	// 注意:
	//   - 应该在应用关闭时调用
	//   - 等待时间过长可能导致关闭超时
	//   - 如果上下文超时,可能有任务未完成
	Shutdown(ctx context.Context) error

	// Running 返回当前正在运行的 goroutine 数量
	// 用途:
	//   - 监控系统负载
	//   - 判断系统是否过载
	//   - 调试和性能分析
	// 返回:
	//   int: 当前活跃的 goroutine 数量
	Running() int

	// Cap 返回协程池的容量
	// 即最大可以同时运行的 goroutine 数量
	// 用途:
	//   - 了解系统配置
	//   - 计算资源利用率(Running/Cap)
	//   - 监控和告警
	// 返回:
	//   int: 协程池容量
	Cap() int
}

// Config 保存调度器的配置
// 包含协程池初始化所需的参数
type Config struct {
	// PoolSize 协程池中的最大 goroutine 数量
	// 这个值决定了系统的最大并发能力
	// 设置考虑:
	//   - 过大:消耗过多内存和 CPU
	//   - 过小:任务会长时间等待
	//   - 推荐:根据任务类型调整
	//     * CPU 密集型: CPU 核心数的 1-2 倍
	//     * IO 密集型: CPU 核心数的 10-100 倍
	// 默认: 10000(适合 IO 密集型任务)
	PoolSize int

	// MaxBlockingTasks 最大阻塞任务数
	// 当协程池满时,可以等待的任务队列大小
	// 设置说明:
	//   - 0: 无限制(任务会一直等待直到有可用 goroutine)
	//   - >0: 限制等待队列大小,超过会返回错误
	// 默认: 0(无限制)
	// 使用场景:
	//   - 防止内存被等待队列耗尽
	//   - 对于非关键任务,可以设置限制并丢弃
	MaxBlockingTasks int

	// ExpiryDuration 空闲 goroutine 的过期时间
	// 超过这个时间没有任务的 goroutine 会被回收
	// 好处:
	//   - 释放不需要的 goroutine,节省内存
	//   - 动态调整池大小,适应负载变化
	// 默认: 1 秒
	// 推荐:
	//   - 高负载系统:1-5 秒(快速回收)
	//   - 低负载系统:10-60 秒(减少创建/销毁)
	ExpiryDuration time.Duration
}

// DefaultConfig 返回一个使用合理默认值的配置
// 这些默认值适合大多数场景
// 返回:
//
//	*Config: 默认配置
//	  - PoolSize: 10000 (支持大量 IO 密集型任务)
//	  - MaxBlockingTasks: 0 (无限制,确保任务不丢失)
//	  - ExpiryDuration: 1 秒 (快速回收空闲 goroutine)
//
// 使用示例:
//
//	cfg := scheduler.DefaultConfig()
//	// 可以根据需要调整
//	cfg.PoolSize = 5000
//	sched, err := scheduler.New(cfg)
func DefaultConfig() *Config {
	return &Config{
		PoolSize:         10000,
		MaxBlockingTasks: 0,
		ExpiryDuration:   time.Second,
	}
}
