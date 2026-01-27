package fusion

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestNativeToken(t *testing.T) {
	// Verify the native token constant is the standard 0xEeee...EEEE address
	assert.Equal(t, "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE", NativeToken)

	// Verify it's a valid address format (42 characters with 0x prefix)
	assert.Len(t, NativeToken, 42)
	assert.Equal(t, "0x", NativeToken[:2])
}

func TestNetworkEnum(t *testing.T) {
	// Verify network IDs match the expected chain IDs
	assert.Equal(t, NetworkEnum(1), ETHEREUM)
	assert.Equal(t, NetworkEnum(137), POLYGON)
	assert.Equal(t, NetworkEnum(56), BINANCE)
	assert.Equal(t, NetworkEnum(42161), ARBITRUM)
	assert.Equal(t, NetworkEnum(43114), AVALANCHE)
	assert.Equal(t, NetworkEnum(10), OPTIMISM)
	assert.Equal(t, NetworkEnum(250), FANTOM)
	assert.Equal(t, NetworkEnum(100), GNOSIS)
	assert.Equal(t, NetworkEnum(8453), COINBASE)
}

func TestChainToWrapper(t *testing.T) {
	// Test that all supported networks have a wrapper address
	supportedNetworks := []NetworkEnum{
		ETHEREUM,
		POLYGON,
		BINANCE,
		ARBITRUM,
		AVALANCHE,
		OPTIMISM,
		FANTOM,
		GNOSIS,
		COINBASE,
	}

	for _, network := range supportedNetworks {
		wrapper, exists := chainToWrapper[network]
		assert.True(t, exists, "Network %d should have a wrapper address", network)
		assert.NotEqual(t, common.Address{}, wrapper, "Wrapper address for network %d should not be zero address", network)
	}
}

func TestChainToWrapper_WellKnownAddresses(t *testing.T) {
	tests := []struct {
		name     string
		network  NetworkEnum
		expected string
	}{
		{
			name:     "Ethereum WETH",
			network:  ETHEREUM,
			expected: "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		},
		{
			name:     "Polygon WMATIC",
			network:  POLYGON,
			expected: "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270",
		},
		{
			name:     "Binance WBNB",
			network:  BINANCE,
			expected: "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c",
		},
		{
			name:     "Arbitrum WETH",
			network:  ARBITRUM,
			expected: "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wrapper := chainToWrapper[tc.network]
			// Compare case-insensitively since addresses are case-insensitive
			assert.Equal(t, common.HexToAddress(tc.expected), wrapper)
		})
	}
}

func TestChainToWrapper_UnsupportedNetwork(t *testing.T) {
	// Test that unsupported networks return zero address
	unsupportedNetwork := NetworkEnum(999999)
	wrapper, exists := chainToWrapper[unsupportedNetwork]
	assert.False(t, exists)
	assert.Equal(t, common.Address{}, wrapper)
}
