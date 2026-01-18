package yaml2go

import (
	"strings"
	"testing"
)

func TestConverter_Convert_BasicTypes(t *testing.T) {
	yamlStr := `
name: "test"
age: 25
price: 19.99
enabled: true
`

	converter := New(nil)
	code, err := converter.Convert(yamlStr)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	// 验证生成的代码包含预期内容
	expected := []string{
		"package main",
		"type Config struct",
		"Name string",
		"Age int64",
		"Price float64",
		"Enabled bool",
		`json:"name"`,
		`yaml:"name"`,
		`mapstructure:"name"`,
	}

	for _, exp := range expected {
		if !strings.Contains(code, exp) {
			t.Errorf("Generated code missing expected content: %s\nGot:\n%s", exp, code)
		}
	}
}

func TestConverter_Convert_NestedStruct(t *testing.T) {
	yamlStr := `
database:
  host: localhost
  port: 5432
server:
  port: 8080
`

	converter := New(&Config{
		PackageName: "config",
		StructName:  "AppConfig",
	})

	code, err := converter.Convert(yamlStr)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	expected := []string{
		"package config",
		"type AppConfig struct",
		"Database struct",
		"Server struct",
		"Host string",
		"Port int64",
	}

	for _, exp := range expected {
		if !strings.Contains(code, exp) {
			t.Errorf("Generated code missing expected content: %s", exp)
		}
	}
}

func TestConverter_Convert_Array(t *testing.T) {
	yamlStr := `
tags:
  - tag1
  - tag2
  - tag3
numbers:
  - 1
  - 2
  - 3
`

	converter := New(nil)
	code, err := converter.Convert(yamlStr)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	expected := []string{
		"Tags []string",
		"Numbers []int64",
	}

	for _, exp := range expected {
		if !strings.Contains(code, exp) {
			t.Errorf("Generated code missing expected content: %s\nGot:\n%s", exp, code)
		}
	}
}

func TestConverter_Convert_ComplexNested(t *testing.T) {
	yamlStr := `
users:
  - name: alice
    age: 30
    roles:
      - admin
      - user
  - name: bob
    age: 25
    roles:
      - user
`

	converter := New(nil)
	code, err := converter.Convert(yamlStr)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	expected := []string{
		"Users []struct",
		"Name string",
		"Age int64",
		"Roles []string",
	}

	for _, exp := range expected {
		if !strings.Contains(code, exp) {
			t.Errorf("Generated code missing expected content: %s", exp)
		}
	}
}

func TestConverter_Convert_WithOmitEmpty(t *testing.T) {
	yamlStr := `
field1: value1
field2: value2
`

	converter := New(&Config{
		OmitEmpty: true,
	})

	code, err := converter.Convert(yamlStr)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if !strings.Contains(code, ",omitempty") {
		t.Error("Generated code should contain 'omitempty' option")
	}
}

func TestConverter_Convert_WithPointers(t *testing.T) {
	yamlStr := `
field: value
`

	converter := New(&Config{
		UsePointer: true,
	})

	code, err := converter.Convert(yamlStr)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if !strings.Contains(code, "*string") {
		t.Error("Generated code should use pointer types")
	}
}

func TestConverter_Convert_CustomTags(t *testing.T) {
	yamlStr := `
field: value
`

	converter := New(&Config{
		Tags: []string{"json", "xml"},
	})

	code, err := converter.Convert(yamlStr)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if !strings.Contains(code, `json:"field"`) {
		t.Error("Generated code should contain json tag")
	}

	if !strings.Contains(code, `xml:"field"`) {
		t.Error("Generated code should contain xml tag")
	}

	if strings.Contains(code, `yaml:"field"`) {
		t.Error("Generated code should not contain yaml tag (not in custom tags)")
	}
}

func TestConverter_Convert_EmptyInput(t *testing.T) {
	converter := New(nil)
	_, err := converter.Convert("")
	if err != ErrEmptyInput {
		t.Errorf("Expected ErrEmptyInput, got %v", err)
	}
}

func TestConverter_Convert_InvalidYAML(t *testing.T) {
	invalidYaml := `
invalid: [
  unclosed
`

	converter := New(nil)
	_, err := converter.Convert(invalidYaml)
	if err == nil || !strings.Contains(err.Error(), "invalid YAML") {
		t.Errorf("Expected invalid YAML error, got %v", err)
	}
}

func TestConverter_SetConfig(t *testing.T) {
	converter := New(nil)

	newConfig := &Config{
		PackageName: "mypackage",
		StructName:  "MyStruct",
	}

	err := converter.SetConfig(newConfig)
	if err != nil {
		t.Fatalf("SetConfig failed: %v", err)
	}

	yamlStr := `field: value`
	code, err := converter.Convert(yamlStr)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if !strings.Contains(code, "package mypackage") {
		t.Error("Config change did not take effect")
	}

	if !strings.Contains(code, "type MyStruct struct") {
		t.Error("Config change did not take effect")
	}
}

func TestConverter_SetConfig_Invalid(t *testing.T) {
	converter := New(nil)

	invalidConfig := &Config{
		IndentStyle: "invalid",
	}

	err := converter.SetConfig(invalidConfig)
	if err == nil {
		t.Error("Expected error for invalid config")
	}
}

func TestFieldType_String(t *testing.T) {
	tests := []struct {
		fieldType FieldType
		expected  string
	}{
		{TypeString, "string"},
		{TypeInt, "int64"},
		{TypeFloat, "float64"},
		{TypeBool, "bool"},
		{TypeSlice, "slice"},
		{TypeStruct, "struct"},
		{TypeMap, "map[string]interface{}"},
		{TypeInterface, "interface{}"},
		{TypeUnknown, "unknown"},
	}

	for _, tt := range tests {
		result := tt.fieldType.String()
		if result != tt.expected {
			t.Errorf("FieldType.String() = %s, expected %s", result, tt.expected)
		}
	}
}

func TestSanitizeFieldName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"my_field", "MyField"},
		{"type", "FieldType"},
		{"interface", "FieldInterface"},
		{"normal_name", "NormalName"},
	}

	for _, tt := range tests {
		result := sanitizeFieldName(tt.input)
		if result != tt.expected {
			t.Errorf("sanitizeFieldName(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestBuildTags(t *testing.T) {
	tags := map[string]string{
		"json": "field_name",
		"yaml": "field_name",
	}

	result := buildTags(tags, false)
	expected := "`json:\"field_name\" yaml:\"field_name\"`"

	if result != expected {
		t.Errorf("buildTags() = %s, expected %s", result, expected)
	}

	// 测试 omitempty
	result = buildTags(tags, true)
	expected = "`json:\"field_name,omitempty\" yaml:\"field_name,omitempty\"`"

	if result != expected {
		t.Errorf("buildTags() with omitempty = %s, expected %s", result, expected)
	}
}

// Benchmark tests
func BenchmarkConverter_Convert_Simple(b *testing.B) {
	yamlStr := `
name: test
age: 25
enabled: true
`
	converter := New(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = converter.Convert(yamlStr)
	}
}

func BenchmarkConverter_Convert_Complex(b *testing.B) {
	yamlStr := `
database:
  host: localhost
  port: 5432
  credentials:
    username: admin
    password: secret
server:
  host: 0.0.0.0
  port: 8080
  timeout: 30
  routes:
    - path: /api
      method: GET
    - path: /health
      method: GET
`
	converter := New(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = converter.Convert(yamlStr)
	}
}
