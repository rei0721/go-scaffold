// Package main 是应用程序的入口点
// 它负责初始化 App 容器并处理操作系统信号的优雅关闭
package main

import (
	"fmt"
	"os"

	"github.com/rei0721/go-scaffold/pkg/cli"
	"github.com/rei0721/go-scaffold/types/constants"
)

func main() {
	// 创建 CLI 应用
	app := cli.NewApp(constants.AppName)
	app.SetVersion(constants.AppVersion)
	app.SetDescription(constants.AppDescription)

	// 注册命令
	app.AddCommand(&AppCommand{})
	app.AddCommand(&InitdbCommand{})
	app.AddCommand(&TestsCommand{})

	// 执行
	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cli.GetExitCode(err))
	}
}
