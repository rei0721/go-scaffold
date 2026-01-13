package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rei0721/rei0721/pkg/sqlgen"
)

func sqlDDLGenMySQL() {
	fmt.Println("================================")
	fmt.Println(" SQLGen MySQL DDL 脚本生成演示")
	fmt.Println("================================")
	fmt.Println()

	// 1. 连接数据库
	fmt.Println("📌 步骤 1: 连接数据库...")
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/testdb")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 2. 创建示例表
	fmt.Println("📌 步骤 2: 创建示例表...")
	createMySQLDemoTables(db)

	// 3. 配置生成器 - 启用 Migration 生成
	fmt.Println("📌 步骤 3: 配置生成器...")
	config := sqlgen.DefaultConfig()
	config.DatabaseType = sqlgen.DatabaseMySQL
	config.Target.Model = true
	config.Target.DAO = true
	config.Target.Migration = true // 启用 DDL 脚本生成
	config.OutputDir = "./generated"

	// 4. 创建生成器
	fmt.Println("📌 步骤 4: 创建生成器...")
	gen, err := sqlgen.NewGeneratorSimple(config)
	if err != nil {
		panic(err)
	}

	// 5. 解析数据库 Schema
	fmt.Println("📌 步骤 5: 解析数据库 Schema...")
	schema, err := gen.Parse(context.Background(), db)
	if err != nil {
		panic(err)
	}

	// 6. 生成代码 (包括 DDL)
	fmt.Println("📌 步骤 6: 生成代码和 DDL 脚本...")
	err = gen.Generate(context.Background(), schema, config.OutputDir)
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println("================================")
	fmt.Println("✅ MySQL DDL 代码生成成功！")
	fmt.Println("================================")
	fmt.Println("📁 输出目录:", config.OutputDir)
	fmt.Println("📊 生成表数:", len(schema.Tables))
	fmt.Println()
	fmt.Println("生成的文件：")
	for _, table := range schema.Tables {
		modelFile := "./generated/models/" + sqlgen.ToSnakeCase(table.Name) + ".go"
		daoFile := "./generated/dao/" + sqlgen.ToSnakeCase(table.Name) + "_dao.go"
		fmt.Println("  -", modelFile)
		fmt.Println("  -", daoFile)
	}
	fmt.Println("  - ./generated/schema.sql (MySQL DDL 脚本)")
	fmt.Println()

	// 7. 单独生成 DDL 演示
	fmt.Println("📌 附加: MySQL DDL 脚本预览...")
	ddlGen := sqlgen.NewDDLGenerator(config)
	ddl, _ := ddlGen.GenerateDDL(schema)
	fmt.Println("\n生成的 MySQL DDL 内容预览:")
	fmt.Println("------------------------------------------")
	fmt.Println(ddl)
}

// createMySQLDemoTables 创建 MySQL 演示测试表
func createMySQLDemoTables(db *sql.DB) {
	db.Exec("DROP TABLE IF EXISTS posts")
	db.Exec("DROP TABLE IF EXISTS users")

	_, err := db.Exec(`
		CREATE TABLE users (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			username VARCHAR(100) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			status TINYINT DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE posts (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			user_id BIGINT NOT NULL,
			title VARCHAR(200) NOT NULL,
			content TEXT,
			view_count INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`)
	if err != nil {
		panic(err)
	}

	fmt.Println("  ✅ 创建表: users")
	fmt.Println("  ✅ 创建表: posts")
}
