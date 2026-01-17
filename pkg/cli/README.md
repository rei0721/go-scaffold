# CLI Package

提供企业级通用 CLI 工具框架,支持标准化命令结构、类型安全的参数绑定和可测试性设计,适用于代码生成器、数据迁移工具、运维脚本等场景。

## 特性

- ✅ **标准化命令结构**: 遵循 POSIX 规范,支持子命令和选项解析
- ✅ **类型安全**: 自动类型转换和验证 (string, int, bool, []string)
- ✅ **可测试性**: 支持 Mock I/O 的测试友好接口
- ✅ **接口化设计**: 便于依赖注入和单元测试
- ✅ **环境变量回退**: Flag 支持从环境变量读取默认值
- ✅ **友好错误提示**: 标准化错误码和帮助信息

## 快速开始

### 安装

```bash
import "github.com/rei0721/go-scaffold/pkg/cli"
```

### 基本使用

```go
package main

import (
    "fmt"
    "os"
    "github.com/rei0721/go-scaffold/pkg/cli"
)

func main() {
    // 1. 创建 CLI 应用
    app := cli.NewApp("mytool")
    app.SetVersion("1.0.0")
    app.SetDescription("My awesome CLI tool")

    // 2. 注册命令
    app.AddCommand(&GenerateCommand{})

    // 3. 执行
    if err := app.Run(os.Args[1:]); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(cli.GetExitCode(err))
    }
}
```

### 定义命令

```go
type GenerateCommand struct{}

func (c *GenerateCommand) Name() string {
    return "generate"
}

func (c *GenerateCommand) Description() string {
    return "Generate code from templates"
}

func (c *GenerateCommand) Usage() string {
    return "generate --model=<name> [--output=<dir>]"
}

func (c *GenerateCommand) Flags() []cli.Flag {
    return []cli.Flag{
        {
            Name:        "model",
            ShortName:   "m",
            Type:        cli.FlagTypeString,
            Required:    true,
            Description: "Model name to generate",
        },
        {
            Name:        "output",
            ShortName:   "o",
            Type:        cli.FlagTypeString,
            Default:     "./models",
            Description: "Output directory",
            EnvVar:      "OUTPUT_DIR", // 支持环境变量
        },
        {
            Name:        "force",
            ShortName:   "f",
            Type:        cli.FlagTypeBool,
            Default:     false,
            Description: "Overwrite existing files",
        },
    }
}

func (c *GenerateCommand) Execute(ctx *cli.Context) error {
    model := ctx.GetString("model")
    output := ctx.GetString("output")
    force := ctx.GetBool("force")

    fmt.Fprintf(ctx.Stdout, "Generating %s to %s (force=%v)\n",
        model, output, force)

    // 业务逻辑...

    return nil
}
```

## API 文档

### App 接口

| 方法                | 说明                         |
| ------------------- | ---------------------------- |
| `NewApp(name)`      | 创建新的 CLI 应用            |
| `SetVersion(v)`     | 设置版本号                   |
| `SetDescription(d)` | 设置描述                     |
| `AddCommand(cmd)`   | 注册子命令                   |
| `Run(args)`         | 执行 CLI                     |
| `RunWithIO(...)`    | 使用自定义 I/O 执行 (测试用) |

### Command 接口

| 方法            | 说明                       |
| --------------- | -------------------------- |
| `Name()`        | 返回命令名称               |
| `Description()` | 返回命令描述 (help 输出用) |
| `Usage()`       | 返回使用说明 (可选)        |
| `Flags()`       | 返回命令支持的选项列表     |
| `Execute(ctx)`  | 执行命令逻辑               |

### Context 方法

| 方法                   | 说明               |
| ---------------------- | ------------------ |
| `GetString(name)`      | 获取字符串类型选项 |
| `GetInt(name)`         | 获取整数类型选项   |
| `GetBool(name)`        | 获取布尔类型选项   |
| `GetStringSlice(name)` | 获取字符串数组选项 |
| `Args`                 | 位置参数列表       |
| `Stdin/Stdout/Stderr`  | I/O 流             |

## Flag 类型

### 支持的类型

```go
const (
    FlagTypeString       // 字符串
    FlagTypeInt          // 整数
    FlagTypeBool         // 布尔值
    FlagTypeStringSlice  // 字符串数组 (逗号分隔)
)
```

### Flag 定义

```go
type Flag struct {
    Name        string        // 长选项名 (--output)
    ShortName   string        // 短选项名 (-o)
    Type        FlagType      // 选项类型
    Required    bool          // 是否必填
    Default     interface{}   // 默认值
    Description string        // 描述信息
    EnvVar      string        // 环境变量名 (用于回退)
}
```

### 示例

```go
{
    Name:        "port",
    ShortName:   "p",
    Type:        cli.FlagTypeInt,
    Default:     8080,
    Description: "Server port",
    EnvVar:      "APP_PORT",  // 从 $APP_PORT 读取默认值
}
```

## 使用示例

### 1. 带依赖注入的命令

