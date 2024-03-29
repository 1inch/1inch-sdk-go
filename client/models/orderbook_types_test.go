package models

import (
	"fmt"
	"testing"

	"github.com/1inch/1inch-sdk-go/internal/helpers/consts/chains"
	"github.com/1inch/1inch-sdk-go/internal/helpers/consts/tokens"

	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/internal/validate"
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
				PrivateKey:   "a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1",
				Maker:        "0x1234567890abcdef1234567890abcdef12345678",
				MakerAsset:   "0x1234567890abcdef1234567890abcdef12345678",
				TakerAsset:   "0x1234567890abcdef1234567890abcdef12345679",
				TakingAmount: "1000000000000000000",
				MakingAmount: "2000000000000000000",
			},
		},
		{
			description: "Missing required parameters",
			params:      CreateOrderParams{},
			expectErrors: []string{
				"'chainId' is required",
				"'privateKey' is required",
				"'maker' is required",
				"'makerAsset' is required",
				"'takerAsset' is required",
				"'takingAmount' is required",
				"'makingAmount' is required",
			},
		},
		{
			description: "Error - MakerAsset is native token",
			params: CreateOrderParams{
				ChainId:      chains.Ethereum,
				PrivateKey:   "a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1",
				Maker:        "0x1234567890abcdef1234567890abcdef12345678",
				MakerAsset:   tokens.NativeToken,
				TakerAsset:   tokens.EthereumUsdc,
				TakingAmount: "1000000000000000000",
				MakingAmount: "2000000000000000000",
			},
			expectErrors: []string{
				"native gas token is not supported as maker or taker asset",
			},
		},
		{
			description: "Error - TakerAsset is native token",
			params: CreateOrderParams{
				ChainId:      chains.Ethereum,
				PrivateKey:   "a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1",
				Maker:        "0x1234567890abcdef1234567890abcdef12345678",
				MakerAsset:   tokens.EthereumUsdc,
				TakerAsset:   tokens.NativeToken,
				TakingAmount: "1000000000000000000",
				MakingAmount: "2000000000000000000",
			},
			expectErrors: []string{
				"native gas token is not supported as maker or taker asset",
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
					Page:       1,
					Limit:      1,
					Statuses:   []float32{1},
					SortBy:     LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortByCreateDateTime,
					TakerAsset: "0x1234567890abcdef1234567890abcdef12345678",
					MakerAsset: "0x1234567890abcdef1234567890abcdef12345678",
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
					Page:       1,
					Limit:      1,
					Statuses:   []float32{1},
					SortBy:     LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortByCreateDateTime,
					TakerAsset: "0x1234567890abcdef1234567890abcdef12345678",
					MakerAsset: "0x1234567890abcdef1234567890abcdef12345678",
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
			},
		},
		{
			description: "Invalid chain id",
			params: GetEventsParams{
				ChainId: -1,
			},
			expectErrors: []string{
				"'chainId': is invalid",
			},
		},
		{
			description: "Invalid limit parameter",
			params: GetEventsParams{
				ChainId: chains.Ethereum,
				LimitOrderV3SubscribedApiControllerGetEventsParams: LimitOrderV3SubscribedApiControllerGetEventsParams{
					Limit: -1,
				}},
			expectErrors: []string{
				"'limit': must be greater than 0",
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
