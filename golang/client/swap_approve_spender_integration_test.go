package client

import (
	"fmt"
	"os"
	"testing"

	"dev-portal-sdk-go/client/swap"
	"dev-portal-sdk-go/helpers"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			spender, resp, err := c.Swap.ApproveSpender()
			assert.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
			assert.Equal(t, tc.expectedOutput.Address, spender.Address)

			helpers.Sleep()
		})
	}
}
