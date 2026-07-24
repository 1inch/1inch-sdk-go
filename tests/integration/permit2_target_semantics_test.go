//go:build integration

package integration

import (
	"math/big"
	"strings"
	"testing"
	"time"

	geth_common "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/v4/constants"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/fusion"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/orderbook"
)

// TestPermit2TargetSemanticsFork pins the protocol semantics of the maker permit's
// leading 20 bytes for the full 352-byte permit2 form: tryPermit dispatches those
// permits to the canonical Permit2 contract by calldata length and ignores the token
// field, so an order encoded with the Permit2 address in place of the maker asset
// still fills. The SDK always encodes the maker asset there; this test demonstrates
// the choice is an SDK/API consistency matter with no on-chain effect for the full
// form.
func TestPermit2TargetSemanticsFork(t *testing.T) {
	env := setupForkEnv(t)

	weth := geth_common.HexToAddress(wethAddress)
	usdc := geth_common.HexToAddress(usdcAddress)
	makingAmount := new(big.Int).Div(big.NewInt(1e18), big.NewInt(10)) // 0.1 WETH
	takingAmount := big.NewInt(100_000_000)                            // 100 USDC

	permitCalldata := buildPermit2Calldata(t, env, makingAmount)

	quote := quoteFixture(env, takingAmount.String())
	orderParams := fusion.OrderParams{
		FromTokenAddress:   wethAddress,
		ToTokenAddress:     usdcAddress,
		Amount:             makingAmount.String(),
		WalletAddress:      strings.ToLower(env.maker.address.Hex()),
		Receiver:           zeroAddress,
		Preset:             fusion.Fast,
		Permit:             permitCalldata,
		IsPermit2:          true,
		AllowPartialFills:  true,
		AllowMultipleFills: true,
	}

	preparedOrder, limitOrder, err := fusion.CreateFusionOrderData(quote, orderParams, env.maker.wallet, 1)
	require.NoError(t, err)

	// Rewrite the permit token field from the maker asset to the Permit2 address,
	// then re-encode the extension, regenerate the salt, and re-sign the order
	doctored := preparedOrder.Order.FusionExtension
	require.True(t, strings.EqualFold(doctored.MakerPermit[:42], wethAddress))
	doctored.MakerPermit = strings.ToLower(geth_common.HexToAddress(constants.Permit2Address).String()) + doctored.MakerPermit[42:]

	orderbookExtension := doctored.ConvertToOrderbookExtension()
	extensionEncoded, err := orderbookExtension.Encode()
	require.NoError(t, err)
	salt, err := orderbook.GenerateSalt(extensionEncoded, nil)
	require.NoError(t, err)

	makerTraits, err := orderbook.DecodeMakerTraits(limitOrder.Data.MakerTraits)
	require.NoError(t, err)
	require.True(t, makerTraits.ShouldUsePermit2)

	doctoredOrder, err := orderbook.CreateLimitOrderMessage(orderbook.CreateOrderParams{
		Wallet:           env.maker.wallet,
		MakerTraits:      makerTraits,
		Extension:        *orderbookExtension,
		ExtensionEncoded: extensionEncoded,
		Salt:             salt,
		Maker:            limitOrder.Data.Maker,
		MakerAsset:       limitOrder.Data.MakerAsset,
		TakerAsset:       limitOrder.Data.TakerAsset,
		TakingAmount:     limitOrder.Data.TakingAmount,
		MakingAmount:     limitOrder.Data.MakingAmount,
		Taker:            limitOrder.Data.Receiver,
	}, 1)
	require.NoError(t, err)

	initMakerWeth := env.balanceOf(t, weth, env.maker.address)
	initMakerUsdc := env.balanceOf(t, usdc, env.maker.address)

	fillData := packFillOrderArgs(t, doctoredOrder, makingAmount, true)
	env.node.setNextBlockTimestamp(t, time.Now().Unix()+5)
	env.node.sendTx(t, env.taker.wallet, env.taker.txBuilder, lopV4Address, fillData, nil)

	finalMakerWeth := env.balanceOf(t, weth, env.maker.address)
	finalMakerUsdc := env.balanceOf(t, usdc, env.maker.address)

	assert.Equal(t, makingAmount.String(), new(big.Int).Sub(initMakerWeth, finalMakerWeth).String(), "maker WETH spent despite doctored permit token field")
	assert.Equal(t, takingAmount.String(), new(big.Int).Sub(finalMakerUsdc, initMakerUsdc).String(), "maker USDC received despite doctored permit token field")
}
