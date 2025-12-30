package daemon

import "context"

// Daemon 定义守护进程的接口
// 所有长期运行的服务都必须实现这个接口
//
// 什么是守护进程(Daemon)?
//   - 守护进程是在后台长期运行的程序
//   - 它们通常在系统启动时启动,在系统关闭时停止
//   - 在我们的项目中,HTTP 服务器、gRPC 服务器、消息队列消费者等都是守护进程
//
// 为什么使用接口?
//   - 抽象:不同类型的服务可以有不同的实现,但都遵循相同的接口
//   - 扩展:添加新的服务类型只需实现这个接口
//   - 测试:可以创建 mock 实现进行单元测试
//   - 统一管理:Manager 可以统一管理所有实现了这个接口的服务
//
// 使用示例:
//
//	// 实现一个简单的守护进程
//	type MyDaemon struct {
//	    name string
//	}
//
//	func (d *MyDaemon) Name() string {
//	    return d.name
//	}
//
//	func (d *MyDaemon) Start(ctx context.Context) error {
//	    // 启动服务的逻辑
//	    // 通常会在这里启动一个长期运行的 goroutine
//	    go func() {
//	        for {
//	            select {
//	            case <-ctx.Done():
//	                // 上下文被取消,退出循环
//	                return
//	            default:
//	                // 执行服务的主要工作
//	            }
//	        }
//	    }()
//	    return nil
//	}
//
//	func (d *MyDaemon) Stop(ctx context.Context) error {
//	    // 停止服务的逻辑
//	    // 进行清理工作,释放资源
//	    return nil
//	}
type Daemon interface {
	// Name 返回守护进程的名称
	// 这个名称用于:
	//   - 日志记录:方便识别是哪个服务
	//   - 监控统计:追踪各个服务的状态
	//   - 错误报告:明确指出出错的服务
	//
	// 命名建议:
	//   - 使用小写字母和连字符
	//   - 简短但有意义
	//   - 例如: "http-server", "grpc-server", "kafka-consumer"
	//
	// 返回:
	//   string: 守护进程的名称
	Name() string

	// Start 启动守护进程
	// 这个方法应该:
	//   1. 初始化服务所需的资源
	//   2. 启动服务(通常在独立的 goroutine 中)
	//   3. 尊重 context 的取消信号
	//
	// 参数:
	//   ctx: 上下文,用于控制启动过程
	//     - 如果 ctx 被取消,应该立即停止启动过程
	//     - 应该将 ctx 传递给服务的主循环,以便接收停止信号
	//
	// 返回:
	//   error: 启动失败时返回错误,成功则返回 nil
	//
	// 注意事项:
	//   - Start 方法应该快速返回,不要阻塞太久
	//   - 如果服务需要长期运行,应该在后台 goroutine 中进行
	//   - 必须检查 ctx.Done() 以支持取消操作
	//   - 启动失败要返回有意义的错误信息
	//
	// 实现示例:
	//   func (d *HTTPDaemon) Start(ctx context.Context) error {
	//       // 在后台启动 HTTP 服务器
	//       go func() {
	//           if err := d.server.ListenAndServe(); err != nil {
	//               log.Error("server error", "error", err)
	//           }
	//       }()
	//
	//       // 等待服务器准备就绪或上下文取消
	//       select {
	//       case <-d.ready:
	//           return nil
	//       case <-ctx.Done():
	//           return ctx.Err()
	//       }
	//   }
	Start(ctx context.Context) error

	// Stop 停止守护进程
	// 这个方法应该:
	//   1. 停止接收新的请求/任务
	//   2. 等待正在处理的请求/任务完成
	//   3. 清理资源(关闭连接、释放内存等)
	//   4. 在超时前完成所有工作
	//
	// 参数:
	//   ctx: 上下文,用于控制停止过程的超时
	//     - 如果 ctx 超时,应该强制停止
	//     - 通常会设置一个合理的超时时间(如 30 秒)
	//
	// 返回:
	//   error: 停止过程中发生错误时返回,成功则返回 nil
	//
	// 注意事项:
	//   - 这是优雅关闭(Graceful Shutdown)的关键
	//   - 必须在 ctx 超时前完成所有清理工作
	//   - 应该先停止接收新请求,再处理现有请求
	//   - 清理所有资源,避免资源泄漏
	//
	// 实现示例:
	//   func (d *HTTPDaemon) Stop(ctx context.Context) error {
	//       // 优雅关闭 HTTP 服务器
	//       // Shutdown 会停止接收新请求,等待现有请求处理完成
	//       if err := d.server.Shutdown(ctx); err != nil {
	//           return fmt.Errorf("failed to shutdown server: %w", err)
	//       }
	//
	//       // 关闭其他资源
	//       d.cleanup()
	//
	//       return nil
	//   }
	Stop(ctx context.Context) error
}
