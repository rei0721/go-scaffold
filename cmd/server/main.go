// Package main 是应用程序的入口点
// 它负责初始化 App 容器并处理操作系统信号的优雅关闭
package main

import (
	"fmt"
	"os"

	"github.com/rei0721/go-scaffold/pkg/cli"
)

const (
	AppCommandName    = "server"
	InitdbCommandName = "initdb"
	VersionNumber     = "0.1.2"
)

func main() {
	// 创建 CLI 应用
	app := cli.NewApp("go-scaffold")
	app.SetVersion(VersionNumber)
	app.SetDescription("This is a go backend scaffold")

	// 注册命令
	app.AddCommand(&AppCommand{})
	app.AddCommand(&InitdbCommand{})

	// 执行
	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cli.GetExitCode(err))
	}
}
