package tokens

import (
	"github.com/1inch/1inch-sdk-go/common"
)

type Client struct {
	api
}

type api struct {
	chainId      uint64
	httpExecutor common.HttpExecutor
}

func NewClient(cfg *Configuration) (*Client, error) {
	c := Client{
		api: cfg.API,
	}
	return &c, nil
}
