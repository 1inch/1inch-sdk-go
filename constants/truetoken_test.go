package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTrueERC20(t *testing.T) {
	const shared = "0xDA0000d4000015A526378bB6faFc650Cea5966F8"

	tests := []struct {
		name     string
		chainID  int
		expected string
		ok       bool
	}{
		{name: "Ethereum", chainID: EthereumChainId, expected: shared, ok: true},
		{name: "Polygon", chainID: PolygonChainId, expected: shared, ok: true},
		{name: "BSC", chainID: BscChainId, expected: shared, ok: true},
		{name: "Arbitrum", chainID: ArbitrumChainId, expected: shared, ok: true},
		{name: "Base", chainID: BaseChainId, expected: shared, ok: true},
		{name: "Optimism", chainID: OptimismChainId, expected: shared, ok: true},
		{name: "Avalanche", chainID: AvalancheChainId, expected: shared, ok: true},
		{name: "Gnosis", chainID: GnosisChainId, expected: shared, ok: true},
		{name: "Fantom", chainID: FantomChainId, expected: shared, ok: true},
		{name: "zkSync Era has a distinct sentinel", chainID: ZkSyncEraChainId, expected: "0xD66097C27eB8dEe404bAC235737932260EdC6f3b", ok: true},
		{name: "Aurora is not a Fusion+ chain", chainID: AuroraChainId, ok: false},
		{name: "Klaytn is not a Fusion+ chain", chainID: KlaytnChainId, ok: false},
		{name: "Unknown chain", chainID: 999999, ok: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			addr, ok := GetTrueERC20(tc.chainID)
			assert.Equal(t, tc.ok, ok)
			if tc.ok {
				assert.Equal(t, tc.expected, addr.Hex())
			}
		})
	}
}

// TestTrueERC20NeverCollidesWithNativeSentinel guards the whole point of the
// sentinel: it must never equal the native-token address, or the source-chain
// taker asset could move real value.
func TestTrueERC20NeverCollidesWithNativeSentinel(t *testing.T) {
	for chainID, addr := range ChainToTrueERC20 {
		require.NotEqualf(t, NativeToken, addr.Hex(), "TRUE_ERC20 for chain %d equals the native sentinel", chainID)
	}
}
