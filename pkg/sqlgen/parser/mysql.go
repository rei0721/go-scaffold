package parser

import (
	"context"
	"database/sql"
	"time"
)

// mysqlParser MySQL 数据库解析器
type mysqlParser struct {
	baseParser
}

// NewMySQLParser 创建 MySQL 解析器
func NewMySQLParser() Parser {
	return &mysqlParser{
		baseParser: baseParser{dbType: DatabaseMySQL},
	}
}

// ParseDatabase 解析整个数据库
func (p *mysqlParser) ParseDatabase(ctx context.Context, db *sql.DB) (*Schema, error) {
	// 获取当前数据库名
	var dbName string
	if err := db.QueryRowContext(ctx, "SELECT DATABASE()").Scan(&dbName); err != nil {
		return nil, &ParseError{Message: "failed to get database name", Cause: err}
	}

	// 获取所有表名
	query := `SELECT TABLE_NAME, TABLE_COMMENT FROM information_schema.TABLES 
              WHERE TABLE_SCHEMA = ? AND TABLE_TYPE = 'BASE TABLE'`
	rows, err := queryRows(ctx, db, query, dbName)
	if err != nil {
		return nil, &ParseError{Message: "failed to query tables", Cause: err}
	}
	defer rows.Close()

	var tables []*Table
	for rows.Next() {
		var tableName, comment string
		if err := rows.Scan(&tableName, &comment); err != nil {
			continue
		}

		table, err := p.ParseTable(ctx, db, tableName)
		if err != nil {
			continue
		}
		table.Comment = comment
		tables = append(tables, table)
	}

	return &Schema{
		Name:         dbName,
		DatabaseType: DatabaseMySQL,
		Tables:       tables,
		ParsedAt:     time.Now(),
	}, nil
}

// ParseTable 解析单个表
func (p *mysqlParser) ParseTable(ctx context.Context, db *sql.DB, tableName string) (*Table, error) {
	table := &Table{
		Name: tableName,
	}

	// 获取当前数据库名
	var dbName string
	if err := db.QueryRowContext(ctx, "SELECT DATABASE()").Scan(&dbName); err != nil {
		return nil, &ParseError{Table: tableName, Message: "failed to get database name", Cause: err}
	}

	// 解析字段
	columns, primaryKeys, err := p.parseColumns(ctx, db, dbName, tableName)
	if err != nil {
		return nil, err
	}
	table.Columns = columns
	table.PrimaryKey = primaryKeys

	// 解析索引
	indexes, err := p.parseIndexes(ctx, db, dbName, tableName)
	if err != nil {
		return nil, err
	}
	table.Indexes = indexes

	return table, nil
}

// parseColumns 解析表字段
func (p *mysqlParser) parseColumns(ctx context.Context, db *sql.DB, dbName, tableName string) ([]*Column, []string, error) {
	query := `SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_DEFAULT, COLUMN_COMMENT, 
              COLUMN_KEY, EXTRA, CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE
              FROM information_schema.COLUMNS 
              WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? 
              ORDER BY ORDINAL_POSITION`

	rows, err := queryRows(ctx, db, query, dbName, tableName)
	if err != nil {
		return nil, nil, &ParseError{Table: tableName, Message: "failed to get columns", Cause: err}
	}
	defer rows.Close()

	var columns []*Column
	var primaryKeys []string

	for rows.Next() {
		var name, dataType, nullable, key, extra string
		var defaultValue, comment sql.NullString
		var charMaxLen, numPrecision, numScale sql.NullInt64

		if err := rows.Scan(&name, &dataType, &nullable, &defaultValue, &comment, &key, &extra, &charMaxLen, &numPrecision, &numScale); err != nil {
			continue
		}

		isNullable := nullable == "YES"
		goType, _ := MapSQLTypeToGo(dataType, isNullable)

		col := &Column{
			Name:            name,
			DataType:        dataType,
			GoType:          goType,
			Nullable:        isNullable,
			IsPrimaryKey:    key == "PRI",
			IsAutoIncrement: extra == "auto_increment",
		}

		if defaultValue.Valid {
			col.DefaultValue = &defaultValue.String
		}
		if comment.Valid {
			col.Comment = comment.String
		}
		if charMaxLen.Valid {
			col.Length = int(charMaxLen.Int64)
		}
		if numPrecision.Valid {
			col.Precision = int(numPrecision.Int64)
		}
		if numScale.Valid {
			col.Scale = int(numScale.Int64)
		}

		columns = append(columns, col)

		if key == "PRI" {
			primaryKeys = append(primaryKeys, name)
		}
	}

	return columns, primaryKeys, nil
}

// parseIndexes 解析表索引
func (p *mysqlParser) parseIndexes(ctx context.Context, db *sql.DB, dbName, tableName string) ([]*Index, error) {
	query := `SELECT INDEX_NAME, COLUMN_NAME, NON_UNIQUE 
              FROM information_schema.STATISTICS 
              WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? 
              ORDER BY INDEX_NAME, SEQ_IN_INDEX`

	rows, err := queryRows(ctx, db, query, dbName, tableName)
	if err != nil {
		return nil, &ParseError{Table: tableName, Message: "failed to get indexes", Cause: err}
	}
	defer rows.Close()

	indexMap := make(map[string]*Index)

	for rows.Next() {
		var indexName, columnName string
		var nonUnique int

		if err := rows.Scan(&indexName, &columnName, &nonUnique); err != nil {
			continue
		}

		if idx, ok := indexMap[indexName]; ok {
			idx.Columns = append(idx.Columns, columnName)
		} else {
			indexMap[indexName] = &Index{
				Name:      indexName,
				Columns:   []string{columnName},
				IsUnique:  nonUnique == 0,
				IsPrimary: indexName == "PRIMARY",
			}
		}
	}

	var indexes []*Index
	for _, idx := range indexMap {
		indexes = append(indexes, idx)
	}

	return indexes, nil
}
