package client

import (
	"context"
	"net/http"

	"dev-portal-sdk-go/client/swap"
)

func (s *SwapService) GetLiquiditySources(ctx context.Context) (*swap.ProtocolsResponse, *http.Response, error) {
	u := "/swap/v5.2/1/liquidity-sources"

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var liquiditySources swap.ProtocolsResponse
	res, err := s.client.Do(ctx, req, &liquiditySources)
	if err != nil {
		return nil, nil, err
	}

	return &liquiditySources, res, nil
}
