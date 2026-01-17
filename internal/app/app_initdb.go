package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rei0721/go-scaffold/internal/models"
	rbacmodels "github.com/rei0721/go-scaffold/pkg/rbac/models"
	"github.com/rei0721/go-scaffold/pkg/sqlgen"
)

// initSqlGenerator 初始化 SQL 生成器
func initSqlGenerator(app *App) error {
	app.Sqlgen = sqlgen.New(&sqlgen.Config{
		Dialect: getDialectFromDriver(app.Config.Database.Driver),
		Pretty:  true,
	})
	return nil
}

// runInitDB 执行数据库初始化
func runInitDB(app *App) error {
	app.Logger.Info("starting database initialization...")

	// 1. 检查锁文件
	lockPath := filepath.Join(app.Config.InitDB.ScriptDir, app.Config.InitDB.LockFile)
	if _, err := os.Stat(lockPath); err == nil {
		app.Logger.Warn("database already initialized (lock file exists)",
			"lock_file", lockPath)
		app.Logger.Info("to reinitialize, delete the lock file and run again")
		return nil
	}

	// 2. 确保脚本目录存在
	if err := os.MkdirAll(app.Config.InitDB.ScriptDir, 0755); err != nil {
		return fmt.Errorf("failed to create script directory: %w", err)
	}

	// 3. 使用 sqlgen 生成 SQL 脚本
	gen := app.Sqlgen

	// 收集所有模型的建表语句
	allModels := []interface{}{
		&models.User{},
		&rbacmodels.Role{},
		&rbacmodels.Permission{},
		&rbacmodels.UserRole{},
		&rbacmodels.RolePermission{},
	}

	var sqlStatements []string
	for _, model := range allModels {
		sql, err := gen.Table(model)
		if err != nil {
			return fmt.Errorf("failed to generate SQL for model: %w", err)
		}
		sqlStatements = append(sqlStatements, sql)
	}

	// 4. 写入 SQL 文件
	scriptPath := filepath.Join(app.Config.InitDB.ScriptDir,
		fmt.Sprintf(ConstantsInitDBScriptFileTemplate, app.Config.InitDB.ScriptFilePrefix, app.Config.Database.Driver))

	fullSQL := strings.Join(sqlStatements, "\n\n")
	if err := os.WriteFile(scriptPath, []byte(fullSQL), 0644); err != nil {
		return fmt.Errorf("failed to write SQL script: %w", err)
	}
	app.Logger.Info("SQL script generated", "path", scriptPath)

	// 5. 执行 SQL 初始化
	db := app.DB.DB()
	for _, sql := range sqlStatements {
		if err := db.Exec(sql).Error; err != nil {
			app.Logger.Error("failed to execute SQL", "error", err, "sql", sql)
			return fmt.Errorf("failed to execute SQL: %w", err)
		}
	}
	app.Logger.Info("database tables created successfully")

	// 6. 创建锁文件
	if err := os.WriteFile(lockPath, []byte("initialized"), 0644); err != nil {
		app.Logger.Warn("failed to create lock file", "error", err)
	}

	app.Logger.Info("database initialization completed successfully")
	return nil
}

// getDialectFromDriver 将数据库驱动转换为 sqlgen 方言
func getDialectFromDriver(driver string) sqlgen.Dialect {
	switch driver {
	case "mysql":
		return sqlgen.MySQL
	case "postgres":
		return sqlgen.PostgreSQL
	case "sqlite":
		return sqlgen.SQLite
	default:
		return sqlgen.MySQL
	}
}
