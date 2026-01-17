package main

import (
	"os"

	"github.com/rei0721/go-scaffold/internal/app"
	"github.com/rei0721/go-scaffold/pkg/cli"
	"github.com/rei0721/go-scaffold/types/constants"
)

// InitdbCommand 数据库初始化命令
type InitdbCommand struct{}

func (c *InitdbCommand) Name() string {
	return "initdb"
}

func (c *InitdbCommand) Description() string {
	return "Initialize database schema and data"
}

func (c *InitdbCommand) Usage() string {
	return "initdb [--config=<path>]"
}

func (c *InitdbCommand) Flags() []cli.Flag {
	return []cli.Flag{
		{
			Name:        "config",
			ShortName:   "c",
			Type:        cli.FlagTypeString,
			Required:    false,
			Default:     constants.DefaultConfigPath,
			Description: "Config file path",
			EnvVar:      "REI_CONFIG_PATH",
		},
	}
}

func (c *InitdbCommand) Execute(ctx *cli.Context) error {
	configPath := ctx.GetString("config")

	// 创建 App 实例（initdb 模式）
	application, err := app.New(app.Options{
		ConfigPath: configPath,
		Mode:       app.ModeInitDB,
	})
	if err != nil {
		os.Stderr.WriteString("failed to initialize application: " + err.Error() + "\n")
		return err
	}

	// initdb 模式下，New 函数已经执行了初始化
	// 这里只需要优雅关闭资源
	defer func() {
		if application.DB != nil {
			_ = application.DB.Close()
		}
		if application.Logger != nil {
			_ = application.Logger.Sync()
		}
	}()

	return nil
}
