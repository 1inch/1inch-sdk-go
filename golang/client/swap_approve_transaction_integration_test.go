package client

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"1inch-sdk-golang/client/swap"
	"1inch-sdk-golang/helpers"
	"1inch-sdk-golang/helpers/consts/tokens"
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
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			transaction, resp, err := c.Swap.ApproveTransaction(context.Background(), tc.params)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
			assert.Equal(t, tc.expectedOutput.To, transaction.To)

			helpers.Sleep()
		})
	}
}
