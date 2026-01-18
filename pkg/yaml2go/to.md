### 基本用法

```go
package main

import (
    "fmt"
    "log"
    "github.com/rei0721/go-scaffold/pkg/yaml2go"
)

func main() {
    yamlStr := `
server:
  required: ${REI_SERVER_REQUIRED:true}
  host: ${REI_SERVER_HOST:localhost}
  port: ${REI_SERVER_PORT:8080}
`

    // 创建转换器（使用默认配置）
    // 为所有配置默认创建统一的方法
    converter := yaml2go.New(nil)

    // 转换 YAML 为 Go 代码
    code, err := converter.Convert(yamlStr)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(code)
}
```

```go
// config.go
type Config struct {
	Server ServerConfig `mapstructure:"server"`
}
```

```go
// server_config.go
type ServerConfig struct {
    Required bool `json:"required" yaml:"required" mapstructure:"required" toml:"required"`
    Host string `json:"host" yaml:"host" mapstructure:"host" toml:"host"`
    Port int64 `json:"port" yaml:"port" mapstructure:"port" toml:"port"`
}

// ValidateName 返回配置名称
func (c *ServerConfig) ValidateName() string {
	return "server"
}

// Validate 为开发者生成一个默认的验证器接口
func (c *ServerConfig) Validate() error {
	return nil
}

// OverrideConfig 使用环境变量覆盖配置
func (cfg *ServerConfig) OverrideConfig() {
	// Host
	if val := os.Getenv("REI_SERVER_HOST"); val != "" {
		if host, err := strconv.Atoi(val); err == nil {
			cfg.Host = host
		}
	}

	// Port
	if val := os.Getenv("REI_SERVER_PORT"); val != "" {
        if port, err := strconv.Atoi(val); err == nil {
            cfg.Port = port
        }
	}

	// Required
	if val := os.Getenv("REI_SERVER_REQUIRED"); val != "" {
		if required, err := strconv.ParseBool(val); err == nil {
			cfg.Required = required
		}
	}
}
```
