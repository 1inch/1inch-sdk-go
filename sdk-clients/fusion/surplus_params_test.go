package fusion

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUint256Max(t *testing.T) {
	// uint256 max = 2^256 - 1
	expected := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
	assert.Equal(t, 0, Uint256Max.Cmp(expected))

	// Should be 78 digits
	assert.Equal(t, 78, len(Uint256Max.String()))
}

func TestSurplusParamsNoFee(t *testing.T) {
	require.NotNil(t, SurplusParamsNoFee)
	assert.Equal(t, 0, SurplusParamsNoFee.EstimatedTakerAmount.Cmp(Uint256Max))
	assert.True(t, SurplusParamsNoFee.ProtocolFee.IsZero())
}

func TestNewSurplusParams(t *testing.T) {
	tests := []struct {
		name                 string
		estimatedTakerAmount *big.Int
		protocolFee          *Bps
		expectError          bool
		errorMsg             string
	}{
		{
			name:                 "Valid - zero protocol fee",
			estimatedTakerAmount: big.NewInt(1000000000000000000),
			protocolFee:          BpsZero,
			expectError:          false,
		},
		{
			name:                 "Valid - 1% protocol fee (100 bps)",
			estimatedTakerAmount: big.NewInt(1000000000000000000),
			protocolFee:          NewBps(big.NewInt(100)),
			expectError:          false,
		},
		{
			name:                 "Valid - 5% protocol fee (500 bps)",
			estimatedTakerAmount: big.NewInt(1000000000000000000),
			protocolFee:          NewBps(big.NewInt(500)),
			expectError:          false,
		},
		{
			name:                 "Valid - max protocol fee (100% = 10000 bps)",
			estimatedTakerAmount: big.NewInt(1000000000000000000),
			protocolFee:          NewBps(big.NewInt(10000)),
			expectError:          false,
		},
		{
			name:                 "Valid - uint256 max estimated amount",
			estimatedTakerAmount: Uint256Max,
			protocolFee:          BpsZero,
			expectError:          false,
		},
		{
			name:                 "Invalid - 0.5% protocol fee (50 bps) - not whole percent",
			estimatedTakerAmount: big.NewInt(1000000000000000000),
			protocolFee:          NewBps(big.NewInt(50)),
			expectError:          true,
			errorMsg:             "only integer percent supported",
		},
		{
			name:                 "Invalid - 1.5% protocol fee (150 bps) - not whole percent",
			estimatedTakerAmount: big.NewInt(1000000000000000000),
			protocolFee:          NewBps(big.NewInt(150)),
			expectError:          true,
			errorMsg:             "only integer percent supported",
		},
		{
			name:                 "Invalid - 0.01% protocol fee (1 bps) - not whole percent",
			estimatedTakerAmount: big.NewInt(1000000000000000000),
			protocolFee:          NewBps(big.NewInt(1)),
			expectError:          true,
			errorMsg:             "only integer percent supported",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := NewSurplusParams(tc.estimatedTakerAmount, tc.protocolFee)
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, 0, tc.estimatedTakerAmount.Cmp(result.EstimatedTakerAmount))
				assert.True(t, tc.protocolFee.Equal(result.ProtocolFee))
			}
		})
	}
}

func TestNewSurplusParams_ImmutabilityOfInput(t *testing.T) {
	// Test that the input big.Int is copied, not referenced
	originalAmount := big.NewInt(1000)
	sp, err := NewSurplusParams(originalAmount, BpsZero)
	require.NoError(t, err)

	// Modify the original
	originalAmount.SetInt64(9999)

	// SurplusParams should still have the original value
	assert.Equal(t, int64(1000), sp.EstimatedTakerAmount.Int64())
}
