package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/1inch/1inch-sdk-go/client/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/helpers"
	"github.com/1inch/1inch-sdk-go/helpers/consts/addresses"
	"github.com/1inch/1inch-sdk-go/helpers/consts/chains"
	"github.com/1inch/1inch-sdk-go/helpers/consts/tokens"
)

func TestApproveAllowanceIntegration(t *testing.T) {

	testcases := []struct {
		description    string
		params         models.ApproveAllowanceParams
		expectedOutput models.AllowanceResponse
	}{
		{
			description: "Get approve spender address",
			params: models.ApproveAllowanceParams{
				ChainId: chains.Ethereum,
				ApproveControllerGetAllowanceParams: models.ApproveControllerGetAllowanceParams{
					TokenAddress:  tokens.EthereumUsdc,
					WalletAddress: addresses.Vitalik,
				},
			},
			expectedOutput: models.AllowanceResponse{
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

			allowance, resp, err := c.SwapApi.GetApproveAllowance(context.Background(), tc.params)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
			assert.Equal(t, tc.expectedOutput.Allowance, allowance.Allowance)
		})
	}
}

func TestApproveSpenderIntegration(t *testing.T) {

	testcases := []struct {
		description    string
		params         models.ApproveSpenderParams
		expectedOutput models.SpenderResponse
	}{
		{
			description: "Get approve spender address",
			params: models.ApproveSpenderParams{
				ChainId: chains.Ethereum,
			},
			expectedOutput: models.SpenderResponse{
				Address: "0x1111111254eeb25477b68fb85ed929f73a960582",
			},
		},
	}

	c, err := NewClient(SimpleEthereumConfig)
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			spender, resp, err := c.SwapApi.GetApproveSpender(context.Background(), tc.params)
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
		params         models.ApproveTransactionParams
		expectedOutput models.ApproveCallDataResponse
	}{
		{
			description: "Get approve spender address",
			params: models.ApproveTransactionParams{
				ChainId: chains.Ethereum,
				ApproveControllerGetCallDataParams: models.ApproveControllerGetCallDataParams{
					TokenAddress: tokens.EthereumUsdc,
					Amount:       "",
				},
			},
			expectedOutput: models.ApproveCallDataResponse{
				To: tokens.EthereumUsdc,
			},
		},
	}

	c, err := NewClient(SimpleEthereumConfig)
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			transaction, resp, err := c.SwapApi.GetApproveTransaction(context.Background(), tc.params)
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
		params           models.GetLiquiditySourcesParams
		expectedProtocol models.ProtocolImage
	}{
		{
			description: "Get approve spender address",
			params: models.GetLiquiditySourcesParams{
				ChainId: chains.Ethereum,
			},
			expectedProtocol: models.ProtocolImage{
				Id:  "UNISWAP_V2",
				Img: "Uniswap V2",
			},
		},
	}

	c, err := NewClient(SimpleEthereumConfig)
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			liquiditySources, resp, err := c.SwapApi.GetLiquiditySources(context.Background(), tc.params)
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
		params        models.GetTokensParams
		expectedToken models.TokenInfo
	}{
		{
			description: "Get approve spender address",
			params: models.GetTokensParams{
				ChainId: chains.Ethereum,
			},
			expectedToken: models.TokenInfo{
				Address:  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				Decimals: 6,
				Eip2612:  true,
				Name:     "USD Coin",
				Symbol:   "USDC",
			},
		},
	}

	c, err := NewClient(SimpleEthereumConfig)
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			tokens, resp, err := c.SwapApi.GetTokens(context.Background(), tc.params)
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
