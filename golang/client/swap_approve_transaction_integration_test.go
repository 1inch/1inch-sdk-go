package client

import (
	"fmt"
	"os"
	"testing"

	"dev-portal-sdk-go/client/swap"
	"dev-portal-sdk-go/helpers"
	"dev-portal-sdk-go/helpers/consts/tokens"
	"github.com/stretchr/testify/assert"
)

func TestApproveTransactionIntegration(t *testing.T) {

	testcases := []struct {
		description    string
		params         swap.ApproveControllerGetCallDataParams
		expectedOutput swap.ApproveCallDataResponse
	}{
		{
			description: "Get approve spender address",
			params: swap.ApproveControllerGetCallDataParams{
				TokenAddress: tokens.EthereumUsdc,
				Amount:       nil,
			},
			expectedOutput: swap.ApproveCallDataResponse{
				To: tokens.EthereumUsdc,
			},
		},
	}

	c, err := NewClient(Config{
		TargetEnvironment: EnvironmentProduction,
		ApiKey:            os.Getenv("DEV_PORTAL_TOKEN"),
	})
	assert.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			transaction, resp, err := c.Swap.ApproveTransaction(tc.params)
			assert.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
			assert.Equal(t, tc.expectedOutput.To, transaction.To)

			helpers.Sleep()
		})
	}
}
