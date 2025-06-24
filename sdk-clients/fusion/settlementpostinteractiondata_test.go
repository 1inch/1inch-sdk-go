package fusion

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeSettlementPostInteractionData(t *testing.T) {
	spid := &SettlementPostInteractionData{
		Whitelist: []WhitelistItem{
			{AddressHalf: "bb839cbe05303d7705fa", Delay: big.NewInt(0)},
		},
		IntegratorFee: &IntegratorFee{
			Ratio:    nil,
			Receiver: common.HexToAddress("0x0000000000000000000000000000000000000000"),
		},
		BankFee:            big.NewInt(10),
		ResolvingStartTime: big.NewInt(11),
		CustomReceiver:     common.HexToAddress("0x0000000000000000000000000000000000000000"),
		AuctionFees: &FeesNew{
			Resolver:   ResolverFee{},
			Integrator: IntegratorFeeNew{},
		},
	}

	details := &AuctionDetails{
		GasCost: GasCostConfigClassFixed{
			GasBumpEstimate:  0,
			GasPriceEstimate: 0,
		},
		StartTime:       1673548149,
		Duration:        180,
		InitialRateBump: 50000,
		Points: []AuctionPointClassFixed{
			{
				Delay:       12,
				Coefficient: 20000,
			},
		},
	}

	surplus := &SurplusParams{
		EstimatedTakerAmount: Uint256Max,
		ProtocolFee:          BpsZero,
	}

	ext := &Extension{
		SettlementContract:  "0x8273f37417da37c4a6c3995e82cf442f87a25d9c",
		AuctionDetails:      details,
		Surplus:             surplus,
		ResolvingStartTime:  big.NewInt(1673548139),
		PostInteractionData: spid,
	}

	data, err := CreateEncodedPostInteractionData(ext)
	require.NoError(t, err)
	require.Equal(t, "0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006463c0516b01bb839cbe05303d7705fa0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00", data)
}

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
			assert.NoError(t, err)
			assert.Equal(t, *data, decoded)
		})
	}
}

//func TestGenerateWhitelist(t *testing.T) {
//	tests := []struct {
//		name     string
//		data     *SettlementSuffixData
//		expected []WhitelistItem
//	}{
//		{
//			name: "Should generate whitelist",
//			data: &SettlementSuffixData{
//				Whitelist:          []AuctionWhitelistItem{{Address: common.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fa"), AllowFrom: big.NewInt(0)}},
//				ResolvingStartTime: big.NewInt(1708117482),
//			},
//			expected: []WhitelistItem{
//				{
//					AddressHalf: "bb839cbe05303d7705fa",
//					Delay:       big.NewInt(0),
//				},
//			},
//		},
//	}
//
//	for _, tc := range tests {
//		t.Run(tc.name, func(t *testing.T) {
//			whitelist, err := GenerateWhitelist(tc.data)
//			require.NoError(t, err)
//			assert.Equal(t, tc.expected, whitelist)
//		})
//	}
//}

func TestGenerateWhitelist(t *testing.T) {
	tests := []struct {
		name               string
		whitelistStrings   []string
		resolvingStartTime *big.Int
		expected           []WhitelistItem
	}{
		{
			name:               "Should generate whitelist",
			whitelistStrings:   []string{"0x00000000219ab540356cbb839cbe05303d7705fa"},
			resolvingStartTime: big.NewInt(1708117482),
			expected: []WhitelistItem{
				{
					AddressHalf: "bb839cbe05303d7705fa",
					Delay:       big.NewInt(0),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			whitelist, err := GenerateWhitelist(tc.whitelistStrings, tc.resolvingStartTime)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, whitelist)
		})
	}
}
