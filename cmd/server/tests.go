package main

import (
	"fmt"
	"log"

	"github.com/rei0721/go-scaffold/pkg/cli"
	"github.com/rei0721/go-scaffold/pkg/yaml2go"
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
	yamlStr := `
server:
  host: localhost
  port: 8080
  required: true

database:
  driver: mysql
  host: localhost
  port: 3306
`

	// 创建转换器
	converter := yaml2go.New(nil)

	// 转换 YAML
	result, err := converter.Convert(yamlStr)
	if err != nil {
		log.Fatal(err)
	}

	// 查看主配置
	fmt.Println("=== config.go ===")
	fmt.Println(result.MainConfig.Content)

	// 查看子配置
	for _, subConfig := range result.SubConfigs {
		fmt.Printf("\n=== %s ===\n", subConfig.FileName)
		fmt.Println(subConfig.Content)
	}
	return nil
}
