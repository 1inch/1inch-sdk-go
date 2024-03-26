package client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/client/models"
	"github.com/1inch/1inch-sdk-go/helpers/consts/addresses"
	"github.com/1inch/1inch-sdk-go/helpers/consts/chains"
	"github.com/1inch/1inch-sdk-go/helpers/consts/tokens"
)

func TestCreateOrder(t *testing.T) {

	endpoint := "/orderbook/v4.0/1/"
	defaultResponse := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[
	{
    	"Success": "true"
	}
]`)
	}

	testcases := []struct {
		description              string
		handlerFunc              func(w http.ResponseWriter, r *http.Request)
		owner                    string
		params                   models.CreateOrderParams
		expectedOutput           string
		expectedErrorDescription string
	}{
		{
			description: "Error - MakerAsset is native token",
			owner:       addresses.Vitalik,
			params: models.CreateOrderParams{
				ChainId:      chains.Ethereum,
				PrivateKey:   os.Getenv("WALLET_KEY"),
				Maker:        addresses.Vitalik,
				MakerAsset:   tokens.NativeToken,
				TakerAsset:   tokens.EthereumUsdc,
				TakingAmount: "1",
				MakingAmount: "1",
			},
			expectedErrorDescription: "native gas token is not supported as maker or taker asset",
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			c, mux, _, teardown, err := setup()
			require.NoError(t, err)
			defer teardown()

			if tc.handlerFunc != nil {
				mux.HandleFunc(endpoint, tc.handlerFunc)
			} else {
				mux.HandleFunc(endpoint, defaultResponse)
			}

			createOrderResponse, _, err := c.OrderbookApi.CreateOrder(context.Background(), tc.params)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				require.Contains(t, err.Error(), tc.expectedErrorDescription)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, createOrderResponse.Success)
		})
	}
}

func TestGetOrdersByCreatorAddress(t *testing.T) {

	endpoint := "/orderbook/v4.0/1/address/"
	defaultSignature := "0x68ca7fe03901373166c85d8e37587cb93bf9befc34c419dedd4d855e05c934e0731194f1b411a878f96a65f9fea7c61cca2e5e7d8c80dba9e35c3126b90b0c961b"
	defaultResponse := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[
	{
    	"Signature": "`+defaultSignature+`"
	}
]`)
	}

	testcases := []struct {
		description              string
		handlerFunc              func(w http.ResponseWriter, r *http.Request)
		owner                    string
		params                   models.GetOrdersByCreatorAddressParams
		expectedOutput           string
		expectedErrorDescription string
	}{
		{
			description: "Success",
			params: models.GetOrdersByCreatorAddressParams{
				ChainId:        chains.Ethereum,
				CreatorAddress: addresses.Vitalik,
			},
			expectedOutput: defaultSignature,
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			c, mux, _, teardown, err := setup()
			require.NoError(t, err)
			defer teardown()

			if tc.handlerFunc != nil {
				mux.HandleFunc(endpoint, tc.handlerFunc)
			} else {
				mux.HandleFunc(endpoint, defaultResponse)
			}

			allOrders, _, err := c.OrderbookApi.GetOrdersByCreatorAddress(context.Background(), tc.params)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				require.Contains(t, err.Error(), tc.expectedErrorDescription)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, allOrders[0].Signature)
		})
	}
}

func TestGetAllOrders(t *testing.T) {

	endpoint := "/orderbook/v4.0/1/all"
	defaultSignature := "0x68ca7fe03901373166c85d8e37587cb93bf9befc34c419dedd4d855e05c934e0731194f1b411a878f96a65f9fea7c61cca2e5e7d8c80dba9e35c3126b90b0c961b"
	defaultResponse := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[
	{
    	"Signature": "`+defaultSignature+`"
	}
]`)
	}

	testcases := []struct {
		description              string
		handlerFunc              func(w http.ResponseWriter, r *http.Request)
		params                   models.GetAllOrdersParams
		expectedOutput           string
		expectedErrorDescription string
	}{
		{
			description: "Success",
			params: models.GetAllOrdersParams{
				ChainId: chains.Ethereum,
				LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams: models.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{},
			},
			expectedOutput: defaultSignature,
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			c, mux, _, teardown, err := setup()
			require.NoError(t, err)
			defer teardown()

			if tc.handlerFunc != nil {
				mux.HandleFunc(endpoint, tc.handlerFunc)
			} else {
				mux.HandleFunc(endpoint, defaultResponse)
			}

			allOrders, _, err := c.OrderbookApi.GetAllOrders(context.Background(), tc.params)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				require.Contains(t, err.Error(), tc.expectedErrorDescription)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, allOrders[0].Signature)
		})
	}
}

func TestGetCount(t *testing.T) {

	endpoint := "/orderbook/v4.0/1/count"
	defaultCount := 1
	defaultResponse := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
    	"count": 1
}`)
	}

	testcases := []struct {
		description              string
		handlerFunc              func(w http.ResponseWriter, r *http.Request)
		params                   models.GetCountParams
		expectedOutput           int
		expectedErrorDescription string
	}{
		{
			description: "Success",
			params: models.GetCountParams{
				ChainId: chains.Ethereum,
				LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams: models.LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams{},
			},
			expectedOutput: defaultCount,
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			c, mux, _, teardown, err := setup()
			require.NoError(t, err)
			defer teardown()

			if tc.handlerFunc != nil {
				mux.HandleFunc(endpoint, tc.handlerFunc)
			} else {
				mux.HandleFunc(endpoint, defaultResponse)
			}

			countResponse, _, err := c.OrderbookApi.GetCount(context.Background(), tc.params)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				require.Contains(t, err.Error(), tc.expectedErrorDescription)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, countResponse.Count)
		})
	}
}

