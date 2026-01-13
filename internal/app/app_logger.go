package app

import (
	"fmt"

	"github.com/rei0721/rei0721/pkg/logger"
)

// initLogger 初始化日志记录器
func initLogger(app *App) error {
	log, err := logger.New(&logger.Config{
		Level:         app.Config.Logger.Level,         // 从配置读取日志级别
		Format:        app.Config.Logger.Format,        // 从配置读取默认日志格式
		ConsoleFormat: app.Config.Logger.ConsoleFormat, // 从配置读取控制台专用格式
		FileFormat:    app.Config.Logger.FileFormat,    // 从配置读取文件专用格式
		Output:        app.Config.Logger.Output,        // 从配置读取输出目标
		FilePath:      app.Config.Logger.FilePath,      // 从配置读取日志文件路径
		MaxSize:       app.Config.Logger.MaxSize,       // 从配置读取日志文件最大大小
		MaxBackups:    app.Config.Logger.MaxBackups,    // 从配置读取日志文件最大备份数
		MaxAge:        app.Config.Logger.MaxAge,        // 从配置读取日志文件最大年龄
	})
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	app.Logger = log
	app.Logger.Info("logger initialized successfully")
	return nil
}
