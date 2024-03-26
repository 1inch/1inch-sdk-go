package client

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/client/models"
	"github.com/1inch/1inch-sdk-go/helpers/consts/addresses"
	"github.com/1inch/1inch-sdk-go/helpers/consts/amounts"
	"github.com/1inch/1inch-sdk-go/helpers/consts/chains"
	"github.com/1inch/1inch-sdk-go/helpers/consts/tokens"
)

func TestApproveAllowance(t *testing.T) {

	endpoint := "/swap/v6.0/1/approve/allowance"
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
		params                   models.ApproveAllowanceParams
		expectedOutput           string
		expectedErrorDescription string
	}{
		{
			description: "Success",
			params: models.ApproveAllowanceParams{
				ChainId: chains.Ethereum,
				ApproveControllerGetAllowanceParams: models.ApproveControllerGetAllowanceParams{
					TokenAddress:  tokens.EthereumWeth,
					WalletAddress: addresses.Vitalik,
				},
			},
			expectedOutput: "0",
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

			allowanceResponse, _, err := c.SwapApi.GetApproveAllowance(context.Background(), tc.params)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				require.Contains(t, err.Error(), tc.expectedErrorDescription)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "0", allowanceResponse.Allowance)
		})
	}
}

func TestApproveTransaction(t *testing.T) {

	endpoint := "/swap/v6.0/1/approve/transaction"
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
		params                   models.ApproveTransactionParams
		expectedOutput           models.ApproveCallDataResponse
		expectedErrorDescription string
	}{
		{
			description: "Success",
			params: models.ApproveTransactionParams{
				ChainId: chains.Ethereum,
				ApproveControllerGetCallDataParams: models.ApproveControllerGetCallDataParams{
					TokenAddress: tokens.EthereumUsdc,
				},
			},
			expectedOutput: models.ApproveCallDataResponse{
				To: tokens.EthereumUsdc,
			},
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

			approveTransactionResponse, _, err := c.SwapApi.GetApproveTransaction(context.Background(), tc.params)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				require.Contains(t, err.Error(), tc.expectedErrorDescription)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedOutput.To, approveTransactionResponse.To)
		})
	}
}

func TestGetQuote(t *testing.T) {

	endpoint := "/swap/v6.0/1/quote"
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
		params                   models.GetQuoteParams
		expectedOutput           string
		expectedErrorDescription string
	}{
		{
			description: "Success - Quote swap WETH to USDC",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tokens.EthereumUsdc, r.URL.Query().Get("src"))
				assert.Equal(t, tokens.EthereumWeth, r.URL.Query().Get("dst"))
				assert.Equal(t, amounts.Ten18, r.URL.Query().Get("amount"))
			},
			params: models.GetQuoteParams{
				ChainId: chains.Ethereum,
				AggregationControllerGetQuoteParams: models.AggregationControllerGetQuoteParams{
					Src:    tokens.EthereumUsdc,
					Dst:    tokens.EthereumWeth,
					Amount: amounts.Ten18,
				},
			},
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

			_, _, err = c.SwapApi.GetQuote(context.Background(), tc.params)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				require.Contains(t, err.Error(), tc.expectedErrorDescription)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestGetSwap(t *testing.T) {

	endpoint := "/swap/v6.0/1/swap"
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
		params                   models.GetSwapParams
		expectedOutput           string
		expectedErrorDescription string
	}{
		{
			description: "Success - Swap WETH to USDC",
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tokens.EthereumUsdc, r.URL.Query().Get("src"))
				assert.Equal(t, tokens.EthereumWeth, r.URL.Query().Get("dst"))
				assert.Equal(t, addresses.Vitalik, r.URL.Query().Get("from"))
				assert.Equal(t, amounts.Ten18, r.URL.Query().Get("amount"))
			},
			params: models.GetSwapParams{
				AggregationControllerGetSwapParams: models.AggregationControllerGetSwapParams{
					Src:      tokens.EthereumUsdc,
					Dst:      tokens.EthereumWeth,
					From:     addresses.Vitalik,
					Amount:   amounts.Ten18,
					Slippage: 0.5,
				},
				ChainId:      chains.Ethereum,
				SkipWarnings: true,
			},
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

			_, _, err = c.SwapApi.GetSwap(context.Background(), tc.params)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				require.Contains(t, err.Error(), tc.expectedErrorDescription)
				return
			}
			require.NoError(t, err)
		})
	}
}
