package fusion

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestSettlementPostInteractionData(t *testing.T) {
	tests := []struct {
		name          string
		data          SettlementSuffixData
		expectedBytes int
	}{
		{
			name: "Should encode/decode with bank fee and whitelist",
			data: SettlementSuffixData{
				BankFee:            big.NewInt(1),
				ResolvingStartTime: big.NewInt(1708117482),
				Whitelist: []AuctionWhitelistItem{
					{
						Address:   common.Address{},
						AllowFrom: big.NewInt(0),
					},
				},
			},
			expectedBytes: 21,
		},
		{
			name: "Should encode/decode with no fees and whitelist",
			data: SettlementSuffixData{
				ResolvingStartTime: big.NewInt(1708117482),
				Whitelist: []AuctionWhitelistItem{
					{
						Address:   common.Address{},
						AllowFrom: big.NewInt(0),
					},
				},
			},
			expectedBytes: 17,
		},
		{
			name: "Should encode/decode with fees and whitelist",
			data: SettlementSuffixData{
				BankFee:            big.NewInt(0),
				ResolvingStartTime: big.NewInt(1708117482),
				Whitelist: []AuctionWhitelistItem{
					{
						Address:   common.Address{},
						AllowFrom: big.NewInt(0),
					},
				},
				IntegratorFee: &IntegratorFee{
					Receiver: common.Address{1},
					Ratio:    big.NewInt(10),
				},
			},
			expectedBytes: 21, // This value might change based on your actual implementation
		},
		{
			name: "Should encode/decode with fees, custom receiver and whitelist",
			data: SettlementSuffixData{
				BankFee:            big.NewInt(0),
				ResolvingStartTime: big.NewInt(1708117482),
				Whitelist: []AuctionWhitelistItem{
					{
						Address:   common.Address{},
						AllowFrom: big.NewInt(0),
					},
				},
				IntegratorFee: &IntegratorFee{
					Receiver: common.Address{1},
					Ratio:    big.NewInt(10),
				},
				CustomReceiver: common.Address{123},
			},
			expectedBytes: 25, // This value might change based on your actual implementation
		},
		{
			name: "Should generate correct whitelist",
			data: SettlementSuffixData{
				ResolvingStartTime: big.NewInt(1708117482),
				Whitelist: []AuctionWhitelistItem{
					{
						Address:   common.Address{2},
						AllowFrom: big.NewInt(1708117482 + 1000),
					},
					{
						Address:   common.Address{},
						AllowFrom: big.NewInt(1708117482 - 10),
					},
					{
						Address:   common.Address{1},
						AllowFrom: big.NewInt(1708117482 + 10),
					},
					{
						Address:   common.Address{3},
						AllowFrom: big.NewInt(1708117482 + 10),
					},
				},
			},
			expectedBytes: 0, // Not used in this case
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data := NewSettlementPostInteractionData(tc.data)
			encoded := data.Encode()

			if tc.expectedBytes != 0 {
				assert.Equal(t, tc.expectedBytes, len(encoded)/2)
				decoded, err := DecodeSettlementPostInteractionData(encoded)
				assert.NoError(t, err)
				assert.Equal(t, data, decoded)
			}

			if tc.name == "Should generate correct whitelist" {
				expectedWhitelist := []WhitelistItem{
					{AddressHalf: "", Delay: big.NewInt(0)},
					{AddressHalf: "00000000000000000001", Delay: big.NewInt(10)},
					{AddressHalf: "00000000000000000003", Delay: big.NewInt(0)},
					{AddressHalf: "00000000000000000002", Delay: big.NewInt(990)},
				}

				assert.Equal(t, expectedWhitelist, data.Whitelist)

				start := big.NewInt(1708117482)
				assert.True(t, data.CanExecuteAt(common.Address{1}, new(big.Int).Add(start, big.NewInt(10))))
				assert.False(t, data.CanExecuteAt(common.Address{1}, new(big.Int).Add(start, big.NewInt(9))))
				assert.True(t, data.CanExecuteAt(common.Address{3}, new(big.Int).Add(start, big.NewInt(10))))
				assert.False(t, data.CanExecuteAt(common.Address{3}, new(big.Int).Add(start, big.NewInt(9))))
				assert.False(t, data.CanExecuteAt(common.Address{2}, new(big.Int).Add(start, big.NewInt(50))))
			}
		})
	}
}
