package client

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk/golang/client/fusion"
	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/addresses"
)

func TestGetOrders(t *testing.T) {

	endpoint := "/fusion/orders/v1.0/1/order/active"
	defaultOrderHash := "0x7b7890f7de294ea020fe1933641fa9f10b8a7975d7670c3c5046e590f973ef8a"
	defaultResponse := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
    "items": [
        {
            "orderHash": "`+defaultOrderHash+`"
        }
	]
}`)
	}

	testcases := []struct {
		description              string
		handlerFunc              func(w http.ResponseWriter, r *http.Request)
		owner                    string
		params                   fusion.OrderApiControllerGetActiveOrdersParams
		expectedOutput           string
		expectedErrorDescription string
	}{
		{
			description:    "Success",
			owner:          addresses.Vitalik,
			params:         fusion.OrderApiControllerGetActiveOrdersParams{},
			expectedOutput: defaultOrderHash,
		},
		{
			description: "Error - Invalid page value",
			owner:       addresses.Vitalik,
			params: fusion.OrderApiControllerGetActiveOrdersParams{
				Page: helpers.GetPtr(float32(0)),
			},
			expectedErrorDescription: "page must be greater than 0",
		},
		{
			description: "Error - Invalid limit value",
			owner:       addresses.Vitalik,
			params: fusion.OrderApiControllerGetActiveOrdersParams{
				Limit: helpers.GetPtr(float32(0)),
			},
			expectedErrorDescription: "limit must be greater than 0",
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			c, mux, _, teardown, err := setup()
			require.NoError(t, err)
			defer teardown()

			if tc.handlerFunc != nil {
				mux.HandleFunc(endpoint, tc.handlerFunc)
			} else {
				mux.HandleFunc(endpoint, defaultResponse)
			}

			allOrders, _, err := c.Fusion.GetOrders(context.Background(), tc.params)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				require.Contains(t, err.Error(), tc.expectedErrorDescription)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, allOrders.Items[0].OrderHash)
		})
	}
}
