package fusion

import (
	"math/big"
	"testing"

	"github.com/1inch/1inch-sdk-go/common/fusionorder"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateWhitelist(t *testing.T) {
	tests := []struct {
		name               string
		whitelistStrings   []string
		resolvingStartTime *big.Int
		expected           []fusionorder.WhitelistItem
		expectError        bool
		errorMsg           string
	}{
		{
			name:               "Should generate whitelist",
			whitelistStrings:   []string{"0x00000000219ab540356cbb839cbe05303d7705fa"},
			resolvingStartTime: big.NewInt(1708117482),
			expected: []fusionorder.WhitelistItem{
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
			expected: []fusionorder.WhitelistItem{
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
			whitelist, err := fusionorder.GenerateWhitelist(tc.whitelistStrings, tc.resolvingStartTime)
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

func TestSettlementPostInteractionData_CanExecuteAt(t *testing.T) {
	// Address half is last 20 hex chars (10 bytes) of address
	resolver1 := common.HexToAddress("0x1234567890123456789012345678901234567890")
	resolver2 := common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")
	resolvingStartTime := big.NewInt(1000)

	whitelist := []fusionorder.WhitelistItem{
		{
			AddressHalf: "12345678901234567890", // last 20 chars of resolver1
			Delay:       big.NewInt(0),
		},
		{
			AddressHalf: "cdefabcdefabcdefabcd", // last 20 chars of resolver2
			Delay:       big.NewInt(100),
		},
	}

	spid := SettlementPostInteractionData{
		Whitelist:          whitelist,
		ResolvingStartTime: resolvingStartTime,
	}

	tests := []struct {
		name          string
		executor      common.Address
		executionTime *big.Int
		expected      bool
	}{
		{
			name:          "First resolver can execute immediately",
			executor:      resolver1,
			executionTime: big.NewInt(1000),
			expected:      true,
		},
		{
			name:          "Second resolver cannot execute before delay",
			executor:      resolver2,
			executionTime: big.NewInt(1050),
			expected:      false,
		},
		{
			name:          "Second resolver can execute after delay",
			executor:      resolver2,
			executionTime: big.NewInt(1100),
			expected:      true,
		},
		{
			name:          "Unknown resolver cannot execute during exclusive period",
			executor:      common.HexToAddress("0x9999999999999999999999999999999999999999"),
			executionTime: big.NewInt(1000),
			expected:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := spid.CanExecuteAt(tc.executor, tc.executionTime)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSettlementPostInteractionData_IsExclusiveResolver(t *testing.T) {
	resolver1 := common.HexToAddress("0x1234567890123456789012345678901234567890")
	resolver2 := common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")

	whitelist := []fusionorder.WhitelistItem{
		{
			AddressHalf: "12345678901234567890", // last 20 chars of resolver1
			Delay:       big.NewInt(0),          // exclusive (delay = 0)
		},
		{
			AddressHalf: "cdefabcdefabcdefabcd", // last 20 chars of resolver2
			Delay:       big.NewInt(100),        // not exclusive (delay > 0)
		},
	}

	spid := SettlementPostInteractionData{
		Whitelist:          whitelist,
		ResolvingStartTime: big.NewInt(1000),
	}

	tests := []struct {
		name     string
		wallet   common.Address
		expected bool
	}{
		{
			name:     "First resolver is exclusive (delay=0)",
			wallet:   resolver1,
			expected: true,
		},
		{
			name:     "Second resolver is not exclusive (delay>0)",
			wallet:   resolver2,
			expected: false,
		},
		{
			name:     "Unknown resolver is not exclusive",
			wallet:   common.HexToAddress("0x9999999999999999999999999999999999999999"),
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := spid.IsExclusiveResolver(tc.wallet)
			assert.Equal(t, tc.expected, result)
		})
	}
}
