package fusionorder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrefix0x(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Already has 0x prefix",
			input:    "0x1234abcd",
			expected: "0x1234abcd",
		},
		{
			name:     "No prefix",
			input:    "1234abcd",
			expected: "0x1234abcd",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "0x",
		},
		{
			name:     "Just 0x",
			input:    "0x",
			expected: "0x",
		},
		{
			name:     "Uppercase 0X",
			input:    "0X1234",
			expected: "0x0X1234", // 0X is not treated as prefix
		},
		{
			name:     "Single character",
			input:    "a",
			expected: "0xa",
		},
		{
			name:     "0x in middle",
			input:    "abc0xdef",
			expected: "0xabc0xdef",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Prefix0x(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
