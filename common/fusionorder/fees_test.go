package fusionorder

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBpsToRatioFormat(t *testing.T) {
	tests := []struct {
		name     string
		bps      *big.Int
		expected *big.Int
	}{
		{
			name:     "Nil bps",
			bps:      nil,
			expected: big.NewInt(0),
		},
		{
			name:     "Zero bps",
			bps:      big.NewInt(0),
			expected: big.NewInt(0),
		},
		{
			name:     "100 bps (1%)",
			bps:      big.NewInt(100),
			expected: big.NewInt(1000), // 100 * 10 = 1000
		},
		{
			name:     "500 bps (5%)",
			bps:      big.NewInt(500),
			expected: big.NewInt(5000), // 500 * 10 = 5000
		},
		{
			name:     "10000 bps (100%)",
			bps:      big.NewInt(10000),
			expected: big.NewInt(100000), // 10000 * 10 = 100000
		},
		{
			name:     "1 bps (0.01%)",
			bps:      big.NewInt(1),
			expected: big.NewInt(10), // 1 * 10 = 10
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := BpsToRatioFormat(tc.bps)
			assert.Equal(t, 0, tc.expected.Cmp(result), "expected %s, got %s", tc.expected.String(), result.String())
		})
	}
}

func TestBpsToRatioFormatDoesNotMutateInput(t *testing.T) {
	original := big.NewInt(100)
	originalCopy := new(big.Int).Set(original)

	_ = BpsToRatioFormat(original)

	// Original should not be mutated
	assert.Equal(t, originalCopy, original)
}
