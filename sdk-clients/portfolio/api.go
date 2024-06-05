package portfolio

import (
	"context"

	"github.com/1inch/1inch-sdk-go/common"
)

func (api *api) GetProtocolsCurrentValue(ctx context.Context, params GetCurrentValuePortfolioV4OverviewProtocolsCurrentValueGetParams) (*GetPortfolioValueResponse, error) {
	u := "/portfolio/portfolio/v4/overview/protocols/current_value"

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

	var response GetPortfolioValueResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetProtocolsProfitAndLoss(ctx context.Context, params GetProfitAndLossPortfolioV4OverviewProtocolsProfitAndLossGetParams) (*GetPortfolioProfitAndLossResponse, error) {
	u := "/portfolio/portfolio/v4/overview/protocols/profit_and_loss"

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

	var response GetPortfolioProfitAndLossResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetProtocolsDetails(ctx context.Context, params GetDetailsPortfolioV4OverviewProtocolsDetailsGetParams) (*GetProtocolsDetailsResponse, error) {
	u := "/portfolio/portfolio/v4/overview/protocols/details"

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

	var response GetProtocolsDetailsResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetTokensCurrentValue(ctx context.Context, params GetCurrentValuePortfolioV4OverviewErc20CurrentValueGetParams) (*GetTokensCurrentValueResponse, error) {
	u := "/portfolio/portfolio/v4/overview/erc20/current_value"

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

	var response GetTokensCurrentValueResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetTokensProfitLoss(ctx context.Context, params GetProfitAndLossPortfolioV4OverviewErc20ProfitAndLossGetParams) (*GetTokensProfitLossResponse, error) {
	u := "/portfolio/portfolio/v4/overview/erc20/profit_and_loss"

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

	var response GetTokensProfitLossResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetTokensDetails(ctx context.Context, params GetDetailsPortfolioV4OverviewErc20DetailsGetParams) (*GetTokensDetailsResponse, error) {
	u := "/portfolio/portfolio/v4/overview/erc20/details"

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
		U:      u,
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
		U:      u,
		Body:   nil,
	}

	var response GetSupportedChainsResponse
	err := api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetCurrentValue(ctx context.Context, params GetCurrentValuePortfolioV4GeneralCurrentValueGetParams) (*GetCurrentValueResponse, error) {
	u := "/portfolio/portfolio/v4/general/current_value"

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

	var response GetCurrentValueResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetProfitLoss(ctx context.Context, params GetProfitAndLossPortfolioV4GeneralProfitAndLossGetParams) (*GetCurrentProfitLossResponse, error) {
	u := "/portfolio/portfolio/v4/general/profit_and_loss"

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

	var response GetCurrentProfitLossResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetValueChart(ctx context.Context, params GetValueChartPortfolioV4GeneralValueChartGetParams) (*GetValueChartResponse, error) {
	u := "/portfolio/portfolio/v4/general/value_chart"

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

	var response GetValueChartResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
