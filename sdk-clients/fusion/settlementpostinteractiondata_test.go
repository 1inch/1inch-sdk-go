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
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			whitelist, err := GenerateWhitelist(tc.whitelistStrings, tc.resolvingStartTime)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, whitelist)
		})
	}
}
