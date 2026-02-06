package fusionorder

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuctionDetailsEncoding_KnownValues verifies auction details encoding against known expected values
func TestAuctionDetailsEncoding_KnownValues(t *testing.T) {
	tests := []struct {
		name           string
		auctionDetails *AuctionDetails
		expectedEncode string
	}{
		{
			name: "Standard auction - no gas cost, single point",
			auctionDetails: &AuctionDetails{
				StartTime:       1673548149,
				Duration:        180,
				InitialRateBump: 50000,
				Points:          []AuctionPointClassFixed{{Coefficient: 20000, Delay: 12}},
				GasCost:         GasCostConfigClassFixed{GasBumpEstimate: 0, GasPriceEstimate: 0},
			},
			// Format: GasCost(7 bytes) + StartTime(4 bytes) + Duration(3 bytes) + InitialRateBump(3 bytes) + PointCount(1 byte) + Points
			// GasCost: 00 00 00 00 00 00 00 (zeros)
			// StartTime: 63c05175
			// Duration: 0000b4 (180)
			// InitialRateBump: 00c350 (50000)
			// PointCount: 01
			// Point 1: 004e20 (20000) + 000c (12)
			expectedEncode: "0000000000000063c051750000b400c35001004e20000c",
		},
		{
			name: "With gas cost",
			auctionDetails: &AuctionDetails{
				StartTime:       1673548149,
				Duration:        180,
				InitialRateBump: 50000,
				Points:          []AuctionPointClassFixed{{Coefficient: 20000, Delay: 12}},
				GasCost:         GasCostConfigClassFixed{GasBumpEstimate: 10000, GasPriceEstimate: 1000000},
			},
			// GasCost: 002710 (10000) + 000f4240 (1000000)
			expectedEncode: "002710000f424063c051750000b400c35001004e20000c",
		},
		{
			name: "Multiple points",
			auctionDetails: &AuctionDetails{
				StartTime:       1673548149,
				Duration:        180,
				InitialRateBump: 50000,
				Points: []AuctionPointClassFixed{
					{Coefficient: 10000, Delay: 10},
					{Coefficient: 5000, Delay: 20},
				},
				GasCost: GasCostConfigClassFixed{GasBumpEstimate: 0, GasPriceEstimate: 0},
			},
			// PointCount: 02
			// Point 1: 002710 (10000) + 000a (10)
			// Point 2: 001388 (5000) + 0014 (20)
			expectedEncode: "0000000000000063c051750000b400c35002002710000a0013880014",
		},
		{
			name: "No points",
			auctionDetails: &AuctionDetails{
				StartTime:       1673548149,
				Duration:        180,
				InitialRateBump: 50000,
				Points:          nil,
				GasCost:         GasCostConfigClassFixed{GasBumpEstimate: 0, GasPriceEstimate: 0},
			},
			// PointCount: 00
			expectedEncode: "0000000000000063c051750000b400c35000",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			encoded := tc.auctionDetails.Encode()
			assert.Equal(t, tc.expectedEncode, encoded, "AuctionDetails encoding mismatch")

			// Verify determinism
			encoded2 := tc.auctionDetails.Encode()
			assert.Equal(t, encoded, encoded2, "AuctionDetails encoding should be deterministic")
		})
	}
}

// TestAuctionDetailsEncoding_WithoutPointCount_KnownValues verifies encoding without point count
func TestAuctionDetailsEncoding_WithoutPointCount_KnownValues(t *testing.T) {
	tests := []struct {
		name           string
		auctionDetails *AuctionDetails
		expectedEncode string
	}{
		{
			name: "FusionPlus style - single point",
			auctionDetails: &AuctionDetails{
				StartTime:       1673548149,
				Duration:        180,
				InitialRateBump: 50000,
				Points:          []AuctionPointClassFixed{{Coefficient: 20000, Delay: 12}},
				GasCost:         GasCostConfigClassFixed{GasBumpEstimate: 0, GasPriceEstimate: 0},
			},
			// Same as Encode but without point count byte
			expectedEncode: "0000000000000063c051750000b400c350004e20000c",
		},
		{
			name: "FusionPlus style - multiple points",
			auctionDetails: &AuctionDetails{
				StartTime:       1673548149,
				Duration:        180,
				InitialRateBump: 50000,
				Points: []AuctionPointClassFixed{
					{Coefficient: 10000, Delay: 10},
					{Coefficient: 5000, Delay: 20},
				},
				GasCost: GasCostConfigClassFixed{GasBumpEstimate: 0, GasPriceEstimate: 0},
			},
			expectedEncode: "0000000000000063c051750000b400c350002710000a0013880014",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			encoded := tc.auctionDetails.EncodeWithoutPointCount()
			assert.Equal(t, tc.expectedEncode, encoded, "AuctionDetails encoding mismatch")
		})
	}
}

