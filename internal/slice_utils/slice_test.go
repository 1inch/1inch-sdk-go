package slice_utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains_String(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		slice    []string
		expected bool
	}{
		{
			name:     "Value exists in slice",
			value:    "apple",
			slice:    []string{"apple", "banana", "cherry"},
			expected: true,
		},
		{
			name:     "Value does not exist in slice",
			value:    "grape",
			slice:    []string{"apple", "banana", "cherry"},
			expected: false,
		},
		{
			name:     "Empty slice",
			value:    "apple",
			slice:    []string{},
			expected: false,
		},
		{
			name:     "Single element slice - match",
			value:    "apple",
			slice:    []string{"apple"},
			expected: true,
		},
		{
			name:     "Single element slice - no match",
			value:    "banana",
			slice:    []string{"apple"},
			expected: false,
		},
		{
			name:     "Empty string in slice",
			value:    "",
			slice:    []string{"", "apple"},
			expected: true,
		},
		{
			name:     "Case sensitive - no match",
			value:    "Apple",
			slice:    []string{"apple", "banana"},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Contains(tc.value, tc.slice)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestContains_Int(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		slice    []int
		expected bool
	}{
		{
			name:     "Value exists",
			value:    5,
			slice:    []int{1, 3, 5, 7, 9},
			expected: true,
		},
		{
			name:     "Value does not exist",
			value:    4,
			slice:    []int{1, 3, 5, 7, 9},
			expected: false,
		},
		{
			name:     "Empty slice",
			value:    5,
			slice:    []int{},
			expected: false,
		},
		{
			name:     "Negative number exists",
			value:    -5,
			slice:    []int{-10, -5, 0, 5, 10},
			expected: true,
		},
		{
			name:     "Zero exists",
			value:    0,
			slice:    []int{-1, 0, 1},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Contains(tc.value, tc.slice)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestContains_Float64(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		slice    []float64
		expected bool
	}{
		{
			name:     "Value exists",
			value:    3.14,
			slice:    []float64{1.0, 2.0, 3.14, 4.0},
			expected: true,
		},
		{
			name:     "Value does not exist",
			value:    2.71,
			slice:    []float64{1.0, 2.0, 3.14, 4.0},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Contains(tc.value, tc.slice)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// Custom type for testing
type CustomType struct {
	ID   int
	Name string
}

func TestContains_CustomType(t *testing.T) {
	a := CustomType{ID: 1, Name: "Alice"}
	b := CustomType{ID: 2, Name: "Bob"}
	c := CustomType{ID: 3, Name: "Charlie"}

	slice := []CustomType{a, b}

	assert.True(t, Contains(a, slice))
	assert.True(t, Contains(b, slice))
	assert.False(t, Contains(c, slice))
}
