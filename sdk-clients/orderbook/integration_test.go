package orderbook

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOrderDataSerialization tests order data JSON serialization
func TestOrderDataSerialization(t *testing.T) {
	order := OrderData{
		MakerAsset:    "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		TakerAsset:    "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		MakingAmount:  "1000000000000000000",
		TakingAmount:  "1420000000",
		Salt:          "12345678901234567890",
		Maker:         "0x1234567890123456789012345678901234567890",
		AllowedSender: "0x0000000000000000000000000000000000000000",
		Receiver:      "0x0000000000000000000000000000000000000000",
		MakerTraits:   "0x4a000000000000000000000000000000000063c0523500000000000000000000",
		Extension:     "0x1234",
	}

	// Serialize
	jsonData, err := json.Marshal(order)
	require.NoError(t, err)

	// Deserialize
	var deserialized OrderData
	err = json.Unmarshal(jsonData, &deserialized)
	require.NoError(t, err)

	// Verify
	assert.Equal(t, order.MakerAsset, deserialized.MakerAsset)
	assert.Equal(t, order.TakerAsset, deserialized.TakerAsset)
	assert.Equal(t, order.MakingAmount, deserialized.MakingAmount)
	assert.Equal(t, order.TakingAmount, deserialized.TakingAmount)
	assert.Equal(t, order.Salt, deserialized.Salt)
	assert.Equal(t, order.MakerTraits, deserialized.MakerTraits)
}

// TestOrderResponseSerialization tests order response deserialization
func TestOrderResponseSerialization(t *testing.T) {
	apiResponse := `{
		"orderHash": "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		"signature": "0x1234567890abcdef",
		"data": {
			"makerAsset": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			"takerAsset": "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
			"makingAmount": "1000000000000000000",
			"takingAmount": "1420000000",
			"salt": "12345678901234567890",
			"maker": "0x1234567890123456789012345678901234567890",
			"allowedSender": "0x0000000000000000000000000000000000000000",
			"receiver": "0x0000000000000000000000000000000000000000",
			"makerTraits": "0x4a000000000000000000000000000000000063c0523500000000000000000000",
			"extension": "0x"
		}
	}`

	var response OrderResponse
	err := json.Unmarshal([]byte(apiResponse), &response)
	require.NoError(t, err)

	assert.Equal(t, "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890", response.OrderHash)
	assert.Equal(t, "0x1234567890abcdef", response.Signature)
	assert.Equal(t, "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", response.Data.MakerAsset)
}

// TestMakerTraitsCreation_Integration tests maker traits creation and encoding
func TestMakerTraitsCreation_Integration(t *testing.T) {
	tests := []struct {
		name              string
		params            MakerTraitsParams
		expectError       bool
		expectPartialFill bool
	}{
		{
			name: "Both fills allowed",
			params: MakerTraitsParams{
				Expiry:             1700000000,
				AllowPartialFills:  true,
				AllowMultipleFills: true,
				HasPostInteraction: true,
				HasExtension:       true,
				Nonce:              0,
			},
			expectError:       false,
			expectPartialFill: true,
		},
		{
			name: "Both fills disabled",
			params: MakerTraitsParams{
				Expiry:             1700000000,
				AllowPartialFills:  false,
				AllowMultipleFills: false,
				HasPostInteraction: true,
				HasExtension:       true,
				Nonce:              12345,
			},
			expectError:       false,
			expectPartialFill: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			traits, err := NewMakerTraits(tc.params)
			if tc.expectError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, traits)
			assert.Equal(t, tc.expectPartialFill, traits.AllowPartialFills)

			// Test encoding
			encoded := traits.Encode()
			assert.NotEmpty(t, encoded)
			assert.True(t, len(encoded) > 2 && encoded[:2] == "0x")
		})
	}
}

// TestExtensionEncodeDecode_Integration tests extension encoding round-trip
func TestExtensionEncodeDecode_Integration(t *testing.T) {
	ext := Extension{
		MakerAssetSuffix: "0x1234",
		TakerAssetSuffix: "0x5678",
		MakingAmountData: "0xabcd",
		TakingAmountData: "0xef01",
		Predicate:        "0x2345",
		MakerPermit:      "0x6789",
		PreInteraction:   "0xabcd",
		PostInteraction:  "0xef01",
	}

	// Encode
	encoded, err := ext.Encode()
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)
	assert.True(t, len(encoded) >= 2 && encoded[:2] == "0x")

	// Decode
	decoded, err := Decode(mustDecodeHexLocal(encoded))
	require.NoError(t, err)

	// Verify round-trip
	assert.Equal(t, ext.MakerAssetSuffix, decoded.MakerAssetSuffix)
	assert.Equal(t, ext.TakerAssetSuffix, decoded.TakerAssetSuffix)
}

