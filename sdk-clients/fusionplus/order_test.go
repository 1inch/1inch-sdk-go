package fusionplus

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPreset(t *testing.T) {
	// Define test presets
	customPreset := &Preset{
		AllowMultipleFills: true,
		AllowPartialFills:  true,
		AuctionDuration:    10,
		AuctionEndAmount:   "1000",
		AuctionStartAmount: "500",
		InitialRateBump:    0.1,
		Points:             []AuctionPoint{},
		StartAuctionIn:     1,
	}

	fastPreset := &Preset{
		AllowMultipleFills: false,
		AllowPartialFills:  false,
		AuctionDuration:    20,
		AuctionEndAmount:   "2000",
		AuctionStartAmount: "1000",
		InitialRateBump:    0.2,
		Points:             []AuctionPoint{},
		StartAuctionIn:     2,
	}

	mediumPreset := &Preset{
		AllowMultipleFills: true,
		AllowPartialFills:  false,
		AuctionDuration:    30,
		AuctionEndAmount:   "3000",
		AuctionStartAmount: "1500",
		InitialRateBump:    0.3,
		Points:             []AuctionPoint{},
		StartAuctionIn:     3,
	}

	slowPreset := &Preset{
		AllowMultipleFills: false,
		AllowPartialFills:  true,
		AuctionDuration:    40,
		AuctionEndAmount:   "4000",
		AuctionStartAmount: "2000",
		InitialRateBump:    0.4,
		Points:             []AuctionPoint{},
		StartAuctionIn:     4,
	}

	// Define test cases
	tests := []struct {
		name       string
		presets    QuotePresets
		presetType GetQuoteOutputRecommendedPreset
		expected   *Preset
		expectErr  bool
	}{
		{
			name: "Get Custom Preset",
			presets: QuotePresets{
				Custom: customPreset,
			},
			presetType: Custom,
			expected:   customPreset,
			expectErr:  false,
		},
		{
			name: "Get Fast Preset",
			presets: QuotePresets{
				Fast: *fastPreset,
			},
			presetType: Fast,
			expected:   fastPreset,
			expectErr:  false,
		},
		{
			name: "Get Medium Preset",
			presets: QuotePresets{
				Medium: *mediumPreset,
			},
			presetType: Medium,
			expected:   mediumPreset,
			expectErr:  false,
		},
		{
			name: "Get Slow Preset",
			presets: QuotePresets{
				Slow: *slowPreset,
			},
			presetType: Slow,
			expected:   slowPreset,
			expectErr:  false,
		},
		{
			name: "Unknown Preset Type",
			presets: QuotePresets{
				Custom: customPreset,
				Fast:   *fastPreset,
				Medium: *mediumPreset,
				Slow:   *slowPreset,
			},
			presetType: GetQuoteOutputRecommendedPreset("unknown"),
			expected:   nil,
			expectErr:  true,
		},
		{
			name: "Nil Presets",
			presets: QuotePresets{
				Custom: nil,
				Fast:   Preset{},
				Medium: Preset{},
				Slow:   Preset{},
			},
			presetType: Custom,
			expected:   nil,
			expectErr:  true,
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GetPreset(tc.presets, tc.presetType)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestCreateAuctionDetails(t *testing.T) {
	tests := []struct {
		name                 string
		preset               *Preset
		additionalWaitPeriod float32
		expected             *AuctionDetails
		expectErr            bool
	}{
		{
			name: "Valid preset with points",
			preset: &Preset{
				Points: []AuctionPoint{
					{Coefficient: 100, Delay: 10},
					{Coefficient: 200, Delay: 20},
				},
				GasCost: GasCostConfig{
					GasBumpEstimate:  1,
					GasPriceEstimate: "100",
				},
				AuctionDuration:    300,
				AuctionEndAmount:   "1000",
				AuctionStartAmount: "500",
				InitialRateBump:    1,
				StartAuctionIn:     5,
			},
			additionalWaitPeriod: 2,
			expected: &AuctionDetails{
				StartTime:       CalcAuctionStartTimeFunc(5, 2),
				Duration:        300,
				InitialRateBump: 1,
				Points: []AuctionPointClassFixed{
					{Coefficient: 100, Delay: 10},
					{Coefficient: 200, Delay: 20},
				},
				GasCost: GasCostConfigClassFixed{
					GasBumpEstimate:  1,
					GasPriceEstimate: 100,
				},
			},
			expectErr: false,
		},
		{
			name: "Invalid gas price estimate",
			preset: &Preset{
				Points: []AuctionPoint{
					{Coefficient: 100, Delay: 10},
				},
				GasCost: GasCostConfig{
					GasBumpEstimate:  1,
					GasPriceEstimate: "invalid",
				},
				AuctionDuration:    300,
				AuctionEndAmount:   "1000",
				AuctionStartAmount: "500",
				InitialRateBump:    0.1,
				StartAuctionIn:     5,
			},
			additionalWaitPeriod: 2,
			expected:             nil,
			expectErr:            true,
		},
		{
			name: "Empty points",
			preset: &Preset{
				Points: []AuctionPoint{},
				GasCost: GasCostConfig{
					GasBumpEstimate:  1,
					GasPriceEstimate: "100",
				},
				AuctionDuration:    300,
				AuctionEndAmount:   "1000",
				AuctionStartAmount: "500",
				InitialRateBump:    1,
				StartAuctionIn:     5,
			},
			additionalWaitPeriod: 2,
			expected: &AuctionDetails{
				StartTime:       CalcAuctionStartTimeFunc(5, 2),
				Duration:        300,
				InitialRateBump: 1,
				Points:          []AuctionPointClassFixed{},
				GasCost: GasCostConfigClassFixed{
					GasBumpEstimate:  1,
					GasPriceEstimate: 100,
				},
			},
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CreateAuctionDetails(tc.preset, tc.additionalWaitPeriod)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestIsNonceRequired(t *testing.T) {
	tests := []struct {
		name               string
		allowPartialFills  bool
		allowMultipleFills bool
		expected           bool
	}{
		{
			name:               "Both allowPartialFills and allowMultipleFills are true",
			allowPartialFills:  true,
			allowMultipleFills: true,
			expected:           false,
		},
		{
			name:               "allowPartialFills is true, allowMultipleFills is false",
			allowPartialFills:  true,
			allowMultipleFills: false,
			expected:           true,
		},
		{
			name:               "allowPartialFills is false, allowMultipleFills is true",
			allowPartialFills:  false,
			allowMultipleFills: true,
			expected:           true,
		},
		{
			name:               "Both allowPartialFills and allowMultipleFills are false",
			allowPartialFills:  false,
			allowMultipleFills: false,
			expected:           true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isNonceRequired(tc.allowPartialFills, tc.allowMultipleFills)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestBpsToRatioFormat(t *testing.T) {
	tests := []struct {
		name     string
		bps      *big.Int
		expected *big.Int
	}{
		{
			name:     "Nil bps",
			bps:      nil,
			expected: big.NewInt(0),
		},
		{
			name:     "Zero bps",
			bps:      big.NewInt(0),
			expected: big.NewInt(0),
		},
		{
			name:     "Positive bps",
			bps:      big.NewInt(500),
			expected: big.NewInt(5000),
		},
		{
			name:     "Large bps",
			bps:      big.NewInt(10000),
			expected: big.NewInt(100000),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := bpsToRatioFormat(tc.bps)
			assert.Equal(t, tc.expected, result)
		})
	}
}
