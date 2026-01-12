package cli

import "fmt"

// Config CLI 配置结构
// 用于管理 CLI 应用的全局配置
// 设计目标:
// - 提供合理的默认值
// - 支持配置验证
// - 易于序列化(可从文件加载)
type Config struct {
	// AppName 应用程序名称
	// 用于帮助文档和错误消息
	AppName string

	// Version 应用程序版本
	// 用于 --version 标志
	Version string

	// Description 应用程序描述
	// 用于帮助文档的头部
	Description string

	// EnableHelp 是否自动处理 --help 标志
	// 默认为 true
	EnableHelp bool

	// EnableVersion 是否自动处理 --version 标志
	// 默认为 true
	EnableVersion bool

	// ExitOnError 解析错误时是否自动退出
	// 如果为 false,错误会返回给调用者
	// 默认为 true
	ExitOnError bool

	// UsageTemplate 自定义用法模板
	// 为空时使用默认模板
	UsageTemplate string

	// HideBuiltinFlags 是否隐藏内置标志(help, version)
	// 在帮助文档中不显示
	// 默认为 false
	HideBuiltinFlags bool
}

// DefaultConfig 返回默认配置
// 提供合理的默认值,适合大多数应用
// 返回:
//
//	*Config: 默认配置
//
// 使用示例:
//
//	config := DefaultConfig()
//	config.AppName = "myapp"
//	config.Version = "1.0.0"
func DefaultConfig() *Config {
	return &Config{
		AppName:          "app",
		Version:          "0.0.0",
		Description:      "",
		EnableHelp:       true,
		EnableVersion:    true,
		ExitOnError:      true,
		UsageTemplate:    "",
		HideBuiltinFlags: false,
	}
}

// Validate 验证配置的有效性
// 检查必需字段和值的合法性
// 返回:
//
//	error: 配置无效时的错误
//
// 使用示例:
//
//	config := &Config{...}
//	if err := config.Validate(); err != nil {
//	  log.Fatal(err)
//	}
func (c *Config) Validate() error {
	// AppName 是必需的
	if c.AppName == "" {
		return fmt.Errorf(ErrMsgConfigInvalid, "app name is required")
	}

	// Version 应该遵循语义化版本(可选检查)
	// 这里只做基本检查
	if c.Version == "" {
		return fmt.Errorf(ErrMsgConfigInvalid, "version is required")
	}

	return nil
}

// Clone 创建配置的副本
// 用于避免意外修改原始配置
// 返回:
//
//	*Config: 配置的副本
func (c *Config) Clone() *Config {
	clone := *c
	return &clone
}

// Merge 合并另一个配置
// 用于从多个来源组合配置
// 非零值会覆盖当前配置
// 参数:
//
//	other: 要合并的配置
func (c *Config) Merge(other *Config) {
	if other.AppName != "" {
		c.AppName = other.AppName
	}
	if other.Version != "" {
		c.Version = other.Version
	}
	if other.Description != "" {
		c.Description = other.Description
	}
	if other.UsageTemplate != "" {
		c.UsageTemplate = other.UsageTemplate
	}
	// 布尔值总是合并
	c.EnableHelp = other.EnableHelp
	c.EnableVersion = other.EnableVersion
	c.ExitOnError = other.ExitOnError
	c.HideBuiltinFlags = other.HideBuiltinFlags
}
