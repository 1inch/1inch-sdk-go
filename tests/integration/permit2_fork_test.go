//go:build integration

package integration

import (
	"context"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	geth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/common"
	"github.com/1inch/1inch-sdk-go/constants"
	transaction_builder "github.com/1inch/1inch-sdk-go/internal/transaction-builder"
	web3_provider "github.com/1inch/1inch-sdk-go/internal/web3-provider"
	"github.com/1inch/1inch-sdk-go/sdk-clients/fusion"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

const minimalErc20ABI = `[
	{"name":"approve","type":"function","inputs":[{"name":"spender","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[{"name":"","type":"bool"}]},
	{"name":"transfer","type":"function","inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[{"name":"","type":"bool"}]},
	{"name":"balanceOf","type":"function","stateMutability":"view","inputs":[{"name":"owner","type":"address"}],"outputs":[{"name":"","type":"uint256"}]},
	{"name":"deposit","type":"function","stateMutability":"payable","inputs":[],"outputs":[]}
]`

type testActor struct {
	wallet    *web3_provider.Wallet
	txBuilder common.TransactionBuilderFactory
	address   geth_common.Address
}

func newActor(t *testing.T, pk, nodeURL string) *testActor {
	t.Helper()
	wallet, err := web3_provider.DefaultWalletProvider(pk, nodeURL, 1)
	require.NoError(t, err)
	return &testActor{
		wallet:    wallet,
		txBuilder: transaction_builder.NewFactory(wallet),
		address:   wallet.Address(),
	}
}

type forkEnv struct {
	node       *forkNode
	erc20      abi.ABI
	deployer   *testActor
	maker      *testActor
	taker      *testActor
	settlement geth_common.Address
}

func (e *forkEnv) balanceOf(t *testing.T, token, owner geth_common.Address) *big.Int {
	t.Helper()
	callData, err := e.erc20.Pack("balanceOf", owner)
	require.NoError(t, err)
	result, err := e.deployer.wallet.Call(context.Background(), token, callData)
	require.NoError(t, err)
	return new(big.Int).SetBytes(result)
}

// setupForkEnv boots the fork, deploys SimpleSettlement, and funds maker (WETH) and taker (USDC)
func setupForkEnv(t *testing.T) *forkEnv {
	t.Helper()

	node := startAnvil(t)

	erc20, err := abi.JSON(strings.NewReader(minimalErc20ABI))
	require.NoError(t, err)

	env := &forkEnv{
		node:     node,
		erc20:    erc20,
		deployer: newActor(t, freshPrivateKey(t), node.url),
		maker:    newActor(t, freshPrivateKey(t), node.url),
		taker:    newActor(t, freshPrivateKey(t), node.url),
	}

	hundredEth := new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18))
	for _, addr := range []string{env.deployer.address.Hex(), env.maker.address.Hex(), env.taker.address.Hex(), usdcDonorAddress} {
		node.setBalance(t, addr, hundredEth)
	}

	// Deploy SimpleSettlement(limitOrderProtocol, accessToken, weth, owner)
	artifact := loadArtifact(t, "testdata/SimpleSettlement.json")
	settlementABI, err := abi.JSON(strings.NewReader(string(artifact.Abi)))
	require.NoError(t, err)
	constructorArgs, err := settlementABI.Constructor.Inputs.Pack(
		geth_common.HexToAddress(lopV4Address),
		geth_common.HexToAddress(accessToken),
		geth_common.HexToAddress(wethAddress),
		env.deployer.address,
	)
	require.NoError(t, err)
	creationCode := append(geth_common.FromHex(artifact.Bytecode), constructorArgs...)
	env.settlement = node.deployContract(t, env.deployer.wallet, env.deployer.txBuilder, creationCode)
	t.Logf("SimpleSettlement deployed at %s", env.settlement.Hex())

	// Maker wraps 1 ETH into WETH
	depositData, err := erc20.Pack("deposit")
	require.NoError(t, err)
	node.sendTx(t, env.maker.wallet, env.maker.txBuilder, wethAddress, depositData, big.NewInt(1e18))

	// Taker receives USDC from a mainnet donor via impersonation
	usdcAmount := big.NewInt(1_000_000_000) // 1000 USDC
	donorBalance := env.balanceOf(t, geth_common.HexToAddress(usdcAddress), geth_common.HexToAddress(usdcDonorAddress))
	require.True(t, donorBalance.Cmp(usdcAmount) >= 0, "USDC donor %s has insufficient balance %s at fork block", usdcDonorAddress, donorBalance)
	transferData, err := erc20.Pack("transfer", env.taker.address, usdcAmount)
	require.NoError(t, err)
	node.sendImpersonated(t, usdcDonorAddress, usdcAddress, transferData)

	// Maker approves WETH to Permit2. The Permit2 -> LOP allowance itself is granted
	// on-chain by the maker permit embedded in the order.
	approveData, err := erc20.Pack("approve", geth_common.HexToAddress(constants.Permit2Address), constants.Uint256Max)
	require.NoError(t, err)
	node.sendTx(t, env.maker.wallet, env.maker.txBuilder, wethAddress, approveData, nil)

	// Taker pays USDC through a regular transferFrom, so approve the LOP directly
	takerApproveData, err := erc20.Pack("approve", geth_common.HexToAddress(lopV4Address), constants.Uint256Max)
	require.NoError(t, err)
	node.sendTx(t, env.taker.wallet, env.taker.txBuilder, usdcAddress, takerApproveData, nil)

	return env
}

