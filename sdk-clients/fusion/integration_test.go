package fusion

import (
	"encoding/json"
	"math/big"
	"testing"

	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
	"github.com/1inch/1inch-sdk-go/sdk-clients/fusionorder"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestQuote creates a realistic quote for testing
func createTestQuote() GetQuoteOutputFixed {
	return GetQuoteOutputFixed{
		QuoteId: "test-quote-id-12345",
		Presets: QuotePresetsClassFixed{
			Fast: PresetClassFixed{
				AllowMultipleFills: true,
				AllowPartialFills:  true,
				AuctionDuration:    180,
				AuctionEndAmount:   "1420000000",
				AuctionStartAmount: "1500000000",
				BankFee:            "0",
				EstP:               0.95,
				ExclusiveResolver:  "",
				GasCost: GasCostConfigClass{
					GasBumpEstimate:  10000,
					GasPriceEstimate: "50000000000",
				},
				InitialRateBump: 50000,
				Points: []AuctionPointClass{
					{Coefficient: 20000, Delay: 12},
					{Coefficient: 10000, Delay: 24},
				},
				StartAuctionIn: 12,
				TokenFee:       "0",
			},
			Medium: PresetClassFixed{
				AllowMultipleFills: true,
				AllowPartialFills:  true,
				AuctionDuration:    360,
				AuctionEndAmount:   "1400000000",
				AuctionStartAmount: "1500000000",
				BankFee:            "0",
				EstP:               0.93,
				ExclusiveResolver:  "",
				GasCost: GasCostConfigClass{
					GasBumpEstimate:  10000,
					GasPriceEstimate: "50000000000",
				},
				InitialRateBump: 60000,
				Points: []AuctionPointClass{
					{Coefficient: 30000, Delay: 24},
					{Coefficient: 15000, Delay: 48},
				},
				StartAuctionIn: 24,
				TokenFee:       "0",
			},
			Slow: PresetClassFixed{
				AllowMultipleFills: true,
				AllowPartialFills:  true,
				AuctionDuration:    600,
				AuctionEndAmount:   "1380000000",
				AuctionStartAmount: "1500000000",
				BankFee:            "0",
				EstP:               0.90,
				ExclusiveResolver:  "",
				GasCost: GasCostConfigClass{
					GasBumpEstimate:  10000,
					GasPriceEstimate: "50000000000",
				},
				InitialRateBump: 70000,
				Points: []AuctionPointClass{
					{Coefficient: 40000, Delay: 60},
					{Coefficient: 20000, Delay: 120},
				},
				StartAuctionIn: 36,
				TokenFee:       "0",
			},
		},
		SettlementAddress: "0x8273f37417da37c4a6c3995e82cf442f87a25d9c",
		Whitelist: []string{
			"0x00000000219ab540356cbb839cbe05303d7705fa",
			"0x1111111111111111111111111111111111111111",
		},
		MarketAmount: "1500000000",
		SurplusFee:   1,
	}
}

// TestSignedOrderInput_Serialization tests that orders serialize correctly for the API
func TestSignedOrderInput_Serialization(t *testing.T) {
	order := SignedOrderInput{
		Extension: "0x1234567890abcdef",
		Order: OrderInput{
			Maker:        "0x1234567890123456789012345678901234567890",
			MakerAsset:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			MakerTraits:  "0x4a000000000000000000000000000000000063c0523500000000000000000000",
			MakingAmount: "1000000000000000000",
			Receiver:     "0x0000000000000000000000000000000000000000",
			Salt:         "12345678901234567890",
			TakerAsset:   "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
			TakingAmount: "1420000000",
		},
		QuoteId:   "test-quote-id",
		Signature: "0xabcdef1234567890",
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(order)
	require.NoError(t, err)

	// Deserialize back
	var deserialized SignedOrderInput
	err = json.Unmarshal(jsonData, &deserialized)
	require.NoError(t, err)

	// Verify round-trip
	assert.Equal(t, order.Extension, deserialized.Extension)
	assert.Equal(t, order.Order.Maker, deserialized.Order.Maker)
	assert.Equal(t, order.Order.MakerAsset, deserialized.Order.MakerAsset)
	assert.Equal(t, order.Order.MakerTraits, deserialized.Order.MakerTraits)
	assert.Equal(t, order.Order.MakingAmount, deserialized.Order.MakingAmount)
	assert.Equal(t, order.Order.Receiver, deserialized.Order.Receiver)
	assert.Equal(t, order.Order.Salt, deserialized.Order.Salt)
	assert.Equal(t, order.Order.TakerAsset, deserialized.Order.TakerAsset)
	assert.Equal(t, order.Order.TakingAmount, deserialized.Order.TakingAmount)
	assert.Equal(t, order.QuoteId, deserialized.QuoteId)
	assert.Equal(t, order.Signature, deserialized.Signature)
}

// TestGetQuoteOutputFixed_Serialization tests quote response deserialization
func TestGetQuoteOutputFixed_Serialization(t *testing.T) {
	// Simulate API response JSON
	apiResponse := `{
		"quoteId": "abc123",
		"settlementAddress": "0x8273f37417da37c4a6c3995e82cf442f87a25d9c",
		"whitelist": ["0x00000000219ab540356cbb839cbe05303d7705fa"],
		"marketAmount": "1500000000",
		"surplusFee": 1,
		"presets": {
			"fast": {
				"allowMultipleFills": true,
				"allowPartialFills": true,
				"auctionDuration": 180,
				"auctionEndAmount": "1420000000",
				"auctionStartAmount": "1500000000",
				"bankFee": "0",
				"estP": 0.95,
				"exclusiveResolver": "",
				"gasCost": {
					"gasBumpEstimate": 10000,
					"gasPriceEstimate": "50000000000"
				},
				"initialRateBump": 50000,
				"points": [{"coefficient": 20000, "delay": 12}],
				"startAuctionIn": 12,
				"tokenFee": "0"
			},
			"medium": {
				"allowMultipleFills": true,
				"allowPartialFills": true,
				"auctionDuration": 360,
				"auctionEndAmount": "1400000000",
				"auctionStartAmount": "1500000000",
				"bankFee": "0",
				"estP": 0.93,
				"exclusiveResolver": "",
				"gasCost": {
					"gasBumpEstimate": 10000,
					"gasPriceEstimate": "50000000000"
				},
				"initialRateBump": 60000,
				"points": [{"coefficient": 30000, "delay": 24}],
				"startAuctionIn": 24,
				"tokenFee": "0"
			},
			"slow": {
				"allowMultipleFills": true,
				"allowPartialFills": true,
				"auctionDuration": 600,
				"auctionEndAmount": "1380000000",
				"auctionStartAmount": "1500000000",
				"bankFee": "0",
				"estP": 0.90,
				"exclusiveResolver": "",
				"gasCost": {
					"gasBumpEstimate": 10000,
					"gasPriceEstimate": "50000000000"
				},
				"initialRateBump": 70000,
				"points": [{"coefficient": 40000, "delay": 60}],
				"startAuctionIn": 36,
				"tokenFee": "0"
			}
		}
	}`

	var quote GetQuoteOutputFixed
	err := json.Unmarshal([]byte(apiResponse), &quote)
	require.NoError(t, err)

	// Verify deserialization
	assert.Equal(t, "abc123", quote.QuoteId)
	assert.Equal(t, "0x8273f37417da37c4a6c3995e82cf442f87a25d9c", quote.SettlementAddress)
	assert.Len(t, quote.Whitelist, 1)
	assert.Equal(t, "1500000000", quote.MarketAmount)
	assert.Equal(t, float32(1), quote.SurplusFee)

	// Verify presets
	assert.Equal(t, float32(180), quote.Presets.Fast.AuctionDuration)
	assert.Equal(t, "1420000000", quote.Presets.Fast.AuctionEndAmount)
	assert.True(t, quote.Presets.Fast.AllowMultipleFills)
	assert.True(t, quote.Presets.Fast.AllowPartialFills)
}

// TestExtensionEncoding_Integration tests that extension encoding works end-to-end
func TestExtensionEncoding_Integration(t *testing.T) {
	// Create a realistic extension
	whitelist := []WhitelistItem{
		{AddressHalf: "bb839cbe05303d7705fa", Delay: big.NewInt(0)},
	}

	postInteractionData := &SettlementPostInteractionData{
		Whitelist:          whitelist,
		ResolvingStartTime: big.NewInt(1673548139),
		CustomReceiver:     common.Address{},
		AuctionFees:        nil,
	}

	auctionDetails := &AuctionDetails{
		StartTime:       1673548149,
		Duration:        180,
		InitialRateBump: 50000,
		Points:          []AuctionPointClassFixed{{Coefficient: 20000, Delay: 12}},
	}

	extension, err := NewExtension(ExtensionParams{
		SettlementContract:  "0x8273f37417da37c4a6c3995e82cf442f87a25d9c",
		AuctionDetails:      auctionDetails,
		PostInteractionData: postInteractionData,
		Asset:               "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		Surplus:             SurplusParamsNoFee,
		ResolvingStartTime:  big.NewInt(1673548139),
	})
	require.NoError(t, err)

	// Convert to orderbook extension and encode
	orderbookExt := extension.ConvertToOrderbookExtension()
	encoded, err := orderbookExt.Encode()
	require.NoError(t, err)

	// Verify the encoding is valid hex
	assert.True(t, len(encoded) > 2 && encoded[:2] == "0x")

	// Verify the extension can generate a salt
	salt, err := extension.GenerateSalt()
	require.NoError(t, err)
	assert.NotNil(t, salt)
	assert.True(t, salt.Sign() > 0)
}

// TestPresetSelection_Integration tests that preset selection works correctly
func TestPresetSelection_Integration(t *testing.T) {
	quote := createTestQuote()

	tests := []struct {
		name           string
		presetType     GetQuoteOutputRecommendedPreset
		expectedAmount string
	}{
		{"Fast preset", Fast, "1420000000"},
		{"Medium preset", Medium, "1400000000"},
		{"Slow preset", Slow, "1380000000"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			preset, err := getPreset(quote.Presets, tc.presetType)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedAmount, preset.AuctionEndAmount)
		})
	}
}

