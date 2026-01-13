package parser

import (
	"context"
	"database/sql"
	"time"
)

// postgresParser PostgreSQL 数据库解析器
type postgresParser struct {
	baseParser
}

// NewPostgresParser 创建 PostgreSQL 解析器
func NewPostgresParser() Parser {
	return &postgresParser{
		baseParser: baseParser{dbType: DatabasePostgres},
	}
}

// ParseDatabase 解析整个数据库
func (p *postgresParser) ParseDatabase(ctx context.Context, db *sql.DB) (*Schema, error) {
	// 获取当前数据库名
	var dbName string
	if err := db.QueryRowContext(ctx, "SELECT current_database()").Scan(&dbName); err != nil {
		return nil, &ParseError{Message: "failed to get database name", Cause: err}
	}

	// 获取所有表名 (public schema)
	query := `SELECT tablename FROM pg_tables WHERE schemaname = 'public' ORDER BY tablename`
	rows, err := queryRows(ctx, db, query)
	if err != nil {
		return nil, &ParseError{Message: "failed to query tables", Cause: err}
	}
	defer rows.Close()

	var tables []*Table
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			continue
		}

		table, err := p.ParseTable(ctx, db, tableName)
		if err != nil {
			continue
		}
		tables = append(tables, table)
	}

	return &Schema{
		Name:         dbName,
		DatabaseType: DatabasePostgres,
		Tables:       tables,
		ParsedAt:     time.Now(),
	}, nil
}

// ParseTable 解析单个表
func (p *postgresParser) ParseTable(ctx context.Context, db *sql.DB, tableName string) (*Table, error) {
	table := &Table{
		Name: tableName,
	}

	// 获取表注释
	var comment sql.NullString
	err := db.QueryRowContext(ctx, `
		SELECT obj_description(oid) 
		FROM pg_class 
		WHERE relname = $1 AND relkind = 'r'
	`, tableName).Scan(&comment)
	if err == nil && comment.Valid {
		table.Comment = comment.String
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

	return table, nil
}

// parseColumns 解析表字段
func (p *postgresParser) parseColumns(ctx context.Context, db *sql.DB, tableName string) ([]*Column, []string, error) {
	query := `
		SELECT 
			a.attname AS column_name,
			pg_catalog.format_type(a.atttypid, a.atttypmod) AS data_type,
			NOT a.attnotnull AS is_nullable,
			pg_get_expr(d.adbin, d.adrelid) AS column_default,
			col_description(c.oid, a.attnum) AS column_comment,
			CASE WHEN pk.contype = 'p' THEN true ELSE false END AS is_primary_key
		FROM pg_attribute a
		JOIN pg_class c ON a.attrelid = c.oid
		LEFT JOIN pg_attrdef d ON a.attrelid = d.adrelid AND a.attnum = d.adnum
		LEFT JOIN pg_constraint pk ON c.oid = pk.conrelid AND a.attnum = ANY(pk.conkey) AND pk.contype = 'p'
		WHERE c.relname = $1
		AND a.attnum > 0
		AND NOT a.attisdropped
		ORDER BY a.attnum`

	rows, err := queryRows(ctx, db, query, tableName)
	if err != nil {
		return nil, nil, &ParseError{Table: tableName, Message: "failed to get columns", Cause: err}
	}
	defer rows.Close()

	var columns []*Column
	var primaryKeys []string

	for rows.Next() {
		var name, dataType string
		var nullable, isPK bool
		var defaultValue, comment sql.NullString

		if err := rows.Scan(&name, &dataType, &nullable, &defaultValue, &comment, &isPK); err != nil {
			continue
		}

		goType, _ := MapSQLTypeToGo(dataType, nullable)

		col := &Column{
			Name:         name,
			DataType:     dataType,
			GoType:       goType,
			Nullable:     nullable,
			IsPrimaryKey: isPK,
		}

		if defaultValue.Valid {
			col.DefaultValue = &defaultValue.String
			// 检查是否是序列 (自增)
			if len(defaultValue.String) > 7 && defaultValue.String[:7] == "nextval" {
				col.IsAutoIncrement = true
			}
		}
		if comment.Valid {
			col.Comment = comment.String
		}

		columns = append(columns, col)

		if isPK {
			primaryKeys = append(primaryKeys, name)
		}
	}

	return columns, primaryKeys, nil
}

// parseIndexes 解析表索引
func (p *postgresParser) parseIndexes(ctx context.Context, db *sql.DB, tableName string) ([]*Index, error) {
	query := `
		SELECT 
			i.relname AS index_name,
			a.attname AS column_name,
			ix.indisunique AS is_unique,
			ix.indisprimary AS is_primary
		FROM pg_index ix
		JOIN pg_class i ON ix.indexrelid = i.oid
		JOIN pg_class t ON ix.indrelid = t.oid
		JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = ANY(ix.indkey)
		WHERE t.relname = $1
		ORDER BY i.relname, a.attnum`

	rows, err := queryRows(ctx, db, query, tableName)
	if err != nil {
		return nil, &ParseError{Table: tableName, Message: "failed to get indexes", Cause: err}
	}
	defer rows.Close()

	indexMap := make(map[string]*Index)

	for rows.Next() {
		var indexName, columnName string
		var isUnique, isPrimary bool

		if err := rows.Scan(&indexName, &columnName, &isUnique, &isPrimary); err != nil {
			continue
		}

		if idx, ok := indexMap[indexName]; ok {
			idx.Columns = append(idx.Columns, columnName)
		} else {
			indexMap[indexName] = &Index{
				Name:      indexName,
				Columns:   []string{columnName},
				IsUnique:  isUnique,
				IsPrimary: isPrimary,
			}
		}
	}

	var indexes []*Index
	for _, idx := range indexMap {
		indexes = append(indexes, idx)
	}

	return indexes, nil
}
