package client

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"dev-portal-sdk-go/client/swap"
	"dev-portal-sdk-go/helpers/consts/addresses"
	"dev-portal-sdk-go/helpers/consts/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApproveAllowance(t *testing.T) {

	endpoint := "/swap/v5.2/1/approve/allowance"
	defaultResponse := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w,
			`{
    "allowance": "0"
}`,
		)
	}

	testcases := []struct {
		description              string
		handlerFunc              func(w http.ResponseWriter, r *http.Request)
		params                   swap.ApproveControllerGetAllowanceParams
		expectedOutput           string
		expectedErrorDescription string
	}{
		{
			description: "Success",
			params: swap.ApproveControllerGetAllowanceParams{
				TokenAddress:  tokens.EthereumWeth,
				WalletAddress: addresses.Vitalik,
			},
			expectedOutput: "0",
		},
		{
			description: "Error - exclude tokenAddress",
			params: swap.ApproveControllerGetAllowanceParams{
				WalletAddress: addresses.Vitalik,
			},
			expectedErrorDescription: "request validation error: tokenAddress is required",
		},
		{
			description: "Error - exclude walletAddress",
			params: swap.ApproveControllerGetAllowanceParams{
				TokenAddress: tokens.EthereumWeth,
			},
			expectedErrorDescription: "request validation error: walletAddress is required",
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			c, apiHandler, _, teardown, err := setup()
			require.NoError(t, err)
			defer teardown()

			if tc.handlerFunc != nil {
				apiHandler.HandleFunc(endpoint, tc.handlerFunc)
			} else {
				apiHandler.HandleFunc(endpoint, defaultResponse)
			}

			allowanceResponse, _, err := c.Swap.ApproveAllowance(context.Background(), tc.params)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				assert.Equal(t, tc.expectedErrorDescription, err.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "0", allowanceResponse.Allowance)
		})
	}
}
