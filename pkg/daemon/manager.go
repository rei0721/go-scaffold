package daemon

import (
	"context"
	"fmt"
	"sync"
)

// Logger 日志接口
// 定义 Manager 需要的日志方法
// 这样可以兼容不同的日志库(zap, logrus, slog 等)
type Logger interface {
	// Info 记录信息级别日志
	Info(msg string, keysAndValues ...interface{})
	// Error 记录错误级别日志
	Error(msg string, keysAndValues ...interface{})
}

// Manager 守护进程管理器
// 负责管理所有注册的守护进程
//
// 职责:
//   - 注册守护进程
//   - 启动所有守护进程
//   - 停止所有守护进程
//   - 记录日志
//   - 处理错误
//
// 为什么需要 Manager?
//   - 统一管理:不需要在多个地方手动启动/停止服务
//   - 并发控制:可以并发启动/停止多个服务,提高效率
//   - 错误处理:统一处理和记录错误
//   - 优雅关闭:确保所有服务都正确关闭
//
// 使用示例:
//
//	// 创建管理器
//	manager := daemon.NewManager(logger)
//
//	// 注册多个守护进程
//	manager.Register(httpDaemon)
//	manager.Register(grpcDaemon)
//	manager.Register(kafkaDaemon)
//
//	// 启动所有守护进程
//	if err := manager.Start(ctx); err != nil {
//	    log.Fatal(err)
//	}
//
//	// 优雅关闭
//	if err := manager.Stop(ctx); err != nil {
//	    log.Error(err)
//	}
type Manager struct {
	// daemons 存储所有注册的守护进程
	// 使用切片而不是 map 的原因:
	//   - 保持注册顺序
	//   - 按顺序启动和停止
	daemons []Daemon

	// logger 日志记录器
	// 用于记录管理器的操作日志
	logger Logger

	// mu 互斥锁
	// 保护 daemons 切片,确保线程安全
	// 为什么需要锁?
	//   - Register 可能在不同 goroutine 中调用
	//   - 需要保证对 daemons 的并发访问是安全的
	mu sync.Mutex
}

// NewManager 创建一个新的守护进程管理器
//
// 参数:
//
//	logger: 日志记录器,用于记录管理器的操作
//	  - 如果为 nil,将会 panic
//	  - 必须实现 Logger 接口
//
// 返回:
//
//	*Manager: 守护进程管理器实例
//
// 使用示例:
//
//	logger := logger.New()
//	manager := daemon.NewManager(logger)
func NewManager(logger Logger) *Manager {
	return &Manager{
		daemons: make([]Daemon, 0),
		logger:  logger,
	}
}

// Register 注册一个守护进程
// 注册后的守护进程会在调用 Start 时启动
//
// 参数:
//
//	daemon: 要注册的守护进程
//	  - 必须实现 Daemon 接口
//	  - 不能为 nil
//
// 注意事项:
//   - 应该在调用 Start 之前注册所有守护进程
//   - 可以注册多个守护进程
//   - 守护进程会按注册顺序启动
//   - 这个方法是线程安全的
//
// 使用示例:
//
//	httpDaemon := daemons.NewHTTPDaemon(":8080", router)
//	manager.Register(httpDaemon)
//
//	grpcDaemon := daemons.NewGRPCDaemon(":9090")
//	manager.Register(grpcDaemon)
func (m *Manager) Register(daemon Daemon) {
	// 加锁保证线程安全
	// 防止多个 goroutine 同时修改 daemons 切片
	m.mu.Lock()
	defer m.mu.Unlock()

	// 将守护进程添加到列表中
	m.daemons = append(m.daemons, daemon)

	// 记录注册日志
	m.logger.Info(MsgDaemonRegistered, "name", daemon.Name())
}