// TestInteractionEncoding_KnownValues verifies interaction encoding against known expected values
func TestInteractionEncoding_KnownValues(t *testing.T) {
	tests := []struct {
		name           string
		target         common.Address
		data           string
		expectedEncode string
	}{
		{
			name:           "Standard interaction",
			target:         common.HexToAddress("0x1234567890123456789012345678901234567890"),
			data:           "0xabcdef",
			expectedEncode: "0x1234567890123456789012345678901234567890abcdef",
		},
		{
			name:           "Interaction with empty data",
			target:         common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
			data:           "0x",
			expectedEncode: "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
		},
		{
			name:           "Interaction with long data",
			target:         common.HexToAddress("0x0000000000000000000000000000000000000001"),
			data:           "0x1234567890abcdef1234567890abcdef",
			expectedEncode: "0x00000000000000000000000000000000000000011234567890abcdef1234567890abcdef",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			interaction, err := NewInteraction(tc.target, tc.data)
			require.NoError(t, err)

			encoded := interaction.Encode()
			assert.Equal(t, tc.expectedEncode, encoded, "Interaction encoding mismatch")

			// Verify round-trip
			decoded, err := DecodeInteraction(encoded)
			require.NoError(t, err)
			assert.Equal(t, interaction.Target, decoded.Target)
			assert.Equal(t, interaction.Data, decoded.Data)
		})
	}
}

// TestBpsToRatioFormat_KnownValues verifies BPS to ratio format conversion
// The function multiplies BPS by 10 (feeBase/bpsBase = 100000/10000 = 10)
func TestBpsToRatioFormat_KnownValues(t *testing.T) {
	tests := []struct {
		name     string
		bps      *big.Int
		expected *big.Int
	}{
		{
			name:     "Nil BPS",
			bps:      nil,
			expected: big.NewInt(0),
		},
		{
			name:     "Zero BPS",
			bps:      big.NewInt(0),
			expected: big.NewInt(0),
		},
		{
			name:     "100 BPS (1%)",
			bps:      big.NewInt(100),
			expected: big.NewInt(1000), // 100 * 10 = 1000
		},
		{
			name:     "1000 BPS (10%)",
			bps:      big.NewInt(1000),
			expected: big.NewInt(10000), // 1000 * 10 = 10000
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := BpsToRatioFormat(tc.bps)
			assert.Equal(t, 0, tc.expected.Cmp(result), "BpsToRatioFormat mismatch: expected %s, got %s",
				tc.expected.String(), result.String())
		})
	}
}

// TestWhitelistItemCreation_KnownValues verifies whitelist item creation
func TestWhitelistItemCreation_KnownValues(t *testing.T) {
	tests := []struct {
		name                string
		addressHalf         string
		delay               *big.Int
		expectedAddressHalf string
		expectedDelay       *big.Int
	}{
		{
			name:                "Standard whitelist item",
			addressHalf:         "bb839cbe05303d7705fa",
			delay:               big.NewInt(0),
			expectedAddressHalf: "bb839cbe05303d7705fa",
			expectedDelay:       big.NewInt(0),
		},
		{
			name:                "Whitelist item with delay",
			addressHalf:         "1234567890abcdef1234",
			delay:               big.NewInt(100),
			expectedAddressHalf: "1234567890abcdef1234",
			expectedDelay:       big.NewInt(100),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			item := NewWhitelistItem(tc.addressHalf, tc.delay)
			assert.Equal(t, tc.expectedAddressHalf, item.AddressHalf)
			assert.Equal(t, 0, tc.expectedDelay.Cmp(item.Delay), "Delay mismatch")
		})
	}
}

