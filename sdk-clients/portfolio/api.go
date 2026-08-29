package portfolio

import (
	"context"

	"github.com/1inch/1inch-sdk-go/v4/common"
)

func (api *api) GetProtocolsCurrentValue(ctx context.Context, params GetProtocolsCurrentValueParams) (*GetProtocolsCurrentValueResponse, error) {
	u := "/portfolio/portfolio/v4/overview/protocols/current_value"

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		Path:   u,
		Body:   nil,
	}

	var response GetProtocolsCurrentValueResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetProtocolsProfitAndLoss(ctx context.Context, params GetProtocolsProfitAndLossParams) (*GetProtocolsProfitAndLossResponse, error) {
	u := "/portfolio/portfolio/v4/overview/protocols/profit_and_loss"

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		Path:   u,
		Body:   nil,
	}

	var response GetProtocolsProfitAndLossResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetProtocolsDetails(ctx context.Context, params GetProtocolsDetailsParams) (*GetProtocolsDetailsResponse, error) {
	u := "/portfolio/portfolio/v4/overview/protocols/details"

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		Path:   u,
		Body:   nil,
	}

	var response GetProtocolsDetailsResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetTokensCurrentValue(ctx context.Context, params GetTokensCurrentValueParams) (*GetTokensCurrentValueResponse, error) {
	u := "/portfolio/portfolio/v4/overview/erc20/current_value"

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		Path:   u,
		Body:   nil,
	}

	var response GetTokensCurrentValueResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetTokensProfitAndLoss(ctx context.Context, params GetTokensProfitAndLossParams) (*GetTokensProfitAndLossResponse, error) {
	u := "/portfolio/portfolio/v4/overview/erc20/profit_and_loss"

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		Path:   u,
		Body:   nil,
	}

	var response GetTokensProfitAndLossResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetTokensDetails(ctx context.Context, params GetTokensDetailsParams) (*GetTokensDetailsResponse, error) {
	u := "/portfolio/portfolio/v4/overview/erc20/details"

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		Path:   u,
		Body:   nil,
	}

	var response GetTokensDetailsResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) IsServiceAvailable(ctx context.Context) (*IsServiceAvailableResponse, error) {
	u := "/portfolio/portfolio/v4/general/is_available"

	payload := common.RequestPayload{
		Method: "GET",
		Params: nil,
		Path:   u,
		Body:   nil,
	}

	var response IsServiceAvailableResponse
	err := api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetSupportedChains(ctx context.Context) (*GetSupportedChainsResponse, error) {
	u := "portfolio/portfolio/v4/general/supported_chains"

	payload := common.RequestPayload{
		Method: "GET",
		Params: nil,
		Path:   u,
		Body:   nil,
	}

	var response GetSupportedChainsResponse
	err := api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetCurrentValue(ctx context.Context, params GetCurrentValueParams) (*GetCurrentValueResponse, error) {
	u := "/portfolio/portfolio/v4/general/current_value"

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		Path:   u,
		Body:   nil,
	}

	var response GetCurrentValueResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetProfitAndLoss(ctx context.Context, params GetProfitAndLossParams) (*GetProfitAndLossResponse, error) {
	u := "/portfolio/portfolio/v4/general/profit_and_loss"

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		Path:   u,
		Body:   nil,
	}

	var response GetProfitAndLossResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetValueChart(ctx context.Context, params GetValueChartParams) (*GetValueChartResponse, error) {
	u := "/portfolio/portfolio/v4/general/value_chart"

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		Path:   u,
		Body:   nil,
	}

	var response GetValueChartResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Deprecated: Use GetTokensProfitAndLoss (standardized on the ProfitAndLoss spelling).
func (api *api) GetTokensProfitLoss(ctx context.Context, params GetTokensProfitAndLossParams) (*GetTokensProfitAndLossResponse, error) {
	return api.GetTokensProfitAndLoss(ctx, params)
}

// Deprecated: Use GetProfitAndLoss.
func (api *api) GetProfitLoss(ctx context.Context, params GetProfitAndLossParams) (*GetProfitAndLossResponse, error) {
	return api.GetProfitAndLoss(ctx, params)
}