// buildPermit2Calldata signs a PermitSingle for makingAmount so the allowance is fully consumed by the fill
func buildPermit2Calldata(t *testing.T, env *forkEnv, makingAmount *big.Int) string {
	t.Helper()
	allowance, err := orderbook.GetPermit2Allowance(
		context.Background(),
		env.maker.wallet,
		env.maker.address,
		geth_common.HexToAddress(wethAddress),
		geth_common.HexToAddress(lopV4Address),
	)
	require.NoError(t, err)

	maxUint48 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 48), big.NewInt(1))
	permitCalldata, err := orderbook.BuildPermit2Calldata(env.maker.wallet, orderbook.Permit2PermitParams{
		Token:       geth_common.HexToAddress(wethAddress),
		Amount:      makingAmount,
		Expiration:  maxUint48,
		Nonce:       allowance.Nonce,
		Spender:     geth_common.HexToAddress(lopV4Address),
		SigDeadline: maxUint48,
	})
	require.NoError(t, err)
	return permitCalldata
}

func quoteFixture(env *forkEnv, takingAmount string) fusion.GetQuoteOutputFixed {
	preset := fusion.PresetClassFixed{
		AllowMultipleFills: true,
		AllowPartialFills:  true,
		AuctionDuration:    120,
		AuctionEndAmount:   takingAmount,
		AuctionStartAmount: takingAmount,
		GasCost:            fusion.GasCostConfigClass{GasBumpEstimate: 0, GasPriceEstimate: "0"},
		InitialRateBump:    0,
		Points:             nil,
		StartAuctionIn:     0,
	}
	return fusion.GetQuoteOutputFixed{
		QuoteId:           "fork-test",
		SettlementAddress: strings.ToLower(env.settlement.Hex()),
		Whitelist:         []string{strings.ToLower(env.taker.address.Hex())},
		MarketAmount:      takingAmount,
		SurplusFee:        0,
		Presets:           fusion.QuotePresetsClassFixed{Fast: preset},
	}
}

