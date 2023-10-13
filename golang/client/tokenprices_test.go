package client

import (
	"fmt"
	"net/http"
	"testing"

	"dev-portal-sdk-go/client/tokenprices"
	"github.com/stretchr/testify/assert"
)

func TestGetTokenPrices(t *testing.T) {

	endpoint := "/price/v1.1/1"
	defaultResponse := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w,
			`{
				"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2": "1584.94014"
			}`,
		)
	}

	testcases := []struct {
		description              string
		handlerFunc              func(w http.ResponseWriter, r *http.Request)
		params                   tokenprices.PricesParameters
		expectedOutput           string
		expectedErrorDescription string
	}{
		{
			description: "Get prices in USD",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				// TODO Can we make these constants?
				assert.Equal(t, "USD", r.URL.Query().Get("currency"))
			},
			params: tokenprices.PricesParameters{
				Currency: tokenprices.CurrencyTypeUSD,
			},
		},
		{
			description: "Get prices in Wei",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Empty(t, r.URL.Query().Get("currency"))
			},
			params: tokenprices.PricesParameters{
				Currency: tokenprices.CurrencyTypeWEI,
			},
		},
		{
			description: "Get prices in Wei (no field)",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Empty(t, r.URL.Query().Get("currency"))
			},
		},
		{
			description: "Fail - provide invalid currency",
			params: tokenprices.PricesParameters{
				Currency: tokenprices.CurrencyType("ok"),
			},
			expectedErrorDescription: "currency value ok is not valid",
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			client, apiHandler, _, teardown := setup()
			defer teardown()

			if tc.handlerFunc != nil {
				apiHandler.HandleFunc(endpoint, tc.handlerFunc)
			} else {
				apiHandler.HandleFunc(endpoint, defaultResponse)
			}

			_, _, err := client.GetTokenPrices(tc.params)
			if tc.expectedErrorDescription != "" {
				assert.Equal(t, tc.expectedErrorDescription, err.Error())
				return
			}
			assert.NoError(t, err)
		})
	}
}
