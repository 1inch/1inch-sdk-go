package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/constants"
)

func TestNormalizeSwapResponse(t *testing.T) {
	tests := []struct {
		name        string
		resp        SwapResponse
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid response",
			resp: SwapResponse{
				Tx: TransactionData{
					To:       "0x6B175474E89094C44Da98b954EedeAC495271d0F",
					Gas:      100000,
					GasPrice: "20000000000",
					Value:    "1000000000000000000",
					Data:     "0x1234567890abcdef",
				},
			},
			expectError: false,
		},
		{
			name: "Valid response - zero value",
			resp: SwapResponse{
				Tx: TransactionData{
					To:       "0x6B175474E89094C44Da98b954EedeAC495271d0F",
					Gas:      100000,
					GasPrice: "20000000000",
					Value:    "0",
					Data:     "0x",
				},
			},
			expectError: false,
		},
		{
			name: "Invalid To address",
			resp: SwapResponse{
				Tx: TransactionData{
					To:       "invalid",
					Gas:      100000,
					GasPrice: "20000000000",
					Value:    "0",
					Data:     "0x",
				},
			},
			expectError: true,
			errorMsg:    "invalid to address",
		},
		{
			name: "Invalid GasPrice - not a number",
			resp: SwapResponse{
				Tx: TransactionData{
					To:       "0x6B175474E89094C44Da98b954EedeAC495271d0F",
					Gas:      100000,
					GasPrice: "invalid",
					Value:    "0",
					Data:     "0x",
				},
			},
			expectError: true,
			errorMsg:    "invalid gas price",
		},
		{
			name: "Invalid Value - not a number",
			resp: SwapResponse{
				Tx: TransactionData{
					To:       "0x6B175474E89094C44Da98b954EedeAC495271d0F",
					Gas:      100000,
					GasPrice: "20000000000",
					Value:    "invalid",
					Data:     "0x",
				},
			},
			expectError: true,
			errorMsg:    "invalid tx value",
		},
		{
			name: "Invalid Data - not hex",
			resp: SwapResponse{
				Tx: TransactionData{
					To:       "0x6B175474E89094C44Da98b954EedeAC495271d0F",
					Gas:      100000,
					GasPrice: "20000000000",
					Value:    "0",
					Data:     "not-hex-data",
				},
			},
			expectError: true,
			errorMsg:    "invalid tx data",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := normalizeSwapResponse(tc.resp)
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, uint64(tc.resp.Tx.Gas), result.TxNormalized.Gas)
				assert.NotNil(t, result.TxNormalized.GasPrice)
				assert.NotNil(t, result.TxNormalized.Value)
				assert.NotNil(t, result.TxNormalized.Data)
			}
		})
	}
}

func TestNormalizeSwapResponse_DataParsing(t *testing.T) {
	resp := SwapResponse{
		Tx: TransactionData{
			To:       "0x6B175474E89094C44Da98b954EedeAC495271d0F",
			Gas:      100000,
			GasPrice: "20000000000",
			Value:    "1000000000000000000",
			Data:     "0x1234567890abcdef",
		},
	}

	result, err := normalizeSwapResponse(resp)
	require.NoError(t, err)

	// Verify the data was correctly decoded
	expectedData := []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef}
	assert.Equal(t, expectedData, result.TxNormalized.Data)

	// Verify gas price
	assert.Equal(t, "20000000000", result.TxNormalized.GasPrice.String())

	// Verify value (1 ETH in wei)
	assert.Equal(t, "1000000000000000000", result.TxNormalized.Value.String())
}

func TestNormalizeApproveCallDataResponse(t *testing.T) {
	tests := []struct {
		name        string
		resp        ApproveCallDataResponse
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid response",
			resp: ApproveCallDataResponse{
				To:       "0x6B175474E89094C44Da98b954EedeAC495271d0F",
				GasPrice: "20000000000",
				Value:    "0",
				Data:     "0x095ea7b3000000000000000000000000",
			},
			expectError: false,
		},
		{
			name: "Invalid To address",
			resp: ApproveCallDataResponse{
				To:       "invalid",
				GasPrice: "20000000000",
				Value:    "0",
				Data:     "0x",
			},
			expectError: true,
			errorMsg:    "invalid to address",
		},
		{
			name: "Invalid GasPrice",
			resp: ApproveCallDataResponse{
				To:       "0x6B175474E89094C44Da98b954EedeAC495271d0F",
				GasPrice: "invalid",
				Value:    "0",
				Data:     "0x",
			},
			expectError: true,
			errorMsg:    "invalid gas price",
		},
		{
			name: "Invalid Value",
			resp: ApproveCallDataResponse{
				To:       "0x6B175474E89094C44Da98b954EedeAC495271d0F",
				GasPrice: "20000000000",
				Value:    "invalid",
				Data:     "0x",
			},
			expectError: true,
			errorMsg:    "invalid value",
		},
		{
			name: "Invalid Data",
			resp: ApproveCallDataResponse{
				To:       "0x6B175474E89094C44Da98b954EedeAC495271d0F",
				GasPrice: "20000000000",
				Value:    "0",
				Data:     "not-hex",
			},
			expectError: true,
			errorMsg:    "invalid data",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := normalizeApproveCallDataResponse(tc.resp)
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				// Verify the gas is set to Erc20ApproveGas constant
				assert.Equal(t, uint64(constants.Erc20ApproveGas), result.TxNormalized.Gas)
			}
		})
	}
}

func TestNormalizeApproveCallDataResponse_UsesConstantGas(t *testing.T) {
	resp := ApproveCallDataResponse{
		To:       "0x6B175474E89094C44Da98b954EedeAC495271d0F",
		GasPrice: "20000000000",
		Value:    "0",
		Data:     "0x",
	}

	result, err := normalizeApproveCallDataResponse(resp)
	require.NoError(t, err)

	// Gas should always be the Erc20ApproveGas constant
	assert.Equal(t, uint64(constants.Erc20ApproveGas), result.TxNormalized.Gas)
}
