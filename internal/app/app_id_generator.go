package app

import (
	"github.com/rei0721/go-scaffold/pkg/utils"
)

// InitIDGenerator 初始化ID生成器
func (app *App) InitIDGenerator() error {
	app.IDGenerator = utils.DefaultSnowflake()
	return nil
}
