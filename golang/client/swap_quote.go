package client

import (
	"context"
	"fmt"
	"net/http"

	"dev-portal-sdk-go/client/tokenprices"
)

func (c Client) GetQuote(params tokenprices.PricesParameters) (string, *http.Response, error) {
	// TODO accept context
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/price/v1.1/1", c.BaseURL), nil)

	err = validateParameters(params)
	if err != nil {
		return "", nil, err
	}

	switch params.Currency {
	case "":
	case "WEI":
	default:
		req.URL.RawQuery += fmt.Sprintf("currency=%s", params.Currency)
	}
	if err != nil {
		return "", nil, err
	}

	var exStr PricesResponse
	res, err := c.Do(context.Background(), req, &exStr)
	if err != nil {
		return "", nil, err
	}

	return fmt.Sprintf("Pirces: %v", exStr), res, nil
}
