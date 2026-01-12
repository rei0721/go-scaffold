package sqlgen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileGenerator 文件生成器
type FileGenerator struct {
	generator *Generator
	outputDir string
}

// NewFileGenerator 创建文件生成器
func NewFileGenerator(generator *Generator, outputDir string) *FileGenerator {
	return &FileGenerator{
		generator: generator,
		outputDir: outputDir,
	}
}

// GenerateFiles 生成 SQL 文件
func (fg *FileGenerator) GenerateFiles(models ...interface{}) error {
	// 确保输出目录存在
	if err := os.MkdirAll(fg.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 为每个模型生成 SQL 文件
	for _, model := range models {
		if err := fg.generateModelFiles(model); err != nil {
			return fmt.Errorf("failed to generate files for model: %w", err)
		}
	}

	// 生成汇总文件
	if err := fg.generateSummaryFile(models...); err != nil {
		return fmt.Errorf("failed to generate summary file: %w", err)
	}

	return nil
}

// generateModelFiles 为单个模型生成 SQL 文件
func (fg *FileGenerator) generateModelFiles(model interface{}) error {
	result, err := fg.generator.GenerateSQL(model)
	if err != nil {
		return err
	}

	// 生成建表 SQL 文件
	if err := fg.writeFile(fmt.Sprintf("%s_create.sql", result.TableName), result.CreateTable); err != nil {
		return err
	}

	// 生成 CRUD SQL 文件
	crudSQL := fg.buildCRUDSQL(result)
	if err := fg.writeFile(fmt.Sprintf("%s_crud.sql", result.TableName), crudSQL); err != nil {
		return err
	}

	return nil
}

// buildCRUDSQL 构建 CRUD SQL 内容
func (fg *FileGenerator) buildCRUDSQL(result *SQLResult) string {
	var content strings.Builder

	content.WriteString(fmt.Sprintf("-- %s 表 CRUD 操作 SQL\n", result.TableName))
	content.WriteString(fmt.Sprintf("-- 生成时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	content.WriteString("-- ==================== 插入操作 ====================\n")
	content.WriteString(result.Insert)
	content.WriteString("\n\n")

	content.WriteString("-- ==================== 查询操作 ====================\n")
	content.WriteString(result.Select)
	content.WriteString("\n\n")

	content.WriteString("-- ==================== 更新操作 ====================\n")
	content.WriteString(result.Update)
	content.WriteString("\n\n")

	content.WriteString("-- ==================== 删除操作 ====================\n")
	content.WriteString(result.Delete)
	content.WriteString("\n")

	return content.String()
}

// generateSummaryFile 生成汇总文件
func (fg *FileGenerator) generateSummaryFile(models ...interface{}) error {
	var content strings.Builder

	content.WriteString("-- 数据库初始化脚本\n")
	content.WriteString(fmt.Sprintf("-- 生成时间: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	content.WriteString("-- 此文件包含所有表的建表语句\n\n")

	// 生成所有建表语句
	for _, model := range models {
		result, err := fg.generator.GenerateSQL(model)
		if err != nil {
			return err
		}

		content.WriteString(fmt.Sprintf("-- ==================== %s 表 ====================\n", result.TableName))
		content.WriteString(result.CreateTable)
		content.WriteString("\n\n")
	}

	return fg.writeFile("init_database.sql", content.String())
}

// writeFile 写入文件
func (fg *FileGenerator) writeFile(filename, content string) error {
	filePath := filepath.Join(fg.outputDir, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}

	return nil
}

// GenerateWithOptions 使用选项生成文件
func (fg *FileGenerator) GenerateWithOptions(options *GenerateOptions, models ...interface{}) error {
	if options == nil {
		options = &GenerateOptions{}
	}

	// 设置默认值
	if options.OutputDir == "" {
		options.OutputDir = fg.outputDir
	}

	// 更新输出目录
	fg.outputDir = options.OutputDir

	// 确保输出目录存在
	if err := os.MkdirAll(fg.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 生成文件
	for _, model := range models {
		result, err := fg.generator.GenerateSQL(model)
		if err != nil {
			return fmt.Errorf("failed to generate SQL for model: %w", err)
		}

		// 根据选项生成不同的文件
		if options.SeparateFiles {
			if err := fg.generateSeparateFiles(result, options); err != nil {
				return err
			}
		} else {
			if err := fg.generateCombinedFile(result, options); err != nil {
				return err
			}
		}
	}

	// 生成汇总文件
	if options.GenerateSummary {
		if err := fg.generateSummaryFile(models...); err != nil {
			return err
		}
	}

	return nil
}

// GenerateOptions 生成选项
type GenerateOptions struct {
	OutputDir       string // 输出目录
	SeparateFiles   bool   // 是否分离文件
	GenerateSummary bool   // 是否生成汇总文件
	IncludeComments bool   // 是否包含注释
}

// generateSeparateFiles 生成分离的文件
func (fg *FileGenerator) generateSeparateFiles(result *SQLResult, options *GenerateOptions) error {
	// 生成建表文件
	createContent := fg.buildFileContent("CREATE TABLE", result.CreateTable, options)
	if err := fg.writeFile(fmt.Sprintf("%s_create.sql", result.TableName), createContent); err != nil {
		return err
	}

	// 生成插入文件
	insertContent := fg.buildFileContent("INSERT", result.Insert, options)
	if err := fg.writeFile(fmt.Sprintf("%s_insert.sql", result.TableName), insertContent); err != nil {
		return err
	}

	// 生成查询文件
	selectContent := fg.buildFileContent("SELECT", result.Select, options)
	if err := fg.writeFile(fmt.Sprintf("%s_select.sql", result.TableName), selectContent); err != nil {
		return err
	}

	// 生成更新文件
	updateContent := fg.buildFileContent("UPDATE", result.Update, options)
	if err := fg.writeFile(fmt.Sprintf("%s_update.sql", result.TableName), updateContent); err != nil {
		return err
	}

	// 生成删除文件
	deleteContent := fg.buildFileContent("DELETE", result.Delete, options)
	if err := fg.writeFile(fmt.Sprintf("%s_delete.sql", result.TableName), deleteContent); err != nil {
		return err
	}

	return nil
}

// generateCombinedFile 生成合并的文件
func (fg *FileGenerator) generateCombinedFile(result *SQLResult, options *GenerateOptions) error {
	content := fg.buildCRUDSQL(result)
	return fg.writeFile(fmt.Sprintf("%s_all.sql", result.TableName), content)
}

// buildFileContent 构建文件内容
func (fg *FileGenerator) buildFileContent(operation, sql string, options *GenerateOptions) string {
	var content strings.Builder

	if options.IncludeComments {
		content.WriteString(fmt.Sprintf("-- %s 操作 SQL\n", operation))
		content.WriteString(fmt.Sprintf("-- 生成时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
	}

	content.WriteString(sql)
	content.WriteString("\n")

	return content.String()
}