// TestSaltGeneration_Integration tests salt generation
func TestSaltGeneration_Integration(t *testing.T) {
	tests := []struct {
		name           string
		extension      string
		customBaseSalt *big.Int
		expectEmpty    bool
	}{
		{
			name:        "Empty extension - timestamp salt",
			extension:   "",
			expectEmpty: false,
		},
		{
			name:        "0x extension - timestamp salt",
			extension:   "0x",
			expectEmpty: false,
		},
		{
			name:        "Non-empty extension - hash-based salt",
			extension:   "0x1234567890abcdef",
			expectEmpty: false,
		},
		{
			name:           "Custom base salt",
			extension:      "0x1234567890abcdef",
			customBaseSalt: big.NewInt(999),
			expectEmpty:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			salt, err := GenerateSalt(tc.extension, tc.customBaseSalt)
			require.NoError(t, err)
			if !tc.expectEmpty {
				assert.NotEmpty(t, salt)
			}
		})
	}
}

// TestSaltGeneration_Deterministic tests that salt generation with custom base is deterministic
func TestSaltGeneration_Deterministic(t *testing.T) {
	extension := "0x1234567890abcdef"
	customBaseSalt := big.NewInt(12345)

	salt1, err := GenerateSalt(extension, customBaseSalt)
	require.NoError(t, err)

	salt2, err := GenerateSalt(extension, customBaseSalt)
	require.NoError(t, err)

	assert.Equal(t, salt1, salt2)
}

// TestTakerTraitsCreation tests taker traits creation
func TestTakerTraitsCreation(t *testing.T) {
	traits := NewTakerTraits(TakerTraitsParams{
		Extension:       "0x",
		MakerAmount:     false,
		UnwrapWETH:      false,
		SkipOrderPermit: false,
		UsePermit2:      false,
		ArgsHasReceiver: false,
	})
	require.NotNil(t, traits)

	// Verify the traits are set correctly
	assert.False(t, traits.MakerAmount)
	assert.False(t, traits.UnwrapWETH)
	assert.False(t, traits.SkipOrderPermit)
	assert.False(t, traits.UsePermit2)
}

// TestGetOrderByHashResponse_Serialization tests order lookup response
func TestGetOrderByHashResponse_Serialization(t *testing.T) {
	apiResponse := `{
		"orderHash": "0xabcdef1234567890",
		"createDateTime": "2024-01-01T00:00:00Z",
		"signature": "0x1234",
		"orderStatus": 1,
		"remainingMakerAmount": "1000000000000000000",
		"makerBalance": "2000000000000000000",
		"makerAllowance": "1000000000000000000",
		"data": {
			"makerAsset": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			"takerAsset": "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
			"makingAmount": "1000000000000000000",
			"takingAmount": "1420000000",
			"salt": "12345",
			"maker": "0x1234567890123456789012345678901234567890",
			"allowedSender": "0x0000000000000000000000000000000000000000",
			"receiver": "0x0000000000000000000000000000000000000000",
			"makerTraits": "0x4a",
			"extension": "0x"
		},
		"makerRate": "1.5",
		"takerRate": "0.67",
		"orderInvalidReason": "",
		"isMakerContract": false,
		"events": ""
	}`

	var response GetOrderByHashResponse
	err := json.Unmarshal([]byte(apiResponse), &response)
	require.NoError(t, err)

	assert.Equal(t, "0xabcdef1234567890", response.OrderHash)
	assert.Equal(t, "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", response.Data.MakerAsset)
}

// TestOrders_Serialization tests orders list response
func TestOrders_Serialization(t *testing.T) {
	apiResponse := `{
		"items": [
			{
				"orderHash": "0xabcdef1234567890",
				"signature": "0x1234",
				"data": {
					"makerAsset": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
					"takerAsset": "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
					"makingAmount": "1000000000000000000",
					"takingAmount": "1420000000",
					"salt": "12345",
					"maker": "0x1234567890123456789012345678901234567890",
					"allowedSender": "0x0000000000000000000000000000000000000000",
					"receiver": "0x0000000000000000000000000000000000000000",
					"makerTraits": "0x4a",
					"extension": "0x"
				}
			}
		]
	}`

	var orders Orders
	err := json.Unmarshal([]byte(apiResponse), &orders)
	require.NoError(t, err)

	require.Len(t, orders.Items, 1)
	assert.Equal(t, "0xabcdef1234567890", orders.Items[0].OrderHash)
}

// TestGetOrderCountResponse_Serialization tests order count response
func TestGetOrderCountResponse_Serialization(t *testing.T) {
	apiResponse := `{"count": 42}`

	var response GetOrderCountResponse
	err := json.Unmarshal([]byte(apiResponse), &response)
	require.NoError(t, err)

	assert.Equal(t, 42, response.Count)
}

// TestCreateOrderResponse_Serialization tests create order response
func TestCreateOrderResponse_Serialization(t *testing.T) {
	apiResponse := `{"success": true}`

	var response CreateOrderResponse
	err := json.Unmarshal([]byte(apiResponse), &response)
	require.NoError(t, err)

	assert.True(t, response.Success)
}

// TestBitmaskOperations tests bitmask operations used in traits
func TestBitmaskOperations(t *testing.T) {
	bm, err := NewBitMask(big.NewInt(0), big.NewInt(8))
	require.NoError(t, err)
	require.NotNil(t, bm)

	// BitMask stores offset and mask
	assert.NotNil(t, bm.Offset)
	assert.NotNil(t, bm.Mask)
}

// Helper function for hex decoding in tests
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
