package fusionplus

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/v4/common/fusionorder"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/orderbook"
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

// TestDecodeExtension_RoundTrip verifies the raw byte entry point reproduces the
// built extension, exercising orderbook.Decode and FromLimitOrderExtension together
func TestDecodeExtension_RoundTrip(t *testing.T) {
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
		Asset:  "0x00000000000000000000000000000000000f1234",
		Permit: "0x3456",

		MakerAssetSuffix: "0x1234",
		TakerAssetSuffix: "0x1234",
		Predicate:        "0x1234",
		PreInteraction:   "0x9abc",
	})
	require.NoError(t, err)

	encoded, err := built.ConvertToOrderbookExtension().Encode()
	require.NoError(t, err)

	decoded, err := DecodeExtension(common.FromHex(encoded))
	require.NoError(t, err)

	assert.Equal(t, "0x00000000000000000000000000000000000f1234", decoded.Asset)
	assert.Equal(t, "0x3456", decoded.Permit)
	assert.Equal(t, built.MakerPermit, decoded.MakerPermit)
	assert.Equal(t, built.MakingAmountData, decoded.MakingAmountData)
	assert.Equal(t, built.PostInteraction, decoded.PostInteraction)
}

// TestDecodeSettlementPostInteractionData_Empty verifies empty input returns an error
// instead of panicking
func TestDecodeSettlementPostInteractionData_Empty(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		errorMsg string
	}{
		{
			name:     "Empty hex data",
			data:     "0x",
			errorMsg: "post interaction data is empty",
		},
		{
			name:     "Missing prefix",
			data:     "",
			errorMsg: "invalid hex string",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecodeSettlementPostInteractionData(tc.data)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.errorMsg)
		})
	}
}

// TestFromLimitOrderExtension_ShortFields verifies malformed extensions with fields
// shorter than a settlement address return errors instead of panicking
func TestFromLimitOrderExtension_ShortFields(t *testing.T) {
	tests := []struct {
		name      string
		extension *orderbook.Extension
	}{
		{
			name:      "Empty fields",
			extension: &orderbook.Extension{MakingAmountData: "0x", TakingAmountData: "0x", PostInteraction: "0x"},
		},
		{
			name: "Post interaction shorter than an address",
			extension: &orderbook.Extension{
				MakingAmountData: "0x0000000000000000000000000000000000005678",
				TakingAmountData: "0x0000000000000000000000000000000000005678",
				PostInteraction:  "0x12",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := FromLimitOrderExtension(tc.extension)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "malformed extension")
		})
	}
}
