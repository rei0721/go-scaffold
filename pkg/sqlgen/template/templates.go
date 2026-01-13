package template

// 内置模板定义

// modelTemplate 模型生成模板
const modelTemplate = `{{ .Header }}

package {{ .PackageName }}
{{ if .Imports }}
import (
{{- range .Imports }}
	"{{ . }}"
{{- end }}
)
{{ end }}
// {{ .StructName }} {{ .Table.Comment }}
type {{ .StructName }} struct {
{{- range .ColumnInfos }}
	{{ .FieldName }} {{ .GoType }}` + " `" + `{{ if $.Tags.JSON }}json:"{{ .JSONTag }}"{{ end }}{{ if $.Tags.GORM }} gorm:"{{ .GORMTag }}"{{ end }}{{ if $.Tags.Validate }} validate:"{{ .ValidateTag }}"{{ end }}` + "`" + `{{ if .Comment }} // {{ .Comment }}{{ end }}
{{- end }}
}

// TableName 返回表名
func ({{ .StructName }}) TableName() string {
	return "{{ .Table.Name }}"
}
`

// daoTemplate DAO 生成模板
const daoTemplate = `{{ .Header }}

package {{ .PackageName }}

import (
	"context"
	"database/sql"
	"fmt"
)

// {{ .StructName }}DAO {{ .Table.Comment }} 数据访问对象
type {{ .StructName }}DAO struct {
	db *sql.DB
}

// New{{ .StructName }}DAO 创建 {{ .StructName }}DAO
func New{{ .StructName }}DAO(db *sql.DB) *{{ .StructName }}DAO {
	return &{{ .StructName }}DAO{db: db}
}

// FindByID 根据 ID 查询
func (d *{{ .StructName }}DAO) FindByID(ctx context.Context, id int64) (*{{ .StructName }}, error) {
	query := ` + "`" + `SELECT {{ range $i, $col := .ColumnInfos }}{{ if $i }}, {{ end }}{{ $col.Name }}{{ end }} FROM {{ .Table.Name }} WHERE id = ?{{ if $.SoftDelete.Enabled }} AND {{ $.SoftDelete.Field }} IS NULL{{ end }}` + "`" + `
	
	row := d.db.QueryRowContext(ctx, query, id)
	
	var m {{ .StructName }}
	err := row.Scan(
{{- range $i, $col := .ColumnInfos }}
		&m.{{ $col.FieldName }},
{{- end }}
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find {{ .StructName }} by id: %w", err)
	}
	
	return &m, nil
}

// Create 创建记录
func (d *{{ .StructName }}DAO) Create(ctx context.Context, m *{{ .StructName }}) error {
	query := ` + "`" + `INSERT INTO {{ .Table.Name }} ({{ range $i, $col := .InsertColumns }}{{ if $i }}, {{ end }}{{ $col.Name }}{{ end }}) 
	            VALUES ({{ range $i, $col := .InsertColumns }}{{ if $i }}, {{ end }}?{{ end }})` + "`" + `
	
	result, err := d.db.ExecContext(ctx, query,
{{- range $i, $col := .InsertColumns }}
		m.{{ $col.FieldName }},
{{- end }}
	)
	if err != nil {
		return fmt.Errorf("failed to create {{ .StructName }}: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	
	m.ID = id
	return nil
}

// Update 更新记录
func (d *{{ .StructName }}DAO) Update(ctx context.Context, m *{{ .StructName }}) error {
	query := ` + "`" + `UPDATE {{ .Table.Name }} 
	            SET {{ range $i, $col := .UpdateColumns }}{{ if $i }}, {{ end }}{{ $col.Name }} = ?{{ end }}
	            WHERE id = ?{{ if $.SoftDelete.Enabled }} AND {{ $.SoftDelete.Field }} IS NULL{{ end }}` + "`" + `
	
	_, err := d.db.ExecContext(ctx, query,
{{- range $i, $col := .UpdateColumns }}
		m.{{ $col.FieldName }},
{{- end }}
		m.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update {{ .StructName }}: %w", err)
	}
	
	return nil
}
{{ if $.SoftDelete.Enabled }}
// Delete 软删除记录
func (d *{{ .StructName }}DAO) Delete(ctx context.Context, id int64) error {
	query := ` + "`" + `UPDATE {{ .Table.Name }} SET {{ $.SoftDelete.Field }} = CURRENT_TIMESTAMP WHERE id = ? AND {{ $.SoftDelete.Field }} IS NULL` + "`" + `
	
	_, err := d.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete {{ .StructName }}: %w", err)
	}
	
	return nil
}
{{ else }}
// Delete 删除记录
func (d *{{ .StructName }}DAO) Delete(ctx context.Context, id int64) error {
	query := ` + "`" + `DELETE FROM {{ .Table.Name }} WHERE id = ?` + "`" + `
	
	_, err := d.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete {{ .StructName }}: %w", err)
	}
	
	return nil
}
{{ end }}
`
