package orderbook

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/1inch/1inch-sdk-go/internal/bigint"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMakerTraitsEncoding_KnownValues verifies MakerTraits encoding against known expected values
func TestMakerTraitsEncoding_KnownValues(t *testing.T) {
	tests := []struct {
		name                string
		makerTraitParams    MakerTraitsParams
		expectedMakerTraits string
	}{
		{
			name: "Extension with expiration",
			makerTraitParams: MakerTraitsParams{
				AllowedSender:      "0x0000000000000000000000000000000000000000",
				ShouldCheckEpoch:   false,
				UsePermit2:         false,
				UnwrapWeth:         false,
				HasExtension:       true,
				HasPreInteraction:  false,
				HasPostInteraction: true,
				AllowPartialFills:  true,
				AllowMultipleFills: true,
				Expiry:             1715201499,
				Nonce:              0,
				Series:             0,
			},
			expectedMakerTraits: "0x4a0000000000000000000000000000000000663be5db00000000000000000000",
		},
		{
			name: "With permit2",
			makerTraitParams: MakerTraitsParams{
				AllowedSender:      "0x0000000000000000000000000000000000000000",
				UsePermit2:         true,
				HasExtension:       true,
				HasPostInteraction: true,
				AllowPartialFills:  true,
				AllowMultipleFills: true,
				Expiry:             1715201499,
			},
			expectedMakerTraits: "0x4b0000000000000000000000000000000000663be5db00000000000000000000",
		},
		{
			name: "With unwrap WETH",
			makerTraitParams: MakerTraitsParams{
				AllowedSender:      "0x0000000000000000000000000000000000000000",
				UnwrapWeth:         true,
				HasExtension:       true,
				HasPostInteraction: true,
				AllowPartialFills:  true,
				AllowMultipleFills: true,
				Expiry:             1715201499,
			},
			expectedMakerTraits: "0x4a8000000000000000000000000000000000663be5db00000000000000000000",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			makerTraits, err := NewMakerTraits(tc.makerTraitParams)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedMakerTraits, makerTraits.Encode(), "MakerTraits encoding mismatch")
		})
	}
}

// TestTakerTraitsEncoding_KnownValues verifies TakerTraits encoding against known expected values
func TestTakerTraitsEncoding_KnownValues(t *testing.T) {
	tests := []struct {
		name               string
		takerTraitParams   TakerTraitsParams
		expectedTraitFlags string
		expectedArgs       string
	}{
		{
			name: "Extension only",
			takerTraitParams: TakerTraitsParams{
				Extension: "0x000000f4000000f4000000f4000000000000000000000000000000000000000045c32fa6df82ead1e2ef74d17b76547eddfaff8900000000000000000000000050c5df26654b5efbdd0c54a062dfa6012933defe000000000000000000000000111111125421ca6dc452d289314280a0f8842a65000000000000000000000000000000000000000000000000002386f26fc1000000000000000000000000000000000000000000000000000000000000663a478b000000000000000000000000000000000000000000000000000000000000001bdf138a0d223e2ef8635075f5fe68efa8a2da1d890fdc3825b754c7ba2083ca0464494f534829f576cd67b966059657c51aaf53edbd6498d51cbd07da8bdb256b",
			},
			expectedTraitFlags: "7440945280133576583328096164017418065923851860621198004784596428783616",
			expectedArgs:       "0x000000f4000000f4000000f4000000000000000000000000000000000000000045c32fa6df82ead1e2ef74d17b76547eddfaff8900000000000000000000000050c5df26654b5efbdd0c54a062dfa6012933defe000000000000000000000000111111125421ca6dc452d289314280a0f8842a65000000000000000000000000000000000000000000000000002386f26fc1000000000000000000000000000000000000000000000000000000000000663a478b000000000000000000000000000000000000000000000000000000000000001bdf138a0d223e2ef8635075f5fe68efa8a2da1d890fdc3825b754c7ba2083ca0464494f534829f576cd67b966059657c51aaf53edbd6498d51cbd07da8bdb256b",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			expectedTakerTraitsBig, err := bigint.FromString(tc.expectedTraitFlags)
			require.NoError(t, err)

			takerTraits := NewTakerTraits(tc.takerTraitParams)
			encoded, err := takerTraits.Encode()
			require.NoError(t, err)

			assert.True(t, expectedTakerTraitsBig.Cmp(encoded.TraitFlags) == 0,
				"TakerTraits encoding mismatch: expected %x, got %x", expectedTakerTraitsBig, encoded.TraitFlags)

			expectedArgs := common.FromHex(tc.expectedArgs)
			assert.Equal(t, expectedArgs, encoded.Args, "TakerTraits args mismatch")
		})
	}
}

