package utils

import (
	"github.com/rei0721/go-scaffold/pkg/i18n"
)

// I18nUtils 国际化工具类，提供多语言消息翻译功能
// 封装了 i18n 库的核心功能，并添加默认语言支持
type I18nUtils struct {
	// _i18n: i18n 库的核心实例，提供多语言翻译能力
	// 前缀下划线表示这是私有字段，仅在内部使用
	_i18n i18n.I18n

	// defaultLanguage: 默认语言代码，例如 "zh-CN"、"en-US" 等
	// 当未指定语言时，使用此默认语言进行翻译
	defaultLanguage string
}

// NewI18nUtils 创建并返回一个新的 I18nUtils 实例
// 参数:
//   - _i18n: i18n 库的实例，提供基础翻译功能
//   - defaultLanguage: 默认使用的语言代码
//
// 返回值:
//   - *I18nUtils: 初始化后的 I18nUtils 指针
//
// 说明: 这是 I18nUtils 的构造函数，用于初始化结构体字段
func NewI18nUtils(_i18n i18n.I18n, defaultLanguage string) *I18nUtils {
	return &I18nUtils{
		_i18n:           _i18n,
		defaultLanguage: defaultLanguage,
	}
}

// T 获取指定消息ID的翻译文本（使用默认语言）
// 参数:
//   - messageID: 要翻译的消息标识符（Key）
//   - templates: 可选的模板参数，用于替换消息中的占位符
//
// 返回值:
//   - string: 翻译后的文本
//
// 说明:
//   - 此方法是 i18n.I18n.T 方法的包装器，自动使用默认语言
//   - templates 参数为可变参数，支持传递多个模板映射或为空
//   - 使用默认语言调用底层 i18n 实例的翻译方法
func (i I18nUtils) T(messageID string, templates ...map[string]interface{}) string {
	// 调用底层 i18n 实例的 T 方法，传入默认语言、消息ID和模板参数
	// 可变参数 templates... 会将参数切片展开传递给被调用函数
	return i._i18n.T(i.defaultLanguage, messageID, templates...)
}
