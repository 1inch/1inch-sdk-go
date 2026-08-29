package fusionplus

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gethCommon "github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk-go/v4/common/fusionorder"
	"github.com/1inch/1inch-sdk-go/v4/constants"
)

// TestConvertToOrderbookExtensionIsPure pins funds-critical behavior: encoding
// the escrow extension must not mutate the receiver. The historical bug
// appended the escrow extra data to the extension in place, so a second call
// double-appended it — a corrupted extension whose salt no longer matches,
// producing an unfillable or wrongly-signed order for any caller converting
// more than once.
func TestConvertToOrderbookExtensionIsPure(t *testing.T) {
	ext := &EscrowExtension{
		Extension: Extension{
			PostInteraction: "deadbeef",
		},
		HashLock:         &HashLock{Value: "0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"},
		DstChainId:       constants.BaseChainId,
		SrcSafetyDeposit: "1000",
		DstSafetyDeposit: "2000",
	}

	first, err := ext.ConvertToOrderbookExtension()
	require.NoError(t, err)
	assert.Equal(t, "deadbeef", ext.PostInteraction, "receiver must not be mutated")

	second, err := ext.ConvertToOrderbookExtension()
	require.NoError(t, err)
	assert.Equal(t, first.PostInteraction, second.PostInteraction, "conversion must be idempotent")
}

// TestEscrowExtensionRoundTrip pins funds-critical behavior: an escrow
// extension must survive encode -> decode with every settlement-critical
// field intact. The historical decoder swapped the source and destination
// safety deposits, rendered deposits in hex where the encoder parses decimal,
// and rendered the hashlock in decimal where the encoder parses hex — so a
// decode-then-re-encode silently corrupted the escrow terms.
func TestEscrowExtensionRoundTrip(t *testing.T) {
	original, err := NewEscrowExtension(EscrowExtensionParams{
		ExtensionParamsPlus: ExtensionParamsPlus{
			SettlementContract: "0x1111111254eeb25477b68fb85ed929f73a960582",
			AuctionDetails: &fusionorder.AuctionDetails{
				StartTime:       1700000000,
				Duration:        180,
				InitialRateBump: 5000,
				Points:          nil,
				GasCost:         fusionorder.GasCostConfig{},
			},
			PostInteractionData: &SettlementPostInteractionData{
				Whitelist: []fusionorder.WhitelistItem{},
				IntegratorFee: &IntegratorFee{
					Ratio:    big.NewInt(0),
					Receiver: gethCommon.Address{},
				},
				BankFee:            big.NewInt(0),
				ResolvingStartTime: big.NewInt(1700000000),
				CustomReceiver:     gethCommon.Address{},
			},
			Asset: "0x82af49447d8a07e3bd95bd0d56f35241523fbab1",
		},
		HashLock:         &HashLock{Value: "0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"},
		DstChainId:       constants.BaseChainId,
		DstToken:         gethCommon.HexToAddress("0x4200000000000000000000000000000000000006"),
		SrcSafetyDeposit: "1000000000000000", // distinct values so a swap is caught
		DstSafetyDeposit: "2000000000000000",
		TimeLocks: TimeLocks{
			SrcWithdrawal: 10, SrcPublicWithdrawal: 20, SrcCancellation: 30, SrcPublicCancellation: 40,
			DstWithdrawal: 50, DstPublicWithdrawal: 60, DstCancellation: 70,
		},
	})
	require.NoError(t, err)

	converted, err := original.ConvertToOrderbookExtension()
	require.NoError(t, err)
	encoded, err := converted.Encode()
	require.NoError(t, err)
	encodedBytes, err := hex.DecodeString(strings.TrimPrefix(encoded, "0x"))
	require.NoError(t, err)

	decoded, err := DecodeEscrowExtension(encodedBytes)
	require.NoError(t, err)

	assert.Equal(t, original.HashLock.Value, decoded.HashLock.Value, "hashlock must round-trip exactly")
	assert.Equal(t, original.DstChainId, decoded.DstChainId)
	assert.Equal(t, original.DstToken, decoded.DstToken)
	assert.Equal(t, original.SrcSafetyDeposit, decoded.SrcSafetyDeposit, "source safety deposit must not be swapped or reformatted")
	assert.Equal(t, original.DstSafetyDeposit, decoded.DstSafetyDeposit, "destination safety deposit must not be swapped or reformatted")
	assert.Equal(t, original.TimeLocks, decoded.TimeLocks)
}

// TestSettlementPostInteractionDataRoundTrip pins funds-critical behavior:
// the settlement post-interaction data carries the resolver whitelist (who
// may fill the order), the integrator fee, and the custom receiver. It must
// survive encode -> decode intact.
func TestSettlementPostInteractionDataRoundTrip(t *testing.T) {
	spid, err := NewSettlementPostInteractionDataWithFees(SettlementSuffixData{
		Whitelist: []fusionorder.AuctionWhitelistItem{
			{Address: gethCommon.HexToAddress("0x1111111254eeb25477b68fb85ed929f73a960582"), AllowFrom: big.NewInt(1700000100)},
			{Address: gethCommon.HexToAddress("0x2222222254eeb25477b68fb85ed929f73a960582"), AllowFrom: big.NewInt(1700000200)},
		},
		IntegratorFee: &IntegratorFee{
			Ratio:    big.NewInt(100),
			Receiver: gethCommon.HexToAddress("0x3333333333333333333333333333333333333333"),
		},
		BankFee:            big.NewInt(7),
		ResolvingStartTime: big.NewInt(1700000000),
		CustomReceiver:     gethCommon.HexToAddress("0x4444444444444444444444444444444444444444"),
	})
	require.NoError(t, err)

	encoded, err := spid.Encode()
	require.NoError(t, err)

	decoded, err := DecodeSettlementPostInteractionData(encoded)
	require.NoError(t, err)

	assert.Equal(t, spid.Whitelist, decoded.Whitelist, "whitelist controls who may fill the order")
	assert.Equal(t, spid.IntegratorFee, decoded.IntegratorFee, "integrator fee routing must round-trip")
	assert.Equal(t, spid.BankFee, decoded.BankFee)
	assert.Equal(t, spid.ResolvingStartTime, decoded.ResolvingStartTime)
	assert.Equal(t, spid.CustomReceiver, decoded.CustomReceiver, "custom receiver is where funds are delivered")
}
