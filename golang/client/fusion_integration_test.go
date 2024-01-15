package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk/golang/client/fusion"
	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/addresses"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/tokens"
)

func TestGetOrdersIntegration(t *testing.T) {

	testcases := []struct {
		description   string
		params        fusion.OrderApiControllerGetActiveOrdersParams
		expectSuccess bool
	}{
		{
			description: "Success",
			params: fusion.OrderApiControllerGetActiveOrdersParams{
				Page:  helpers.GetPtr(float32(1)),
				Limit: helpers.GetPtr(float32(1)),
			},
			expectSuccess: true,
		},
	}

	c, err := NewClient(SimpleEthereumConfig)
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			t.Cleanup(func() {
				helpers.Sleep()
			})

			getOrdersResponse, resp, err := c.Fusion.GetOrders(context.Background(), tc.params)
			require.NoError(t, err)
			require.Equal(t, 200, resp.StatusCode)

			require.NotNil(t, getOrdersResponse)
			require.True(t, len(getOrdersResponse.Items) > 0)
		})
	}
}

func TestGetSettlementContractIntegration(t *testing.T) {

	testcases := []struct {
		description   string
		expectSuccess bool
	}{
		{
			description:   "Success",
			expectSuccess: true,
		},
	}

	c, err := NewClient(SimpleEthereumConfig)
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			t.Cleanup(func() {
				helpers.Sleep()
			})

			settlementContractResponse, resp, err := c.Fusion.GetSettlementContract(context.Background())
			require.NoError(t, err)
			require.Equal(t, 200, resp.StatusCode)

			require.NotNil(t, settlementContractResponse)
			require.True(t, len(settlementContractResponse.Address) > 0)
		})
	}
}

func TestGetQuoteIntegration(t *testing.T) {

	testcases := []struct {
		description   string
		params        fusion.QuoterControllerGetQuoteParams
		expectSuccess bool
	}{
		{
			description: "Success",
			params: fusion.QuoterControllerGetQuoteParams{
				FromTokenAddress: tokens.EthereumWeth,
				ToTokenAddress:   tokens.EthereumUsdc,
				Amount:           1000000000000000000,
				WalletAddress:    addresses.Vitalik,
				EnableEstimate:   false,
				Fee:              nil,
				IsLedgerLive:     false,
				Permit:           nil,
			},
			expectSuccess: true,
		},
	}

	c, err := NewClient(SimpleEthereumConfig)
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			t.Cleanup(func() {
				helpers.Sleep()
			})

			getQuoteResponse, resp, err := c.Fusion.GetQuote(context.Background(), tc.params)
			require.NoError(t, err)
			require.Equal(t, 200, resp.StatusCode)

			require.NotNil(t, getQuoteResponse)
			require.Equal(t, tokens.EthereumUsdc, getQuoteResponse.FeeToken)
		})
	}
}

func TestGetQuoteWithCustomPresetsIntegration(t *testing.T) {

	testcases := []struct {
		description   string
		params        fusion.QuoterControllerGetQuoteParams
		body          fusion.CustomPresetInput
		expectSuccess bool
	}{
		{
			description: "Success",
			params: fusion.QuoterControllerGetQuoteParams{
				FromTokenAddress: tokens.EthereumWeth,
				ToTokenAddress:   tokens.EthereumUsdc,
				Amount:           1000000000000000000,
				WalletAddress:    addresses.Vitalik,
				EnableEstimate:   false,
				Fee:              nil,
				IsLedgerLive:     false,
				Permit:           nil,
			},
			body: fusion.CustomPresetInput{
				AuctionDuration:    600,
				AuctionStartAmount: "2075314488",
				AuctionEndAmount:   "2021177846",
			},
			expectSuccess: true,
		},
	}

	c, err := NewClient(SimpleEthereumConfig)
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			t.Cleanup(func() {
				helpers.Sleep()
			})

			getQuoteResponse, resp, err := c.Fusion.GetQuoteWithCustomPresets(context.Background(), tc.params, tc.body)
			require.NoError(t, err)
			require.Equal(t, 201, resp.StatusCode) // incorrect code, but this is what the API returns

			helpers.PrettyPrintStruct(getQuoteResponse)

			require.NotNil(t, getQuoteResponse)
			require.Equal(t, tokens.EthereumUsdc, getQuoteResponse.FeeToken)
		})
	}
}
