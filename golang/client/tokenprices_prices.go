package client

import (
	"context"
	"net/http"

	"dev-portal-sdk-go/client/tokenprices"
)

type PricesResponse map[string]string

func (tp *TokenPricesService) GetPrices(ctx context.Context, params tokenprices.ChainControllerByAddressesParams) (*PricesResponse, *http.Response, error) {
	u := "/price/v1.1/1"

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	// If nothing is set, remove the field from the struct
	// A blank input is required by the API to set the response currency to wei
	if params.Currency != nil && *params.Currency == "" {
		params.Currency = nil
	}

	u, err = addOptions(u, params)
	if err != nil {
		return nil, nil, err
	}

	req, err := tp.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var exStr PricesResponse
	res, err := tp.client.Do(ctx, req, &exStr)
	if err != nil {
		return nil, nil, err
	}

	return &exStr, res, nil
}
