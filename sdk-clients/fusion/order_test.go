package fusion

import (
	"errors"
	"math/big"
	"testing"

	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

func TestGetPreset(t *testing.T) {
	customPreset := &PresetClassFixed{
		AllowMultipleFills: true,
		AllowPartialFills:  true,
		AuctionDuration:    10.0,
		AuctionEndAmount:   "1000",
		AuctionStartAmount: "500",
		BankFee:            "5",
		EstP:               0.1,
		ExclusiveResolver:  "resolver",
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

	fastPreset := PresetClassFixed{
		AllowMultipleFills: false,
		AllowPartialFills:  false,
		AuctionDuration:    20.0,
		AuctionEndAmount:   "2000",
		AuctionStartAmount: "1000",
		BankFee:            "10",
		EstP:               0.2,
		ExclusiveResolver:  "resolver",
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

	mediumPreset := PresetClassFixed{
		AllowMultipleFills: true,
		AllowPartialFills:  false,
		AuctionDuration:    30.0,
		AuctionEndAmount:   "3000",
		AuctionStartAmount: "1500",
		BankFee:            "15",
		EstP:               0.3,
		ExclusiveResolver:  "resolver",
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

	slowPreset := PresetClassFixed{
		AllowMultipleFills: false,
		AllowPartialFills:  true,
		AuctionDuration:    40.0,
		AuctionEndAmount:   "4000",
		AuctionStartAmount: "2000",
		BankFee:            "20",
		EstP:               0.4,
		ExclusiveResolver:  "resolver",
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

	presets := QuotePresetsClassFixed{
		Custom: customPreset,
		Fast:   fastPreset,
		Medium: mediumPreset,
		Slow:   slowPreset,
	}

	tests := []struct {
		name       string
		presetType GetQuoteOutputRecommendedPreset
		expected   *PresetClassFixed
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
		preset               *PresetClassFixed
		additionalWaitPeriod float32
		expected             *AuctionDetails
		expectErr            bool
	}{
		{
			name: "Valid Preset",
			preset: &PresetClassFixed{
				AllowMultipleFills: true,
				AllowPartialFills:  true,
				AuctionDuration:    60.0,
				AuctionEndAmount:   "1000",
				AuctionStartAmount: "500",
				BankFee:            "5",
				EstP:               0.1,
				ExclusiveResolver:  "resolver",
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
			preset: &PresetClassFixed{
				AllowMultipleFills: true,
				AllowPartialFills:  true,
				AuctionDuration:    60.0,
				AuctionEndAmount:   "1000",
				AuctionStartAmount: "500",
				BankFee:            "5",
				EstP:               0.1,
				ExclusiveResolver:  "resolver",
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
				Fees: FeesOld{
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
				Fees: FeesOld{
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
				Fees: FeesOld{
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
		//{
		//	name: "Valid Details and Order Info with non-zero Delay",
		//	details: Details{
		//		ResolvingStartTime: big.NewInt(1622548800), // Example timestamp
		//		Fees: FeesOld{
		//			IntFee: IntegratorFee{
		//				Ratio:    big.NewInt(100),
		//				Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
		//			},
		//			BankFee: big.NewInt(200),
		//		},
		//		Whitelist: []AuctionWhitelistItem{
		//			{
		//				Address:   common.HexToAddress("0x0000000000000000000000000000000000000002"),
		//				AllowFrom: big.NewInt(1622549800),
		//			},
		//		},
		//	},
		//	orderInfo: FusionOrderV4{
		//		Receiver: "0x0000000000000000000000000000000000000003",
		//	},
		//	expected: &SettlementPostInteractionData{
		//		Whitelist: []WhitelistItem{
		//			{
		//				AddressHalf: "00000000000000000002",
		//				Delay:       big.NewInt(1000),
		//			},
		//		},
		//		IntegratorFee: &IntegratorFee{
		//			Ratio:    big.NewInt(100),
		//			Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
		//		},
		//		BankFee:            big.NewInt(200),
		//		ResolvingStartTime: big.NewInt(1622548800),
		//		CustomReceiver:     common.HexToAddress("0x0000000000000000000000000000000000000003"),
		//	},
		//	expectErr: false,
		//},
		//{
		//	name: "Valid Details and Order Info without Resolving Start Time",
		//	details: Details{
		//		ResolvingStartTime: nil,
		//		Fees: FeesOld{
		//			IntFee: IntegratorFee{
		//				Ratio:    big.NewInt(100),
		//				Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
		//			},
		//			BankFee: big.NewInt(200),
		//		},
		//		Whitelist: []AuctionWhitelistItem{
		//			{
		//				Address:   common.HexToAddress("0x0000000000000000000000000000000000000002"),
		//				AllowFrom: big.NewInt(1622548800),
		//			},
		//		},
		//	},
		//	orderInfo: FusionOrderV4{
		//		Receiver: "0x0000000000000000000000000000000000000003",
		//	},
		//	expected: &SettlementPostInteractionData{
		//		Whitelist: []WhitelistItem{
		//			{
		//				AddressHalf: "00000000000000000002",
		//				Delay:       big.NewInt(0),
		//			},
		//		},
		//		IntegratorFee: &IntegratorFee{
		//			Ratio:    big.NewInt(100),
		//			Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
		//		},
		//		BankFee:            big.NewInt(200),
		//		ResolvingStartTime: big.NewInt(timeNow()), // This will be dynamically set
		//		CustomReceiver:     common.HexToAddress("0x0000000000000000000000000000000000000003"),
		//	},
		//	expectErr: false,
		//},
		//{
		//	name: "Delay too large",
		//	details: Details{
		//		ResolvingStartTime: big.NewInt(1622548800), // Example timestamp
		//		Fees: FeesOld{
		//			IntFee: IntegratorFee{
		//				Ratio:    big.NewInt(100),
		//				Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
		//			},
		//			BankFee: big.NewInt(200),
		//		},
		//		Whitelist: []AuctionWhitelistItem{
		//			{
		//				Address:   common.HexToAddress("0x0000000000000000000000000000000000000002"),
		//				AllowFrom: big.NewInt(1622649800),
		//			},
		//		},
		//	},
		//	orderInfo: FusionOrderV4{
		//		Receiver: "0x0000000000000000000000000000000000000003",
		//	},
		//	expected: &SettlementPostInteractionData{
		//		Whitelist: []WhitelistItem{
		//			{
		//				AddressHalf: "00000000000000000002",
		//				Delay:       big.NewInt(1000),
		//			},
		//		},
		//		IntegratorFee: &IntegratorFee{
		//			Ratio:    big.NewInt(100),
		//			Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
		//		},
		//		BankFee:            big.NewInt(200),
		//		ResolvingStartTime: big.NewInt(1622548800),
		//		CustomReceiver:     common.HexToAddress("0x0000000000000000000000000000000000000003"),
		//	},
		//	expectErr:   true,
		//	expectedErr: fmt.Errorf("delay too big - %d must be less than %d", 101000, uint16Max),
		//},
		//{
		//	name: "Whitelist empty",
		//	details: Details{
		//		ResolvingStartTime: nil,
		//		Fees: FeesOld{
		//			IntFee: IntegratorFee{
		//				Ratio:    big.NewInt(100),
		//				Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
		//			},
		//			BankFee: big.NewInt(200),
		//		},
		//		Whitelist: []AuctionWhitelistItem{},
		//	},
		//	orderInfo: FusionOrderV4{
		//		Receiver: "0x0000000000000000000000000000000000000003",
		//	},
		//	expectErr:   true,
		//	expectedErr: errors.New("whitelist cannot be empty"),
		//},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			var whitelistStrings []string
			for _, whitelistItem := range tc.details.Whitelist {
				whitelistStrings = append(whitelistStrings, whitelistItem.Address.Hex())
			}

			// TODO this does not track AllowFrom anymore. Need to refactor these tests
			whitelist, err := GenerateWhitelist(whitelistStrings, tc.details.ResolvingStartTime)
			require.NoError(t, err)

			result, err := CreateSettlementPostInteractionData(tc.details, whitelist, tc.orderInfo)
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

var extensionContract = "0x8273f37417da37c4a6c3995e82cf442f87a25d9c"

//func TestCreateFusionOrder(t *testing.T) {
//
//	details := Details{
//		Auction: &AuctionDetails{
//			StartTime:       1673548149,
//			Duration:        180,
//			InitialRateBump: 50000,
//			Points:          []AuctionPointClassFixed{{Coefficient: 20000, Delay: 12}},
//		},
//		Fees: Fees{
//			IntFee:  IntegratorFee{},
//			BankFee: nil,
//		},
//		Whitelist:          []AuctionWhitelistItem{{Address: common.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fa"), AllowFrom: big.NewInt(0)}},
//		ResolvingStartTime: big.NewInt(1673548139),
//	}
//
//	postInteractionSufixDataNew := &SettlementSuffixData{
//		Whitelist:          details.Whitelist,
//		ResolvingStartTime: details.ResolvingStartTime,
//		CustomReceiver:     common.Address{},
//		Fees:               FeesNew{},
//	}
//
//	whitelist, err := GenerateWhitelist(postInteractionSufixDataNew)
//	require.NoError(t, err)
//
//	postInteractionDataNew := SettlementPostInteractionData{
//		Whitelist:          whitelist,
//		IntegratorFee:      &details.FeesOld.IntFee,
//		BankFee:            details.FeesOld.BankFee,
//		ResolvingStartTime: details.ResolvingStartTime,
//		CustomReceiver:     common.Address{},
//		AuctionFees:        FeesNew{},
//	}
//
//	extensionTmp, err := NewExtension(ExtensionParams{
//		SettlementContract:  extensionContract,
//		AuctionDetails:      details.Auction,
//		PostInteractionData: nil,
//		Asset:               "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
//		Surplus:             SurplusParamsNoFee,
//		ResolvingStartTime:  details.ResolvingStartTime,
//	})
//	require.NoError(t, err)
//
//	encodedPostInteractionData, err := postInteractionDataNew.EncodeNew(*extensionTmp)
//	require.NoError(t, err)
//
//	//postInteractionData, err := CreateSettlementPostInteractionData(details, FusionOrderV4{})
//	//if err != nil {
//	//	require.NoError(t, err)
//	//}
//
//	amountData, err := extensionTmp.BuildAmountGetterData(postInteractionDataNew.Whitelist, true)
//	require.NoError(t, err)
//
//	extension, err := NewExtension(ExtensionParams{
//		SettlementContract:         extensionContract,
//		AuctionDetails:             details.Auction,
//		PostInteractionData:        nil,
//		PostInteractionDataEncoded: encodedPostInteractionData,
//		Asset:                      "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
//		Surplus:                    SurplusParamsNoFee,
//		AmountData:                 amountData,
//	})
//	require.NoError(t, err)
//
//	extensionEncoded, err := extension.ConvertToOrderbookExtension().Encode()
//
//	var receiver = "0x0000000000000000000000000000000000000000" // needs to be explicitly set
//	extra := ExtraParams{
//		AllowPartialFills:    true,
//		AllowMultipleFills:   true,
//		OrderExpirationDelay: 12,
//		EnablePermit2:        false,
//	}
//
//	makerTraits, err := CreateMakerTraits(details, extra)
//	require.NoError(t, err)
//
//	var baseSalt int64 = 10
//	// Save the original function
//	originalBigIntMaxFunc := random_number_generation.BigIntMaxFunc
//	// Monkey patch the function
//	random_number_generation.BigIntMaxFunc = func(max *big.Int) (*big.Int, error) {
//		return big.NewInt(baseSalt), nil
//	}
//	// Restore the original function after the test
//	defer func() {
//		random_number_generation.BigIntMaxFunc = originalBigIntMaxFunc
//	}()
//	salt, err := extension.GenerateSalt()
//	require.NoError(t, err)
//
//	orderData := orderbook.OrderData{
//		MakerAsset:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
//		TakerAsset:   "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
//		MakingAmount: "1000000000000000000",
//		TakingAmount: "1420000000",
//		Salt:         salt.String(),
//		Maker:        "0x00000000219ab540356cbb839cbe05303d7705fa",
//		Receiver:     getReceiver(details.FeesNew, extensionContract, receiver),
//		MakerTraits:  makerTraits.Encode(),
//		Extension:    extensionEncoded,
//	}
//
//	order := FusionOrderV4{
//		Maker:        orderData.Maker,
//		MakerAsset:   orderData.MakerAsset,
//		MakerTraits:  orderData.MakerTraits,
//		MakingAmount: orderData.MakingAmount,
//		Receiver:     orderData.Receiver,
//		Salt:         orderData.Salt,
//		TakerAsset:   orderData.TakerAsset,
//		TakingAmount: orderData.TakingAmount,
//	}
//
//	assert.Equal(t, "0x00000000219ab540356cbb839cbe05303d7705fa", order.Maker)
//	assert.Equal(t, "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", order.MakerAsset)
//	assert.Equal(t, "1000000000000000000", order.MakingAmount)
//	assert.Equal(t, "0x0000000000000000000000000000000000000000", order.Receiver)
//	assert.Equal(t, "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", order.TakerAsset)
//	assert.Equal(t, "1420000000", order.TakingAmount)
//	assert.Equal(t, "0x4a000000000000000000000000000000000063c0523500000000000000000000", order.MakerTraits)
//	assert.Equal(t, "15367444032321112222254918300472459458657037823940", order.Salt)
//	assert.Equal(t, "0x8273f37417Da37c4A6c3995E82Cf442f87a25D9c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006463c0516b01bb839cbe05303d7705fa0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00", extension.PostInteraction)
//}

func TestCreateFusionOrderTdd(t *testing.T) {
	tests := []struct {
		name                    string
		details                 Details
		expected                FusionOrderV4
		expectedPostInteraction string
	}{
		{
			name: "basic fusion order (no fees)",
			details: Details{
				Auction: &AuctionDetails{
					StartTime:       1673548149,
					Duration:        180,
					InitialRateBump: 50000,
					Points:          []AuctionPointClassFixed{{Coefficient: 20000, Delay: 12}},
				},
				Whitelist:          []AuctionWhitelistItem{{Address: common.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fa"), AllowFrom: big.NewInt(0)}},
				ResolvingStartTime: big.NewInt(1673548139),
			},
			expected: FusionOrderV4{
				Maker:        "0x00000000219ab540356cbb839cbe05303d7705fa",
				MakerAsset:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
				MakingAmount: "1000000000000000000",
				Receiver:     "0x0000000000000000000000000000000000000000",
				TakerAsset:   "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				TakingAmount: "1420000000",
				MakerTraits:  "0x4a000000000000000000000000000000000063c0523500000000000000000000",
				Salt:         "15367444032321112222254918300472459458657037823940",
			},
			expectedPostInteraction: "0x8273f37417da37c4a6c3995e82cf442f87a25d9c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006463c0516b01bb839cbe05303d7705fa0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00",
		},
		{
			name: "fusion order with integrator fees",
			details: Details{
				Auction: &AuctionDetails{
					StartTime:       1673548149,
					Duration:        180,
					InitialRateBump: 50000,
					Points:          []AuctionPointClassFixed{{Coefficient: 20000, Delay: 12}},
				},
				FeesNew: &FeesNew{
					Integrator: IntegratorFeeNew{
						Integrator: "0x0000000000000000000000000000000000000001",
						Protocol:   "0x0000000000000000000000000000000000000002",
						Fee:        FromPercent(1, GetDefaultBase()),
						Share:      FromPercent(50, GetDefaultBase()),
					},
				},
				Whitelist: []AuctionWhitelistItem{
					{Address: common.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fa"), AllowFrom: big.NewInt(0)},
				},
				ResolvingStartTime: big.NewInt(1673548139),
			},
			expected: FusionOrderV4{
				Maker:        "0x00000000219ab540356cbb839cbe05303d7705fa",
				MakerAsset:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
				MakingAmount: "1000000000000000000",
				Receiver:     "0x8273f37417da37c4a6c3995e82cf442f87a25d9c", // extension contract address because fees exist
				TakerAsset:   "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				TakingAmount: "1420000000",
				MakerTraits:  "0x4a000000000000000000000000000000000063c0523500000000000000000000",
				Salt:         "15367444032321112222254918300472459458657037823940",
			},
			expectedPostInteraction: "0x8273f37417da37c4a6c3995e82cf442f87a25d9c000000000000000000000000000000000000000001000000000000000000000000000000000000000203e83200006463c0516b01bb839cbe05303d7705fa0000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postInteractionSufixDataNew := &SettlementSuffixData{
				Whitelist:          tt.details.Whitelist,
				ResolvingStartTime: tt.details.ResolvingStartTime,
			}

			var whitelistStrings []string
			for _, whitelistItem := range tt.details.Whitelist {
				whitelistStrings = append(whitelistStrings, whitelistItem.Address.Hex())
			}

			whitelist, err := GenerateWhitelist(whitelistStrings, tt.details.ResolvingStartTime)
			require.NoError(t, err)

			extra := ExtraParams{
				AllowPartialFills:    true,
				AllowMultipleFills:   true,
				OrderExpirationDelay: 12,
				EnablePermit2:        false,
			}

			makerTraits, err := CreateMakerTraits(tt.details, extra)
			require.NoError(t, err)

			// When fees are present, use the settlement contract as the custom receiver
			postInteractionData := &SettlementPostInteractionData{
				Whitelist:          whitelist,
				IntegratorFee:      &tt.details.Fees.IntFee,
				BankFee:            tt.details.Fees.BankFee,
				ResolvingStartTime: tt.details.ResolvingStartTime,
				CustomReceiver:     common.Address{},
				AuctionFees:        tt.details.FeesNew,
			}

			//extensionTmp, err := NewExtension(ExtensionParams{
			//	SettlementContract:  extensionContract,
			//	AuctionDetails:      tt.details.Auction,
			//	PostInteractionData: nil,
			//	Asset:               "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			//	Surplus:             SurplusParamsNoFee,
			//	ResolvingStartTime:  tt.details.ResolvingStartTime,
			//})
			//// TODO put this in the constructor
			//extensionTmp.Fees = tt.details.FeesNew
			//
			//require.NoError(t, err)

			//encodedPostInteractionData, err := postInteractionData.EncodeNew(*extensionTmp)
			//require.NoError(t, err)

			extension, err := NewExtension(ExtensionParams{
				SettlementContract:  extensionContract,
				AuctionDetails:      tt.details.Auction,
				PostInteractionData: postInteractionData,
				//PostInteractionDataEncoded: encodedPostInteractionData,
				Asset:              "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
				Surplus:            SurplusParamsNoFee,
				ResolvingStartTime: tt.details.ResolvingStartTime,
			})
			require.NoError(t, err)

			extensionEncoded, err := extension.ConvertToOrderbookExtension().Encode()
			require.NoError(t, err)

			baseSalt := int64(10)
			originalBigIntMaxFunc := random_number_generation.BigIntMaxFunc
			random_number_generation.BigIntMaxFunc = func(max *big.Int) (*big.Int, error) {
				return big.NewInt(baseSalt), nil
			}
			defer func() { random_number_generation.BigIntMaxFunc = originalBigIntMaxFunc }()

			salt, err := extension.GenerateSalt()
			require.NoError(t, err)

			orderData := orderbook.OrderData{
				MakerAsset:   tt.expected.MakerAsset,
				TakerAsset:   tt.expected.TakerAsset,
				MakingAmount: tt.expected.MakingAmount,
				TakingAmount: tt.expected.TakingAmount,
				Salt:         salt.String(),
				Maker:        tt.expected.Maker,
				Receiver:     getReceiver(tt.details.FeesNew, extensionContract, postInteractionSufixDataNew.CustomReceiver.Hex()),
				MakerTraits:  makerTraits.Encode(),
				Extension:    extensionEncoded,
			}

			order := FusionOrderV4{
				Maker:        orderData.Maker,
				MakerAsset:   orderData.MakerAsset,
				MakerTraits:  orderData.MakerTraits,
				MakingAmount: orderData.MakingAmount,
				Receiver:     orderData.Receiver,
				Salt:         orderData.Salt,
				TakerAsset:   orderData.TakerAsset,
				TakingAmount: orderData.TakingAmount,
			}

			assert.Equal(t, tt.expected.Maker, order.Maker)
			assert.Equal(t, tt.expected.MakerAsset, order.MakerAsset)
			assert.Equal(t, tt.expected.MakingAmount, order.MakingAmount)
			assert.Equal(t, tt.expected.Receiver, order.Receiver)
			assert.Equal(t, tt.expected.TakerAsset, order.TakerAsset)
			assert.Equal(t, tt.expected.TakingAmount, order.TakingAmount)
			assert.Equal(t, tt.expected.MakerTraits, order.MakerTraits)
			//assert.Equal(t, tt.expected.Salt, order.Salt)
			assert.Equal(t, tt.expectedPostInteraction, extension.PostInteraction)
		})
	}
}
