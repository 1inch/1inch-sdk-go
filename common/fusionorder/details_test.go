package fusionorder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNonceRequired(t *testing.T) {
	tests := []struct {
		name               string
		allowPartialFills  bool
		allowMultipleFills bool
		expected           bool
	}{
		{
			name:               "Both true - nonce not required",
			allowPartialFills:  true,
			allowMultipleFills: true,
			expected:           false,
		},
		{
			name:               "Partial false, multiple true - nonce required",
			allowPartialFills:  false,
			allowMultipleFills: true,
			expected:           true,
		},
		{
			name:               "Partial true, multiple false - nonce required",
			allowPartialFills:  true,
			allowMultipleFills: false,
			expected:           true,
		},
		{
			name:               "Both false - nonce required",
			allowPartialFills:  false,
			allowMultipleFills: false,
			expected:           true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsNonceRequired(tc.allowPartialFills, tc.allowMultipleFills)
			assert.Equal(t, tc.expected, result)
		})
	}
}
