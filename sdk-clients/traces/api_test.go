package traces

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

func TestGetSyncedInterval(t *testing.T) {
	ctx := context.Background()

	mockedResp := ReadSyncedIntervalResponseDto{
		From: 1,
		To:   19868703,
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
		chainId:      constants.EthereumChainId,
	}

	prices, err := api.GetSyncedInterval(ctx)
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

func TestGetBlockTraceByNumber(t *testing.T) {
	ctx := context.Background()

	mockedResp := CoreBuiltinBlockTracesDto{
		BlockHash:      "0xd5480dba100c54b6d545baa1ca5e70a0d13983cac46f0c8a9a3a0e049e531f28",
		BlockTimestamp: "0x64771d1b",
		Traces: []CoreBuiltinTransactionRootSuccessTraceDto{
			CoreBuiltinTransactionRootSuccessTraceDto{
				Calls: nil, // No nested calls provided
				Error: "",
				Events: []CoreBuiltinTraceLogDto{
					CoreBuiltinTraceLogDto{
						Topics:   nil,
						Contract: "0xdac17f958d2ee523a2206206994597c13d831ec7",
						Data:     "0x0000000000000000000000000000000000000000000000000000000059682f00",
					},
				},
				From:                 "0x05adb055fc6ab6e0eff6537229d93086544d7d12",
				Gas:                  "38330",
				GasHex:               "0xea16",
				GasPrice:             "0x123474949b",
				GasUsed:              41297,
				Input:                "0xa9059cbb000000000000000000000000afe8edba6577170d3f422c05d91c117b36825df00000000000000000000000000000000000000000000000000000000059682f00",
				MaxFeePerGas:         "0x1240be3e8a",
				MaxPriorityFeePerGas: "0xb45267ad1",
				Nonce:                "0x3",
				Output:               "", // No output provided
				RevertReason:         "", // No revert reason provided
				To:                   "0xdac17f958d2ee523a2206206994597c13d831ec7",
				TxHash:               "0x1c59f28b7fbbbe9b5b0257facd2231acdac9d7cf52422dbd13bcb463d4426d2b",
				Type:                 "CALL",
				Value:                "0x0",
			},
		},
		Type:    "CUSTOM",
		Version: "10.0.1",
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
		chainId:      constants.EthereumChainId,
	}

	params := GetBlockTraceByNumberParam(17378177)

	prices, err := api.GetBlockTraceByNumber(ctx, params)
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

func TestGetTxTraceByNumberAndHash(t *testing.T) {
	ctx := context.Background()

	mockedResp := TransactionTraceResponse{
		TransactionTrace: TransactionTrace{
			TxHash:       "0x1c59f28b7fbbbe9b5b0257facd2231acdac9d7cf52422dbd13bcb463d4426d2b",
			Nonce:        "0x3",
			GasPrice:     "0x123474949b",
			Type:         "CALL",
			From:         "0x05adb055fc6ab6e0eff6537229d93086544d7d12",
			To:           "0xdac17f958d2ee523a2206206994597c13d831ec7",
			GasLimit:     38330,
			GasActual:    36497,
			GasHex:       "0xea16",
			GasUsed:      41297,
			IntrinsicGas: 0,
			GasRefund:    4800,
			Input:        "0xa9059cbb000000000000000000000000afe8edba6577170d3f422c05d91c117b36825df00000000000000000000000000000000000000000000000000000000059682f00",
			Calls:        []Call{},
			Logs: []Log{
				{
					Topics:   []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0x5adb055fc6ab6e0eff6537229d93086544d7d12", "0xafe8edba6577170d3f422c05d91c117b36825df0"},
					Contract: "0xdac17f958d2ee523a2206206994597c13d831ec7",
					Data:     "0x0000000000000000000000000000000000000000000000000000000059682f00",
				},
			},
			Status:               "STOPPED",
			Storage:              []StorageItem{},
			Value:                "0x0",
			MaxFeePerGas:         "0x1240be3e8a",
			MaxPriorityFeePerGas: "0xb45267ad1",
			Depth:                0,
		},
		Type: "CUSTOM",
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
		chainId:      constants.EthereumChainId,
	}

	params := GetTxTraceByNumberAndHashParams{
		BlockNumber:     17378176,
		TransactionHash: "0x16897e492b2e023d8f07be9e925f2c15a91000ef11a01fc71e70f75050f1e03c",
	}

	prices, err := api.GetTxTraceByNumberAndHash(ctx, params)
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

func TestGetTxTraceByNumberAndOffset(t *testing.T) {
	ctx := context.Background()

	mockedResp := TransactionTraceResponse{
		TransactionTrace: TransactionTrace{
			TxHash:       "0x1c59f28b7fbbbe9b5b0257facd2231acdac9d7cf52422dbd13bcb463d4426d2b",
			Nonce:        "0x3",
			GasPrice:     "0x123474949b",
			Type:         "CALL",
			From:         "0x05adb055fc6ab6e0eff6537229d93086544d7d12",
			To:           "0xdac17f958d2ee523a2206206994597c13d831ec7",
			GasLimit:     38330,
			GasActual:    36497,
			GasHex:       "0xea16",
			GasUsed:      41297,
			IntrinsicGas: 0,
			GasRefund:    4800,
			Input:        "0xa9059cbb000000000000000000000000afe8edba6577170d3f422c05d91c117b36825df00000000000000000000000000000000000000000000000000000000059682f00",
			Calls:        []Call{},
			Logs: []Log{
				{
					Topics:   []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0x5adb055fc6ab6e0eff6537229d93086544d7d12", "0xafe8edba6577170d3f422c05d91c117b36825df0"},
					Contract: "0xdac17f958d2ee523a2206206994597c13d831ec7",
					Data:     "0x0000000000000000000000000000000000000000000000000000000059682f00",
				},
			},
			Status:               "STOPPED",
			Storage:              []StorageItem{},
			Value:                "0x0",
			MaxFeePerGas:         "0x1240be3e8a",
			MaxPriorityFeePerGas: "0xb45267ad1",
			Depth:                0,
		},
		Type: "CUSTOM",
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
		chainId:      constants.EthereumChainId,
	}

	params := GetTxTraceByNumberAndOffsetParams{
		BlockNumber: 17378176,
		Offset:      1,
	}

	prices, err := api.GetTxTraceByNumberAndOffset(ctx, params)
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
