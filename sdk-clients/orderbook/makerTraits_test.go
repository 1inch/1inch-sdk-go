package orderbook

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakerTraitsEncode(t *testing.T) {

	tests := []struct {
		name                string
		makerTraitParams    MakerTraitsParams
		expectedMakerTraits string
	}{
		{
			name: "Extension, expiration",
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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			makerTraits, err := NewMakerTraits(tc.makerTraitParams)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedMakerTraits, makerTraits.Encode())
		})
	}
}

func TestDecodeMakerTraits(t *testing.T) {
	tests := []struct {
		name        string
		encoded     string
		expected    *MakerTraits
		expectError bool
		errorMsg    string
	}{
		{
			name:    "Extension, post interaction, multiple fills, expiration",
			encoded: "0x4a0000000000000000000000000000000000663be5db00000000000000000000",
			expected: &MakerTraits{
				Expiry:              1715201499,
				NeedPostinteraction: true,
				HasExtension:        true,
				AllowPartialFills:   true,
				AllowMultipleFills:  true,
			},
		},
		{
			name:    "With permit2 flag",
			encoded: "0x4b0000000000000000000000000000000000663be5db00000000000000000000",
			expected: &MakerTraits{
				Expiry:              1715201499,
				NeedPostinteraction: true,
				HasExtension:        true,
				ShouldUsePermit2:    true,
				AllowPartialFills:   true,
				AllowMultipleFills:  true,
			},
		},
		{
			name:        "Invalid hex",
			encoded:     "0xzz",
			expectError: true,
			errorMsg:    "invalid maker traits hex",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := DecodeMakerTraits(tc.encoded)
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestMakerTraitsEncodeDecodeRoundTrip(t *testing.T) {
	tests := []struct {
		name   string
		traits *MakerTraits
	}{
		{
			name: "Typical fusion order traits",
			traits: &MakerTraits{
				Expiry:              1715201499,
				NeedPostinteraction: true,
				HasExtension:        true,
				AllowPartialFills:   true,
				AllowMultipleFills:  true,
			},
		},
		{
			name: "Permit2 with nonce and unwrap",
			traits: &MakerTraits{
				Expiry:              1715201499,
				Nonce:               42,
				NeedPostinteraction: true,
				HasExtension:        true,
				ShouldUsePermit2:    true,
				ShouldUnwrapWeth:    true,
				AllowPartialFills:   true,
				AllowMultipleFills:  true,
			},
		},
		{
			name: "Bit invalidator mode with epoch check and pre interaction",
			traits: &MakerTraits{
				Expiry:              1715201499,
				Nonce:               7,
				Series:              3,
				NoPartialFills:      true,
				NeedPostinteraction: true,
				NeedPreinteraction:  true,
				NeedEpochCheck:      true,
				HasExtension:        true,
			},
		},
		{
			name: "With allowed sender tail",
			traits: &MakerTraits{
				AllowedSender:       "44839cbe05303d7705fa",
				Expiry:              1715201499,
				NeedPostinteraction: true,
				HasExtension:        true,
				AllowPartialFills:   true,
				AllowMultipleFills:  true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			decoded, err := DecodeMakerTraits(tc.traits.Encode())
			require.NoError(t, err)
			// The encoding stores only the low 80 bits of the sender, zero padded
			expected := *tc.traits
			if expected.AllowedSender != "" {
				expected.AllowedSender = fmt.Sprintf("%020s", expected.AllowedSender)
			}
			assert.Equal(t, &expected, decoded)
		})
	}
}

func TestDecodeMakerTraitsRejectsSignedHex(t *testing.T) {
	tests := []struct {
		name    string
		encoded string
	}{
		{name: "Negative hex", encoded: "-1"},
		{name: "Negative prefixed hex", encoded: "0x-1"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecodeMakerTraits(tc.encoded)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid maker traits hex")
		})
	}
}
