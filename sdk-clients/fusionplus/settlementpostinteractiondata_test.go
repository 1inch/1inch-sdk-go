package fusionplus

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestSettlementPostInteractionData_CanExecuteAt(t *testing.T) {
	// Address half is last 20 hex chars (10 bytes) of address
	resolver1 := common.HexToAddress("0x1234567890123456789012345678901234567890")
	resolver2 := common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")
	resolvingStartTime := big.NewInt(1000)

	whitelist := []WhitelistItem{
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

	whitelist := []WhitelistItem{
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
