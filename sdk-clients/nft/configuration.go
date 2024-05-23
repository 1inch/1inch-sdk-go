package nft

import (
	http_executor "github.com/1inch/1inch-sdk-go/internal/http-executor"
)

type Configuration struct {
	ApiKey string
	ApiURL string
	API    api
}

type ConfigurationParams struct {
	ApiUrl string
	ApiKey string
}

func NewConfiguration(params ConfigurationParams) (*Configuration, error) {
	executor, err := http_executor.DefaultHttpClient(params.ApiUrl, params.ApiKey)
	if err != nil {
		return nil, err
	}

	a := api{
		httpExecutor: executor,
	}

	return &Configuration{
		ApiURL: params.ApiUrl,
		ApiKey: params.ApiKey,
		API:    a,
	}, nil
}