// packFillOrderArgs builds the LOP v4 fillOrderArgs calldata for a locally created order
func packFillOrderArgs(t *testing.T, limitOrder *orderbook.Order, amount *big.Int, makerAmountMode bool) []byte {
	t.Helper()

	routerABI, err := abi.JSON(strings.NewReader(constants.AggregationRouterV6ABI))
	require.NoError(t, err)

	salt, ok := new(big.Int).SetString(limitOrder.Data.Salt, 10)
	if !ok {
		salt, ok = new(big.Int).SetString(strings.TrimPrefix(limitOrder.Data.Salt, "0x"), 16)
		require.True(t, ok, "invalid salt: %s", limitOrder.Data.Salt)
	}
	makingAmount, ok := new(big.Int).SetString(limitOrder.Data.MakingAmount, 10)
	require.True(t, ok)
	takingAmount, ok := new(big.Int).SetString(limitOrder.Data.TakingAmount, 10)
	require.True(t, ok)
	makerTraits, err := hexutil.DecodeBig(limitOrder.Data.MakerTraits)
	require.NoError(t, err)

	orderTuple := orderbook.NormalizedLimitOrderData{
		Salt:         salt,
		Maker:        orderbook.AddressStringToBigInt(limitOrder.Data.Maker),
		Receiver:     orderbook.AddressStringToBigInt(limitOrder.Data.Receiver),
		MakerAsset:   orderbook.AddressStringToBigInt(limitOrder.Data.MakerAsset),
		TakerAsset:   orderbook.AddressStringToBigInt(limitOrder.Data.TakerAsset),
		MakingAmount: makingAmount,
		TakingAmount: takingAmount,
		MakerTraits:  makerTraits,
	}

	compact, err := orderbook.CompressSignature(limitOrder.Signature[2:])
	require.NoError(t, err)
	var r, vs [32]byte
	copy(r[:], compact.R)
	copy(vs[:], compact.VS)

	takerTraits := orderbook.NewTakerTraits(orderbook.TakerTraitsParams{
		Extension: limitOrder.Data.Extension,
	})
	encodedTraits, err := takerTraits.Encode()
	require.NoError(t, err)
	traitFlags := encodedTraits.TraitFlags
	if makerAmountMode {
		// TakerTraits.Encode does not handle the amount mode flag, so set bit 255 directly

		traitFlags = new(big.Int).Or(traitFlags, new(big.Int).Lsh(big.NewInt(1), orderbook.MakerAmountFlag))
	}

	fillData, err := routerABI.Pack("fillOrderArgs", orderTuple, r, vs, amount, traitFlags, encodedTraits.Args)
	require.NoError(t, err)
	return fillData
}

