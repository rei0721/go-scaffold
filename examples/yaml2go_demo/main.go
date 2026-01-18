package main

import (
	"fmt"
	"log"

	togo "github.com/rei0721/go-scaffold/pkg/yaml2go"
)

func main() {
	// 示例 1: 基本用法
	fmt.Println("=== 示例 1: 基本用法 ===")
	yamlStr1 := `
name: MyApp
version: 1.0.0
port: 8080
debug: true
`
	converter := togo.New(nil)
	code, err := converter.Convert(yamlStr1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(code)
	fmt.Println()

	// 示例 2: 嵌套结构
	fmt.Println("=== 示例 2: 嵌套结构 ===")
	yamlStr2 := `
database:
  host: localhost
  port: 5432
  credentials:
    username: admin
    password: secret
server:
  host: 0.0.0.0
  port: 8080
  timeout: 30
`
	converter2 := togo.New(&togo.Config{
		PackageName: "config",
		StructName:  "AppConfig",
	})
	code2, err := converter2.Convert(yamlStr2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(code2)
	fmt.Println()

	// 示例 3: 数组和自定义标签
	fmt.Println("=== 示例 3: 数组和自定义标签 ===")
	yamlStr3 := `
users:
  - name: alice
    age: 30
    roles:
      - admin
      - user
  - name: bob
    age: 25
    roles:
      - user
`
	converter3 := togo.New(&togo.Config{
		PackageName: "model",
		StructName:  "UserConfig",
		Tags:        []string{"json", "yaml"},
		OmitEmpty:   true,
	})
	code3, err := converter3.Convert(yamlStr3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(code3)
}
