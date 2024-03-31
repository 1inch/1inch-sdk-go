package aggregation

import (
	"github.com/1inch/1inch-sdk-go/internal/common"
	"github.com/1inch/1inch-sdk-go/internal/http_executor"
	"github.com/1inch/1inch-sdk-go/internal/web3_provider"
)

type api struct {
	httpExecutor common.HttpExecutor
}

type Client struct {
	api
	Wallet common.Wallet
}

func NewClient(cfg *Configuration) (*Client, error) {
	executor := http_executor.DefaultHttpClient(cfg.ApiURL, cfg.ApiKey)
	api := api{
		httpExecutor: &executor,
	}

	c := Client{
		api: api,
	}

	if cfg.WalletConfig != nil {
		w, err := web3_provider.DefaultWalletProvider(cfg.WalletConfig.PrivateKey, cfg.WalletConfig.NodeURL, cfg.ChainId)
		if err != nil {
			return nil, err
		}
		c.Wallet = w
	}

	return &c, nil
}
