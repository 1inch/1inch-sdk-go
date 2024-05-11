package spotprices

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/1inch/1inch-sdk-go/common"
)

// GetPricesForWhitelistedTokens Get Prices for whitelisted tokens
func (api *api) GetPricesForWhitelistedTokens(ctx context.Context, params GetWhitelistedTokensPricesParams) (*PricesForWhitelistedTokensResponse, error) {
	u := fmt.Sprintf("/price/v1.1/%d", api.chainId)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		U:      u,
		Body:   nil,
	}

	var response PricesForWhitelistedTokensResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetCustomCurrenciesList  Get List of custom currencies
func (api *api) GetCustomCurrenciesList(ctx context.Context) (*CurrenciesResponseDto, error) {
	u := fmt.Sprintf("/price/v1.1/%d/currencies", api.chainId)

	payload := common.RequestPayload{
		Method: "GET",
		Params: nil,
		U:      u,
		Body:   nil,
	}

	var response CurrenciesResponseDto
	err := api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetPricesForRequestedTokens Get prices for requested tokens (do not support too many tokens since it is HTTP GET)
func (api *api) GetPricesForRequestedTokens(ctx context.Context, params GetPricesRequestDto) (*PricesForRequestedTokensResponse, error) {
	u := fmt.Sprintf("/price/v1.1/%d/%s", api.chainId, strings.Join(params.Tokens, ","))

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: struct {
			Currency GetPricesRequestDtoCurrency `json:"currency"`
		}{
			Currency: params.Currency,
		},
		U:    u,
		Body: nil,
	}

	var response PricesForRequestedTokensResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetPricesForRequestedTokensLarge Get prices for requested tokens
func (api *api) GetPricesForRequestedTokensLarge(ctx context.Context, params GetPricesRequestDto) (*PricesForRequestedTokensResponse, error) {
	u := fmt.Sprintf("/price/v1.1/%d", api.chainId)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "POST",
		Params: nil,
		U:      u,
		Body:   body,
	}

	var response PricesForRequestedTokensResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
