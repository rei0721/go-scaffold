package main

import (
	"context"
	"database/sql"

	"github.com/rei0721/rei0721/pkg/sqlgen"
)

func sqlModelGenSqlite() {

	println("================================")
	println(" SQLGen 代码生成演示")
	println("================================")
	println()

	// 1. 连接数据库
	println("📌 步骤 1: 连接数据库...")
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 2. 创建示例表
	println("📌 步骤 2: 创建示例表...")
	createTestTables(db)

	// 3. 配置生成器
	println("📌 步骤 3: 配置生成器...")
	config := sqlgen.DefaultConfig()
	config.DatabaseType = sqlgen.DatabaseSQLite
	config.Target.Model = true
	config.Target.DAO = true
	config.OutputDir = "./generated"

	// 4. 创建生成器
	println("📌 步骤 4: 创建生成器...")
	gen, err := sqlgen.NewGeneratorSimple(config)
	if err != nil {
		panic(err)
	}

	// 5. 解析数据库 Schema
	println("📌 步骤 5: 解析数据库 Schema...")
	schema, err := gen.Parse(context.Background(), db)
	if err != nil {
		panic(err)
	}

	// 6. 生成代码
	println("📌 步骤 6: 生成代码...")
	err = gen.Generate(context.Background(), schema, config.OutputDir)
	if err != nil {
		panic(err)
	}

	println()
	println("================================")
	println("✅ 代码生成成功！")
	println("================================")
	println("📁 输出目录:", config.OutputDir)
	println("📊 生成表数:", len(schema.Tables))
	println()
	println("生成的文件：")
	for _, table := range schema.Tables {
		modelFile := "./generated/models/" + sqlgen.ToSnakeCase(table.Name) + ".go"
		daoFile := "./generated/dao/" + sqlgen.ToSnakeCase(table.Name) + "_dao.go"
		println("  -", modelFile)
		println("  -", daoFile)
	}
	println()
}

// createTestTables 创建示例测试表
func createTestTables(db *sql.DB) {
	// 删除已存在的表
	db.Exec("DROP TABLE IF EXISTS posts")
	db.Exec("DROP TABLE IF EXISTS users")

	// 创建 users 表
	_, err := db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(100) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			status INTEGER DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		panic(err)
	}

	// 创建 posts 表
	_, err = db.Exec(`
		CREATE TABLE posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title VARCHAR(200) NOT NULL,
			content TEXT,
			view_count INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		panic(err)
	}

	println("  ✅ 创建表: users")
	println("  ✅ 创建表: posts")
}
