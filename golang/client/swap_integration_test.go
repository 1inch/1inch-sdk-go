package client

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"1inch-sdk-golang/client/swap"
	"1inch-sdk-golang/helpers"
	"1inch-sdk-golang/helpers/consts/addresses"
	"1inch-sdk-golang/helpers/consts/tokens"
)

func TestApproveAllowanceIntegration(t *testing.T) {

	testcases := []struct {
		description    string
		params         swap.ApproveControllerGetAllowanceParams
		expectedOutput swap.AllowanceResponse
	}{
		{
			description: "Get approve spender address",
			params: swap.ApproveControllerGetAllowanceParams{
				TokenAddress:  tokens.EthereumUsdc,
				WalletAddress: addresses.Vitalik,
			},
			expectedOutput: swap.AllowanceResponse{
				Allowance: "0",
			},
		},
	}

	c, err := NewClient(Config{
		TargetEnvironment: EnvironmentProduction,
		DevPortalApiKey:   os.Getenv("DEV_PORTAL_TOKEN"),
	})
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			allowance, resp, err := c.Swap.ApproveAllowance(context.Background(), tc.params)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
			assert.Equal(t, tc.expectedOutput.Allowance, allowance.Allowance)

			helpers.Sleep()
		})
	}
}

func TestApproveSpenderIntegration(t *testing.T) {

	testcases := []struct {
		description    string
		expectedOutput swap.SpenderResponse
	}{
		{
			description: "Get approve spender address",
			expectedOutput: swap.SpenderResponse{
				Address: "0x1111111254eeb25477b68fb85ed929f73a960582",
			},
		},
	}

	c, err := NewClient(Config{
		TargetEnvironment: EnvironmentProduction,
		DevPortalApiKey:   os.Getenv("DEV_PORTAL_TOKEN"),
	})
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			spender, resp, err := c.Swap.ApproveSpender(context.Background())
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
		params         swap.ApproveControllerGetCallDataParams
		expectedOutput swap.ApproveCallDataResponse
	}{
		{
			description: "Get approve spender address",
			params: swap.ApproveControllerGetCallDataParams{
				TokenAddress: tokens.EthereumUsdc,
				Amount:       nil,
			},
			expectedOutput: swap.ApproveCallDataResponse{
				To: tokens.EthereumUsdc,
			},
		},
	}

	c, err := NewClient(Config{
		TargetEnvironment: EnvironmentProduction,
		DevPortalApiKey:   os.Getenv("DEV_PORTAL_TOKEN"),
	})
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
		expectedProtocol swap.ProtocolImage
	}{
		{
			description: "Get approve spender address",
			expectedProtocol: swap.ProtocolImage{
				Id:  "UNISWAP_V2",
				Img: "Uniswap V2",
			},
		},
	}

	c, err := NewClient(Config{
		TargetEnvironment: EnvironmentProduction,
		DevPortalApiKey:   os.Getenv("DEV_PORTAL_TOKEN"),
	})
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			liquiditySources, resp, err := c.Swap.GetLiquiditySources(context.Background())
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
		expectedToken swap.TokenInfo
	}{
		{
			description: "Get approve spender address",
			expectedToken: swap.TokenInfo{
				Address:  "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				Decimals: 6,
				Eip2612:  helpers.BoolPtr(true),
				Name:     "USD Coin",
				Symbol:   "USDC",
			},
		},
	}

	c, err := NewClient(Config{
		TargetEnvironment: EnvironmentProduction,
		DevPortalApiKey:   os.Getenv("DEV_PORTAL_TOKEN"),
	})
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			tokens, resp, err := c.Swap.GetTokens(context.Background())
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
