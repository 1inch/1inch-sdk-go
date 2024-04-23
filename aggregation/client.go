package aggregation

import (
	"github.com/1inch/1inch-sdk-go/common"
)

type Client struct {
	api
	Wallet    common.Wallet
	TxBuilder common.TransactionBuilderFactory
}

type api struct {
	chainId      uint64
	httpExecutor common.HttpExecutor
}

func NewClient(cfg *Configuration) (*Client, error) {
	c := Client{
		api: cfg.API,
	}

	if cfg.WalletConfiguration != nil {
		c.Wallet = cfg.WalletConfiguration.Wallet
		c.TxBuilder = cfg.WalletConfiguration.TxBuilder
	}

	return &c, nil
}
