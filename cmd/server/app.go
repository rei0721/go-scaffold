package main

import (
	"fmt"

	"github.com/rei0721/go-scaffold/pkg/cli"
	"github.com/rei0721/go-scaffold/types/constants"
)

type AppCommand struct{}

func (c *AppCommand) Name() string {
	return constants.AppServerCommandName
}

func (c *AppCommand) Description() string {
	return "Run server"
}

func (c *AppCommand) Usage() string {
	return fmt.Sprintf("%s [--config=<name>]", constants.AppServerCommandName)
}

func (c *AppCommand) Flags() []cli.Flag {
	return []cli.Flag{
		{
			Name:        "config",
			ShortName:   "c",
			Type:        cli.FlagTypeString,
			Required:    false,
			Default:     constants.AppDefaultConfigPath,
			Description: "Config file path",
			EnvVar:      "REI_CONFIG_PATH", // 支持环境变量
		},
	}
}

func (c *AppCommand) Execute(ctx *cli.Context) error {
	configPath := ctx.GetString("config")

	// 业务逻辑...
	runApp(configPath)

	return nil
}
