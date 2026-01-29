package fusionplus

import (
	"testing"

	"github.com/1inch/1inch-sdk-go/common/fusionorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuctionDetails(t *testing.T) {
	tests := []struct {
		name    string
		details fusionorder.AuctionDetails
	}{
		{
			name: "Encode/Decode AuctionDetails",
			details: fusionorder.AuctionDetails{
				Duration:        180,
				StartTime:       1673548149,
				InitialRateBump: 50000,
				Points: []fusionorder.AuctionPointClassFixed{
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
			// FusionPlus uses encoding without point count byte
			encoded := tc.details.EncodeWithoutPointCount()
			decoded, err := fusionorder.DecodeAuctionDetails(encoded)
			require.NoError(t, err)
			assert.Equal(t, tc.details, *decoded)
		})
	}
}
