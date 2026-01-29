package fusion

import (
	"math/big"
	"testing"

	"github.com/1inch/1inch-sdk-go/common/fusionorder"
	"github.com/1inch/1inch-sdk-go/internal/addresses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegratorFeeZero(t *testing.T) {
	require.NotNil(t, IntegratorFeeZero)
	assert.Equal(t, addresses.ZeroAddress, IntegratorFeeZero.Integrator)
	assert.Equal(t, addresses.ZeroAddress, IntegratorFeeZero.Protocol)
	assert.True(t, IntegratorFeeZero.Fee.IsZero())
	assert.True(t, IntegratorFeeZero.Share.IsZero())
}

func TestNewIntegratorFee(t *testing.T) {
	validAddress1 := "0x6B175474E89094C44Da98b954EedeAC495271d0F"
	validAddress2 := "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"

	tests := []struct {
		name        string
		integrator  string
		protocol    string
		fee         *Bps
		share       *Bps
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid - zero fee with zero addresses",
			integrator:  addresses.ZeroAddress,
			protocol:    addresses.ZeroAddress,
			fee:         fusionorder.BpsZero,
			share:       fusionorder.BpsZero,
			expectError: false,
		},
		{
			name:        "Valid - non-zero fee with valid addresses",
			integrator:  validAddress1,
			protocol:    validAddress2,
			fee:         fusionorder.MustNewBps(big.NewInt(100)),
			share:       fusionorder.MustNewBps(big.NewInt(5000)),
			expectError: false,
		},
		{
			name:        "Invalid - zero fee but non-zero share",
			integrator:  addresses.ZeroAddress,
			protocol:    addresses.ZeroAddress,
			fee:         fusionorder.BpsZero,
			share:       fusionorder.MustNewBps(big.NewInt(100)),
			expectError: true,
			errorMsg:    "zero fee requires zero integrator share",
		},
		{
			name:        "Invalid - zero fee but non-zero integrator",
			integrator:  validAddress1,
			protocol:    addresses.ZeroAddress,
			fee:         fusionorder.BpsZero,
			share:       fusionorder.BpsZero,
			expectError: true,
			errorMsg:    "zero fee requires zero integrator address",
		},
		{
			name:        "Invalid - zero fee but non-zero protocol",
			integrator:  addresses.ZeroAddress,
			protocol:    validAddress2,
			fee:         fusionorder.BpsZero,
			share:       fusionorder.BpsZero,
			expectError: true,
			errorMsg:    "zero fee requires zero protocol address",
		},
		{
			name:        "Invalid - non-zero fee with zero integrator",
			integrator:  addresses.ZeroAddress,
			protocol:    validAddress2,
			fee:         fusionorder.MustNewBps(big.NewInt(100)),
			share:       fusionorder.MustNewBps(big.NewInt(5000)),
			expectError: true,
			errorMsg:    "non-zero fee requires non-zero integrator and protocol addresses",
		},
		{
			name:        "Invalid - non-zero fee with zero protocol",
			integrator:  validAddress1,
			protocol:    addresses.ZeroAddress,
			fee:         fusionorder.MustNewBps(big.NewInt(100)),
			share:       fusionorder.MustNewBps(big.NewInt(5000)),
			expectError: true,
			errorMsg:    "non-zero fee requires non-zero integrator and protocol addresses",
		},
		{
			name:        "Invalid - both addresses zero but fee non-zero",
			integrator:  addresses.ZeroAddress,
			protocol:    addresses.ZeroAddress,
			fee:         fusionorder.MustNewBps(big.NewInt(100)),
			share:       fusionorder.MustNewBps(big.NewInt(5000)),
			expectError: true,
			errorMsg:    "non-zero fee requires non-zero integrator and protocol addresses",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := NewIntegratorFee(tc.integrator, tc.protocol, tc.fee, tc.share)
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tc.integrator, result.Integrator)
				assert.Equal(t, tc.protocol, result.Protocol)
				assert.True(t, tc.fee.Equal(result.Fee))
				assert.True(t, tc.share.Equal(result.Share))
			}
		})
	}
}

func TestIntegratorFee_String(t *testing.T) {
	validAddress1 := "0x6B175474E89094C44Da98b954EedeAC495271d0F"
	validAddress2 := "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"

	fee, err := NewIntegratorFee(validAddress1, validAddress2, fusionorder.MustNewBps(big.NewInt(100)), fusionorder.MustNewBps(big.NewInt(5000)))
	require.NoError(t, err)

	result := fee.String()
	assert.Contains(t, result, "IntegratorFee")
	assert.Contains(t, result, validAddress1)
	assert.Contains(t, result, validAddress2)
	assert.Contains(t, result, "100")
	assert.Contains(t, result, "5000")
}

func TestIntegratorFee_ZeroFeeString(t *testing.T) {
	result := IntegratorFeeZero.String()
	assert.Contains(t, result, "IntegratorFee")
	assert.Contains(t, result, addresses.ZeroAddress)
	assert.Contains(t, result, "0")
}
