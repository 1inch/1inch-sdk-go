package fusion

import (
	"math/big"
	"strings"
	"testing"

	"github.com/1inch/1inch-sdk-go/common/fusionorder"
	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMakerTraitsEncoding_KnownValues verifies MakerTraits encoding against known expected values
// This tests fusion-specific CreateMakerTraits wrapper with fusion's ExtraParams type
func TestMakerTraitsEncoding_KnownValues(t *testing.T) {
	tests := []struct {
		name           string
		details        Details
		extraParams    ExtraParams
		expectedEncode string
	}{
		{
			name: "Standard fusion order - partial and multiple fills allowed",
			details: Details{
				Auction: &fusionorder.AuctionDetails{
					StartTime: 1673548149,
					Duration:  180,
				},
			},
			extraParams: ExtraParams{
				Nonce:                nil,
				AllowPartialFills:    true,
				AllowMultipleFills:   true,
				OrderExpirationDelay: 12,
			},
			expectedEncode: "0x4a000000000000000000000000000000000063c0523500000000000000000000",
		},
		{
			name: "No partial fills - requires nonce",
			details: Details{
				Auction: &fusionorder.AuctionDetails{
					StartTime: 1673548149,
					Duration:  180,
				},
			},
			extraParams: ExtraParams{
				Nonce:                big.NewInt(12345),
				AllowPartialFills:    false,
				AllowMultipleFills:   false,
				OrderExpirationDelay: 12,
			},
			expectedEncode: "0x8a000000000000000000000000000030390063c0523500000000000000000000",
		},
		{
			name: "With unwrap WETH flag",
			details: Details{
				Auction: &fusionorder.AuctionDetails{
					StartTime: 1673548149,
					Duration:  180,
				},
			},
			extraParams: ExtraParams{
				AllowPartialFills:    true,
				AllowMultipleFills:   true,
				OrderExpirationDelay: 12,
				unwrapWeth:           true,
			},
			expectedEncode: "0x4a800000000000000000000000000000000063c0523500000000000000000000",
		},
		{
			name: "With permit2 flag",
			details: Details{
				Auction: &fusionorder.AuctionDetails{
					StartTime: 1673548149,
					Duration:  180,
				},
			},
			extraParams: ExtraParams{
				AllowPartialFills:    true,
				AllowMultipleFills:   true,
				OrderExpirationDelay: 12,
				EnablePermit2:        true,
			},
			expectedEncode: "0x4b000000000000000000000000000000000063c0523500000000000000000000",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			makerTraits, err := CreateMakerTraits(tc.details, tc.extraParams)
			require.NoError(t, err)

			encoded := makerTraits.Encode()
			assert.Equal(t, tc.expectedEncode, encoded, "MakerTraits encoding mismatch")

			assert.True(t, strings.HasPrefix(encoded, "0x"), "MakerTraits should start with 0x")
			assert.Equal(t, 66, len(encoded), "MakerTraits should be 32 bytes (66 chars with 0x)")
		})
	}
}

// TestMakerTraitsEncoding_FlagVariations verifies different flag combinations produce different encodings
func TestMakerTraitsEncoding_FlagVariations(t *testing.T) {
	baseDetails := Details{
		Auction: &fusionorder.AuctionDetails{
			StartTime: 1673548149,
			Duration:  180,
		},
	}

	tests := []struct {
		name        string
		extraParams ExtraParams
	}{
		{
			name: "Standard",
			extraParams: ExtraParams{
				AllowPartialFills:    true,
				AllowMultipleFills:   true,
				OrderExpirationDelay: 12,
			},
		},
		{
			name: "With unwrap WETH",
			extraParams: ExtraParams{
				AllowPartialFills:    true,
				AllowMultipleFills:   true,
				OrderExpirationDelay: 12,
				unwrapWeth:           true,
			},
		},
		{
			name: "With permit2",
			extraParams: ExtraParams{
				AllowPartialFills:    true,
				AllowMultipleFills:   true,
				OrderExpirationDelay: 12,
				EnablePermit2:        true,
			},
		},
	}

	encodings := make(map[string]string)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			makerTraits, err := CreateMakerTraits(baseDetails, tc.extraParams)
			require.NoError(t, err)
			encodings[tc.name] = makerTraits.Encode()
		})
	}

	assert.NotEqual(t, encodings["Standard"], encodings["With unwrap WETH"], "unwrapWeth flag should change encoding")
	assert.NotEqual(t, encodings["Standard"], encodings["With permit2"], "permit2 flag should change encoding")
	assert.NotEqual(t, encodings["With unwrap WETH"], encodings["With permit2"], "different flags should produce different encodings")
}