// TestSignatureCompression_KnownValues verifies signature compression against known expected values
func TestSignatureCompression_KnownValues(t *testing.T) {
	tests := []struct {
		name            string
		signature       string
		expectedRValue  string
		expectedVSValue string
	}{
		{
			name:            "v = 0x1b (27)",
			signature:       "4ca2e082038ead998ff272272153b02643135a47bf8a820e5cebb5303763a3211117316de39fbe43781e78e60db312c9173ee7da13fec73ef8366af0d89c9f3c1b",
			expectedRValue:  "4ca2e082038ead998ff272272153b02643135a47bf8a820e5cebb5303763a321",
			expectedVSValue: "1117316de39fbe43781e78e60db312c9173ee7da13fec73ef8366af0d89c9f3c",
		},
		{
			name:            "v = 0x1c (28)",
			signature:       "2fac11bfe002d84bd0837f6efc88688bf4a35309bb5cfde80f740105ddbc9e024e552465e5087d9997739ba467e161c9752364d16cebaf9afd9f8e1a8f22addc1c",
			expectedRValue:  "2fac11bfe002d84bd0837f6efc88688bf4a35309bb5cfde80f740105ddbc9e02",
			expectedVSValue: "ce552465e5087d9997739ba467e161c9752364d16cebaf9afd9f8e1a8f22addc",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			compactSignature, err := CompressSignature(tc.signature)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedRValue, fmt.Sprintf("%x", compactSignature.R), "R value mismatch")
			assert.Equal(t, tc.expectedVSValue, fmt.Sprintf("%x", compactSignature.VS), "VS value mismatch")
		})
	}
}

// TestSaltGeneration_KnownValues verifies salt generation against known expected values
func TestSaltGeneration_KnownValues(t *testing.T) {
	extension := Extension{
		MakerAssetSuffix: "0x",
		TakerAssetSuffix: "0x",
		MakingAmountData: "0x2ad5004c60e16e54d5007c80ce329adde5b51ef5000000000000006859e6260000b401bf92000000000000640ac0866635457d36ab318d0000000000000000000066593d4e7d3a5f55167fd18bd45f0b94f54a968f000000000000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d972ad4499f120902631a95770895ad27ad6b0d95",
		TakingAmountData: "0x2ad5004c60e16e54d5007c80ce329adde5b51ef5000000000000006859e6260000b401bf92000000000000640ac0866635457d36ab318d0000000000000000000066593d4e7d3a5f55167fd18bd45f0b94f54a968f000000000000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d972ad4499f120902631a95770895ad27ad6b0d95",
		Predicate:        "0x",
		MakerPermit:      "0x",
		PreInteraction:   "0x",
		PostInteraction:  "0x2ad5004c60e16e54d5007c80ce329adde5b51ef500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000646859e6150ac0866635457d36ab318d000000000000000000000000000066593d4e7d3a5f55167f0000d18bd45f0b94f54a968f0000000000000000000000000000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d0000972ad4499f120902631a000095770895ad27ad6b0d9500000000000000000000000000000000000000000000000000000000000000075dec5a",
	}

	tests := []struct {
		name        string
		baseSalt    *big.Int
		expectedLow string // Last 40 chars (160 bits) of the salt
	}{
		{
			name:        "No base salt - verify extension hash portion",
			baseSalt:    nil,
			expectedLow: "743b07ed0eae652cb39033bb6e4a3c7fa8662b5c",
		},
		{
			name:        "With base salt 0",
			baseSalt:    big.NewInt(0),
			expectedLow: "743b07ed0eae652cb39033bb6e4a3c7fa8662b5c",
		},
		{
			name:        "With base salt 11111111",
			baseSalt:    big.NewInt(11111111),
			expectedLow: "743b07ed0eae652cb39033bb6e4a3c7fa8662b5c",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := extension.Encode()
			require.NoError(t, err)

			salt, err := GenerateSalt(encoded, tc.baseSalt)
			require.NoError(t, err)

			// Verify the last 40 characters match (extension hash portion)
			saltLow := salt[len(salt)-40:]
			assert.Equal(t, tc.expectedLow, saltLow, "Salt extension hash portion mismatch")
		})
	}
}

