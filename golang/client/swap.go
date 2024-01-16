package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/1inch/1inch-sdk/golang/client/swap"
)

type SwapService service

// ApproveAllowance returns the allowance the 1inch router has to spend a token on behalf of a wallet
func (s *SwapService) ApproveAllowance(ctx context.Context, params swap.ApproveControllerGetAllowanceParams) (*swap.AllowanceResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/approve/allowance", s.client.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	u, err = addQueryParameters(u, params)
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
func (s *SwapService) ApproveSpender(ctx context.Context) (*swap.SpenderResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/approve/spender", s.client.ChainId)

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
func (s *SwapService) ApproveTransaction(ctx context.Context, params swap.ApproveControllerGetCallDataParams) (*swap.ApproveCallDataResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/approve/transaction", s.client.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	u, err = addQueryParameters(u, params)
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
func (s *SwapService) GetLiquiditySources(ctx context.Context) (*swap.ProtocolsResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/liquidity-sources", s.client.ChainId)

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
func (s *SwapService) GetQuote(ctx context.Context, params swap.AggregationControllerGetQuoteParams) (*swap.QuoteResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/quote", s.client.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	u, err = addQueryParameters(u, params)
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

// TODO temporarily adding a bool to the function call until config refactor

// GetSwapData returns a swap quote with transaction data that can be used to execute a swap through the Aggregation Protocol
func (s *SwapService) GetSwapData(ctx context.Context, params swap.AggregationControllerGetSwapParams, skipWarnings bool) (*swap.SwapResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/swap", s.client.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	u, err = addQueryParameters(u, params)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var swapResponse swap.SwapResponse // TODO the struct does not match the response well
	res, err := s.client.Do(ctx, req, &swapResponse)
	if err != nil {
		return nil, nil, err
	}

	swapResponse.FromToken = &swap.TokenInfo{
		Address: params.Src,
	}
	swapResponse.ToToken = &swap.TokenInfo{
		Address: params.Dst,
	}

	if !skipWarnings {
		err = swap.ConfirmSwapDataWithUser(&swapResponse, params.Amount, params.Slippage, s.client.EthClient)
		if err != nil {
			return nil, nil, err
		}
	}

	return &swapResponse, res, nil
}

// GetTokens returns all tokens officially tracked by the 1inch Aggregation Protocol for a given chain
func (s *SwapService) GetTokens(ctx context.Context) (*swap.TokensResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/tokens", s.client.ChainId)

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
