package client

import (
	"context"
	"fmt"
	helpers2 "github.com/1inch/1inch-sdk-go/internal/helpers"
	"github.com/1inch/1inch-sdk-go/internal/helpers/consts/amounts"
	"github.com/1inch/1inch-sdk-go/internal/helpers/consts/chains"
	"github.com/1inch/1inch-sdk-go/internal/helpers/consts/tokens"
	"os"
	"testing"

	"github.com/1inch/1inch-sdk-go/client/models"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/helpers"
)

func TestCreateOrderIntegration(t *testing.T) {

	testcases := []struct {
		description   string
		orderRequest  models.CreateOrderParams
		expectedError string
	}{
		{
			description: "Success",
			orderRequest: models.CreateOrderParams{
				ChainId:      chains.Polygon,
				PrivateKey:   os.Getenv("WALLET_KEY_EMPTY"),
				Maker:        os.Getenv("WALLET_ADDRESS_EMPTY"),
				MakerAsset:   tokens.PolygonFrax,
				TakerAsset:   tokens.PolygonWeth,
				TakingAmount: amounts.Ten6,
				MakingAmount: amounts.Ten18,
				Taker:        "0x0000000000000000000000000000000000000000",
				SkipWarnings: true,
			},
		},
		{
			description: "Failure - No approval with error",
			orderRequest: models.CreateOrderParams{
				ChainId:                        chains.Polygon,
				PrivateKey:                     os.Getenv("WALLET_KEY_EMPTY"),
				Maker:                          os.Getenv("WALLET_ADDRESS_EMPTY"),
				MakerAsset:                     tokens.PolygonDai,
				TakerAsset:                     tokens.PolygonWeth,
				TakingAmount:                   amounts.Ten6,
				MakingAmount:                   amounts.Ten18,
				Taker:                          "0x0000000000000000000000000000000000000000",
				SkipWarnings:                   true,
				EnableOnchainApprovalsIfNeeded: false,
			},
			expectedError: models.ErrorFailWhenApprovalIsNeeded.Error(),
		},
	}

	c, err := NewClient(models.ClientConfig{
		DevPortalApiKey: helpers2.GetenvSafe("DEV_PORTAL_TOKEN"),
		Web3HttpProviders: []models.Web3Provider{
			{
				ChainId: chains.Polygon,
				Url:     helpers2.GetenvSafe("WEB_3_HTTP_PROVIDER_URL_WITH_KEY_POLYGON"),
			},
		},
	})
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {

			t.Cleanup(func() {
				helpers.Sleep()
			})

			orderResponse, resp, err := c.OrderbookApi.CreateOrder(context.Background(), tc.orderRequest)
			if tc.expectedError != "" {
				require.Equal(t, err.Error(), tc.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, 201, resp.StatusCode)

				require.NotNil(t, orderResponse)
				require.True(t, orderResponse.Success)
			}
		})
	}
}

func TestGetAllOrdersIntegration(t *testing.T) {

	testcases := []struct {
		description string
		params      models.GetAllOrdersParams
	}{
		{
			description: "Get all orders",
			params: models.GetAllOrdersParams{
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

			orders, resp, err := c.OrderbookApi.GetAllOrders(context.Background(), tc.params)
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
		params      models.GetCountParams
	}{
		{
			description: "Get counts",
			params: models.GetCountParams{
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

			countResponse, resp, err := c.OrderbookApi.GetCount(context.Background(), tc.params)
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
		params      models.GetEventsParams
	}{
		{
			description: "Get events",
			params: models.GetEventsParams{
				ChainId: chains.Ethereum,
				LimitOrderV3SubscribedApiControllerGetEventsParams: models.LimitOrderV3SubscribedApiControllerGetEventsParams{
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

			events, resp, err := c.OrderbookApi.GetEvents(context.Background(), tc.params)
			require.NoError(t, err)
			require.Equal(t, 200, resp.StatusCode)

			require.NotNil(t, events)
			require.Greater(t, len(events), 0)
		})
	}
}
