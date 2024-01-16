//go:build e2e
// +build e2e

package actions

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/1inch/1inch-sdk/golang/client"
	"github.com/1inch/1inch-sdk/golang/client/swap"
	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/tokens"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/web3providers"
	"github.com/stretchr/testify/require"
)

func swapTokens(config client.Config, swapParams swap.AggregationControllerGetSwapParams) error {

	// Create the 1inch client
	c, err := client.NewClient(config)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}

	// Swap tokens
	err = SwapTokens(c, swapParams, true)
	if err != nil {
		log.Fatalf("Failed to swap tokens: %v", err)
	}

	return nil
}

func TestSwapTokensE2E(t *testing.T) {

	testcases := []struct {
		description    string
		config         client.Config
		swapParams     swap.AggregationControllerGetSwapParams
		expectedOutput string
	}{
		{
			description: "Polygon - Swap 0.01 FRAX for USDC - Permit1",
			config: client.Config{
				DevPortalApiKey:  os.Getenv("DEV_PORTAL_TOKEN"),
				WalletKey:        os.Getenv("WALLET_KEY"),
				Web3HttpProvider: os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY_POLYGON"),
				ChainId:          chains.Polygon,
			},
			swapParams: swap.AggregationControllerGetSwapParams{
				Src:      tokens.PolygonFrax,
				Dst:      tokens.PolygonUsdc,
				Amount:   "10000000000000000",
				From:     os.Getenv("WALLET_ADDRESS"),
				Slippage: 0.5,
			},
		},
		{
			description: "Polygon - Swap 0.01 DAI for USDC - Approval",
			config: client.Config{
				DevPortalApiKey:  os.Getenv("DEV_PORTAL_TOKEN"),
				WalletKey:        os.Getenv("WALLET_KEY"),
				Web3HttpProvider: os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY_POLYGON"),
				ChainId:          chains.Polygon,
			},
			swapParams: swap.AggregationControllerGetSwapParams{
				Src:      tokens.PolygonDai,
				Dst:      tokens.PolygonUsdc,
				Amount:   "10000000000000000",
				From:     os.Getenv("WALLET_ADDRESS"),
				Slippage: 0.5,
			},
		},
		{
			description: "Arbitrum - Swap 0.01 USDC for DAI - Permit1",
			config: client.Config{
				DevPortalApiKey:  os.Getenv("DEV_PORTAL_TOKEN"),
				WalletKey:        os.Getenv("WALLET_KEY"),
				Web3HttpProvider: web3providers.Arbitrum,
				ChainId:          chains.Arbitrum,
			},
			swapParams: swap.AggregationControllerGetSwapParams{
				Src:      tokens.ArbitrumUsdc,
				Dst:      tokens.ArbitrumDai,
				Amount:   "10000",
				From:     os.Getenv("WALLET_ADDRESS"),
				Slippage: 0.5,
			},
		},
		{
			description: "Arbitrum - Swap $0.01 worth of ETH for USDC - Permit1",
			config: client.Config{
				DevPortalApiKey:  os.Getenv("DEV_PORTAL_TOKEN"),
				WalletKey:        os.Getenv("WALLET_KEY"),
				Web3HttpProvider: web3providers.Arbitrum,
				ChainId:          chains.Arbitrum,
			},
			swapParams: swap.AggregationControllerGetSwapParams{
				Src:      tokens.NativeToken,
				Dst:      tokens.ArbitrumUsdc,
				Amount:   "3974301376798",
				From:     os.Getenv("WALLET_ADDRESS"),
				Slippage: 0.5,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			t.Cleanup(func() {
				helpers.Sleep()
			})

			err := swapTokens(tc.config, tc.swapParams)
			require.NoError(t, err)
		})
	}
}
