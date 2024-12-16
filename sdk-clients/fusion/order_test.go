package fusion

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

func TestGetPreset(t *testing.T) {
	customPreset := &PresetClass{
		AllowMultipleFills: true,
		AllowPartialFills:  true,
		AuctionDuration:    10.0,
		AuctionEndAmount:   "1000",
		AuctionStartAmount: "500",
		BankFee:            "5",
		EstP:               0.1,
		ExclusiveResolver:  map[string]interface{}{"resolver": "value"},
		GasCost: GasCostConfigClass{
			GasBumpEstimate:  1.0,
			GasPriceEstimate: "100",
		},
		InitialRateBump: 0.2,
		Points: []AuctionPointClass{
			{Coefficient: 1.0, Delay: 2.0},
		},
		StartAuctionIn: 1.0,
		TokenFee:       "1",
	}

	fastPreset := PresetClass{
		AllowMultipleFills: false,
		AllowPartialFills:  false,
		AuctionDuration:    20.0,
		AuctionEndAmount:   "2000",
		AuctionStartAmount: "1000",
		BankFee:            "10",
		EstP:               0.2,
		ExclusiveResolver:  map[string]interface{}{"resolver": "value"},
		GasCost: GasCostConfigClass{
			GasBumpEstimate:  2.0,
			GasPriceEstimate: "200",
		},
		InitialRateBump: 0.4,
		Points: []AuctionPointClass{
			{Coefficient: 2.0, Delay: 4.0},
		},
		StartAuctionIn: 2.0,
		TokenFee:       "2",
	}

	mediumPreset := PresetClass{
		AllowMultipleFills: true,
		AllowPartialFills:  false,
		AuctionDuration:    30.0,
		AuctionEndAmount:   "3000",
		AuctionStartAmount: "1500",
		BankFee:            "15",
		EstP:               0.3,
		ExclusiveResolver:  map[string]interface{}{"resolver": "value"},
		GasCost: GasCostConfigClass{
			GasBumpEstimate:  3.0,
			GasPriceEstimate: "300",
		},
		InitialRateBump: 0.6,
		Points: []AuctionPointClass{
			{Coefficient: 3.0, Delay: 6.0},
		},
		StartAuctionIn: 3.0,
		TokenFee:       "3",
	}

	slowPreset := PresetClass{
		AllowMultipleFills: false,
		AllowPartialFills:  true,
		AuctionDuration:    40.0,
		AuctionEndAmount:   "4000",
		AuctionStartAmount: "2000",
		BankFee:            "20",
		EstP:               0.4,
		ExclusiveResolver:  map[string]interface{}{"resolver": "value"},
		GasCost: GasCostConfigClass{
			GasBumpEstimate:  4.0,
			GasPriceEstimate: "400",
		},
		InitialRateBump: 0.8,
		Points: []AuctionPointClass{
			{Coefficient: 4.0, Delay: 8.0},
		},
		StartAuctionIn: 4.0,
		TokenFee:       "4",
	}

	presets := QuotePresetsClass{
		Custom: customPreset,
		Fast:   fastPreset,
		Medium: mediumPreset,
		Slow:   slowPreset,
	}

	tests := []struct {
		name       string
		presetType GetQuoteOutputRecommendedPreset
		expected   *PresetClass
		expectErr  bool
	}{
		{
			name:       "Get Custom Preset",
			presetType: Custom,
			expected:   customPreset,
			expectErr:  false,
		},
		{
			name:       "Get Fast Preset",
			presetType: Fast,
			expected:   &fastPreset,
			expectErr:  false,
		},
		{
			name:       "Get Medium Preset",
			presetType: Medium,
			expected:   &mediumPreset,
			expectErr:  false,
		},
		{
			name:       "Get Slow Preset",
			presetType: Slow,
			expected:   &slowPreset,
			expectErr:  false,
		},
		{
			name:       "Unknown Preset Type",
			presetType: GetQuoteOutputRecommendedPreset("Unknown"),
			expected:   nil,
			expectErr:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := getPreset(presets, tc.presetType)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestCreateAuctionDetails(t *testing.T) {
	tests := []struct {
		name                 string
		preset               *PresetClass
		additionalWaitPeriod float32
		expected             *AuctionDetails
		expectErr            bool
	}{
		{
			name: "Valid Preset",
			preset: &PresetClass{
				AllowMultipleFills: true,
				AllowPartialFills:  true,
				AuctionDuration:    60.0,
				AuctionEndAmount:   "1000",
				AuctionStartAmount: "500",
				BankFee:            "5",
				EstP:               0.1,
				ExclusiveResolver:  map[string]interface{}{"resolver": "value"},
				GasCost: GasCostConfigClass{
					GasBumpEstimate:  1.0,
					GasPriceEstimate: "100",
				},
				InitialRateBump: 2,
				Points: []AuctionPointClass{
					{Coefficient: 1.0, Delay: 2.0},
				},
				StartAuctionIn: 5.0,
				TokenFee:       "1",
			},
			additionalWaitPeriod: 10.0,
			expected: &AuctionDetails{
				StartTime:       CalcAuctionStartTimeFunc(5, 10),
				Duration:        60,
				InitialRateBump: 2,
				Points: []AuctionPointClassFixed{
					{Coefficient: 1, Delay: 2},
				},
				GasCost: GasCostConfigClassFixed{
					GasBumpEstimate:  1,
					GasPriceEstimate: 100,
				},
			},
			expectErr: false,
		},
		{
			name: "Invalid Gas Price Estimate",
			preset: &PresetClass{
				AllowMultipleFills: true,
				AllowPartialFills:  true,
				AuctionDuration:    60.0,
				AuctionEndAmount:   "1000",
				AuctionStartAmount: "500",
				BankFee:            "5",
				EstP:               0.1,
				ExclusiveResolver:  map[string]interface{}{"resolver": "value"},
				GasCost: GasCostConfigClass{
					GasBumpEstimate:  1.0,
					GasPriceEstimate: "invalid",
				},
				InitialRateBump: 0.2,
				Points: []AuctionPointClass{
					{Coefficient: 1.0, Delay: 2.0},
				},
				StartAuctionIn: 5.0,
				TokenFee:       "1",
			},
			additionalWaitPeriod: 10.0,
			expected:             nil,
			expectErr:            true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CreateAuctionDetails(tc.preset, tc.additionalWaitPeriod)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestBpsToRatioFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    *big.Int
		expected *big.Int
	}{
		{
			name:     "Nil input",
			input:    nil,
			expected: big.NewInt(0),
		},
		{
			name:     "Zero input",
			input:    big.NewInt(0),
			expected: big.NewInt(0),
		},
		{
			name:     "Positive input",
			input:    big.NewInt(5),
			expected: big.NewInt(50), // 5 * 100_000 / 10_000
		},
		{
			name:     "Negative input",
			input:    big.NewInt(-5),
			expected: big.NewInt(-50), // -5 * 100_000 / 10_000
		},
		{
			name:     "Large input",
			input:    big.NewInt(100_000),
			expected: big.NewInt(1_000_000), // 100_000 * 100_000 / 10_000
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := bpsToRatioFormat(tc.input)
			require.NotNil(t, result)
			assert.Equal(t, tc.expected, result)
		})
	}
}

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
			expectedErr: errors.New("nonce required when partial fill or multiple fill disallowed"),
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