// TestAuctionDetailsCreation_Integration tests auction details creation
func TestAuctionDetailsCreation_Integration(t *testing.T) {
	// Mock auction start time for deterministic tests
	originalCalcAuctionStartTimeFunc := fusionorder.CalcAuctionStartTimeFunc
	fusionorder.CalcAuctionStartTimeFunc = func(startAuctionIn uint32, additionalWaitPeriod uint32) uint32 {
		return 1700000000 + startAuctionIn + additionalWaitPeriod
	}
	defer func() { fusionorder.CalcAuctionStartTimeFunc = originalCalcAuctionStartTimeFunc }()

	// Create a preset with gas price that fits in uint32
	preset := &PresetClassFixed{
		AllowMultipleFills: true,
		AllowPartialFills:  true,
		AuctionDuration:    180,
		AuctionEndAmount:   "1420000000",
		AuctionStartAmount: "1500000000",
		GasCost: GasCostConfigClass{
			GasBumpEstimate:  10000,
			GasPriceEstimate: "1000000000", // 1 gwei - fits in uint32
		},
		InitialRateBump: 50000,
		Points: []AuctionPointClass{
			{Coefficient: 20000, Delay: 12},
			{Coefficient: 10000, Delay: 24},
		},
		StartAuctionIn: 12,
	}

	details, err := CreateAuctionDetails(preset, 10)
	require.NoError(t, err)

	// Verify auction details
	assert.Equal(t, uint32(1700000000+12+10), details.StartTime) // base + startAuctionIn + additionalWait
	assert.Equal(t, uint32(180), details.Duration)
	assert.Equal(t, uint32(50000), details.InitialRateBump)
	assert.Len(t, details.Points, 2)
}

