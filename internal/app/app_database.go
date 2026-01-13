package app

import (
	"fmt"

	"github.com/rei0721/rei0721/pkg/database"
)

// initDatabase 初始化数据库连接
func initDatabase(app *App) error {
	db, err := database.New(&database.Config{
		Driver:       database.Driver(app.Config.Database.Driver),
		Host:         app.Config.Database.Host,
		Port:         app.Config.Database.Port,
		User:         app.Config.Database.User,
		Password:     app.Config.Database.Password,
		DBName:       app.Config.Database.DBName,
		MaxOpenConns: app.Config.Database.MaxOpenConns,
		MaxIdleConns: app.Config.Database.MaxIdleConns,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	app.DB = db
	app.Logger.Info("database connected successfully")
	return nil
}
