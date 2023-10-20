package client

import (
	"context"
	"fmt"
	"net/http"

	"dev-portal-sdk-go/client/swap"
)

func (s *SwapService) GetSwap(params swap.AggregationControllerGetSwapParams) (*swap.SwapResponse, *http.Response, error) {
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
	res, err := s.client.Do(context.Background(), req, &swap)
	if err != nil {
		return nil, nil, err
	}

	return &swap, res, nil
}
