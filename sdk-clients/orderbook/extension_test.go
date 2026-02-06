package orderbook

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name             string
		extension        Extension
		expectedEncoding string
	}{
		{
			name: "Simple Limit Order 1",
			extension: Extension{
				MakerAssetSuffix: "0x01",
				TakerAssetSuffix: "0x02",
				MakingAmountData: "0x03",
				TakingAmountData: "0x04",
				Predicate:        "0x05",
				MakerPermit:      "0x06",
				PreInteraction:   "0x07",
				PostInteraction:  "0x08",
			},
			expectedEncoding: "0x00000008000000070000000600000005000000040000000300000002000000010102030405060708",
		},
		{
			name: "Realistic Order 1",
			extension: Extension{
				MakerAssetSuffix: "0x",
				TakerAssetSuffix: "0x",
				MakingAmountData: "0xfb2809A5314473E1165f6B58018E20ed8F07B84000000000000000666cdf850000b400c45c00688b007e",
				TakingAmountData: "0xfb2809A5314473E1165f6B58018E20ed8F07B84000000000000000666cdf850000b400c45c00688b007e",
				Predicate:        "0x",
				MakerPermit:      "0x",
				PreInteraction:   "0x",
				PostInteraction:  "0xfb2809A5314473E1165f6B58018E20ed8F07B840666cdf74c0866635457d36ab318d0000f3a44b7b0d08f4e198b80000000000000000000000000000d18bd45f0b94f54a968f0000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d000000000000000000000000000040",
			},
			expectedEncoding: "0x000000cd000000540000005400000054000000540000002a0000000000000000fb2809A5314473E1165f6B58018E20ed8F07B84000000000000000666cdf850000b400c45c00688b007efb2809A5314473E1165f6B58018E20ed8F07B84000000000000000666cdf850000b400c45c00688b007efb2809A5314473E1165f6B58018E20ed8F07B840666cdf74c0866635457d36ab318d0000f3a44b7b0d08f4e198b80000000000000000000000000000d18bd45f0b94f54a968f0000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d000000000000000000000000000040",
		},
		{
			name: "Realistic Order 2",
			extension: Extension{
				MakerAssetSuffix: "0x",
				TakerAssetSuffix: "0x",
				MakingAmountData: "0xfb2809A5314473E1165f6B58018E20ed8F07B8400000000000000067217a910000b401a70b",
				TakingAmountData: "0xfb2809A5314473E1165f6B58018E20ed8F07B8400000000000000067217a910000b401a70b",
				Predicate:        "0x",
				MakerPermit:      "0x",
				PreInteraction:   "0x",
				PostInteraction:  "0xfb2809A5314473E1165f6B58018E20ed8F07B84067217a80c0866635457d36ab318d00002385c09fca8e96142deb0000000000000000000000000000d18bd45f0b94f54a968f0000000000000000000000000000ade19567bb538035ed360000fff14e3bbfd616d555650000000000000000000000000000000000000000000000000000f3a44b7b0d08f4e198b80000c976bf098c4dba0a061d000000000000000000000000000060",
			},
			expectedEncoding: "0x000000f30000004a0000004a0000004a0000004a000000250000000000000000fb2809A5314473E1165f6B58018E20ed8F07B8400000000000000067217a910000b401a70bfb2809A5314473E1165f6B58018E20ed8F07B8400000000000000067217a910000b401a70bfb2809A5314473E1165f6B58018E20ed8F07B84067217a80c0866635457d36ab318d00002385c09fca8e96142deb0000000000000000000000000000d18bd45f0b94f54a968f0000000000000000000000000000ade19567bb538035ed360000fff14e3bbfd616d555650000000000000000000000000000000000000000000000000000f3a44b7b0d08f4e198b80000c976bf098c4dba0a061d000000000000000000000000000060",
		},
		{
			name: "Simple Limit Order 2",
			extension: Extension{
				MakerAssetSuffix: "0x01",
				TakerAssetSuffix: "0x0222",
				MakingAmountData: "0x033344",
				TakingAmountData: "0x04",
				Predicate:        "0x05",
				MakerPermit:      "0x06",
				PreInteraction:   "0x075533",
				PostInteraction:  "0x4A7F9C3B2D8E1F5A6B4C3D2E1F0A9B8C7D6E5F4A3B2C1D0E9F8A7B6C5D4E3F2A1B0C9D8E7F6A5B4C3D2E1F0A9B8C7D6E5F4A3B2C1D0E9F8A7B6C5D4E3F2A1B0C9D8E7F6A5B4C3D2E1F0A9B8C7D6E5F4A3B2C1D0E9F8A7B6C5D4E3F2A1B0C9D8E7F6A5B4C3D2E1F0A9B8C7D6E5F4A3B2C1D0E9F8A7B6C5D4E3F2A1B0C9D8E7F6A5B4C3D2E1F0A9B8C7D6E5F4A3B2C1D0E9F8A7B6C5D4E3F2A1B0C9D8E7F6A5B4C3D2E1F0A9B8C7D6E5F4A3B2C1D0E9F8A7B6C5D4E3F2A1B0C9D8E7F6A5B4C3D2E1F0A9B8C7D6E5F4A3B2C1D0E9F87",
			},
			expectedEncoding: "0x000000da0000000c0000000900000008000000070000000600000003000000010102220333440405060755334A7F9C3B2D8E1F5A6B4C3D2E1F0A9B8C7D6E5F4A3B2C1D0E9F8A7B6C5D4E3F2A1B0C9D8E7F6A5B4C3D2E1F0A9B8C7D6E5F4A3B2C1D0E9F8A7B6C5D4E3F2A1B0C9D8E7F6A5B4C3D2E1F0A9B8C7D6E5F4A3B2C1D0E9F8A7B6C5D4E3F2A1B0C9D8E7F6A5B4C3D2E1F0A9B8C7D6E5F4A3B2C1D0E9F8A7B6C5D4E3F2A1B0C9D8E7F6A5B4C3D2E1F0A9B8C7D6E5F4A3B2C1D0E9F8A7B6C5D4E3F2A1B0C9D8E7F6A5B4C3D2E1F0A9B8C7D6E5F4A3B2C1D0E9F8A7B6C5D4E3F2A1B0C9D8E7F6A5B4C3D2E1F0A9B8C7D6E5F4A3B2C1D0E9F87",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.extension.Encode()
			require.NoError(t, err)
			assert.Equal(t, strings.ToLower(tc.expectedEncoding), strings.ToLower(result))
		})
	}
}

