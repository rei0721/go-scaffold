package sqlgen

import (
	"errors"
	"fmt"
)

// 预定义错误
var (
	// ErrDatabaseConnection 数据库连接失败
	ErrDatabaseConnection = errors.New("database connection failed")
	// ErrTableNotFound 表不存在
	ErrTableNotFound = errors.New("table not found")
	// ErrInvalidConfig 配置无效
	ErrInvalidConfig = errors.New("invalid configuration")
	// ErrTemplateRender 模板渲染失败
	ErrTemplateRender = errors.New("template render failed")
	// ErrFileWrite 文件写入失败
	ErrFileWrite = errors.New("file write failed")
	// ErrUnsupportedDatabase 不支持的数据库类型
	ErrUnsupportedDatabase = errors.New("unsupported database type")
	// ErrParseSchema Schema 解析失败
	ErrParseSchema = errors.New("schema parse failed")
)

// ParseError 解析错误
type ParseError struct {
	Table   string
	Column  string
	Message string
	Cause   error
}

func (e *ParseError) Error() string {
	if e.Column != "" {
		return fmt.Sprintf("parse error on table %s column %s: %s", e.Table, e.Column, e.Message)
	}
	if e.Table != "" {
		return fmt.Sprintf("parse error on table %s: %s", e.Table, e.Message)
	}
	return fmt.Sprintf("parse error: %s", e.Message)
}

func (e *ParseError) Unwrap() error {
	return e.Cause
}

// GenerateError 生成错误
type GenerateError struct {
	Table   string
	File    string
	Message string
	Cause   error
}

func (e *GenerateError) Error() string {
	if e.File != "" {
		return fmt.Sprintf("generate error for file %s: %s", e.File, e.Message)
	}
	if e.Table != "" {
		return fmt.Sprintf("generate error for table %s: %s", e.Table, e.Message)
	}
	return fmt.Sprintf("generate error: %s", e.Message)
}

func (e *GenerateError) Unwrap() error {
	return e.Cause
}

// ConfigError 配置错误
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("config error on field %s: %s", e.Field, e.Message)
}
