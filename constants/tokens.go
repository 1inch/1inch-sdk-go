package constants

import "github.com/ethereum/go-ethereum/common"

// ZeroAddress is the Ethereum zero address (0x0000...0000)
const ZeroAddress = "0x0000000000000000000000000000000000000000"

// NativeToken is the address used to represent the native token (ETH, MATIC, etc.)
// Uses EIP-55 checksummed format for consistency
const NativeToken = "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"

// Deprecated: Use the chain-id int constants (EthereumChainId, …) directly.
// NetworkEnum duplicated the chain-id constants and was an incomplete parallel
// representation; the chain→address maps are now keyed by int chain id to match
// ChainToTrueERC20/GetTrueERC20.
type NetworkEnum int

const (
	// Deprecated: Use EthereumChainId.
	NetworkEthereum NetworkEnum = EthereumChainId
	// Deprecated: Use PolygonChainId.
	NetworkPolygon NetworkEnum = PolygonChainId
	// Deprecated: Use BscChainId.
	NetworkBinance NetworkEnum = BscChainId
	// Deprecated: Use ArbitrumChainId.
	NetworkArbitrum NetworkEnum = ArbitrumChainId
	// Deprecated: Use AvalancheChainId.
	NetworkAvalanche NetworkEnum = AvalancheChainId
	// Deprecated: Use OptimismChainId.
	NetworkOptimism NetworkEnum = OptimismChainId
	// Deprecated: Use FantomChainId.
	NetworkFantom NetworkEnum = FantomChainId
	// Deprecated: Use GnosisChainId.
	NetworkGnosis NetworkEnum = GnosisChainId
	// Deprecated: Use BaseChainId.
	NetworkBase NetworkEnum = BaseChainId
)

// ChainToWrapper maps a chain id to its wrapped native token address. Keyed by
// int chain id, consistent with ChainToTrueERC20.
var ChainToWrapper = map[int]common.Address{
	EthereumChainId:  common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"), // WETH
	BscChainId:       common.HexToAddress("0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c"), // WBNB
	PolygonChainId:   common.HexToAddress("0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270"), // WMATIC
	ArbitrumChainId:  common.HexToAddress("0x82af49447d8a07e3bd95bd0d56f35241523fbab1"), // WETH
	AvalancheChainId: common.HexToAddress("0xb31f66aa3c1e785363f0875a1b74e27b85fd66c7"), // WAVAX
	GnosisChainId:    common.HexToAddress("0xe91d153e0b41518a2ce8dd3d7944fa863463a97d"), // WXDAI
	BaseChainId:      common.HexToAddress("0x4200000000000000000000000000000000000006"), // WETH
	OptimismChainId:  common.HexToAddress("0x4200000000000000000000000000000000000006"), // WETH
	FantomChainId:    common.HexToAddress("0x21be370d5312f44cb42ce377bc9b8a0cef1a4c83"), // WFTM
}

// GetWrappedToken returns the wrapped native token address for a chain id,
// matching the GetTrueERC20 signature.
func GetWrappedToken(chainID int) (common.Address, bool) {
	addr, ok := ChainToWrapper[chainID]
	return addr, ok
}
