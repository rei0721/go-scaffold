package sqlgen

import (
	"testing"
)

// TestToSnakeCase 测试蛇形命名转换
func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"UserName", "user_name"},
		{"userName", "user_name"},
		{"user_name", "user_name"},
		{"HTTPRequest", "http_request"},
		{"ID", "id"},
		{"", ""},
	}

	for _, tt := range tests {
		result := ToSnakeCase(tt.input)
		if result != tt.expected {
			t.Errorf("ToSnakeCase(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

// TestToCamelCase 测试驼峰命名转换
func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user_name", "userName"},
		{"USER_NAME", "userName"},
		{"user-name", "userName"},
		{"UserName", "username"},
		{"", ""},
	}

	for _, tt := range tests {
		result := ToCamelCase(tt.input)
		if result != tt.expected {
			t.Errorf("ToCamelCase(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

// TestToPascalCase 测试帕斯卡命名转换
func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user_name", "UserName"},
		{"USER_NAME", "UserName"},
		{"user-name", "UserName"},
		{"", ""},
	}

	for _, tt := range tests {
		result := ToPascalCase(tt.input)
		if result != tt.expected {
			t.Errorf("ToPascalCase(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

// TestToPlural 测试复数转换
func TestToPlural(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user", "users"},
		{"class", "classes"},
		{"box", "boxes"},
		{"company", "companies"},
		{"", ""},
	}

	for _, tt := range tests {
		result := ToPlural(tt.input)
		if result != tt.expected {
			t.Errorf("ToPlural(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

// TestDefaultConfig 测试默认配置
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.OutputDir != DefaultOutputDir {
		t.Errorf("DefaultConfig().OutputDir = %q, want %q", cfg.OutputDir, DefaultOutputDir)
	}

	if cfg.PackageName != DefaultPackageName {
		t.Errorf("DefaultConfig().PackageName = %q, want %q", cfg.PackageName, DefaultPackageName)
	}

	if !cfg.Target.Model {
		t.Error("DefaultConfig().Target.Model should be true")
	}

	if !cfg.Tags.JSON {
		t.Error("DefaultConfig().Tags.JSON should be true")
	}
}

// TestConfigValidate 测试配置验证
func TestConfigValidate(t *testing.T) {
	// 有效配置
	cfg := DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("Config.Validate() returned error for valid config: %v", err)
	}

	// 缺少 OutputDir
	cfg2 := DefaultConfig()
	cfg2.OutputDir = ""
	if err := cfg2.Validate(); err == nil {
		t.Error("Config.Validate() should return error when OutputDir is empty")
	}

	// 无效数据库类型
	cfg3 := DefaultConfig()
	cfg3.DatabaseType = "invalid"
	if err := cfg3.Validate(); err == nil {
		t.Error("Config.Validate() should return error for invalid database type")
	}
}
