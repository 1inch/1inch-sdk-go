package fusionorder

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDefaultBase(t *testing.T) {
	base := GetDefaultBase()
	assert.Equal(t, big.NewInt(1), base)
}

func TestNewBps(t *testing.T) {
	tests := []struct {
		name        string
		value       *big.Int
		shouldError bool
	}{
		{
			name:        "Valid - zero",
			value:       big.NewInt(0),
			shouldError: false,
		},
		{
			name:        "Valid - 100 bps (1%)",
			value:       big.NewInt(100),
			shouldError: false,
		},
		{
			name:        "Valid - 5000 bps (50%)",
			value:       big.NewInt(5000),
			shouldError: false,
		},
		{
			name:        "Valid - 10000 bps (100%)",
			value:       big.NewInt(10000),
			shouldError: false,
		},
		{
			name:        "Invalid - negative",
			value:       big.NewInt(-1),
			shouldError: true,
		},
		{
			name:        "Invalid - exceeds 10000",
			value:       big.NewInt(10001),
			shouldError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bps, err := NewBps(tc.value)
			if tc.shouldError {
				assert.Error(t, err)
				assert.Nil(t, bps)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.value, bps.Value())
			}
		})
	}
}

func TestMustNewBps(t *testing.T) {
	// Valid value should not panic
	assert.NotPanics(t, func() {
		bps := MustNewBps(big.NewInt(100))
		assert.Equal(t, big.NewInt(100), bps.Value())
	})

	// Invalid value should panic
	assert.Panics(t, func() {
		MustNewBps(big.NewInt(-1))
	})
}

func TestFromPercent(t *testing.T) {
	tests := []struct {
		name     string
		percent  float64
		base     *big.Int
		expected *big.Int
	}{
		{
			name:     "1% with base 1",
			percent:  1,
			base:     big.NewInt(1),
			expected: big.NewInt(100),
		},
		{
			name:     "0.5% with base 1",
			percent:  0.5,
			base:     big.NewInt(1),
			expected: big.NewInt(50),
		},
		{
			name:     "10% with base 1",
			percent:  10,
			base:     big.NewInt(1),
			expected: big.NewInt(1000),
		},
		{
			name:     "0% with base 1",
			percent:  0,
			base:     big.NewInt(1),
			expected: big.NewInt(0),
		},
		{
			name:     "1% with base 2",
			percent:  1,
			base:     big.NewInt(2),
			expected: big.NewInt(50),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bps, err := FromPercent(tc.percent, tc.base)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, bps.Value())
		})
	}
}

func TestFromFraction(t *testing.T) {
	tests := []struct {
		name     string
		fraction float64
		base     *big.Int
		expected *big.Int
	}{
		{
			name:     "0.01 (1%) with base 1",
			fraction: 0.01,
			base:     big.NewInt(1),
			expected: big.NewInt(100),
		},
		{
			name:     "0.005 (0.5%) with base 1",
			fraction: 0.005,
			base:     big.NewInt(1),
			expected: big.NewInt(50),
		},
		{
			name:     "0.1 (10%) with base 1",
			fraction: 0.1,
			base:     big.NewInt(1),
			expected: big.NewInt(1000),
		},
		{
			name:     "0 with base 1",
			fraction: 0,
			base:     big.NewInt(1),
			expected: big.NewInt(0),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bps, err := FromFraction(tc.fraction, tc.base)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, bps.Value())
		})
	}
}

func TestBpsEqual(t *testing.T) {
	bps1 := MustNewBps(big.NewInt(100))
	bps2 := MustNewBps(big.NewInt(100))
	bps3 := MustNewBps(big.NewInt(200))

	assert.True(t, bps1.Equal(bps2))
	assert.False(t, bps1.Equal(bps3))
}

func TestBpsIsZero(t *testing.T) {
	zero := MustNewBps(big.NewInt(0))
	nonZero := MustNewBps(big.NewInt(100))

	assert.True(t, zero.IsZero())
	assert.False(t, nonZero.IsZero())
	assert.True(t, BpsZero.IsZero())
}

func TestBpsToPercent(t *testing.T) {
	tests := []struct {
		name     string
		bps      *Bps
		base     *big.Int
		expected float64
	}{
		{
			name:     "100 bps = 1%",
			bps:      MustNewBps(big.NewInt(100)),
			base:     big.NewInt(1),
			expected: 1.0,
		},
		{
			name:     "5000 bps = 50%",
			bps:      MustNewBps(big.NewInt(5000)),
			base:     big.NewInt(1),
			expected: 50.0,
		},
		{
			name:     "0 bps = 0%",
			bps:      MustNewBps(big.NewInt(0)),
			base:     big.NewInt(1),
			expected: 0.0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.bps.ToPercent(tc.base)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestBpsToFraction(t *testing.T) {
	tests := []struct {
		name     string
		bps      *Bps
		base     *big.Int
		expected *big.Int
	}{
		{
			name:     "100 bps with base 100",
			bps:      MustNewBps(big.NewInt(100)),
			base:     big.NewInt(100),
			expected: big.NewInt(1), // 100 * 100 / 10000 = 1
		},
		{
			name:     "5000 bps with base 100",
			bps:      MustNewBps(big.NewInt(5000)),
			base:     big.NewInt(100),
			expected: big.NewInt(50), // 5000 * 100 / 10000 = 50
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.bps.ToFraction(tc.base)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestBpsString(t *testing.T) {
	bps := MustNewBps(big.NewInt(100))
	assert.Equal(t, "100", bps.String())
}

func TestBpsValue(t *testing.T) {
	original := big.NewInt(100)
	bps := MustNewBps(original)
	
	// Value should return a copy, not the original
	value := bps.Value()
	require.Equal(t, original, value)
	
	// Modifying the returned value should not affect the Bps
	value.SetInt64(999)
	assert.Equal(t, big.NewInt(100), bps.Value())
}
