package fusion

import (
	"testing"

	"github.com/1inch/1inch-sdk-go/common/fusionorder"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

// These tests verify that fusionorder.NativeToken and related functions work correctly
// when accessed through the fusionorder package (previously exported from fusion)

func TestNativeToken(t *testing.T) {
	// Verify the native token constant is the standard 0xEeee...EEEE address
	assert.Equal(t, "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE", fusionorder.NativeToken)

	// Verify it's a valid address format (42 characters with 0x prefix)
	assert.Len(t, fusionorder.NativeToken, 42)
	assert.Equal(t, "0x", fusionorder.NativeToken[:2])
}

func TestNetworkEnum(t *testing.T) {
	// Verify network IDs match the expected chain IDs
	assert.Equal(t, fusionorder.NetworkEnum(1), fusionorder.ETHEREUM)
	assert.Equal(t, fusionorder.NetworkEnum(137), fusionorder.POLYGON)
	assert.Equal(t, fusionorder.NetworkEnum(56), fusionorder.BINANCE)
	assert.Equal(t, fusionorder.NetworkEnum(42161), fusionorder.ARBITRUM)
	assert.Equal(t, fusionorder.NetworkEnum(43114), fusionorder.AVALANCHE)
	assert.Equal(t, fusionorder.NetworkEnum(10), fusionorder.OPTIMISM)
	assert.Equal(t, fusionorder.NetworkEnum(250), fusionorder.FANTOM)
	assert.Equal(t, fusionorder.NetworkEnum(100), fusionorder.GNOSIS)
	assert.Equal(t, fusionorder.NetworkEnum(8453), fusionorder.COINBASE)
}

func TestChainToWrapper(t *testing.T) {
	// Test that all supported networks have a wrapper address
	supportedNetworks := []fusionorder.NetworkEnum{
		fusionorder.ETHEREUM,
		fusionorder.POLYGON,
		fusionorder.BINANCE,
		fusionorder.ARBITRUM,
		fusionorder.AVALANCHE,
		fusionorder.OPTIMISM,
		fusionorder.FANTOM,
		fusionorder.GNOSIS,
		fusionorder.COINBASE,
	}

	for _, network := range supportedNetworks {
		wrapper, exists := fusionorder.ChainToWrapper[network]
		assert.True(t, exists, "Network %d should have a wrapper address", network)
		assert.NotEqual(t, common.Address{}, wrapper, "Wrapper address for network %d should not be zero address", network)
	}
}

func TestChainToWrapper_WellKnownAddresses(t *testing.T) {
	tests := []struct {
		name     string
		network  fusionorder.NetworkEnum
		expected string
	}{
		{
			name:     "Ethereum WETH",
			network:  fusionorder.ETHEREUM,
			expected: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		},
		{
			name:     "Polygon WMATIC",
			network:  fusionorder.POLYGON,
			expected: "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270",
		},
		{
			name:     "Binance WBNB",
			network:  fusionorder.BINANCE,
			expected: "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c",
		},
		{
			name:     "Arbitrum WETH",
			network:  fusionorder.ARBITRUM,
			expected: "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wrapper := fusionorder.ChainToWrapper[tc.network]
			// Compare case-insensitively since addresses are case-insensitive
			assert.Equal(t, common.HexToAddress(tc.expected), wrapper)
		})
	}
}

func TestChainToWrapper_UnsupportedNetwork(t *testing.T) {
	// Test that unsupported networks return zero address
	unsupportedNetwork := fusionorder.NetworkEnum(999999)
	wrapper, exists := fusionorder.ChainToWrapper[unsupportedNetwork]
	assert.False(t, exists)
	assert.Equal(t, common.Address{}, wrapper)
}
