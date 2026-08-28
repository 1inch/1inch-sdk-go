package fusionplus

import (
	"context"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/v4/common"
	"github.com/1inch/1inch-sdk-go/v4/constants"
)

// mainnetTrueERC20 is the Ethereum Fusion+ TRUE_ERC20 sentinel.
const mainnetTrueERC20 = "0xda0000d4000015a526378bb6fafc650cea5966f8"

// collisionToken is a real ERC-20 (1INCH). It is the type of destination token
// that can move real value on the source chain. This happens if the code writes
// it into the source-chain taker asset.
const collisionToken = "0x111111111117dc0aa78b770fa6a738034120c302"

// TestFusionPlusSourceTakerAssetIsTrueERC20 checks the funds-critical rule.
// The source-chain order must use the chain's TRUE_ERC20 sentinel as its taker
// asset. It must never use the destination token. The escrow extension must
// hold the real destination token without a change. This includes the native
// 0xEeee… sentinel. The defect is a destination token in the source taker
// asset. This test guards against that defect.
func TestFusionPlusSourceTakerAssetIsTrueERC20(t *testing.T) {
	tests := []struct {
		name       string
		dstToken   string
		wantDstHex string // expected escrow extension DstToken (EIP-55)
	}{
		{
			name:       "ERC-20 destination",
			dstToken:   collisionToken,
			wantDstHex: gethcommon.HexToAddress(collisionToken).Hex(),
		},
		{
			name:       "native destination stays the native sentinel",
			dstToken:   constants.NativeToken,
			wantDstHex: gethcommon.HexToAddress(constants.NativeToken).Hex(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			prepared, err := CreateFusionPlusOrderData(
				takerAssetTestQuoteParams(tc.dstToken),
				takerAssetTestQuote(),
				takerAssetTestOrderParams(),
				&stubWallet{addr: gethcommon.HexToAddress("0x4444444444444444444444444444444444444444")},
				constants.EthereumChainId,
			)
			require.NoError(t, err)

			signedTaker := strings.ToLower(prepared.LimitOrder.Data.TakerAsset)
			assert.Equal(t, mainnetTrueERC20, signedTaker,
				"source-chain taker asset must be the TRUE_ERC20 sentinel")
			assert.NotEqual(t, strings.ToLower(tc.dstToken), signedTaker,
				"source-chain taker asset must never be the destination token")

			assert.Equal(t, tc.wantDstHex, prepared.Order.EscExtension.DstToken.Hex(),
				"escrow extension must carry the real destination token")
		})
	}
}

// TestFusionPlusUnsupportedSourceChainFailsLoudly checks the behavior for a
// source chain that has no TRUE_ERC20 sentinel. Such a chain is not a Fusion+
// chain. Order construction must fail with an error. It must not sign an order
// that has a wrong taker asset.
func TestFusionPlusUnsupportedSourceChainFailsLoudly(t *testing.T) {
	_, err := CreateFusionPlusOrderData(
		takerAssetTestQuoteParams(collisionToken),
		takerAssetTestQuote(),
		takerAssetTestOrderParams(),
		&stubWallet{addr: gethcommon.HexToAddress("0x4444444444444444444444444444444444444444")},
		constants.AuroraChainId,
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not supported")
}

func takerAssetTestQuoteParams(dstToken string) QuoteParams {
	return QuoteParams{
		SrcChain:        constants.EthereumChainId,
		DstChain:        constants.BscChainId,
		SrcTokenAddress: "0x5555555555555555555555555555555555555555",
		DstTokenAddress: dstToken,
		Amount:          "1000000",
		WalletAddress:   "0x4444444444444444444444444444444444444444",
		Fee:             0,
	}
}

func takerAssetTestOrderParams() OrderParams {
	return OrderParams{
		HashLock: &HashLock{Value: "0x00000000000000000000000000000000000000000000000000000000000000ff"},
		Preset:   Fast,
		Receiver: constants.ZeroAddress,
		Nonce:    big.NewInt(1),
	}
}

func takerAssetTestQuote() *Quote {
	return &Quote{
		QuoteId:          "local-only",
		SrcEscrowFactory: "0x3333333333333333333333333333333333333333",
		DstEscrowFactory: "0x3333333333333333333333333333333333333333",
		SrcSafetyDeposit: "17",
		DstSafetyDeposit: "34",
		Presets: QuotePresets{
			Fast: Preset{
				AllowMultipleFills: true,
				AllowPartialFills:  true,
				AuctionDuration:    180,
				AuctionEndAmount:   "900000",
				AuctionStartAmount: "950000",
				GasCost:            GasCostConfig{GasBumpEstimate: 1, GasPriceEstimate: "1"},
				InitialRateBump:    1,
				Points:             []AuctionPoint{},
				StartAuctionIn:     1,
			},
		},
		TimeLocks: TimeLocks{
			DstCancellation:       7,
			DstPublicWithdrawal:   6,
			DstWithdrawal:         5,
			SrcPublicCancellation: 4,
			SrcCancellation:       3,
			SrcPublicWithdrawal:   2,
			SrcWithdrawal:         1,
		},
		Whitelist: []string{"0x2222222222222222222222222222222222222222"},
	}
}

// stubWallet is a common.Wallet that does nothing, for local order
// construction. SignBytes returns a 65-byte signature so the limit-order
// signature step completes.
type stubWallet struct{ addr gethcommon.Address }

func (w *stubWallet) Call(context.Context, gethcommon.Address, []byte) ([]byte, error) {
	return nil, nil
}
func (w *stubWallet) Nonce(context.Context) (uint64, error)          { return 0, nil }
func (w *stubWallet) Address() gethcommon.Address                    { return w.addr }
func (w *stubWallet) Balance(context.Context) (*big.Int, error)      { return big.NewInt(0), nil }
func (w *stubWallet) GetGasTipCap(context.Context) (*big.Int, error) { return big.NewInt(0), nil }
func (w *stubWallet) GetGasPrice(context.Context) (*big.Int, error)  { return big.NewInt(0), nil }
func (w *stubWallet) GetGasEstimate(context.Context, ethereum.CallMsg) (uint64, error) {
	return 0, nil
}
func (w *stubWallet) Sign(tx *types.Transaction) (*types.Transaction, error) { return tx, nil }
func (w *stubWallet) SignBytes([]byte) ([]byte, error)                       { return make([]byte, 65), nil }
func (w *stubWallet) BroadcastTransaction(context.Context, *types.Transaction) error {
	return nil
}
func (w *stubWallet) TransactionReceipt(context.Context, gethcommon.Hash) (*types.Receipt, error) {
	return nil, nil
}
func (w *stubWallet) GetContractDetailsForPermit(context.Context, gethcommon.Address, gethcommon.Address, *big.Int, int64) (*common.ContractPermitData, error) {
	return nil, nil
}
func (w *stubWallet) GetContractDetailsForPermitDaiLike(context.Context, gethcommon.Address, gethcommon.Address, int64) (*common.ContractPermitDataDaiLike, error) {
	return nil, nil
}
func (w *stubWallet) TokenPermit(common.ContractPermitData) (string, error) { return "", nil }
func (w *stubWallet) TokenPermitDaiLike(common.ContractPermitDataDaiLike) (string, error) {
	return "", nil
}
func (w *stubWallet) IsEIP1559Applicable() bool { return false }
func (w *stubWallet) ChainId() int64            { return int64(constants.EthereumChainId) }
