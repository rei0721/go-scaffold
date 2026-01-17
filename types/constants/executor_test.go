package constants

import (
	"testing"

	"github.com/rei0721/go-scaffold/pkg/executor"
)

// TestPoolNameConstants 验证池名常量的类型和值
func TestPoolNameConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant executor.PoolName
		expected string
	}{
		{"PoolHTTP", PoolHTTP, "http"},
		{"PoolDatabase", PoolDatabase, "database"},
		{"PoolCache", PoolCache, "cache"},
		{"PoolLogger", PoolLogger, "logger"},
		{"PoolBackground", PoolBackground, "background"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 验证常量类型
			var _ executor.PoolName = tt.constant

			// 验证常量值
			if string(tt.constant) != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

// TestPoolNameUniqueness 验证池名常量的唯一性
func TestPoolNameUniqueness(t *testing.T) {
	pools := []executor.PoolName{
		PoolHTTP,
		PoolDatabase,
		PoolCache,
		PoolLogger,
		PoolBackground,
	}

	seen := make(map[string]bool)
	for _, pool := range pools {
		poolStr := string(pool)
		if seen[poolStr] {
			t.Errorf("Duplicate pool name found: %s", poolStr)
		}
		seen[poolStr] = true
	}

	// 验证定义了5个池
	if len(pools) != 5 {
		t.Errorf("Expected 5 pools, got %d", len(pools))
	}
}
