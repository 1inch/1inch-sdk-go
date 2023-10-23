package client

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"1inch-sdk-golang/client/swap"
	"1inch-sdk-golang/helpers/consts/amounts"
	"1inch-sdk-golang/helpers/consts/tokens"
)

func TestGetQuote(t *testing.T) {

	endpoint := "/swap/v5.2/1/quote"
	defaultResponse := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w,
			`{
				"toAmount": "639290156773079"
			}`,
		)
	}

	testcases := []struct {
		description              string
		handlerFunc              func(w http.ResponseWriter, r *http.Request)
		params                   swap.AggregationControllerGetQuoteParams
		expectedOutput           string
		expectedErrorDescription string
	}{
		{
			description: "Quote swap WETH to USDC",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tokens.EthereumUsdc, r.URL.Query().Get("src"))
				assert.Equal(t, tokens.EthereumWeth, r.URL.Query().Get("dst"))
				assert.Equal(t, amounts.Ten18, r.URL.Query().Get("amount"))
			},
			params: swap.AggregationControllerGetQuoteParams{
				Src:    tokens.EthereumUsdc,
				Dst:    tokens.EthereumWeth,
				Amount: amounts.Ten18,
			},
		},
		{
			description: "Error - exclude src",
			params: swap.AggregationControllerGetQuoteParams{
				Dst:    tokens.EthereumWeth,
				Amount: amounts.Ten18,
			},
			expectedErrorDescription: "request validation error: src is required",
		},
		{
			description: "Error - exclude dst",
			params: swap.AggregationControllerGetQuoteParams{
				Src:    tokens.EthereumUsdc,
				Amount: amounts.Ten18,
			},
			expectedErrorDescription: "request validation error: dst is required",
		},
		{
			description: "Error - exclude amount",
			params: swap.AggregationControllerGetQuoteParams{
				Src: tokens.EthereumUsdc,
				Dst: tokens.EthereumWeth,
			},
			expectedErrorDescription: "request validation error: amount is required",
		},
		{
			description: "Error - src and dst identical",
			params: swap.AggregationControllerGetQuoteParams{
				Src:    tokens.EthereumUsdc,
				Dst:    tokens.EthereumUsdc,
				Amount: amounts.Ten18,
			},
			expectedErrorDescription: "request validation error: src and dst tokens must be different",
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

			_, _, err = c.Swap.GetQuote(context.Background(), tc.params)
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
