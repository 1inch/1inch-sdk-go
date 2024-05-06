package aggregation

import (
	"github.com/1inch/1inch-sdk-go/common"
)

type Client struct {
	api
	Wallet    common.Wallet
	TxBuilder common.TransactionBuilderFactory
}

type ClientOnlyAPI struct {
	api
}

type api struct {
	chainId      uint64
	httpExecutor common.HttpExecutor
}

func NewClient(cfg *Configuration) (*Client, error) {
	c := Client{
		api: cfg.APIConfiguration.API,
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
