package fusionplus

import "github.com/ethereum/go-ethereum/common"

const NativeToken = "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"

type NetworkEnum int

const (
	ETHEREUM  NetworkEnum = 1
	POLYGON   NetworkEnum = 137
	BINANCE   NetworkEnum = 56
	ARBITRUM  NetworkEnum = 42161
	AVALANCHE NetworkEnum = 43114
	OPTIMISM  NetworkEnum = 10
	FANTOM    NetworkEnum = 250
	GNOSIS    NetworkEnum = 100
	COINBASE  NetworkEnum = 8453
)

var chainToWrapper = map[NetworkEnum]common.Address{
	ETHEREUM:  common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
	BINANCE:   common.HexToAddress("0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c"),
	POLYGON:   common.HexToAddress("0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270"),
	ARBITRUM:  common.HexToAddress("0x82af49447d8a07e3bd95bd0d56f35241523fbab1"),
	AVALANCHE: common.HexToAddress("0xb31f66aa3c1e785363f0875a1b74e27b85fd66c7"),
	GNOSIS:    common.HexToAddress("0xe91d153e0b41518a2ce8dd3d7944fa863463a97d"),
	COINBASE:  common.HexToAddress("0x4200000000000000000000000000000000000006"),
	OPTIMISM:  common.HexToAddress("0x4200000000000000000000000000000000000006"),
	FANTOM:    common.HexToAddress("0x21be370d5312f44cb42ce377bc9b8a0cef1a4c83"),
}
