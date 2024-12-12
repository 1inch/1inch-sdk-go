package fusion

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//func TestSettlementPostInteractionDataDecode(t *testing.T) {
//	tests := []struct {
//		name   string
//		data   string
//		expect SettlementPostInteractionData
//	}{
//		{
//			name: "Should decode",
//			data: "6656b877b09498030ae3416b66dc0000db05a6a504f04d92e79d00000c989d73cf0bd5f83b660000d18bd45f0b94f54a968f0000d61b892b2ad6249011850000d0847e80c0b823a65ce70000901f8f650d76dcc657d1000038ad1723a873d05effcbdc57dcf7d00458d6a8c763558d5af7522bf6ad2d3e253d000000000000000000000000000000000000000000000000000000000000a4b1000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000064000000000000000000000000000000c80000000000000003000000020000000100000004000000030000000200000001",
//			expect: SettlementPostInteractionData{
//				Whitelist: []WhitelistItem{
//					{
//						AddressHalf: "b09498030ae3416b66dc",
//						Delay:       big.NewInt(0),
//					},
//					{
//						AddressHalf: "db05a6a504f04d92e79d",
//						Delay:       big.NewInt(0),
//					},
//					{
//						AddressHalf: "0c989d73cf0bd5f83b66",
//						Delay:       big.NewInt(0),
//					},
//					{
//						AddressHalf: "d18bd45f0b94f54a968f",
//						Delay:       big.NewInt(0),
//					},
//					{
//						AddressHalf: "d61b892b2ad624901185",
//						Delay:       big.NewInt(0),
//					},
//					{
//						AddressHalf: "d0847e80c0b823a65ce7",
//						Delay:       big.NewInt(0),
//					},
//					{
//						AddressHalf: "901f8f650d76dcc657d1",
//						Delay:       big.NewInt(0),
//					},
//				},
//				BankFee:            nil,
//				ResolvingStartTime: big.NewInt(1716959351),
//			},
//		},
//	}
//
//	for _, tc := range tests {
//		t.Run(tc.name, func(t *testing.T) {
//			data, err := Decode(tc.data)
//			require.NoError(t, err)
//			assert.Equal(t, tc.expect, data)
//		})
//	}
//
//}

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

			encoded := data.Encode()
			if tc.expectedBytes != 0 {
				assert.Equal(t, tc.expectedBytes, len(encoded[2:])/2)
			}

			decoded, err := Decode(encoded)
			assert.NoError(t, err)
			assert.Equal(t, *data, decoded)
		})
	}
}
