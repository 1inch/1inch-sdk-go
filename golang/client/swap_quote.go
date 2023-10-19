package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"dev-portal-sdk-go/client/swap"
)

func (c Client) GetQuote(params swap.AggregationControllerGetQuoteParams) (*swap.QuoteResponse, *http.Response, error) {
	u := "/swap/v5.2/1/quote"

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

	var quote swap.QuoteResponse
	res, err := c.Do(context.Background(), req, &quote)
	if err != nil {
		return nil, nil, err
	}

	return &quote, res, nil
}

func getQuoteAddQueryParameters(query url.Values, params swap.AggregationControllerGetQuoteParams) url.Values {
	query.Add("src", params.Src)
	query.Add("dst", params.Dst)
	query.Add("amount", params.Amount)

	if params.Protocols != nil {
		query.Add("protocols", *params.Protocols)
	}
	if params.Fee != nil {
		query.Add("fee", fmt.Sprintf("%f", *params.Fee))
	}
	if params.GasPrice != nil {
		query.Add("gasPrice", *params.GasPrice)
	}
	if params.ComplexityLevel != nil {
		query.Add("complexityLevel", fmt.Sprintf("%f", *params.ComplexityLevel))
	}
	if params.Parts != nil {
		query.Add("parts", fmt.Sprintf("%f", *params.Parts))
	}
	if params.MainRouteParts != nil {
		query.Add("mainRouteParts", fmt.Sprintf("%f", *params.MainRouteParts))
	}
	if params.GasLimit != nil {
		query.Add("gasLimit", fmt.Sprintf("%f", *params.GasLimit))
	}
	if params.IncludeTokensInfo != nil {
		query.Add("includeTokensInfo", fmt.Sprintf("%v", *params.IncludeTokensInfo))
	}
	if params.IncludeProtocols != nil {
		query.Add("includeProtocols", fmt.Sprintf("%v", *params.IncludeProtocols))
	}
	if params.IncludeGas != nil {
		query.Add("includeGas", fmt.Sprintf("%v", *params.IncludeGas))
	}
	if params.ConnectorTokens != nil {
		query.Add("connectorTokens", *params.ConnectorTokens)
	}

	return query
}
