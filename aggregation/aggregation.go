package aggregation

import (
	"net/url"

	"github.com/1inch/1inch-sdk-go/internal/common"
	"github.com/1inch/1inch-sdk-go/internal/http_executor"
	"github.com/1inch/1inch-sdk-go/internal/web3_provider"
)

type api struct {
	httpExecutor common.HttpExecutor
}

type Client struct {
	api
	Wallet web3_provider.Wallet
}

// todo: not done
func DefaultClient() *Client {
	// todo: move to input params, that will be validated before
	u, _ := url.Parse("https://api.1inch.dev")
	executor := http_executor.DefaultHttpClient(u, "")
	api := api{
		httpExecutor: &executor,
	}
	c := Client{
		api,
	}
	return &c
}

func example() {

	// allowance()
	// approve()
	// sign()
	// broadcast()

	// permit()

	// swapData()
	// sign()
	// broadcast()
}

// manipulate txs
// each action -> return
