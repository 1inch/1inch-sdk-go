package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/1inch/1inch-sdk-go/client/models"
)

type SwapService service

// GetApproveAllowance returns the allowance the 1inch router has to spend a token on behalf of a wallet
func (s *SwapService) GetApproveAllowance(ctx context.Context, params models.ApproveAllowanceParams) (*models.AllowanceResponse, *http.Response, error) {
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

	var allowanceResponse models.AllowanceResponse
	res, err := s.client.Do(ctx, req, &allowanceResponse)
	if err != nil {
		return nil, nil, err
	}

	return &allowanceResponse, res, nil
}

// GetApproveSpender returns the address of the 1inch router contract
func (s *SwapService) GetApproveSpender(ctx context.Context, params models.ApproveSpenderParams) (*models.SpenderResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/approve/spender", params.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var spender models.SpenderResponse
	res, err := s.client.Do(ctx, req, &spender)
	if err != nil {
		return nil, nil, err
	}

	return &spender, res, nil
}

// GetApproveTransaction returns the transaction data for approving the 1inch router to spend a token on behalf of a wallet
func (s *SwapService) GetApproveTransaction(ctx context.Context, params models.ApproveTransactionParams) (*models.ApproveCallDataResponse, *http.Response, error) {
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

	var approveCallData models.ApproveCallDataResponse
	res, err := s.client.Do(ctx, req, &approveCallData)
	if err != nil {
		return nil, nil, err
	}

	return &approveCallData, res, nil
}

// GetLiquiditySources returns all liquidity sources tracked by the 1inch Aggregation Protocol for a given chain
func (s *SwapService) GetLiquiditySources(ctx context.Context, params models.GetLiquiditySourcesParams) (*models.ProtocolsResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/liquidity-sources", params.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var liquiditySources models.ProtocolsResponse
	res, err := s.client.Do(ctx, req, &liquiditySources)
	if err != nil {
		return nil, nil, err
	}

	return &liquiditySources, res, nil
}

// GetQuote returns the quote for a potential swap through the Aggregation Protocol
func (s *SwapService) GetQuote(ctx context.Context, params models.GetQuoteParams) (*models.QuoteResponse, *http.Response, error) {
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

	var quote models.QuoteResponse
	res, err := s.client.Do(ctx, req, &quote)
	if err != nil {
		return nil, nil, err
	}

	return &quote, res, nil
}

// GetSwap returns a swap quote with transaction data that can be used to execute a swap through the Aggregation Protocol
func (s *SwapService) GetSwap(ctx context.Context, params models.GetSwapParams) (*models.SwapResponse, *http.Response, error) {
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

	var swapResponse models.SwapResponse
	res, err := s.client.Do(ctx, req, &swapResponse)
	if err != nil {
		return nil, nil, err
	}
	return &swapResponse, res, nil
}

// GetTokens returns all tokens officially tracked by the 1inch Aggregation Protocol for a given chain
func (s *SwapService) GetTokens(ctx context.Context, params models.GetTokensParams) (*models.TokensResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/tokens", params.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var tokens models.TokensResponse
	res, err := s.client.Do(ctx, req, &tokens)
	if err != nil {
		return nil, nil, err
	}

	return &tokens, res, nil
}
