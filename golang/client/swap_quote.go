package client

import (
	"context"
	"fmt"
	"net/http"

	"dev-portal-sdk-go/client/swap"
)

func (c Client) GetQuote(params swap.AggregationControllerGetQuoteParams) (string, *http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/swap/v5.2/1/quote", c.BaseURL), nil)

	err = params.Validate()
	if err != nil {
		return "", nil, err
	}

	query := req.URL.Query()
	query.Add("src", params.Src)
	query.Add("dst", params.Dst)
	query.Add("amount", params.Amount)
	req.URL.RawQuery = query.Encode()

	var quote swap.QuoteResponse
	res, err := c.Do(context.Background(), req, &quote)
	if err != nil {
		return "", nil, err
	}

	return fmt.Sprintf("Quote: %v", quote), res, nil
}
