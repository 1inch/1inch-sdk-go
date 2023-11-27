package client

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"1inch-sdk-golang/client/orderbook"
	"1inch-sdk-golang/helpers"
)

// TODO This test requires a private key to be set in the environment and sign the limit order
//func TestCreateOrderIntegration(t *testing.T) {
//
//	testcases := []struct {
//		description   string
//		orderRequest  orderbook.OrderRequest
//		expectSuccess bool
//	}{
//		{
//			description: "Success",
//			orderRequest: orderbook.OrderRequest{
//				SourceWallet: os.Getenv("WALLET_ADDRESS"),
//				FromToken:    tokens.PolygonDai,
//				ToToken:      tokens.PolygonWeth,
//				TakingAmount: 1000000000000000000,
//				MakingAmount: 1000000000000000000,
//				Receiver:     "0x0000000000000000000000000000000000000000",
//			},
//			expectSuccess: true,
//		},
//	}
//
//	c, err := NewClient(Config{
//		TargetEnvironment: EnvironmentProduction,
//		ChainId:           chains.Polygon,
//		DevPortalApiKey:   helpers.GetenvSafe("DEV_PORTAL_TOKEN"),
//		WalletKey:         helpers.GetenvSafe("WALLET_KEY"),
//	})
//	require.NoError(t, err)
//
//	for _, tc := range testcases {
//		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {
//
//			t.Cleanup(func() {
//				helpers.Sleep()
//			})
//
//			orderResponse, resp, err := c.Orderbook.CreateOrder(context.Background(), tc.orderRequest)
//			require.NoError(t, err)
//			require.Equal(t, 201, resp.StatusCode)
//
//			require.NotNil(t, orderResponse)
//			require.True(t, orderResponse.Success)
//		})
//	}
//}

func TestGetAllOrdersIntegration(t *testing.T) {

	testcases := []struct {
		description string
	}{
		{
			description: "Get all orders",
		},
	}

	c, err := NewClient(Config{
		TargetEnvironment: EnvironmentProduction,
		DevPortalApiKey:   os.Getenv("DEV_PORTAL_TOKEN"),
	})
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			t.Cleanup(func() {
				helpers.Sleep()
			})

			orders, resp, err := c.Orderbook.GetAllOrders(context.Background(), orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{})
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
	}{
		{
			description: "Get counts",
		},
	}

	c, err := NewClient(Config{
		TargetEnvironment: EnvironmentProduction,
		DevPortalApiKey:   os.Getenv("DEV_PORTAL_TOKEN"),
	})
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			t.Cleanup(func() {
				helpers.Sleep()
			})

			countResponse, resp, err := c.Orderbook.GetCount(context.Background(), orderbook.LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams{})
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
	}{
		{
			description: "Get all events",
		},
	}

	c, err := NewClient(Config{
		TargetEnvironment: EnvironmentProduction,
		DevPortalApiKey:   os.Getenv("DEV_PORTAL_TOKEN"),
	})
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			t.Cleanup(func() {
				helpers.Sleep()
			})

			events, resp, err := c.Orderbook.GetEvents(context.Background(), orderbook.LimitOrderV3SubscribedApiControllerGetEventsParams{
				Limit: 1,
			})
			require.NoError(t, err)
			require.Equal(t, 200, resp.StatusCode)

			require.NotNil(t, events)
			require.Greater(t, len(events), 0)
		})
	}
}
