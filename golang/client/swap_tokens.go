package client

import (
	"context"
	"net/http"

	"dev-portal-sdk-go/client/swap"
)

func (s *SwapService) GetTokens(ctx context.Context) (*swap.TokensResponse, *http.Response, error) {
	u := "/swap/v5.2/1/tokens"

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var tokens swap.TokensResponse
	res, err := s.client.Do(ctx, req, &tokens)
	if err != nil {
		return nil, nil, err
	}

	return &tokens, res, nil
}
