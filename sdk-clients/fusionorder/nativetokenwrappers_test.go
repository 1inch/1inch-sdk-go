package fusionorder

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestNativeTokenConstant(t *testing.T) {
	assert.Equal(t, "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE", NativeToken)
}

func TestNetworkEnumValues(t *testing.T) {
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

func TestGetWrappedToken(t *testing.T) {
	tests := []struct {
		name     string
		chainID  NetworkEnum
		expected common.Address
		found    bool
	}{
		{
			name:     "Ethereum WETH",
			chainID:  ETHEREUM,
			expected: common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
			found:    true,
		},
		{
			name:     "Polygon WMATIC",
			chainID:  POLYGON,
			expected: common.HexToAddress("0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"),
			found:    true,
		},
		{
			name:     "BSC WBNB",
			chainID:  BINANCE,
			expected: common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"),
			found:    true,
		},
		{
			name:     "Arbitrum WETH",
			chainID:  ARBITRUM,
			expected: common.HexToAddress("0x82aF49447D8a07e3bd95BD0d56f35241523fBab1"),
			found:    true,
		},
		{
			name:     "Avalanche WAVAX",
			chainID:  AVALANCHE,
			expected: common.HexToAddress("0xB31f66AA3C1e785363F0875A1B74E27b85FD66c7"),
			found:    true,
		},
		{
			name:     "Optimism WETH",
			chainID:  OPTIMISM,
			expected: common.HexToAddress("0x4200000000000000000000000000000000000006"),
			found:    true,
		},
		{
			name:     "Fantom WFTM",
			chainID:  FANTOM,
			expected: common.HexToAddress("0x21be370D5312f44cB42ce377BC9b8a0cEF1A4C83"),
			found:    true,
		},
		{
			name:     "Gnosis WXDAI",
			chainID:  GNOSIS,
			expected: common.HexToAddress("0xe91D153E0b41518A2Ce8Dd3D7944Fa863463a97d"),
			found:    true,
		},
		{
			name:     "Coinbase WETH",
			chainID:  COINBASE,
			expected: common.HexToAddress("0x4200000000000000000000000000000000000006"),
			found:    true,
		},
		{
			name:     "Unknown chain",
			chainID:  NetworkEnum(9999),
			expected: common.Address{},
			found:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, found := GetWrappedToken(tc.chainID)
			assert.Equal(t, tc.found, found)
			if tc.found {
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestChainToWrapperMap(t *testing.T) {
	// Verify all expected chains are in the map
	expectedChains := []NetworkEnum{
		ETHEREUM, POLYGON, BINANCE, ARBITRUM, AVALANCHE,
		OPTIMISM, FANTOM, GNOSIS, COINBASE,
	}

	for _, chain := range expectedChains {
		_, exists := ChainToWrapper[chain]
		assert.True(t, exists, "Chain %d should exist in ChainToWrapper map", chain)
	}
}
