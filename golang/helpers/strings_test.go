package helpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsScientificNotation(t *testing.T) {
	testCases := []struct {
		description string
		input       string
		expected    bool
	}{
		{
			description: "Valid scientific notation with positive exponent",
			input:       "1e+18",
			expected:    true,
		},
		{
			description: "Valid scientific notation with negative exponent",
			input:       "2.3e-4",
			expected:    true,
		},
		{
			description: "Invalid format - not scientific notation",
			input:       "not_scientific",
			expected:    false,
		},
		{
			description: "Invalid format - regular decimal number",
			input:       "3.14",
			expected:    false,
		},
		{
			description: "Valid scientific notation without sign in exponent",
			input:       "1e18",
			expected:    true,
		},
		{
			description: "Valid scientific notation with negative base and positive exponent",
			input:       "-3.5E+10",
			expected:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result := IsScientificNotation(tc.input)
			assert.Equal(t, tc.expected, result, fmt.Sprintf("TestIsScientificNotation failed for input: %s", tc.input))
		})
	}
}

func TestExpandScientificNotation(t *testing.T) {
	testCases := []struct {
		description    string
		input          string
		expectedOutput string
		expectedError  bool
	}{
		{
			description:    "Valid scientific notation with positive exponent",
			input:          "1e+18",
			expectedOutput: "1000000000000000000",
			expectedError:  false,
		},
		{
			description:    "Valid scientific notation with negative exponent",
			input:          "2.3e-4",
			expectedOutput: "0.00023",
			expectedError:  false,
		},
		{
			description:   "Invalid format - not a number",
			input:         "not_a_number",
			expectedError: true,
		},
		{
			description:    "Regular number without scientific notation",
			input:          "12345",
			expectedOutput: "12345",
			expectedError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			output, err := ExpandScientificNotation(tc.input)

			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}
