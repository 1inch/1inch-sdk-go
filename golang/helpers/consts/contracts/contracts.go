package contracts

import (
	"fmt"

	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
)

const AggregationRouterV5 = "0x1111111254eeb25477b68fb85ed929f73a960582"
const AggregationV5RouterZkSyncEra = "0x6e2B76966cbD9cF4cC2Fa0D76d24d5241E0ABC2F"
const AggregationRouterV5Name = "1inch Aggregation Router"
const AggregationRouterV5VersionNumber = "5"

const SeriesNonceManager = "0x303389f541ff2d620e42832f180a08e767b28e10"
const SeriesNonceManagerPolygon = "0xa5eb255EF45dFb48B5d133d08833DEF69871691D"
const SeriesNonceManagerBsc = "0x58ce0e6ef670c9a05622f4188faa03a9e12ee2e4"

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

// TODO add nonce manager contracts for all supported chains

func GetSeriesNonceManagerFromChainId(chainId int) (string, error) {
	switch chainId {
	case chains.Ethereum:
		return SeriesNonceManager, nil
	case chains.Polygon:
		return SeriesNonceManagerPolygon, nil
	case chains.Bsc:
		return SeriesNonceManagerBsc, nil
	default:
		return "", fmt.Errorf("unrecognized chain id: %d", chainId)
	}
}
