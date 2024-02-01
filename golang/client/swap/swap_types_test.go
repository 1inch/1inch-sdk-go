package swap

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk/golang/client/validate"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
)

func TestSwapTokensParams_Validate(t *testing.T) {
	testCases := []struct {
		description  string
		params       SwapTokensParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: SwapTokensParams{
				ChainId:       chains.Ethereum,
				PublicAddress: "0x1234567890abcdef1234567890abcdef12345678",
				WalletKey:     "a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1",
				AggregationControllerGetSwapParams: AggregationControllerGetSwapParams{
					Src:      "0x1234567890abcdef1234567890abcdef12345678",
					Dst:      "0x1234567890abcdef1234567890abcdef12345679",
					Amount:   "10000",
					From:     "0x1234567890abcdef1234567890abcdef12345678",
					Slippage: 0.5,
				},
			},
		},
		{
			description: "Missing required parameters",
			params:      SwapTokensParams{},
			expectErrors: []string{
				"'chainId' is required",
				"'publicAddress' is required",
				"'walletKey' is required",
				"'src' is required",
				"'dst' is required",
				"'amount' is required",
				"'from' is required",
				"'slippage' is required",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := tc.params.Validate()

			if len(tc.expectErrors) > 0 {
				require.Error(t, err)
				for _, expectedError := range tc.expectErrors {
					require.Contains(t, err.Error(), expectedError, "Error message should contain the expected text")
				}
				require.Equal(t, len(tc.expectErrors), validate.GetValidatorErrorsCount(err), "The number of errors returned should match the length of the expected errors: %s\n", err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestApproveAllowanceParams_Validate(t *testing.T) {
	testCases := []struct {
		description  string
		params       ApproveAllowanceParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: ApproveAllowanceParams{
				ChainId: chains.Ethereum,
				ApproveControllerGetAllowanceParams: ApproveControllerGetAllowanceParams{
					TokenAddress:  "0x1234567890abcdef1234567890abcdef12345678",
					WalletAddress: "0x1234567890abcdef1234567890abcdef12345678",
				},
			},
		},
		{
			description: "Missing required parameters",
			params:      ApproveAllowanceParams{},
			expectErrors: []string{
				"'chainId' is required",
				"'tokenAddress' is required",
				"'walletAddress' is required",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := tc.params.Validate()

			if len(tc.expectErrors) > 0 {
				require.Error(t, err)
				for _, expectedError := range tc.expectErrors {
					require.Contains(t, err.Error(), expectedError, "Error message should contain the expected text")
				}
				require.Equal(t, len(tc.expectErrors), validate.GetValidatorErrorsCount(err), "The number of errors returned should match the length of the expected errors: %s\n", err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestApproveSpenderParams_Validate(t *testing.T) {
	testCases := []struct {
		description  string
		params       ApproveSpenderParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: ApproveSpenderParams{
				ChainId: chains.Ethereum,
			},
		},
		{
			description: "Missing required parameters",
			params:      ApproveSpenderParams{},
			expectErrors: []string{
				"'chainId' is required",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := tc.params.Validate()

			if len(tc.expectErrors) > 0 {
				require.Error(t, err)
				for _, expectedError := range tc.expectErrors {
					require.Contains(t, err.Error(), expectedError, "Error message should contain the expected text")
				}
				require.Equal(t, len(tc.expectErrors), validate.GetValidatorErrorsCount(err), "The number of errors returned should match the length of the expected errors: %s\n", err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestApproveTransactionParams_Validate(t *testing.T) {
	testCases := []struct {
		description  string
		params       ApproveTransactionParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: ApproveTransactionParams{
				ChainId: chains.Ethereum,
				ApproveControllerGetCallDataParams: ApproveControllerGetCallDataParams{
					TokenAddress: "0x1234567890abcdef1234567890abcdef12345678",
				},
			},
		},
		{
			description: "Missing required parameters",
			params:      ApproveTransactionParams{},
			expectErrors: []string{
				"'chainId' is required",
				"'tokenAddress' is required",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := tc.params.Validate()

			if len(tc.expectErrors) > 0 {
				require.Error(t, err)
				for _, expectedError := range tc.expectErrors {
					require.Contains(t, err.Error(), expectedError, "Error message should contain the expected text")
				}
				require.Equal(t, len(tc.expectErrors), validate.GetValidatorErrorsCount(err), "The number of errors returned should match the length of the expected errors: %s\n", err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetLiquiditySourcesParams_Validate(t *testing.T) {
	testCases := []struct {
		description  string
		params       GetLiquiditySourcesParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: GetLiquiditySourcesParams{
				ChainId: chains.Ethereum,
			},
		},
		{
			description: "Missing required parameters",
			params:      GetLiquiditySourcesParams{},
			expectErrors: []string{
				"'chainId' is required",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := tc.params.Validate()

			if len(tc.expectErrors) > 0 {
				require.Error(t, err)
				for _, expectedError := range tc.expectErrors {
					require.Contains(t, err.Error(), expectedError, "Error message should contain the expected text")
				}
				require.Equal(t, len(tc.expectErrors), validate.GetValidatorErrorsCount(err), "The number of errors returned should match the length of the expected errors: %s\n", err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetQuoteParams_Validate(t *testing.T) {
	testCases := []struct {
		description  string
		params       GetQuoteParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: GetQuoteParams{
				ChainId: chains.Ethereum,
				AggregationControllerGetQuoteParams: AggregationControllerGetQuoteParams{
					Src:    "0x1234567890abcdef1234567890abcdef12345678",
					Dst:    "0x1234567890abcdef1234567890abcdef12345679",
					Amount: "10000",
				},
			},
		},
		{
			description: "Missing required parameters",
			params:      GetQuoteParams{},
			expectErrors: []string{
				"'chainId' is required",
				"'src' is required",
				"'dst' is required",
				"'amount' is required",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := tc.params.Validate()

			if len(tc.expectErrors) > 0 {
				require.Error(t, err)
				for _, expectedError := range tc.expectErrors {
					require.Contains(t, err.Error(), expectedError, "Error message should contain the expected text")
				}
				require.Equal(t, len(tc.expectErrors), validate.GetValidatorErrorsCount(err), "The number of errors returned should match the length of the expected errors: %s\n", err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetSwapDataParams_Validate(t *testing.T) {
	testCases := []struct {
		description  string
		params       GetSwapDataParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: GetSwapDataParams{
				ChainId: chains.Ethereum,
				AggregationControllerGetSwapParams: AggregationControllerGetSwapParams{
					Src:      "0x1234567890abcdef1234567890abcdef12345678",
					Dst:      "0x1234567890abcdef1234567890abcdef12345679",
					Amount:   "10000",
					From:     "0x1234567890abcdef1234567890abcdef12345678",
					Slippage: 0.5,
				},
			},
		},
		{
			description: "Missing required parameters",
			params:      GetSwapDataParams{},
			expectErrors: []string{
				"'chainId' is required",
				"'src' is required",
				"'dst' is required",
				"'amount' is required",
				"'from' is required",
				"'slippage' is required",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := tc.params.Validate()

			if len(tc.expectErrors) > 0 {
				require.Error(t, err)
				for _, expectedError := range tc.expectErrors {
					require.Contains(t, err.Error(), expectedError, "Error message should contain the expected text")
				}
				require.Equal(t, len(tc.expectErrors), validate.GetValidatorErrorsCount(err), "The number of errors returned should match the length of the expected errors: %s\n", err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetTokensParams_Validate(t *testing.T) {
	testCases := []struct {
		description  string
		params       GetTokensParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: GetTokensParams{
				ChainId: chains.Ethereum,
			},
		},
		{
			description: "Missing required parameters",
			params:      GetTokensParams{},
			expectErrors: []string{
				"'chainId' is required",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := tc.params.Validate()

			if len(tc.expectErrors) > 0 {
				require.Error(t, err)
				for _, expectedError := range tc.expectErrors {
					require.Contains(t, err.Error(), expectedError, "Error message should contain the expected text")
				}
				require.Equal(t, len(tc.expectErrors), validate.GetValidatorErrorsCount(err), "The number of errors returned should match the length of the expected errors: %s\n", err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
