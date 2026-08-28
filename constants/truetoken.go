package constants

import "github.com/ethereum/go-ethereum/common"

// ChainToTrueERC20 maps a chain id to its Fusion+ "true ERC20" sentinel token.
//
// On the SOURCE chain of a Fusion+ order, the order's taker asset must be this
// sentinel — a token whose transfer is a no-op — not the real destination
// token. The destination token is carried only by the escrow extension and is
// delivered on the destination chain. Putting the destination token in the
// source-chain taker asset can cause a real, unintended transfer on the source
// chain when that address is also a live ERC-20 there.
//
// A chain absent from this map is not supported by Fusion+; order construction
// for it fails loudly rather than signing an order with a wrong taker asset.
//
// Source: @1inch/cross-chain-sdk deployments.ts (TRUE_ERC20 map), verified
// 2026-08-28. Every supported chain shares one deterministically-deployed
// address except zkSync Era.
var ChainToTrueERC20 = map[int]common.Address{
	EthereumChainId:  common.HexToAddress("0xda0000d4000015a526378bb6fafc650cea5966f8"),
	PolygonChainId:   common.HexToAddress("0xda0000d4000015a526378bb6fafc650cea5966f8"),
	BscChainId:       common.HexToAddress("0xda0000d4000015a526378bb6fafc650cea5966f8"),
	ArbitrumChainId:  common.HexToAddress("0xda0000d4000015a526378bb6fafc650cea5966f8"),
	AvalancheChainId: common.HexToAddress("0xda0000d4000015a526378bb6fafc650cea5966f8"),
	OptimismChainId:  common.HexToAddress("0xda0000d4000015a526378bb6fafc650cea5966f8"),
	FantomChainId:    common.HexToAddress("0xda0000d4000015a526378bb6fafc650cea5966f8"),
	GnosisChainId:    common.HexToAddress("0xda0000d4000015a526378bb6fafc650cea5966f8"),
	BaseChainId:      common.HexToAddress("0xda0000d4000015a526378bb6fafc650cea5966f8"),
	ZkSyncEraChainId: common.HexToAddress("0xd66097c27eb8dee404bac235737932260edc6f3b"),
}

// GetTrueERC20 returns the Fusion+ true-ERC20 sentinel for a chain and whether
// the chain is a supported Fusion+ chain.
func GetTrueERC20(chainID int) (common.Address, bool) {
	addr, ok := ChainToTrueERC20[chainID]
	return addr, ok
}
