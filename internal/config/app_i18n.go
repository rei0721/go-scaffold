package config

import (
	"errors"
	"os"
	"strings"
)

// I18nConfig 国际化配置
// 支持多语言
type I18nConfig struct {
	// Default 默认语言
	// 当请求的语言不支持时使用
	// 例如: en, zh-CN, ja
	Default string `mapstructure:"default"`

	// Supported 支持的语言列表
	// 必须包含 Default 语言
	// 例如: ["en", "zh-CN", "ja"]
	Supported []string `mapstructure:"supported"`

	// MessagesDir 语言文件目录
	// 包含所有语言的翻译文件
	// 目录结构: MessagesDir/{lang}.yaml
	// 例如: ./configs/locales/en.yaml, ./configs/locales/zh-CN.yaml
	MessagesDir string `mapstructure:"messages_dir"`
}

func (c *I18nConfig) ValidateName() string {
	return AppI18nName
}

func (c *I18nConfig) ValidateRequired() bool {
	return true
}

// Validate 验证国际化配置
// 实现 Configurable 接口
func (c *I18nConfig) Validate() error {
	// 验证默认语言
	if c.Default == "" {
		return errors.New("default locale is required")
	}

	// 验证支持的语言列表
	if len(c.Supported) == 0 {
		return errors.New("at least one supported locale is required")
	}

	// 确保默认语言在支持列表中
	found := false
	for _, s := range c.Supported {
		if s == c.Default {
			found = true
			break
		}
	}
	if !found {
		// 默认语言必须是支持的语言之一
		// 否则会导致运行时错误
		return errors.New("default locale must be in supported list")
	}

	return nil
}

// overrideI18nConfig 使用环境变量覆盖国际化配置
func overrideI18nConfig(cfg *I18nConfig) {
	// Default
	if val := os.Getenv(EnvI18nDefault); val != "" {
		cfg.Default = val
	}

	// Supported
	// 环境变量格式: "zh-CN,en-US,ja-JP"
	// 解析为: ["zh-CN", "en-US", "ja-JP"]
	if val := os.Getenv(EnvI18nSupported); val != "" {
		langs := strings.Split(val, DefaultSeparator)
		// 去除空白
		var supported []string
		for _, lang := range langs {
			trimmed := strings.TrimSpace(lang)
			if trimmed != "" {
				supported = append(supported, trimmed)
			}
		}
		if len(supported) > 0 {
			cfg.Supported = supported
		}
	}
}
