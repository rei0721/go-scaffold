package app

import (
	"fmt"

	"github.com/rei0721/go-scaffold/pkg/i18n"
)

// initI18n 初始化i18n
func initI18n(app *App) error {
	// 初始化i18n
	i18nCfg := &i18n.Config{
		DefaultLanguage:    ConstantsI18nDefaultLanguage,
		SupportedLanguages: ConstantsI18nSupportedLanguages,
		MessagesDir:        ConstantsI18nMessagesDir,
	}
	i18nApp, i18nErr := i18n.New(i18nCfg)
	if i18nErr != nil {
		return fmt.Errorf("failed to create i18n: %w", i18nErr)
	}
	app.I18n = i18nApp
	return nil
}
