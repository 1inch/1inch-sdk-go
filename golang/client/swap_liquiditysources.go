package client

import (
	"context"
	"net/http"

	"dev-portal-sdk-go/client/swap"
)

func (c Client) GetLiquiditySources() (*swap.ProtocolsResponse, *http.Response, error) {
	u := "/swap/v5.2/1/liquidity-sources"

	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var liquiditySources swap.ProtocolsResponse
	res, err := c.Do(context.Background(), req, &liquiditySources)
	if err != nil {
		return nil, nil, err
	}

	return &liquiditySources, res, nil
}
