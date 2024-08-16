package orderbook

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakerTraitsEncode(t *testing.T) {

	tests := []struct {
		name                string
		makerTraitParams    MakerTraitsParams
		expectedMakerTraits string
	}{
		{
			name: "Extension, expiration",
			makerTraitParams: MakerTraitsParams{
				AllowedSender:      "0x0000000000000000000000000000000000000000",
				ShouldCheckEpoch:   false,
				UsePermit2:         false,
				UnwrapWeth:         false,
				HasExtension:       true,
				HasPreInteraction:  false,
				HasPostInteraction: false,
				AllowPartialFills:  true,
				AllowMultipleFills: true,
				Expiry:             1715201499,
				Nonce:              0,
				Series:             0,
			},
			expectedMakerTraits: "0x420000000000000000000000000000000000663be5db00000000000000000000",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			makerTraits, err := NewMakerTraits(tc.makerTraitParams)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedMakerTraits, makerTraits.Encode())
		})
	}
}
