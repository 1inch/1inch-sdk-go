package hexadecimal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsHexBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid hex with 0x prefix",
			input:    "0x1234567890abcdef",
			expected: true,
		},
		{
			name:     "Valid hex without 0x prefix",
			input:    "1234567890abcdef",
			expected: true,
		},
		{
			name:     "Valid hex uppercase",
			input:    "0xABCDEF",
			expected: true,
		},
		{
			name:     "Valid hex mixed case",
			input:    "0xAbCdEf123456",
			expected: true,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: true, // empty hex decodes to empty bytes
		},
		{
			name:     "Just 0x prefix",
			input:    "0x",
			expected: true, // empty hex decodes to empty bytes
		},
		{
			name:     "Invalid hex - odd length",
			input:    "0x123",
			expected: false,
		},
		{
			name:     "Invalid hex - non-hex characters",
			input:    "0xGHIJKL",
			expected: false,
		},
		{
			name:     "Invalid hex - spaces",
			input:    "0x12 34",
			expected: false,
		},
		{
			name:     "Valid Ethereum address",
			input:    "0x6B175474E89094C44Da98b954EedeAC495271d0F",
			expected: true,
		},
		{
			name:     "Valid 32-byte hash",
			input:    "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsHexBytes(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestTrim0x(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "With 0x prefix",
			input:    "0x1234567890abcdef",
			expected: "1234567890abcdef",
		},
		{
			name:     "Without 0x prefix",
			input:    "1234567890abcdef",
			expected: "1234567890abcdef",
		},
		{
			name:     "Just 0x",
			input:    "0x",
			expected: "",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "0X uppercase prefix - should not trim",
			input:    "0X1234",
			expected: "0X1234",
		},
		{
			name:     "Multiple 0x prefixes - only first trimmed",
			input:    "0x0x1234",
			expected: "0x1234",
		},
		{
			name:     "Ethereum address",
			input:    "0x6B175474E89094C44Da98b954EedeAC495271d0F",
			expected: "6B175474E89094C44Da98b954EedeAC495271d0F",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Trim0x(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
