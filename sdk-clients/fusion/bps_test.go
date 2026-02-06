package fusion

import (
	"math/big"
	"testing"

	"github.com/1inch/1inch-sdk-go/common/fusionorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDefaultBase(t *testing.T) {
	base := fusionorder.GetDefaultBase()
	require.NotNil(t, base)
	assert.Equal(t, int64(1), base.Int64())
}

func TestNewBps(t *testing.T) {
	tests := []struct {
		name        string
		value       *big.Int
		expectError bool
	}{
		{
			name:        "Valid - zero",
			value:       big.NewInt(0),
			expectError: false,
		},
		{
			name:        "Valid - 100 bps (1%)",
			value:       big.NewInt(100),
			expectError: false,
		},
		{
			name:        "Valid - 5000 bps (50%)",
			value:       big.NewInt(5000),
			expectError: false,
		},
		{
			name:        "Valid - 10000 bps (100%)",
			value:       big.NewInt(10000),
			expectError: false,
		},
		{
			name:        "Invalid - negative",
			value:       big.NewInt(-1),
			expectError: true,
		},
		{
			name:        "Invalid - exceeds 10000",
			value:       big.NewInt(10001),
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bps, err := fusionorder.NewBps(tc.value)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, bps)
			} else {
				require.NoError(t, err)
				require.NotNil(t, bps)
				assert.Equal(t, tc.value.String(), bps.String())
			}
		})
	}
}

func TestFromPercent(t *testing.T) {
	tests := []struct {
		name     string
		percent  float64
		base     *big.Int
		expected string
	}{
		{
			name:     "1% with base 1",
			percent:  1,
			base:     big.NewInt(1),
			expected: "100",
		},
		{
			name:     "0.5% with base 1",
			percent:  0.5,
			base:     big.NewInt(1),
			expected: "50",
		},
		{
			name:     "10% with base 1",
			percent:  10,
			base:     big.NewInt(1),
			expected: "1000",
		},
		{
			name:     "50% with base 1",
			percent:  50,
			base:     big.NewInt(1),
			expected: "5000",
		},
		{
			name:     "100% with base 1",
			percent:  100,
			base:     big.NewInt(1),
			expected: "10000",
		},
		{
			name:     "1% with base 2",
			percent:  1,
			base:     big.NewInt(2),
			expected: "50",
		},
		{
			name:     "0% with base 1",
			percent:  0,
			base:     big.NewInt(1),
			expected: "0",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bps, err := fusionorder.FromPercent(tc.percent, tc.base)
			require.NoError(t, err)
			require.NotNil(t, bps)
			assert.Equal(t, tc.expected, bps.String())
		})
	}
}

func TestFromFraction(t *testing.T) {
	tests := []struct {
		name     string
		fraction float64
		base     *big.Int
		expected string
	}{
		{
			name:     "0.01 (1%) with base 1",
			fraction: 0.01,
			base:     big.NewInt(1),
			expected: "100",
		},
		{
			name:     "0.005 (0.5%) with base 1",
			fraction: 0.005,
			base:     big.NewInt(1),
			expected: "50",
		},
		{
			name:     "0.1 (10%) with base 1",
			fraction: 0.1,
			base:     big.NewInt(1),
			expected: "1000",
		},
		{
			name:     "0.5 (50%) with base 1",
			fraction: 0.5,
			base:     big.NewInt(1),
			expected: "5000",
		},
		{
			name:     "1.0 (100%) with base 1",
			fraction: 1.0,
			base:     big.NewInt(1),
			expected: "10000",
		},
		{
			name:     "0 with base 1",
			fraction: 0,
			base:     big.NewInt(1),
			expected: "0",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bps, err := fusionorder.FromFraction(tc.fraction, tc.base)
			require.NoError(t, err)
			require.NotNil(t, bps)
			assert.Equal(t, tc.expected, bps.String())
		})
	}
}

