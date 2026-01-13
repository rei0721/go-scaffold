package sqlgen

import (
	"context"
	"database/sql"

	"github.com/rei0721/rei0721/pkg/sqlgen/parser"
	"github.com/rei0721/rei0721/pkg/sqlgen/template"
)

// parserAdapter 适配器，将 parser.Parser 适配为 sqlgen.Parser
type parserAdapter struct {
	p parser.Parser
}

func (a *parserAdapter) ParseDatabase(ctx context.Context, db *sql.DB) (*Schema, error) {
	schema, err := a.p.ParseDatabase(ctx, db)
	if err != nil {
		return nil, err
	}
	return convertSchema(schema), nil
}

func (a *parserAdapter) ParseTable(ctx context.Context, db *sql.DB, tableName string) (*Table, error) {
	table, err := a.p.ParseTable(ctx, db, tableName)
	if err != nil {
		return nil, err
	}
	return convertTable(table), nil
}

func (a *parserAdapter) GetDatabaseType() DatabaseType {
	return DatabaseType(a.p.GetDatabaseType())
}

// templateAdapter 适配器，将 template.TemplateEngine 适配为 sqlgen.TemplateEngine
type templateAdapter struct {
	t template.TemplateEngine
}

func (a *templateAdapter) Render(tmplName string, data interface{}) (string, error) {
	return a.t.Render(tmplName, data)
}

func (a *templateAdapter) RegisterFunc(name string, fn interface{}) error {
	return a.t.RegisterFunc(name, fn)
}

func (a *templateAdapter) LoadTemplate(name string, content string) error {
	return a.t.LoadTemplate(name, content)
}

func (a *templateAdapter) LoadTemplateFile(name string, path string) error {
	return a.t.LoadTemplateFile(name, path)
}

// convertSchema 转换 parser.Schema 为 sqlgen.Schema
func convertSchema(s *parser.Schema) *Schema {
	if s == nil {
		return nil
	}

	var tables []*Table
	for _, t := range s.Tables {
		tables = append(tables, convertTable(t))
	}

	return &Schema{
		Name:         s.Name,
		DatabaseType: DatabaseType(s.DatabaseType),
		Tables:       tables,
		ParsedAt:     s.ParsedAt,
	}
}

// convertTable 转换 parser.Table 为 sqlgen.Table
func convertTable(t *parser.Table) *Table {
	if t == nil {
		return nil
	}

	var columns []*Column
	for _, c := range t.Columns {
		columns = append(columns, convertColumn(c))
	}

	var indexes []*Index
	for _, i := range t.Indexes {
		indexes = append(indexes, convertIndex(i))
	}

	var foreignKeys []*ForeignKey
	for _, fk := range t.ForeignKeys {
		foreignKeys = append(foreignKeys, convertForeignKey(fk))
	}

	return &Table{
		Name:        t.Name,
		Comment:     t.Comment,
		Columns:     columns,
		PrimaryKey:  t.PrimaryKey,
		Indexes:     indexes,
		ForeignKeys: foreignKeys,
	}
}

// convertColumn 转换字段
func convertColumn(c *parser.Column) *Column {
	if c == nil {
		return nil
	}

	return &Column{
		Name:            c.Name,
		DataType:        c.DataType,
		GoType:          c.GoType,
		Nullable:        c.Nullable,
		DefaultValue:    c.DefaultValue,
		Comment:         c.Comment,
		IsPrimaryKey:    c.IsPrimaryKey,
		IsAutoIncrement: c.IsAutoIncrement,
		Length:          c.Length,
		Precision:       c.Precision,
		Scale:           c.Scale,
	}
}

// convertIndex 转换索引
func convertIndex(i *parser.Index) *Index {
	if i == nil {
		return nil
	}

	return &Index{
		Name:      i.Name,
		Columns:   i.Columns,
		IsUnique:  i.IsUnique,
		IsPrimary: i.IsPrimary,
	}
}

// convertForeignKey 转换外键
func convertForeignKey(fk *parser.ForeignKey) *ForeignKey {
	if fk == nil {
		return nil
	}

	return &ForeignKey{
		Name:      fk.Name,
		Column:    fk.Column,
		RefTable:  fk.RefTable,
		RefColumn: fk.RefColumn,
		OnDelete:  fk.OnDelete,
		OnUpdate:  fk.OnUpdate,
	}
}

// NewGeneratorSimple 创建生成器的简化版本（自动设置适配器）
func NewGeneratorSimple(config *Config) (Generator, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	// 创建 parser
	p, err := parser.NewParser(parser.DatabaseType(config.DatabaseType))
	if err != nil {
		return nil, err
	}

	// 创建 template
	t := template.NewEngine()

	return &generator{
		config:   config,
		parser:   &parserAdapter{p: p},
		template: &templateAdapter{t: t},
		writer:   NewFileWriter(),
	}, nil
}
