package parser

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// sqliteParser SQLite 数据库解析器
type sqliteParser struct {
	baseParser
}

// NewSQLiteParser 创建 SQLite 解析器
func NewSQLiteParser() Parser {
	return &sqliteParser{
		baseParser: baseParser{dbType: DatabaseSQLite},
	}
}

// ParseDatabase 解析整个数据库
func (p *sqliteParser) ParseDatabase(ctx context.Context, db *sql.DB) (*Schema, error) {
	// 获取所有表名
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name`
	rows, err := queryRows(ctx, db, query)
	if err != nil {
		return nil, &ParseError{Message: "failed to query tables", Cause: err}
	}
	defer rows.Close()

	var tableNames []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, &ParseError{Message: "failed to scan table name", Cause: err}
		}
		tableNames = append(tableNames, name)
	}

	if err := rows.Err(); err != nil {
		return nil, &ParseError{Message: "error iterating tables", Cause: err}
	}

	// 解析每个表
	var tables []*Table
	for _, tableName := range tableNames {
		table, err := p.ParseTable(ctx, db, tableName)
		if err != nil {
			// 记录错误但继续处理其他表
			continue
		}
		tables = append(tables, table)
	}

	return &Schema{
		Name:         "main", // SQLite 默认数据库名
		DatabaseType: DatabaseSQLite,
		Tables:       tables,
		ParsedAt:     time.Now(),
	}, nil
}

// ParseTable 解析单个表
func (p *sqliteParser) ParseTable(ctx context.Context, db *sql.DB, tableName string) (*Table, error) {
	table := &Table{
		Name: tableName,
	}

	// 解析字段
	columns, primaryKeys, err := p.parseColumns(ctx, db, tableName)
	if err != nil {
		return nil, err
	}
	table.Columns = columns
	table.PrimaryKey = primaryKeys

	// 解析索引
	indexes, err := p.parseIndexes(ctx, db, tableName)
	if err != nil {
		return nil, err
	}
	table.Indexes = indexes

	// 解析外键
	foreignKeys, err := p.parseForeignKeys(ctx, db, tableName)
	if err != nil {
		return nil, err
	}
	table.ForeignKeys = foreignKeys

	return table, nil
}

// parseColumns 解析表字段
func (p *sqliteParser) parseColumns(ctx context.Context, db *sql.DB, tableName string) ([]*Column, []string, error) {
	query := fmt.Sprintf(`PRAGMA table_info("%s")`, tableName)
	rows, err := queryRows(ctx, db, query)
	if err != nil {
		return nil, nil, &ParseError{Table: tableName, Message: "failed to get columns", Cause: err}
	}
	defer rows.Close()

	var columns []*Column
	var primaryKeys []string

	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull, pk int
		var defaultValue sql.NullString

		if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk); err != nil {
			return nil, nil, &ParseError{
				Table:   tableName,
				Column:  name,
				Message: "failed to scan column",
				Cause:   err,
			}
		}

		// 映射 Go 类型
		nullable := notNull == 0
		goType, _ := MapSQLTypeToGo(dataType, nullable)

		col := &Column{
			Name:         name,
			DataType:     dataType,
			GoType:       goType,
			Nullable:     nullable,
			IsPrimaryKey: pk > 0,
		}

		if defaultValue.Valid {
			col.DefaultValue = &defaultValue.String
		}

		// 检查是否自增 (SQLite 中 INTEGER PRIMARY KEY 自动自增)
		if pk > 0 && strings.EqualFold(dataType, "INTEGER") {
			col.IsAutoIncrement = true
		}

		columns = append(columns, col)

		if pk > 0 {
			primaryKeys = append(primaryKeys, name)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, nil, &ParseError{Table: tableName, Message: "error iterating columns", Cause: err}
	}

	return columns, primaryKeys, nil
}

// parseIndexes 解析表索引
func (p *sqliteParser) parseIndexes(ctx context.Context, db *sql.DB, tableName string) ([]*Index, error) {
	query := fmt.Sprintf(`PRAGMA index_list("%s")`, tableName)
	rows, err := queryRows(ctx, db, query)
	if err != nil {
		return nil, &ParseError{Table: tableName, Message: "failed to get indexes", Cause: err}
	}
	defer rows.Close()

	var indexes []*Index

	for rows.Next() {
		var seq int
		var name, origin string
		var unique, partial int

		if err := rows.Scan(&seq, &name, &unique, &origin, &partial); err != nil {
			continue // 跳过解析失败的索引
		}

		// 获取索引列
		indexColumns, err := p.parseIndexColumns(ctx, db, name)
		if err != nil {
			continue
		}

		indexes = append(indexes, &Index{
			Name:      name,
			Columns:   indexColumns,
			IsUnique:  unique == 1,
			IsPrimary: origin == "pk",
		})
	}

	return indexes, nil
}

// parseIndexColumns 解析索引包含的字段
func (p *sqliteParser) parseIndexColumns(ctx context.Context, db *sql.DB, indexName string) ([]string, error) {
	query := fmt.Sprintf(`PRAGMA index_info("%s")`, indexName)
	rows, err := queryRows(ctx, db, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var seqno, cid int
		var name sql.NullString

		if err := rows.Scan(&seqno, &cid, &name); err != nil {
			continue
		}
		if name.Valid {
			columns = append(columns, name.String)
		}
	}

	return columns, nil
}

// parseForeignKeys 解析外键
func (p *sqliteParser) parseForeignKeys(ctx context.Context, db *sql.DB, tableName string) ([]*ForeignKey, error) {
	query := fmt.Sprintf(`PRAGMA foreign_key_list("%s")`, tableName)
	rows, err := queryRows(ctx, db, query)
	if err != nil {
		return nil, &ParseError{Table: tableName, Message: "failed to get foreign keys", Cause: err}
	}
	defer rows.Close()

	var foreignKeys []*ForeignKey

	for rows.Next() {
		var id, seq int
		var table, from, to, onUpdate, onDelete, match string

		if err := rows.Scan(&id, &seq, &table, &from, &to, &onUpdate, &onDelete, &match); err != nil {
			continue
		}

		foreignKeys = append(foreignKeys, &ForeignKey{
			Name:      fmt.Sprintf("fk_%s_%s", tableName, from),
			Column:    from,
			RefTable:  table,
			RefColumn: to,
			OnUpdate:  onUpdate,
			OnDelete:  onDelete,
		})
	}

	return foreignKeys, nil
}
