package fusion

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/common"
	web3_provider "github.com/1inch/1inch-sdk-go/internal/web3-provider"
)

// capturingHttpExecutor records every request and replays queued responses
type capturingHttpExecutor struct {
	Payloads  []common.RequestPayload
	Responses []any
}

func (m *capturingHttpExecutor) ExecuteRequest(ctx context.Context, payload common.RequestPayload, v any) error {
	m.Payloads = append(m.Payloads, payload)
	if len(m.Responses) == 0 {
		return nil
	}
	response := m.Responses[0]
	m.Responses = m.Responses[1:]
	if response != nil && v != nil {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Ptr || rv.IsNil() {
			return fmt.Errorf("v must be a non-nil pointer")
		}
		reflect.Indirect(rv).Set(reflect.ValueOf(response))
	}
	return nil
}

func TestPlaceOrderFromParams(t *testing.T) {
	testPrivateKey := "d8d1f95deb28949ea0ecc4e9a0decf89e98422c2d76ab6e5f736792a388c56c7"
	wallet, err := web3_provider.DefaultWalletOnlyProvider(testPrivateKey, 1)
	require.NoError(t, err)

	quote := GetQuoteOutputFixed{
		QuoteId:           "test-quote-id",
		SettlementAddress: extensionContract,
		Whitelist:         []string{"0x00000000219ab540356cbb839cbe05303d7705fa"},
		MarketAmount:      "1420000000",
		SurplusFee:        0,
		Presets: QuotePresetsClassFixed{
			Fast: PresetClassFixed{
				AllowMultipleFills: true,
				AllowPartialFills:  true,
				AuctionDuration:    180,
				AuctionEndAmount:   "1420000000",
				AuctionStartAmount: "1500000000",
				GasCost:            GasCostConfigClass{GasBumpEstimate: 0, GasPriceEstimate: "0"},
				InitialRateBump:    50000,
				Points:             []AuctionPointClass{{Coefficient: 20000, Delay: 12}},
				StartAuctionIn:     0,
			},
		},
	}

	executor := &capturingHttpExecutor{Responses: []any{quote, nil}}
	client := &Client{
		api:    api{chainId: 1, httpExecutor: executor},
		Wallet: wallet,
	}

	orderParams := OrderParams{
		FromTokenAddress:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		ToTokenAddress:     "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		Amount:             "1000000000000000000",
		WalletAddress:      strings.ToLower(wallet.Address().Hex()),
		Receiver:           "0x0000000000000000000000000000000000000000",
		Preset:             Fast,
		Permit:             "0xdeadbeef01020304",
		IsPermit2:          true,
		AllowPartialFills:  true,
		AllowMultipleFills: true,
	}

	orderHash, err := client.PlaceOrderFromParams(context.Background(), orderParams)
	require.NoError(t, err)
	assert.NotEmpty(t, orderHash)
	require.Len(t, executor.Payloads, 2)

	// The quote request must carry the permit settings from the single input
	quoteParams, ok := executor.Payloads[0].Params.(QuoterControllerGetQuoteParamsFixed)
	require.True(t, ok, "first request must be the quote request")
	assert.Equal(t, "true", quoteParams.IsPermit2)
	assert.Equal(t, orderParams.Permit, quoteParams.Permit)
	assert.Equal(t, orderParams.Amount, quoteParams.Amount)
	assert.True(t, quoteParams.EnableEstimate)
	assert.True(t, quoteParams.Surplus)

	// The submitted order must embed the permit in its extension
	submission := executor.Payloads[1]
	assert.Contains(t, submission.U, "/order/submit")
	assert.Contains(t, string(submission.Body), "deadbeef01020304")
}