// TestMakerTraitsCreation_KnownValues verifies maker traits creation via shared function
func TestMakerTraitsCreation_KnownValues(t *testing.T) {
	tests := []struct {
		name           string
		params         MakerTraitsParams
		expectedEncode string
	}{
		{
			name: "Standard order - partial and multiple fills allowed",
			params: MakerTraitsParams{
				AuctionStartTime:     1673548149,
				AuctionDuration:      180,
				OrderExpirationDelay: 12,
				Nonce:                nil,
				AllowPartialFills:    true,
				AllowMultipleFills:   true,
				UnwrapWeth:           false,
				EnablePermit2:        false,
			},
			// Deadline = 1673548149 + 180 + 12 = 1673548341 = 0x63c05235
			expectedEncode: "0x4a000000000000000000000000000000000063c0523500000000000000000000",
		},
		{
			name: "Order with nonce",
			params: MakerTraitsParams{
				AuctionStartTime:     1673548149,
				AuctionDuration:      180,
				OrderExpirationDelay: 12,
				Nonce:                big.NewInt(12345),
				AllowPartialFills:    false,
				AllowMultipleFills:   false,
				UnwrapWeth:           false,
				EnablePermit2:        false,
			},
			expectedEncode: "0x8a000000000000000000000000000030390063c0523500000000000000000000",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			makerTraits, err := CreateMakerTraits(tc.params)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedEncode, makerTraits.Encode(), "MakerTraits encoding mismatch")
		})
	}
}

// TestBpsConversion_KnownValues verifies basis points conversions
func TestBpsConversion_KnownValues(t *testing.T) {
	defaultBase := GetDefaultBase()

	tests := []struct {
		name        string
		percent     float64
		base        *big.Int
		expectedBps *big.Int
	}{
		{
			name:        "1% with default base",
			percent:     1,
			base:        defaultBase,
			expectedBps: big.NewInt(100),
		},
		{
			name:        "0.5% with default base",
			percent:     0.5,
			base:        defaultBase,
			expectedBps: big.NewInt(50),
		},
		{
			name:        "10% with default base",
			percent:     10,
			base:        defaultBase,
			expectedBps: big.NewInt(1000),
		},
		{
			name:        "0% with default base",
			percent:     0,
			base:        defaultBase,
			expectedBps: big.NewInt(0),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bps := MustFromPercent(tc.percent, tc.base)
			assert.Equal(t, 0, tc.expectedBps.Cmp(bps.Value()), "Bps value mismatch")
		})
	}
}

// TestWhitelistGeneration_KnownValues verifies whitelist address half extraction
func TestWhitelistGeneration_KnownValues(t *testing.T) {
	tests := []struct {
		name                string
		addresses           []string
		resolvingStartTime  *big.Int
		expectedAddressHalf string
	}{
		{
			name:                "ETH2 deposit contract",
			addresses:           []string{"0x00000000219ab540356cbb839cbe05303d7705fa"},
			resolvingStartTime:  big.NewInt(1708117482),
			expectedAddressHalf: "bb839cbe05303d7705fa",
		},
		{
			name:                "Simple incremental address",
			addresses:           []string{"0x1234567890123456789012345678901234567890"},
			resolvingStartTime:  big.NewInt(1000000),
			expectedAddressHalf: "12345678901234567890",
		},
		{
			name:                "All F's address",
			addresses:           []string{"0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"},
			resolvingStartTime:  big.NewInt(1000000),
			expectedAddressHalf: "ffffffffffffffffffff",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			whitelist, err := GenerateWhitelist(tc.addresses, tc.resolvingStartTime)
			require.NoError(t, err)
			require.Len(t, whitelist, 1)

			// Verify the address half matches the last 20 characters of the address (lowercase)
			expectedHalf := strings.ToLower(tc.addresses[0][len(tc.addresses[0])-20:])
			assert.Equal(t, expectedHalf, whitelist[0].AddressHalf, "Address half mismatch")

			// First item should always have delay of 0
			assert.Equal(t, big.NewInt(0), whitelist[0].Delay, "First item delay should be 0")
		})
	}
}

