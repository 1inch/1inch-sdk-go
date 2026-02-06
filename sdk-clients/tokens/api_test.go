package tokens

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/common"
	"github.com/1inch/1inch-sdk-go/constants"
)

type MockHttpExecutor struct {
	Called      bool
	ExecuteErr  error
	ResponseObj any
}

func (m *MockHttpExecutor) ExecuteRequest(ctx context.Context, payload common.RequestPayload, v any) error {
	m.Called = true
	if m.ExecuteErr != nil {
		return m.ExecuteErr
	}

	// Copy the mock response object to v
	if m.ResponseObj != nil && v != nil {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Ptr || rv.IsNil() {
			return fmt.Errorf("v must be a non-nil pointer")
		}
		reflect.Indirect(rv).Set(reflect.ValueOf(m.ResponseObj))
	}
	return nil
}

func getPtr[T any](v T) *T {
	return &v
}

func TestGetPricesForWhitelistedTokens(t *testing.T) {
	ctx := context.Background()

	mockedResp := []ProviderTokenDtoFixed{
		{
			Address:         "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984",
			ChainId:         1,
			Decimals:        18,
			DisplayedSymbol: nil,
			Eip2612:         getPtr(true),
			IsFoT:           nil,
			LogoURI:         getPtr("https://tokens.1inch.io/0x1f9840a85d5af5bf1d1762f925bdaddc4201f984.png"),
			Name:            "Uniswap",
			Providers: []string{
				"1inch",
				"Arb Whitelist Era",
				"CMC DeFi",
				"CoinGecko",
				"Compound",
				"Curve Token List",
				"Defiprime",
				"Furucombo",
				"Gemini Token List",
				"Kleros Tokens",
				"MyCrypto Token List",
				"Trust Wallet Assets",
				"Uniswap Labs Default",
				"Zerion",
			},
			Symbol: "UNI",
			Tags: []TagDto{
				{
					Value:    "defi",
					Provider: "Defiprime",
				},
				{
					Value:    "farming",
					Provider: "Zerion",
				},
				{
					Value:    "tokens",
					Provider: "1inch",
				},
			},
		},
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
		chainId:      constants.EthereumChainId,
	}

	params := SearchControllerSearchAllChainsParams{
		Query:              "UNI",
		IgnoreListed:       false,
		OnlyPositiveRating: false,
		Limit:              float32(1),
	}

	tokens, err := api.SearchTokenAllChains(ctx, params)
	require.NoError(t, err)
	require.NotNil(t, tokens)

	if !mockExecutor.Called {
		t.Errorf("Expected ExecuteRequest to be called")
	}
	if !reflect.DeepEqual(tokens, mockedResp) {
		t.Errorf("Expected swap to be %+v, got %+v", tokens, mockedResp)
	}
	require.NoError(t, err)
}