func TestGetEvent(t *testing.T) {

	endpoint := "/orderbook/v4.0/1/events/"
	defaultResponse := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
        	"id": 48608667
}`)
	}

	testcases := []struct {
		description              string
		handlerFunc              func(w http.ResponseWriter, r *http.Request)
		params                   models.GetEventParams
		expectedOutput           int
		expectedErrorDescription string
	}{
		{
			description:    "Success",
			expectedOutput: 48608667,
			params: models.GetEventParams{
				ChainId:   chains.Ethereum,
				OrderHash: "123",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			c, mux, _, teardown, err := setup()
			require.NoError(t, err)
			defer teardown()

			if tc.handlerFunc != nil {
				mux.HandleFunc(endpoint, tc.handlerFunc)
			} else {
				mux.HandleFunc(endpoint, defaultResponse)
			}

			eventResponse, _, err := c.OrderbookApi.GetEvent(context.Background(), tc.params)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				require.Contains(t, err.Error(), tc.expectedErrorDescription)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, eventResponse.Id)
		})
	}
}

func TestGetEvents(t *testing.T) {

	endpoint := "/orderbook/v4.0/1/events"
	defaultResponse := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[
		{
        	"id": 48608667
		}
]`)
	}

	testcases := []struct {
		description              string
		handlerFunc              func(w http.ResponseWriter, r *http.Request)
		params                   models.GetEventsParams
		expectedOutput           int
		expectedErrorDescription string
	}{
		{
			description: "Success",
			params: models.GetEventsParams{
				ChainId: chains.Ethereum,
				LimitOrderV3SubscribedApiControllerGetEventsParams: models.LimitOrderV3SubscribedApiControllerGetEventsParams{
					Limit: 1,
				},
			},
			expectedOutput: 48608667,
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			c, mux, _, teardown, err := setup()
			require.NoError(t, err)
			defer teardown()

			if tc.handlerFunc != nil {
				mux.HandleFunc(endpoint, tc.handlerFunc)
			} else {
				mux.HandleFunc(endpoint, defaultResponse)
			}

			eventsResponse, _, err := c.OrderbookApi.GetEvents(context.Background(), tc.params)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				require.Contains(t, err.Error(), tc.expectedErrorDescription)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, eventsResponse[0].Id)
		})
	}
}

// TODO fix tests once docs are available
//func TestGetActiveOrdersWithPermit(t *testing.T) {
//
//	endpoint := "/orderbook/v4.0/1/has-active-orders-with-permit/"
//	defaultResponse := func(w http.ResponseWriter, r *http.Request) {
//		fmt.Fprint(w, `[
//		{
//        	"id": 48608667
//		}
//]`)
//	}
//
//	testcases := []struct {
//		description              string
//		handlerFunc              func(w http.ResponseWriter, r *http.Request)
//		params                   orderbook.GetActiveOrdersWithPermitParams
//		expectedOutput           int
//		expectedErrorDescription string
//	}{
//		{
//			description:    "Success",
//			wallet:         addresses.Vitalik,
//			token:          tokens.EthereumUsdc,
//			expectedOutput: 48608667,
//		},
//	}
//
//	for _, tc := range testcases {
//		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {
//
//			c, mux, _, teardown, err := setup()
//			require.NoError(t, err)
//			defer teardown()
//
//			if tc.handlerFunc != nil {
//				mux.HandleFunc(endpoint, tc.handlerFunc)
//			} else {
//				mux.HandleFunc(endpoint, defaultResponse)
//			}
//
//			ordersResponse, _, err := c.OrderbookApi.GetActiveOrdersWithPermit(context.Background(), tc.params)
//			if tc.expectedErrorDescription != "" {
//				if err == nil {
//					assert.FailNow(t, "Expected error message, but error was nil")
//				}
//				require.Contains(t, err.Error(), tc.expectedErrorDescription)
//				return
//			}
//			require.NoError(t, err)
//			assert.Equal(t, tc.expectedOutput, ordersResponse[0].Signature) // TODO use right values
//		})
//	}
//}
