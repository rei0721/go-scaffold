### 通用 CLI 工具库 (General-Purpose CLI Framework)

#### 1. 组件定位

为当前系统 `rei0721` 定制，在 `pkg` 基础库层封装**企业级通用 CLI 工具框架**。该组件旨在作为项目级的命令行基础设施，提供统一的、可扩展的 CLI 构建能力，适用于代码生成器、数据迁移工具、运维脚本及开发辅助工具等多种场景。

#### 2. 核心能力 (Core Capabilities)

- **标准化命令结构 (Standardized Command Structure)**：
  提供统一的命令注册、参数解析、子命令嵌套能力，遵循 POSIX 规范，支持短选项（`-h`）和长选项（`--help`），确保用户体验一致性。

- **类型安全的参数绑定 (Type-Safe Argument Binding)**：
  支持将命令行参数自动绑定到 Go 结构体，内置类型验证（`int`、`bool`、`string`、`duration` 等），减少手动解析和类型转换的样板代码。

- **交互式增强 (Interactive Enhancements)**：
  内置交互式输入能力（如确认提示、选择列表、进度条），支持彩色输出和格式化表格，提升 CLI 工具的用户友好度。

- **统一错误处理 (Unified Error Handling)**：
  提供标准化的错误码体系（0=成功，1=通用错误，2=参数错误等），支持错误堆栈追踪和友好的错误提示，便于问题诊断和脚本集成。

- **可测试性设计 (Testability-First Design)**：
  支持命令的单元测试（通过 Mock I/O），提供 `Execute(args []string)` 等测试友好接口，确保 CLI 逻辑可被自动化测试覆盖。

- **插件化架构 (Plugin Architecture)**：
  支持通过接口动态注册命令，便于扩展新功能而无需修改核心代码，遵循开闭原则（OCP）。

#### 3. 设计约束 (Design Constraints)

> 为了保证组件的通用性、一致性与可维护性，需遵守以下约束：

**3.1 架构规范**

- **包路径**：`pkg/cli`。
- **零全局依赖**：严禁使用全局变量存储命令或配置，所有命令必须通过 `NewApp()` 或 `RegisterCommand()` 显式注册。
- **接口定义**：
  对外暴露的核心接口必须包含以下能力：

```go
// App CLI 应用接口
type App interface {
    // Name 返回应用名称
    Name() string
    // Version 返回应用版本
    Version() string
    // AddCommand 注册子命令
    AddCommand(cmd Command) error
    // Run 执行 CLI,解析参数并路由到对应命令
    Run(args []string) error
}

// Command 命令接口
type Command interface {
    // Name 命令名称(如 "generate", "migrate")
    Name() string
    // Description 命令描述(用于 help 输出)
    Description() string
    // Flags 返回命令支持的选项列表
    Flags() []Flag
    // Execute 命令执行逻辑
    Execute(ctx *Context) error
}

// Context 命令执行上下文
type Context struct {
    // Args 位置参数
    Args []string
    // Flags 解析后的选项值
    Flags map[string]interface{}
    // Stdin/Stdout/Stderr IO 流(可 Mock)
    Stdin  io.Reader
    Stdout io.Writer
    Stderr io.Writer
}
```

**3.2 参数解析规范**

- **多级子命令**：支持嵌套子命令（如 `app db migrate up`），最多支持 3 级嵌套。
- **参数类型**：

  - **位置参数**（Positional Args）：按顺序解析，支持必选/可选标记。
  - **选项参数**（Flags）：支持短选项（`-v`）、长选项（`--verbose`）、带值选项（`--output=file.txt`）和布尔开关（`--force`）。
  - **环境变量回退**：支持从环境变量读取默认值（如 `--port` 默认读取 `$APP_PORT`）。

- **验证机制**：
  - 必填参数检查（Required）
  - 值约束验证（如端口范围 `1-65535`）
  - 互斥选项检查（如 `--json` 与 `--yaml` 不可同时使用）

**3.3 输出规范**

- **标准化输出**：

  - **正常输出**：写入 `Stdout`，可被管道捕获。
  - **错误输出**：写入 `Stderr`，包含错误码和堆栈（`--debug` 模式）。
  - **日志输出**：集成 `pkg/logger`，支持通过 `--log-level` 动态调整。

- **格式化能力**：
  - **表格输出**：支持自动对齐的表格（用于 `list` 类命令）。
  - **JSON/YAML 输出**：支持结构化输出（`--output=json`）。
  - **彩色输出**：支持彩色高亮（成功=绿色，错误=红色，警告=黄色），可通过 `--no-color` 禁用。

**3.4 交互式能力**

- **确认提示**：

  ```go
  if !ctx.ConfirmPrompt("Are you sure to delete?") {
      return ErrCancelled
  }
  ```

- **选择列表**：

  ```go
  choice := ctx.SelectPrompt("Choose environment:", []string{"dev", "test", "prod"})
  ```

- **进度条**：
  ```go
  progress := ctx.NewProgressBar(total)
  for i := 0; i < total; i++ {
      // 处理任务...
      progress.Increment()
  }
  ```

**3.5 错误码标准**

遵循 Unix 退出码约定：

| 退出码 | 含义              | 使用场景                 |
| ------ | ----------------- | ------------------------ |
| `0`    | 成功              | 命令正常完成             |
| `1`    | 通用错误          | 运行时错误（如网络失败） |
| `2`    | 参数错误          | 命令行参数非法           |
| `3`    | 配置错误          | 配置文件无效             |
| `130`  | 用户中断 (Ctrl+C) | 捕获 SIGINT 信号         |

**3.6 测试规范**

