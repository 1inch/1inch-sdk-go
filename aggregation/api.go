package aggregation

import (
	"context"
	"fmt"
	"github.com/1inch/1inch-sdk-go/internal/helpers"
	"github.com/google/go-querystring/query"
	"net/http"
	"net/url"
	"reflect"
)

// GetApproveAllowance returns the allowance the 1inch router has to spend a token on behalf of a wallet
func (api *apiActions) GetApproveAllowance(ctx context.Context, params ApproveAllowanceParams) (*AllowanceResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/approve/allowance", params.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	u, err = getQueryParameters(u, params.ApproveControllerGetAllowanceParams)
	if err != nil {
		return nil, nil, err
	}

	req, err := api.httpClient.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var allowanceResponse AllowanceResponse
	res, err := s.client.Do(ctx, req, &allowanceResponse)
	if err != nil {
		return nil, nil, err
	}

	return &allowanceResponse, res, nil
}

// GetApproveSpender returns the address of the 1inch router contract
func (s *apiActions) GetApproveSpender(ctx context.Context, params ApproveSpenderParams) (*SpenderResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/approve/spender", params.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var spender SpenderResponse
	res, err := s.client.Do(ctx, req, &spender)
	if err != nil {
		return nil, nil, err
	}

	return &spender, res, nil
}

// GetApproveTransaction returns the transaction data for approving the 1inch router to spend a token on behalf of a wallet
func (s *apiActions) GetApproveTransaction(ctx context.Context, params ApproveTransactionParams) (*ApproveCallDataResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/approve/transaction", params.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	u, err = getQueryParameters(u, params.ApproveControllerGetCallDataParams)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var approveCallData ApproveCallDataResponse
	res, err := s.client.Do(ctx, req, &approveCallData)
	if err != nil {
		return nil, nil, err
	}

	return &approveCallData, res, nil
}

// GetLiquiditySources returns all liquidity sources tracked by the 1inch Aggregation Protocol for a given chain
func (s *apiActions) GetLiquiditySources(ctx context.Context, params GetLiquiditySourcesParams) (*ProtocolsResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/liquidity-sources", params.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var liquiditySources ProtocolsResponse
	res, err := s.client.Do(ctx, req, &liquiditySources)
	if err != nil {
		return nil, nil, err
	}

	return &liquiditySources, res, nil
}

// GetQuote returns the quote for a potential swap through the Aggregation Protocol
func (s *apiActions) GetQuote(ctx context.Context, params GetQuoteParams) (*QuoteResponse, *http.Response, error) {
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

	u, err = getQueryParameters(u, params.AggregationControllerGetQuoteParams)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var quote QuoteResponse
	res, err := s.client.Do(ctx, req, &quote)
	if err != nil {
		return nil, nil, err
	}

	return &quote, res, nil
}

// GetSwap returns a swap quote with transaction data that can be used to execute a swap through the Aggregation Protocol
func (s *apiActions) GetSwap(ctx context.Context, params GetSwapParams) (*SwapResponse, *http.Response, error) {
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

	u, err = getQueryParameters(u, params.AggregationControllerGetSwapParams)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var swapResponse SwapResponse
	res, err := s.client.Do(ctx, req, &swapResponse)
	if err != nil {
		return nil, nil, err
	}
	return &swapResponse, res, nil
}

// GetTokens returns all tokens officially tracked by the 1inch Aggregation Protocol for a given chain
func (s *apiActions) GetTokens(ctx context.Context, params GetTokensParams) (*TokensResponse, *http.Response, error) {
	u := fmt.Sprintf("/swap/v5.2/%d/tokens", params.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var tokens TokensResponse
	res, err := s.client.Do(ctx, req, &tokens)
	if err != nil {
		return nil, nil, err
	}

	return &tokens, res, nil
}

// addQueryParameters adds the parameters in the struct params as URL query parameters to s.
// params must be a struct whose fields may contain "url" tags.
func getQueryParameters(s string, params interface{}) (string, error) {
	v := reflect.ValueOf(params)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(params)
	if err != nil {
		return s, err
	}

	for k, v := range qs {
		if helpers.IsScientificNotation(v[0]) {
			expanded, err := helpers.ExpandScientificNotation(v[0])
			if err != nil {
				return "", fmt.Errorf("failed to expand scientific notation for parameter %v with a value of %v: %v", k, v, err)
			}
			v[0] = expanded
		}
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
