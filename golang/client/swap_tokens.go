package client

import (
	"context"
	"net/http"

	"dev-portal-sdk-go/client/swap"
)

func (c Client) GetTokens() (*swap.TokensResponse, *http.Response, error) {
	u := "/swap/v5.2/1/tokens"

	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var tokens swap.TokensResponse
	res, err := c.Do(context.Background(), req, &tokens)
	if err != nil {
		return nil, nil, err
	}

	return &tokens, res, nil
}