```go
type MigrateCommand struct {
    db      database.Database
    logger  logger.Logger
}

func NewMigrateCommand(db database.Database, log logger.Logger) *MigrateCommand {
    return &MigrateCommand{
        db:     db,
        logger: log,
    }
}

func (c *MigrateCommand) Execute(ctx *cli.Context) error {
    c.logger.Info("Running database migration")

    // 使用注入的依赖
    if err := c.db.Migrate(); err != nil {
        return fmt.Errorf("migration failed: %w", err)
    }

    fmt.Fprintln(ctx.Stdout, "✅ Migration completed")
    return nil
}

// 使用
func main() {
    db := database.New(cfg)
    logger := logger.Default()

    app := cli.NewApp("dbmigrate")
    app.AddCommand(NewMigrateCommand(db, logger))

   os.Exit(app.Run(os.Args[1:]))
}
```

### 2. 错误处理

```go
func (c *MyCommand) Execute(ctx *cli.Context) error {
    // 参数错误 (退出码 2)
    if invalidInput {
        return &cli.UsageError{
            Command: "generate",
            Message: "invalid model name",
        }
    }

    // 执行错误 (退出码 1)
    if err := doWork(); err != nil {
        return &cli.CommandError{
            Command: "generate",
            Message: "failed to generate",
            Cause:   err,
        }
    }

    return nil
}
```

### 3. 测试命令

```go
func TestGenerateCommand(t *testing.T) {
    app := cli.NewApp("test")
    app.AddCommand(&GenerateCommand{})

    // Mock I/O
    var stdout bytes.Buffer
    var stderr bytes.Buffer

    err := app.RunWithIO(
        []string{"generate", "--model", "User", "--output", "./out"},
        nil,       // stdin
        &stdout,   // stdout
        &stderr,   // stderr
    )

    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    output := stdout.String()
    if !strings.Contains(output, "User") {
        t.Errorf("expected output to contain 'User', got: %s", output)
    }
}
```

## 命令行用法

### 查看帮助

```bash
# 应用帮助
$ mytool --help

# 命令帮助
$ mytool generate --help

# 版本信息
$ mytool --version
```

### 使用选项

```bash
# 长选项
$ mytool generate --model=User --output=./models

# 短选项
$ mytool generate -m User -o ./models

# 布尔选项
$ mytool generate --model=User --force

# 数组选项 (逗号分隔)
$ mytool list --tags=go,cli,tool
```

### 环境变量回退

```bash
# 设置环境变量
$ export OUTPUT_DIR=/tmp/models
$ export APP_PORT=9090

# 使用环境变量默认值
$ mytool generate --model=User
# output 将使用 $OUTPUT_DIR 的值
```

## 错误码

遵循 Unix 退出码约定：

| 退出码 | 常量            | 含义              | 使用场景         |
| ------ | --------------- | ----------------- | ---------------- |
| `0`    | ExitSuccess     | 成功              | 命令正常完成     |
| `1`    | ExitError       | 通用错误          | 运行时错误       |
| `2`    | ExitUsage       | 参数错误          | 命令行参数非法   |
| `3`    | ExitConfig      | 配置错误          | 配置文件无效     |
| `130`  | ExitInterrupted | 用户中断 (Ctrl+C) | 捕获 SIGINT 信号 |

### 获取退出码

```go
if err := app.Run(os.Args[1:]); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(cli.GetExitCode(err))  // 自动提取错误码
}
```

## 最佳实践

### 1. 使用依赖注入

```go
// ✅ 推荐：通过构造函数注入
type MyCommand struct {
    db     database.Database
    logger logger.Logger
}

func NewMyCommand(db database.Database, log logger.Logger) *MyCommand {
    return &MyCommand{db: db, logger: log}
}

// ❌ 避免：使用全局变量
var globalDB database.Database

type MyCommand struct{}
```

### 2. 验证参数

```go
func (c *MyCommand) Execute(ctx *cli.Context) error {
    port := ctx.GetInt("port")

    // 验证范围
    if port < 1 || port > 65535 {
        return &cli.UsageError{
            Message: "port must be between 1 and 65535",
        }
    }

    // ...
}
```

### 3. 友好的输出

```go
func (c *MyCommand) Execute(ctx *cli.Context) error {
    // 进度提示
    fmt.Fprintf(ctx.Stdout, "Processing...\n")

    // 成功消息
    fmt.Fprintf(ctx.Stdout, "✅ Operation completed\n")

    // 警告 (写入 stderr)
    fmt.Fprintf(ctx.Stderr, "⚠️  Warning: deprecated flag\n")

    return nil
}
```

### 4. 结构化错误

```go
// ✅ 使用结构化错误
return &cli.CommandError{
    Command: "migrate",
    Message: "database connection failed",
    Cause:   err,
}

// ❌ 避免裸字符串
return errors.New("error")
```

## 项目结构

```
pkg/cli/
├── cli.go          # 核心接口 (App, Command, Context, Flag)
├── app.go          # App 实现
├── flag.go         # Flag 解析器
├── constants.go    # 错误码和常量
├── errors.go       # 错误类型
├── doc.go          # Go doc 文档
└── README.md       # 本文档
```

## 依赖项

- Go 标准库 `flag` - 参数解析
- 无第三方依赖

## 相关资源

- [Go flag 包文档](https://pkg.go.dev/flag)
- [POSIX 参数约定](https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap12.html)

## 许可证

本项目使用 MIT 许可证。
