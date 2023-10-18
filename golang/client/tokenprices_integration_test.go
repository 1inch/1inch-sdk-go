package client

import (
	"fmt"
	"testing"

	"dev-portal-sdk-go/client/spotprice"
	"dev-portal-sdk-go/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGetTokenPricesIntegration(t *testing.T) {

	testcases := []struct {
		description              string
		currency                 spotprice.ChainControllerByAddressesParamsCurrency
		expectedOutput           string
		expectedErrorCode        int
		expectedErrorDescription string
	}{
		{
			description: "Get prices in USD",
			currency:    spotprice.USD,
		},
		{
			description: "Get prices in Wei (no field)",
			currency:    "",
		},
	}

	c := NewClient(nil)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			priceParameters := spotprice.ChainControllerByAddressesParams{
				Currency: &tc.currency,
			}

			message, resp, err := c.GetTokenPrices(priceParameters)
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