// TestExtensionKeccak256_Deterministic verifies that extension hash is deterministic
func TestExtensionKeccak256_Deterministic(t *testing.T) {
	whitelist := []fusionorder.WhitelistItem{
		{AddressHalf: "bb839cbe05303d7705fa", Delay: big.NewInt(0)},
	}

	postInteractionData := &SettlementPostInteractionData{
		Whitelist:          whitelist,
		ResolvingStartTime: big.NewInt(1673548139),
		CustomReceiver:     common.Address{},
		AuctionFees:        nil,
	}

	auctionDetails := &fusionorder.AuctionDetails{
		StartTime:       1673548149,
		Duration:        180,
		InitialRateBump: 50000,
		Points:          []fusionorder.AuctionPointClassFixed{{Coefficient: 20000, Delay: 12}},
		GasCost:         fusionorder.GasCostConfigClassFixed{GasBumpEstimate: 10000, GasPriceEstimate: 1000000},
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

	hash1, err := extension.Keccak256()
	require.NoError(t, err)

	hash2, err := extension.Keccak256()
	require.NoError(t, err)

	assert.Equal(t, hash1.Cmp(hash2), 0, "Extension hash should be deterministic")
	assert.NotNil(t, hash1, "Hash should not be nil")
	assert.True(t, hash1.Sign() >= 0, "Hash should be non-negative")
}

// TestOrderDataConsistency verifies that order data fields are consistent
func TestOrderDataConsistency(t *testing.T) {
	originalBigIntMaxFunc := random_number_generation.BigIntMaxFunc
	random_number_generation.BigIntMaxFunc = func(max *big.Int) (*big.Int, error) {
		return big.NewInt(12345678), nil
	}
	defer func() { random_number_generation.BigIntMaxFunc = originalBigIntMaxFunc }()

	settlementAddress := "0x8273f37417da37c4a6c3995e82cf442f87a25d9c"

	whitelist := []fusionorder.WhitelistItem{
		{AddressHalf: "bb839cbe05303d7705fa", Delay: big.NewInt(0)},
	}

	postInteractionData := &SettlementPostInteractionData{
		Whitelist:          whitelist,
		ResolvingStartTime: big.NewInt(1673548139),
		CustomReceiver:     common.Address{},
		AuctionFees:        nil,
	}

	auctionDetails := &fusionorder.AuctionDetails{
		StartTime:       1673548149,
		Duration:        180,
		InitialRateBump: 50000,
		Points:          []fusionorder.AuctionPointClassFixed{{Coefficient: 20000, Delay: 12}},
		GasCost:         fusionorder.GasCostConfigClassFixed{GasBumpEstimate: 10000, GasPriceEstimate: 1000000},
	}

	extension, err := NewExtension(ExtensionParams{
		SettlementContract:  settlementAddress,
		AuctionDetails:      auctionDetails,
		PostInteractionData: postInteractionData,
		Asset:               "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		Surplus:             SurplusParamsNoFee,
		ResolvingStartTime:  big.NewInt(1673548139),
	})
	require.NoError(t, err)

	details := Details{
		Auction:            auctionDetails,
		ResolvingStartTime: big.NewInt(1673548139),
		Whitelist:          []fusionorder.AuctionWhitelistItem{{Address: common.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fa"), AllowFrom: big.NewInt(0)}},
	}

	extraParams := ExtraParams{
		Nonce:                big.NewInt(12345),
		AllowPartialFills:    true,
		AllowMultipleFills:   true,
		OrderExpirationDelay: 12,
	}

	makerTraits, err := CreateMakerTraits(details, extraParams)
	require.NoError(t, err)

	orderInfo := FusionOrderV4{
		Maker:        "0x1234567890123456789012345678901234567890",
		MakerAsset:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		TakerAsset:   "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		MakingAmount: "1000000000000000000",
		TakingAmount: "1420000000",
		Receiver:     "0x9876543210987654321098765432109876543210",
	}

	params := CreateOrderDataParams{
		Extension:           extension,
		SettlementAddress:   settlementAddress,
		PostInteractionData: postInteractionData,
		orderInfo:           orderInfo,
		Details:             details,
		ExtraParams:         extraParams,
		MakerTraits:         makerTraits,
	}

	order1, err := CreateOrder(params)
	require.NoError(t, err)

	order2, err := CreateOrder(params)
	require.NoError(t, err)

	assert.Equal(t, order1.Inner.MakerAsset, order2.Inner.MakerAsset)
	assert.Equal(t, order1.Inner.TakerAsset, order2.Inner.TakerAsset)
	assert.Equal(t, order1.Inner.Salt, order2.Inner.Salt, "Salt should be deterministic with mocked random")
	assert.Equal(t, order1.Inner.MakerTraits, order2.Inner.MakerTraits)
}

// TestAuctionDetailsEncoding_Fusion_KnownValues tests fusion-specific AuctionDetails.Encode()
func TestAuctionDetailsEncoding_Fusion_KnownValues(t *testing.T) {
	tests := []struct {
		name           string
		auctionDetails *fusionorder.AuctionDetails
		expectedEncode string
	}{
		{
			name: "Standard auction - no gas cost",
			auctionDetails: &fusionorder.AuctionDetails{
				StartTime:       1673548149,
				Duration:        180,
				InitialRateBump: 50000,
				Points:          []fusionorder.AuctionPointClassFixed{{Coefficient: 20000, Delay: 12}},
				GasCost:         fusionorder.GasCostConfigClassFixed{GasBumpEstimate: 0, GasPriceEstimate: 0},
			},
			expectedEncode: "0000000000000063c051750000b400c35001004e20000c",
		},
		{
			name: "With gas cost",
			auctionDetails: &fusionorder.AuctionDetails{
				StartTime:       1673548149,
				Duration:        180,
				InitialRateBump: 50000,
				Points:          []fusionorder.AuctionPointClassFixed{{Coefficient: 20000, Delay: 12}},
				GasCost:         fusionorder.GasCostConfigClassFixed{GasBumpEstimate: 10000, GasPriceEstimate: 1000000},
			},
			expectedEncode: "002710000f424063c051750000b400c35001004e20000c",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			encoded := tc.auctionDetails.Encode()
			assert.Equal(t, tc.expectedEncode, encoded, "fusionorder.AuctionDetails encoding mismatch")
		})
	}
}

// TestOrderbookExtensionConversion verifies fusion extension converts correctly
func TestOrderbookExtensionConversion(t *testing.T) {
	whitelist := []fusionorder.WhitelistItem{
		{AddressHalf: "bb839cbe05303d7705fa", Delay: big.NewInt(0)},
	}

	postInteractionData := &SettlementPostInteractionData{
		Whitelist:          whitelist,
		ResolvingStartTime: big.NewInt(1673548139),
		CustomReceiver:     common.Address{},
		AuctionFees:        nil,
	}

	auctionDetails := &fusionorder.AuctionDetails{
		StartTime:       1673548149,
		Duration:        180,
		InitialRateBump: 50000,
		Points:          []fusionorder.AuctionPointClassFixed{{Coefficient: 20000, Delay: 12}},
		GasCost:         fusionorder.GasCostConfigClassFixed{GasBumpEstimate: 10000, GasPriceEstimate: 1000000},
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

	obExtension := extension.ConvertToOrderbookExtension()
	require.NotNil(t, obExtension)

	encoded, err := obExtension.Encode()
	require.NoError(t, err)
	require.NotEmpty(t, encoded)

	assert.True(t, len(encoded) >= 2 && encoded[:2] == "0x", "Encoded extension should start with 0x")

	decoded, err := orderbook.Decode(mustDecodeHexLocal(encoded))
	require.NoError(t, err)

	assert.Equal(t, obExtension.MakerAssetSuffix, decoded.MakerAssetSuffix)
	assert.Equal(t, obExtension.TakerAssetSuffix, decoded.TakerAssetSuffix)
}

// TestBuildAmountGetterData_KnownValues verifies amount getter data construction (fusion-specific)
func TestBuildAmountGetterData_KnownValues(t *testing.T) {
	tests := []struct {
		name             string
		auctionDetails   *fusionorder.AuctionDetails
		whitelist        []fusionorder.WhitelistItem
		resolvingTime    *big.Int
		forAmountGetters bool
		expected         string
	}{
		{
			name: "Basic auction details with forAmountGetters true",
			auctionDetails: &fusionorder.AuctionDetails{
				GasCost:         fusionorder.GasCostConfigClassFixed{GasBumpEstimate: 0, GasPriceEstimate: 0},
				StartTime:       1673548149,
				Duration:        180,
				InitialRateBump: 50000,
				Points:          []fusionorder.AuctionPointClassFixed{{Delay: 12, Coefficient: 20000}},
			},
			whitelist:        []fusionorder.WhitelistItem{{AddressHalf: "bb839cbe05303d7705fa", Delay: big.NewInt(0)}},
			resolvingTime:    big.NewInt(1673548139),
			forAmountGetters: true,
			expected:         "0x0000000000000063c051750000b400c35001004e20000c00000000006401bb839cbe05303d7705fa",
		},
		{
			name: "Basic auction details with forAmountGetters false",
			auctionDetails: &fusionorder.AuctionDetails{
				GasCost:         fusionorder.GasCostConfigClassFixed{GasBumpEstimate: 0, GasPriceEstimate: 0},
				StartTime:       1673548149,
				Duration:        180,
				InitialRateBump: 50000,
				Points:          []fusionorder.AuctionPointClassFixed{{Delay: 12, Coefficient: 20000}},
			},
			whitelist:        []fusionorder.WhitelistItem{{AddressHalf: "bb839cbe05303d7705fa", Delay: big.NewInt(0)}},
			resolvingTime:    big.NewInt(1673548139),
			forAmountGetters: false,
			expected:         "0x00000000006463c0516b01bb839cbe05303d7705fa0000",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			params := &BuildAmountGetterDataParams{
				AuctionDetails:     tc.auctionDetails,
				ResolvingStartTime: tc.resolvingTime,
				PostInteractionData: &SettlementPostInteractionData{
					Whitelist:          tc.whitelist,
					ResolvingStartTime: tc.resolvingTime,
					CustomReceiver:     common.Address{},
					AuctionFees:        nil,
				},
			}

			result, err := BuildAmountGetterData(params, tc.forAmountGetters)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result, "BuildAmountGetterData mismatch")
		})
	}
}

// TestExtensionCreation_KnownValues verifies fusion extension fields are constructed correctly
func TestExtensionCreation_KnownValues(t *testing.T) {
	tests := []struct {
		name                     string
		params                   ExtensionParams
		expectedMakingAmountData string
		expectedTakingAmountData string
		expectedMakerPermit      string
		expectedPostInteraction  string
	}{
		{
			name: "Extension with surplus params",
			params: ExtensionParams{
				SettlementContract: "0x5678",
				AuctionDetails: &fusionorder.AuctionDetails{
					StartTime:       0,
					Duration:        0,
					InitialRateBump: 0,
					Points:          nil,
					GasCost:         fusionorder.GasCostConfigClassFixed{},
				},
				PostInteractionData: &SettlementPostInteractionData{
					Whitelist: []fusionorder.WhitelistItem{},
					AuctionFees: &FeesIntegratorAndResolver{
						Resolver:   ResolverFee{},
						Integrator: IntegratorFee{},
					},
					ResolvingStartTime: big.NewInt(0),
					CustomReceiver:     common.Address{},
				},
				Asset:              "0x1234",
				Permit:             "0x3456",
				MakerAssetSuffix:   "0x1234",
				TakerAssetSuffix:   "0x1234",
				Predicate:          "0x1234",
				PreInteraction:     "0x5678",
				Surplus:            &SurplusParams{EstimatedTakerAmount: big.NewInt(1), ProtocolFee: fusionorder.MustFromPercent(1, fusionorder.GetDefaultBase())},
				ResolvingStartTime: big.NewInt(0),
			},
			expectedMakingAmountData: "0x000000000000000000000000000000000000567800000000000000000000000000000000000000000000006400",
			expectedTakingAmountData: "0x000000000000000000000000000000000000567800000000000000000000000000000000000000000000006400",
			expectedMakerPermit:      "0x00000000000000000000000000000000000012343456",
			expectedPostInteraction:  "0x000000000000000000000000000000000000567800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000000000000000101",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ext, err := NewExtension(tc.params)
			require.NoError(t, err)
			require.NotNil(t, ext)

			assert.Equal(t, tc.expectedMakingAmountData, ext.MakingAmountData, "MakingAmountData mismatch")
			assert.Equal(t, tc.expectedTakingAmountData, ext.TakingAmountData, "TakingAmountData mismatch")
			assert.Equal(t, tc.expectedMakerPermit, ext.MakerPermit, "MakerPermit mismatch")
			assert.Equal(t, tc.expectedPostInteraction, ext.PostInteraction, "PostInteraction mismatch")
		})
	}
}

