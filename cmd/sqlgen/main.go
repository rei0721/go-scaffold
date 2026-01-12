// Package main 提供了 SQL 生成工具的命令行接口
// 用于根据 GORM 模型生成数据库 SQL 语句
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rei0721/rei0721/internal/models"
	"github.com/rei0721/rei0721/pkg/sqlgen"
)

func main() {
	// 定义命令行参数
	var (
		dialect   = flag.String("dialect", "postgres", "数据库方言 (postgres, mysql, sqlite)")
		outputDir = flag.String("output", "./sql", "输出目录")
		separate  = flag.Bool("separate", false, "是否生成分离的 SQL 文件")
		summary   = flag.Bool("summary", true, "是否生成汇总文件")
		comments  = flag.Bool("comments", true, "是否包含注释")
		help      = flag.Bool("help", false, "显示帮助信息")
	)

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// 创建输出目录
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("创建输出目录失败: %v", err)
	}

	// 选择数据库方言
	var dialectImpl sqlgen.Dialect
	switch *dialect {
	case "postgres":
		dialectImpl = sqlgen.NewPostgresDialect()
	case "mysql":
		dialectImpl = sqlgen.NewMySQLDialect()
	case "sqlite":
		dialectImpl = sqlgen.NewSQLiteDialect()
	default:
		log.Fatalf("不支持的数据库方言: %s", *dialect)
	}

	// 创建生成器
	generator := sqlgen.New(dialectImpl)
	fileGenerator := sqlgen.NewFileGenerator(generator, *outputDir)

	// 定义要生成 SQL 的模型
	models := []interface{}{
		models.User{},
		// 在这里添加更多模型
	}

	// 生成选项
	options := &sqlgen.GenerateOptions{
		OutputDir:       *outputDir,
		SeparateFiles:   *separate,
		GenerateSummary: *summary,
		IncludeComments: *comments,
	}

	// 生成 SQL 文件
	if err := fileGenerator.GenerateWithOptions(options, models...); err != nil {
		log.Fatalf("生成 SQL 文件失败: %v", err)
	}

	fmt.Printf("✅ SQL 文件生成成功!\n")
	fmt.Printf("📁 输出目录: %s\n", *outputDir)
	fmt.Printf("🗄️  数据库方言: %s\n", *dialect)

	// 显示生成的文件
	showGeneratedFiles(*outputDir)
}

// showHelp 显示帮助信息
func showHelp() {
	fmt.Println("SQL 生成工具 - 根据 GORM 模型生成数据库 SQL 语句")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  go run cmd/sqlgen/main.go [选项]")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -dialect string")
	fmt.Println("        数据库方言 (postgres, mysql, sqlite) (默认: postgres)")
	fmt.Println("  -output string")
	fmt.Println("        输出目录 (默认: ./sql)")
	fmt.Println("  -separate")
	fmt.Println("        是否生成分离的 SQL 文件 (默认: false)")
	fmt.Println("  -summary")
	fmt.Println("        是否生成汇总文件 (默认: true)")
	fmt.Println("  -comments")
	fmt.Println("        是否包含注释 (默认: true)")
	fmt.Println("  -help")
	fmt.Println("        显示帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  # 生成 PostgreSQL SQL 文件")
	fmt.Println("  go run cmd/sqlgen/main.go -dialect postgres")
	fmt.Println()
	fmt.Println("  # 生成 MySQL SQL 文件到指定目录")
	fmt.Println("  go run cmd/sqlgen/main.go -dialect mysql -output ./mysql_sql")
	fmt.Println()
	fmt.Println("  # 生成分离的 SQLite SQL 文件")
	fmt.Println("  go run cmd/sqlgen/main.go -dialect sqlite -separate")
}

// showGeneratedFiles 显示生成的文件
func showGeneratedFiles(outputDir string) {
	fmt.Println("\n📄 生成的文件:")

	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".sql" {
			relPath, _ := filepath.Rel(outputDir, path)
			fmt.Printf("  - %s\n", relPath)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("⚠️  列出文件时出错: %v\n", err)
	}
}
