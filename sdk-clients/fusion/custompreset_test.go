package fusion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuoterCustomPresetRequest_Validate(t *testing.T) {
	tests := []struct {
		name         string
		customPreset CustomPreset
		expectedErr  string
	}{
		{
			name: "auctionStartAmount should be valid",
			customPreset: CustomPreset{
				AuctionDuration:    180,
				AuctionStartAmount: "ama bad string",
				AuctionEndAmount:   "50000",
				Points: []CustomPresetPoint{
					{ToTokenAmount: "90000", Delay: 20},
					{ToTokenAmount: "110000", Delay: 40},
				},
			},
			expectedErr: "invalid auction start amount: ama bad string",
		},
		{
			name: "auctionEndAmount should be valid",
			customPreset: CustomPreset{
				AuctionDuration:    180,
				AuctionStartAmount: "100000",
				AuctionEndAmount:   "ama bad string",
				Points: []CustomPresetPoint{
					{ToTokenAmount: "90000", Delay: 20},
					{ToTokenAmount: "110000", Delay: 40},
				},
			},
			expectedErr: "invalid auction end amount: ama bad string",
		},
		{
			name: "auctionDuration should be valid",
			customPreset: CustomPreset{
				AuctionDuration:    0,
				AuctionStartAmount: "100000",
				AuctionEndAmount:   "50000",
				Points: []CustomPresetPoint{
					{ToTokenAmount: "90000", Delay: 20},
					{ToTokenAmount: "110000", Delay: 40},
				},
			},
			expectedErr: "invalid auction duration: expected positive integer, got 0",
		},
		{
			name: "points should be in range",
			customPreset: CustomPreset{
				AuctionDuration:    180,
				AuctionStartAmount: "100000",
				AuctionEndAmount:   "50000",
				Points: []CustomPresetPoint{
					{ToTokenAmount: "90000", Delay: 20},
					{ToTokenAmount: "110000", Delay: 40},
				},
			},
			expectedErr: "point at index 1 out of auction range [50000, 100000]: 110000",
		},
		{
			name: "points should be in range (below minimum)",
			customPreset: CustomPreset{
				AuctionDuration:    180,
				AuctionStartAmount: "100000",
				AuctionEndAmount:   "50000",
				Points: []CustomPresetPoint{
					{ToTokenAmount: "40000", Delay: 20},
					{ToTokenAmount: "70000", Delay: 40},
				},
			},
			expectedErr: "point at index 0 out of auction range [50000, 100000]: 40000",
		},
		{
			name: "points should be an array of valid amounts",
			customPreset: CustomPreset{
				AuctionDuration:    180,
				AuctionStartAmount: "100000",
				AuctionEndAmount:   "50000",
				Points: []CustomPresetPoint{
					{ToTokenAmount: "ama bad string", Delay: 20},
					{ToTokenAmount: "70000", Delay: 40},
				},
			},
			expectedErr: "invalid point amount at index 0: ama bad string",
		},
		{
			name: "valid custom preset",
			customPreset: CustomPreset{
				AuctionDuration:    180,
				AuctionStartAmount: "100000",
				AuctionEndAmount:   "50000",
				Points: []CustomPresetPoint{
					{ToTokenAmount: "80000", Delay: 20},
					{ToTokenAmount: "60000", Delay: 40},
				},
			},
			expectedErr: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.customPreset.Validate()
			if tc.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err.Error())
			}
		})
	}
}
