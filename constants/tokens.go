package constants

import "github.com/ethereum/go-ethereum/common"

// ZeroAddress is the Ethereum zero address (0x0000...0000)
const ZeroAddress = "0x0000000000000000000000000000000000000000"

// NativeToken is the address used to represent the native token (ETH, MATIC, etc.)
// Uses EIP-55 checksummed format for consistency
const NativeToken = "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"

// NetworkEnum represents supported network chain IDs for wrapped token lookups
type NetworkEnum int

const (
	NetworkEthereum  NetworkEnum = EthereumChainId
	NetworkPolygon   NetworkEnum = PolygonChainId
	NetworkBinance   NetworkEnum = BscChainId
	NetworkArbitrum  NetworkEnum = ArbitrumChainId
	NetworkAvalanche NetworkEnum = AvalancheChainId
	NetworkOptimism  NetworkEnum = OptimismChainId
	NetworkFantom    NetworkEnum = FantomChainId
	NetworkGnosis    NetworkEnum = GnosisChainId
	NetworkBase      NetworkEnum = BaseChainId
)

// ChainToWrapper maps chain IDs to their wrapped native token addresses
var ChainToWrapper = map[NetworkEnum]common.Address{
	NetworkEthereum:  common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"), // WETH
	NetworkBinance:   common.HexToAddress("0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c"), // WBNB
	NetworkPolygon:   common.HexToAddress("0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270"), // WMATIC
	NetworkArbitrum:  common.HexToAddress("0x82af49447d8a07e3bd95bd0d56f35241523fbab1"), // WETH
	NetworkAvalanche: common.HexToAddress("0xb31f66aa3c1e785363f0875a1b74e27b85fd66c7"), // WAVAX
	NetworkGnosis:    common.HexToAddress("0xe91d153e0b41518a2ce8dd3d7944fa863463a97d"), // WXDAI
	NetworkBase:      common.HexToAddress("0x4200000000000000000000000000000000000006"), // WETH
	NetworkOptimism:  common.HexToAddress("0x4200000000000000000000000000000000000006"), // WETH
	NetworkFantom:    common.HexToAddress("0x21be370d5312f44cb42ce377bc9b8a0cef1a4c83"), // WFTM
}

// GetWrappedToken returns the wrapped token address for a given chain
func GetWrappedToken(chainID NetworkEnum) (common.Address, bool) {
	addr, ok := ChainToWrapper[chainID]
	return addr, ok
}
