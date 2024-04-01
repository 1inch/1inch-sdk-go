package aggregation

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/1inch/1inch-sdk-go/aggregation/models"
	"github.com/1inch/1inch-sdk-go/internal/common"
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

func DefaultConfiguration(nodeUrl string, privateKey string, chainId uint64, apiUrl string, apiKey string) (*Configuration, error) {
	executor, err := http_executor.DefaultHttpClient(apiUrl, apiKey)
	if err != nil {
		return nil, err
	}

	a := api{
		httpExecutor: executor,
	}

	walletCfg, err := DefaultWalletConfiguration(nodeUrl, privateKey, chainId)
	if err != nil {
		return nil, err
	}
	return &Configuration{
		WalletConfiguration: walletCfg,
		API:                 a,
	}, nil
}

func DefaultWalletConfiguration(nodeUrl string, privateKey string, chainId uint64) (*WalletConfiguration, error) {
	w, err := web3_provider.DefaultWalletProvider(privateKey, nodeUrl, chainId)
	if err != nil {
		return nil, err
	}

	return &WalletConfiguration{
		Wallet: w,
	}, nil
}

func testSDK() {
	config, err := DefaultConfiguration("", "", 1, "", "")
	if err != nil {
		return
	}
	client, err := NewClient(config)

	swapData, err := client.GetSwap(context.Background(), models.GetSwapParams{
		ChainId:                            1,
		SkipWarnings:                       false,
		AggregationControllerGetSwapParams: models.AggregationControllerGetSwapParams{},
	})
	if err != nil {
		return
	}

	signedTx, err := client.Wallet.Sign(&types.Transaction{})
	if err != nil {
		return
	}

	err = client.Wallet.BroadcastTransaction(context.Background(), signedTx)
	if err != nil {
		return
	}
}
