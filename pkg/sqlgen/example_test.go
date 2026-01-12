package sqlgen_test

import (
	"fmt"
	"log"
	"os"

	"github.com/rei0721/rei0721/internal/models"
	"github.com/rei0721/rei0721/pkg/sqlgen"
)

// Example_generator 演示如何使用 SQL 生成器
func Example_generator() {
	// 创建 PostgreSQL 方言的生成器
	dialect := sqlgen.NewPostgresDialect()
	generator := sqlgen.New(dialect)

	// 生成 User 模型的 SQL
	result, err := generator.GenerateSQL(models.User{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("表名:", result.TableName)
	fmt.Println("\n建表 SQL:")
	fmt.Println(result.CreateTable)
	fmt.Println("\n插入 SQL:")
	fmt.Println(result.Insert)
	fmt.Println("\n查询 SQL:")
	fmt.Println(result.Select)
	fmt.Println("\n更新 SQL:")
	fmt.Println(result.Update)
	fmt.Println("\n删除 SQL:")
	fmt.Println(result.Delete)
}

// Example_fileGenerator 演示如何生成 SQL 文件
func Example_fileGenerator() {
	// 创建 MySQL 方言的生成器
	dialect := sqlgen.NewMySQLDialect()
	generator := sqlgen.New(dialect)

	// 创建文件生成器
	outputDir := "./generated_sql"
	fileGenerator := sqlgen.NewFileGenerator(generator, outputDir)

	// 生成选项
	options := &sqlgen.GenerateOptions{
		OutputDir:       outputDir,
		SeparateFiles:   true, // 生成分离的文件
		GenerateSummary: true, // 生成汇总文件
		IncludeComments: true, // 包含注释
	}

	// 定义要生成的模型
	models := []interface{}{
		models.User{},
		// 可以添加更多模型
	}

	// 生成文件
	if err := fileGenerator.GenerateWithOptions(options, models...); err != nil {
		log.Fatal(err)
	}

	fmt.Println("SQL 文件生成完成!")

	// 清理测试文件
	os.RemoveAll(outputDir)
}

// Example_multipleDialects 演示如何为多个数据库生成 SQL
func Example_multipleDialects() {
	// 定义模型
	user := models.User{}

	// 支持的数据库方言
	dialects := map[string]sqlgen.Dialect{
		"postgres": sqlgen.NewPostgresDialect(),
		"mysql":    sqlgen.NewMySQLDialect(),
		"sqlite":   sqlgen.NewSQLiteDialect(),
	}

	// 为每个数据库生成 SQL
	for name, dialect := range dialects {
		fmt.Printf("\n=== %s ===\n", name)

		generator := sqlgen.New(dialect)
		result, err := generator.GenerateSQL(user)
		if err != nil {
			log.Printf("生成 %s SQL 失败: %v", name, err)
			continue
		}

		fmt.Printf("建表 SQL:\n%s\n", result.CreateTable)
	}
}
