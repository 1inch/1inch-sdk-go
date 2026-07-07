package fusion

import (
	"context"

	"github.com/1inch/1inch-sdk-go/common"
)

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
		api:    cfg.APIConfiguration.API,
		Wallet: cfg.WalletConfiguration.Wallet,
	}
	return &c, nil
}

// PlaceOrderFromParams fetches a quote for the given order params and places the
// resulting order in one call, so settings like Permit and IsPermit2 are supplied
// once and propagate to both the quote request and the order
func (c *Client) PlaceOrderFromParams(ctx context.Context, orderParams OrderParams) (string, error) {
	isPermit2 := ""
	if orderParams.IsPermit2 {
		isPermit2 = "true"
	}
	quote, err := c.GetQuote(ctx, QuoterControllerGetQuoteParamsFixed{
		FromTokenAddress: orderParams.FromTokenAddress,
		ToTokenAddress:   orderParams.ToTokenAddress,
		Amount:           orderParams.Amount,
		WalletAddress:    orderParams.WalletAddress,
		EnableEstimate:   true,
		IsPermit2:        isPermit2,
		Permit:           orderParams.Permit,
		Surplus:          true,
	})
	if err != nil {
		return "", err
	}
	return c.PlaceOrder(ctx, *quote, orderParams, c.Wallet)
}
