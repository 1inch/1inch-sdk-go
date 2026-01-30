package fusionplus

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
// This tests fusionplus-specific CreateMakerTraits wrapper with fusionplus's ExtraParams type
func TestMakerTraitsEncoding_KnownValues(t *testing.T) {
	tests := []struct {
		name           string
		details        Details
		extraParams    ExtraParams
		expectedEncode string
	}{
		{
			name: "Standard fusionplus order - partial and multiple fills allowed",
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
		{
			name: "Different deadline - longer duration",
			details: Details{
				Auction: &fusionorder.AuctionDetails{
					StartTime: 1673548149,
					Duration:  3600,
				},
			},
			extraParams: ExtraParams{
				AllowPartialFills:    true,
				AllowMultipleFills:   true,
				OrderExpirationDelay: 60,
			},
			expectedEncode: "0x4a000000000000000000000000000000000063c05fc100000000000000000000",
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
			name: "With permit2",
			extraParams: ExtraParams{
				AllowPartialFills:    true,
				AllowMultipleFills:   true,
				OrderExpirationDelay: 12,
				EnablePermit2:        true,
			},
		},
		{
			name: "No partial fills with nonce",
			extraParams: ExtraParams{
				Nonce:                big.NewInt(12345),
				AllowPartialFills:    false,
				AllowMultipleFills:   false,
				OrderExpirationDelay: 12,
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

	assert.NotEqual(t, encodings["Standard"], encodings["With permit2"], "permit2 flag should change encoding")
	assert.NotEqual(t, encodings["Standard"], encodings["No partial fills with nonce"], "partial fills flag should change encoding")
	assert.NotEqual(t, encodings["With permit2"], encodings["No partial fills with nonce"], "different flags should produce different encodings")
}

// TestAuctionDetailsEncoding_FusionPlus_KnownValues tests fusionplus-specific EncodeWithoutPointCount()
func TestAuctionDetailsEncoding_FusionPlus_KnownValues(t *testing.T) {
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
			expectedEncode: "0000000000000063c051750000b400c350004e20000c",
		},
		{
			name: "With multiple points",
			auctionDetails: &fusionorder.AuctionDetails{
				StartTime:       1673548149,
				Duration:        180,
				InitialRateBump: 50000,
				Points: []fusionorder.AuctionPointClassFixed{
					{Coefficient: 10000, Delay: 10},
					{Coefficient: 5000, Delay: 20},
				},
				GasCost: fusionorder.GasCostConfigClassFixed{GasBumpEstimate: 0, GasPriceEstimate: 0},
			},
			expectedEncode: "0000000000000063c051750000b400c350002710000a0013880014",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			encoded := tc.auctionDetails.EncodeWithoutPointCount()
			assert.Equal(t, tc.expectedEncode, encoded, "fusionorder.AuctionDetails encoding mismatch")

			encoded2 := tc.auctionDetails.EncodeWithoutPointCount()
			assert.Equal(t, encoded, encoded2, "fusionorder.AuctionDetails encoding should be deterministic")
		})
	}
}

// TestHashLock_KnownValues verifies HashLock computation against known expected values
func TestHashLock_KnownValues(t *testing.T) {
	tests := []struct {
		name     string
		secret   string
		expected string
	}{
		{
			name:     "Single fill hashlock",
			secret:   "0x531d1d2d7a594f1c7e413b074c7b693161486b5c495d457748144a01795c6a45",
			expected: "0x9f65fdcf781d4320c2dde70da02a1fe916d595dc1817149cc4758fd6a4bfd830",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ForSingleFill(tc.secret)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result.Value, "HashLock mismatch")
		})
	}
}

// TestHashLock_MultipleFills_KnownValues verifies merkle root hashlock computation
func TestHashLock_MultipleFills_KnownValues(t *testing.T) {
	tests := []struct {
		name     string
		secrets  []string
		expected string
	}{
		{
			name: "Three secrets merkle root",
			secrets: []string{
				"0x531d1d2d7a594f1c7e413b074c7b693161486b5c495d457748144a01795c6a45",
				"0x657812136b5000651d5e18516d764b5e661a681c760d3c3c4c15751020757823",
				"0x62071a322351281f04756576270c362a6e5b395e3b0f68027f231141555c3d43",
			},
			expected: "0x000292766d9172e4b4983ee4d4b6d511cdbcbef175c7e3e1b1554d513e1ab724",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			leaves, err := GetMerkleLeaves(tc.secrets)
			require.NoError(t, err)
			result, err := ForMultipleFills(leaves)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result.Value, "HashLock merkle root mismatch")
		})
	}
}

