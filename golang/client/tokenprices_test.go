package client

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"1inch-sdk-golang/client/tokenprices"
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
		params                   tokenprices.ChainControllerByAddressesParams
		expectedOutput           string
		expectedErrorDescription string
	}{
		{
			description: "Get prices in USD",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, string(tokenprices.USD), r.URL.Query().Get("currency"))
			},
			params: tokenprices.ChainControllerByAddressesParams{
				Currency: tokenprices.GetCurrencyParameter(tokenprices.USD),
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
			params: tokenprices.ChainControllerByAddressesParams{
				Currency: tokenprices.GetCurrencyParameter("ok"),
			},
			expectedErrorDescription: "currency value ok is not valid",
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			c, apiHandler, _, teardown, err := setup()
			require.NoError(t, err)
			defer teardown()

			if tc.handlerFunc != nil {
				apiHandler.HandleFunc(endpoint, tc.handlerFunc)
			} else {
				apiHandler.HandleFunc(endpoint, defaultResponse)
			}

			_, _, err = c.TokenPrices.GetPrices(context.Background(), tc.params)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				assert.Equal(t, tc.expectedErrorDescription, err.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}
