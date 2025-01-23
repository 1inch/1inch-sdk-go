package fusionplus

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			name: "Should encode/decode with bank fee and whitelist with multiple entries",
			data: SettlementSuffixData{
				BankFee:            big.NewInt(1),
				ResolvingStartTime: big.NewInt(1708117482),
				Whitelist: []AuctionWhitelistItem{
					{
						Address:   common.HexToAddress("0x7a28c1b1478581b9e1293fc1c20449e2ed3efec9"),
						AllowFrom: big.NewInt(1),
					},
					{
						Address:   common.HexToAddress("0x7a28c1b1478581b9e1293fc1c20449e2ed3efec9"),
						AllowFrom: big.NewInt(2),
					},
				},
			},
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
		},
		{
			name: "Should encode/decode with fees, custom receiver and whitelist",
			data: SettlementSuffixData{
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
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := NewSettlementPostInteractionData(tc.data)
			require.NoError(t, err)

			encoded, err := data.Encode()
			require.NoError(t, err)
			if tc.expectedBytes != 0 {
				assert.Equal(t, tc.expectedBytes, len(encoded[2:])/2)
			}
			decoded, err := Decode(encoded)

			if decoded.BankFee != nil {
				decoded.BankFee = big.NewInt(decoded.BankFee.Int64())
			}
			if data.BankFee != nil {
				data.BankFee = big.NewInt(data.BankFee.Int64())
			}

			assert.NoError(t, err)
			assert.Equal(t, data.ResolvingStartTime.Cmp(decoded.ResolvingStartTime), 0)
			assert.Equal(t, data.BankFee.Cmp(decoded.BankFee), 0)

			for i, expectedItem := range data.Whitelist {
				assert.Equal(t, expectedItem.AddressHalf, decoded.Whitelist[i].AddressHalf)
				assert.Equal(t, expectedItem.Delay.Cmp(decoded.Whitelist[i].Delay), 0)
			}

			assert.Equal(t, data.IntegratorFee, decoded.IntegratorFee)
			assert.Equal(t, data.CustomReceiver, decoded.CustomReceiver)
			assert.Equal(t, data.IntegratorFee, decoded.IntegratorFee)
			assert.Equal(t, data.CustomReceiver, decoded.CustomReceiver)
		})
	}
}
