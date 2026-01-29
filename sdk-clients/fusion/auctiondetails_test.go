package fusion

import (
	"testing"

	"github.com/1inch/1inch-sdk-go/common/fusionorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeAuctionDetails(t *testing.T) {
	tests := []struct {
		name    string
		details AuctionDetails
	}{
		{
			name: "Encode/Decode AuctionDetails",
			details: AuctionDetails{
				Duration:        180,
				StartTime:       1673548149,
				InitialRateBump: 50000,
				Points: []AuctionPointClassFixed{
					{
						Delay:       10,
						Coefficient: 10000,
					},
					{
						Delay:       20,
						Coefficient: 5000,
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			encoded := tc.details.Encode()
			decoded, err := fusionorder.DecodeLegacyAuctionDetails(encoded)
			require.NoError(t, err)
			assert.Equal(t, tc.details, *decoded)
		})
	}
}

func TestEncodeAuctionDetails(t *testing.T) {
	tests := []struct {
		name     string
		details  AuctionDetails
		expected string
	}{
		{
			name: "Encode AuctionDetails",
			details: AuctionDetails{
				GasCost: GasCostConfigClassFixed{
					GasBumpEstimate:  0,
					GasPriceEstimate: 0,
				},
				StartTime:       1673548149,
				Duration:        180,
				InitialRateBump: 50000,
				Points: []AuctionPointClassFixed{
					{
						Delay:       12,
						Coefficient: 20000,
					},
				},
			},
			expected: "0000000000000063c051750000b400c35001004e20000c",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			encoded := tc.details.Encode()
			assert.Equal(t, tc.expected, encoded)
		})
	}
}

func TestIsNonceRequired(t *testing.T) {
	tests := []struct {
		name                string
		allowPartialFills   bool
		allowMultipleFills  bool
		expectedNonceResult bool
	}{
		{
			name:                "Both allowPartialFills and allowMultipleFills are true",
			allowPartialFills:   true,
			allowMultipleFills:  true,
			expectedNonceResult: false,
		},
		{
			name:                "allowPartialFills is false, allowMultipleFills is true",
			allowPartialFills:   false,
			allowMultipleFills:  true,
			expectedNonceResult: true,
		},
		{
			name:                "allowPartialFills is true, allowMultipleFills is false",
			allowPartialFills:   true,
			allowMultipleFills:  false,
			expectedNonceResult: true,
		},
		{
			name:                "Both allowPartialFills and allowMultipleFills are false",
			allowPartialFills:   false,
			allowMultipleFills:  false,
			expectedNonceResult: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := fusionorder.IsNonceRequired(tc.allowPartialFills, tc.allowMultipleFills)
			assert.Equal(t, tc.expectedNonceResult, result)
		})
	}
}
