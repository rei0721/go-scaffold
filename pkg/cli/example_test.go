package cli_test

import (
	"bytes"
	"context"
	"fmt"

	"github.com/rei0721/rei0721/pkg/cli"
)

// Example_basicUsage 演示基本的 CLI 工具创建
func Example_basicUsage() {
	// 定义命令
	type GreetCommand struct {
		name *string
	}

	cmd := &GreetCommand{}

	// 实现 Command 接口
	cmd.Name = func() string {
		return "greet"
	}
	cmd.Description = func() string {
		return "向用户问候"
	}
	cmd.Run = func(ctx context.Context, args []string) error {
		fmt.Printf("Hello, %s!\n", *cmd.name)
		return nil
	}

	// 创建运行器
	runner := cli.NewSimpleRunner(&simpleCommand{
		name:        "greet",
		description: "向用户问候",
		runFunc: func(ctx context.Context, args []string) error {
			fmt.Println("Hello, World!")
			return nil
		},
	})

	// 运行
	exitCode := runner.Run(context.Background(), []string{"greet"})
	fmt.Printf("Exit code: %d\n", exitCode)

	// Output:
	// Hello, World!
	// Exit code: 0
}

// Example_withParameters 演示带参数的 CLI 工具
func Example_withParameters() {
	type ConvertCommand struct {
		input  *string
		output *string
	}

	cmd := &ConvertCommand{}
	runner := cli.NewSimpleRunner(&simpleCommand{
		name:        "convert",
		description: "转换文件",
		runFunc: func(ctx context.Context, args []string) error {
			fmt.Printf("Converting %s to %s\n", *cmd.input, *cmd.output)
			return nil
		},
	})

	// 定义参数
	parser := runner.GetParser()
	cmd.input = parser.String("input", "i", "", "输入文件")
	cmd.output = parser.String("output", "o", "", "输出文件")

	// 模拟命令行参数
	exitCode := runner.Run(context.Background(), []string{"convert", "-i", "input.txt", "-o", "output.txt"})
	fmt.Printf("Exit code: %d\n", exitCode)

	// Output:
	// Converting input.txt to output.txt
	// Exit code: 0
}

// Example_validation 演示参数验证
func Example_validation() {
	runner := cli.NewSimpleRunner(&simpleCommand{
		name:        "server",
		description: "启动服务器",
		runFunc: func(ctx context.Context, args []string) error {
			fmt.Println("Server started")
			return nil
		},
	})

	parser := runner.GetParser()
	port := parser.Int("port", "p", 8080, "监听端口")

	// 添加验证器
	parser.SetValidator("port", func(v interface{}) error {
		p := v.(int)
		if p < 1 || p > 65535 {
			return fmt.Errorf("端口必须在 1-65535 之间")
		}
		return nil
	})

	// 使用有效端口
	exitCode := runner.Run(context.Background(), []string{"server", "--port", "8080"})
	fmt.Printf("Valid port - Exit code: %d\n", exitCode)

	// Output:
	// Server started
	// Valid port - Exit code: 0
}

// Example_helpFormatter 演示自定义帮助格式
func Example_helpFormatter() {
	cmd := &simpleCommand{
		name:        "myapp",
		description: "我的应用程序",
		runFunc: func(ctx context.Context, args []string) error {
			return nil
		},
	}

	var buf bytes.Buffer
	runner := cli.NewSimpleRunner(cmd, cli.WithStdout(&buf))

	// 显示帮助
	runner.Run(context.Background(), []string{"myapp", "--help"})

	// 帮助信息会包含命令名称和描述
	fmt.Println("Help displayed")

	// Output:
	// Help displayed
}

// Example_exitCodes 演示不同的退出码使用
func Example_exitCodes() {
	// 成功
	fmt.Printf("Success: %d\n", cli.ExitSuccess)

	// 错误
	fmt.Printf("Error: %d\n", cli.ExitError)

	// 无效使用
	fmt.Printf("Invalid usage: %d\n", cli.ExitInvalidUsage)

	// Output:
	// Success: 0
	// Error: 1
	// Invalid usage: 2
}

// Example_constants 演示常量的使用
func Example_constants() {
	// 使用指示器
	fmt.Printf("%s 操作成功\n", cli.IndicatorSuccess)
	fmt.Printf("%s 操作失败\n", cli.IndicatorError)
	fmt.Printf("%s 警告信息\n", cli.IndicatorWarning)

	// Output:
	// ✅ 操作成功
	// ❌ 操作失败
	// ⚠️ 警告信息
}

// simpleCommand 是一个简单的命令实现,用于示例
type simpleCommand struct {
	name        string
	description string
	runFunc     func(context.Context, []string) error
}

func (c *simpleCommand) Name() string {
	return c.name
}

func (c *simpleCommand) Description() string {
	return c.description
}

func (c *simpleCommand) Run(ctx context.Context, args []string) error {
	return c.runFunc(ctx, args)
}
