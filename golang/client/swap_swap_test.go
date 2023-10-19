package client

import (
	"fmt"
	"net/http"
	"testing"

	"dev-portal-sdk-go/client/swap"
	"dev-portal-sdk-go/helpers/consts/addresses"
	"dev-portal-sdk-go/helpers/consts/amounts"
	"dev-portal-sdk-go/helpers/consts/tokens"
	"github.com/stretchr/testify/assert"
)

func TestGetSwap(t *testing.T) {

	endpoint := "/swap/v5.2/1/swap"
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
		params                   swap.AggregationControllerGetSwapParams
		expectedOutput           string
		expectedErrorDescription string
	}{
		{
			description: "Swap WETH to USDC",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tokens.EthereumUsdc, r.URL.Query().Get("src"))
				assert.Equal(t, tokens.EthereumWeth, r.URL.Query().Get("dst"))
				assert.Equal(t, addresses.Vitalik, r.URL.Query().Get("from"))
				assert.Equal(t, amounts.Ten18, r.URL.Query().Get("amount"))
			},
			params: swap.AggregationControllerGetSwapParams{
				Src:    tokens.EthereumUsdc,
				Dst:    tokens.EthereumWeth,
				From:   addresses.Vitalik,
				Amount: amounts.Ten18,
			},
		},
		{
			description: "Swap Error - exclude src",
			params: swap.AggregationControllerGetSwapParams{
				Dst:    tokens.EthereumWeth,
				From:   addresses.Vitalik,
				Amount: amounts.Ten18,
			},
			expectedErrorDescription: "request validation error: src is required",
		},
		{
			description: "Swap Error - exclude dst",
			params: swap.AggregationControllerGetSwapParams{
				Src:    tokens.EthereumUsdc,
				From:   addresses.Vitalik,
				Amount: amounts.Ten18,
			},
			expectedErrorDescription: "request validation error: dst is required",
		},
		{
			description: "Swap Error - exclude amount",
			params: swap.AggregationControllerGetSwapParams{
				Src:  tokens.EthereumUsdc,
				Dst:  tokens.EthereumWeth,
				From: addresses.Vitalik,
			},
			expectedErrorDescription: "request validation error: amount is required",
		},
		{
			description: "Swap Error - exclude from",
			params: swap.AggregationControllerGetSwapParams{
				Src:    tokens.EthereumUsdc,
				Dst:    tokens.EthereumWeth,
				Amount: amounts.Ten18,
			},
			expectedErrorDescription: "request validation error: from is required",
		},
		{
			description: "Swap Error - src and dst identical",
			params: swap.AggregationControllerGetSwapParams{
				Src:    tokens.EthereumUsdc,
				Dst:    tokens.EthereumUsdc,
				Amount: amounts.Ten18,
				From:   addresses.Vitalik,
			},
			expectedErrorDescription: "request validation error: src and dst tokens must be different",
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

			_, _, err = client.GetSwap(tc.params)
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
