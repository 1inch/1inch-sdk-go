package aggregation

import (
	"context"

	"github.com/1inch/1inch-sdk-go/aggregation/models"
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

	ctx := context.Background()

	swapData, err := client.GetSwap(ctx, models.AggregationControllerGetSwapParams{})
	if err != nil {
		return
	}

	nonce, err := client.Wallet.Nonce(ctx)
	if err != nil {
		return
	}

	gasTip, err := client.Wallet.GetGasTipCap(ctx)
	if err != nil {
		return
	}

	gasFee, err := client.Wallet.GetGasFeeCap(ctx)
	if err != nil {
		return
	}

	tx, err := client.BuildSwapTransaction(swapData, nonce, gasTip, gasFee)
	if err != nil {
		return
	}

	signedTx, err := client.Wallet.Sign(tx)
	if err != nil {
		return
	}

	err = client.Wallet.BroadcastTransaction(ctx, signedTx)
	if err != nil {
		return
	}
}
