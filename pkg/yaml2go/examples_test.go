package yaml2go_test

import (
	"fmt"
	"log"

	togo "github.com/rei0721/go-scaffold/pkg/yaml2go"
)

// Example_basic 演示基本用法
func Example_basic() {
	yamlStr := `
name: John Doe
age: 30
email: john@example.com
`

	converter := togo.New(nil)
	code, err := converter.Convert(yamlStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(code)
}

// Example_nestedStruct 演示嵌套结构转换
func Example_nestedStruct() {
	yamlStr := `
database:
  host: localhost
  port: 5432
server:
  port: 8080
  timeout: 30
`

	converter := togo.New(&togo.Config{
		PackageName: "config",
		StructName:  "AppConfig",
	})

	code, err := converter.Convert(yamlStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(code)
}

// Example_customTags 演示自定义标签
func Example_customTags() {
	yamlStr := `
user_id: 123
user_name: alice
`

	converter := togo.New(&togo.Config{
		Tags:      []string{"json", "yaml"},
		OmitEmpty: true,
	})

	code, err := converter.Convert(yamlStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(code)
}

// Example_array 演示数组类型转换
func Example_array() {
	yamlStr := `
tags:
  - golang
  - yaml
  - converter
users:
  - name: alice
    age: 30
  - name: bob
    age: 25
`

	converter := togo.New(nil)
	code, err := converter.Convert(yamlStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(code)
}

// Example_withPointers 演示使用指针类型
func Example_withPointers() {
	yamlStr := `
host: localhost
port: 8080
`

	converter := togo.New(&togo.Config{
		UsePointer: true,
		OmitEmpty:  true,
	})

	code, err := converter.Convert(yamlStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(code)
}

// Example_saveToFile 演示保存到文件
func Example_saveToFile() {
	yamlStr := `
app_name: MyApp
version: 1.0.0
`

	converter := togo.New(&togo.Config{
		PackageName: "config",
		StructName:  "AppConfig",
	})

	err := converter.ConvertToFile(yamlStr, "config/types.go")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Config saved to config/types.go")
}

// Example_viperIntegration 演示与 Viper 集成
func Example_viperIntegration() {
	yamlStr := `
database:
  host: localhost
  port: 5432
  username: admin
  password: secret
`

	// 生成配置结构体（开发阶段）
	converter := togo.New(&togo.Config{
		PackageName: "config",
		StructName:  "DatabaseConfig",
		Tags:        []string{"mapstructure", "json", "yaml"},
	})

	code, err := converter.Convert(yamlStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("// Generated config structure for Viper:")
	fmt.Println(code)
	fmt.Println()
	fmt.Println("// Usage with Viper:")
	fmt.Println("// var cfg config.DatabaseConfig")
	fmt.Println("// viper.Unmarshal(&cfg)")
}

// Example_complexConfig 演示复杂配置转换
func Example_complexConfig() {
	yamlStr := `
app:
  name: MyApp
  version: 1.0.0
  features:
    - auth
    - logging
    - metrics
database:
  primary:
    host: localhost
    port: 5432
  replicas:
    - host: replica1
      port: 5432
    - host: replica2
      port: 5432
`

	converter := togo.New(&togo.Config{
		PackageName: "config",
		StructName:  "Config",
		AddComments: true,
	})

	code, err := converter.Convert(yamlStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(code)
}
