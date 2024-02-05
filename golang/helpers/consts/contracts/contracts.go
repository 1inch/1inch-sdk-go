package contracts

import (
	"fmt"

	"github.com/svanas/1inch-sdk/golang/helpers"
	"github.com/svanas/1inch-sdk/golang/helpers/consts/chains"
)

const AggregationRouterV5 = "0x1111111254eeb25477b68fb85ed929f73a960582"
const AggregationV5RouterZkSyncEra = "0x6e2B76966cbD9cF4cC2Fa0D76d24d5241E0ABC2F"
const AggregationRouterV5Name = "1inch Aggregation Router"
const AggregationRouterV5VersionNumber = "5"

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
