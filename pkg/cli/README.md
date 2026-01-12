# CLI - 命令行工具开发框架

## 概述

CLI 是一个统一的命令行工具开发框架,提供参数解析、命令执行、帮助文档生成等功能,让你专注于业务逻辑。

### 特性

- ✅ **统一常量定义** - 避免魔法字符串,规范退出码、消息和标志
- ✅ **类型安全** - 强类型的参数访问(String/Int/Bool/StringSlice)
- ✅ **参数验证** - 内置验证器,支持自定义验证逻辑
- ✅ **清晰的抽象** - Command、Parser、Runner 接口分离
- ✅ **易于测试** - 可替换 IO,支持依赖注入
- ✅ **详细注释** - 完整的中文注释,适合初学者

## 快速开始

### 1. 创建命令

```go
package main

import (
    "context"
    "fmt"
    "github.com/rei0721/rei0721/pkg/cli"
)

type GreetCommand struct {
    name *string
}

func (c *GreetCommand) Name() string {
    return "greet"
}

func (c *GreetCommand) Description() string {
    return "向用户问候"
}

func (c *GreetCommand) Run(ctx context.Context, args []string) error {
    fmt.Printf("Hello, %s!\n", *c.name)
    return nil
}
```

### 2. 创建运行器并定义参数

```go
func main() {
    // 创建命令
    cmd := &GreetCommand{}

    // 创建运行器
    runner := cli.NewSimpleRunner(cmd)

    // 定义参数
    parser := runner.GetParser()
    cmd.name = parser.String("name", "n", "World", "要问候的名字")

    // 运行
    exitCode := runner.Run(context.Background(), os.Args)
    os.Exit(exitCode)
}
```

### 3. 使用

```bash
# 基本使用
go run main.go --name Alice
# 输出: Hello, Alice!

# 使用短标志
go run main.go -n Bob
# 输出: Hello, Bob!

# 显示帮助
go run main.go -h
go run main.go --help
```

## API 文档

### 常量

#### 退出码常量

| 常量                | 值  | 说明           |
| ------------------- | --- | -------------- |
| `ExitSuccess`       | 0   | 成功退出       |
| `ExitError`         | 1   | 一般错误       |
| `ExitInvalidUsage`  | 2   | 无效使用       |
| `ExitInvalidConfig` | 3   | 无效配置       |
| `ExitIOError`       | 4   | IO 错误        |
| `ExitInterrupted`   | 130 | 被中断(Ctrl+C) |

#### 默认标志

| 常量                 | 值        | 说明         |
| -------------------- | --------- | ------------ |
| `DefaultHelpFlag`    | "help"    | 帮助标志     |
| `DefaultVersionFlag` | "version" | 版本标志     |
| `DefaultVerboseFlag` | "verbose" | 详细输出标志 |

### Command 接口

```go
type Command interface {
    Name() string
    Description() string
    Run(ctx context.Context, args []string) error
}
```

**方法说明:**

- `Name()` - 返回命令名称,用于调用和帮助文档
- `Description()` - 返回命令描述,用于帮助文档
- `Run()` - 执行命令逻辑

### Parser 参数解析器

#### 定义参数

```go
parser := cli.NewParser("myapp")

// 字符串参数
output := parser.String("output", "o", "./output", "输出目录")

// 整数参数
port := parser.Int("port", "p", 8080, "监听端口")

// 布尔参数
verbose := parser.Bool("verbose", "v", false, "详细输出")

// 字符串切片参数
tags := parser.StringSlice("tag", "t", "", "标签(可多次指定)")
```

#### 参数验证

```go
// 设置必需参数
parser.SetRequired("input")

// 自定义验证器
parser.SetValidator("port", func(v interface{}) error {
    port := v.(int)
    if port < 1 || port > 65535 {
        return fmt.Errorf("端口必须在 1-65535 之间")
    }
    return nil
})

// 解析参数
if err := parser.Parse(os.Args[1:]); err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(cli.ExitInvalidUsage)
}
```

#### 获取位置参数

```go
parser.Parse(os.Args[1:])

// 获取所有位置参数
args := parser.Args()

// 获取第一个位置参数
firstArg := parser.Arg(0)

// 获取位置参数数量
count := parser.NArg()
```

### Runner 运行器

#### SimpleRunner

```go
// 创建运行器
runner := cli.NewSimpleRunner(cmd)

// 配置选项
runner := cli.NewSimpleRunner(cmd,
    cli.WithConfig(config),
    cli.WithStdout(os.Stdout),
    cli.WithStderr(os.Stderr),
)

// 获取解析器
parser := runner.GetParser()

// 运行命令
exitCode := runner.Run(context.Background(), os.Args)
```

### Config 配置

```go
// 使用默认配置
config := cli.DefaultConfig()
config.AppName = "myapp"
config.Version = "1.0.0"
config.Description = "我的应用"

// 验证配置
if err := config.Validate(); err != nil {
    log.Fatal(err)
}
```

