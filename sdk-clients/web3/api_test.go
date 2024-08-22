package web3

import (
	"context"
	"encoding/json"
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

func TestPerformRpcCallAgainstFullNode(t *testing.T) {
	ctx := context.Background()

	mockedResp := map[string]interface{}{
		"result": "0x10",
		"id":     "1",
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
		chainId:      constants.EthereumChainId,
	}

	params := PerformRpcCallAgainstFullNodeParams{
		PostChainIdJSONBody: PostChainIdJSONBody{
			Jsonrpc: "2.0",
			Method:  "eth_blockNumber",
			Params:  []string{},
			Id:      "1",
		},
	}

	response, err := api.PerformRpcCallAgainstFullNode(ctx, params)
	require.NoError(t, err)
	require.NotNil(t, response)

	if !mockExecutor.Called {
		t.Errorf("Expected ExecuteRequest to be called")
	}

	// Validate the request payload
	expectedURL := fmt.Sprintf("web3/%d", api.chainId)
	expectedBody, _ := json.Marshal(params.PostChainIdJSONBody)
	expectedPayload := common.RequestPayload{
		Method: "POST",
		Params: nil,
		U:      expectedURL,
		Body:   expectedBody,
	}

	require.Equal(t, expectedPayload.U, expectedURL)
	require.Equal(t, string(expectedPayload.Body), string(expectedBody))

	// Validate the response
	require.Equal(t, mockedResp, response)
}

func TestPerformRpcCall(t *testing.T) {
	ctx := context.Background()

	mockedResp := map[string]interface{}{
		"result": "0x10",
		"id":     "1",
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
		chainId:      constants.EthereumChainId,
	}

	params := PerformRpcCallParams{
		PostChainIdNodeTypeParamsNodeType: "full",
		PostChainIdJSONBody: PostChainIdJSONBody{
			Jsonrpc: "2.0",
			Method:  "eth_blockNumber",
			Params:  nil,
			Id:      "1",
		},
	}

	response, err := api.PerformRpcCall(ctx, params)
	require.NoError(t, err)
	require.NotNil(t, response)

	if !mockExecutor.Called {
		t.Errorf("Expected ExecuteRequest to be called")
	}

	// Validate the request payload
	expectedURL := fmt.Sprintf("web3/%d/%s", api.chainId, params.PostChainIdNodeTypeParamsNodeType)
	expectedBody, _ := json.Marshal(params)
	expectedPayload := common.RequestPayload{
		Method: "POST",
		Params: nil,
		U:      expectedURL,
		Body:   expectedBody,
	}

	require.Equal(t, expectedPayload.U, expectedURL)
	require.Equal(t, string(expectedPayload.Body), string(expectedBody))

	// Validate the response
	require.Equal(t, mockedResp, response)
}
