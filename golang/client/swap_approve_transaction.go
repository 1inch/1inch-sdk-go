package client

import (
	"context"
	"fmt"
	"net/http"

	"dev-portal-sdk-go/client/swap"
)

func (s *SwapService) ApproveTransaction(params swap.ApproveControllerGetCallDataParams) (*swap.ApproveCallDataResponse, *http.Response, error) {
	u := "/swap/v5.2/1/approve/transaction"

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

	var approveCallData swap.ApproveCallDataResponse
	res, err := s.client.Do(context.Background(), req, &approveCallData)
	if err != nil {
		return nil, nil, err
	}

	return &approveCallData, res, nil
}
