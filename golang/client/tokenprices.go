package client

import (
	"context"
	"fmt"
	"net/http"

	"dev-portal-sdk-go/client/spotprice"
)

type PricesResponse map[string]string

func (c Client) GetTokenPrices(params spotprice.ChainControllerByAddressesParams) (*PricesResponse, *http.Response, error) {
	// TODO accept context
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/price/v1.1/1", c.BaseURL), nil)

	err = params.Validate()
	if err != nil {
		return nil, nil, err
	}

	if params.Currency != nil && *params.Currency != "" {
		req.URL.RawQuery += fmt.Sprintf("currency=%s", string(*params.Currency))
	}
	if err != nil {
		return nil, nil, err
	}

	var exStr PricesResponse
	res, err := c.Do(context.Background(), req, &exStr)
	if err != nil {
		return nil, nil, err
	}

	return &exStr, res, nil
}
