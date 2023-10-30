package client

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"1inch-sdk-golang/client/tokenprices"
	"1inch-sdk-golang/helpers"
)

func TestGetTokenPricesIntegration(t *testing.T) {

	testcases := []struct {
		description              string
		currency                 tokenprices.ChainControllerByAddressesParamsCurrency
		expectedOutput           string
		expectedErrorCode        int
		expectedErrorDescription string
	}{
		{
			description: "Get prices in USD",
			currency:    tokenprices.USD,
		},
		{
			description: "Get prices in Wei (exclude currency field)",
		},
	}

	c, err := NewClient(Config{
		TargetEnvironment: EnvironmentProduction,
		DevPortalApiKey:   os.Getenv("DEV_PORTAL_TOKEN"),
	})
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			priceParameters := tokenprices.ChainControllerByAddressesParams{
				Currency: &tc.currency,
			}

			message, resp, err := c.TokenPrices.GetPrices(context.Background(), priceParameters)
			if tc.expectedErrorCode != 0 {
				if resp == nil {
					assert.FailNow(t, "Response should not be nil")
				}
				assert.Equal(t, tc.expectedErrorCode, resp.StatusCode)
				assert.Equal(t, tc.expectedErrorDescription, err.Error())
				return
			}

			require.NoError(t, err)
			assert.NotEqual(t, "", message)

			helpers.Sleep()
		})
	}
}
