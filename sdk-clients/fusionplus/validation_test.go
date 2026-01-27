package fusionplus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderApiControllerGetActiveOrdersParams_Validate(t *testing.T) {
	tests := []struct {
		name        string
		params      OrderApiControllerGetActiveOrdersParams
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid params - empty (defaults)",
			params:      OrderApiControllerGetActiveOrdersParams{},
			expectError: false,
		},
		{
			name: "Valid params - with page and limit",
			params: OrderApiControllerGetActiveOrdersParams{
				Page:  1,
				Limit: 10,
			},
			expectError: false,
		},
		{
			name: "Invalid page - negative",
			params: OrderApiControllerGetActiveOrdersParams{
				Page: -1,
			},
			expectError: true,
			errorMsg:    "Page",
		},
		{
			name: "Invalid limit - negative",
			params: OrderApiControllerGetActiveOrdersParams{
				Limit: -1,
			},
			expectError: true,
			errorMsg:    "Limit",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.params.Validate()
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestQuoterControllerGetQuoteParamsFixed_Validate(t *testing.T) {
	validAddress := "0x6B175474E89094C44Da98b954EedeAC495271d0F"

	tests := []struct {
		name        string
		params      QuoterControllerGetQuoteParamsFixed
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid params",
			params: QuoterControllerGetQuoteParamsFixed{
				SrcTokenAddress: validAddress,
				DstTokenAddress: validAddress,
				WalletAddress:   validAddress,
				SrcChain:        1,
				DstChain:        137,
				Amount:          "1000000000000000000",
			},
			expectError: false,
		},
		{
			name: "Missing SrcTokenAddress",
			params: QuoterControllerGetQuoteParamsFixed{
				DstTokenAddress: validAddress,
				WalletAddress:   validAddress,
				SrcChain:        1,
				DstChain:        137,
				Amount:          "1000000000000000000",
			},
			expectError: true,
			errorMsg:    "SrcTokenAddress",
		},
		{
			name: "Missing WalletAddress",
			params: QuoterControllerGetQuoteParamsFixed{
				SrcTokenAddress: validAddress,
				DstTokenAddress: validAddress,
				SrcChain:        1,
				DstChain:        137,
				Amount:          "1000000000000000000",
			},
			expectError: true,
			errorMsg:    "WalletAddress",
		},
		{
			name: "Missing Amount",
			params: QuoterControllerGetQuoteParamsFixed{
				SrcTokenAddress: validAddress,
				DstTokenAddress: validAddress,
				WalletAddress:   validAddress,
				SrcChain:        1,
				DstChain:        137,
			},
			expectError: true,
			errorMsg:    "Amount",
		},
		{
			name: "Invalid Amount",
			params: QuoterControllerGetQuoteParamsFixed{
				SrcTokenAddress: validAddress,
				DstTokenAddress: validAddress,
				WalletAddress:   validAddress,
				SrcChain:        1,
				DstChain:        137,
				Amount:          "invalid",
			},
			expectError: true,
			errorMsg:    "Amount",
		},
		{
			name: "Invalid SrcChain - zero",
			params: QuoterControllerGetQuoteParamsFixed{
				SrcTokenAddress: validAddress,
				DstTokenAddress: validAddress,
				WalletAddress:   validAddress,
				SrcChain:        0,
				DstChain:        137,
				Amount:          "1000000000000000000",
			},
			expectError: true,
			errorMsg:    "SrcChain",
		},
		{
			name: "Invalid DstChain - zero",
			params: QuoterControllerGetQuoteParamsFixed{
				SrcTokenAddress: validAddress,
				DstTokenAddress: validAddress,
				WalletAddress:   validAddress,
				SrcChain:        1,
				DstChain:        0,
				Amount:          "1000000000000000000",
			},
			expectError: true,
			errorMsg:    "DstChain",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.params.Validate()
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// NOTE: QuoterControllerGetQuoteWithCustomPresetsParams has a type mismatch bug - the Amount field
// is float32 but the validation function CheckBigIntRequired expects a string. This causes
// validation to always fail for this field.
func TestQuoterControllerGetQuoteWithCustomPresetsParams_Validate(t *testing.T) {
	validAddress := "0x6B175474E89094C44Da98b954EedeAC495271d0F"

	tests := []struct {
		name        string
		params      QuoterControllerGetQuoteWithCustomPresetsParams
		expectError bool
		errorMsg    string
	}{
		// This test documents the validation bug - the Amount field is float32
		// but validation expects string for BigInt
		{
			name: "Type mismatch causes validation error for Amount field",
			params: QuoterControllerGetQuoteWithCustomPresetsParams{
				SrcTokenAddress: validAddress,
				DstTokenAddress: validAddress,
				WalletAddress:   validAddress,
				SrcChain:        1,
				DstChain:        137,
				Amount:          1000000000000000000,
			},
			expectError: true,
			errorMsg:    "must be a string",
		},
		{
			name:        "Missing all required fields",
			params:      QuoterControllerGetQuoteWithCustomPresetsParams{},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.params.Validate()
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOrderParams_Validate(t *testing.T) {
	validAddress := "0x6B175474E89094C44Da98b954EedeAC495271d0F"

	tests := []struct {
		name        string
		params      OrderParams
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid params",
			params: OrderParams{
				Receiver: validAddress,
				Preset:   "fast",
			},
			expectError: false,
		},
		{
			name: "Missing Receiver",
			params: OrderParams{
				Preset: "fast",
			},
			expectError: true,
			errorMsg:    "Receiver",
		},
		{
			name: "Missing Preset",
			params: OrderParams{
				Receiver: validAddress,
			},
			expectError: true,
			errorMsg:    "Preset",
		},
		{
			name: "Empty Preset",
			params: OrderParams{
				Receiver: validAddress,
				Preset:   "",
			},
			expectError: true,
			errorMsg:    "Preset",
		},
		{
			name: "Valid with permit",
			params: OrderParams{
				Receiver: validAddress,
				Preset:   "medium",
				Permit:   "0x1234567890abcdef",
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.params.Validate()
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