// Start 启动所有注册的守护进程
// 所有守护进程会并发启动,提高启动速度
//
// 参数:
//
//	ctx: 上下文,用于控制启动过程
//	  - 如果 ctx 被取消,会停止启动过程
//	  - 建议使用 context.Background() 或有超时的 context
//
// 返回:
//
//	error: 如果任何一个守护进程启动失败,返回错误
//	  - 一旦有守护进程失败,会立即返回
//	  - 其他守护进程可能已经启动成功
//
// 工作流程:
//  1. 记录开始启动的日志
//  2. 为每个守护进程创建独立的 goroutine
//  3. 并发调用每个守护进程的 Start 方法
//  4. 等待所有守护进程启动完成或出错
//  5. 记录启动结果
//
// 注意事项:
//   - 如果某个守护进程启动失败,其他守护进程可能已经启动
//   - 出错时应该调用 Stop 清理已启动的守护进程
//   - 这个方法会阻塞直到所有守护进程启动完成
//
// 使用示例:
//
//	ctx := context.Background()
//	if err := manager.Start(ctx); err != nil {
//	    log.Error("failed to start daemons", "error", err)
//	    // 清理已启动的守护进程
//	    manager.Stop(context.Background())
//	    return err
//	}
func (m *Manager) Start(ctx context.Context) error {
	m.logger.Info(MsgAllDaemonsStarting, "count", len(m.daemons))

	// 创建错误通道
	// 缓冲大小设置为守护进程数量,避免 goroutine 阻塞
	// 为什么使用通道?
	//   - 可以从多个 goroutine 收集错误
	//   - 一旦有错误就能立即知道
	errChan := make(chan error, len(m.daemons))

	// 并发启动所有守护进程
	// 为什么要并发?
	//   - 提高启动速度
	//   - 多个服务可以同时初始化
	for _, daemon := range m.daemons {
		// 注意:必须在循环内部创建新变量
		// 避免闭包陷阱(所有 goroutine 共享同一个 daemon 变量)
		d := daemon

		// 在独立的 goroutine 中启动守护进程
		go func() {
			m.logger.Info(MsgDaemonStarting, "name", d.Name())

			// 调用守护进程的 Start 方法
			if err := d.Start(ctx); err != nil {
				// 启动失败,记录错误日志
				m.logger.Error(MsgDaemonStartFailed, "name", d.Name(), "error", err)
				// 将错误发送到通道
				errChan <- fmt.Errorf(ErrMsgDaemonStartFailed, d.Name(), err)
			} else {
				// 启动成功,记录信息日志
				m.logger.Info(MsgDaemonStarted, "name", d.Name())
			}
		}()
	}

	// 等待所有守护进程启动完成或出错
	// 使用 select 同时监听错误和上下文取消
	select {
	case err := <-errChan:
		// 有守护进程启动失败
		return err
	case <-ctx.Done():
		// 上下文被取消(可能是超时)
		return ctx.Err()
	default:
		// 所有守护进程启动成功
		m.logger.Info(MsgAllDaemonsStarted)
		return nil
	}
}

// Stop 停止所有守护进程
// 所有守护进程会并发停止,提高关闭速度
//
// 参数:
//
//	ctx: 上下文,用于控制停止过程的超时
//	  - 建议设置超时(如 30 秒)
//	  - 超时后会强制返回,可能有守护进程未完全停止
//
// 返回:
//
//	error: 如果停止过程超时,返回错误
//	  - 即使返回错误,也会尽力停止所有守护进程
//
// 工作流程:
//  1. 记录开始停止的日志
//  2. 为每个守护进程创建独立的 goroutine
//  3. 并发调用每个守护进程的 Stop 方法
//  4. 等待所有守护进程停止完成或超时
//  5. 记录停止结果
//
// 注意事项:
//   - 必须设置合理的超时时间
//   - 超时后可能有守护进程未完全停止
//   - 即使某个守护进程停止失败,其他守护进程仍会继续停止
//   - 这是优雅关闭的关键步骤
//
// 使用示例:
//
//	// 创建带 30 秒超时的上下文
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//
//	if err := manager.Stop(ctx); err != nil {
//	    log.Error("failed to stop daemons", "error", err)
//	}
func (m *Manager) Stop(ctx context.Context) error {
	m.logger.Info(MsgAllDaemonsStopping, "count", len(m.daemons))

	// 使用 WaitGroup 等待所有守护进程停止
	// WaitGroup 用于同步多个 goroutine
	var wg sync.WaitGroup

	// 创建错误通道收集停止过程中的错误
	errChan := make(chan error, len(m.daemons))

	// 并发停止所有守护进程
	for _, daemon := range m.daemons {
		// 增加 WaitGroup 计数
		wg.Add(1)

		// 避免闭包陷阱
		d := daemon

		// 在独立的 goroutine 中停止守护进程
		go func() {
			// goroutine 结束时减少 WaitGroup 计数
			defer wg.Done()

			m.logger.Info(MsgDaemonStopping, "name", d.Name())

			// 调用守护进程的 Stop 方法
			if err := d.Stop(ctx); err != nil {
				// 停止失败,记录错误
				m.logger.Error(MsgDaemonStopFailed, "name", d.Name(), "error", err)
				errChan <- err
			} else {
				// 停止成功
				m.logger.Info(MsgDaemonStopped, "name", d.Name())
			}
		}()
	}

	// 创建一个通道,用于通知所有守护进程已停止
	done := make(chan struct{})

	// 在后台等待所有守护进程停止
	go func() {
		// 等待所有 goroutine 完成
		wg.Wait()
		// 关闭 done 通道,表示所有守护进程已停止
		close(done)
	}()

	// 等待所有守护进程停止或超时
	select {
	case <-done:
		// 所有守护进程已成功停止
		m.logger.Info(MsgAllDaemonsStopped)
		return nil
	case <-ctx.Done():
		// 上下文超时或被取消
		// 注意:即使超时,后台的 goroutine 仍会继续执行
		return fmt.Errorf("%s: %w", MsgStopTimeout, ctx.Err())
	}
}
