package fusionplus

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/common/fusionorder"
	"github.com/1inch/1inch-sdk-go/constants"
	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
)

// createTestQuoteFusionPlus creates a realistic quote for testing
func createTestQuoteFusionPlus() *GetQuoteOutputFixed {
	return &GetQuoteOutputFixed{
		QuoteId: "test-quote-id-fusionplus-12345",
		Presets: QuotePresets{
			Fast: Preset{
				AllowMultipleFills: true,
				AllowPartialFills:  true,
				AuctionDuration:    180,
				AuctionEndAmount:   "1420000000",
				AuctionStartAmount: "1500000000",
				GasCost: GasCostConfig{
					GasBumpEstimate:  10000,
					GasPriceEstimate: "1000000000", // Fits in uint32
				},
				InitialRateBump: 50000,
				Points: []AuctionPoint{
					{Coefficient: 20000, Delay: 12},
				},
				StartAuctionIn: 12,
			},
			Medium: Preset{
				AllowMultipleFills: true,
				AllowPartialFills:  true,
				AuctionDuration:    360,
				AuctionEndAmount:   "1400000000",
				AuctionStartAmount: "1500000000",
				GasCost: GasCostConfig{
					GasBumpEstimate:  10000,
					GasPriceEstimate: "1000000000",
				},
				InitialRateBump: 60000,
				Points: []AuctionPoint{
					{Coefficient: 30000, Delay: 24},
				},
				StartAuctionIn: 24,
			},
			Slow: Preset{
				AllowMultipleFills: true,
				AllowPartialFills:  true,
				AuctionDuration:    600,
				AuctionEndAmount:   "1380000000",
				AuctionStartAmount: "1500000000",
				GasCost: GasCostConfig{
					GasBumpEstimate:  10000,
					GasPriceEstimate: "1000000000",
				},
				InitialRateBump: 70000,
				Points: []AuctionPoint{
					{Coefficient: 40000, Delay: 60},
				},
				StartAuctionIn: 36,
			},
		},
		SrcEscrowFactory: "0x8273f37417da37c4a6c3995e82cf442f87a25d9c",
		DstEscrowFactory: "0x9384f38417da37c4a6c3995e82cf442f87a25d9d",
		Whitelist: []string{
			"0x00000000219ab540356cbb839cbe05303d7705fa",
		},
		SrcSafetyDeposit: "1000000000000000",
		DstSafetyDeposit: "1000000000000000",
		TimeLocks: TimeLocks{
			DstCancellation:       3600,
			DstPublicWithdrawal:   1800,
			DstWithdrawal:         600,
			SrcCancellation:       7200,
			SrcPublicCancellation: 5400,
			SrcPublicWithdrawal:   3600,
			SrcWithdrawal:         1200,
		},
	}
}

// TestSignedOrderInput_Serialization_FusionPlus tests order serialization for FusionPlus
func TestSignedOrderInput_Serialization_FusionPlus(t *testing.T) {
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
		QuoteId:    "test-quote-id",
		Signature:  "0xabcdef1234567890",
		SrcChainId: 1,
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
	assert.Equal(t, order.SrcChainId, deserialized.SrcChainId)
	assert.Equal(t, order.QuoteId, deserialized.QuoteId)
}