// TestWhitelistGeneration_Validation verifies whitelist generation produces correct structure
func TestWhitelistGeneration_Validation(t *testing.T) {
	tests := []struct {
		name               string
		addresses          []string
		resolvingStartTime *big.Int
		expectedLen        int
	}{
		{
			name:               "Single address",
			addresses:          []string{"0x00000000219ab540356cbb839cbe05303d7705fa"},
			resolvingStartTime: big.NewInt(1708117482),
			expectedLen:        1,
		},
		{
			name: "Multiple addresses",
			addresses: []string{
				"0x00000000219ab540356cbb839cbe05303d7705fa",
				"0x1234567890123456789012345678901234567890",
			},
			resolvingStartTime: big.NewInt(1708117482),
			expectedLen:        2,
		},
		{
			name:               "Checksummed address",
			addresses:          []string{"0x6B175474E89094C44Da98b954EedeAC495271d0F"},
			resolvingStartTime: big.NewInt(1708117482),
			expectedLen:        1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Generate whitelist twice to test determinism
			whitelist1, err := GenerateWhitelist(tc.addresses, tc.resolvingStartTime)
			require.NoError(t, err)
			require.Len(t, whitelist1, tc.expectedLen, "Whitelist should have expected length")

			whitelist2, err := GenerateWhitelist(tc.addresses, tc.resolvingStartTime)
			require.NoError(t, err)

			// Verify determinism
			for i := range whitelist1 {
				assert.Equal(t, whitelist1[i].AddressHalf, whitelist2[i].AddressHalf, "Whitelist generation should be deterministic")
				assert.Equal(t, 0, whitelist1[i].Delay.Cmp(whitelist2[i].Delay), "Delay should be deterministic")
			}

			// Verify structure
			assert.Equal(t, big.NewInt(0), whitelist1[0].Delay, "First whitelist item should have 0 delay")

			for i, item := range whitelist1 {
				assert.Len(t, item.AddressHalf, 20, "Address half should be 20 hex chars (10 bytes)")
				assert.Equal(t, strings.ToLower(item.AddressHalf), item.AddressHalf, "Address half should be lowercase")
				assert.NotNil(t, item.Delay, "Delay should not be nil for item %d", i)
			}
		})
	}
}

// TestWhitelistDelayCalculation_KnownValues verifies delay calculations
func TestWhitelistDelayCalculation_KnownValues(t *testing.T) {
	tests := []struct {
		name           string
		items          []AuctionWhitelistItem
		resolvingStart *big.Int
		expectedDelays []*big.Int
	}{
		{
			name: "Three addresses with increasing delays",
			items: []AuctionWhitelistItem{
				{Address: common.HexToAddress("0x1111111111111111111111111111111111111111"), AllowFrom: big.NewInt(1000000)},
				{Address: common.HexToAddress("0x2222222222222222222222222222222222222222"), AllowFrom: big.NewInt(1000100)},
				{Address: common.HexToAddress("0x3333333333333333333333333333333333333333"), AllowFrom: big.NewInt(1000200)},
			},
			resolvingStart: big.NewInt(1000000),
			expectedDelays: []*big.Int{big.NewInt(0), big.NewInt(100), big.NewInt(100)},
		},
		{
			name: "Two addresses with same time",
			items: []AuctionWhitelistItem{
				{Address: common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), AllowFrom: big.NewInt(1000000)},
				{Address: common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"), AllowFrom: big.NewInt(1000000)},
			},
			resolvingStart: big.NewInt(1000000),
			expectedDelays: []*big.Int{big.NewInt(0), big.NewInt(0)},
		},
		{
			name: "Single resolver",
			items: []AuctionWhitelistItem{
				{Address: common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"), AllowFrom: big.NewInt(1000000)},
			},
			resolvingStart: big.NewInt(1000000),
			expectedDelays: []*big.Int{big.NewInt(0)},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			whitelist, err := GenerateWhitelistFromItems(tc.items, tc.resolvingStart)
			require.NoError(t, err)
			require.Len(t, whitelist, len(tc.expectedDelays))

			for i, expected := range tc.expectedDelays {
				assert.Equal(t, 0, expected.Cmp(whitelist[i].Delay),
					"Delay mismatch at index %d: expected %s, got %s", i, expected, whitelist[i].Delay)
			}
		})
	}
}
