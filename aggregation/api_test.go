package aggregation

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/1inch/1inch-sdk-go/internal/common"
	"github.com/1inch/1inch-sdk-go/internal/helpers/consts/chains"
)

type MockHttpExecutor struct {
	Called      bool
	ExecuteErr  error
	ResponseObj interface{}
}

func (m *MockHttpExecutor) ExecuteRequest(ctx context.Context, payload common.RequestPayload, v interface{}) error {
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

func TestGetQuote(t *testing.T) {
	ctx := context.Background()

	mockedResp := QuoteResponse{
		FromToken: &TokenInfo{
			Address:  "0x6b175474e89094c44da98b954eedeac495271d0f",
			Symbol:   "DAI",
			Name:     "Dai Stablecoin",
			Decimals: 18,
			LogoURI:  "https://tokens.1inch.io/0x6b175474e89094c44da98b954eedeac495271d0f.png",
			Tags: []string{
				"PEG:USD",
				"tokens",
			},
		},
		Gas:      181416,
		ToAmount: "289424403260095",
		ToToken: &TokenInfo{
			Address:  "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			Symbol:   "WETH",
			Name:     "Wrapped Ether",
			Decimals: 18,
			LogoURI:  "https://tokens.1inch.io/0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2.png",
			Tags: []string{
				"PEG:ETH",
				"tokens",
			},
		},
		Protocols: [][][]SelectedProtocol{
			{
				{
					{
						FromTokenAddress: "0x6b175474e89094c44da98b954eedeac495271d0f",
						Name:             "SUSHI",
						Part:             100,
						ToTokenAddress:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
					},
				},
			},
		},
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}
	api := api{httpExecutor: mockExecutor}

	params := GetQuoteParams{
		ChainId: chains.Ethereum,
		AggregationControllerGetQuoteParams: AggregationControllerGetQuoteParams{
			Src:               "0x6b175474e89094c44da98b954eedeac495271d0f",
			Dst:               "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			Amount:            "1000000000000000000",
			IncludeTokensInfo: true,
			IncludeGas:        true,
			IncludeProtocols:  true,
		},
	}

	quote, err := api.GetQuote(ctx, params)
	if err != nil {
		t.Fatalf("GetQuote returned an error: %v", err)
	}

	if !mockExecutor.Called {
		t.Errorf("Expected ExecuteRequest to be called")
	}

	expectedQuote := mockedResp
	if !reflect.DeepEqual(*quote, expectedQuote) {
		t.Errorf("Expected quote to be %+v, got %+v", expectedQuote, *quote)
	}
}
