package orderbook

import (
	"fmt"

	"github.com/1inch/1inch-sdk-go/common"
	http_executor "github.com/1inch/1inch-sdk-go/internal/http-executor"
	web3_provider "github.com/1inch/1inch-sdk-go/internal/web3-provider"
)

type Configuration struct {
	WalletConfiguration *WalletConfiguration
	APIConfiguration    *ConfigurationAPI
}

type ConfigurationAPI struct {
	ApiKey string
	ApiURL string

	API api
}

type WalletConfiguration struct {
	PrivateKey string
	NodeURL    string

	Wallet    common.Wallet
	TxBuilder common.TransactionBuilderFactory
}

type ConfigurationParams struct {
	NodeUrl    string
	PrivateKey string
	ChainId    uint64
	ApiUrl     string
	ApiKey     string
}

func NewConfiguration(params ConfigurationParams) (*Configuration, error) {
	apiCfg, err := NewConfigurationAPI(params.ChainId, params.ApiUrl, params.ApiKey)
	if err != nil {
		return nil, err
	}
	walletCfg, err := NewConfigurationWallet(params.PrivateKey)
	if err != nil {
		return nil, err
	}
	return &Configuration{
		WalletConfiguration: walletCfg,
		APIConfiguration:    apiCfg,
	}, nil
}

func NewConfigurationAPI(chainId uint64, apiUrl string, apiKey string) (*ConfigurationAPI, error) {
	executor, err := http_executor.DefaultHttpClient(apiUrl, apiKey)
	if err != nil {
		return nil, err
	}

	a := api{
		chainId:      chainId,
		httpExecutor: executor,
	}

	return &ConfigurationAPI{
		ApiURL: apiUrl,
		ApiKey: apiKey,

		API: a,
	}, nil
}

func NewConfigurationWallet(privateKey string) (*WalletConfiguration, error) {
	w, err := web3_provider.DefaultWalletOnlyProvider(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %v", err)
	}
	return &WalletConfiguration{
		Wallet: w,
	}, nil
}