// TestSaltGeneration_Fusion_KnownValues verifies fusion Extension.GenerateSalt()
func TestSaltGeneration_Fusion_KnownValues(t *testing.T) {
	originalBigIntMaxFunc := random_number_generation.BigIntMaxFunc
	random_number_generation.BigIntMaxFunc = func(max *big.Int) (*big.Int, error) {
		return big.NewInt(123456), nil
	}
	defer func() { random_number_generation.BigIntMaxFunc = originalBigIntMaxFunc }()

	tests := []struct {
		name      string
		extension *Extension
		expected  string
	}{
		{
			name: "Extension with all fields",
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
			expected: "180431909497609865807168059378624320943465639784996571",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.extension.GenerateSalt()
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result.String(), "Salt generation mismatch")
		})
	}
}

// Helper function to decode hex string to bytes
func mustDecodeHexLocal(s string) []byte {
	if len(s) >= 2 && s[:2] == "0x" {
		s = s[2:]
	}
	if len(s)%2 != 0 {
		s = "0" + s
	}
	b := make([]byte, len(s)/2)
	for i := 0; i < len(b); i++ {
		b[i] = hexCharToByteLocal(s[2*i])<<4 | hexCharToByteLocal(s[2*i+1])
	}
	return b
}

func hexCharToByteLocal(c byte) byte {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}
