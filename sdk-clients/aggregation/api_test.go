package aggregation

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestGetQuote(t *testing.T) {
	ctx := context.Background()

	mockedResp := QuoteResponse{
		SrcToken: &TokenInfo{
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
		Gas:       181416,
		DstAmount: "289424403260095",
		DstToken: &TokenInfo{
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
	api := api{httpExecutor: mockExecutor, chainId: constants.EthereumChainId}

	params := GetQuoteParams{
		Src:               "0x6b175474e89094c44da98b954eedeac495271d0f",
		Dst:               "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		Amount:            "1000000000000000000",
		IncludeTokensInfo: true,
		IncludeGas:        true,
		IncludeProtocols:  true,
	}

	quote, err := api.GetQuote(ctx, params)
	require.NoError(t, err)

	if !mockExecutor.Called {
		t.Errorf("Expected ExecuteRequest to be called")
	}

	expectedQuote := mockedResp
	if !reflect.DeepEqual(*quote, expectedQuote) {
		t.Errorf("Expected quote to be %+v, got %+v", expectedQuote, *quote)
	}
}

func TestGetSwap(t *testing.T) {
	ctx := context.Background()

	mockedResp := SwapResponse{
		SrcToken: &TokenInfo{
			Address:  "0x5a98fcbea516cf06857215779fd812ca3bef1b32",
			Symbol:   "LDO",
			Name:     "Lido DAO Token",
			Decimals: 18,
			LogoURI:  "https://tokens.1inch.io/0x5a98fcbea516cf06857215779fd812ca3bef1b32.png",
			Tags: []string{
				"tokens",
			},
		},
		DstAmount: "6",
		DstToken: &TokenInfo{
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
						FromTokenAddress: "0x5a98fcbea516cf06857215779fd812ca3bef1b32",
						Name:             "SUSHI",
						Part:             100,
						ToTokenAddress:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
					},
				},
			},
		},
		Tx: TransactionData{
			Data:     "0x0502b1c50000000000000000000000005a98fcbea516cf06857215779fd812ca3bef1b32000000000000000000000000000000000000000000000000000000000000271000000000000000000000000000000000000000000000000000000000000000050000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000100000000000000003b6d0340c558f600b34a5f69dd2f0d06cb8a88d829b7420ade8bb62d",
			From:     "0x083fc10ce7e97cafbae0fe332a9c4384c5f54e45",
			Gas:      257615,
			GasPrice: "22800337026",
			To:       "0x1111111254eeb25477b68fb85ed929f73a960582",
			Value:    "0",
		},
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
		chainId:      constants.EthereumChainId,
	}

	params := GetSwapParams{
		Src:               "0x5a98fcbea516cf06857215779fd812ca3bef1b32",
		Dst:               "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		Amount:            "10000",
		Slippage:          1,
		From:              "0x083fc10ce7e97cafbae0fe332a9c4384c5f54e45",
		IncludeTokensInfo: true,
		IncludeGas:        true,
		IncludeProtocols:  true,
	}

	swap, err := api.GetSwap(ctx, params)
	require.NoError(t, err)

	if !mockExecutor.Called {
		t.Errorf("Expected ExecuteRequest to be called")
	}

	expectedSwap, err := normalizeSwapResponse(mockedResp)
	if !reflect.DeepEqual(swap, expectedSwap) {
		t.Errorf("Expected swap to be %+v, got %+v", expectedSwap, swap)
	}
	require.NoError(t, err)
}