- **命令可测试性**：

  ```go
  func TestGenerateCommand(t *testing.T) {
      app := cli.NewApp("test")
      app.AddCommand(&GenerateCommand{})

      // Mock I/O
      var stdout bytes.Buffer
      err := app.RunWithIO(
          []string{"generate", "--model", "User"},
          nil,           // stdin
          &stdout,       // stdout
          io.Discard,    // stderr
      )

      assert.NoError(t, err)
      assert.Contains(t, stdout.String(), "Generated successfully")
  }
  ```

- **单元测试覆盖率要求**：核心逻辑覆盖率 ≥ 80%。

**3.7 配置文件支持**

- **可选配置文件**：支持从 `.cli-config.yaml` 读取默认选项（优先级：命令行 > 环境变量 > 配置文件 > 默认值）。
- **配置示例**：
  ```yaml
  # .cli-config.yaml
  defaults:
    output: json
    log-level: info

  generate:
    template-dir: ./templates
    output-dir: ./generated
  ```

**3.8 依赖注入规范**

- **接口隔离**：业务逻辑通过 `Command` 接口注册，不依赖具体 CLI 框架实现。
- **依赖注入示例**：
  ```go
  // cmd/sqlgen/main.go
  func main() {
      app := cli.NewApp("sqlgen")
      app.Version("1.0.0")

      // 注入依赖
      db := database.New(dbConfig)
      logger := logger.New(logConfig)

      // 注册命令
      app.AddCommand(commands.NewGenerateCommand(db, logger))
      app.AddCommand(commands.NewMigrateCommand(db, logger))

      os.Exit(app.Run(os.Args[1:]))
  }
  ```

#### 4. 技术选型建议

**4.1 推荐基础库**

可基于以下成熟库进行封装：

| 库                       | 用途       | 理由                        |
| ------------------------ | ---------- | --------------------------- |
| `spf13/cobra`            | 命令框架   | 业界标准，kubectl/hugo 在用 |
| `spf13/pflag`            | 参数解析   | 兼容 POSIX，类型安全        |
| `AlecAivazis/survey`     | 交互式提示 | 丰富的交互组件              |
| `fatih/color`            | 彩色输出   | 简洁易用                    |
| `olekukonko/tablewriter` | 表格输出   | 自动对齐，支持多种样式      |

**4.2 最小依赖原则**

如果不需要复杂子命令，可仅使用 `flag` 标准库 + 自定义封装，减少外部依赖。

#### 5. 使用示例

**5.1 定义命令**

```go
// internal/commands/generate.go
type GenerateCommand struct {
    db     database.Database
    logger logger.Logger
}

func (c *GenerateCommand) Name() string {
    return "generate"
}

func (c *GenerateCommand) Description() string {
    return "Generate code from database schema"
}

func (c *GenerateCommand) Flags() []cli.Flag {
    return []cli.Flag{
        {Name: "model", Type: "string", Required: true, Description: "Model name"},
        {Name: "output", Type: "string", Default: "./models", Description: "Output directory"},
        {Name: "force", Type: "bool", Default: false, Description: "Overwrite existing files"},
    }
}

func (c *GenerateCommand) Execute(ctx *cli.Context) error {
    model := ctx.GetString("model")
    output := ctx.GetString("output")
    force := ctx.GetBool("force")

    if !force && fileExists(output) {
        if !ctx.ConfirmPrompt("File exists. Overwrite?") {
            return cli.ErrCancelled
        }
    }

    c.logger.Info("Generating model", "name", model, "output", output)

    // 业务逻辑...

    fmt.Fprintf(ctx.Stdout, "✅ Generated %s successfully\n", model)
    return nil
}
```

**5.2 注册与运行**

```go
// cmd/sqlgen/main.go
func main() {
    app := cli.NewApp("sqlgen")
    app.Version("1.0.0")
    app.Description("SQL code generator for rei0721")

    // 初始化依赖
    logger := logger.Default()
    db, _ := database.New(dbConfig)

    // 注册命令
    app.AddCommand(&commands.GenerateCommand{db: db, logger: logger})
    app.AddCommand(&commands.MigrateCommand{db: db, logger: logger})

    // 运行
    if err := app.Run(os.Args[1:]); err != nil {
        logger.Error("command failed", "error", err)
        os.Exit(1)
    }
}
```

**5.3 使用示例**

```bash
# 生成代码
$ sqlgen generate --model=User --output=./models

# 带确认提示
$ sqlgen generate --model=User --force

# 子命令
$ sqlgen migrate up --step=1

# 查看帮助
$ sqlgen --help
$ sqlgen generate --help

# JSON 输出
$ sqlgen list --output=json

# 调试模式
$ sqlgen generate --model=User --log-level=debug
```

#### 6. 质量要求

- **文档完整性**：每个命令必须提供详细的 `--help` 输出和 `README.md` 使用示例。
- **错误提示友好**：错误信息应包含建议的修复方法（如参数错误时提示 `Run 'command --help' for usage`）。
- **向后兼容**：命令行接口一旦发布，不得删除或修改已有选项，仅可新增或标记 `deprecated`。
- **性能要求**：CLI 启动时间应 < 100ms（不含业务逻辑），参数解析开销可忽略不计。

#### 7. 可选增强特性

- **自动补全**：生成 Bash/Zsh 自动补全脚本。
- **版本检查**：支持检查新版本（`--check-update`）。
- **遥测数据**：可选的匿名使用统计（需用户明确 opt-in）。
- **插件系统**：支持从外部 `.so` 文件加载插件命令。

---

### 总结

本组件旨在为 `rei0721` 项目提供**统一、规范、易测试**的 CLI 工具构建能力，确保所有命令行工具（代码生成器、数据迁移、运维脚本等）遵循一致的用户体验和代码规范，降低维护成本并提升开发效率。