## 使用场景

### 场景 1: 带参数的工具

```go
type ConvertCommand struct {
    input  *string
    output *string
    format *string
}

func (c *ConvertCommand) Name() string {
    return "convert"
}

func (c *ConvertCommand) Description() string {
    return "转换文件格式"
}

func main() {
    cmd := &ConvertCommand{}
    runner := cli.NewSimpleRunner(cmd)

    parser := runner.GetParser()
    cmd.input = parser.String("input", "i", "", "输入文件")
    cmd.output = parser.String("output", "o", "", "输出文件")
    cmd.format = parser.String("format", "f", "json", "输出格式")

    parser.SetRequired("input")
    parser.SetValidator("format", func(v interface{}) error {
        format := v.(string)
        if format != "json" && format != "yaml" {
            return fmt.Errorf("格式必须是 json 或 yaml")
        }
        return nil
    })

    exitCode := runner.Run(context.Background(), os.Args)
    os.Exit(exitCode)
}
```

### 场景 2: 文件处理工具

```go
type ProcessCommand struct {
    files   *[]string
    workers *int
    verbose *bool
}

func (c *ProcessCommand) Run(ctx context.Context, args []string) error {
    if *c.verbose {
        fmt.Printf("处理 %d 个文件,使用 %d 个工作线程\n",
            len(*c.files), *c.workers)
    }

    // 处理文件...
    for _, file := range *c.files {
        fmt.Printf("处理: %s\n", file)
    }

    return nil
}

func main() {
    cmd := &ProcessCommand{}
    runner := cli.NewSimpleRunner(cmd)
    parser := runner.GetParser()

    cmd.files = parser.StringSlice("file", "f", "", "要处理的文件")
    cmd.workers = parser.Int("workers", "w", 4, "工作线程数")
    cmd.verbose = parser.Bool("verbose", "v", false, "详细输出")

    exitCode := runner.Run(context.Background(), os.Args)
    os.Exit(exitCode)
}
```

### 场景 3: 支持取消的长期任务

```go
func (c *MyCommand) Run(ctx context.Context, args []string) error {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for i := 0; i < 30; i++ {
        select {
        case <-ctx.Done():
            fmt.Println("任务被取消")
            return ctx.Err()
        case <-ticker.C:
            fmt.Printf("进度: %d/30\n", i+1)
        }
    }

    return nil
}
```

## 最佳实践

### 1. 使用常量避免魔法值

```go
// ✅ 好的做法
os.Exit(cli.ExitSuccess)
fmt.Println(cli.IndicatorSuccess, "完成")

// ❌ 不好的做法
os.Exit(0)
fmt.Println("✅", "完成")
```

### 2. 验证用户输入

```go
// ✅ 验证端口范围
parser.SetValidator("port", func(v interface{}) error {
    if port := v.(int); port < 1 || port > 65535 {
        return fmt.Errorf("端口必须在 1-65535 之间")
    }
    return nil
})

// ✅ 验证文件存在
parser.SetValidator("input", func(v interface{}) error {
    path := v.(string)
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return fmt.Errorf("文件不存在: %s", path)
    }
    return nil
})
```

### 3. 提供清晰的错误消息

```go
// ✅ 好的错误消息
if err != nil {
    return fmt.Errorf("读取配置文件 %s 失败: %w", path, err)
}

// ❌ 模糊的错误消息
if err != nil {
    return err
}
```

### 4. 使用上下文支持取消

```go
// ✅ 检查上下文取消
func (c *MyCommand) Run(ctx context.Context, args []string) error {
    for _, item := range items {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            process(item)
        }
    }
    return nil
}
```

### 5. 编写可测试的命令

```go
// ✅ 可测试的设计
type MyCommand struct {
    stdout io.Writer
    input  *string
}

func (c *MyCommand) Run(ctx context.Context, args []string) error {
    fmt.Fprintf(c.stdout, "处理: %s\n", *c.input)
    return nil
}

// 测试
func TestMyCommand(t *testing.T) {
    var buf bytes.Buffer
    cmd := &MyCommand{stdout: &buf}
    input := "test.txt"
    cmd.input = &input

    err := cmd.Run(context.Background(), nil)
    assert.NoError(t, err)
    assert.Contains(t, buf.String(), "处理: test.txt")
}
```

## 项目结构

```
pkg/cli/
├── constants.go      # 常量定义(退出码、消息、默认值)
├── types.go         # 类型定义(Command、Context等)
├── parser.go        # 参数解析器
├── runner.go        # 命令运行器
├── config.go        # 配置管理
├── doc.go           # 包文档
├── README.md        # 本文档
└── example_test.go  # 示例代码
```

## 参考链接

- [Go flag 包](https://pkg.go.dev/flag)
- [CLI 最佳实践](https://clig.dev/)
- [12 Factor CLI Apps](https://medium.com/@jdxcode/12-factor-cli-apps-dd3c227a0e46)
