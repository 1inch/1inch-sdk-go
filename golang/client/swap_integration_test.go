package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk/golang/client/swap"
	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/addresses"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/tokens"
)

func TestApproveAllowanceIntegration(t *testing.T) {

	testcases := []struct {
		description    string
		params         swap.ApproveAllowanceParams
		expectedOutput swap.AllowanceResponse
	}{
		{
			description: "Get approve spender address",
			params: swap.ApproveAllowanceParams{
				RequestParams: swap.RequestParams{
					ChainId: chains.Ethereum,
				},
				ApproveControllerGetAllowanceParams: swap.ApproveControllerGetAllowanceParams{
					TokenAddress:  tokens.EthereumUsdc,
					WalletAddress: addresses.Vitalik,
				},
			},
			expectedOutput: swap.AllowanceResponse{
				Allowance: "0",
			},
		},
	}

	c, err := NewClient(SimpleEthereumConfig)
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			t.Cleanup(func() {
				helpers.Sleep()
			})

			allowance, resp, err := c.Swap.ApproveAllowance(context.Background(), tc.params)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
			assert.Equal(t, tc.expectedOutput.Allowance, allowance.Allowance)
		})
	}
}

func TestApproveSpenderIntegration(t *testing.T) {

	testcases := []struct {
		description    string
		params         swap.ApproveSpenderParams
		expectedOutput swap.SpenderResponse
	}{
		{
			description: "Get approve spender address",
			params: swap.ApproveSpenderParams{
				RequestParams: swap.RequestParams{
					ChainId: chains.Ethereum,
				},
			},
			expectedOutput: swap.SpenderResponse{
				Address: "0x1111111254eeb25477b68fb85ed929f73a960582",
			},
		},
	}

	c, err := NewClient(SimpleEthereumConfig)
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			spender, resp, err := c.Swap.ApproveSpender(context.Background(), tc.params)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
			assert.Equal(t, tc.expectedOutput.Address, spender.Address)

			helpers.Sleep()
		})
	}
}

func TestApproveTransactionIntegration(t *testing.T) {

	testcases := []struct {
		description    string
		params         swap.ApproveTransactionParams
		expectedOutput swap.ApproveCallDataResponse
	}{
		{
			description: "Get approve spender address",
			params: swap.ApproveTransactionParams{
				RequestParams: swap.RequestParams{
					ChainId: chains.Ethereum,
				},
				ApproveControllerGetCallDataParams: swap.ApproveControllerGetCallDataParams{
					TokenAddress: tokens.EthereumUsdc,
					Amount:       nil,
				},
			},
			expectedOutput: swap.ApproveCallDataResponse{
				To: tokens.EthereumUsdc,
			},
		},
	}

	c, err := NewClient(SimpleEthereumConfig)
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			transaction, resp, err := c.Swap.ApproveTransaction(context.Background(), tc.params)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
			assert.Equal(t, tc.expectedOutput.To, transaction.To)

			helpers.Sleep()
		})
	}
}

func TestGetLiquiditySourcesIntegration(t *testing.T) {

	testcases := []struct {
		description      string
		params           swap.GetLiquiditySourcesParams
		expectedProtocol swap.ProtocolImage
	}{
		{
			description: "Get approve spender address",
			params: swap.GetLiquiditySourcesParams{
				RequestParams: swap.RequestParams{
					ChainId: chains.Ethereum,
				},
			},
			expectedProtocol: swap.ProtocolImage{
				Id:  "UNISWAP_V2",
				Img: "Uniswap V2",
			},
		},
	}

	c, err := NewClient(SimpleEthereumConfig)
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			liquiditySources, resp, err := c.Swap.GetLiquiditySources(context.Background(), tc.params)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)

			found := false
			for _, liquiditySource := range liquiditySources.Protocols {
				if liquiditySource.Id == tc.expectedProtocol.Id {
					found = true
					break
				}
			}
			assert.True(t, found, fmt.Sprintf("expected to find %s in liquidity sources list, but did not", tc.expectedProtocol.Id))
			helpers.Sleep()
		})
	}
}

func TestGetTokensIntegration(t *testing.T) {

	testcases := []struct {
		description   string
		params        swap.GetTokensParams
		expectedToken swap.TokenInfo
	}{
		{
			description: "Get approve spender address",
			params: swap.GetTokensParams{
				RequestParams: swap.RequestParams{
					ChainId: chains.Ethereum,
				},
			},
			expectedToken: swap.TokenInfo{
				Address:  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				Decimals: 6,
				Eip2612:  helpers.GetPtr(true),
				Name:     "USD Coin",
				Symbol:   "USDC",
			},
		},
	}

	c, err := NewClient(SimpleEthereumConfig)
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			tokens, resp, err := c.Swap.GetTokens(context.Background(), tc.params)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)

			found := false
			for _, tok := range tokens.Tokens {
				if tok.Address == tc.expectedToken.Address {
					found = true
					break
				}
			}
			assert.True(t, found, fmt.Sprintf("expected to find %s in tokens list, but did not", tc.expectedToken.Address))
			helpers.Sleep()
		})
	}
}
