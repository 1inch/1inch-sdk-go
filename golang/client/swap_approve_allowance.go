package client

import (
	"context"
	"fmt"
	"net/http"

	"dev-portal-sdk-go/client/swap"
)

func (c Client) ApproveAllowance(params swap.ApproveControllerGetAllowanceParams) (*swap.AllowanceResponse, *http.Response, error) {
	u := "/swap/v5.2/1/approve/allowance"

	err := params.Validate()
	if err != nil {
		return nil, nil, fmt.Errorf("request validation error: %v", err)
	}

	u, err = addOptions(u, params)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var allowanceResponse swap.AllowanceResponse
	res, err := c.Do(context.Background(), req, &allowanceResponse)
	if err != nil {
		return nil, nil, err
	}

	return &allowanceResponse, res, nil
}
