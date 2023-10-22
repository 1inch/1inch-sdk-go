package client

import (
	"context"
	"fmt"
	"os"
	"testing"

	"dev-portal-sdk-go/client/swap"
	"dev-portal-sdk-go/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
		ApiKey:            os.Getenv("DEV_PORTAL_TOKEN"),
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
