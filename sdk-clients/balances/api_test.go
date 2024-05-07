package balances

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

func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }
func boolPtr(b bool) *bool    { return &b }

func TestGetBalancesAndAllowancesByWalletAddressList(t *testing.T) {
	ctx := context.Background()

	mockedResp := BalancesAndAllowancesByWalletAddressListResponse{
		"0xfc1e690f61efd961294b3e1ce3313fbd8aa4f85d": TokenDetails{
			Balance:   "0",
			Allowance: "0",
		},
		"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee": TokenDetails{
			Balance:   "548417674176835310649",
			Allowance: "115792089237316195423570985008687907853269984665640564039457584007913129639935",
		},
		"0x06af07097c9eeb7fd685c692751d5c66db49c215": TokenDetails{
			Balance:   "0",
			Allowance: "0",
		},
		"0xf5dce57282a584d2746faf1593d3121fcac444dc": TokenDetails{
			Balance:   "0",
			Allowance: "0",
		},
		"0x4ddc2d193948926d02f9b1fe9e1daa0718270ed5": TokenDetails{
			Balance:   "0",
			Allowance: "0",
		},
		"0x39aa39c021dfbae8fac545936693ac917d5e7563": TokenDetails{
			Balance:   "61790806808",
			Allowance: "0",
		},
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
		chainId:      constants.EthereumChainId,
	}

	params := BalancesAndAllowancesByWalletAddressListParams{
		Wallet:  "0x083fc10cE7e97CaFBaE0fE332a9c4384c5f54E45",
		Spender: "0x111111125421ca6dc452d289314280a0f8842a65",
	}

	balances, err := api.GetBalancesAndAllowancesByWalletAddressList(ctx, params)
	require.NoError(t, err)
	require.NotNil(t, balances)

	if !mockExecutor.Called {
		t.Errorf("Expected ExecuteRequest to be called")
	}
	if !reflect.DeepEqual(balances, &mockedResp) {
		t.Errorf("Expected swap to be %+v, got %+v", balances, mockedResp)
	}
	require.NoError(t, err)
}

func TestGetBalancesAndAllowances(t *testing.T) {
	ctx := context.Background()

	mockedResp := AggregatedBalancesAndAllowancesResponse{
		{
			Address:  strPtr("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"),
			Decimals: intPtr(18),
			IsCustom: boolPtr(false),
			LogoURI:  strPtr("https://tokens.1inch.io/0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee.png"),
			Name:     strPtr("Ether"),
			Symbol:   strPtr("ETH"),
			Tags:     &[]string{"native", "PEG:ETH"},
			Tracked:  boolPtr(true),
			Type:     strPtr("ethereum"),
			Wallets: &map[string]struct {
				Allowance *string `json:"allowance,omitempty"`
				Balance   *string `json:"balance,omitempty"`
			}{
				"0x1C17622cfa9B6fD2043A76DfC39A5B5a109aa708": {
					Balance:   strPtr("18920076738417670657"),
					Allowance: strPtr("115792089237316195423570985008687907853269984665640564039457584007913129639935"),
				},
				"0x28C6c06298d514Db089934071355E5743bf21d60": {
					Balance:   strPtr("45363705428251046849847"),
					Allowance: strPtr("115792089237316195423570985008687907853269984665640564039457584007913129639935"),
				},
			},
		},
		{
			Address:  strPtr("0x0d02755a5700414b26ff040e1de35d337df56218"),
			Decimals: intPtr(18),
			IsCustom: boolPtr(false),
			LogoURI:  strPtr("https://tokens.1inch.io/0x0d02755a5700414b26ff040e1de35d337df56218.png"),
			Name:     strPtr("Bend Token"),
			Symbol:   strPtr("BEND"),
			Tracked:  boolPtr(true),
			Type:     strPtr("token"),
			Wallets: &map[string]struct {
				Allowance *string `json:"allowance,omitempty"`
				Balance   *string `json:"balance,omitempty"`
			}{
				"0x1C17622cfa9B6fD2043A76DfC39A5B5a109aa708": {
					Balance:   strPtr("2000000000000000000"),
					Allowance: strPtr("0"),
				},
				"0x28C6c06298d514Db089934071355E5743bf21d60": {
					Balance:   strPtr("5000000000000000000"),
					Allowance: strPtr("0"),
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

	params := BalancesAndAllowancesParams{
		Wallets:     []string{"0x1C17622cfa9B6fD2043A76DfC39A5B5a109aa708", "0x28C6c06298d514Db089934071355E5743bf21d60"},
		FilterEmpty: true,
		Spender:     "0x58b6a8a3302369daec383334672404ee733ab239",
	}

	balances, err := api.GetBalancesAndAllowances(ctx, params)
	require.NoError(t, err)
	require.NotNil(t, balances)

	if !mockExecutor.Called {
		t.Errorf("Expected ExecuteRequest to be called")
	}
	if !reflect.DeepEqual(balances, &mockedResp) {
		t.Errorf("Expected swap to be %+v, got %+v", balances, mockedResp)
	}
	require.NoError(t, err)
}
