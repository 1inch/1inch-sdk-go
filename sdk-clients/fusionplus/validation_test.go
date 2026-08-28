package fusionplus

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/v4/constants"
)

func TestOrderApiControllerGetActiveOrdersParams_Validate(t *testing.T) {
	tests := []struct {
		name        string
		params      GetActiveOrdersParams
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid params - empty (defaults)",
			params:      GetActiveOrdersParams{},
			expectError: false,
		},
		{
			name: "Valid params - with page and limit",
			params: GetActiveOrdersParams{
				Page:  1,
				Limit: 10,
			},
			expectError: false,
		},
		{
			name: "Invalid page - negative",
			params: GetActiveOrdersParams{
				Page: -1,
			},
			expectError: true,
			errorMsg:    "Page",
		},
		{
			name: "Invalid limit - negative",
			params: GetActiveOrdersParams{
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

func TestQuoterControllerGetQuoteParams_Validate(t *testing.T) {
	validAddress := "0x6B175474E89094C44Da98b954EedeAC495271d0F"

	tests := []struct {
		name        string
		params      QuoteParams
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid params",
			params: QuoteParams{
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
			params: QuoteParams{
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
			params: QuoteParams{
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
			params: QuoteParams{
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
			params: QuoteParams{
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
			// Aurora's chain id (1313161554) exceeds float32's 24-bit integer precision;
			// it must survive validation unrounded now that chain ids are ints.
			name: "Valid params - Aurora chain id above 2^24",
			params: QuoteParams{
				SrcTokenAddress: validAddress,
				DstTokenAddress: validAddress,
				WalletAddress:   validAddress,
				SrcChain:        constants.AuroraChainId,
				DstChain:        constants.AuroraChainId,
				Amount:          "1000000000000000000",
			},
			expectError: false,
		},
		{
			name: "Invalid SrcChain - zero",
			params: QuoteParams{
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
			params: QuoteParams{
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

func TestQuoterControllerGetQuoteWithCustomPresetsParams_Validate(t *testing.T) {
	validAddress := "0x6B175474E89094C44Da98b954EedeAC495271d0F"

	tests := []struct {
		name        string
		params      CustomPresetQuoteParams
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid params",
			params: CustomPresetQuoteParams{
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
			params: CustomPresetQuoteParams{
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
			name: "Missing Amount",
			params: CustomPresetQuoteParams{
				SrcTokenAddress: validAddress,
				DstTokenAddress: validAddress,
				WalletAddress:   validAddress,
				SrcChain:        1,
				DstChain:        137,
				Amount:          "",
			},
			expectError: true,
			errorMsg:    "Amount",
		},
		{
			name: "Invalid chain",
			params: CustomPresetQuoteParams{
				SrcTokenAddress: validAddress,
				DstTokenAddress: validAddress,
				WalletAddress:   validAddress,
				SrcChain:        999999,
				DstChain:        137,
				Amount:          "1000000000000000000",
			},
			expectError: true,
			errorMsg:    "SrcChain",
		},
		{
			name:        "Missing all required fields",
			params:      CustomPresetQuoteParams{},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.params.Validate()
			if tc.expectError {
				require.Error(t, err)
				if tc.errorMsg != "" {
					assert.Contains(t, err.Error(), tc.errorMsg)
				}
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
			name: "Compact permit2 form rejected",
			params: OrderParams{
				Receiver:  validAddress,
				Preset:    "fast",
				Permit:    "0x" + strings.Repeat("11", 96),
				IsPermit2: true,
			},
			expectError: true,
			errorMsg:    "compact permit2",
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
