package portfolio

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/common"
)

type MockHttpExecutor struct {
	Called      bool
	ExecuteErr  error
	ResponseObj any
}

func (m *MockHttpExecutor) ExecuteRequest(ctx context.Context, payload common.RequestPayload, v any) error {
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

func TestGetCurrentValue(t *testing.T) {
	ctx := context.Background()

	// JSON data
	data := `{
	    "result": [
	        {
	            "address": "0xd8da6bf26964af9d7eed9e03e53415d37aa96045",
	            "value_usd": 2746515.5134213683
	        }
	    ],
	    "system": {
	        "click_time": 0.2215709686279297,
	        "node_time": 0.4157066345214844,
	        "microservices_time": 0.03647613525390625,
	        "redis_time": 0.0,
	        "total_time": 0.43703126907348633
	    }
	}`

	// Create an instance of the struct
	var mockedResp GetCurrentValueResponse

	// Unmarshal the JSON data into the struct
	err := json.Unmarshal([]byte(data), &mockedResp)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON: %v", err)
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
	}

	params := GetCurrentValuePortfolioV4GeneralCurrentValueGetParams{
		Addresses: []string{"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"},
		ChainId:   1,
	}

	prices, err := api.GetCurrentValue(ctx, params)
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
