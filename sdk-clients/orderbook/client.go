package orderbook

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/1inch/1inch-sdk-go/common"
	"github.com/1inch/1inch-sdk-go/constants"
)

type Client struct {
	api
	Wallet              common.Wallet
	TxBuilder           common.TransactionBuilderFactory
	AggregationRouterV6 *abi.ABI
	SeriesNonceManager  *abi.ABI
}

type ClientOnlyAPI struct {
	api
}

type api struct {
	chainId      uint64
	httpExecutor common.HttpExecutor
}

func NewClient(cfg *Configuration) (*Client, error) {

	aggregationRouterV6, err := abi.JSON(strings.NewReader(constants.AggregationRouterV6ABI))
	if err != nil {
		return nil, err
	}

	seriesNonceManagerABI, err := abi.JSON(strings.NewReader(constants.SeriesNonceManagerABI))
	if err != nil {
		return nil, err
	}

	c := Client{
		api:                 cfg.APIConfiguration.API,
		AggregationRouterV6: &aggregationRouterV6,
		SeriesNonceManager:  &seriesNonceManagerABI,
	}

	if cfg.WalletConfiguration != nil {
		c.Wallet = cfg.WalletConfiguration.Wallet
		c.TxBuilder = cfg.WalletConfiguration.TxBuilder
	}

	return &c, nil
}

func NewClientOnlyAPI(cfg *ConfigurationAPI) (*ClientOnlyAPI, error) {
	c := ClientOnlyAPI{
		api: cfg.API,
	}

	return &c, nil
}
