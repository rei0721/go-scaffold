package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// app CLI 应用实现
type app struct {
	name        string
	version     string
	description string
	commands    map[string]Command
	mu          sync.RWMutex
}

// NewApp 创建新的 CLI 应用
func NewApp(name string) App {
	return &app{
		name:     name,
		commands: make(map[string]Command),
	}
}

// Name 返回应用名称
func (a *app) Name() string {
	return a.name
}

// Version 返回应用版本
func (a *app) Version() string {
	return a.version
}

// Description 返回应用描述
func (a *app) Description() string {
	return a.description
}

// SetVersion 设置版本号
func (a *app) SetVersion(version string) {
	a.version = version
}

// SetDescription 设置描述
func (a *app) SetDescription(desc string) {
	a.description = desc
}

// AddCommand 注册子命令
func (a *app) AddCommand(cmd Command) error {
	if cmd == nil {
		return fmt.Errorf("command cannot be nil")
	}

	cmdName := cmd.Name()
	if cmdName == "" {
		return fmt.Errorf("command name cannot be empty")
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if _, exists := a.commands[cmdName]; exists {
		return fmt.Errorf("%s: %s", ErrMsgDuplicateCommand, cmdName)
	}

	a.commands[cmdName] = cmd
	return nil
}

// Run 执行 CLI
func (a *app) Run(args []string) error {
	return a.RunWithIO(args, os.Stdin, os.Stdout, os.Stderr)
}

// RunWithIO 执行 CLI，使用自定义 I/O
func (a *app) RunWithIO(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	// 没有参数或只有 help 选项，显示帮助
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		a.printHelp(stdout)
		return nil
	}

	// 版本选项
	if args[0] == "--version" || args[0] == "-v" {
		a.printVersion(stdout)
		return nil
	}

	// 查找命令
	cmdName := args[0]
	a.mu.RLock()
	cmd, exists := a.commands[cmdName]
	a.mu.RUnlock()

	if !exists {
		return &UsageError{
			Message: fmt.Sprintf("%s: %s", ErrMsgCommandNotFound, cmdName),
		}
	}

	// 解析命令选项
	parser := newFlagParser(cmdName, cmd.Flags())
	remainingArgs, err := parser.parse(args[1:])
	if err != nil {
		return err
	}

	// 创建执行上下文
	ctx := &Context{
		Args:   remainingArgs,
		Flags:  parser.getValues(),
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	}

	// 执行命令
	if err := cmd.Execute(ctx); err != nil {
		return &CommandError{
			Command: cmdName,
			Message: "execution failed",
			Cause:   err,
		}
	}

	return nil
}

// printHelp 打印帮助信息
func (a *app) printHelp(w io.Writer) {
	fmt.Fprintf(w, "%s", a.name)
	if a.version != "" {
		fmt.Fprintf(w, " v%s", a.version)
	}
	fmt.Fprintln(w)

	if a.description != "" {
		fmt.Fprintf(w, "\n%s\n", a.description)
	}

	fmt.Fprintln(w, "\nUsage:")
	fmt.Fprintf(w, "  %s [command] [flags]\n", a.name)

	fmt.Fprintln(w, "\nAvailable Commands:")
	a.mu.RLock()
	defer a.mu.RUnlock()

	if len(a.commands) == 0 {
		fmt.Fprintln(w, "  (no commands registered)")
		return
	}

	// 找到最长的命令名，用于对齐
	maxLen := 0
	for name := range a.commands {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	for name, cmd := range a.commands {
		padding := strings.Repeat(" ", maxLen-len(name)+2)
		fmt.Fprintf(w, "  %s%s%s\n", name, padding, cmd.Description())
	}

	fmt.Fprintln(w, "\nFlags:")
	fmt.Fprintln(w, "  -h, --help       Show help information")
	fmt.Fprintln(w, "  -v, --version    Show version information")

	fmt.Fprintf(w, "\nRun '%s [command] --help' for more information on a command.\n", a.name)
}

// printVersion 打印版本信息
func (a *app) printVersion(w io.Writer) {
	if a.version != "" {
		fmt.Fprintf(w, "%s version %s\n", a.name, a.version)
	} else {
		fmt.Fprintf(w, "%s (version not set)\n", a.name)
	}
}
