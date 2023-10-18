package client

import (
	"context"
	"fmt"
	"net/http"

	"dev-portal-sdk-go/client/spotprice"
)

type PricesResponse map[string]string

func (c Client) GetTokenPrices(params spotprice.PricesParameters) (*PricesResponse, *http.Response, error) {
	// TODO accept context
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/price/v1.1/1", c.BaseURL), nil)

	err = validateParameters(params)
	if err != nil {
		return nil, nil, err
	}

	switch params.Currency {
	case "":
	case "WEI":
	default:
		req.URL.RawQuery += fmt.Sprintf("currency=%s", params.Currency)
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

func validateParameters(params spotprice.PricesParameters) error {
	if !contains(spotprice.CurrencyTypeValues, params.Currency) {
		return fmt.Errorf("currency value %s is not valid", params.Currency)
	}
	return nil
}

// TODO Make a helpers class
func contains(slice []spotprice.CurrencyType, item spotprice.CurrencyType) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
