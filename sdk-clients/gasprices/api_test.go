package gasprices

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

func TestGetPricesForRequestedTokensLarge(t *testing.T) {
	ctx := context.Background()

	mockedResp := PricesForRequestedTokensResponse{
		"0x0d8775f648430679a709e98d2b0cb6250d2887ef": "1000000000000000000",
		"0x58b6a8a3302369daec383334672404ee733ab239": "1000000000000000000",
		"0x320623b8e4ff03373931769a31fc52a4e78b5d70": "1000000000000000000",
		"0x71ab77b7dbb4fa7e017bc15090b2163221420282": "1000000000000000000",
		"0x256d1fce1b1221e8398f65f9b36033ce50b2d497": "1000000000000000000",
		"0x85f17cf997934a597031b2e18a9ab6ebd4b9f6a4": "1000000000000000000",
		"0x55c08ca52497e2f1534b59e2917bf524d4765257": "1000000000000000000",
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
		chainId:      constants.EthereumChainId,
	}

	const (
		tokenAddress1 = "0x0d8775f648430679a709e98d2b0cb6250d2887ef"
		tokenAddress2 = "0x58b6a8a3302369daec383334672404ee733ab239"
		tokenAddress3 = "0x320623b8e4ff03373931769a31fc52a4e78b5d70"
		tokenAddress4 = "0x71ab77b7dbb4fa7e017bc15090b2163221420282"
		tokenAddress5 = "0x256d1fce1b1221e8398f65f9b36033ce50b2d497"
		tokenAddress6 = "0x85f17cf997934a597031b2e18a9ab6ebd4b9f6a4"
		tokenAddress7 = "0x55c08ca52497e2f1534b59e2917bf524d4765257"
	)

	params := GetPricesRequestDto{
		Currency: GetPricesRequestDtoCurrency(USD),
		Tokens:   []string{tokenAddress1, tokenAddress2, tokenAddress3, tokenAddress4, tokenAddress5, tokenAddress6, tokenAddress7},
	}

	prices, err := api.GetPricesForRequestedTokensLarge(ctx, params)
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

func TestGetCustomCurrenciesList(t *testing.T) {
	ctx := context.Background()

	mockedResp := CurrenciesResponseDto{
		Codes: []string{
			"USD",
			"EUR",
		},
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
		chainId:      constants.EthereumChainId,
	}

	prices, err := api.GetCustomCurrenciesList(ctx)
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
