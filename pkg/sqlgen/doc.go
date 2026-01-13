// Package sqlgen 提供企业级 SQL 代码生成工具库
//
// 概述
//
// sqlgen 是一个通用的 SQL 代码生成器，支持从数据库 Schema 自动生成
// 类型安全的 Go 模型、DAO 层代码和查询构建器。
//
// 主要特性:
//   - 多数据库支持: MySQL, PostgreSQL, SQLite
//   - 类型安全: 自动将 SQL 类型映射为 Go 类型
//   - 灵活模板: 内置模板 + 支持自定义模板
//   - 企业特性: 软删除、时间戳、乐观锁
//   - 原子写入: 确保生成过程中断不会损坏文件
//
// 快速开始
//
//	import (
//	    "context"
//	    "database/sql"
//	    "github.com/rei0721/rei0721/pkg/sqlgen"
//	    _ "github.com/mattn/go-sqlite3"
//	)
//
//	func main() {
//	    db, _ := sql.Open("sqlite3", "mydb.db")
//	    defer db.Close()
//
//	    config := sqlgen.DefaultConfig()
//	    config.Target.Model = true
//	    config.Target.DAO = true
//
//	    gen, _ := sqlgen.NewGenerator(config)
//	    schema, _ := gen.Parse(context.Background(), db)
//	    gen.Generate(context.Background(), schema, "./generated")
//	}
//
// 配置
//
// 使用 Config 结构体配置生成器行为:
//   - DatabaseType: 数据库类型 (mysql/postgres/sqlite)
//   - OutputDir: 输出目录
//   - PackageName: 生成的包名
//   - Target: 生成目标 (Model/DAO/Query)
//   - Tags: 标签选项 (JSON/GORM/Validate)
//
// 架构
//
//	pkg/sqlgen/
//	├── parser/      # 数据库 Schema 解析器
//	├── template/    # 模板引擎
//	└── *.go         # 核心类型和接口
//
// 参见
//
// 详细文档请参考 pkg/sqlgen/README.md
package sqlgen
