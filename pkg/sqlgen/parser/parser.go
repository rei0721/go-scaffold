// Package parser 提供数据库 Schema 解析能力
// 支持 MySQL, PostgreSQL, SQLite 多种数据库
package parser

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// DatabaseType 表示数据库类型 (本地定义避免循环引用)
type DatabaseType string

const (
	DatabaseMySQL    DatabaseType = "mysql"
	DatabasePostgres DatabaseType = "postgres"
	DatabaseSQLite   DatabaseType = "sqlite"
)

// Parser 数据库 Schema 解析器接口 (本地定义)
type Parser interface {
	ParseDatabase(ctx context.Context, db *sql.DB) (*Schema, error)
	ParseTable(ctx context.Context, db *sql.DB, tableName string) (*Table, error)
	GetDatabaseType() DatabaseType
}

// Schema 表示数据库 Schema
type Schema struct {
	Name         string
	DatabaseType DatabaseType
	Tables       []*Table
	ParsedAt     time.Time
}

// Table 表示数据库表结构
type Table struct {
	Name        string
	Comment     string
	Columns     []*Column
	PrimaryKey  []string
	Indexes     []*Index
	ForeignKeys []*ForeignKey
}

// Column 表示表字段
type Column struct {
	Name            string
	DataType        string
	GoType          string
	Nullable        bool
	DefaultValue    *string
	Comment         string
	IsPrimaryKey    bool
	IsAutoIncrement bool
	Length          int
	Precision       int
	Scale           int
}

// Index 表示索引
type Index struct {
	Name      string
	Columns   []string
	IsUnique  bool
	IsPrimary bool
}

// ForeignKey 表示外键约束
type ForeignKey struct {
	Name      string
	Column    string
	RefTable  string
	RefColumn string
	OnDelete  string
	OnUpdate  string
}

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

// NewParser 根据数据库类型创建解析器
func NewParser(dbType DatabaseType) (Parser, error) {
	switch dbType {
	case DatabaseMySQL:
		return NewMySQLParser(), nil
	case DatabasePostgres:
		return NewPostgresParser(), nil
	case DatabaseSQLite:
		return NewSQLiteParser(), nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

// baseParser 基础解析器（提供通用方法）
type baseParser struct {
	dbType DatabaseType
}

// GetDatabaseType 返回数据库类型
func (p *baseParser) GetDatabaseType() DatabaseType {
	return p.dbType
}

// queryRows 执行查询并返回结果
func queryRows(ctx context.Context, db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	return rows, nil
}

// queryRow 执行查询并返回单行结果
func queryRow(ctx context.Context, db *sql.DB, query string, args ...interface{}) *sql.Row {
	return db.QueryRowContext(ctx, query, args...)
}
