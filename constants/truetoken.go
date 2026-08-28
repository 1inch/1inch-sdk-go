package constants

import "github.com/ethereum/go-ethereum/common"

// ChainToTrueERC20 maps a chain id to its Fusion+ TRUE_ERC20 sentinel token.
//
// The taker asset of a Fusion+ order on the source chain must be this sentinel.
// The sentinel is a token that does nothing when a resolver transfers it. It is
// not the real destination token. The escrow extension carries the real
// destination token, and the swap delivers it on the destination chain.
//
// Do not put the destination token in the source-chain taker asset. The
// destination token can also be a live ERC-20 on the source chain. In that
// condition, the transfer moves a real and unwanted token.
//
// A chain that is not in this map is not a Fusion+ chain. For such a chain,
// order construction fails with an error. It does not sign an order that has
// the wrong taker asset.
//
// Source: @1inch/cross-chain-sdk deployments.ts (TRUE_ERC20 map), verified
// 2026-08-28. All supported chains use one deterministic address, except
// zkSync Era.
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

// GetTrueERC20 returns the Fusion+ TRUE_ERC20 sentinel for a chain. The second
// result is true when the chain is a Fusion+ chain.
func GetTrueERC20(chainID int) (common.Address, bool) {
	addr, ok := ChainToTrueERC20[chainID]
	return addr, ok
}
