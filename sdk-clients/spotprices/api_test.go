package spotprices

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/common"
	"github.com/1inch/1inch-sdk-go/constants"
)

type MockHttpExecutor struct {
	Called      bool
	ExecuteErr  error
	ResponseObj interface{}
}

func (m *MockHttpExecutor) ExecuteRequest(ctx context.Context, payload common.RequestPayload, v interface{}) error {
	m.Called = true
	if m.ExecuteErr != nil {
		return m.ExecuteErr
	}

	// Copy the mock response object to v
	if m.ResponseObj != nil && v != nil {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Ptr || rv.IsNil() {
			return fmt.Errorf("v must be a non-nil pointer")
		}
		reflect.Indirect(rv).Set(reflect.ValueOf(m.ResponseObj))
	}
	return nil
}

func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }
func boolPtr(b bool) *bool    { return &b }

func TestGetPricesForWhitelistedTokens(t *testing.T) {
	ctx := context.Background()

	mockedResp := PricesForWhitelistedTokensResponse{
		"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2": "1000000000000000000",
		"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee": "1000000000000000000",
		"0xc3d688b66703497daa19211eedff47f25384cdc3": "344516328099167",
		"0x320623b8e4ff03373931769a31fc52a4e78b5d70": "2136128317246",
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
		chainId:      constants.EthereumChainId,
	}

	params := GetWhitelistedTokensPricesParams{
		Currency: GetWhitelistedTokensPricesParamsCurrency(USD),
	}

	prices, err := api.GetPricesForWhitelistedTokens(ctx, params)
	require.NoError(t, err)
	require.NotNil(t, prices)

	if !mockExecutor.Called {
		t.Errorf("Expected ExecuteRequest to be called")
	}
	if !reflect.DeepEqual(prices, &mockedResp) {
		t.Errorf("Expected swap to be %+v, got %+v", prices, mockedResp)
	}
	require.NoError(t, err)
}

func TestGetPricesForRequestedTokens(t *testing.T) {
	ctx := context.Background()

	mockedResp := PricesForRequestedTokensResponse{
		"0x0d8775f648430679a709e98d2b0cb6250d2887ef": "1000000000000000000",
		"0x58b6a8a3302369daec383334672404ee733ab239": "1000000000000000000",
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
		chainId:      constants.EthereumChainId,
	}

	params := GetPricesRequestDto{
		Currency: GetPricesRequestDtoCurrency(USD),
		Tokens:   []string{"0x0d8775f648430679a709e98d2b0cb6250d2887ef", "0x58b6a8a3302369daec383334672404ee733ab239"},
	}

	prices, err := api.GetPricesForRequestedTokens(ctx, params)
	require.NoError(t, err)
	require.NotNil(t, prices)

	if !mockExecutor.Called {
		t.Errorf("Expected ExecuteRequest to be called")
	}
	if !reflect.DeepEqual(prices, &mockedResp) {
		t.Errorf("Expected swap to be %+v, got %+v", prices, mockedResp)
	}
	require.NoError(t, err)
}
