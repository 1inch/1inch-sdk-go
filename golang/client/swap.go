package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/1inch/1inch-sdk/golang/client/swap"
)

type SwapService service

// ApproveAllowance returns the allowance the 1inch router has to spend a token on behalf of a wallet
func (s *SwapService) ApproveAllowance(ctx context.Context, params swap.ApproveAllowanceParams) (*swap.AllowanceResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/approve/allowance", params.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	u, err = addQueryParameters(u, params.ApproveControllerGetAllowanceParams)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var allowanceResponse swap.AllowanceResponse
	res, err := s.client.Do(ctx, req, &allowanceResponse)
	if err != nil {
		return nil, nil, err
	}

	return &allowanceResponse, res, nil
}

// ApproveSpender returns the address of the 1inch router contract
func (s *SwapService) ApproveSpender(ctx context.Context, params swap.ApproveSpenderParams) (*swap.SpenderResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/approve/spender", params.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var spender swap.SpenderResponse
	res, err := s.client.Do(ctx, req, &spender)
	if err != nil {
		return nil, nil, err
	}

	return &spender, res, nil
}

// ApproveTransaction returns the transaction data for approving the 1inch router to spend a token on behalf of a wallet
func (s *SwapService) ApproveTransaction(ctx context.Context, params swap.ApproveTransactionParams) (*swap.ApproveCallDataResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/approve/transaction", params.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	u, err = addQueryParameters(u, params.ApproveControllerGetCallDataParams)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var approveCallData swap.ApproveCallDataResponse
	res, err := s.client.Do(ctx, req, &approveCallData)
	if err != nil {
		return nil, nil, err
	}

	return &approveCallData, res, nil
}

// GetLiquiditySources returns all liquidity sources tracked by the 1inch Aggregation Protocol for a given chain
func (s *SwapService) GetLiquiditySources(ctx context.Context, params swap.GetLiquiditySourcesParams) (*swap.ProtocolsResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/liquidity-sources", params.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var liquiditySources swap.ProtocolsResponse
	res, err := s.client.Do(ctx, req, &liquiditySources)
	if err != nil {
		return nil, nil, err
	}

	return &liquiditySources, res, nil
}

// GetQuote returns the quote for a potential swap through the Aggregation Protocol
func (s *SwapService) GetQuote(ctx context.Context, params swap.GetQuoteParams) (*swap.QuoteResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/quote", params.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	// Token info is used by certain parts of the SDK and more info by default is helpful to integrators
	// Because we use generated structs with concrete types, extra data is forced on regardless of what the user passes in.
	params.IncludeTokensInfo = true
	params.IncludeGas = true
	params.IncludeProtocols = true

	u, err = addQueryParameters(u, params.AggregationControllerGetQuoteParams)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var quote swap.QuoteResponse
	res, err := s.client.Do(ctx, req, &quote)
	if err != nil {
		return nil, nil, err
	}

	return &quote, res, nil
}

// GetSwapData returns a swap quote with transaction data that can be used to execute a swap through the Aggregation Protocol
func (s *SwapService) GetSwapData(ctx context.Context, params swap.GetSwapDataParams) (*swap.SwapResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/swap", params.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	// Token info is used by certain parts of the SDK and more info by default is helpful to integrators
	// Because we use generated structs with concrete types, extra data is forced on regardless of what the user passes in.
	params.IncludeTokensInfo = true
	params.IncludeGas = true
	params.IncludeProtocols = true

	u, err = addQueryParameters(u, params.AggregationControllerGetSwapParams)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var swapResponse swap.SwapResponse
	res, err := s.client.Do(ctx, req, &swapResponse)
	if err != nil {
		return nil, nil, err
	}

	if !params.SkipWarnings {
		err = swap.ConfirmSwapDataWithUser(&swapResponse, params.Amount, params.Slippage)
		if err != nil {
			return nil, nil, err
		}
	}

	return &swapResponse, res, nil
}

// GetTokens returns all tokens officially tracked by the 1inch Aggregation Protocol for a given chain
func (s *SwapService) GetTokens(ctx context.Context, params swap.GetTokensParams) (*swap.TokensResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/tokens", params.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var tokens swap.TokensResponse
	res, err := s.client.Do(ctx, req, &tokens)
	if err != nil {
		return nil, nil, err
	}

	return &tokens, res, nil
}