func TestBps_Equal(t *testing.T) {
	tests := []struct {
		name     string
		bps1     *Bps
		bps2     *Bps
		expected bool
	}{
		{
			name:     "Equal - both zero",
			bps1:     fusionorder.MustNewBps(big.NewInt(0)),
			bps2:     fusionorder.MustNewBps(big.NewInt(0)),
			expected: true,
		},
		{
			name:     "Equal - same non-zero value",
			bps1:     fusionorder.MustNewBps(big.NewInt(100)),
			bps2:     fusionorder.MustNewBps(big.NewInt(100)),
			expected: true,
		},
		{
			name:     "Not equal - different values",
			bps1:     fusionorder.MustNewBps(big.NewInt(100)),
			bps2:     fusionorder.MustNewBps(big.NewInt(200)),
			expected: false,
		},
		{
			name:     "Not equal - zero vs non-zero",
			bps1:     fusionorder.MustNewBps(big.NewInt(0)),
			bps2:     fusionorder.MustNewBps(big.NewInt(100)),
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.bps1.Equal(tc.bps2)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestBps_IsZero(t *testing.T) {
	tests := []struct {
		name     string
		bps      *Bps
		expected bool
	}{
		{
			name:     "Is zero",
			bps:      fusionorder.MustNewBps(big.NewInt(0)),
			expected: true,
		},
		{
			name:     "Is not zero - small value",
			bps:      fusionorder.MustNewBps(big.NewInt(1)),
			expected: false,
		},
		{
			name:     "Is not zero - large value",
			bps:      fusionorder.MustNewBps(big.NewInt(10000)),
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.bps.IsZero()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestBps_ToPercent(t *testing.T) {
	tests := []struct {
		name     string
		bps      *Bps
		base     *big.Int
		expected float64
	}{
		{
			name:     "0 bps to percent",
			bps:      fusionorder.MustNewBps(big.NewInt(0)),
			base:     big.NewInt(1),
			expected: 0,
		},
		{
			name:     "100 bps (1%) to percent",
			bps:      fusionorder.MustNewBps(big.NewInt(100)),
			base:     big.NewInt(1),
			expected: 1,
		},
		{
			name:     "5000 bps (50%) to percent",
			bps:      fusionorder.MustNewBps(big.NewInt(5000)),
			base:     big.NewInt(1),
			expected: 50,
		},
		{
			name:     "10000 bps (100%) to percent",
			bps:      fusionorder.MustNewBps(big.NewInt(10000)),
			base:     big.NewInt(1),
			expected: 100,
		},
		{
			name:     "100 bps with base 2",
			bps:      fusionorder.MustNewBps(big.NewInt(100)),
			base:     big.NewInt(2),
			expected: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.bps.ToPercent(tc.base)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestBps_ToFraction(t *testing.T) {
	tests := []struct {
		name     string
		bps      *Bps
		base     *big.Int
		expected *big.Int
	}{
		{
			name:     "0 bps to fraction",
			bps:      fusionorder.MustNewBps(big.NewInt(0)),
			base:     big.NewInt(1),
			expected: big.NewInt(0),
		},
		{
			name:     "10000 bps (100%) to fraction",
			bps:      fusionorder.MustNewBps(big.NewInt(10000)),
			base:     big.NewInt(1),
			expected: big.NewInt(1),
		},
		{
			name:     "5000 bps (50%) to fraction with base 2",
			bps:      fusionorder.MustNewBps(big.NewInt(5000)),
			base:     big.NewInt(2),
			expected: big.NewInt(1), // 5000 * 2 / 10000 = 1
		},
		{
			name:     "100 bps (1%) to fraction with large base",
			bps:      fusionorder.MustNewBps(big.NewInt(100)),
			base:     big.NewInt(10000),
			expected: big.NewInt(100), // 100 * 10000 / 10000 = 100
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.bps.ToFraction(tc.base)
			assert.Equal(t, 0, tc.expected.Cmp(result))
		})
	}
}

func TestBps_String(t *testing.T) {
	tests := []struct {
		name     string
		bps      *Bps
		expected string
	}{
		{
			name:     "Zero",
			bps:      fusionorder.MustNewBps(big.NewInt(0)),
			expected: "0",
		},
		{
			name:     "100 bps",
			bps:      fusionorder.MustNewBps(big.NewInt(100)),
			expected: "100",
		},
		{
			name:     "10000 bps",
			bps:      fusionorder.MustNewBps(big.NewInt(10000)),
			expected: "10000",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.bps.String()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestBpsZero(t *testing.T) {
	require.NotNil(t, fusionorder.BpsZero)
	assert.True(t, fusionorder.BpsZero.IsZero())
	assert.Equal(t, "0", fusionorder.BpsZero.String())
}
