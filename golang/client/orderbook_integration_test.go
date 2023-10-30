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

			orders, resp, err := c.Orderbook.GetAllOrders(context.Background(), orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{})
			require.NoError(t, err)
			require.Equal(t, 200, resp.StatusCode)

			require.NotNil(t, orders)
			require.Greater(t, len(orders), 0)
			helpers.Sleep()
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

			countResponse, resp, err := c.Orderbook.GetCount(context.Background(), orderbook.LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams{})
			require.NoError(t, err)
			require.Equal(t, 200, resp.StatusCode)

			require.NotNil(t, countResponse)
			require.NotEqual(t, countResponse.Count, 0)
			helpers.Sleep()
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

			events, resp, err := c.Orderbook.GetEvents(context.Background(), orderbook.LimitOrderV3SubscribedApiControllerGetEventsParams{
				Limit: 1,
			})
			require.NoError(t, err)
			require.Equal(t, 200, resp.StatusCode)

			require.NotNil(t, events)
			require.Greater(t, len(events), 0)
			helpers.Sleep()
		})
	}
}