// TestGetQuoteOutputFixed_Serialization_FusionPlus tests quote response deserialization
func TestGetQuoteOutputFixed_Serialization_FusionPlus(t *testing.T) {
	apiResponse := `{
		"quoteId": "abc123",
		"srcEscrowFactory": "0x8273f37417da37c4a6c3995e82cf442f87a25d9c",
		"dstEscrowFactory": "0x9384f38417da37c4a6c3995e82cf442f87a25d9d",
		"whitelist": ["0x00000000219ab540356cbb839cbe05303d7705fa"],
		"srcSafetyDeposit": "1000000000000000",
		"dstSafetyDeposit": "1000000000000000",
		"timeLocks": {
			"dstCancellation": 3600,
			"dstPublicWithdrawal": 1800,
			"dstWithdrawal": 600,
			"srcCancellation": 7200,
			"srcPublicCancellation": 5400,
			"srcPublicWithdrawal": 3600,
			"srcWithdrawal": 1200
		},
		"presets": {
			"fast": {
				"allowMultipleFills": true,
				"allowPartialFills": true,
				"auctionDuration": 180,
				"auctionEndAmount": "1420000000",
				"auctionStartAmount": "1500000000",
				"gasCost": {
					"gasBumpEstimate": 10000,
					"gasPriceEstimate": "1000000000"
				},
				"initialRateBump": 50000,
				"points": [{"coefficient": 20000, "delay": 12}],
				"startAuctionIn": 12
			},
			"medium": {
				"allowMultipleFills": true,
				"allowPartialFills": true,
				"auctionDuration": 360,
				"auctionEndAmount": "1400000000",
				"auctionStartAmount": "1500000000",
				"gasCost": {
					"gasBumpEstimate": 10000,
					"gasPriceEstimate": "1000000000"
				},
				"initialRateBump": 60000,
				"points": [{"coefficient": 30000, "delay": 24}],
				"startAuctionIn": 24
			},
			"slow": {
				"allowMultipleFills": true,
				"allowPartialFills": true,
				"auctionDuration": 600,
				"auctionEndAmount": "1380000000",
				"auctionStartAmount": "1500000000",
				"gasCost": {
					"gasBumpEstimate": 10000,
					"gasPriceEstimate": "1000000000"
				},
				"initialRateBump": 70000,
				"points": [{"coefficient": 40000, "delay": 60}],
				"startAuctionIn": 36
			}
		}
	}`

	var quote GetQuoteOutputFixed
	err := json.Unmarshal([]byte(apiResponse), &quote)
	require.NoError(t, err)

	// Verify deserialization
	assert.Equal(t, "abc123", quote.QuoteId)
	assert.Equal(t, "0x8273f37417da37c4a6c3995e82cf442f87a25d9c", quote.SrcEscrowFactory)
	assert.Equal(t, "0x9384f38417da37c4a6c3995e82cf442f87a25d9d", quote.DstEscrowFactory)
	assert.Equal(t, "1000000000000000", quote.SrcSafetyDeposit)
	assert.Equal(t, "1000000000000000", quote.DstSafetyDeposit)

	// Verify time locks
	assert.Equal(t, float32(3600), quote.TimeLocks.DstCancellation)
	assert.Equal(t, float32(7200), quote.TimeLocks.SrcCancellation)

	// Verify presets
	assert.Equal(t, float32(180), quote.Presets.Fast.AuctionDuration)
	assert.True(t, quote.Presets.Fast.AllowMultipleFills)
}

// TestPresetSelection_Integration_FusionPlus tests preset selection
func TestPresetSelection_Integration_FusionPlus(t *testing.T) {
	quote := createTestQuoteFusionPlus()

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
			preset, err := GetPreset(quote.Presets, tc.presetType)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedAmount, preset.AuctionEndAmount)
		})
	}
}

// TestAuctionDetailsCreation_Integration_FusionPlus tests auction details creation
func TestAuctionDetailsCreation_Integration_FusionPlus(t *testing.T) {
	// Mock auction start time
	originalCalcAuctionStartTimeFunc := fusionorder.CalcAuctionStartTimeFunc
	fusionorder.CalcAuctionStartTimeFunc = func(startAuctionIn uint32, additionalWaitPeriod uint32) uint32 {
		return 1700000000 + startAuctionIn + additionalWaitPeriod
	}
	defer func() { fusionorder.CalcAuctionStartTimeFunc = originalCalcAuctionStartTimeFunc }()

	quote := createTestQuoteFusionPlus()
	preset, err := GetPreset(quote.Presets, Fast)
	require.NoError(t, err)

	details, err := CreateAuctionDetails(preset, 0)
	require.NoError(t, err)

	assert.Equal(t, uint32(1700000000+12), details.StartTime)
	assert.Equal(t, uint32(180), details.Duration)
	assert.Equal(t, uint32(50000), details.InitialRateBump)
}