// TestWhitelistGeneration_Integration tests whitelist generation
func TestWhitelistGeneration_Integration(t *testing.T) {
	quote := createTestQuote()
	resolvingStartTime := big.NewInt(1673548139)

	whitelist, err := GenerateWhitelist(quote.Whitelist, resolvingStartTime)
	require.NoError(t, err)

	// Should have same number of items
	assert.Len(t, whitelist, len(quote.Whitelist))

	// Each item should have an address half (last 10 bytes = 20 hex chars)
	for _, item := range whitelist {
		assert.Len(t, item.AddressHalf, 20)
		assert.NotNil(t, item.Delay)
	}
}

// TestSettlementPostInteractionData_Integration tests post interaction data creation
func TestSettlementPostInteractionData_Integration(t *testing.T) {
	quote := createTestQuote()
	resolvingStartTime := big.NewInt(1673548139)

	whitelist, err := GenerateWhitelist(quote.Whitelist, resolvingStartTime)
	require.NoError(t, err)

	details := Details{
		ResolvingStartTime: resolvingStartTime,
		FeesIntAndRes:      nil,
	}

	orderInfo := FusionOrderV4{
		Receiver: "0x1234567890123456789012345678901234567890",
	}

	postInteraction, err := CreateSettlementPostInteractionData(details, whitelist, orderInfo)
	require.NoError(t, err)

	assert.Equal(t, whitelist, postInteraction.Whitelist)
	assert.Equal(t, resolvingStartTime, postInteraction.ResolvingStartTime)
	assert.Equal(t, common.HexToAddress(orderInfo.Receiver), postInteraction.CustomReceiver)
}

