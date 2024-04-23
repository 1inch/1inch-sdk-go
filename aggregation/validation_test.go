package aggregation

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/internal/validate"
)

const ethereumUsdc = "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"

func TestSwapTokensParams_Validate(t *testing.T) {
	testCases := []struct {
		description  string
		params       GetSwapParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: GetSwapParams{
				Src:      "0x1234567890abcdef1234567890abcdef12345678",
				Dst:      "0x1234567890abcdef1234567890abcdef12345679",
				Amount:   "10000",
				From:     "0x1234567890abcdef1234567890abcdef12345678",
				Slippage: 0.5,
			},
		},
		{
			description: "Missing required parameters",
			params:      GetSwapParams{},
			expectErrors: []string{
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
		params       GetAllowanceParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: GetAllowanceParams{
				TokenAddress:  "0x1234567890abcdef1234567890abcdef12345678",
				WalletAddress: "0x1234567890abcdef1234567890abcdef12345678",
			},
		},
		{
			description: "Missing required parameters",
			params:      GetAllowanceParams{},
			expectErrors: []string{
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

func TestApproveTransactionParams_Validate(t *testing.T) {
	testCases := []struct {
		description  string
		params       GetApproveParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: GetApproveParams{
				TokenAddress: "0x1234567890abcdef1234567890abcdef12345678",
			},
		},
		{
			description: "Missing required parameters",
			params:      GetApproveParams{},
			expectErrors: []string{
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

func TestGetQuoteParams_Validate(t *testing.T) {
	testCases := []struct {
		description  string
		params       GetQuoteParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: GetQuoteParams{
				Src:    "0x1234567890abcdef1234567890abcdef12345678",
				Dst:    "0x1234567890abcdef1234567890abcdef12345679",
				Amount: "10000",
			},
		},
		{
			description: "Missing required parameters",
			params:      GetQuoteParams{},
			expectErrors: []string{
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
		params       GetSwapParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: GetSwapParams{
				Src:      "0x1234567890abcdef1234567890abcdef12345678",
				Dst:      "0x1234567890abcdef1234567890abcdef12345679",
				Amount:   "10000",
				From:     "0x1234567890abcdef1234567890abcdef12345678",
				Slippage: 0.5,
			},
		},
		{
			description: "Missing required parameters",
			params:      GetSwapParams{},
			expectErrors: []string{
				"'src' is required",
				"'dst' is required",
				"'amount' is required",
				"'from' is required",
				"'slippage' is required",
			},
		},
		{
			description: "Error - src and dst tokens are identical",
			params: GetSwapParams{
				Src:      ethereumUsdc,
				Dst:      ethereumUsdc,
				Amount:   "10000",
				From:     "0x1234567890abcdef1234567890abcdef12345678",
				Slippage: 0.5,
			},
			expectErrors: []string{
				"src and dst tokens must be different",
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
