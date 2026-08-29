package tokens

import (
	"context"
	"fmt"

	"github.com/1inch/1inch-sdk-go/v4/common"
)

// SearchTokenAllChains Get Tokens that match the provided search criteria across all chains
func (api *api) SearchTokenAllChains(ctx context.Context, params SearchAllChainsParams) ([]ProviderTokenDto, error) {
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

	var response []ProviderTokenDto
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// SearchTokenSingleChain Get Tokens that match the provided search criteria on a specific chain
func (api *api) SearchTokenSingleChain(ctx context.Context, params SearchSingleChainParams) ([]ProviderTokenDto, error) {
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

	var response []ProviderTokenDto
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (api *api) GetWhitelistedTokens(ctx context.Context, params GetWhitelistedTokensParams) (map[string]ProviderTokenDto, error) {
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

func (api *api) GetWhitelistedTokensAsList(ctx context.Context, params GetWhitelistedTokensParams) (*TokenListResponseDto, error) {
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

func (api *api) GetCustomTokens(ctx context.Context, params GetCustomTokensParams) (map[string]ProviderTokenDto, error) {
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

func (api *api) GetCustomToken(ctx context.Context, params GetCustomTokenParams) (*ProviderTokenDto, error) {
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

	var response ProviderTokenDto
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Deprecated: Use GetWhitelistedTokens (added the Get prefix used by every other method).
func (api *api) WhitelistedTokens(ctx context.Context, params GetWhitelistedTokensParams) (map[string]ProviderTokenDto, error) {
	return api.GetWhitelistedTokens(ctx, params)
}

// Deprecated: Use GetWhitelistedTokensAsList.
func (api *api) WhitelistedTokensAsList(ctx context.Context, params GetWhitelistedTokensParams) (*TokenListResponseDto, error) {
	return api.GetWhitelistedTokensAsList(ctx, params)
}
