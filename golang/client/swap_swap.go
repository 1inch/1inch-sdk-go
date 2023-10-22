package client

import (
	"context"
	"fmt"
	"net/http"

	"1inch-sdk-golang/client/swap"
)

func (s *SwapService) GetSwap(ctx context.Context, params swap.AggregationControllerGetSwapParams) (*swap.SwapResponse, *http.Response, error) {
	u := "/swap/v5.2/1/swap"

	err := params.Validate()
	if err != nil {
		return nil, nil, fmt.Errorf("request validation error: %v", err)
	}

	u, err = addOptions(u, params)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var swap swap.SwapResponse
	res, err := s.client.Do(ctx, req, &swap)
	if err != nil {
		return nil, nil, err
	}

	return &swap, res, nil
}
