package main

import (
	"fmt"

	"github.com/rei0721/go-scaffold/pkg/cli"
	"github.com/rei0721/go-scaffold/types/constants"
)

// TestsCommand 测试命令
type TestsCommand struct{}

func (c *TestsCommand) Name() string {
	return constants.AppTestsCommandName
}

func (c *TestsCommand) Description() string {
	return "Run tests"
}

func (c *TestsCommand) Usage() string {
	return fmt.Sprintf("%s", constants.AppTestsCommandName)
}

func (c *TestsCommand) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (c *TestsCommand) Execute(ctx *cli.Context) error {
	return nil
}
