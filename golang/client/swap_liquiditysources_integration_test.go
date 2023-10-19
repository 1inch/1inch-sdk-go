package client

import (
	"fmt"
	"os"
	"testing"

	"dev-portal-sdk-go/client/swap"
	"dev-portal-sdk-go/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGetLiquiditySourcesIntegration(t *testing.T) {

	testcases := []struct {
		description      string
		expectedProtocol swap.ProtocolImage
	}{
		{
			description: "Get approve spender address",
			expectedProtocol: swap.ProtocolImage{
				Id:       "UNISWAP_V2",
				Img:      "Uniswap V2",
				ImgColor: "https://cdn.1inch.io/liquidity-sources-logo/uniswap.png",
				Title:    "https://cdn.1inch.io/liquidity-sources-logo/uniswap_color.png",
			},
		},
	}

	c, err := NewClient(Config{
		TargetEnvironment: EnvironmentProduction,
		ApiKey:            os.Getenv("DEV_PORTAL_TOKEN"),
	})
	assert.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			liquiditySources, resp, err := c.GetLiquiditySources()
			assert.NoError(t, err)
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
