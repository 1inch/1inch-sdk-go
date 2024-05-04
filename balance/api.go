package balance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/1inch/1inch-sdk-go/common"
)

// GetAggregatedBalancesAndAllowances Get balances and allowances by spender for list of EVM addresses
func (api *api) GetAggregatedBalancesAndAllowances(ctx context.Context, params AggregatedBalancesAndAllowancesParams) (*AggregatedBalancesAndAllowancesResponse, error) {
	u := fmt.Sprintf("/balance/v1.2/%d/aggregatedBalancesAndAllowances/%s", api.chainId, params.Spender)

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

	var response AggregatedBalancesAndAllowancesResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetBalancesByWalletAddress Get balances of tokens for walletAddress for default token list (1inch tokens list)
func (api *api) GetBalancesByWalletAddress(ctx context.Context, params BalancesByWalletAddressParams) (*BalancesByWalletAddressResponse, error) {
	u := fmt.Sprintf("/balance/v1.2/%d/balances/%s", api.chainId, params.WalletAddress)

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

	var response BalancesByWalletAddressResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetBalancesOfCustomTokensByWalletAddress Get balances of custom tokens for walletAddress
// Takes wallet address and provided tokens and provides balance of each token
func (api *api) GetBalancesOfCustomTokensByWalletAddress(ctx context.Context, params BalancesOfCustomTokensByWalletAddressParams) (*BalancesOfCustomTokensByWalletAddressResponse, error) {
	u := fmt.Sprintf("/balance/v1.2/%d/balances/%s", api.chainId, params.WalletAddress)

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

	var response BalancesOfCustomTokensByWalletAddressResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetBalancesOfCustomTokensByWalletAddress Get balances of custom tokens for walletAddress
// Takes wallet address and provided tokens and provides balance of each token
func (api *api) GetBalancesOfCustomTokensByWalletAddressesList(ctx context.Context, params BalancesOfCustomTokensByWalletAddressesListParams) (*BalancesOfCustomTokensByWalletAddressesListResponse, error) {
	u := fmt.Sprintf("/balance/v1.2/%d/balances/multiple/walletsAndTokens", api.chainId)

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

	var response BalancesOfCustomTokensByWalletAddressesListResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
