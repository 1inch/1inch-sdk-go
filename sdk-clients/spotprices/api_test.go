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

	mockedResp := PricesForWhitelistedTokensResponse{}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
		chainId:      constants.EthereumChainId,
	}

	params := GetWhitelistedTokensPricesParams{}

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
