package client

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/svanas/1inch-sdk/golang/client/swap"
	"github.com/svanas/1inch-sdk/golang/helpers/consts/addresses"
	"github.com/svanas/1inch-sdk/golang/helpers/consts/amounts"
	"github.com/svanas/1inch-sdk/golang/helpers/consts/chains"
	"github.com/svanas/1inch-sdk/golang/helpers/consts/tokens"
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
		params                   swap.ApproveAllowanceParams
		expectedOutput           string
		expectedErrorDescription string
	}{
		{
			description: "Success",
			params: swap.ApproveAllowanceParams{
				ChainId: chains.Ethereum,
				ApproveControllerGetAllowanceParams: swap.ApproveControllerGetAllowanceParams{
					TokenAddress:  tokens.EthereumWeth,
					WalletAddress: addresses.Vitalik,
				},
			},
			expectedOutput: "0",
		},
		{
			description: "Error - Exclude tokenAddress",
			params: swap.ApproveAllowanceParams{
				ChainId: chains.Ethereum,
				ApproveControllerGetAllowanceParams: swap.ApproveControllerGetAllowanceParams{
					WalletAddress: addresses.Vitalik,
				},
			},
			expectedErrorDescription: `'tokenAddress' is required`,
		},
		{
			description: "Error - Exclude walletAddress",
			params: swap.ApproveAllowanceParams{
				ChainId: chains.Ethereum,
				ApproveControllerGetAllowanceParams: swap.ApproveControllerGetAllowanceParams{
					TokenAddress: tokens.EthereumWeth,
				},
			},
			expectedErrorDescription: `'walletAddress' is required`,
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

			allowanceResponse, _, err := c.Swap.ApproveAllowance(context.Background(), tc.params)
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
		params                   swap.ApproveTransactionParams
		expectedOutput           swap.ApproveCallDataResponse
		expectedErrorDescription string
	}{
		{
			description: "Success",
			params: swap.ApproveTransactionParams{
				ChainId: chains.Ethereum,
				ApproveControllerGetCallDataParams: swap.ApproveControllerGetCallDataParams{
					TokenAddress: tokens.EthereumUsdc,
				},
			},
			expectedOutput: swap.ApproveCallDataResponse{
				To: tokens.EthereumUsdc,
			},
		},
		{
			description: "Error - exclude tokenAddress",
			params: swap.ApproveTransactionParams{
				ChainId: chains.Ethereum,
			},
			expectedErrorDescription: `'tokenAddress' is required`,
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

			approveTransactionResponse, _, err := c.Swap.ApproveTransaction(context.Background(), tc.params)
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
		params                   swap.GetQuoteParams
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
			params: swap.GetQuoteParams{
				ChainId: chains.Ethereum,
				AggregationControllerGetQuoteParams: swap.AggregationControllerGetQuoteParams{
					Src:    tokens.EthereumUsdc,
					Dst:    tokens.EthereumWeth,
					Amount: amounts.Ten18,
				},
			},
		},
		{
			description: "Error - exclude src",
			params: swap.GetQuoteParams{
				ChainId: chains.Ethereum,
				AggregationControllerGetQuoteParams: swap.AggregationControllerGetQuoteParams{
					Dst:    tokens.EthereumWeth,
					Amount: amounts.Ten18,
				},
			},
			expectedErrorDescription: `'src' is required`,
		},
		{
			description: "Error - exclude Dst",
			params: swap.GetQuoteParams{
				ChainId: chains.Ethereum,
				AggregationControllerGetQuoteParams: swap.AggregationControllerGetQuoteParams{
					Src:    tokens.EthereumUsdc,
					Amount: amounts.Ten18,
				},
			},
			expectedErrorDescription: `'dst' is required`,
		},
		{
			description: "Error - exclude amount",
			params: swap.GetQuoteParams{
				ChainId: chains.Ethereum,
				AggregationControllerGetQuoteParams: swap.AggregationControllerGetQuoteParams{
					Src: tokens.EthereumUsdc,
					Dst: tokens.EthereumWeth,
				},
			},
			expectedErrorDescription: `'amount' is required`,
		},
		{
			description: "Error - src and dst identical",
			params: swap.GetQuoteParams{
				ChainId: chains.Ethereum,
				AggregationControllerGetQuoteParams: swap.AggregationControllerGetQuoteParams{
					Src:    tokens.EthereumUsdc,
					Dst:    tokens.EthereumUsdc,
					Amount: amounts.Ten18,
				},
			},
			expectedErrorDescription: "src and dst tokens must be different",
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

			_, _, err = c.Swap.GetQuote(context.Background(), tc.params)
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
		params                   swap.GetSwapDataParams
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
			params: swap.GetSwapDataParams{
				AggregationControllerGetSwapParams: swap.AggregationControllerGetSwapParams{
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
		{
			description: "Error - exclude src",
			params: swap.GetSwapDataParams{
				AggregationControllerGetSwapParams: swap.AggregationControllerGetSwapParams{
					Dst:      tokens.EthereumWeth,
					From:     addresses.Vitalik,
					Amount:   amounts.Ten18,
					Slippage: 0.5,
				},
				ChainId:      chains.Ethereum,
				SkipWarnings: true,
			},
			expectedErrorDescription: `'src' is required`,
		},
		{
			description: "Error - exclude dst",
			params: swap.GetSwapDataParams{
				AggregationControllerGetSwapParams: swap.AggregationControllerGetSwapParams{
					Src:      tokens.EthereumUsdc,
					From:     addresses.Vitalik,
					Amount:   amounts.Ten18,
					Slippage: 0.5,
				},
				ChainId:      chains.Ethereum,
				SkipWarnings: true,
			},
			expectedErrorDescription: `'dst' is required`,
		},
		{
			description: "Error - exclude amount",
			params: swap.GetSwapDataParams{
				AggregationControllerGetSwapParams: swap.AggregationControllerGetSwapParams{
					Src:      tokens.EthereumUsdc,
					Dst:      tokens.EthereumWeth,
					From:     addresses.Vitalik,
					Slippage: 0.5,
				},
				ChainId:      chains.Ethereum,
				SkipWarnings: true,
			},
			expectedErrorDescription: `'amount' is required`,
		},
		{
			description: "Error - exclude from",
			params: swap.GetSwapDataParams{
				AggregationControllerGetSwapParams: swap.AggregationControllerGetSwapParams{
					Src:      tokens.EthereumUsdc,
					Dst:      tokens.EthereumWeth,
					Amount:   amounts.Ten18,
					Slippage: 0.5,
				},
				ChainId:      chains.Ethereum,
				SkipWarnings: true,
			},
			expectedErrorDescription: `'from' is required`,
		},
		{
			description: "Error - src and dst identical",
			params: swap.GetSwapDataParams{
				AggregationControllerGetSwapParams: swap.AggregationControllerGetSwapParams{
					Src:      tokens.EthereumUsdc,
					Dst:      tokens.EthereumUsdc,
					Amount:   amounts.Ten18,
					From:     addresses.Vitalik,
					Slippage: 0.5,
				},
				ChainId:      chains.Ethereum,
				SkipWarnings: true,
			},
			expectedErrorDescription: "src and dst tokens must be different",
		},
		{
			description: "Error - slippage is required",
			params: swap.GetSwapDataParams{
				AggregationControllerGetSwapParams: swap.AggregationControllerGetSwapParams{
					Src:    tokens.EthereumUsdc,
					Dst:    tokens.EthereumWeth,
					Amount: amounts.Ten18,
					From:   addresses.Vitalik,
				},
				ChainId:      chains.Ethereum,
				SkipWarnings: true,
			},
			expectedErrorDescription: `'slippage' is required`,
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

			_, _, err = c.Swap.GetSwapData(context.Background(), tc.params)
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