func TestGetApproveTransaction(t *testing.T) {
	ctx := context.Background()

	mockedResp := ApproveCallDataResponse{
		Data:     "0x095ea7b30000000000000000000000001111111254eeb25477b68fb85ed929f73a9605820000000000000000000000000000000000000000000000000de0b6b3a7640000",
		GasPrice: "17836383853",
		To:       "0x5a98fcbea516cf06857215779fd812ca3bef1b32",
		Value:    "0",
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
		chainId:      constants.EthereumChainId,
	}

	params := GetApproveParams{
		TokenAddress: "0x5a98fcbea516cf06857215779fd812ca3bef1b32",
		Amount:       "1000000000000000000",
	}

	approveData, err := api.GetApproveTransaction(ctx, params)
	if err != nil {
		t.Fatalf("GetApproveTransaction returned an error: %v", err)
	}

	if !mockExecutor.Called {
		t.Errorf("Expected ExecuteRequest to be called")
	}

	expectedApproveData, err := normalizeApproveCallDataResponse(mockedResp)
	if !reflect.DeepEqual(approveData, expectedApproveData) {
		t.Errorf("Expected approve data to be %+v, got %+v", expectedApproveData, approveData)
	}
	require.NoError(t, err)

}

func TestGetApproveSpender(t *testing.T) {
	ctx := context.Background()

	mockedResp := SpenderResponse{
		Address: "0x1111111254eeb25477b68fb85ed929f73a960582",
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
		chainId:      constants.EthereumChainId,
	}

	spender, err := api.GetApproveSpender(ctx)
	require.NoError(t, err)
	require.True(t, mockExecutor.Called, "ExecuteRequest should be called")
	require.NotNil(t, spender, "Spender response should not be nil")
	require.Equal(t, mockedResp.Address, spender.Address, "The returned address should match the expected address")
}

func TestGetApproveAllowance(t *testing.T) {
	ctx := context.Background()

	mockedResp := AllowanceResponse{
		Allowance: "1000000000000000000",
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := &api{
		httpExecutor: mockExecutor,
		chainId:      constants.EthereumChainId,
	}

	params := GetAllowanceParams{
		TokenAddress:  "0x5a98fcbea516cf06857215779fd812ca3bef1b32",
		WalletAddress: "0x083fc10ce7e97cafbae0fe332a9c4384c5f54e45",
	}

	allowance, err := api.GetApproveAllowance(ctx, params)
	assert.NoError(t, err)
	assert.NotNil(t, allowance)
	assert.Equal(t, mockedResp.Allowance, allowance.Allowance)
}

func TestGetTokens(t *testing.T) {
	ctx := context.Background()

	mockedResp := TokensResponse{
		Tokens: map[string]TokenInfo{
			"0x5a98fcbea516cf06857215779fd812ca3bef1b32": {
				Address:  "0x5a98fcbea516cf06857215779fd812ca3bef1b32",
				Symbol:   "LDO",
				Name:     "Lido DAO Token",
				Decimals: 18,
				LogoURI:  "https://tokens.1inch.io/0x5a98fcbea516cf06857215779fd812ca3bef1b32.png",
				Tags: []string{
					"tokens",
				},
			},
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2": {
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
		},
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := &api{
		httpExecutor: mockExecutor,
		chainId:      constants.EthereumChainId,
	}

	tokensResponse, err := api.GetTokens(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, tokensResponse)
	assert.Equal(t, mockedResp.Tokens, tokensResponse.Tokens, "The returned tokens should match the expected tokens")
}

func TestGetLiquiditySources(t *testing.T) {
	ctx := context.Background()

	mockedResp := ProtocolsResponse{
		Protocols: []ProtocolImage{
			{Id: "UNISWAP_V1", Img: "https://cdn.1inch.io/liquidity-sources-logo/uniswap.png", ImgColor: "https://cdn.1inch.io/liquidity-sources-logo/uniswap_color.png", Title: "Uniswap V1"},
			{Id: "UNISWAP_V2", Img: "https://cdn.1inch.io/liquidity-sources-logo/uniswap.png", ImgColor: "https://cdn.1inch.io/liquidity-sources-logo/uniswap_color.png", Title: "Uniswap V2"},
		},
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}
	api := &api{
		httpExecutor: mockExecutor,
		chainId:      constants.EthereumChainId,
	}

	liquiditySources, err := api.GetLiquiditySources(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, liquiditySources)
	assert.Equal(t, mockedResp.Protocols, liquiditySources.Protocols, "The returned protocols should match the expected protocols")

}
