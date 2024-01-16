package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimplifyValue(t *testing.T) {
	testCases := []struct {
		description    string
		input          string
		decimalPlaces  int
		expectedResult string
	}{
		{
			description:    "Integer greater than zero",
			input:          "100",
			decimalPlaces:  2,
			expectedResult: "1",
		},
		{
			description:    "Decimal",
			input:          "100",
			decimalPlaces:  3,
			expectedResult: "0.1",
		},
		{
			description:    "Number greater than one with decimals",
			input:          "101",
			decimalPlaces:  2,
			expectedResult: "1.01",
		},
		{
			description:    "Large number greater than one with decimals",
			input:          "10000000000001000010",
			decimalPlaces:  18,
			expectedResult: "10.00000000000100001",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result := SimplifyValue(tc.input, tc.decimalPlaces)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
