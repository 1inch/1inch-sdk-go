package client

import (
	"fmt"
	"os"
	"testing"

	"dev-portal-sdk-go/client/tokenprices"
	"dev-portal-sdk-go/helpers"
	"github.com/stretchr/testify/assert"
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
		ApiKey:            os.Getenv("DEV_PORTAL_TOKEN"),
	})
	assert.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			priceParameters := tokenprices.ChainControllerByAddressesParams{
				Currency: &tc.currency,
			}

			message, resp, err := c.TokenPrices.GetPrices(priceParameters)
			if tc.expectedErrorCode != 0 {
				if resp == nil {
					assert.FailNow(t, "Response should not be nil")
				}
				assert.Equal(t, tc.expectedErrorCode, resp.StatusCode)
				assert.Equal(t, tc.expectedErrorDescription, err.Error())
				return
			}

			assert.NoError(t, err)
			assert.NotEqual(t, "", message)

			helpers.Sleep()
		})
	}
}
