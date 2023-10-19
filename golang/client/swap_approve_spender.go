package client

import (
	"context"
	"net/http"

	"dev-portal-sdk-go/client/swap"
)

func (c Client) ApproveSpender() (*swap.SpenderResponse, *http.Response, error) {
	u := "/swap/v5.2/1/approve/spender"

	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var spender swap.SpenderResponse
	res, err := c.Do(context.Background(), req, &spender)
	if err != nil {
		return nil, nil, err
	}

	return &spender, res, nil
}
