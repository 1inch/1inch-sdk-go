package fusionorder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomPresetValidate(t *testing.T) {
	tests := []struct {
		name        string
		preset      CustomPreset
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid preset without points",
			preset: CustomPreset{
				AuctionDuration:    180,
				AuctionStartAmount: "1000000000000000000",
				AuctionEndAmount:   "900000000000000000",
				Points:             nil,
			},
			expectError: false,
		},
		{
			name: "Valid preset with points",
			preset: CustomPreset{
				AuctionDuration:    180,
				AuctionStartAmount: "1000000000000000000",
				AuctionEndAmount:   "900000000000000000",
				Points: []CustomPresetPoint{
					{ToTokenAmount: "950000000000000000", Delay: 60},
				},
			},
			expectError: false,
		},
		{
			name: "Invalid auction start amount",
			preset: CustomPreset{
				AuctionDuration:    180,
				AuctionStartAmount: "not-a-number",
				AuctionEndAmount:   "900000000000000000",
			},
			expectError: true,
		errorMsg:    "invalid auction start amount",
	},
	{
		name: "Invalid auction end amount",
		preset: CustomPreset{
			AuctionDuration:    180,
			AuctionStartAmount: "1000000000000000000",
			AuctionEndAmount:   "not-a-number",
		},
		expectError: true,
		errorMsg:    "invalid auction end amount",
	},
	{
		name: "Zero duration",
		preset: CustomPreset{
			AuctionDuration:    0,
			AuctionStartAmount: "1000000000000000000",
			AuctionEndAmount:   "900000000000000000",
		},
		expectError: true,
		errorMsg:    "invalid auction duration: expected positive integer",
	},
	{
		name: "Negative duration",
		preset: CustomPreset{
			AuctionDuration:    -1,
			AuctionStartAmount: "1000000000000000000",
			AuctionEndAmount:   "900000000000000000",
		},
		expectError: true,
		errorMsg:    "invalid auction duration: expected positive integer",
	},
	{
		name: "Point amount above start amount",
		preset: CustomPreset{
			AuctionDuration:    180,
			AuctionStartAmount: "1000000000000000000",
			AuctionEndAmount:   "900000000000000000",
			Points: []CustomPresetPoint{
				{ToTokenAmount: "1100000000000000000", Delay: 60}, // Above start
			},
		},
		expectError: true,
		errorMsg:    "out of auction range",
	},
	{
		name: "Point amount below end amount",
		preset: CustomPreset{
			AuctionDuration:    180,
			AuctionStartAmount: "1000000000000000000",
			AuctionEndAmount:   "900000000000000000",
			Points: []CustomPresetPoint{
				{ToTokenAmount: "800000000000000000", Delay: 60}, // Below end
			},
		},
		expectError: true,
		errorMsg:    "out of auction range",
	},
	{
		name: "Invalid point amount",
		preset: CustomPreset{
			AuctionDuration:    180,
			AuctionStartAmount: "1000000000000000000",
			AuctionEndAmount:   "900000000000000000",
			Points: []CustomPresetPoint{
				{ToTokenAmount: "not-a-number", Delay: 60},
			},
		},
		expectError: true,
		errorMsg:    "invalid point amount",
	},
		{
			name: "Multiple valid points",
			preset: CustomPreset{
				AuctionDuration:    180,
				AuctionStartAmount: "1000000000000000000",
				AuctionEndAmount:   "900000000000000000",
				Points: []CustomPresetPoint{
					{ToTokenAmount: "980000000000000000", Delay: 30},
					{ToTokenAmount: "950000000000000000", Delay: 60},
					{ToTokenAmount: "920000000000000000", Delay: 90},
				},
			},
			expectError: false,
		},
		{
			name: "Point at exact start amount",
			preset: CustomPreset{
				AuctionDuration:    180,
				AuctionStartAmount: "1000000000000000000",
				AuctionEndAmount:   "900000000000000000",
				Points: []CustomPresetPoint{
					{ToTokenAmount: "1000000000000000000", Delay: 60}, // Exactly at start
				},
			},
			expectError: false,
		},
		{
			name: "Point at exact end amount",
			preset: CustomPreset{
				AuctionDuration:    180,
				AuctionStartAmount: "1000000000000000000",
				AuctionEndAmount:   "900000000000000000",
				Points: []CustomPresetPoint{
					{ToTokenAmount: "900000000000000000", Delay: 60}, // Exactly at end
				},
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.preset.Validate()
			
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseGasPriceEstimate(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    uint32
		expectError bool
	}{
		{
			name:        "Valid gas price",
			input:       "1000000000",
			expected:    1000000000,
			expectError: false,
		},
		{
			name:        "Zero gas price",
			input:       "0",
			expected:    0,
			expectError: false,
		},
		{
			name:        "Empty string",
			input:       "",
			expected:    0,
			expectError: false,
		},
		{
			name:        "Invalid - not a number",
			input:       "not-a-number",
			expected:    0,
			expectError: true,
		},
		{
			name:        "Invalid - negative",
			input:       "-100",
			expected:    0,
			expectError: true,
		},
		{
			name:        "Max uint32 value",
			input:       "4294967295",
			expected:    4294967295,
			expectError: false,
		},
		{
			name:        "Exceeds uint32",
			input:       "4294967296",
			expected:    0,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseGasPriceEstimate(tc.input)
			
			if tc.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestPresetTypes(t *testing.T) {
	assert.Equal(t, PresetType("custom"), PresetCustom)
	assert.Equal(t, PresetType("fast"), PresetFast)
	assert.Equal(t, PresetType("medium"), PresetMedium)
	assert.Equal(t, PresetType("slow"), PresetSlow)
}
