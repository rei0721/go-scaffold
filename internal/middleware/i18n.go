package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rei0721/rei0721/pkg/i18n"
)

// I18n 中间件提取并存储用户的语言偏好
func I18n(i18nApp i18n.I18n) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Accept-Language 头部获取语言
		lang := c.GetHeader(i18n.LanguageHeader)

		// 如果语言不支持,使用默认语言
		if lang == "" || !i18nApp.IsSupported(lang) {
			lang = i18nApp.GetDefaultLanguage()
		}

		// 存储到上下文
		c.Set("lang", lang)
		c.Set("i18n", i18nApp)

		c.Next()
	}
}
