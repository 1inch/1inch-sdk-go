package constants

import (
	"fmt"
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

func GetSeriesNonceManagerFromChainId(chainId int) (string, error) {
	switch chainId {
	case ArbitrumChainId:
		return SeriesNonceManagerArbitrum, nil
	case AuroraChainId:
		return SeriesNonceManagerAurora, nil
	case AvalancheChainId:
		return SeriesNonceManagerAvalanche, nil
	case BaseChainId:
		return SeriesNonceManagerBase, nil
	case BscChainId:
		return SeriesNonceManagerBsc, nil
	case EthereumChainId:
		return SeriesNonceManagerEthereum, nil
	case FantomChainId:
		return SeriesNonceManagerFantom, nil
	case GnosisChainId:
		return SeriesNonceManagerGnosis, nil
	case KlaytnChainId:
		return SeriesNonceManagerKlaytn, nil
	case OptimismChainId:
		return SeriesNonceManagerOptimism, nil
	case PolygonChainId:
		return SeriesNonceManagerPolygon, nil
	case ZkSyncEraChainId:
		return "", fmt.Errorf("zksync contract unknown") // TODO get this contract
	default:
		return "", fmt.Errorf("unrecognized chain id: %d", chainId)
	}
}
