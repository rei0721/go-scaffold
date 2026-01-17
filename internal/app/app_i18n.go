package app

import (
	"fmt"

	"github.com/rei0721/go-scaffold/pkg/i18n"
	"github.com/rei0721/go-scaffold/pkg/utils"
)

// initI18n 初始化i18n
func initI18n(app *App) error {
	// 初始化i18n
	i18nCfg := &i18n.Config{
		DefaultLanguage:    app.Config.I18n.Default,
		SupportedLanguages: app.Config.I18n.Supported,
		MessagesDir:        app.Config.I18n.MessagesDir,
	}
	i18nApp, i18nErr := i18n.New(i18nCfg)
	if i18nErr != nil {
		return fmt.Errorf("failed to create i18n: %w", i18nErr)
	}
	app.I18n = i18nApp
	app.I18nUtils = utils.NewI18nUtils(i18nApp, app.Config.I18n.Default)
	return nil
}

func (a *App) UI18n(messageID string, templates ...map[string]interface{}) string {
	return a.I18nUtils.T(messageID, templates...)
}
