// Package cli 提供统一的命令行工具开发框架
//
// # 概述
//
// cli 包封装了命令行工具开发的常用功能,包括:
//   - 参数解析和验证
//   - 命令执行和错误处理
//   - 帮助文档生成
//   - 版本信息管理
//
// # 设计目标
//
//   - 统一常量定义 - 避免魔法字符串,便于维护
//   - 清晰的抽象层 - Command、Parser、Runner 接口
//   - 类型安全 - 强类型的参数访问
//   - 易于测试 - 依赖注入,可替换 IO
//   - 详细注释 - 每个导出类型都有完整说明
//
// # 核心概念
//
// Command (命令)
//
// Command 接口定义了一个可执行的命令。每个命令需要实现:
//   - Name() - 命令名称
//   - Description() - 命令描述
//   - Run() - 执行逻辑
//
// Parser (参数解析器)
//
// Parser 提供类型安全的命令行参数解析:
//   - String/Int/Bool - 基本类型
//   - StringSlice - 切片类型
//   - SetValidator - 自定义验证
//   - SetRequired - 必需参数
//
// Runner (运行器)
//
// Runner 管理命令的完整生命周期:
//   - 解析参数
//   - 处理内置标志(help, version)
//   - 执行命令
//   - 错误处理
//
// # 快速开始
//
// 创建一个简单的 CLI 工具:
//
//	package main
//
//	import (
//	    "context"
//	    "fmt"
//	    "os"
//
//	    "github.com/rei0721/rei0721/pkg/cli"
//	)
//
//	// 定义命令
//	type GreetCommand struct {
//	    name *string
//	}
//
//	func (c *GreetCommand) Name() string {
//	    return "greet"
//	}
//
//	func (c *GreetCommand) Description() string {
//	    return "向用户问候"
//	}
//
//	func (c *GreetCommand) Run(ctx context.Context, args []string) error {
//	    fmt.Printf("Hello, %s!\n", *c.name)
//	    return nil
//	}
//
//	func main() {
//	    // 创建命令
//	    cmd := &GreetCommand{}
//
//	    // 创建运行器
//	    runner := cli.NewSimpleRunner(cmd)
//
//	    // 定义参数
//	    parser := runner.GetParser()
//	    cmd.name = parser.String("name", "n", "World", "要问候的名字")
//
//	    // 运行
//	    exitCode := runner.Run(context.Background(), os.Args)
//	    os.Exit(exitCode)
//	}
//
// 使用:
//
//	go run main.go --name Alice
//	# 输出: Hello, Alice!
//
//	go run main.go -h
//	# 显示帮助信息
//
// # 常量定义
//
// 退出码常量:
//   - ExitSuccess (0) - 成功
//   - ExitError (1) - 一般错误
//   - ExitInvalidUsage (2) - 无效使用
//   - ExitInvalidConfig (3) - 无效配置
//   - ExitIOError (4) - IO 错误
//
// 默认标志:
//   - DefaultHelpFlag ("help")
//   - DefaultVersionFlag ("version")
//   - DefaultVerboseFlag ("verbose")
//
// 日志消息:
//   - MsgCommandStarting
//   - MsgCommandCompleted
//   - MsgCommandFailed
//
// 错误消息:
//   - ErrMsgInvalidFlag
//   - ErrMsgMissingRequired
//   - ErrMsgCommandFailedFmt
//
// # 最佳实践
//
// 1. 使用常量避免魔法字符串
//
//	// ✅ 好的做法
//	os.Exit(cli.ExitSuccess)
//
//	// ❌ 不好的做法
//	os.Exit(0)
//
// 2. 验证用户输入
//
//	port := parser.Int("port", "p", 8080, "监听端口")
//	parser.SetValidator("port", func(v interface{}) error {
//	    if p := v.(int); p < 1 || p > 65535 {
//	        return fmt.Errorf("端口必须在 1-65535 之间")
//	    }
//	    return nil
//	})
//
// 3. 使用上下文支持取消
//
//	func (c *MyCommand) Run(ctx context.Context, args []string) error {
//	    select {
//	    case <-ctx.Done():
//	        return ctx.Err()
//	    default:
//	        // 执行工作...
//	    }
//	    return nil
//	}
//
// 4. 提供清晰的错误消息
//
//	if err != nil {
//	    return fmt.Errorf("读取配置失败: %w", err)
//	}
//
// # 参考链接
//
//   - Go flag 包: https://pkg.go.dev/flag
//   - CLI 最佳实践: https://clig.dev/
package cli
