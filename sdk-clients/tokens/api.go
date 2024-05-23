package tokens

import (
	"context"
	"fmt"

	"github.com/1inch/1inch-sdk-go/common"
)

// SearchTokenAllChains Get Tokens that match the provided search criteria across all chains
func (api *api) SearchTokenAllChains(ctx context.Context, params SearchControllerSearchAllChainsParams) ([]ProviderTokenDtoFixed, error) {
	u := "/token/v1.2/search"

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

	var response []ProviderTokenDtoFixed
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SearchTokenSingleChain Get Tokens that match the provided search criteria on a specfic chain
func (api *api) SearchTokenSingleChain(ctx context.Context, params SearchControllerSearchSingleChainParams) ([]ProviderTokenDtoFixed, error) {
	u := fmt.Sprintf("/token/v1.2/%d/search", api.chainId)

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

	var response []ProviderTokenDtoFixed
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (api *api) WhitelistedTokens(ctx context.Context, params TokenListControllerTokensParams) (map[string]ProviderTokenDto, error) {
	u := fmt.Sprintf("/token/v1.2/%d", api.chainId)

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

	var response map[string]ProviderTokenDto
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (api *api) WhitelistedTokensAsList(ctx context.Context, params TokenListControllerTokensParams) (*TokenListResponseDto, error) {
	u := fmt.Sprintf("/token/v1.2/%d/token-list", api.chainId)

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

	var response TokenListResponseDto
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetCustomTokens(ctx context.Context, params CustomTokensControllerGetTokensInfoParams) (map[string]ProviderTokenDto, error) {
	u := fmt.Sprintf("/token/v1.2/%d/custom", api.chainId)

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

	var response map[string]ProviderTokenDto
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (api *api) GetCustomToken(ctx context.Context, params CustomTokensControllerGetTokenInfoParams) (*ProviderTokenDtoFixed, error) {
	u := fmt.Sprintf("/token/v1.2/%d/custom/%s", api.chainId, params.Address)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: nil,
		U:      u,
		Body:   nil,
	}

	var response ProviderTokenDtoFixed
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