// TestMakerTraitsCreation_Integration tests maker traits creation
func TestMakerTraitsCreation_Integration(t *testing.T) {
	// Note: The SDK has specific constraints:
	// - If AllowPartialFills is false, AllowMultipleFills must also be false
	// - If AllowMultipleFills is false, AllowPartialFills must also be false
	// This means partial/multiple fills must both be true or both be false
	tests := []struct {
		name              string
		allowPartialFills bool
		allowMultipleFills bool
		nonce             *big.Int
		expectError       bool
	}{
		{
			name:              "Both fills allowed - no nonce required",
			allowPartialFills: true,
			allowMultipleFills: true,
			nonce:             nil,
			expectError:       false,
		},
		{
			name:              "Both fills disabled with nonce - valid",
			allowPartialFills: false,
			allowMultipleFills: false,
			nonce:             big.NewInt(12345),
			expectError:       false,
		},
		{
			name:              "Both fills disabled without nonce - error",
			allowPartialFills: false,
			allowMultipleFills: false,
			nonce:             nil,
			expectError:       true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			details := Details{
				Auction: &AuctionDetails{
					StartTime: 1700000000,
					Duration:  180,
				},
			}

			extraParams := ExtraParams{
				Nonce:                tc.nonce,
				AllowPartialFills:    tc.allowPartialFills,
				AllowMultipleFills:   tc.allowMultipleFills,
				OrderExpirationDelay: 12,
			}

			makerTraits, err := CreateMakerTraits(details, extraParams)

			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "nonce required")
			} else {
				require.NoError(t, err)
				assert.NotNil(t, makerTraits)
				assert.Equal(t, tc.allowPartialFills, makerTraits.AllowPartialFills)
				assert.Equal(t, tc.allowMultipleFills, makerTraits.AllowMultipleFills)
			}
		})
	}
}

