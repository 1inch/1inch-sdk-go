package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk/golang/client/tokenprices"
	"github.com/1inch/1inch-sdk/golang/helpers"
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

	c, err := NewClient(SimpleEthereumConfig)
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			t.Cleanup(func() {
				helpers.Sleep()
			})

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
		})
	}
}
