package fusion

import (
	"fmt"

	"github.com/1inch/1inch-sdk-go/common"
	http_executor "github.com/1inch/1inch-sdk-go/internal/http-executor"
	web3_provider "github.com/1inch/1inch-sdk-go/internal/web3-provider"
)

type Configuration struct {
	WalletConfiguration *ConfigurationWallet
	APIConfiguration    *ConfigurationAPI
}

type ConfigurationAPI struct {
	ApiKey string
	ApiURL string

	API api
}

type ConfigurationWallet struct {
	PrivateKey string
	Wallet     common.Wallet
}

type ConfigurationParams struct {
	ChainId    uint64
	ApiUrl     string
	ApiKey     string
	PrivateKey string
}

func NewConfiguration(params ConfigurationParams) (*Configuration, error) {
	executor, err := http_executor.DefaultHttpClient(params.ApiUrl, params.ApiKey)
	if err != nil {
		return nil, err
	}

	a := api{
		chainId:      params.ChainId,
		httpExecutor: executor,
	}

	walletCfg, err := NewConfigurationWallet(params.PrivateKey, params.ChainId)
	if err != nil {
		return nil, err
	}

	return &Configuration{
		WalletConfiguration: walletCfg,
		APIConfiguration: &ConfigurationAPI{
			ApiURL: params.ApiUrl,
			ApiKey: params.ApiKey,
			API:    a,
		},
	}, nil
}

func NewConfigurationWallet(privateKey string, chainId uint64) (*ConfigurationWallet, error) {
	if privateKey == "" {
		return nil, fmt.Errorf("private key is required")
	}
	w, err := web3_provider.DefaultWalletOnlyProvider(privateKey, chainId)
	if err != nil {
		return nil, err
	}
	return &ConfigurationWallet{
		Wallet: w,
	}, nil
}
