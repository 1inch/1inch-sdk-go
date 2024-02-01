//go:build ignoretest
// +build ignoretest

package client

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk/golang/client/orderbook"
	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/addresses"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/tokens"
)

func TestCreateOrderE2E(t *testing.T) {

	testcases := []struct {
		description       string
		config            Config
		createOrderParams orderbook.CreateOrderParams
		expectedOutput    string
	}{
		{
			description: "Polygon - Create limit order offering 1 DAI for 1 USDC",
			config: Config{
				DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
				Web3HttpProviders: []Web3ProviderConfig{
					{
						ChainId: chains.Polygon,
						Url:     os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY_POLYGON"),
					},
				},
			},
			createOrderParams: orderbook.CreateOrderParams{
				ChainId:      chains.Polygon,
				WalletKey:    os.Getenv("WALLET_KEY"),
				SourceWallet: os.Getenv("WALLET_ADDRESS"),
				FromToken:    tokens.PolygonDai,
				ToToken:      tokens.PolygonUsdc,
				MakingAmount: "1000000000000000000",
				TakingAmount: "1000000",
				Receiver:     addresses.Zero,
				SkipWarnings: true,
			},
		},
		{
			description: "Ethereum - Create limit order offering 1 DAI for 1 USDC",
			config: Config{
				DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
				Web3HttpProviders: []Web3ProviderConfig{
					{
						ChainId: chains.Ethereum,
						Url:     os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY"),
					},
				},
			},
			createOrderParams: orderbook.CreateOrderParams{
				ChainId:      chains.Ethereum,
				WalletKey:    os.Getenv("WALLET_KEY"),
				SourceWallet: os.Getenv("WALLET_ADDRESS"),
				FromToken:    tokens.EthereumUsdc,
				ToToken:      tokens.EthereumDai,
				MakingAmount: "1000000",
				TakingAmount: "1000000000000000000",
				Receiver:     addresses.Zero,
				SkipWarnings: true,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			t.Cleanup(func() {
				helpers.Sleep()
			})

			c, err := NewClient(tc.config)
			require.NoError(t, err)

			_, _, err = c.Orderbook.CreateOrder(context.Background(), tc.createOrderParams)
			require.NoError(t, err)
		})
	}
}
