package fusion

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/internal/addresses"
)

func TestResolverFeeZero(t *testing.T) {
	require.NotNil(t, ResolverFeeZero)
	assert.Equal(t, addresses.ZeroAddress, ResolverFeeZero.Receiver)
	assert.True(t, ResolverFeeZero.Fee.IsZero())
	assert.True(t, ResolverFeeZero.WhitelistDiscount.IsZero())
}

func TestNewResolverFee(t *testing.T) {
	validAddress := "0x6B175474E89094C44Da98b954EedeAC495271d0F"

	tests := []struct {
		name              string
		receiver          string
		fee               *Bps
		whitelistDiscount *Bps
		expectError       bool
		errorMsg          string
	}{
		{
			name:              "Valid - zero fee with zero address",
			receiver:          addresses.ZeroAddress,
			fee:               BpsZero,
			whitelistDiscount: BpsZero,
			expectError:       false,
		},
		{
			name:              "Valid - zero fee with empty receiver",
			receiver:          "",
			fee:               BpsZero,
			whitelistDiscount: BpsZero,
			expectError:       false,
		},
		{
			name:              "Valid - non-zero fee with valid receiver",
			receiver:          validAddress,
			fee:               NewBps(big.NewInt(100)),
			whitelistDiscount: NewBps(big.NewInt(100)), // 1% discount
			expectError:       false,
		},
		{
			name:              "Valid - fee with zero whitelist discount",
			receiver:          validAddress,
			fee:               NewBps(big.NewInt(100)),
			whitelistDiscount: BpsZero,
			expectError:       false,
		},
		{
			name:              "Invalid - non-zero fee with zero address",
			receiver:          addresses.ZeroAddress,
			fee:               NewBps(big.NewInt(100)),
			whitelistDiscount: BpsZero,
			expectError:       true,
			errorMsg:          "fee must be zero if receiver is zero address",
		},
		{
			name:              "Invalid - non-zero fee with empty receiver",
			receiver:          "",
			fee:               NewBps(big.NewInt(100)),
			whitelistDiscount: BpsZero,
			expectError:       true,
			errorMsg:          "fee must be zero if receiver is zero address",
		},
		{
			name:              "Invalid - zero fee with non-zero receiver",
			receiver:          validAddress,
			fee:               BpsZero,
			whitelistDiscount: BpsZero,
			expectError:       true,
			errorMsg:          "receiver must be zero address if fee is zero",
		},
		{
			name:              "Invalid - zero fee with non-zero whitelist discount",
			receiver:          addresses.ZeroAddress,
			fee:               BpsZero,
			whitelistDiscount: NewBps(big.NewInt(100)),
			expectError:       true,
			errorMsg:          "whitelist discount must be zero if fee is zero",
		},
		{
			name:              "Invalid - whitelist discount not percent precision (50 bps = 0.5%)",
			receiver:          validAddress,
			fee:               NewBps(big.NewInt(100)),
			whitelistDiscount: NewBps(big.NewInt(50)), // 0.5% - not whole percent
			expectError:       true,
			errorMsg:          "whitelist discount must have percent precision",
		},
		{
			name:              "Valid - whitelist discount with percent precision (200 bps = 2%)",
			receiver:          validAddress,
			fee:               NewBps(big.NewInt(100)),
			whitelistDiscount: NewBps(big.NewInt(200)), // 2%
			expectError:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := NewResolverFee(tc.receiver, tc.fee, tc.whitelistDiscount)
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tc.receiver, result.Receiver)
				assert.True(t, tc.fee.Equal(result.Fee))
				assert.True(t, tc.whitelistDiscount.Equal(result.WhitelistDiscount))
			}
		})
	}
}

func TestResolverFee_String(t *testing.T) {
	validAddress := "0x6B175474E89094C44Da98b954EedeAC495271d0F"

	fee, err := NewResolverFee(validAddress, NewBps(big.NewInt(100)), NewBps(big.NewInt(200)))
	require.NoError(t, err)

	result := fee.String()
	assert.Contains(t, result, "ResolverFee")
	assert.Contains(t, result, validAddress)
	assert.Contains(t, result, "100")
	assert.Contains(t, result, "200")
}
