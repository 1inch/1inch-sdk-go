package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsBigInt(t *testing.T) {
	testCases := []struct {
		description string
		input       string
		expectError bool
	}{
		{
			description: "Max big value",
			input:       "115792089237316195423570985008687907853269984665640564039457584007913129639935",
		},
		{
			description: "Max big value + 1",
			input:       "115792089237316195423570985008687907853269984665640564039457584007913129639936",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			_, err := BigIntFromString(tc.input)
			require.NoError(t, err)
		})
	}
}
