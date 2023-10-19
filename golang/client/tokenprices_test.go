package client

import (
	"fmt"
	"net/http"
	"testing"

	"dev-portal-sdk-go/client/spotprice"
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
		params                   spotprice.ChainControllerByAddressesParams
		expectedOutput           string
		expectedErrorDescription string
	}{
		{
			description: "Get prices in USD",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, string(spotprice.USD), r.URL.Query().Get("currency"))
			},
			params: spotprice.ChainControllerByAddressesParams{
				Currency: spotprice.GetCurrencyType(spotprice.USD),
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
			params: spotprice.ChainControllerByAddressesParams{
				Currency: spotprice.GetCurrencyType("ok"),
			},
			expectedErrorDescription: "currency value ok is not valid",
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			client, apiHandler, _, teardown, err := setup()
			assert.NoError(t, err)
			defer teardown()

			if tc.handlerFunc != nil {
				apiHandler.HandleFunc(endpoint, tc.handlerFunc)
			} else {
				apiHandler.HandleFunc(endpoint, defaultResponse)
			}

			_, _, err = client.GetTokenPrices(tc.params)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				assert.Equal(t, tc.expectedErrorDescription, err.Error())
				return
			}
			assert.NoError(t, err)
		})
	}
}
