package fusion

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuctionDetails(t *testing.T) {
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
			decoded, err := DecodeAuctionDetails(encoded)
			require.NoError(t, err)
			assert.Equal(t, tc.details, decoded)
		})
	}
}
