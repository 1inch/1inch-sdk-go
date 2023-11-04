package client

import (
	"context"
	"fmt"
	"net/http"

	"1inch-sdk-golang/client/swap"
)

type SwapService service

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

func (s *SwapService) GetSwap(ctx context.Context, params swap.AggregationControllerGetSwapParams) (*swap.SwapResponse, *http.Response, error) {
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

	var swap swap.SwapResponse
	res, err := s.client.Do(ctx, req, &swap)
	if err != nil {
		return nil, nil, err
	}

	return &swap, res, nil
}

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
