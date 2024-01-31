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
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/tokens"
)

func TestCreateOrderIntegration(t *testing.T) {

	testcases := []struct {
		description   string
		orderRequest  orderbook.OrderRequest
		expectSuccess bool
	}{
		{
			description: "Success",
			orderRequest: orderbook.OrderRequest{
				ChainId:      chains.Polygon,
				WalletKey:    os.Getenv("WALLET_KEY"),
				SourceWallet: os.Getenv("WALLET_ADDRESS"),
				FromToken:    tokens.PolygonDai,
				ToToken:      tokens.PolygonWeth,
				TakingAmount: "1000000000000000000",
				MakingAmount: "1000000000000000000",
				Receiver:     "0x0000000000000000000000000000000000000000",
				SkipWarnings: true,
			},
			expectSuccess: true,
		},
	}

	c, err := NewClient(Config{
		DevPortalApiKey: helpers.GetenvSafe("DEV_PORTAL_TOKEN"),
		Web3HttpProviders: []Web3ProviderConfig{
			{
				ChainId: chains.Polygon,
				Url:     helpers.GetenvSafe("WEB_3_HTTP_PROVIDER_URL_WITH_KEY_POLYGON"),
			},
		},
	})
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			t.Cleanup(func() {
				helpers.Sleep()
			})

			orderResponse, resp, err := c.Orderbook.CreateOrder(context.Background(), tc.orderRequest)
			require.NoError(t, err)
			require.Equal(t, 201, resp.StatusCode)

			require.NotNil(t, orderResponse)
			require.True(t, orderResponse.Success)
		})
	}
}

func TestGetAllOrdersIntegration(t *testing.T) {

	testcases := []struct {
		description string
		params      orderbook.GetAllOrdersParams
	}{
		{
			description: "Get all orders",
			params: orderbook.GetAllOrdersParams{
				ChainId: chains.Ethereum,
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

			orders, resp, err := c.Orderbook.GetAllOrders(context.Background(), tc.params)
			require.NoError(t, err)
			require.Equal(t, 200, resp.StatusCode)

			require.NotNil(t, orders)
			require.Greater(t, len(orders), 0)
		})
	}
}

func TestGetCountIntegration(t *testing.T) {

	testcases := []struct {
		description string
		params      orderbook.GetCountParams
	}{
		{
			description: "Get counts",
			params: orderbook.GetCountParams{
				ChainId: chains.Ethereum,
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

			countResponse, resp, err := c.Orderbook.GetCount(context.Background(), tc.params)
			require.NoError(t, err)
			require.Equal(t, 200, resp.StatusCode)

			require.NotNil(t, countResponse)
			require.NotEqual(t, countResponse.Count, 0)
		})
	}
}

func TestGetEventsIntegration(t *testing.T) {

	testcases := []struct {
		description string
		params      orderbook.GetEventsParams
	}{
		{
			description: "Get events",
			params: orderbook.GetEventsParams{
				ChainId: chains.Ethereum,
				LimitOrderV3SubscribedApiControllerGetEventsParams: orderbook.LimitOrderV3SubscribedApiControllerGetEventsParams{
					Limit: 1,
				},
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

			events, resp, err := c.Orderbook.GetEvents(context.Background(), tc.params)
			require.NoError(t, err)
			require.Equal(t, 200, resp.StatusCode)

			require.NotNil(t, events)
			require.Greater(t, len(events), 0)
		})
	}
}
