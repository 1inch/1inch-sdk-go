package orderbook

import (
	"fmt"
	"testing"

	web3_provider "github.com/1inch/1inch-sdk-go/v5/internal/web3-provider"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/v5/internal/validate"
)

var wallet, _ = web3_provider.DefaultWalletOnlyProvider("965e092fdfc08940d2bd05c7b5c7e1c51e283e92c7f52bbf1408973ae9a9acb7", 137)

func TestCreateOrderParams_Validate(t *testing.T) {
	testCases := []struct {
		description  string
		params       CreateOrderParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: CreateOrderParams{
				Wallet:       wallet,
				Maker:        "0x1234567890abcdef1234567890abcdef12345678",
				MakerAsset:   "0x1234567890abcdef1234567890abcdef12345678",
				TakerAsset:   "0x1234567890abcdef1234567890abcdef12345679",
				TakingAmount: "1000000000000000000",
				MakingAmount: "2000000000000000000",
				Taker:        "0x1234567890abcdef1234567890abcdef12345678",
			},
		},
		{
			description: "Missing required parameters",
			params:      CreateOrderParams{},
			expectErrors: []string{
				"'maker' is required",
				"'makerAsset' is required",
				"'takerAsset' is required",
				"'takingAmount' is required",
				"'makingAmount' is required",
				"'taker' is required",
			},
		},
		{
			description: "Error - MakerAsset is native token",
			params: CreateOrderParams{
				Wallet:       wallet,
				Maker:        "0x1234567890abcdef1234567890abcdef12345678",
				MakerAsset:   "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
				TakerAsset:   "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				TakingAmount: "1000000000000000000",
				MakingAmount: "2000000000000000000",
				Taker:        "0x1234567890abcdef1234567890abcdef12345678",
			},
			expectErrors: []string{
				"unsupported: native gas token as maker or taker asset",
			},
		},
		{
			description: "Error - TakerAsset is native token",
			params: CreateOrderParams{
				Wallet:       wallet,
				Maker:        "0x1234567890abcdef1234567890abcdef12345678",
				MakerAsset:   "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				TakerAsset:   "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
				TakingAmount: "1000000000000000000",
				MakingAmount: "2000000000000000000",
				Taker:        "0x1234567890abcdef1234567890abcdef12345678",
			},
			expectErrors: []string{
				"unsupported: native gas token as maker or taker asset",
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
				CreatorAddress: "0x1234567890abcdef1234567890abcdef12345678",
				LimitOrdersQueryParams: LimitOrdersQueryParams{
					Page:       1,
					Limit:      1,
					Statuses:   []float32{1},
					SortBy:     LimitOrdersQueryParamsSortByCreateDateTime,
					TakerAsset: "0x1234567890abcdef1234567890abcdef12345678",
					MakerAsset: "0x1234567890abcdef1234567890abcdef12345678",
				},
			},
		},
		{
			description: "Missing required parameters",
			params:      GetOrdersByCreatorAddressParams{},
			expectErrors: []string{
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
				LimitOrdersQueryParams: LimitOrdersQueryParams{
					Page:       1,
					Limit:      1,
					Statuses:   []float32{1},
					SortBy:     LimitOrdersQueryParamsSortByCreateDateTime,
					TakerAsset: "0x1234567890abcdef1234567890abcdef12345678",
					MakerAsset: "0x1234567890abcdef1234567890abcdef12345678",
				},
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

func TestGetOrderParams_Validate(t *testing.T) {
	testCases := []struct {
		description  string
		params       GetOrderParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: GetOrderParams{
				OrderHash: "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			},
		},
		{
			description: "Missing required parameters",
			params:      GetOrderParams{},
			expectErrors: []string{
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
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetFeeInfoParams_Validate(t *testing.T) {
	validAddress := "0x1234567890abcdef1234567890abcdef12345678"

	testCases := []struct {
		description  string
		params       GetFeeInfoParams
		expectErrors []string
	}{
		{
			description: "Valid parameters",
			params: GetFeeInfoParams{
				MakerAmount: "1000000000000000000",
				MakerAsset:  validAddress,
				TakerAmount: "1000000000000000000",
				TakerAsset:  "0xabcdef1234567890abcdef1234567890abcdef12",
			},
		},
		{
			description: "Missing required parameters",
			params:      GetFeeInfoParams{},
			expectErrors: []string{
				"'makerAmount' is required",
				"'makerAsset' is required",
				"'takerAmount' is required",
				"'takerAsset' is required",
			},
		},
		{
			description: "Invalid MakerAmount",
			params: GetFeeInfoParams{
				MakerAmount: "invalid",
				MakerAsset:  validAddress,
				TakerAmount: "1000000000000000000",
				TakerAsset:  "0xabcdef1234567890abcdef1234567890abcdef12",
			},
			expectErrors: []string{
				"'makerAmount'",
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

// NOTE: GetOrderCountParams.Statuses has a type mismatch bug - the field is []OrderStatus
// but the validation function CheckStatusesStrings expects []string. This causes validation
// to always fail for this field. The test below documents this behavior.
func TestGetOrderCountParams_Validate(t *testing.T) {
	validAddress := "0x1234567890abcdef1234567890abcdef12345678"

	testCases := []struct {
		description  string
		params       GetOrderCountParams
		expectErrors []string
	}{
		{
			description: "Valid params with all statuses",
			params: GetOrderCountParams{
				Statuses:   []OrderStatus{ValidOrders, TemporarilyInvalidOrders, InvalidOrders},
				MakerAsset: validAddress,
				TakerAsset: "0xabcdef1234567890abcdef1234567890abcdef12",
			},
			expectErrors: nil,
		},
		{
			description: "Valid params with single status",
			params: GetOrderCountParams{
				Statuses:   []OrderStatus{ValidOrders},
				MakerAsset: validAddress,
				TakerAsset: "0xabcdef1234567890abcdef1234567890abcdef12",
			},
			expectErrors: nil,
		},
		{
			description: "Invalid status value",
			params: GetOrderCountParams{
				Statuses:   []OrderStatus{ValidOrders, 4}, // 4 is not a valid status
				MakerAsset: validAddress,
				TakerAsset: "0xabcdef1234567890abcdef1234567890abcdef12",
			},
			expectErrors: []string{"statuses"},
		},
		{
			description: "Duplicate statuses",
			params: GetOrderCountParams{
				Statuses:   []OrderStatus{ValidOrders, ValidOrders},
				MakerAsset: validAddress,
				TakerAsset: "0xabcdef1234567890abcdef1234567890abcdef12",
			},
			expectErrors: []string{"duplicates"},
		},
		{
			description: "Missing MakerAsset",
			params: GetOrderCountParams{
				Statuses:   []OrderStatus{ValidOrders},
				TakerAsset: validAddress,
			},
			expectErrors: []string{"makerAsset"},
		},
		{
			description: "Missing TakerAsset",
			params: GetOrderCountParams{
				Statuses:   []OrderStatus{ValidOrders},
				MakerAsset: validAddress,
			},
			expectErrors: []string{"takerAsset"},
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
			} else {
				require.NoError(t, err)
			}
		})
	}
}
