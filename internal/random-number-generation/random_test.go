package random_number_generation

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBigIntMax(t *testing.T) {
	tests := []struct {
		name string
		max  *big.Int
	}{
		{
			name: "Small max value",
			max:  big.NewInt(100),
		},
		{
			name: "Large max value",
			max:  new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), // 10^18
		},
		{
			name: "Max value 1",
			max:  big.NewInt(1),
		},
		{
			name: "Max value 2",
			max:  big.NewInt(2),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := BigIntMax(tc.max)
			require.NoError(t, err)
			require.NotNil(t, result)

			// Result should be >= 0
			assert.GreaterOrEqual(t, result.Cmp(big.NewInt(0)), 0)

			// Result should be < max
			assert.Less(t, result.Cmp(tc.max), 0)
		})
	}
}

func TestBigIntMax_RandomnessDistribution(t *testing.T) {
	max := big.NewInt(10)
	counts := make(map[int64]int)

	// Generate many random numbers
	iterations := 1000
	for i := 0; i < iterations; i++ {
		result, err := BigIntMax(max)
		require.NoError(t, err)
		counts[result.Int64()]++
	}

	// All values from 0 to 9 should appear at least once with high probability
	// This is a weak test but ensures basic distribution
	for i := int64(0); i < 10; i++ {
		assert.Greater(t, counts[i], 0, "Value %d should appear at least once in %d iterations", i, iterations)
	}
}

func TestBigIntMaxFunc_CanBeMocked(t *testing.T) {
	// Save the original function
	originalFunc := BigIntMaxFunc
	defer func() { BigIntMaxFunc = originalFunc }()

	// Mock the function
	mockedValue := big.NewInt(42)
	BigIntMaxFunc = func(max *big.Int) (*big.Int, error) {
		return mockedValue, nil
	}

	result, err := BigIntMaxFunc(big.NewInt(100))
	require.NoError(t, err)
	assert.Equal(t, mockedValue, result)
}

func TestBigIntMax_VeryLargeMax(t *testing.T) {
	// Test with uint256 max (2^256 - 1)
	max := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

	result, err := BigIntMax(max)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Result should be >= 0
	assert.GreaterOrEqual(t, result.Cmp(big.NewInt(0)), 0)

	// Result should be < max
	assert.Less(t, result.Cmp(max), 0)
}
