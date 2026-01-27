package orderbook

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewExtension(t *testing.T) {
	tests := []struct {
		name        string
		params      ExtensionParams
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid - empty extension",
			params: ExtensionParams{
				MakerAssetData:  "",
				TakerAssetData:  "",
				GetMakingAmount: "",
				GetTakingAmount: "",
				Predicate:       "",
				Permit:          "",
				PreInteraction:  "",
				PostInteraction: "",
			},
			expectError: false,
		},
		{
			name: "Valid - with all fields",
			params: ExtensionParams{
				MakerAsset:      "0x6B175474E89094C44Da98b954EedeAC495271d0F",
				MakerAssetData:  "0x1234",
				TakerAssetData:  "0x5678",
				GetMakingAmount: "0xabcd",
				GetTakingAmount: "0xef01",
				Predicate:       "0x2345",
				Permit:          "0x6789",
				PreInteraction:  "0xabcd",
				PostInteraction: "0xef01",
			},
			expectError: false,
		},
		{
			name: "Invalid - Permit without MakerAsset",
			params: ExtensionParams{
				Permit: "0x1234",
			},
			expectError: true,
			errorMsg:    "MakerAsset",
		},
		{
			name: "Invalid - MakerAsset without Permit",
			params: ExtensionParams{
				MakerAsset: "0x6B175474E89094C44Da98b954EedeAC495271d0F",
			},
			expectError: true,
			errorMsg:    "Permit",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := NewExtension(tc.params)
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}

func TestExtension_Encode(t *testing.T) {
	tests := []struct {
		name      string
		extension Extension
		expectErr bool
	}{
		{
			name: "Empty extension",
			extension: Extension{
				MakerAssetSuffix: "",
				TakerAssetSuffix: "",
				MakingAmountData: "",
				TakingAmountData: "",
				Predicate:        "",
				MakerPermit:      "",
				PreInteraction:   "",
				PostInteraction:  "",
			},
			expectErr: false,
		},
		{
			name: "Extension with 0x prefix fields",
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
			expectErr: false,
		},
		{
			name: "Extension without 0x prefix",
			extension: Extension{
				MakerAssetSuffix: "1234",
				TakerAssetSuffix: "5678",
				MakingAmountData: "",
				TakingAmountData: "",
				Predicate:        "",
				MakerPermit:      "",
				PreInteraction:   "",
				PostInteraction:  "",
			},
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.extension.Encode()
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				// Result should start with 0x
				assert.True(t, len(result) >= 2 && result[:2] == "0x")
			}
		})
	}
}

func TestExtension_EncodeDecodeRoundTrip(t *testing.T) {
	// Create an extension
	original := Extension{
		MakerAssetSuffix: "0x1234",
		TakerAssetSuffix: "0x5678",
		MakingAmountData: "0xabcdef",
		TakingAmountData: "0x123456",
		Predicate:        "0xaabbcc",
		MakerPermit:      "0xddeeff",
		PreInteraction:   "0x112233",
		PostInteraction:  "0x445566",
	}

	// Encode it
	encoded, err := original.Encode()
	require.NoError(t, err)
	require.NotEmpty(t, encoded)

	// Decode the encoded string (skip the 0x prefix)
	decoded, err := Decode(mustDecodeHex(encoded))
	require.NoError(t, err)

	// Compare fields
	assert.Equal(t, original.MakerAssetSuffix, decoded.MakerAssetSuffix)
	assert.Equal(t, original.TakerAssetSuffix, decoded.TakerAssetSuffix)
	assert.Equal(t, original.MakingAmountData, decoded.MakingAmountData)
	assert.Equal(t, original.TakingAmountData, decoded.TakingAmountData)
	assert.Equal(t, original.Predicate, decoded.Predicate)
	assert.Equal(t, original.MakerPermit, decoded.MakerPermit)
	assert.Equal(t, original.PreInteraction, decoded.PreInteraction)
	assert.Equal(t, original.PostInteraction, decoded.PostInteraction)
}

func TestExtension_EncodeDecodeRoundTrip_EmptyFields(t *testing.T) {
	// Create an extension with empty fields
	original := Extension{
		MakerAssetSuffix: "0x",
		TakerAssetSuffix: "0x",
		MakingAmountData: "0x",
		TakingAmountData: "0x",
		Predicate:        "0x",
		MakerPermit:      "0x",
		PreInteraction:   "0x",
		PostInteraction:  "0x",
	}

	// Encode it
	encoded, err := original.Encode()
	require.NoError(t, err)
	require.NotEmpty(t, encoded)

	// Decode the encoded string
	decoded, err := Decode(mustDecodeHex(encoded))
	require.NoError(t, err)

	// For empty fields, both should decode to "0x"
	assert.Equal(t, "0x", decoded.MakerAssetSuffix)
	assert.Equal(t, "0x", decoded.TakerAssetSuffix)
}

func TestDecode_InvalidData(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		errorMsg string
	}{
		{
			name:     "Empty data",
			data:     []byte{},
			errorMsg: "failed to read offsets",
		},
		{
			name:     "Data too short for offsets",
			data:     []byte{0x01, 0x02, 0x03},
			errorMsg: "failed to read offsets",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Decode(tc.data)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.errorMsg)
		})
	}
}

func TestExtension_EncodeDeterministic(t *testing.T) {
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

	// Encode multiple times
	encoded1, err := ext.Encode()
	require.NoError(t, err)

	encoded2, err := ext.Encode()
	require.NoError(t, err)

	encoded3, err := ext.Encode()
	require.NoError(t, err)

	// All encodings should be identical
	assert.Equal(t, encoded1, encoded2)
	assert.Equal(t, encoded2, encoded3)
}

// Helper function to decode hex string to bytes
func mustDecodeHex(s string) []byte {
	if len(s) >= 2 && s[:2] == "0x" {
		s = s[2:]
	}
	if len(s)%2 != 0 {
		s = "0" + s
	}
	b := make([]byte, len(s)/2)
	for i := 0; i < len(b); i++ {
		b[i] = hexCharToByte(s[2*i])<<4 | hexCharToByte(s[2*i+1])
	}
	return b
}

func hexCharToByte(c byte) byte {
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
