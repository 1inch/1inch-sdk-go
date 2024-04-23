package aggregation

import (
	"context"
	"fmt"

	"github.com/1inch/1inch-sdk-go/common"
)

// GetApproveAllowance returns the allowance the 1inch router has to spend a token on behalf of a wallet
func (api *api) GetApproveAllowance(ctx context.Context, params ApproveAllowanceParams) (*AllowanceResponse, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/approve/allowance", api.chainId)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params.ApproveControllerGetAllowanceParams,
		U:      u,
		Body:   nil,
	}

	var allowanceResponse AllowanceResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &allowanceResponse)
	if err != nil {
		return nil, err
	}

	return &allowanceResponse, nil
}

// GetApproveSpender returns the address of the 1inch router contract
func (api *api) GetApproveSpender(ctx context.Context, params ApproveSpenderParams) (*SpenderResponse, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/approve/spender", api.chainId)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		U:      u,
		Params: nil,
		Body:   nil,
	}

	var spender SpenderResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &spender)
	if err != nil {
		return nil, err
	}

	return &spender, nil
}

// GetApproveTransaction returns the transaction data for approving the 1inch router to spend a token on behalf of a wallet
func (api *api) GetApproveTransaction(ctx context.Context, params ApproveTransactionParams) (*ApproveCallDataResponse, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/approve/transaction", api.chainId)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params.ApproveControllerGetCallDataParams,
		U:      u,
		Body:   nil,
	}

	var approveCallData ApproveCallDataResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &approveCallData)
	if err != nil {
		return nil, err
	}
	return &approveCallData, nil
}

// GetLiquiditySources returns all liquidity sources tracked by the 1inch Aggregation Protocol for a given chain
func (api *api) GetLiquiditySources(ctx context.Context, params GetLiquiditySourcesParams) (*ProtocolsResponse, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/liquidity-sources", api.chainId)

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

	var liquiditySources ProtocolsResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &liquiditySources)
	if err != nil {
		return nil, err
	}

	return &liquiditySources, nil
}

// GetQuote returns the quote for a potential swap through the Aggregation Protocol
func (api *api) GetQuote(ctx context.Context, params GetQuoteParams) (*QuoteResponse, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/quote", api.chainId)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params.AggregationControllerGetQuoteParams,
		U:      u,
	}

	var quote QuoteResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &quote)
	if err != nil {
		return nil, err
	}

	return &quote, nil
}

// GetSwap returns a swap quote with transaction data that can be used to execute a swap through the Aggregation Protocol
func (api *api) GetSwap(ctx context.Context, params AggregationControllerGetSwapParams) (*SwapResponseExtended, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/swap", api.chainId)

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

	var swapResponse SwapResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &swapResponse)
	if err != nil {
		return nil, err
	}
	return normalizeSwapResponse(swapResponse)
}

// GetTokens returns all tokens officially tracked by the 1inch Aggregation Protocol for a given chain
func (api *api) GetTokens(ctx context.Context, params GetTokensParams) (*TokensResponse, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/tokens", api.chainId)

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

	var tokens TokensResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &tokens)
	if err != nil {
		return nil, err
	}

	return &tokens, nil
}
