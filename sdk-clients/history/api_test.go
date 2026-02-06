package history

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

func TestGetHistoryEventsByAddress(t *testing.T) {
	ctx := context.Background()

	mockedResp := EventsByAddressResponse{
		Items: []Item{
			{
				TimeMs:  1696151195140,
				Address: "0x266e77ce9034a023056ea2845cb6a20517f6fdb7",
				Type:    0,
				Rating:  "reliable",
				Details: Details{
					TxHash:       "0x00b61f593833012bde2bcb6b209a79fe9de47cc28f102509145fbbca4e1c6566",
					ChainID:      1,
					BlockNumber:  18254685,
					BlockTimeSec: 1696151195,
					Status:       "completed",
					Type:         "Send",
					TokenActions: []TokenAction{
						{
							Address:     "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
							Standard:    "Native",
							FromAddress: "0x266e77ce9034a023056ea2845cb6a20517f6fdb7",
							ToAddress:   "0xae0ee0a63a2ce6baeeffe56e7714fb4efe48d419",
							Amount:      "4935656027779822",
							Direction:   "Out",
						},
					},
					FromAddress:  "0x266e77ce9034a023056ea2845cb6a20517f6fdb7",
					ToAddress:    "0xae0ee0a63a2ce6baeeffe56e7714fb4efe48d419",
					OrderInBlock: 140,
					Nonce:        9,
					FeeInWei:     "750577247363861",
				},
				ID:                      "306262793328640",
				EventOrderInTransaction: 0,
			},
		},
		CacheCounter: 1,
	}

	mockExecutor := &MockHttpExecutor{
		ResponseObj: mockedResp,
	}

	api := api{
		httpExecutor: mockExecutor,
	}

	params := EventsByAddressParams{
		Address: "0x266E77cE9034a023056ea2845CB6A20517F6FDB7",
		ChainId: constants.EthereumChainId,
	}

	prices, err := api.GetHistoryEventsByAddress(ctx, params)
	require.NoError(t, err)
	require.NotNil(t, prices)

	if !mockExecutor.Called {
		t.Errorf("Expected ExecuteRequest to be called")
	}
	if !reflect.DeepEqual(prices, &mockedResp) {
		t.Errorf("Expected swap to be %+v, got %+v", prices, mockedResp)
	}
	require.NoError(t, err)
}
