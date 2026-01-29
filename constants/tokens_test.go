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
	// Verify NetworkEnum values match their corresponding chain IDs
	assert.Equal(t, NetworkEnum(EthereumChainId), NetworkEthereum)
	assert.Equal(t, NetworkEnum(PolygonChainId), NetworkPolygon)
	assert.Equal(t, NetworkEnum(BscChainId), NetworkBinance)
	assert.Equal(t, NetworkEnum(ArbitrumChainId), NetworkArbitrum)
	assert.Equal(t, NetworkEnum(AvalancheChainId), NetworkAvalanche)
	assert.Equal(t, NetworkEnum(OptimismChainId), NetworkOptimism)
	assert.Equal(t, NetworkEnum(FantomChainId), NetworkFantom)
	assert.Equal(t, NetworkEnum(GnosisChainId), NetworkGnosis)
	assert.Equal(t, NetworkEnum(BaseChainId), NetworkBase)
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
			chainID:  NetworkEthereum,
			expected: common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
			found:    true,
		},
		{
			name:     "Polygon WMATIC",
			chainID:  NetworkPolygon,
			expected: common.HexToAddress("0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"),
			found:    true,
		},
		{
			name:     "BSC WBNB",
			chainID:  NetworkBinance,
			expected: common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"),
			found:    true,
		},
		{
			name:     "Arbitrum WETH",
			chainID:  NetworkArbitrum,
			expected: common.HexToAddress("0x82aF49447D8a07e3bd95BD0d56f35241523fBab1"),
			found:    true,
		},
		{
			name:     "Avalanche WAVAX",
			chainID:  NetworkAvalanche,
			expected: common.HexToAddress("0xB31f66AA3C1e785363F0875A1B74E27b85FD66c7"),
			found:    true,
		},
		{
			name:     "Optimism WETH",
			chainID:  NetworkOptimism,
			expected: common.HexToAddress("0x4200000000000000000000000000000000000006"),
			found:    true,
		},
		{
			name:     "Fantom WFTM",
			chainID:  NetworkFantom,
			expected: common.HexToAddress("0x21be370D5312f44cB42ce377BC9b8a0cEF1A4C83"),
			found:    true,
		},
		{
			name:     "Gnosis WXDAI",
			chainID:  NetworkGnosis,
			expected: common.HexToAddress("0xe91D153E0b41518A2Ce8Dd3D7944Fa863463a97d"),
			found:    true,
		},
		{
			name:     "Base WETH",
			chainID:  NetworkBase,
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
		NetworkEthereum, NetworkPolygon, NetworkBinance, NetworkArbitrum, NetworkAvalanche,
		NetworkOptimism, NetworkFantom, NetworkGnosis, NetworkBase,
	}

	for _, chain := range expectedChains {
		_, exists := ChainToWrapper[chain]
		require.True(t, exists, "Chain %d should exist in ChainToWrapper map", chain)
	}
}
