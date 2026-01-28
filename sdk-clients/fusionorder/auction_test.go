package fusionorder

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuctionDetails(t *testing.T) {
	tests := []struct {
		name            string
		startTime       uint32
		duration        uint32
		initialRateBump uint32
		gasCost         GasCostConfigClassFixed
		shouldError     bool
	}{
		{
			name:            "Valid auction details",
			startTime:       1673548149,
			duration:        180,
			initialRateBump: 50000,
			gasCost: GasCostConfigClassFixed{
				GasBumpEstimate:  10000,
				GasPriceEstimate: 1000000000,
			},
			shouldError: false,
		},
		{
			name:            "Duration exceeds uint24 max",
			startTime:       1673548149,
			duration:        Uint24Max + 1,
			initialRateBump: 50000,
			gasCost: GasCostConfigClassFixed{
				GasBumpEstimate:  10000,
				GasPriceEstimate: 1000000000,
			},
			shouldError: true,
		},
		{
			name:            "InitialRateBump exceeds uint24 max",
			startTime:       1673548149,
			duration:        180,
			initialRateBump: Uint24Max + 1,
			gasCost: GasCostConfigClassFixed{
				GasBumpEstimate:  10000,
				GasPriceEstimate: 1000000000,
			},
			shouldError: true,
		},
		{
			name:            "GasBumpEstimate exceeds uint24 max",
			startTime:       1673548149,
			duration:        180,
			initialRateBump: 50000,
			gasCost: GasCostConfigClassFixed{
				GasBumpEstimate:  Uint24Max + 1,
				GasPriceEstimate: 1000000000,
			},
			shouldError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := NewAuctionDetails(tc.startTime, tc.duration, tc.initialRateBump, nil, tc.gasCost)
			if tc.shouldError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.startTime, result.StartTime)
				assert.Equal(t, tc.duration, result.Duration)
				assert.Equal(t, tc.initialRateBump, result.InitialRateBump)
			}
		})
	}
}

func TestAuctionDetailsEncodeDecode(t *testing.T) {
	tests := []struct {
		name    string
		details AuctionDetails
	}{
		{
			name: "Basic auction details without points",
			details: AuctionDetails{
				StartTime:       1673548149,
				Duration:        180,
				InitialRateBump: 50000,
				Points:          nil, // Decoded points will be nil, not empty slice
				GasCost: GasCostConfigClassFixed{
					GasBumpEstimate:  0,
					GasPriceEstimate: 0,
				},
			},
		},
		{
			name: "Auction details with points",
			details: AuctionDetails{
				StartTime:       1673548149,
				Duration:        180,
				InitialRateBump: 50000,
				Points: []AuctionPointClassFixed{
					{Coefficient: 10000, Delay: 10},
					{Coefficient: 5000, Delay: 20},
				},
				GasCost: GasCostConfigClassFixed{
					GasBumpEstimate:  10000,
					GasPriceEstimate: 1000000000,
				},
			},
		},
		{
			name: "Auction details with single point",
			details: AuctionDetails{
				StartTime:       1700000000,
				Duration:        300,
				InitialRateBump: 100000,
				Points: []AuctionPointClassFixed{
					{Coefficient: 25000, Delay: 30},
				},
				GasCost: GasCostConfigClassFixed{
					GasBumpEstimate:  5000,
					GasPriceEstimate: 500000000,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Test Encode -> DecodeLegacyAuctionDetails roundtrip
			encoded := tc.details.Encode()
			decoded, err := DecodeLegacyAuctionDetails(encoded)
			require.NoError(t, err)
			assert.Equal(t, tc.details, *decoded)
		})
	}
}

func TestAuctionDetailsEncodeWithoutPointCount(t *testing.T) {
	details := AuctionDetails{
		StartTime:       1673548149,
		Duration:        180,
		InitialRateBump: 50000,
		Points: []AuctionPointClassFixed{
			{Coefficient: 10000, Delay: 10},
			{Coefficient: 5000, Delay: 20},
		},
		GasCost: GasCostConfigClassFixed{
			GasBumpEstimate:  0,
			GasPriceEstimate: 0,
		},
	}

	// Test EncodeWithoutPointCount -> DecodeAuctionDetails roundtrip
	encoded := details.EncodeWithoutPointCount()
	decoded, err := DecodeAuctionDetails(encoded)
	require.NoError(t, err)
	assert.Equal(t, details, *decoded)
}

func TestDecodeAuctionDetailsErrors(t *testing.T) {
	tests := []struct {
		name        string
		data        string
		expectError bool
	}{
		{
			name:        "Invalid hex",
			data:        "not-hex-data",
			expectError: true,
		},
		{
			name:        "Data too short",
			data:        "0102030405",
			expectError: true,
		},
		{
			name:        "Incomplete point data",
			data:        "0000000000000063c051750000b400c35000270f", // 17 bytes header + 3 bytes (incomplete point)
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecodeAuctionDetails(tc.data)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCalcAuctionStartTime(t *testing.T) {
	// Test that the function returns a time in the future
	startAuctionIn := uint32(60)      // 60 seconds
	additionalWait := uint32(30)      // 30 seconds
	
	before := time.Now().Unix()
	result := CalcAuctionStartTime(startAuctionIn, additionalWait)
	after := time.Now().Unix()
	
	// Result should be current time + startAuctionIn + additionalWait
	expectedMin := uint32(before) + startAuctionIn + additionalWait
	expectedMax := uint32(after) + startAuctionIn + additionalWait
	
	assert.GreaterOrEqual(t, result, expectedMin)
	assert.LessOrEqual(t, result, expectedMax)
}

func TestCalcAuctionStartTimeZeroDelays(t *testing.T) {
	before := time.Now().Unix()
	result := CalcAuctionStartTime(0, 0)
	after := time.Now().Unix()
	
	assert.GreaterOrEqual(t, result, uint32(before))
	assert.LessOrEqual(t, result, uint32(after))
}
