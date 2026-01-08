package repository

import (
	"testing"
)

func TestGenerateCacheKey(t *testing.T) {
	repo := &cacheAdminAreaRepository{}

	tests := []struct {
		name     string
		prefix   string
		parts    []interface{}
		expected string
	}{
		{
			name:     "nil tolerance",
			prefix:   "admin_area",
			parts:    []interface{}{int32(1), 123, (*float64)(nil)},
			expected: "admin_area:1:123:<nil>",
		},
		{
			name:     "zero tolerance",
			prefix:   "admin_area",
			parts:    []interface{}{int32(1), 123, floatPtr(0.0)},
			expected: "admin_area:1:123:0.0000000000",
		},
		{
			name:     "non-zero tolerance",
			prefix:   "admin_area",
			parts:    []interface{}{int32(1), 123, floatPtr(0.001)},
			expected: "admin_area:1:123:0.0010000000",
		},
		{
			name:     "same tolerance value different pointers",
			prefix:   "admin_area",
			parts:    []interface{}{int32(1), 123, floatPtr(0.001)},
			expected: "admin_area:1:123:0.0010000000",
		},
		{
			name:     "list with nil tolerance",
			prefix:   "admin_area:list",
			parts:    []interface{}{int32(1), (*float64)(nil)},
			expected: "admin_area:list:1:<nil>",
		},
		{
			name:     "code with tolerance",
			prefix:   "admin_area:code",
			parts:    []interface{}{int32(1), "TH", floatPtr(0.001)},
			expected: "admin_area:code:1:TH:0.0010000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.generateCacheKey(tt.prefix, tt.parts...)
			if result != tt.expected {
				t.Errorf("generateCacheKey() = %v, want %v", result, tt.expected)
			}
		})
	}

	// Test that same tolerance value generates same cache key
	t.Run("consistency check", func(t *testing.T) {
		key1 := repo.generateCacheKey("admin_area", int32(1), 123, floatPtr(0.001))
		key2 := repo.generateCacheKey("admin_area", int32(1), 123, floatPtr(0.001))
		if key1 != key2 {
			t.Errorf("Same tolerance value should generate same cache key: key1=%v, key2=%v", key1, key2)
		}
	})
}

func floatPtr(f float64) *float64 {
	return &f
}