func TestFusionOrderPermit2Fork(t *testing.T) {
	env := setupForkEnv(t)

	weth := geth_common.HexToAddress(wethAddress)
	usdc := geth_common.HexToAddress(usdcAddress)
	makingAmount := new(big.Int).Div(big.NewInt(1e18), big.NewInt(10)) // 0.1 WETH
	takingAmount := big.NewInt(100_000_000)                            // 100 USDC

	t.Run("permit2 order fills through Permit2", func(t *testing.T) {
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

		// The first 20 bytes of the maker permit carry the token the permit applies to
		makerPermit := preparedOrder.Order.FusionExtension.MakerPermit
		require.True(t,
			strings.EqualFold(makerPermit[:42], wethAddress),
			"maker permit token field = %s, want the maker asset", makerPermit[:42])

		initBalances := map[string]*big.Int{
			"makerWeth": env.balanceOf(t, weth, env.maker.address),
			"takerWeth": env.balanceOf(t, weth, env.taker.address),
			"makerUsdc": env.balanceOf(t, usdc, env.maker.address),
			"takerUsdc": env.balanceOf(t, usdc, env.taker.address),
		}

		fillData := packFillOrderArgs(t, limitOrder, makingAmount, true)

		// Ensure the fill block timestamp is at or after the auction/resolving start time
		env.node.setNextBlockTimestamp(t, time.Now().Unix()+5)
		env.node.sendTx(t, env.taker.wallet, env.taker.txBuilder, lopV4Address, fillData, nil)

		finalBalances := map[string]*big.Int{
			"makerWeth": env.balanceOf(t, weth, env.maker.address),
			"takerWeth": env.balanceOf(t, weth, env.taker.address),
			"makerUsdc": env.balanceOf(t, usdc, env.maker.address),
			"takerUsdc": env.balanceOf(t, usdc, env.taker.address),
		}

		assert.Equal(t, makingAmount.String(), new(big.Int).Sub(initBalances["makerWeth"], finalBalances["makerWeth"]).String(), "maker WETH spent")
		assert.Equal(t, makingAmount.String(), new(big.Int).Sub(finalBalances["takerWeth"], initBalances["takerWeth"]).String(), "taker WETH received")
		assert.Equal(t, takingAmount.String(), new(big.Int).Sub(finalBalances["makerUsdc"], initBalances["makerUsdc"]).String(), "maker USDC received")
		assert.Equal(t, takingAmount.String(), new(big.Int).Sub(initBalances["takerUsdc"], finalBalances["takerUsdc"]).String(), "taker USDC spent")

		// The finite Permit2 allowance granted by the embedded permit must be fully
		// consumed, proving maker funds moved through Permit2
		finalAllowance, err := orderbook.GetPermit2Allowance(context.Background(), env.maker.wallet, env.maker.address, weth, geth_common.HexToAddress(lopV4Address))
		require.NoError(t, err)
		assert.Equal(t, "0", finalAllowance.Amount.String(), "Permit2 allowance fully consumed")
	})

	// Before the fix, the fusion path silently dropped OrderParams.Permit and IsPermit2,
	// producing an order with no maker permit and no USE_PERMIT2 traits bit. For a maker
	// who only approved Permit2 (never the LOP directly), such an order cannot be filled.
	// This is the user-visible failure the fix resolves.
	//
	// Note the permit interaction TARGET itself is not observable on-chain: the LOP's
	// tryPermit dispatches 352-byte permits to the canonical Permit2 address regardless
	// of the encoded target. The target fix matters for SDK, API, and decoder consistency
	// and is asserted above and in the unit tests.
	t.Run("pre-fix behavior (permit dropped) fails to fill", func(t *testing.T) {
		quote := quoteFixture(env, takingAmount.String())
		orderParams := fusion.OrderParams{
			FromTokenAddress:   wethAddress,
			ToTokenAddress:     usdcAddress,
			Amount:             makingAmount.String(),
			WalletAddress:      strings.ToLower(env.maker.address.Hex()),
			Receiver:           zeroAddress,
			Preset:             fusion.Fast,
			Permit:             "", // pre-fix: permit never reached the order
			IsPermit2:          false,
			AllowPartialFills:  true,
			AllowMultipleFills: true,
		}

		_, limitOrder, err := fusion.CreateFusionOrderData(quote, orderParams, env.maker.wallet, 1)
		require.NoError(t, err)

		fillData := packFillOrderArgs(t, limitOrder, makingAmount, true)

		env.node.setNextBlockTimestamp(t, time.Now().Unix()+10)
		status := env.node.trySendTx(t, env.taker.wallet, env.taker.txBuilder, lopV4Address, fillData)
		assert.Equal(t, uint64(0), status, "fill without the embedded permit must revert for a Permit2-only maker")
	})

	// Resolvers routinely partial-fill fusion orders during auctions. The permit is applied
	// on the first fill; subsequent fills spend the remaining Permit2 allowance.
	t.Run("permit2 order supports partial fills", func(t *testing.T) {
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

		_, limitOrder, err := fusion.CreateFusionOrderData(quote, orderParams, env.maker.wallet, 1)
		require.NoError(t, err)

		initMakerWeth := env.balanceOf(t, weth, env.maker.address)
		initTakerWeth := env.balanceOf(t, weth, env.taker.address)

		halfAmount := new(big.Int).Div(makingAmount, big.NewInt(2))
		fillData := packFillOrderArgs(t, limitOrder, halfAmount, true)

		env.node.setNextBlockTimestamp(t, time.Now().Unix()+15)
		env.node.sendTx(t, env.taker.wallet, env.taker.txBuilder, lopV4Address, fillData, nil)
		env.node.sendTx(t, env.taker.wallet, env.taker.txBuilder, lopV4Address, fillData, nil)

		finalMakerWeth := env.balanceOf(t, weth, env.maker.address)
		finalTakerWeth := env.balanceOf(t, weth, env.taker.address)

		assert.Equal(t, makingAmount.String(), new(big.Int).Sub(initMakerWeth, finalMakerWeth).String(), "maker WETH spent across both fills")
		assert.Equal(t, makingAmount.String(), new(big.Int).Sub(finalTakerWeth, initTakerWeth).String(), "taker WETH received across both fills")

		finalAllowance, err := orderbook.GetPermit2Allowance(context.Background(), env.maker.wallet, env.maker.address, weth, geth_common.HexToAddress(lopV4Address))
		require.NoError(t, err)
		assert.Equal(t, "0", finalAllowance.Amount.String(), "Permit2 allowance fully consumed across both fills")
	})

	// The compact 96-byte permit is reconstructed on-chain from the maker permit's token
	// field. On the currently deployed router that reconstruction leaves dirty upper
	// bits in the amount slot, which Permit2's uint160 calldata validation rejects, so
	// the fill reverts. This subtest pins that behavior; if a future router deployment
	// accepts compact permits, this test will flag it by failing.
	t.Run("compact permit2 order reverts on current router", func(t *testing.T) {
		allowance, err := orderbook.GetPermit2Allowance(
			context.Background(),
			env.maker.wallet,
			env.maker.address,
			weth,
			geth_common.HexToAddress(lopV4Address),
		)
		require.NoError(t, err)

		maxUint48 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 48), big.NewInt(1))
		permitCalldata, err := orderbook.BuildPermit2CalldataCompact(env.maker.wallet, orderbook.Permit2PermitParams{
			Token:       weth,
			Amount:      makingAmount,
			Expiration:  maxUint48,
			Nonce:       allowance.Nonce,
			Spender:     geth_common.HexToAddress(lopV4Address),
			SigDeadline: maxUint48,
		})
		require.NoError(t, err)

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

		_, limitOrder, err := fusion.CreateFusionOrderData(quote, orderParams, env.maker.wallet, 1)
		require.NoError(t, err)

		fillData := packFillOrderArgs(t, limitOrder, makingAmount, true)
		env.node.setNextBlockTimestamp(t, time.Now().Unix()+18)
		status := env.node.trySendTx(t, env.taker.wallet, env.taker.txBuilder, lopV4Address, fillData)
		assert.Equal(t, uint64(0), status, "compact permit fill is expected to revert on the current router")
	})

	// Regression: ordinary orders without any permit must be unaffected by the permit wiring
	t.Run("order without permit fills with direct allowance", func(t *testing.T) {
		approveData, err := env.erc20.Pack("approve", geth_common.HexToAddress(lopV4Address), constants.Uint256Max)
		require.NoError(t, err)
		env.node.sendTx(t, env.maker.wallet, env.maker.txBuilder, wethAddress, approveData, nil)

		quote := quoteFixture(env, takingAmount.String())
		orderParams := fusion.OrderParams{
			FromTokenAddress:   wethAddress,
			ToTokenAddress:     usdcAddress,
			Amount:             makingAmount.String(),
			WalletAddress:      strings.ToLower(env.maker.address.Hex()),
			Receiver:           zeroAddress,
			Preset:             fusion.Fast,
			AllowPartialFills:  true,
			AllowMultipleFills: true,
		}

		_, limitOrder, err := fusion.CreateFusionOrderData(quote, orderParams, env.maker.wallet, 1)
		require.NoError(t, err)

		initMakerWeth := env.balanceOf(t, weth, env.maker.address)

		fillData := packFillOrderArgs(t, limitOrder, makingAmount, true)
		env.node.setNextBlockTimestamp(t, time.Now().Unix()+20)
		env.node.sendTx(t, env.taker.wallet, env.taker.txBuilder, lopV4Address, fillData, nil)

		finalMakerWeth := env.balanceOf(t, weth, env.maker.address)
		assert.Equal(t, makingAmount.String(), new(big.Int).Sub(initMakerWeth, finalMakerWeth).String(), "maker WETH spent on plain order")
	})
}
