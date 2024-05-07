package aggregation

import (
	"github.com/1inch/1inch-sdk-go/common"
	"github.com/1inch/1inch-sdk-go/internal/http-executor"
	"github.com/1inch/1inch-sdk-go/internal/transaction-builder"
	"github.com/1inch/1inch-sdk-go/internal/web3-provider"
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
	NodeURL    string

	Wallet    common.Wallet
	TxBuilder common.TransactionBuilderFactory
}

func NewConfiguration(nodeUrl string, privateKey string, chainId uint64, apiUrl string, apiKey string) (*Configuration, error) {
	apiCfg, err := NewConfigurationAPI(chainId, apiUrl, apiKey)
	if err != nil {
		return nil, err
	}
	walletCfg, err := NewConfigurationWallet(nodeUrl, privateKey, chainId)
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

func NewConfigurationWallet(nodeUrl string, privateKey string, chainId uint64) (*ConfigurationWallet, error) {
	w, err := web3_provider.DefaultWalletProvider(privateKey, nodeUrl, chainId)
	if err != nil {
		return nil, err
	}

	f := transaction_builder.NewFactory(w)
	return &ConfigurationWallet{
		Wallet:    w,
		TxBuilder: f,
	}, nil
}