// TestSaltGeneration_Integration tests salt generation with extension
func TestSaltGeneration_Integration(t *testing.T) {
	// Mock random number generation for deterministic tests
	originalBigIntMaxFunc := random_number_generation.BigIntMaxFunc
	random_number_generation.BigIntMaxFunc = func(max *big.Int) (*big.Int, error) {
		return big.NewInt(12345678), nil
	}
	defer func() { random_number_generation.BigIntMaxFunc = originalBigIntMaxFunc }()

	whitelist := []WhitelistItem{
		{AddressHalf: "bb839cbe05303d7705fa", Delay: big.NewInt(0)},
	}

	postInteractionData := &SettlementPostInteractionData{
		Whitelist:          whitelist,
		ResolvingStartTime: big.NewInt(1673548139),
		CustomReceiver:     common.Address{},
		AuctionFees:        nil,
	}

	auctionDetails := &AuctionDetails{
		StartTime:       1673548149,
		Duration:        180,
		InitialRateBump: 50000,
		Points:          []AuctionPointClassFixed{{Coefficient: 20000, Delay: 12}},
	}

	extension, err := NewExtension(ExtensionParams{
		SettlementContract:  "0x8273f37417da37c4a6c3995e82cf442f87a25d9c",
		AuctionDetails:      auctionDetails,
		PostInteractionData: postInteractionData,
		Asset:               "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		Surplus:             SurplusParamsNoFee,
		ResolvingStartTime:  big.NewInt(1673548139),
	})
	require.NoError(t, err)

	// Generate salt twice - should be deterministic with mocked random
	salt1, err := extension.GenerateSalt()
	require.NoError(t, err)

	salt2, err := extension.GenerateSalt()
	require.NoError(t, err)

	// With deterministic random, salts should be the same
	assert.Equal(t, salt1, salt2)
}

// TestOrderResponseSerialization tests that order responses deserialize correctly
func TestOrderResponseSerialization(t *testing.T) {
	apiResponse := `{
		"orderHash": "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		"status": "filled",
		"order": {
			"maker": "0x1234567890123456789012345678901234567890",
			"makerAsset": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			"makingAmount": "1000000000000000000",
			"receiver": "0x0000000000000000000000000000000000000000",
			"takerAsset": "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
			"takingAmount": "1420000000"
		}
	}`

	var response OrderResponse
	err := json.Unmarshal([]byte(apiResponse), &response)
	require.NoError(t, err)

	assert.Equal(t, "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890", response.OrderHash)
	assert.Equal(t, "filled", string(response.Status))
}

// TestFeesIntegration tests fee handling in orders
func TestFeesIntegration(t *testing.T) {
	// Test with integrator fees
	fees := &FeesIntegratorAndResolver{
		Integrator: IntegratorFee{
			Integrator: "0x0000000000000000000000000000000000000001",
			Protocol:   "0x0000000000000000000000000000000000000002",
			Fee:        fusionorder.MustFromPercent(1, fusionorder.GetDefaultBase()),
			Share:      fusionorder.MustFromPercent(50, fusionorder.GetDefaultBase()),
		},
		Resolver: ResolverFee{},
	}

	settlementAddress := "0x8273f37417da37c4a6c3995e82cf442f87a25d9c"
	receiver := "0x1234567890123456789012345678901234567890"

	// When fees exist, receiver should be settlement address
	actualReceiver := getReceiver(fees, settlementAddress, receiver)
	assert.Equal(t, settlementAddress, actualReceiver)

	// When no fees, receiver should be the original
	actualReceiver = getReceiver(nil, settlementAddress, receiver)
	assert.Equal(t, receiver, actualReceiver)
}

// TestNativeTokenWrapping tests native token to wrapped token conversion
func TestNativeTokenWrapping(t *testing.T) {
	tests := []struct {
		name            string
		chainId         fusionorder.NetworkEnum
		expectedWrapper string
	}{
		{"Ethereum WETH", fusionorder.ETHEREUM, "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"},
		{"Polygon WMATIC", fusionorder.POLYGON, "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"},
		{"Arbitrum WETH", fusionorder.ARBITRUM, "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1"},
		{"Binance WBNB", fusionorder.BINANCE, "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wrapper, exists := fusionorder.ChainToWrapper[tc.chainId]
			assert.True(t, exists)
			assert.Equal(t, common.HexToAddress(tc.expectedWrapper), wrapper)
		})
	}
}