// TestEscrowExtraDataEncoding_KnownValues verifies escrow extra data encoding
func TestEscrowExtraDataEncoding_KnownValues(t *testing.T) {
	tests := []struct {
		name            string
		extraData       *EscrowExtraData
		expectedEncoded string
	}{
		{
			name: "Standard escrow extra data",
			extraData: &EscrowExtraData{
				HashLock: &HashLock{
					Value: "ad1723a873d05effcbdc57dcf7d00458d6a8c763558d5af7522bf6ad2d3e253d",
				},
				DstChainId:       42161,
				DstToken:         common.HexToAddress("0x0000000000000000000000000000000000000001"),
				SrcSafetyDeposit: big.NewInt(100),
				DstSafetyDeposit: big.NewInt(200),
				TimeLocks: &TimeLocks{
					DstCancellation:       3,
					DstPublicWithdrawal:   2,
					DstWithdrawal:         1,
					SrcPublicCancellation: 4,
					SrcCancellation:       3,
					SrcPublicWithdrawal:   2,
					SrcWithdrawal:         1,
				},
			},
			expectedEncoded: "ad1723a873d05effcbdc57dcf7d00458d6a8c763558d5af7522bf6ad2d3e253d000000000000000000000000000000000000000000000000000000000000a4b1000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000064000000000000000000000000000000c80000000000000003000000020000000100000004000000030000000200000001",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := encodeExtraData(tc.extraData)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedEncoded, common.Bytes2Hex(encoded), "EscrowExtraData encoding mismatch")
		})
	}
}

// TestExtensionPlusConversion verifies ExtensionPlus converts correctly to orderbook extension
func TestExtensionPlusConversion(t *testing.T) {
	tests := []struct {
		name      string
		extension ExtensionPlus
	}{
		{
			name: "Standard extension",
			extension: ExtensionPlus{
				MakerAssetSuffix: "0x1234",
				TakerAssetSuffix: "0x5678",
				MakingAmountData: "0xabcd",
				TakingAmountData: "0xef01",
				Predicate:        "0x2345",
				MakerPermit:      "0x6789",
				PreInteraction:   "0xabcd",
				PostInteraction:  "0xef01",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			obExtension := tc.extension.ConvertToOrderbookExtension()
			require.NotNil(t, obExtension)

			encoded, err := obExtension.Encode()
			require.NoError(t, err)
			require.NotEmpty(t, encoded)

			assert.True(t, len(encoded) >= 2 && encoded[:2] == "0x", "Encoded extension should start with 0x")

			decoded, err := orderbook.Decode(mustDecodeHexLocal(encoded))
			require.NoError(t, err)

			assert.Equal(t, obExtension.MakerAssetSuffix, decoded.MakerAssetSuffix)
			assert.Equal(t, obExtension.TakerAssetSuffix, decoded.TakerAssetSuffix)
		})
	}
}

