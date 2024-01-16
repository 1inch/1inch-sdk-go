//go:build e2e
// +build e2e

package client

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/1inch/1inch-sdk/golang/client/orderbook"
	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/addresses"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/tokens"
	"github.com/stretchr/testify/require"
)

func createOrder(config Config, createOrderParams orderbook.OrderRequest) (*orderbook.CreateOrderResponse, error) {

	// Create the 1inch client
	c, err := NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	// Execute orders request
	createOrderResponse, _, err := c.Orderbook.CreateOrder(context.Background(), createOrderParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %v", err)
	}

	return createOrderResponse, err
}

func TestCreateOrderE2E(t *testing.T) {

	testcases := []struct {
		description       string
		config            Config
		createOrderParams orderbook.OrderRequest
		expectedOutput    string
	}{
		{
			description: "Polygon - Create limit order offering 1 DAI for 1 USDC",
			config: Config{
				DevPortalApiKey:  os.Getenv("DEV_PORTAL_TOKEN"),
				Web3HttpProvider: os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY_POLYGON"),
				WalletKey:        os.Getenv("WALLET_KEY"),
				ChainId:          chains.Polygon,
			},
			createOrderParams: orderbook.OrderRequest{
				SourceWallet: os.Getenv("WALLET_ADDRESS"),
				FromToken:    tokens.PolygonDai,
				ToToken:      tokens.PolygonUsdc,
				MakingAmount: 1000000000000000000,
				TakingAmount: 1000000,
				Receiver:     addresses.Zero,
				SkipWarnings: true,
			},
		},
		{
			description: "Ethereum - Create limit order offering 1 DAI for 1 USDC",
			config: Config{
				DevPortalApiKey:  os.Getenv("DEV_PORTAL_TOKEN"),
				Web3HttpProvider: os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY"),
				WalletKey:        os.Getenv("WALLET_KEY"),
				ChainId:          chains.Ethereum,
			},
			createOrderParams: orderbook.OrderRequest{
				SourceWallet: os.Getenv("WALLET_ADDRESS"),
				FromToken:    tokens.EthereumUsdc,
				ToToken:      tokens.EthereumDai,
				MakingAmount: 1000000,
				TakingAmount: 1000000000000000000,
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

			_, err := createOrder(tc.config, tc.createOrderParams)
			require.NoError(t, err)
		})
	}
}
