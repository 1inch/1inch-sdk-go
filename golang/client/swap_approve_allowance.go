package client

import (
	"context"
	"fmt"
	"net/http"

	"1inch-sdk-golang/client/swap"
)

func (s *SwapService) ApproveAllowance(ctx context.Context, params swap.ApproveControllerGetAllowanceParams) (*swap.AllowanceResponse, *http.Response, error) {
	u := "/swap/v5.2/1/approve/allowance"

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

	var allowanceResponse swap.AllowanceResponse
	res, err := s.client.Do(ctx, req, &allowanceResponse)
	if err != nil {
		return nil, nil, err
	}

	return &allowanceResponse, res, nil
}