// TestDecodeExtension contains all unit tests for the DecodeEscrowExtension function.
func TestDecodeExtension(t *testing.T) {
	tests := []struct {
		name          string
		hexInput      string
		expected      *Extension
		expectingErr  bool
		errorContains string
	}{
		{
			name:     "Successful Decoding",
			hexInput: "00000008000000070000000600000005000000040000000300000002000000010102050604030708",
			expected: &Extension{
				MakerAssetSuffix: "0x01",
				TakerAssetSuffix: "0x02",
				MakingAmountData: "0x05",
				TakingAmountData: "0x06",
				Predicate:        "0x04",
				MakerPermit:      "0x03",
				PreInteraction:   "0x07",
				PostInteraction:  "0x08",
			},
			expectingErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert hex string to bytes
			data, err := hex.DecodeString(tt.hexInput)
			if err != nil {
				t.Fatalf("Failed to convert hex to bytes: %v", err)
			}

			// Decode the data
			decoded, err := Decode(data)

			if tt.expectingErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s' but got '%s'", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !extensionsEqual(decoded, tt.expected) {
					t.Errorf("Decoded Extension does not match expected.\nGot: %+v\nExpected: %+v", decoded, tt.expected)
				}
			}
		})
	}
}

func extensionsEqual(a, b *Extension) bool {
	return a.MakerAssetSuffix == b.MakerAssetSuffix &&
		a.TakerAssetSuffix == b.TakerAssetSuffix &&
		a.MakingAmountData == b.MakingAmountData &&
		a.TakingAmountData == b.TakingAmountData &&
		a.Predicate == b.Predicate &&
		a.MakerPermit == b.MakerPermit &&
		a.PreInteraction == b.PreInteraction &&
		a.PostInteraction == b.PostInteraction
	// a.CustomData == b.CustomData
}
