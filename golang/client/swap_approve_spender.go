package client

import (
	"context"
	"net/http"

	"1inch-sdk-golang/client/swap"
)

func (s *SwapService) ApproveSpender(ctx context.Context) (*swap.SpenderResponse, *http.Response, error) {
	u := "/swap/v5.2/1/approve/spender"

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var spender swap.SpenderResponse
	res, err := s.client.Do(ctx, req, &spender)
	if err != nil {
		return nil, nil, err
	}

	return &spender, res, nil
}
