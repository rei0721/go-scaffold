// Package main 是应用程序的入口点
// 它负责初始化 App 容器并处理操作系统信号的优雅关闭
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rei0721/rei0721/internal/app"
	"github.com/rei0721/rei0721/pkg/cli"
	"github.com/rei0721/rei0721/types/constants"
)

const (
	// defaultConfigPath 是配置文件的默认路径
	// 当环境变量 CONFIG_PATH 未设置时使用此路径
	defaultConfigPath = constants.DefaultConfigPath

	// shutdownTimeout 是优雅关闭的最大等待时间
	// 设置为 30 秒以确保所有正在处理的请求能够完成
	// 超过此时间后将强制关闭,避免进程无限期等待
	shutdownTimeout = constants.ShutdownTimeout
)

type DevCommand struct{}

func (c *DevCommand) Name() string {
	return "dev"
}

func (c *DevCommand) Description() string {
	return "Dev mode"
}

func (c *DevCommand) Usage() string {
	return "dev [--config=<name>]"
}

func (c *DevCommand) Flags() []cli.Flag {
	return []cli.Flag{
		{
			Name:        "config",
			ShortName:   "c",
			Type:        cli.FlagTypeString,
			Default:     "configs/config.yaml",
			Description: "Config file path",
			EnvVar:      "REI_CONFIG_PATH", // 支持环境变量
		},
	}
}

func (c *DevCommand) Execute(ctx *cli.Context) error {
	configPath := ctx.GetString("config")

	fmt.Fprintf(ctx.Stdout, "Dev Config Path = %s\n",
		configPath)

	// 业务逻辑...
	run()

	return nil
}

func main() {
	// 创建 CLI 应用
	app := cli.NewApp("go-scaffold")
	app.SetVersion("0.1.2")
	app.SetDescription("This is a go backend scaffold")

	// 注册命令
	app.AddCommand(&DevCommand{})

	// 执行
	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cli.GetExitCode(err))
	}

}

func run() {

	// 1. 获取配置文件路径
	// 首先尝试从环境变量 `constants.EnvConfigPathName` 读取配置文件路径
	// 这样做是为了支持不同环境(开发、测试、生产)使用不同的配置文件
	configPath := os.Getenv(constants.EnvConfigPathName)
	if configPath == "" {
		// 如果环境变量未设置,使用默认路径
		// 这提供了一个合理的默认值,简化了本地开发
		configPath = defaultConfigPath
	}

	// 2. 初始化应用程序容器
	// app.New() 会按照依赖顺序初始化所有组件
	// 使用依赖注入容器模式来管理组件生命周期
	application, err := app.New(app.Options{
		ConfigPath: configPath,
	})
	if err != nil {
		// 此时 logger 可能还未初始化,因此直接使用 os.Stderr 输出错误
		// 使用 os.Stderr 而不是 os.Stdout,确保错误信息不会被重定向丢失
		os.Stderr.WriteString("failed to initialize application: " + err.Error() + "\n")
		// 退出码 1 表示程序异常退出,便于运维监控和部署脚本判断
		os.Exit(1)
	}

	// 3. 创建信号接收通道
	// 缓冲区大小为 1,确保至少能接收一个信号而不被阻塞
	// 这样可以避免在处理第一个信号时丢失后续信号
	quit := make(chan os.Signal, 1)

	// 4. 注册要监听的操作系统信号
	// SIGINT: 用户按下 Ctrl+C 时发送的中断信号
	// SIGTERM: 系统或容器编排工具(如 Docker/K8s)发送的终止信号
	// 监听这些信号是实现优雅关闭的关键,确保应用能够正确清理资源
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 5. 在新的 goroutine 中启动 HTTP 服务器
	// 使用 goroutine 使得主线程可以继续执行,监听关闭信号
	// 创建缓冲通道用于接收服务器运行时的错误
	errChan := make(chan error, 1)
	go func() {
		// application.Run() 会阻塞直到服务器停止
		// 如果启动失败或运行时出错,将错误发送到 errChan
		if err := application.Run(); err != nil {
			errChan <- err
		}
	}()

	// 6. 等待关闭信号或服务器错误
	// 使用 select 同时监听多个通道,哪个先有数据就处理哪个
	// 这是 Go 语言中常用的并发模式
	select {
	case sig := <-quit:
		// 收到操作系统信号(用户中断或系统终止)
		// 记录日志以便问题追踪和审计
		application.Logger.Info("received shutdown signal", "signal", sig.String())
	case err := <-errChan:
		// 服务器运行时发生错误(例如端口已被占用)
		// 记录错误日志,便于排查问题
		application.Logger.Error("server error", "error", err)
	}

	// 7. 创建带超时的上下文
	// 上下文用于控制优雅关闭的最大等待时间
	// 如果超过 shutdownTimeout(30秒)仍未完成关闭,将强制退出
	// 这防止了因某些请求或连接无法正常关闭导致的程序挂起
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	// defer 确保在函数返回时调用 cancel,释放上下文相关资源
	// 这是 Go 中管理资源的最佳实践
	defer cancel()

	// 8. 执行优雅关闭
	// 优雅关闭的顺序: HTTP服务器 → 调度器 → 数据库连接 → 日志同步
	// 这个顺序确保:
	// - 先停止接收新请求
	// - 等待正在执行的任务完成
	// - 关闭数据库连接
	// - 最后同步日志,确保所有日志都写入磁盘
	if err := application.Shutdown(ctx); err != nil {
		// 如果优雅关闭失败,记录错误并以非零状态退出
		// 退出码 1 表示程序异常退出
		application.Logger.Error("shutdown error", "error", err)
		os.Exit(1)
	}

	// 9. 记录成功退出的日志
	// 到达这里说明所有资源都已正确清理,程序正常退出
	application.Logger.Info("application exited gracefully")
}
