package client

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk/golang/client/orderbook"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/addresses"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/tokens"
)

func TestCreateOrder(t *testing.T) {

	endpoint := "/orderbook/v3.0/1/"
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
		params                   orderbook.CreateOrderParams
		expectedOutput           string
		expectedErrorDescription string
	}{
		{
			description: "Error - MakerAsset is native token",
			owner:       addresses.Vitalik,
			params: orderbook.CreateOrderParams{
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

			createOrderResponse, _, err := c.Orderbook.CreateOrder(context.Background(), tc.params)
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

	endpoint := "/orderbook/v3.0/1/address/"
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
		params                   orderbook.GetOrdersByCreatorAddressParams
		expectedOutput           string
		expectedErrorDescription string
	}{
		{
			description: "Success",
			params: orderbook.GetOrdersByCreatorAddressParams{
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

			allOrders, _, err := c.Orderbook.GetOrdersByCreatorAddress(context.Background(), tc.params)
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

	endpoint := "/orderbook/v3.0/1/all"
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
		params                   orderbook.GetAllOrdersParams
		expectedOutput           string
		expectedErrorDescription string
	}{
		{
			description: "Success",
			params: orderbook.GetAllOrdersParams{
				ChainId: chains.Ethereum,
				LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams: orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{},
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

			allOrders, _, err := c.Orderbook.GetAllOrders(context.Background(), tc.params)
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

	endpoint := "/orderbook/v3.0/1/count"
	defaultCount := 1
	defaultResponse := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
    	"count": 1
}`)
	}

	testcases := []struct {
		description              string
		handlerFunc              func(w http.ResponseWriter, r *http.Request)
		params                   orderbook.GetCountParams
		expectedOutput           int
		expectedErrorDescription string
	}{
		{
			description: "Success",
			params: orderbook.GetCountParams{
				ChainId: chains.Ethereum,
				LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams: orderbook.LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams{},
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

			countResponse, _, err := c.Orderbook.GetCount(context.Background(), tc.params)
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

	endpoint := "/orderbook/v3.0/1/events/"
	defaultResponse := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
        	"id": 48608667
}`)
	}

	testcases := []struct {
		description              string
		handlerFunc              func(w http.ResponseWriter, r *http.Request)
		params                   orderbook.GetEventParams
		expectedOutput           int
		expectedErrorDescription string
	}{
		{
			description:    "Success",
			expectedOutput: 48608667,
			params: orderbook.GetEventParams{
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

			eventResponse, _, err := c.Orderbook.GetEvent(context.Background(), tc.params)
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

	endpoint := "/orderbook/v3.0/1/events"
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
		params                   orderbook.GetEventsParams
		expectedOutput           int
		expectedErrorDescription string
	}{
		{
			description: "Success",
			params: orderbook.GetEventsParams{
				ChainId: chains.Ethereum,
				LimitOrderV3SubscribedApiControllerGetEventsParams: orderbook.LimitOrderV3SubscribedApiControllerGetEventsParams{
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

			eventsResponse, _, err := c.Orderbook.GetEvents(context.Background(), tc.params)
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
//	endpoint := "/orderbook/v3.0/1/has-active-orders-with-permit/"
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
//			ordersResponse, _, err := c.Orderbook.GetActiveOrdersWithPermit(context.Background(), tc.params)
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
