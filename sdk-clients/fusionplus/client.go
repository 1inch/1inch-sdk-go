package fusionplus

import (
	"github.com/1inch/1inch-sdk-go/common"
)

type Client struct {
	api
	Wallet common.Wallet
}

type api struct {
	chainId      uint64
	httpExecutor common.HttpExecutor
}

func NewClient(cfg *Configuration) (*Client, error) {
	c := Client{
		api:    cfg.APIConfiguration.API,
		Wallet: cfg.WalletConfiguration.Wallet,
	}
	return &c, nil
}
