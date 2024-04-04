package aggregation

import (
	"github.com/1inch/1inch-sdk-go/common"

	"github.com/1inch/1inch-sdk-go/internal/http_executor"
	"github.com/1inch/1inch-sdk-go/internal/web3_provider"
)

type Configuration struct {
	WalletConfiguration *WalletConfiguration
	ChainId             uint64

	ApiKey string
	ApiURL string

	API api
}

type WalletConfiguration struct {
	PrivateKey string
	NodeURL    string

	Wallet common.Wallet
}

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
		api: cfg.API,
	}

	if cfg.WalletConfiguration != nil {
		c.Wallet = cfg.WalletConfiguration.Wallet
	}

	return &c, nil
}

func NewDefaultConfiguration(nodeUrl string, privateKey string, chainId uint64, apiUrl string, apiKey string) (*Configuration, error) {
	executor, err := http_executor.DefaultHttpClient(apiUrl, apiKey)
	if err != nil {
		return nil, err
	}

	a := api{
		chainId:      chainId,
		httpExecutor: executor,
	}

	walletCfg, err := NewDefaultWalletConfiguration(nodeUrl, privateKey, chainId)
	if err != nil {
		return nil, err
	}
	return &Configuration{
		WalletConfiguration: walletCfg,
		API:                 a,
	}, nil
}

func NewDefaultWalletConfiguration(nodeUrl string, privateKey string, chainId uint64) (*WalletConfiguration, error) {
	w, err := web3_provider.DefaultWalletProvider(privateKey, nodeUrl, chainId)
	if err != nil {
		return nil, err
	}

	return &WalletConfiguration{
		Wallet: w,
	}, nil
}
