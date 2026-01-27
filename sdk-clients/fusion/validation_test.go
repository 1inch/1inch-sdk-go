package fusion

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
				FromTokenAddress: validAddress,
				ToTokenAddress:   validAddress,
				Amount:           "1000000000000000000",
			},
			expectError: false,
		},
		{
			name: "Missing FromTokenAddress",
			params: QuoterControllerGetQuoteParamsFixed{
				ToTokenAddress: validAddress,
				Amount:         "1000000000000000000",
			},
			expectError: true,
			errorMsg:    "FromTokenAddress",
		},
		{
			name: "Missing ToTokenAddress",
			params: QuoterControllerGetQuoteParamsFixed{
				FromTokenAddress: validAddress,
				Amount:           "1000000000000000000",
			},
			expectError: true,
			errorMsg:    "ToTokenAddress",
		},
		{
			name: "Missing Amount",
			params: QuoterControllerGetQuoteParamsFixed{
				FromTokenAddress: validAddress,
				ToTokenAddress:   validAddress,
			},
			expectError: true,
			errorMsg:    "Amount",
		},
		{
			name: "Invalid Amount - not a number",
			params: QuoterControllerGetQuoteParamsFixed{
				FromTokenAddress: validAddress,
				ToTokenAddress:   validAddress,
				Amount:           "invalid",
			},
			expectError: true,
			errorMsg:    "Amount",
		},
		{
			name: "Invalid FromTokenAddress",
			params: QuoterControllerGetQuoteParamsFixed{
				FromTokenAddress: "invalid",
				ToTokenAddress:   validAddress,
				Amount:           "1000000000000000000",
			},
			expectError: true,
			errorMsg:    "FromTokenAddress",
		},
		{
			name: "Valid with permit",
			params: QuoterControllerGetQuoteParamsFixed{
				FromTokenAddress: validAddress,
				ToTokenAddress:   validAddress,
				Amount:           "1000000000000000000",
				Permit:           "0x1234567890abcdef",
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

func TestQuoterControllerGetQuoteWithCustomPresetsParamsFixed_Validate(t *testing.T) {
	validAddress := "0x6B175474E89094C44Da98b954EedeAC495271d0F"

	tests := []struct {
		name        string
		params      QuoterControllerGetQuoteWithCustomPresetsParamsFixed
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid params",
			params: QuoterControllerGetQuoteWithCustomPresetsParamsFixed{
				FromTokenAddress: validAddress,
				ToTokenAddress:   validAddress,
				Amount:           "1000000000000000000",
				WalletAddress:    validAddress,
			},
			expectError: false,
		},
		{
			name: "Missing WalletAddress",
			params: QuoterControllerGetQuoteWithCustomPresetsParamsFixed{
				FromTokenAddress: validAddress,
				ToTokenAddress:   validAddress,
				Amount:           "1000000000000000000",
			},
			expectError: true,
			errorMsg:    "WalletAddress",
		},
		{
			name:        "Missing all required fields",
			params:      QuoterControllerGetQuoteWithCustomPresetsParamsFixed{},
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

func TestPlaceOrderBody_Validate(t *testing.T) {
	validAddress := "0x6B175474E89094C44Da98b954EedeAC495271d0F"

	tests := []struct {
		name        string
		body        PlaceOrderBody
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid body",
			body: PlaceOrderBody{
				Maker:        validAddress,
				MakerAsset:   validAddress,
				MakingAmount: "1000000000000000000",
				Receiver:     validAddress,
			},
			expectError: false,
		},
		{
			name: "Missing Maker",
			body: PlaceOrderBody{
				MakerAsset:   validAddress,
				MakingAmount: "1000000000000000000",
				Receiver:     validAddress,
			},
			expectError: true,
			errorMsg:    "Maker",
		},
		{
			name: "Missing MakerAsset",
			body: PlaceOrderBody{
				Maker:        validAddress,
				MakingAmount: "1000000000000000000",
				Receiver:     validAddress,
			},
			expectError: true,
			errorMsg:    "MakerAsset",
		},
		{
			name: "Missing MakingAmount",
			body: PlaceOrderBody{
				Maker:      validAddress,
				MakerAsset: validAddress,
				Receiver:   validAddress,
			},
			expectError: true,
			errorMsg:    "MakingAmount",
		},
		{
			name: "Missing Receiver",
			body: PlaceOrderBody{
				Maker:        validAddress,
				MakerAsset:   validAddress,
				MakingAmount: "1000000000000000000",
			},
			expectError: true,
			errorMsg:    "Receiver",
		},
		{
			name: "Invalid MakingAmount",
			body: PlaceOrderBody{
				Maker:        validAddress,
				MakerAsset:   validAddress,
				MakingAmount: "invalid",
				Receiver:     validAddress,
			},
			expectError: true,
			errorMsg:    "MakingAmount",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.body.Validate()
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
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
				Receiver:         validAddress,
				WalletAddress:    validAddress,
				FromTokenAddress: validAddress,
				ToTokenAddress:   validAddress,
				Amount:           "1000000000000000000",
				Preset:           "fast",
			},
			expectError: false,
		},
		{
			name: "Missing Receiver",
			params: OrderParams{
				WalletAddress:    validAddress,
				FromTokenAddress: validAddress,
				ToTokenAddress:   validAddress,
				Amount:           "1000000000000000000",
				Preset:           "fast",
			},
			expectError: true,
			errorMsg:    "Receiver",
		},
		{
			name: "Missing WalletAddress",
			params: OrderParams{
				Receiver:         validAddress,
				FromTokenAddress: validAddress,
				ToTokenAddress:   validAddress,
				Amount:           "1000000000000000000",
				Preset:           "fast",
			},
			expectError: true,
			errorMsg:    "WalletAddress",
		},
		{
			name: "Missing Preset",
			params: OrderParams{
				Receiver:         validAddress,
				WalletAddress:    validAddress,
				FromTokenAddress: validAddress,
				ToTokenAddress:   validAddress,
				Amount:           "1000000000000000000",
			},
			expectError: true,
			errorMsg:    "Preset",
		},
		{
			name: "Empty preset",
			params: OrderParams{
				Receiver:         validAddress,
				WalletAddress:    validAddress,
				FromTokenAddress: validAddress,
				ToTokenAddress:   validAddress,
				Amount:           "1000000000000000000",
				Preset:           "",
			},
			expectError: true,
			errorMsg:    "Preset",
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
