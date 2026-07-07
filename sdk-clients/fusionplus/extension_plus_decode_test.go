package fusionplus

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/common/fusionorder"
	"github.com/ethereum/go-ethereum/common"
)

// TestFromLimitOrderExtension_RoundTrip verifies a built extension decodes back
// correctly, including the maker permit token and data and the post-interaction data
func TestFromLimitOrderExtension_RoundTrip(t *testing.T) {
	tests := []struct {
		name           string
		asset          string
		permit         string
		expectedAsset  string
		expectedPermit string
	}{
		{
			name:           "Permit round-trips with token field",
			asset:          "0x00000000000000000000000000000000000f1234",
			permit:         "0x3456",
			expectedAsset:  "0x00000000000000000000000000000000000f1234",
			expectedPermit: "0x3456",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			built, err := NewExtensionPlus(ExtensionParamsPlus{
				SettlementContract: "0x5678",
				AuctionDetails: &fusionorder.AuctionDetails{
					StartTime:       0,
					Duration:        0,
					InitialRateBump: 0,
					Points:          nil,
					GasCost:         fusionorder.GasCostConfigClassFixed{},
				},
				PostInteractionData: &SettlementPostInteractionData{
					Whitelist: []fusionorder.WhitelistItem{},
					IntegratorFee: &IntegratorFee{
						Ratio:    big.NewInt(0),
						Receiver: common.Address{},
					},
					BankFee:            big.NewInt(0),
					ResolvingStartTime: big.NewInt(0),
					CustomReceiver:     common.Address{},
				},
				Asset:  tc.asset,
				Permit: tc.permit,

				MakerAssetSuffix: "0x1234",
				TakerAssetSuffix: "0x1234",
				Predicate:        "0x1234",
				PreInteraction:   "pre",
			})
			require.NoError(t, err)

			decoded, err := FromLimitOrderExtension(built.ConvertToOrderbookExtension())
			require.NoError(t, err)

			assert.Equal(t, tc.expectedAsset, decoded.Asset)
			assert.Equal(t, tc.expectedPermit, decoded.Permit)
			assert.Equal(t, built.MakerPermit, decoded.MakerPermit)
		})
	}
}