// TestExtensionRoundTrip_KnownValues verifies extension encoding and decoding
func TestExtensionRoundTrip_KnownValues(t *testing.T) {
	tests := []struct {
		name      string
		extension Extension
	}{
		{
			name: "Simple extension",
			extension: Extension{
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
		{
			name: "Empty extension",
			extension: Extension{
				MakerAssetSuffix: "0x",
				TakerAssetSuffix: "0x",
				MakingAmountData: "0x",
				TakingAmountData: "0x",
				Predicate:        "0x",
				MakerPermit:      "0x",
				PreInteraction:   "0x",
				PostInteraction:  "0x",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Encode
			encoded, err := tc.extension.Encode()
			require.NoError(t, err)
			require.True(t, len(encoded) >= 2 && encoded[:2] == "0x", "Encoded extension should start with 0x")

			// Decode and verify round-trip
			decoded, err := Decode(mustDecodeHex(encoded))
			require.NoError(t, err)

			assert.Equal(t, tc.extension.MakerAssetSuffix, decoded.MakerAssetSuffix)
			assert.Equal(t, tc.extension.TakerAssetSuffix, decoded.TakerAssetSuffix)
			assert.Equal(t, tc.extension.MakingAmountData, decoded.MakingAmountData)
			assert.Equal(t, tc.extension.TakingAmountData, decoded.TakingAmountData)
			assert.Equal(t, tc.extension.Predicate, decoded.Predicate)
			assert.Equal(t, tc.extension.MakerPermit, decoded.MakerPermit)
			assert.Equal(t, tc.extension.PreInteraction, decoded.PreInteraction)
			assert.Equal(t, tc.extension.PostInteraction, decoded.PostInteraction)
		})
	}
}

// TestBitmaskOperations_KnownValues verifies bitmask operations
func TestBitmaskOperations_KnownValues(t *testing.T) {
	tests := []struct {
		name           string
		startBit       int64
		endBit         int64
		expectedString string
	}{
		{
			name:           "Simple mask bits 4-8",
			startBit:       4,
			endBit:         8,
			expectedString: "0xf0",
		},
		{
			name:           "Single bit mask",
			startBit:       0,
			endBit:         1,
			expectedString: "0x1",
		},
		{
			name:           "Full byte mask",
			startBit:       0,
			endBit:         8,
			expectedString: "0xff",
		},
		{
			name:           "Two byte mask",
			startBit:       0,
			endBit:         16,
			expectedString: "0xffff",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bitmask, err := NewBitMask(big.NewInt(tc.startBit), big.NewInt(tc.endBit))
			require.NoError(t, err)
			assert.Equal(t, tc.expectedString, bitmask.ToString(), "Bitmask string mismatch")
		})
	}
}
