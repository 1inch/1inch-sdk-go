package client

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk/golang/client/orderbook"
	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/addresses"
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
		params                   orderbook.OrderRequest
		expectedOutput           string
		expectedErrorDescription string
	}{
		// TODO unsure of the value of validation tests when using the github.com/go-playground/validator/v10 library
		{
			description: "Error - missing fromToken",
			owner:       addresses.Vitalik,
			params: orderbook.OrderRequest{
				ToToken:      tokens.PolygonWeth,
				TakingAmount: "100",
				MakingAmount: "100",
			},
			expectedErrorDescription: `'FromToken' failed on the 'required' tag`,
		},
		{
			description: "Error - missing fromToken",
			owner:       addresses.Vitalik,
			params: orderbook.OrderRequest{
				FromToken:    tokens.PolygonDai,
				TakingAmount: "100",
				MakingAmount: "100",
			},
			expectedErrorDescription: `'ToToken' failed on the 'required' tag`,
		},
		// TODO commenting these out until we make validation logic for big ints represented as strings
		//{
		//	description: "Error - TakingAmount negative",
		//	owner:       addresses.Vitalik,
		//	params: orderbook.OrderRequest{
		//		FromToken:    tokens.PolygonDai,
		//		ToToken:      tokens.PolygonWeth,
		//		TakingAmount: "-1",
		//		MakingAmount: "100",
		//	},
		//	expectedErrorDescription: `'TakingAmount' failed on the 'gt' tag`,
		//},
		//{
		//	description: "Error - MakingAmount negative",
		//	owner:       addresses.Vitalik,
		//	params: orderbook.OrderRequest{
		//		FromToken:    tokens.PolygonDai,
		//		ToToken:      tokens.PolygonWeth,
		//		TakingAmount: "100",
		//		MakingAmount: "-1",
		//	},
		//	expectedErrorDescription: `'MakingAmount' failed on the 'gt' tag`,
		//},
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
		params                   orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams
		expectedOutput           string
		expectedErrorDescription string
	}{
		{
			description:    "Success",
			owner:          addresses.Vitalik,
			params:         orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{},
			expectedOutput: defaultSignature,
		},
		{
			description: "Error - Missing owner address",
			params: orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{
				Page: helpers.GetPtr(float32(0)),
			},
			expectedErrorDescription: "address must be a valid Ethereum address",
		},
		{
			description: "Error - Invalid owner address",
			owner:       "0x0",
			params: orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{
				Page: helpers.GetPtr(float32(0)),
			},
			expectedErrorDescription: "address must be a valid Ethereum address",
		},
		{
			description: "Error - Invalid page value",
			owner:       addresses.Vitalik,
			params: orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{
				Page: helpers.GetPtr(float32(0)),
			},
			expectedErrorDescription: "page must be greater than 0",
		},
		{
			description: "Error - Invalid limit value",
			owner:       addresses.Vitalik,
			params: orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{
				Limit: helpers.GetPtr(float32(0)),
			},
			expectedErrorDescription: "limit must be greater than 0",
		},
		{
			description: "Error - Invalid status",
			owner:       addresses.Vitalik,
			params: orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{
				Page:     helpers.GetPtr(float32(1)),
				Limit:    helpers.GetPtr(float32(2)),
				Statuses: &[]float32{0, 2},
			},
			expectedErrorDescription: "statuses can only contain 1, 2, and/or 3",
		},
		{
			description: "Error - Invalid sortBy",
			owner:       addresses.Vitalik,
			params: orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{
				Page:   helpers.GetPtr(float32(1)),
				Limit:  helpers.GetPtr(float32(2)),
				SortBy: orderbook.GetSortByParameter("invalid"),
			},
			expectedErrorDescription: "sortBy can only contain createDateTime, takerRate, makerRate, makerAmount, or takerAmount",
		},
		{
			description: "Error - Invalid takerAddress",
			owner:       addresses.Vitalik,
			params: orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{
				Page:       helpers.GetPtr(float32(1)),
				Limit:      helpers.GetPtr(float32(2)),
				TakerAsset: helpers.GetPtr("invalid"),
			},
			expectedErrorDescription: "takerAsset must be a valid Ethereum address",
		},
		{
			description: "Error - Invalid makerAddress",
			owner:       addresses.Vitalik,
			params: orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{
				Page:       helpers.GetPtr(float32(1)),
				Limit:      helpers.GetPtr(float32(2)),
				MakerAsset: helpers.GetPtr("invalid"),
			},
			expectedErrorDescription: "makerAsset must be a valid Ethereum address",
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

			allOrders, _, err := c.Orderbook.GetOrdersByCreatorAddress(context.Background(), tc.owner, tc.params)
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
		params                   orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams
		expectedOutput           string
		expectedErrorDescription string
	}{
		{
			description:    "Success",
			params:         orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{},
			expectedOutput: defaultSignature,
		},
		{
			description: "Error - Invalid page value",
			params: orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{
				Page: helpers.GetPtr(float32(0)),
			},
			expectedErrorDescription: "page must be greater than 0",
		},
		{
			description: "Error - Invalid limit value",
			params: orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{
				Limit: helpers.GetPtr(float32(0)),
			},
			expectedErrorDescription: "limit must be greater than 0",
		},
		{
			description: "Error - Invalid status",
			params: orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{
				Page:     helpers.GetPtr(float32(1)),
				Limit:    helpers.GetPtr(float32(2)),
				Statuses: &[]float32{0, 2},
			},
			expectedErrorDescription: "statuses can only contain 1, 2, and/or 3",
		},
		{
			description: "Error - Invalid sortBy",
			params: orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{
				Page:   helpers.GetPtr(float32(1)),
				Limit:  helpers.GetPtr(float32(2)),
				SortBy: orderbook.GetSortByParameter("invalid"),
			},
			expectedErrorDescription: "sortBy can only contain createDateTime, takerRate, makerRate, makerAmount, or takerAmount",
		},
		{
			description: "Error - Invalid takerAddress",
			params: orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{
				Page:       helpers.GetPtr(float32(1)),
				Limit:      helpers.GetPtr(float32(2)),
				TakerAsset: helpers.GetPtr("invalid"),
			},
			expectedErrorDescription: "takerAsset must be a valid Ethereum address",
		},
		{
			description: "Error - Invalid makerAddress",
			params: orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{
				Page:       helpers.GetPtr(float32(1)),
				Limit:      helpers.GetPtr(float32(2)),
				MakerAsset: helpers.GetPtr("invalid"),
			},
			expectedErrorDescription: "makerAsset must be a valid Ethereum address",
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
		params                   orderbook.LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams
		expectedOutput           int
		expectedErrorDescription string
	}{
		{
			description:    "Success",
			params:         orderbook.LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams{},
			expectedOutput: defaultCount,
		},
		{
			description: "Error - Duplicate statuses",
			params: orderbook.LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams{
				Statuses: []string{"1", "1"},
			},
			expectedErrorDescription: "statuses must not contain duplicates",
		},
		{
			description: "Error - Invalid status",
			params: orderbook.LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams{
				Statuses: []string{"4"},
			},
			expectedErrorDescription: "statuses can only contain 1, 2, and/or 3",
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
		orderHash                string
		expectedOutput           int
		expectedErrorDescription string
	}{
		{
			description:    "Success",
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

			eventResponse, _, err := c.Orderbook.GetEvent(context.Background(), tc.orderHash)
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
		params                   orderbook.LimitOrderV3SubscribedApiControllerGetEventsParams
		expectedOutput           int
		expectedErrorDescription string
	}{
		{
			description: "Success",
			params: orderbook.LimitOrderV3SubscribedApiControllerGetEventsParams{
				Limit: 1,
			},
			expectedOutput: 48608667,
		},
		{
			description: "Error - Limit too small",
			params: orderbook.LimitOrderV3SubscribedApiControllerGetEventsParams{
				Limit: 0,
			},
			expectedErrorDescription: "limit must be greater than 0",
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
func TestGetActiveOrdersWithPermit(t *testing.T) {

	endpoint := "/orderbook/v3.0/1/has-active-orders-with-permit/"
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
		wallet                   string
		token                    string
		expectedOutput           int
		expectedErrorDescription string
	}{
		//{
		//	description:    "Success",
		//	wallet:         addresses.Vitalik,
		//	token:          tokens.EthereumUsdc,
		//	expectedOutput: 48608667,
		//},
		{
			description:              "Error - Invalid wallet address",
			wallet:                   "0x123",
			token:                    tokens.EthereumUsdc,
			expectedErrorDescription: "wallet must be a valid Ethereum address",
		},
		{
			description:              "Error - Invalid token address",
			wallet:                   addresses.Vitalik,
			token:                    "0x123",
			expectedErrorDescription: "token must be a valid Ethereum address",
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

			ordersResponse, _, err := c.Orderbook.GetActiveOrdersWithPermit(context.Background(), tc.wallet, tc.token)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				require.Contains(t, err.Error(), tc.expectedErrorDescription)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, ordersResponse[0].Signature) // TODO use right values
		})
	}
}
