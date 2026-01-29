package hexadecimal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsHexBytes_KnownValues verifies hex validation against known values
func TestIsHexBytes_KnownValues(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valid hex strings
		{name: "Empty with 0x prefix", input: "0x", expected: true},
		{name: "Simple hex", input: "0x1234", expected: true},
		{name: "All lowercase", input: "0xabcdef", expected: true},
		{name: "All uppercase", input: "0xABCDEF", expected: true},
		{name: "Mixed case", input: "0xAbCdEf", expected: true},
		{name: "Full address", input: "0x1234567890123456789012345678901234567890", expected: true},
		{name: "Without 0x prefix", input: "1234", expected: true},
		{name: "Empty string", input: "", expected: true},

		// Invalid hex strings
		{name: "Invalid char g", input: "0x123g", expected: false},
		{name: "Invalid char z", input: "0xz123", expected: false},
		{name: "Spaces", input: "0x12 34", expected: false},
		{name: "Special chars", input: "0x12@34", expected: false},
		{name: "Odd length without prefix", input: "123", expected: false},
		{name: "Odd length with prefix", input: "0x123", expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsHexBytes(tc.input)
			assert.Equal(t, tc.expected, result, "IsHexBytes(%q) mismatch", tc.input)
		})
	}
}

// TestTrim0x_KnownValues verifies 0x prefix trimming
func TestTrim0x_KnownValues(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "With 0x prefix", input: "0x1234", expected: "1234"},
		{name: "Without prefix", input: "1234", expected: "1234"},
		{name: "Empty with prefix", input: "0x", expected: ""},
		{name: "Empty string", input: "", expected: ""},
		{name: "Just 0x", input: "0x", expected: ""},
		{name: "Multiple 0x", input: "0x0x1234", expected: "0x1234"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Trim0x(tc.input)
			assert.Equal(t, tc.expected, result, "Trim0x(%q) mismatch", tc.input)
		})
	}
}
