// Package daemon 提供守护进程管理功能
// 用于统一管理所有长期运行的服务(如 HTTP 服务器、gRPC 服务器、消息队列消费者等)
//
// # 设计目标
//
// - 统一管理长期运行的服务
// - 支持优雅启动和关闭
// - 提供服务状态监控
// - 便于扩展新的服务类型
//
// # 核心概念
//
// Daemon(守护进程):
//   - 长期运行的后台服务
//   - 通常不会自动退出,除非收到停止信号
//   - 例如: HTTP 服务器、gRPC 服务器、Kafka 消费者等
//
// Manager(管理器):
//   - 统一管理多个 Daemon
//   - 负责启动、停止所有 Daemon
//   - 处理错误和超时
//
// # 使用示例
//
// 创建管理器:
//
//	manager := daemon.NewManager(logger)
//
// 注册守护进程:
//
//	httpDaemon := &HTTPDaemon{...}
//	manager.Register(httpDaemon)
//
// 启动所有守护进程:
//
//	ctx := context.Background()
//	if err := manager.Start(ctx); err != nil {
//	    log.Fatal(err)
//	}
//
// 优雅关闭:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//	if err := manager.Stop(ctx); err != nil {
//	    log.Error("shutdown error", "error", err)
//	}
//
// # 与其他包的区别
//
// pkg/daemon:
//   - 管理**长期运行的服务**(HTTP、gRPC、Kafka 等)
//   - 服务通常不会自动退出
//   - 需要优雅关闭机制
//
// pkg/scheduler:
//   - 管理**短期异步任务**(发送邮件、记录日志等)
//   - 任务完成后 goroutine 会返回池中
//   - 基于 ants 协程池
//
// internal/service:
//   - **业务逻辑层**
//   - 处理具体的业务功能
//   - 不涉及服务生命周期管理
//
// # 最佳实践
//
// 1. 所有长期运行的服务都应该实现 Daemon 接口
// 2. 在 App 初始化时注册所有 Daemon
// 3. 使用 Manager 统一启动和停止
// 4. 设置合理的超时时间
// 5. 在 Daemon 实现中正确处理 context 取消信号
//
// # 线程安全
//
// Manager 是线程安全的,可以在多个 goroutine 中安全调用
package daemon
