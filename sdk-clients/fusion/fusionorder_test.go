package fusion

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

func TestCreateMakerTraits(t *testing.T) {
	tests := []struct {
		name        string
		details     Details
		extraParams ExtraParams
		expected    *orderbook.MakerTraits
		expectErr   bool
		expectedErr error
	}{
		{
			name: "Valid Maker Traits",
			details: Details{
				Auction: &AuctionDetails{
					StartTime: 1000,
					Duration:  2000,
				},
				Fees: Fees{
					IntFee: IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.Address{},
					},
					BankFee: big.NewInt(200),
				},
			},
			extraParams: ExtraParams{
				Nonce:                big.NewInt(1),
				Permit:               "permit",
				AllowPartialFills:    true,
				AllowMultipleFills:   true,
				OrderExpirationDelay: 3000,
				EnablePermit2:        true,
				Source:               "source",
				unwrapWeth:           true,
			},
			expected: &orderbook.MakerTraits{
				AllowedSender:       "",
				Expiry:              6000,
				Nonce:               1,
				Series:              0,
				NoPartialFills:      false,
				NeedPostinteraction: true,
				NeedPreinteraction:  false,
				NeedEpochCheck:      false,
				HasExtension:        true,
				ShouldUsePermit2:    true,
				ShouldUnwrapWeth:    true,
				AllowPartialFills:   true,
				AllowMultipleFills:  true,
			},
			expectErr: false,
		},
		{
			name: "Invalid Maker Traits - No Nonce",
			details: Details{
				Auction: &AuctionDetails{
					StartTime: 1000,
					Duration:  2000,
				},
				Fees: Fees{
					IntFee: IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.Address{},
					},
					BankFee: big.NewInt(200),
				},
			},
			extraParams: ExtraParams{
				Nonce:                big.NewInt(0),
				Permit:               "permit",
				AllowPartialFills:    false,
				AllowMultipleFills:   false,
				OrderExpirationDelay: 3000,
				EnablePermit2:        true,
				Source:               "source",
				unwrapWeth:           true,
			},
			expected:    nil,
			expectErr:   true,
			expectedErr: errors.New("nonce required, when partial fill or multiple fill disallowed"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CreateMakerTraits(tc.details, tc.extraParams)
			if tc.expectErr {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestCreateSettlementPostInteractionData(t *testing.T) {
	tests := []struct {
		name        string
		details     Details
		orderInfo   FusionOrderV4
		expected    *SettlementPostInteractionData
		expectErr   bool
		expectedErr error
	}{
		{
			name: "Valid Details and Order Info with Resolving Start Time",
			details: Details{
				ResolvingStartTime: big.NewInt(1622548800), // Example timestamp
				Fees: Fees{
					IntFee: IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
					},
					BankFee: big.NewInt(200),
				},
				Whitelist: []AuctionWhitelistItem{
					{
						Address:   common.HexToAddress("0x0000000000000000000000000000000000000002"),
						AllowFrom: big.NewInt(1622548800),
					},
				},
			},
			orderInfo: FusionOrderV4{
				Receiver: "0x0000000000000000000000000000000000000003",
			},
			expected: &SettlementPostInteractionData{
				Whitelist: []WhitelistItem{
					{
						AddressHalf: "00000000000000000002",
						Delay:       big.NewInt(0),
					},
				},
				IntegratorFee: &IntegratorFee{
					Ratio:    big.NewInt(100),
					Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
				},
				BankFee:            big.NewInt(200),
				ResolvingStartTime: big.NewInt(1622548800),
				CustomReceiver:     common.HexToAddress("0x0000000000000000000000000000000000000003"),
			},
			expectErr: false,
		},
		{
			name: "Valid Details and Order Info with non-zero Delay",
			details: Details{
				ResolvingStartTime: big.NewInt(1622548800), // Example timestamp
				Fees: Fees{
					IntFee: IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
					},
					BankFee: big.NewInt(200),
				},
				Whitelist: []AuctionWhitelistItem{
					{
						Address:   common.HexToAddress("0x0000000000000000000000000000000000000002"),
						AllowFrom: big.NewInt(1622549800),
					},
				},
			},
			orderInfo: FusionOrderV4{
				Receiver: "0x0000000000000000000000000000000000000003",
			},
			expected: &SettlementPostInteractionData{
				Whitelist: []WhitelistItem{
					{
						AddressHalf: "00000000000000000002",
						Delay:       big.NewInt(1000),
					},
				},
				IntegratorFee: &IntegratorFee{
					Ratio:    big.NewInt(100),
					Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
				},
				BankFee:            big.NewInt(200),
				ResolvingStartTime: big.NewInt(1622548800),
				CustomReceiver:     common.HexToAddress("0x0000000000000000000000000000000000000003"),
			},
			expectErr: false,
		},
		{
			name: "Valid Details and Order Info without Resolving Start Time",
			details: Details{
				ResolvingStartTime: nil,
				Fees: Fees{
					IntFee: IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
					},
					BankFee: big.NewInt(200),
				},
				Whitelist: []AuctionWhitelistItem{
					{
						Address:   common.HexToAddress("0x0000000000000000000000000000000000000002"),
						AllowFrom: big.NewInt(1622548800),
					},
				},
			},
			orderInfo: FusionOrderV4{
				Receiver: "0x0000000000000000000000000000000000000003",
			},
			expected: &SettlementPostInteractionData{
				Whitelist: []WhitelistItem{
					{
						AddressHalf: "00000000000000000002",
						Delay:       big.NewInt(0),
					},
				},
				IntegratorFee: &IntegratorFee{
					Ratio:    big.NewInt(100),
					Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
				},
				BankFee:            big.NewInt(200),
				ResolvingStartTime: big.NewInt(timeNow()), // This will be dynamically set
				CustomReceiver:     common.HexToAddress("0x0000000000000000000000000000000000000003"),
			},
			expectErr: false,
		},
		{
			name: "Delay too large",
			details: Details{
				ResolvingStartTime: big.NewInt(1622548800), // Example timestamp
				Fees: Fees{
					IntFee: IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
					},
					BankFee: big.NewInt(200),
				},
				Whitelist: []AuctionWhitelistItem{
					{
						Address:   common.HexToAddress("0x0000000000000000000000000000000000000002"),
						AllowFrom: big.NewInt(1622649800),
					},
				},
			},
			orderInfo: FusionOrderV4{
				Receiver: "0x0000000000000000000000000000000000000003",
			},
			expected: &SettlementPostInteractionData{
				Whitelist: []WhitelistItem{
					{
						AddressHalf: "00000000000000000002",
						Delay:       big.NewInt(1000),
					},
				},
				IntegratorFee: &IntegratorFee{
					Ratio:    big.NewInt(100),
					Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
				},
				BankFee:            big.NewInt(200),
				ResolvingStartTime: big.NewInt(1622548800),
				CustomReceiver:     common.HexToAddress("0x0000000000000000000000000000000000000003"),
			},
			expectErr:   true,
			expectedErr: fmt.Errorf("delay too big - %d must be less than %d", 101000, uint16Max),
		},
		{
			name: "Whitelist empty",
			details: Details{
				ResolvingStartTime: nil,
				Fees: Fees{
					IntFee: IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
					},
					BankFee: big.NewInt(200),
				},
				Whitelist: []AuctionWhitelistItem{},
			},
			orderInfo: FusionOrderV4{
				Receiver: "0x0000000000000000000000000000000000000003",
			},
			expectErr:   true,
			expectedErr: errors.New("whitelist cannot be empty"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CreateSettlementPostInteractionData(tc.details, tc.orderInfo)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				// Setting the dynamic field to the expected result for comparison
				if tc.details.ResolvingStartTime == nil {
					tc.expected.ResolvingStartTime = result.ResolvingStartTime
				}
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestCreateExtension(t *testing.T) {
	tests := []struct {
		name      string
		params    CreateExtensionParams
		expected  *Extension
		expectErr bool
	}{
		{
			name: "Valid Parameters with Permit",
			params: CreateExtensionParams{
				settlementAddress: "0x0000000000000000000000000000000000000001",
				postInteractionData: &SettlementPostInteractionData{
					Whitelist: []WhitelistItem{
						{
							AddressHalf: "abcdef",
							Delay:       big.NewInt(1000),
						},
					},
					IntegratorFee: &IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000002"),
					},
					BankFee:            big.NewInt(200),
					ResolvingStartTime: big.NewInt(1622548800),
					CustomReceiver:     common.HexToAddress("0x0000000000000000000000000000000000000003"),
				},
				orderInfo: FusionOrderV4{
					MakerAsset: "0x0000000000000000000000000000000000000004",
					Receiver:   "0x0000000000000000000000000000000000000005",
				},
				details: Details{
					Auction: &AuctionDetails{
						StartTime: 1000,
						Duration:  2000,
					},
				},
				extraParams: ExtraParams{
					Permit: "0xabcdef",
				},
			},
			expected: &Extension{
				MakingAmountData: "0x0000000000000000000000000000000000000001" + "00000000000000000003e80007d0000000",
				TakingAmountData: "0x0000000000000000000000000000000000000001" + "00000000000000000003e80007d0000000",
				PostInteraction:  "0x0000000000000000000000000000000000000001000000c800640000000000000000000000000000000000000002000000000000000000000000000000000000000360b62140abcdef03e80f",
				MakerPermit:      "0x0000000000000000000000000000000000000004abcdef",
			},
			expectErr: false,
		},
		{
			name: "Valid Parameters without Permit",
			params: CreateExtensionParams{
				settlementAddress: "0x0000000000000000000000000000000000000001",
				postInteractionData: &SettlementPostInteractionData{
					Whitelist: []WhitelistItem{
						{
							AddressHalf: "abcdef",
							Delay:       big.NewInt(1000),
						},
					},
					IntegratorFee: &IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000002"),
					},
					BankFee:            big.NewInt(200),
					ResolvingStartTime: big.NewInt(1622548800),
					CustomReceiver:     common.HexToAddress("0x0000000000000000000000000000000000000003"),
				},
				orderInfo: FusionOrderV4{
					MakerAsset: "0x0000000000000000000000000000000000000004",
					Receiver:   "0x0000000000000000000000000000000000000005",
				},
				details: Details{
					Auction: &AuctionDetails{
						StartTime: 1000,
						Duration:  2000,
					},
				},
				extraParams: ExtraParams{},
			},
			expected: &Extension{
				MakingAmountData: "0x0000000000000000000000000000000000000001" + "00000000000000000003e80007d0000000",
				TakingAmountData: "0x0000000000000000000000000000000000000001" + "00000000000000000003e80007d0000000",
				PostInteraction:  "0x0000000000000000000000000000000000000001000000c800640000000000000000000000000000000000000002000000000000000000000000000000000000000360b62140abcdef03e80f",
			},
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CreateExtension(tc.params)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestCreateOrder(t *testing.T) {
	tests := []struct {
		name       string
		staticSalt string
		params     CreateOrderDataParams
		expected   *Order
		expectErr  bool
	}{
		{
			name:       "Valid Order with Integrator Fee",
			staticSalt: "180431658011416401710119735245975317914670388782711199",
			params: CreateOrderDataParams{
				settlementAddress: "0x0000000000000000000000000000000000000001",
				postInteractionData: &SettlementPostInteractionData{
					IntegratorFee: &IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000002"),
					},
					BankFee: big.NewInt(200),
				},
				extension: &Extension{
					MakerAssetSuffix: "suffix1",
					TakerAssetSuffix: "suffix2",
					MakingAmountData: "data1",
					TakingAmountData: "data2",
					Predicate:        "predicate",
					MakerPermit:      "permit",
					PreInteraction:   "pre",
					PostInteraction:  "post",
					CustomData:       "custom",
				},
				orderInfo: FusionOrderV4{
					Maker:        "0x0000000000000000000000000000000000000003",
					MakerAsset:   "0x0000000000000000000000000000000000000004",
					TakerAsset:   "0x0000000000000000000000000000000000000005",
					MakingAmount: "1000",
					TakingAmount: "2000",
					Receiver:     "0x0000000000000000000000000000000000000006",
				},
				details: Details{
					Auction: &AuctionDetails{
						StartTime: 1000,
						Duration:  2000,
					},
				},
				extraParams: ExtraParams{
					Nonce: big.NewInt(1),
				},
				makerTraits: &orderbook.MakerTraits{
					AllowedSender:      "0x0000000000000000000000000000000000000007",
					Expiry:             5000,
					Nonce:              1,
					AllowPartialFills:  true,
					AllowMultipleFills: true,
				},
			},
			expected: &Order{
				FusionExtension: &Extension{
					MakerAssetSuffix: "suffix1",
					TakerAssetSuffix: "suffix2",
					MakingAmountData: "data1",
					TakingAmountData: "data2",
					Predicate:        "predicate",
					MakerPermit:      "permit",
					PreInteraction:   "pre",
					PostInteraction:  "post",
					CustomData:       "custom",
				},
				Inner: orderbook.OrderData{
					Maker:        "0x0000000000000000000000000000000000000003",
					MakerAsset:   "0x0000000000000000000000000000000000000004",
					TakerAsset:   "0x0000000000000000000000000000000000000005",
					MakingAmount: "1000",
					TakingAmount: "2000",
					Salt:         "1e24059a92eed490b75f51b98fbf2e143c8e9712a419f59a92eed490b75f51b98fbf2e143c8e9712a419f",
					MakerTraits:  "0x4000000000000000000000000000000001000000138800000000000000000007",
					Receiver:     "0x0000000000000000000000000000000000000001", // Address of settlementAddress because Integrator Fee is set
					Extension:    "343845d3ef4b5505456e95d059a92eed490b75f51b98fbf2e143c8e9712a419f",
				},
				SettlementExtension: common.HexToAddress("0x0000000000000000000000000000000000000001"),
				OrderInfo: FusionOrderV4{
					Maker:        "0x0000000000000000000000000000000000000003",
					MakerAsset:   "0x0000000000000000000000000000000000000004",
					TakerAsset:   "0x0000000000000000000000000000000000000005",
					MakingAmount: "1000",
					TakingAmount: "2000",
					Receiver:     "0x0000000000000000000000000000000000000006",
				},
				AuctionDetails: &AuctionDetails{
					StartTime: 1000,
					Duration:  2000,
				},
				PostInteractionData: &SettlementPostInteractionData{
					IntegratorFee: &IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000002"),
					},
					BankFee: big.NewInt(200),
				},
				Extra: ExtraData{
					UnwrapWETH:           false,
					Nonce:                big.NewInt(1),
					Permit:               "",
					AllowPartialFills:    false,
					AllowMultipleFills:   false,
					OrderExpirationDelay: 0,
					EnablePermit2:        false,
					Source:               "",
				},
			},
			expectErr: false,
		},
		{
			name:       "Valid Order without Integrator Fee",
			staticSalt: "180431658011416401710119735245975317914670388782711199",
			params: CreateOrderDataParams{
				settlementAddress: "0x0000000000000000000000000000000000000001",
				postInteractionData: &SettlementPostInteractionData{
					IntegratorFee: &IntegratorFee{
						Ratio:    big.NewInt(0),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000002"),
					},
					BankFee: big.NewInt(200),
				},
				extension: &Extension{
					MakerAssetSuffix: "suffix1",
					TakerAssetSuffix: "suffix2",
					MakingAmountData: "data1",
					TakingAmountData: "data2",
					Predicate:        "predicate",
					MakerPermit:      "permit",
					PreInteraction:   "pre",
					PostInteraction:  "post",
					CustomData:       "custom",
				},
				orderInfo: FusionOrderV4{
					Maker:        "0x0000000000000000000000000000000000000003",
					MakerAsset:   "0x0000000000000000000000000000000000000004",
					TakerAsset:   "0x0000000000000000000000000000000000000005",
					MakingAmount: "1000",
					TakingAmount: "2000",
					Receiver:     "0x0000000000000000000000000000000000000006",
				},
				details: Details{
					Auction: &AuctionDetails{
						StartTime: 1000,
						Duration:  2000,
					},
				},
				extraParams: ExtraParams{
					Nonce: big.NewInt(1),
				},
				makerTraits: &orderbook.MakerTraits{
					AllowedSender:      "0x0000000000000000000000000000000000000007",
					Expiry:             5000,
					Nonce:              1,
					AllowPartialFills:  true,
					AllowMultipleFills: true,
				},
			},
			expected: &Order{
				FusionExtension: &Extension{
					MakerAssetSuffix: "suffix1",
					TakerAssetSuffix: "suffix2",
					MakingAmountData: "data1",
					TakingAmountData: "data2",
					Predicate:        "predicate",
					MakerPermit:      "permit",
					PreInteraction:   "pre",
					PostInteraction:  "post",
					CustomData:       "custom",
				},
				Inner: orderbook.OrderData{
					Maker:        "0x0000000000000000000000000000000000000003",
					MakerAsset:   "0x0000000000000000000000000000000000000004",
					TakerAsset:   "0x0000000000000000000000000000000000000005",
					MakingAmount: "1000",
					TakingAmount: "2000",
					Salt:         "1e24059a92eed490b75f51b98fbf2e143c8e9712a419f59a92eed490b75f51b98fbf2e143c8e9712a419f",
					MakerTraits:  "0x4000000000000000000000000000000001000000138800000000000000000007",
					Receiver:     "0x0000000000000000000000000000000000000006", // Address of orderInfo.Receiver because Integrator Fee is not set
					Extension:    "343845d3ef4b5505456e95d059a92eed490b75f51b98fbf2e143c8e9712a419f",
				},
				SettlementExtension: common.HexToAddress("0x0000000000000000000000000000000000000001"),
				OrderInfo: FusionOrderV4{
					Maker:        "0x0000000000000000000000000000000000000003",
					MakerAsset:   "0x0000000000000000000000000000000000000004",
					TakerAsset:   "0x0000000000000000000000000000000000000005",
					MakingAmount: "1000",
					TakingAmount: "2000",
					Receiver:     "0x0000000000000000000000000000000000000006",
				},
				AuctionDetails: &AuctionDetails{
					StartTime: 1000,
					Duration:  2000,
				},
				PostInteractionData: &SettlementPostInteractionData{
					IntegratorFee: &IntegratorFee{
						Ratio:    big.NewInt(0),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000002"),
					},
					BankFee: big.NewInt(200),
				},
				Extra: ExtraData{
					UnwrapWETH:           false,
					Nonce:                big.NewInt(1),
					Permit:               "",
					AllowPartialFills:    false,
					AllowMultipleFills:   false,
					OrderExpirationDelay: 0,
					EnablePermit2:        false,
					Source:               "",
				},
			},
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			originalRandBigIntFunc := random_number_generation.BigIntMaxFunc

			staticSalt, err := BigIntFromString(tc.staticSalt)
			require.NoError(t, err)
			random_number_generation.BigIntMaxFunc = func(b *big.Int) (*big.Int, error) {
				return staticSalt, nil
			}
			result, err := CreateOrder(tc.params)
			random_number_generation.BigIntMaxFunc = originalRandBigIntFunc
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}
