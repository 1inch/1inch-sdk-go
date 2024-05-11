package spotprices

import (
	"github.com/1inch/1inch-sdk-go/internal/http-executor"
)

type Configuration struct {
	ApiKey string
	ApiURL string
	API    api
}

func NewConfiguration(chainId uint64, apiUrl string, apiKey string) (*Configuration, error) {
	executor, err := http_executor.DefaultHttpClient(apiUrl, apiKey)
	if err != nil {
		return nil, err
	}

	a := api{
		chainId:      chainId,
		httpExecutor: executor,
	}

	return &Configuration{
		ApiURL: apiUrl,
		ApiKey: apiKey,
		API:    a,
	}, nil
}
