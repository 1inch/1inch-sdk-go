package client

import (
	"fmt"
	"net/http"
	"testing"

	"dev-portal-sdk-go/client/swap"
	"dev-portal-sdk-go/helpers/consts/tokens"
	"github.com/stretchr/testify/assert"
)

func TestApproveTransaction(t *testing.T) {

	endpoint := "/swap/v5.2/1/approve/transaction"
	defaultResponse := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w,
			`{
    "data": "0x095ea7b30000000000000000000000001111111254eeb25477b68fb85ed929f73a960582ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
    "gasPrice": "12161498847",
    "to": "`+tokens.EthereumUsdc+`",
    "value": "0"
}`,
		)
	}

	testcases := []struct {
		description              string
		handlerFunc              func(w http.ResponseWriter, r *http.Request)
		params                   swap.ApproveControllerGetCallDataParams
		expectedOutput           swap.ApproveCallDataResponse
		expectedErrorDescription string
	}{
		{
			description: "Success",
			params: swap.ApproveControllerGetCallDataParams{
				TokenAddress: tokens.EthereumUsdc,
			},
			expectedOutput: swap.ApproveCallDataResponse{
				To: tokens.EthereumUsdc,
			},
		},
		{
			description:              "Error - exclude tokenAddress",
			params:                   swap.ApproveControllerGetCallDataParams{},
			expectedErrorDescription: "request validation error: tokenAddress is required",
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			client, apiHandler, _, teardown, err := setup()
			assert.NoError(t, err)
			defer teardown()

			if tc.handlerFunc != nil {
				apiHandler.HandleFunc(endpoint, tc.handlerFunc)
			} else {
				apiHandler.HandleFunc(endpoint, defaultResponse)
			}

			approveTransactionResponse, _, err := client.ApproveTransaction(tc.params)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				assert.Equal(t, tc.expectedErrorDescription, err.Error())
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedOutput.To, approveTransactionResponse.To)
		})
	}
}
