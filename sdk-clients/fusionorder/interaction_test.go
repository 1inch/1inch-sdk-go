package fusionorder

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInteraction(t *testing.T) {
	tests := []struct {
		name        string
		target      common.Address
		data        string
		expectError bool
	}{
		{
			name:        "Valid interaction",
			target:      common.HexToAddress("0x1234567890123456789012345678901234567890"),
			data:        "0xabcdef",
			expectError: false,
		},
		{
			name:        "Valid interaction with empty data",
			target:      common.HexToAddress("0x1234567890123456789012345678901234567890"),
			data:        "0x",
			expectError: false,
		},
		{
			name:        "Invalid hex data",
			target:      common.HexToAddress("0x1234567890123456789012345678901234567890"),
			data:        "not-hex",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			interaction, err := NewInteraction(tc.target, tc.data)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.target, interaction.Target)
				assert.Equal(t, tc.data, interaction.Data)
			}
		})
	}
}

func TestInteraction_Encode(t *testing.T) {
	target := common.HexToAddress("0x1234567890123456789012345678901234567890")
	data := "0xabcdef"

	interaction, err := NewInteraction(target, data)
	require.NoError(t, err)

	encoded := interaction.Encode()
	// Should be lowercase target + data without 0x prefix
	assert.Equal(t, "0x1234567890123456789012345678901234567890abcdef", encoded)
}

func TestDecodeInteraction(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		checkTarget string
		checkData   string
	}{
		{
			name:        "Valid interaction",
			input:       "0x1234567890123456789012345678901234567890abcdef",
			expectError: false,
			checkTarget: "0x1234567890123456789012345678901234567890",
			checkData:   "0xabcdef",
		},
		{
			name:        "Invalid hex",
			input:       "not-hex-data",
			expectError: true,
		},
		{
			name:        "Too short",
			input:       "0x1234",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			interaction, err := DecodeInteraction(tc.input)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, common.HexToAddress(tc.checkTarget), interaction.Target)
				assert.Equal(t, tc.checkData, interaction.Data)
			}
		})
	}
}

func TestInteraction_EncodeDecodeRoundTrip(t *testing.T) {
	target := common.HexToAddress("0xABCDEF1234567890ABCDEF1234567890ABCDEF12")
	data := "0x1234567890"

	original, err := NewInteraction(target, data)
	require.NoError(t, err)

	encoded := original.Encode()
	decoded, err := DecodeInteraction(encoded)
	require.NoError(t, err)

	assert.Equal(t, original.Target, decoded.Target)
	assert.Equal(t, original.Data, decoded.Data)
}
