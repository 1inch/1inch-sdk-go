package fusion

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateWhitelist(t *testing.T) {
	tests := []struct {
		name               string
		whitelistStrings   []string
		resolvingStartTime *big.Int
		expected           []WhitelistItem
		expectError        bool
		errorMsg           string
	}{
		{
			name:               "Should generate whitelist",
			whitelistStrings:   []string{"0x00000000219ab540356cbb839cbe05303d7705fa"},
			resolvingStartTime: big.NewInt(1708117482),
			expected: []WhitelistItem{
				{
					AddressHalf: "bb839cbe05303d7705fa",
					Delay:       big.NewInt(0),
				},
			},
			expectError: false,
		},
		{
			name:               "Should generate whitelist with multiple addresses",
			whitelistStrings:   []string{"0x00000000219ab540356cbb839cbe05303d7705fa", "0x1234567890123456789012345678901234567890"},
			resolvingStartTime: big.NewInt(1708117482),
			expected: []WhitelistItem{
				{
					AddressHalf: "bb839cbe05303d7705fa",
					Delay:       big.NewInt(0),
				},
				{
					AddressHalf: "12345678901234567890",
					Delay:       big.NewInt(0),
				},
			},
			expectError: false,
		},
		{
			name:               "Empty whitelist should return error",
			whitelistStrings:   []string{},
			resolvingStartTime: big.NewInt(1708117482),
			expected:           nil,
			expectError:        true,
			errorMsg:           "whitelist cannot be empty",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			whitelist, err := GenerateWhitelist(tc.whitelistStrings, tc.resolvingStartTime)
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, whitelist)
			}
		})
	}
}
