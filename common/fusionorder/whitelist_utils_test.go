package fusionorder

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestCanExecuteAt(t *testing.T) {
	// Create a test whitelist - AddressHalf is last 20 hex chars (10 bytes) of address
	whitelist := []WhitelistItem{
		{AddressHalf: "1234567890abcdef1234", Delay: big.NewInt(0)},
		{AddressHalf: "abcdef12345678901234", Delay: big.NewInt(100)},
	}
	resolvingStartTime := big.NewInt(1000)

	tests := []struct {
		name          string
		executor      common.Address
		executionTime *big.Int
		expected      bool
	}{
		{
			name:          "First resolver can execute immediately",
			executor:      common.HexToAddress("0x00000000000000000000001234567890abcdef1234"),
			executionTime: big.NewInt(1000),
			expected:      true,
		},
		{
			name:          "Second resolver cannot execute before delay",
			executor:      common.HexToAddress("0x0000000000000000000000abcdef12345678901234"),
			executionTime: big.NewInt(1050),
			expected:      false,
		},
		{
			name:          "Second resolver can execute after delay",
			executor:      common.HexToAddress("0x0000000000000000000000abcdef12345678901234"),
			executionTime: big.NewInt(1100),
			expected:      true,
		},
		{
			name:          "Unknown executor cannot execute",
			executor:      common.HexToAddress("0x0000000000000000000000ffffffffffffffffffff"),
			executionTime: big.NewInt(2000),
			expected:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CanExecuteAt(whitelist, resolvingStartTime, tc.executor, tc.executionTime)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCanExecuteAt_EmptyWhitelist(t *testing.T) {
	whitelist := []WhitelistItem{}
	resolvingStartTime := big.NewInt(1000)
	executor := common.HexToAddress("0x1234567890123456789012345678901234567890")
	executionTime := big.NewInt(2000)

	result := CanExecuteAt(whitelist, resolvingStartTime, executor, executionTime)
	assert.False(t, result, "Empty whitelist should return false")
}

func TestIsExclusiveResolver(t *testing.T) {
	tests := []struct {
		name      string
		whitelist []WhitelistItem
		wallet    common.Address
		expected  bool
	}{
		{
			name: "Single resolver is exclusive",
			whitelist: []WhitelistItem{
				{AddressHalf: "1234567890abcdef1234", Delay: big.NewInt(0)},
			},
			wallet:   common.HexToAddress("0x00000000000000000000001234567890abcdef1234"),
			expected: true,
		},
		{
			name: "Single resolver - different wallet not exclusive",
			whitelist: []WhitelistItem{
				{AddressHalf: "1234567890abcdef1234", Delay: big.NewInt(0)},
			},
			wallet:   common.HexToAddress("0x0000000000000000000000ffffffffffffffffffff"),
			expected: false,
		},
		{
			name: "Multiple resolvers with different delays - first is exclusive",
			whitelist: []WhitelistItem{
				{AddressHalf: "1234567890abcdef1234", Delay: big.NewInt(0)},
				{AddressHalf: "abcdef12345678901234", Delay: big.NewInt(100)},
			},
			wallet:   common.HexToAddress("0x00000000000000000000001234567890abcdef1234"),
			expected: true,
		},
		{
			name: "Multiple resolvers with same delay - no exclusive",
			whitelist: []WhitelistItem{
				{AddressHalf: "1234567890abcdef1234", Delay: big.NewInt(0)},
				{AddressHalf: "abcdef12345678901234", Delay: big.NewInt(0)},
			},
			wallet:   common.HexToAddress("0x00000000000000000000001234567890abcdef1234"),
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsExclusiveResolver(tc.whitelist, tc.wallet)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsExclusiveResolver_EmptyWhitelist(t *testing.T) {
	whitelist := []WhitelistItem{}
	wallet := common.HexToAddress("0x1234567890123456789012345678901234567890")

	result := IsExclusiveResolver(whitelist, wallet)
	assert.False(t, result, "Empty whitelist should return false")
}
