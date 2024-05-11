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

func TestGetGasPriceEIP1559(t *testing.T) {
	ctx := context.Background()

	mockedResp := Eip1559GasPriceResponse{
		BaseFee: "10000000",
		High: Eip1559GasValueResponse{
			MaxFeePerGas:         "11400000",
			MaxPriorityFeePerGas: "21400000",
		},
		Instant: Eip1559GasValueResponse{
			MaxFeePerGas:         "11500000",
			MaxPriorityFeePerGas: "21500000",
		},
		Low: Eip1559GasValueResponse{
			MaxFeePerGas:         "11500000",
			MaxPriorityFeePerGas: "21500000",
		},
		Medium: Eip1559GasValueResponse{
			MaxFeePerGas:         "11500000",
			MaxPriorityFeePerGas: "21500000",
		},
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
		chainId:      constants.EthereumChainId,
	}

	prices, err := api.GetGasPriceEIP1559(ctx)
	require.NoError(t, err)
	require.NotNil(t, prices)

	if !mockExecutor.Called {
		t.Errorf("Expected ExecuteRequest to be called")
	}
	if !reflect.DeepEqual(prices, &mockedResp) {
		t.Errorf("Expected swap to be %+v, got %+v", prices, mockedResp)
	}
	require.NoError(t, err)

	_, err = api.GetGasPriceLegacy(ctx)
	require.Error(t, err)
}

func TestGetGasPriceLegacy(t *testing.T) {
	ctx := context.Background()

	mockedResp := GetGasPriceLegacyResponse{
		Standard: "77000000",
		Fast:     "91000000",
		Instant:  "105000000",
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
		chainId:      constants.AuroraChainId,
	}

	prices, err := api.GetGasPriceLegacy(ctx)
	require.NoError(t, err)
	require.NotNil(t, prices)

	if !mockExecutor.Called {
		t.Errorf("Expected ExecuteRequest to be called")
	}
	if !reflect.DeepEqual(prices, &mockedResp) {
		t.Errorf("Expected swap to be %+v, got %+v", prices, mockedResp)
	}
	require.NoError(t, err)

	_, err = api.GetGasPriceEIP1559(ctx)
	require.Error(t, err)
}
