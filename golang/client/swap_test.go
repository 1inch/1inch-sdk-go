package client

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"1inch-sdk-golang/client/swap"
	"1inch-sdk-golang/helpers/consts/addresses"
	"1inch-sdk-golang/helpers/consts/amounts"
	"1inch-sdk-golang/helpers/consts/tokens"
)

func TestApproveAllowance(t *testing.T) {

	endpoint := "/swap/v5.2/1/approve/allowance"
	defaultResponse := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w,
			`{
    "allowance": "0"
}`,
		)
	}

	testcases := []struct {
		description              string
		handlerFunc              func(w http.ResponseWriter, r *http.Request)
		params                   swap.ApproveControllerGetAllowanceParams
		expectedOutput           string
		expectedErrorDescription string
	}{
		{
			description: "Success",
			params: swap.ApproveControllerGetAllowanceParams{
				TokenAddress:  tokens.EthereumWeth,
				WalletAddress: addresses.Vitalik,
			},
			expectedOutput: "0",
		},
		{
			description: "Error - exclude tokenAddress",
			params: swap.ApproveControllerGetAllowanceParams{
				WalletAddress: addresses.Vitalik,
			},
			expectedErrorDescription: "request validation error: tokenAddress is required",
		},
		{
			description: "Error - exclude walletAddress",
			params: swap.ApproveControllerGetAllowanceParams{
				TokenAddress: tokens.EthereumWeth,
			},
			expectedErrorDescription: "request validation error: walletAddress is required",
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

			allowanceResponse, _, err := c.Swap.ApproveAllowance(context.Background(), tc.params)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				assert.Equal(t, tc.expectedErrorDescription, err.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "0", allowanceResponse.Allowance)
		})
	}
}

func TestApproveTransaction(t *testing.T) {

	endpoint := "/swap/v5.2/1/approve/transaction"
	defaultResponse := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w,
			`{
    "data": "0x095ea7b30000000000000000000000001111111254eeb25477b68fb85ed929f73a960582ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
    "gasPrice": "12161498847",
    "to": "`+tokens.EthereumUsdc+`",
    "value": "0"
}`,
		)
	}

	testcases := []struct {
		description              string
		handlerFunc              func(w http.ResponseWriter, r *http.Request)
		params                   swap.ApproveControllerGetCallDataParams
		expectedOutput           swap.ApproveCallDataResponse
		expectedErrorDescription string
	}{
		{
			description: "Success",
			params: swap.ApproveControllerGetCallDataParams{
				TokenAddress: tokens.EthereumUsdc,
			},
			expectedOutput: swap.ApproveCallDataResponse{
				To: tokens.EthereumUsdc,
			},
		},
		{
			description:              "Error - exclude tokenAddress",
			params:                   swap.ApproveControllerGetCallDataParams{},
			expectedErrorDescription: "request validation error: tokenAddress is required",
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

			approveTransactionResponse, _, err := c.Swap.ApproveTransaction(context.Background(), tc.params)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				assert.Equal(t, tc.expectedErrorDescription, err.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedOutput.To, approveTransactionResponse.To)
		})
	}
}

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
			description: "Error - exclude src",
			params: swap.AggregationControllerGetSwapParams{
				Dst:    tokens.EthereumWeth,
				From:   addresses.Vitalik,
				Amount: amounts.Ten18,
			},
			expectedErrorDescription: "request validation error: src is required",
		},
		{
			description: "Error - exclude dst",
			params: swap.AggregationControllerGetSwapParams{
				Src:    tokens.EthereumUsdc,
				From:   addresses.Vitalik,
				Amount: amounts.Ten18,
			},
			expectedErrorDescription: "request validation error: dst is required",
		},
		{
			description: "Error - exclude amount",
			params: swap.AggregationControllerGetSwapParams{
				Src:  tokens.EthereumUsdc,
				Dst:  tokens.EthereumWeth,
				From: addresses.Vitalik,
			},
			expectedErrorDescription: "request validation error: amount is required",
		},
		{
			description: "Error - exclude from",
			params: swap.AggregationControllerGetSwapParams{
				Src:    tokens.EthereumUsdc,
				Dst:    tokens.EthereumWeth,
				Amount: amounts.Ten18,
			},
			expectedErrorDescription: "request validation error: from is required",
		},
		{
			description: "Error - src and dst identical",
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

			c, apiHandler, _, teardown, err := setup()
			require.NoError(t, err)
			defer teardown()

			if tc.handlerFunc != nil {
				apiHandler.HandleFunc(endpoint, tc.handlerFunc)
			} else {
				apiHandler.HandleFunc(endpoint, defaultResponse)
			}

			_, _, err = c.Swap.GetSwap(context.Background(), tc.params)
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
