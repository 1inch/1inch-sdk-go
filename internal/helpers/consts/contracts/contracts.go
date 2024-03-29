package contracts

import (
	"fmt"

	"github.com/1inch/1inch-sdk-go/internal/helpers"
	"github.com/1inch/1inch-sdk-go/internal/helpers/consts/chains"
)

const AggregationRouterV5 = "0x1111111254eeb25477b68fb85ed929f73a960582" // Contract address is identical for all chains except zkSync
const AggregationV5RouterZkSyncEra = "0x6e2B76966cbD9cF4cC2Fa0D76d24d5241E0ABC2F"
const AggregationRouterV5Name = "1inch Aggregation Router"
const AggregationRouterV5VersionNumber = "5"

// Series Nonce Manager contract addresses are taken from limit-order-protocol/deployments

const SeriesNonceManagerArbitrum = "0xD7936052D1e096d48C81Ef3918F9Fd6384108480"
const SeriesNonceManagerAurora = "0x7F069df72b7A39bCE9806e3AfaF579E54D8CF2b9"
const SeriesNonceManagerAvalanche = "0x2EC255797FEF7669fA243509b7a599121148FFba"
const SeriesNonceManagerBase = "0xD9Cc0A957cAC93135596f98c20Fbaca8Bf515909"
const SeriesNonceManagerBsc = "0x58ce0e6ef670c9a05622f4188faa03a9e12ee2e4"
const SeriesNonceManagerEthereum = "0x303389f541ff2d620e42832f180a08e767b28e10"
const SeriesNonceManagerFantom = "0x7871769b3816b23dB12E83a482aAc35F1FD35D4B"
const SeriesNonceManagerGnosis = "0x11431a89893025D2a48dCA4EddC396f8C8117187"
const SeriesNonceManagerKlaytn = "0x7871769b3816b23dB12E83a482aAc35F1FD35D4B"
const SeriesNonceManagerOptimism = "0x32d12a25f539E341089050E2d26794F041fC9dF8"
const SeriesNonceManagerPolygon = "0xa5eb255EF45dFb48B5d133d08833DEF69871691D"

func Get1inchRouterFromChainId(chainId int) (string, error) {
	if helpers.Contains(chainId, chains.ValidChainIds) {
		if chainId == chains.ZkSyncEra {
			return AggregationV5RouterZkSyncEra, nil
		} else {
			return AggregationRouterV5, nil
		}
	} else {
		return "", fmt.Errorf("unrecognized chain id: %d", chainId)
	}
}

func GetSeriesNonceManagerFromChainId(chainId int) (string, error) {
	switch chainId {
	case chains.Arbitrum:
		return SeriesNonceManagerArbitrum, nil
	case chains.Aurora:
		return SeriesNonceManagerAurora, nil
	case chains.Avalanche:
		return SeriesNonceManagerAvalanche, nil
	case chains.Base:
		return SeriesNonceManagerBase, nil
	case chains.Bsc:
		return SeriesNonceManagerBsc, nil
	case chains.Ethereum:
		return SeriesNonceManagerEthereum, nil
	case chains.Fantom:
		return SeriesNonceManagerFantom, nil
	case chains.Gnosis:
		return SeriesNonceManagerGnosis, nil
	case chains.Klaytn:
		return SeriesNonceManagerKlaytn, nil
	case chains.Optimism:
		return SeriesNonceManagerOptimism, nil
	case chains.Polygon:
		return SeriesNonceManagerPolygon, nil
	case chains.ZkSyncEra:
		return "", fmt.Errorf("zksync contract unknown") // TODO get this contract
	default:
		return "", fmt.Errorf("unrecognized chain id: %d", chainId)
	}
}
