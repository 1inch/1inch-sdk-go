package orderbook

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/1inch/1inch-sdk-go/common"
	"github.com/1inch/1inch-sdk-go/constants"

	http_executor "github.com/1inch/1inch-sdk-go/internal/http-executor"
	transaction_builder "github.com/1inch/1inch-sdk-go/internal/transaction-builder"
	web3_provider "github.com/1inch/1inch-sdk-go/internal/web3-provider"
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

	Wallet    common.Wallet
	TxBuilder common.TransactionBuilderFactory
}

type Client struct {
	api
	AggregationRouterV6 *abi.ABI
	SeriesNonceManager  *abi.ABI
	Wallet              common.Wallet
	TxBuilder           common.TransactionBuilderFactory
}

type api struct {
	chainId      uint64
	httpExecutor common.HttpExecutor
}

func NewClient(cfg *Configuration) (*Client, error) {

	aggregationRouterV6, err := abi.JSON(strings.NewReader(constants.AggregationRouterV6ABI))
	if err != nil {
		return nil, err
	}

	seriesNonceManagerABI, err := abi.JSON(strings.NewReader(constants.SeriesNonceManagerABI)) // Make a generic version of this ABI
	if err != nil {
		return nil, err
	}

	c := Client{
		api:                 cfg.API,
		AggregationRouterV6: &aggregationRouterV6,
		SeriesNonceManager:  &seriesNonceManagerABI,
	}

	if cfg.WalletConfiguration != nil {
		c.Wallet = cfg.WalletConfiguration.Wallet
		c.TxBuilder = cfg.WalletConfiguration.TxBuilder
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

	f := transaction_builder.NewFactory(w)
	return &WalletConfiguration{
		Wallet:    w,
		TxBuilder: f,
	}, nil
}