// TestSaltGeneration_FusionPlus_KnownValues verifies EscrowExtension.GenerateSalt() with known values
func TestSaltGeneration_FusionPlus_KnownValues(t *testing.T) {
	originalBigIntMaxFunc := random_number_generation.BigIntMaxFunc
	random_number_generation.BigIntMaxFunc = func(max *big.Int) (*big.Int, error) {
		return big.NewInt(123456), nil
	}
	defer func() { random_number_generation.BigIntMaxFunc = originalBigIntMaxFunc }()

	tests := []struct {
		name      string
		extension *EscrowExtension
		expected  string
	}{
		{
			name: "Extension with all fields",
			extension: &EscrowExtension{
				ExtensionPlus: ExtensionPlus{
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
			},
			expected: "180431178743033967347942937469468920088249224033532329",
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

// TestSaltGeneration_Deterministic verifies salt generation is deterministic with mocked random
func TestSaltGeneration_Deterministic(t *testing.T) {
	originalBigIntMaxFunc := random_number_generation.BigIntMaxFunc
	random_number_generation.BigIntMaxFunc = func(max *big.Int) (*big.Int, error) {
		return big.NewInt(123456), nil
	}
	defer func() { random_number_generation.BigIntMaxFunc = originalBigIntMaxFunc }()

	extension := &EscrowExtension{
		ExtensionPlus: ExtensionPlus{
			MakerAssetSuffix: "0x1234",
			TakerAssetSuffix: "0x5678",
			MakingAmountData: "0xabcd",
			TakingAmountData: "0xef01",
			Predicate:        "0x2345",
			MakerPermit:      "0x6789",
			PreInteraction:   "0xabcd",
			PostInteraction:  "0xef01",
		},
	}

	salt1, err := extension.GenerateSalt()
	require.NoError(t, err)

	salt2, err := extension.GenerateSalt()
	require.NoError(t, err)

	assert.Equal(t, 0, salt1.Cmp(salt2), "Salt generation should be deterministic with mocked random")
}

// TestMerkleProof_KnownValues verifies merkle proof generation
func TestMerkleProof_KnownValues(t *testing.T) {
	tests := []struct {
		name          string
		secrets       []string
		leafIndex     int
		expectedProof []string
	}{
		{
			name: "Proof for first leaf",
			secrets: []string{
				"0x6466643931343237333333313437633162386632316365646666323931643738",
				"0x3131353932633266343034343466363562333230313837353438356463616130",
				"0x6634376135663837653765303462346261616566383430303662303336386635",
			},
			leafIndex:     0,
			expectedProof: []string{"0x540daf363747246d40b31da95b3ef1c1497e22e9a56b70d117c835839822c95f"},
		},
		{
			name: "Proof for first leaf (different secrets)",
			secrets: []string{
				"0x531d1d2d7a594f1c7e413b074c7b693161486b5c495d457748144a01795c6a45",
				"0x657812136b5000651d5e18516d764b5e661a681c760d3c3c4c15751020757823",
				"0x62071a322351281f04756576270c362a6e5b395e3b0f68027f231141555c3d43",
			},
			leafIndex:     0,
			expectedProof: []string{"0xb19c79aa34d58e459ce8119c301a24f7a01b8080ced7f3d608093e9e67624729"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			leaves, err := GetMerkleLeaves(tc.secrets)
			require.NoError(t, err)

			proof, err := GetProof(leaves, tc.leafIndex)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedProof, proof, "Merkle proof mismatch")
		})
	}
}

// TestTimeLocks_Encoding verifies time locks encoding is deterministic
func TestTimeLocks_Encoding(t *testing.T) {
	tests := []struct {
		name      string
		timeLocks *TimeLocks
	}{
		{
			name: "Standard timelocks",
			timeLocks: &TimeLocks{
				DstCancellation:       3,
				DstPublicWithdrawal:   2,
				DstWithdrawal:         1,
				SrcPublicCancellation: 4,
				SrcCancellation:       3,
				SrcPublicWithdrawal:   2,
				SrcWithdrawal:         1,
			},
		},
		{
			name: "Larger timelocks",
			timeLocks: &TimeLocks{
				DstCancellation:       3600,
				DstPublicWithdrawal:   1800,
				DstWithdrawal:         900,
				SrcPublicCancellation: 7200,
				SrcCancellation:       3600,
				SrcPublicWithdrawal:   1800,
				SrcWithdrawal:         900,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			extraData := &EscrowExtraData{
				HashLock: &HashLock{
					Value: "0000000000000000000000000000000000000000000000000000000000000000",
				},
				DstChainId:       1,
				DstToken:         common.HexToAddress("0x0000000000000000000000000000000000000001"),
				SrcSafetyDeposit: big.NewInt(0),
				DstSafetyDeposit: big.NewInt(0),
				TimeLocks:        tc.timeLocks,
			}

			encoded1, err := encodeExtraData(extraData)
			require.NoError(t, err)

			encoded2, err := encodeExtraData(extraData)
			require.NoError(t, err)

			assert.Equal(t, encoded1, encoded2, "TimeLocks encoding should be deterministic")
		})
	}
}

// TestEscrowExtension_KnownFields verifies escrow extension field construction
func TestEscrowExtension_KnownFields(t *testing.T) {
	tests := []struct {
		name                     string
		params                   EscrowExtensionParams
		expectedMakingAmountData string
		expectedTakingAmountData string
		expectedPostInteraction  string
	}{
		{
			name: "Basic escrow extension",
			params: EscrowExtensionParams{
				ExtensionParamsPlus: ExtensionParamsPlus{
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
						IntegratorFee: &IntegratorFee{
							Ratio:    big.NewInt(0),
							Receiver: common.Address{},
						},
						BankFee:            big.NewInt(0),
						ResolvingStartTime: big.NewInt(0),
						CustomReceiver:     common.Address{},
					},
					Asset:            "0x1234",
					Permit:           "0x3456",
					MakerAssetSuffix: "0x1234",
					TakerAssetSuffix: "0x1234",
					Predicate:        "0x1234",
					PreInteraction:   "pre",
				},
			},
			expectedMakingAmountData: "0x00000000000000000000000000000000000056780000000000000000000000000000000000",
			expectedTakingAmountData: "0x00000000000000000000000000000000000056780000000000000000000000000000000000",
			expectedPostInteraction:  "0x00000000000000000000000000000000000056780000000000",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ext, err := NewEscrowExtension(tc.params)
			require.NoError(t, err)
			require.NotNil(t, ext)

			assert.Equal(t, tc.expectedMakingAmountData, ext.MakingAmountData, "MakingAmountData mismatch")
			assert.Equal(t, tc.expectedTakingAmountData, ext.TakingAmountData, "TakingAmountData mismatch")
			assert.Equal(t, tc.expectedPostInteraction, ext.PostInteraction, "PostInteraction mismatch")
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
