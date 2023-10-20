package client

import (
	"context"
	"fmt"
	"net/http"

	"dev-portal-sdk-go/client/swap"
)

func (s *SwapService) GetQuote(params swap.AggregationControllerGetQuoteParams) (*swap.QuoteResponse, *http.Response, error) {
	u := "/swap/v5.2/1/quote"

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

	var quote swap.QuoteResponse
	res, err := s.client.Do(context.Background(), req, &quote)
	if err != nil {
		return nil, nil, err
	}

	return &quote, res, nil
}