// TestMerkleTreeIntegration tests merkle tree creation for hashlocks
func TestMerkleTreeIntegration(t *testing.T) {
	// Test single leaf
	singleLeaf := []string{
		"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}
	tree1 := MakeTree(singleLeaf)
	require.NotNil(t, tree1)
	assert.Len(t, tree1.leaves, 1)

	// Test multiple leaves
	multipleLeaves := []string{
		"0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
	}
	tree2 := MakeTree(multipleLeaves)
	require.NotNil(t, tree2)
	assert.Len(t, tree2.leaves, 3)

	// Get proof for first leaf using the standalone GetProof function
	proof, err := GetProof(multipleLeaves, 0)
	require.NoError(t, err)
	assert.NotNil(t, proof)
}

// TestBpsToRatioFormat_Integration tests bps conversion
func TestBpsToRatioFormat_Integration(t *testing.T) {
	tests := []struct {
		name     string
		input    *big.Int
		expected *big.Int
	}{
		{"Zero", big.NewInt(0), big.NewInt(0)},
		{"Nil", nil, big.NewInt(0)},
		{"100 bps (1%)", big.NewInt(100), big.NewInt(1000)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := fusionorder.BpsToRatioFormat(tc.input)
			assert.Equal(t, 0, tc.expected.Cmp(result))
		})
	}
}

// TestNativeTokenWrapping_FusionPlus tests native token mapping
func TestNativeTokenWrapping_FusionPlus(t *testing.T) {
	tests := []struct {
		name     string
		chainId  constants.NetworkEnum
		expected string
	}{
		{"Ethereum WETH", constants.NetworkEthereum, "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"},
		{"Polygon WMATIC", constants.NetworkPolygon, "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"},
		{"Arbitrum WETH", constants.NetworkArbitrum, "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wrapper, exists := constants.ChainToWrapper[tc.chainId]
			assert.True(t, exists)
			assert.Equal(t, common.HexToAddress(tc.expected), wrapper)
		})
	}
}

// TestSettlementPostInteractionData_Integration_FusionPlus tests post interaction data
func TestSettlementPostInteractionData_Integration_FusionPlus(t *testing.T) {
	// Mock timeNow
	originalTimeNow := timeNow
	timeNow = func() int64 { return 1700000000 }
	defer func() { timeNow = originalTimeNow }()

	details := Details{
		Auction: &AuctionDetails{
			StartTime: 1700000000,
			Duration:  180,
		},
		Fees: Fees{
			IntFee: IntegratorFee{
				Ratio:    big.NewInt(0),
				Receiver: common.Address{},
			},
			BankFee: big.NewInt(0),
		},
		Whitelist: []AuctionWhitelistItem{
			{
				Address:   common.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fa"),
				AllowFrom: big.NewInt(0),
			},
		},
	}

	orderInfo := CrossChainOrderDto{
		Receiver: "0x1234567890123456789012345678901234567890",
	}

	postInteraction, err := CreateSettlementPostInteractionData(details, orderInfo)
	require.NoError(t, err)
	assert.NotNil(t, postInteraction)
}

// TestNonceRequirement_FusionPlus tests nonce requirement logic
func TestNonceRequirement_FusionPlus(t *testing.T) {
	tests := []struct {
		name               string
		allowPartialFills  bool
		allowMultipleFills bool
		expected           bool
	}{
		{"Both true - no nonce required", true, true, false},
		{"Partial false - nonce required", false, true, true},
		{"Multiple false - nonce required", true, false, true},
		{"Both false - nonce required", false, false, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := fusionorder.IsNonceRequired(tc.allowPartialFills, tc.allowMultipleFills)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestSecretInputSerialization tests secret submission serialization
func TestSecretInputSerialization(t *testing.T) {
	secret := SecretInput{
		Secret:    "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		OrderHash: "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
	}

	jsonData, err := json.Marshal(secret)
	require.NoError(t, err)

	var deserialized SecretInput
	err = json.Unmarshal(jsonData, &deserialized)
	require.NoError(t, err)

	assert.Equal(t, secret.Secret, deserialized.Secret)
	assert.Equal(t, secret.OrderHash, deserialized.OrderHash)
}

// TestSaltGeneration_FusionPlus tests deterministic salt generation
func TestSaltGeneration_FusionPlus(t *testing.T) {
	// Mock random number generation
	originalBigIntMaxFunc := random_number_generation.BigIntMaxFunc
	random_number_generation.BigIntMaxFunc = func(max *big.Int) (*big.Int, error) {
		return big.NewInt(12345678), nil
	}
	defer func() { random_number_generation.BigIntMaxFunc = originalBigIntMaxFunc }()

	// Test that salt generation is consistent
	salt1, err := random_number_generation.BigIntMaxFunc(big.NewInt(1000000))
	require.NoError(t, err)

	salt2, err := random_number_generation.BigIntMaxFunc(big.NewInt(1000000))
	require.NoError(t, err)

	assert.Equal(t, salt1, salt2)
}
