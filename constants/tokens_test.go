package constants

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZeroAddressConstant(t *testing.T) {
	// Verify the zero address constant
	assert.Equal(t, "0x0000000000000000000000000000000000000000", ZeroAddress)

	// Verify it's a valid address format (42 characters with 0x prefix)
	assert.Len(t, ZeroAddress, 42)
	assert.Equal(t, "0x", ZeroAddress[:2])

	// Verify all characters after 0x are zeros
	for _, c := range ZeroAddress[2:] {
		assert.Equal(t, '0', c)
	}
}

func TestNativeTokenConstant(t *testing.T) {
	// Verify the native token constant is the standard checksummed address
	assert.Equal(t, "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE", NativeToken)

	// Verify it's a valid address format (42 characters with 0x prefix)
	assert.Len(t, NativeToken, 42)
	assert.Equal(t, "0x", NativeToken[:2])
}

func TestNetworkEnumValues(t *testing.T) {
	// The deprecated NetworkEnum members still resolve to their chain ids.
	assert.Equal(t, EthereumChainId, int(NetworkEthereum))
	assert.Equal(t, PolygonChainId, int(NetworkPolygon))
	assert.Equal(t, BscChainId, int(NetworkBinance))
	assert.Equal(t, ArbitrumChainId, int(NetworkArbitrum))
	assert.Equal(t, AvalancheChainId, int(NetworkAvalanche))
	assert.Equal(t, OptimismChainId, int(NetworkOptimism))
	assert.Equal(t, FantomChainId, int(NetworkFantom))
	assert.Equal(t, GnosisChainId, int(NetworkGnosis))
	assert.Equal(t, BaseChainId, int(NetworkBase))
}

func TestGetWrappedToken(t *testing.T) {
	tests := []struct {
		name     string
		chainID  int
		expected common.Address
		found    bool
	}{
		{
			name:     "Ethereum WETH",
			chainID:  EthereumChainId,
			expected: common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
			found:    true,
		},
		{
			name:     "Polygon WMATIC",
			chainID:  PolygonChainId,
			expected: common.HexToAddress("0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"),
			found:    true,
		},
		{
			name:     "BSC WBNB",
			chainID:  BscChainId,
			expected: common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"),
			found:    true,
		},
		{
			name:     "Arbitrum WETH",
			chainID:  ArbitrumChainId,
			expected: common.HexToAddress("0x82aF49447D8a07e3bd95BD0d56f35241523fBab1"),
			found:    true,
		},
		{
			name:     "Avalanche WAVAX",
			chainID:  AvalancheChainId,
			expected: common.HexToAddress("0xB31f66AA3C1e785363F0875A1B74E27b85FD66c7"),
			found:    true,
		},
		{
			name:     "Optimism WETH",
			chainID:  OptimismChainId,
			expected: common.HexToAddress("0x4200000000000000000000000000000000000006"),
			found:    true,
		},
		{
			name:     "Fantom WFTM",
			chainID:  FantomChainId,
			expected: common.HexToAddress("0x21be370D5312f44cB42ce377BC9b8a0cEF1A4C83"),
			found:    true,
		},
		{
			name:     "Gnosis WXDAI",
			chainID:  GnosisChainId,
			expected: common.HexToAddress("0xe91D153E0b41518A2Ce8Dd3D7944Fa863463a97d"),
			found:    true,
		},
		{
			name:     "Base WETH",
			chainID:  BaseChainId,
			expected: common.HexToAddress("0x4200000000000000000000000000000000000006"),
			found:    true,
		},
		{
			name:     "Unknown chain",
			chainID:  9999,
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
	expectedChains := []int{
		EthereumChainId, PolygonChainId, BscChainId, ArbitrumChainId, AvalancheChainId,
		OptimismChainId, FantomChainId, GnosisChainId, BaseChainId,
	}

	for _, chain := range expectedChains {
		_, exists := ChainToWrapper[chain]
		require.True(t, exists, "Chain %d should exist in ChainToWrapper map", chain)
	}
}
