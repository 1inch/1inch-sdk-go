package client

import (
	"context"
	"fmt"
	"os"
	"testing"

	"dev-portal-sdk-go/client/swap"
	"dev-portal-sdk-go/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApproveSpenderIntegration(t *testing.T) {

	testcases := []struct {
		description    string
		expectedOutput swap.SpenderResponse
	}{
		{
			description: "Get approve spender address",
			expectedOutput: swap.SpenderResponse{
				Address: "0x1111111254eeb25477b68fb85ed929f73a960582",
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

			spender, resp, err := c.Swap.ApproveSpender(context.Background())
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
			assert.Equal(t, tc.expectedOutput.Address, spender.Address)

			helpers.Sleep()
		})
	}
}
