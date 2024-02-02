package orderbook

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk/golang/client/validate"
	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
)

func TestCreateOrderParams_Validate(t *testing.T) {
	testCases := []struct {
		description  string
		params       CreateOrderParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: CreateOrderParams{
				ChainId:      chains.Ethereum,
				WalletKey:    "a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1",
				SourceWallet: "0x1234567890abcdef1234567890abcdef12345678",
				FromToken:    "0x1234567890abcdef1234567890abcdef12345678",
				ToToken:      "0x1234567890abcdef1234567890abcdef12345678",
				TakingAmount: "1000000000000000000",
				MakingAmount: "2000000000000000000",
				Receiver:     "0x1234567890abcdef1234567890abcdef12345678",
			},
		},
		{
			description: "Missing required parameters",
			params:      CreateOrderParams{},
			expectErrors: []string{
				"'chainId' is required",
				"'walletKey' is required",
				"'sourceWallet' is required",
				"'fromToken' is required",
				"'toToken' is required",
				"'takingAmount' is required",
				"'makingAmount' is required",
				"'receiver' is required",
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
				require.Equal(t, len(tc.expectErrors), validate.GetValidatorErrorsCount(err), "The number of errors returned should match the length of the expected errors")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetOrdersByCreatorAddressParams_Validate(t *testing.T) {
	testCases := []struct {
		description  string
		params       GetOrdersByCreatorAddressParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: GetOrdersByCreatorAddressParams{
				ChainId:        chains.Ethereum,
				CreatorAddress: "0x1234567890abcdef1234567890abcdef12345678",
				LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams: LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{
					Page:       helpers.GetPtr(float32(1)),
					Limit:      helpers.GetPtr(float32(1)),
					Statuses:   helpers.GetPtr([]float32{1}),
					SortBy:     helpers.GetPtr(LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortByCreateDateTime),
					TakerAsset: helpers.GetPtr("0x1234567890abcdef1234567890abcdef12345678"),
					MakerAsset: helpers.GetPtr("0x1234567890abcdef1234567890abcdef12345678"),
				},
			},
		},
		{
			description: "Missing required parameters",
			params:      GetOrdersByCreatorAddressParams{},
			expectErrors: []string{
				"'chainId' is required",
				"'creatorAddress' is required",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := tc.params.Validate()

			fmt.Printf("Errors: %v\n", err)

			if len(tc.expectErrors) > 0 {
				require.Error(t, err)
				for _, expectedError := range tc.expectErrors {
					require.Contains(t, err.Error(), expectedError, "Error message should contain the expected text")
				}
				require.Equal(t, len(tc.expectErrors), validate.GetValidatorErrorsCount(err), "The number of errors returned should match the length of the expected errors")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetAllOrdersParams_Validate(t *testing.T) {
	testCases := []struct {
		description  string
		params       GetAllOrdersParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: GetAllOrdersParams{
				ChainId: chains.Ethereum,
				LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams: LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{
					Page:       helpers.GetPtr(float32(1)),
					Limit:      helpers.GetPtr(float32(1)),
					Statuses:   helpers.GetPtr([]float32{1}),
					SortBy:     helpers.GetPtr(LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortByCreateDateTime),
					TakerAsset: helpers.GetPtr("0x1234567890abcdef1234567890abcdef12345678"),
					MakerAsset: helpers.GetPtr("0x1234567890abcdef1234567890abcdef12345678"),
				},
			},
		},
		{
			description: "Missing required parameters",
			params:      GetAllOrdersParams{},
			expectErrors: []string{
				"'chainId' is required",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := tc.params.Validate()

			fmt.Printf("Errors: %v\n", err)

			if len(tc.expectErrors) > 0 {
				require.Error(t, err)
				for _, expectedError := range tc.expectErrors {
					require.Contains(t, err.Error(), expectedError, "Error message should contain the expected text")
				}
				require.Equal(t, len(tc.expectErrors), validate.GetValidatorErrorsCount(err), "The number of errors returned should match the length of the expected errors")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetCountParams_Validate(t *testing.T) {
	testCases := []struct {
		description  string
		params       GetCountParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: GetCountParams{
				ChainId: chains.Ethereum,
				LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams: LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams{
					Statuses: []string{"1"},
				},
			},
		},
		{
			description: "Missing required parameters",
			params:      GetCountParams{},
			expectErrors: []string{
				"'chainId' is required",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := tc.params.Validate()

			fmt.Printf("Errors: %v\n", err)

			if len(tc.expectErrors) > 0 {
				require.Error(t, err)
				for _, expectedError := range tc.expectErrors {
					require.Contains(t, err.Error(), expectedError, "Error message should contain the expected text")
				}
				require.Equal(t, len(tc.expectErrors), validate.GetValidatorErrorsCount(err), "The number of errors returned should match the length of the expected errors")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetEventParams_Validate(t *testing.T) {
	testCases := []struct {
		description  string
		params       GetEventParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: GetEventParams{
				ChainId:   chains.Ethereum,
				OrderHash: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef12",
			},
		},
		{
			description: "Missing required parameters",
			params:      GetEventParams{},
			expectErrors: []string{
				"'chainId' is required",
				"'orderHash' is required",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := tc.params.Validate()

			fmt.Printf("Errors: %v\n", err)

			if len(tc.expectErrors) > 0 {
				require.Error(t, err)
				for _, expectedError := range tc.expectErrors {
					require.Contains(t, err.Error(), expectedError, "Error message should contain the expected text")
				}
				require.Equal(t, len(tc.expectErrors), validate.GetValidatorErrorsCount(err), "The number of errors returned should match the length of the expected errors")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetEventsParams_Validate(t *testing.T) {
	testCases := []struct {
		description  string
		params       GetEventsParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: GetEventsParams{
				ChainId: chains.Ethereum,
				LimitOrderV3SubscribedApiControllerGetEventsParams: LimitOrderV3SubscribedApiControllerGetEventsParams{
					Limit: 1,
				},
			},
		},
		{
			description: "Missing required parameters",
			params:      GetEventsParams{},
			expectErrors: []string{
				"'chainId' is required",
				"'limit': must be greater than 0", // TODO is this what I want to check here?
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := tc.params.Validate()

			fmt.Printf("Errors: %v\n", err)

			if len(tc.expectErrors) > 0 {
				require.Error(t, err)
				for _, expectedError := range tc.expectErrors {
					require.Contains(t, err.Error(), expectedError, "Error message should contain the expected text")
				}
				require.Equal(t, len(tc.expectErrors), validate.GetValidatorErrorsCount(err), "The number of errors returned should match the length of the expected errors")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetActiveOrdersWithPermitParams_Validate(t *testing.T) {
	testCases := []struct {
		description  string
		params       GetActiveOrdersWithPermitParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: GetActiveOrdersWithPermitParams{
				ChainId: chains.Ethereum,
				Wallet:  "a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1",
				Token:   "0x1234567890abcdef1234567890abcdef12345678",
			},
		},
		{
			description: "Missing required parameters",
			params:      GetActiveOrdersWithPermitParams{},
			expectErrors: []string{
				"'chainId' is required",
				"'wallet' is required",
				"'token' is required",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := tc.params.Validate()

			fmt.Printf("Errors: %v\n", err)

			if len(tc.expectErrors) > 0 {
				require.Error(t, err)
				for _, expectedError := range tc.expectErrors {
					require.Contains(t, err.Error(), expectedError, "Error message should contain the expected text")
				}
				require.Equal(t, len(tc.expectErrors), validate.GetValidatorErrorsCount(err), "The number of errors returned should match the length of the expected errors")
			} else {
				require.NoError(t, err)
			}
		})
	}
}
