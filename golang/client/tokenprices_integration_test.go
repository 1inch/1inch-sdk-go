package client

import (
	"fmt"
	"testing"

	"dev-portal-sdk-go/client/tokenprices"
	"github.com/stretchr/testify/assert"
)

func TestGetTokenPricesIntegration(t *testing.T) {

	testcases := []struct {
		description              string
		currency                 tokenprices.CurrencyType
		expectedOutput           string
		expectedErrorCode        int
		expectedErrorDescription string
	}{
		{
			description: "Get prices in USD",
			currency:    tokenprices.CurrencyTypeUSD,
		},
		{
			description: "Get prices in Wei",
			currency:    tokenprices.CurrencyTypeWEI,
		},
		{
			description: "Get prices in Wei (no field)",
			currency:    "",
		},
		{
			description:              "Fail - provide invalid currency",
			currency:                 tokenprices.CurrencyType("ok"),
			expectedErrorCode:        400,
			expectedErrorDescription: "Currency ok not supported",
		},
	}

	c := NewClient()

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			priceParameters := tokenprices.PricesParameters{
				Currency: tc.currency,
			}

			message, _, err := c.GetTokenPrices(priceParameters)
			if tc.expectedErrorCode != 0 {
				assert.Equal(t, tc.expectedErrorCode, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEqual(t, "", message)
		})
	}
}
