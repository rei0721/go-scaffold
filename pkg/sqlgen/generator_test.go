package sqlgen

import (
	"strings"
	"testing"

	"github.com/rei0721/rei0721/internal/models"
)

func TestPostgresDialect(t *testing.T) {
	dialect := NewPostgresDialect()
	generator := New(dialect)

	result, err := generator.GenerateSQL(models.User{})
	if err != nil {
		t.Fatalf("生成 SQL 失败: %v", err)
	}

	// 验证表名
	if result.TableName != "users" {
		t.Errorf("期望表名 'users', 得到 '%s'", result.TableName)
	}

	// 验证建表 SQL 包含必要的字段
	createSQL := result.CreateTable
	expectedFields := []string{"id", "username", "email", "status", "created_at", "updated_at", "deleted_at"}
	for _, field := range expectedFields {
		if !strings.Contains(createSQL, field) {
			t.Errorf("建表 SQL 缺少字段: %s", field)
		}
	}

	// 验证索引
	if !strings.Contains(createSQL, "UNIQUE INDEX") {
		t.Error("建表 SQL 缺少唯一索引")
	}

	// 验证插入 SQL
	if !strings.Contains(result.Insert, "INSERT INTO users") {
		t.Error("插入 SQL 格式不正确")
	}

	// 验证查询 SQL
	if !strings.Contains(result.Select, "SELECT") {
		t.Error("查询 SQL 格式不正确")
	}

	// 验证更新 SQL
	if !strings.Contains(result.Update, "UPDATE users SET") {
		t.Error("更新 SQL 格式不正确")
	}

	// 验证删除 SQL
	if !strings.Contains(result.Delete, "UPDATE users SET deleted_at") {
		t.Error("删除 SQL 应该是软删除")
	}
}

func TestMySQLDialect(t *testing.T) {
	dialect := NewMySQLDialect()
	generator := New(dialect)

	result, err := generator.GenerateSQL(models.User{})
	if err != nil {
		t.Fatalf("生成 SQL 失败: %v", err)
	}

	// 验证 MySQL 特有的语法
	createSQL := result.CreateTable
	if !strings.Contains(createSQL, "ENGINE=InnoDB") {
		t.Error("MySQL 建表 SQL 应该包含 ENGINE=InnoDB")
	}

	if !strings.Contains(createSQL, "utf8mb4") {
		t.Error("MySQL 建表 SQL 应该包含 utf8mb4 字符集")
	}

	// 验证反引号
	if !strings.Contains(createSQL, "`users`") {
		t.Error("MySQL 建表 SQL 应该使用反引号")
	}
}

func TestSQLiteDialect(t *testing.T) {
	dialect := NewSQLiteDialect()
	generator := New(dialect)

	result, err := generator.GenerateSQL(models.User{})
	if err != nil {
		t.Fatalf("生成 SQL 失败: %v", err)
	}

	// 验证 SQLite 特有的语法
	createSQL := result.CreateTable
	if !strings.Contains(createSQL, "AUTOINCREMENT") {
		t.Error("SQLite 建表 SQL 应该包含 AUTOINCREMENT")
	}

	// 验证数据类型
	if !strings.Contains(createSQL, "INTEGER") {
		t.Error("SQLite 建表 SQL 应该使用 INTEGER 类型")
	}
}

func TestToSnakeCase(t *testing.T) {
	generator := &Generator{}

	testCases := []struct {
		input    string
		expected string
	}{
		{"ID", "id"},
		{"UserID", "user_id"},
		{"CreatedAt", "created_at"},
		{"UpdatedAt", "updated_at"},
		{"DeletedAt", "deleted_at"},
		{"Username", "username"},
		{"Email", "email"},
		{"Status", "status"},
		{"URL", "url"},
		{"HTTP", "http"},
		{"API", "api"},
	}

	for _, tc := range testCases {
		result := generator.toSnakeCase(tc.input)
		if result != tc.expected {
			t.Errorf("toSnakeCase(%s) = %s, 期望 %s", tc.input, result, tc.expected)
		}
	}
}

func TestFileGenerator(t *testing.T) {
	dialect := NewPostgresDialect()
	generator := New(dialect)
	fileGenerator := NewFileGenerator(generator, "./test_output")

	// 测试生成选项
	options := &GenerateOptions{
		OutputDir:       "./test_output",
		SeparateFiles:   false,
		GenerateSummary: true,
		IncludeComments: true,
	}

	models := []interface{}{
		models.User{},
	}

	err := fileGenerator.GenerateWithOptions(options, models...)
	if err != nil {
		t.Fatalf("生成文件失败: %v", err)
	}

	// 清理测试文件
	// 注意: 在实际测试中，你可能想要验证文件内容而不是直接删除
	// os.RemoveAll("./test_output")
}

// BenchmarkGenerateSQL 性能测试
func BenchmarkGenerateSQL(b *testing.B) {
	dialect := NewPostgresDialect()
	generator := New(dialect)
	user := models.User{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := generator.GenerateSQL(user)
		if err != nil {
			b.Fatal(err)
		}
	}
}
