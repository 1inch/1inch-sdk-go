//go:build integration

package integration

import (
	"context"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	geth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/v4/constants"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/aggregation"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/fusion"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/fusionplus"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/orderbook"
)

// Production canaries place real dust-sized trades through the production API,
// covering each allowance mechanism the SDK supports:
//
//	Fusion (Base):        direct approval, EIP-2612 permit, Permit2 permit
//	Aggregation (Arbitrum): direct approval, EIP-2612 permit, Permit2 allowance
//	Fusion+ (Base <-> Arbitrum): direct approval (cross-chain permits are not
//	  offered by the API, so probing them would fail by design)
//
// They run only when the canary secrets are present and alternate direction based
// on current balances, so the same funds recycle indefinitely. Fusion and Fusion+
// fills are gasless for the maker; gas is needed for one-time approvals and for
// aggregation swaps.
const canaryApiUrl = "https://api.1inch.com"

const permit2ApproveABI = `[{"inputs":[{"internalType":"address","name":"token","type":"address"},{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint160","name":"amount","type":"uint160"},{"internalType":"uint48","name":"expiration","type":"uint48"}],"name":"approve","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

type canaryChain struct {
	name       string
	chainId    uint64
	rpcEnv     string
	weth       string
	usdc       string
	wethAmount *big.Int
	usdcAmount *big.Int
}

var canaryBase = canaryChain{
	name:       "base",
	chainId:    8453,
	rpcEnv:     "CANARY_BASE_RPC_URL",
	weth:       "0x4200000000000000000000000000000000000006",
	usdc:       "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
	wethAmount: big.NewInt(200_000_000_000_000), // 0.0002 WETH
	usdcAmount: big.NewInt(500_000),             // 0.5 USDC
}

var canaryArbitrum = canaryChain{
	name:       "arbitrum",
	chainId:    42161,
	rpcEnv:     "CANARY_ARBITRUM_RPC_URL",
	weth:       "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1",
	usdc:       "0xaf88d065e77c8cC2239327C5EDb3A432268e5831",
	wethAmount: big.NewInt(200_000_000_000_000), // 0.0002 WETH
	usdcAmount: big.NewInt(500_000),             // 0.5 USDC
}

// canaryFusionPlusAmount is the USDC amount bridged per Fusion+ run; cross-chain
// orders carry safety deposit economics, so it is larger than the single-chain dust
var canaryFusionPlusAmount = big.NewInt(1_500_000) // 1.5 USDC

type canaryActor struct {
	chain     canaryChain
	apiKey    string
	walletKey string
	rpcUrl    string
	orderbook *orderbook.Client
	erc20     abi.ABI
	owner     geth_common.Address
	router    geth_common.Address
	permit2   geth_common.Address
}

// newCanaryActor loads the canary secrets and builds an RPC-connected client for the
// chain, skipping the test when any secret is missing
func newCanaryActor(t *testing.T, chain canaryChain) *canaryActor {
	t.Helper()
	apiKey := os.Getenv("DEV_PORTAL_TOKEN")
	walletKey := os.Getenv("CANARY_WALLET_KEY")
	rpcUrl := os.Getenv(chain.rpcEnv)
	if apiKey == "" || walletKey == "" || rpcUrl == "" {
		t.Skipf("set DEV_PORTAL_TOKEN, CANARY_WALLET_KEY, and %s to run this canary", chain.rpcEnv)
	}

	config, err := orderbook.NewConfiguration(orderbook.ConfigurationParams{
		NodeUrl:    rpcUrl,
		PrivateKey: walletKey,
		ChainId:    chain.chainId,
		ApiUrl:     canaryApiUrl,
		ApiKey:     apiKey,
	})
	require.NoError(t, err)
	client, err := orderbook.NewClient(config)
	require.NoError(t, err)

	erc20, err := abi.JSON(strings.NewReader(minimalErc20ABI))
	require.NoError(t, err)

	return &canaryActor{
		chain:     chain,
		apiKey:    apiKey,
		walletKey: walletKey,
		rpcUrl:    rpcUrl,
		orderbook: client,
		erc20:     erc20,
		owner:     client.Wallet.Address(),
		router:    geth_common.HexToAddress(constants.AggregationRouterV6),
		permit2:   geth_common.HexToAddress(constants.Permit2Address),
	}
}

func (a *canaryActor) balance(t *testing.T, token geth_common.Address) *big.Int {
	t.Helper()
	callData, err := a.erc20.Pack("balanceOf", a.owner)
	require.NoError(t, err)
	result, err := a.orderbook.Wallet.Call(context.Background(), token, callData)
	require.NoError(t, err)
	return new(big.Int).SetBytes(result)
}

func (a *canaryActor) erc20Allowance(t *testing.T, token, spender geth_common.Address) *big.Int {
	t.Helper()
	allowanceData, err := a.erc20.Pack("allowance", a.owner, spender)
	require.NoError(t, err)
	result, err := a.orderbook.Wallet.Call(context.Background(), token, allowanceData)
	require.NoError(t, err)
	return new(big.Int).SetBytes(result)
}

func (a *canaryActor) sendTx(t *testing.T, to geth_common.Address, data []byte) {
	t.Helper()
	ctx := context.Background()
	// Public RPCs load-balance across nodes that can lag a block behind, so a
	// freshly mined transaction may not be reflected in the next nonce query yet;
	// rebuild for a fresh nonce and retry when the broadcast is rejected as stale
	var lastErr error
	for attempt := 0; attempt < 4; attempt++ {
		if attempt > 0 {
			time.Sleep(5 * time.Second)
		}
		tx, err := a.orderbook.TxBuilder.New().
			SetData(data).
			SetTo(&to).
			SetGasFeeCap(a.feeCapWithHeadroom(t)).
			Build(ctx)
		if err == nil {
			var signedTx *types.Transaction
			signedTx, err = a.orderbook.Wallet.Sign(tx)
			require.NoError(t, err)
			err = a.orderbook.Wallet.BroadcastTransaction(ctx, signedTx)
			if err == nil {
				a.waitForReceipt(t, signedTx.Hash(), 3*time.Minute)
				return
			}
		}
		lastErr = err
		require.True(t, retryableBroadcastError(lastErr), "failed to send transaction: %v", lastErr)
		t.Logf("transient send failure, rebuilding and retrying: %v", lastErr)
	}
	t.Fatalf("failed to send transaction after retries: %v", lastErr)
}

// feeCapWithHeadroom doubles the node's suggested gas price. The transaction
// builder defaults the fee cap to the bare suggestion, which chains with a
// volatile base fee (Arbitrum) reject whenever the base fee ticks up before
// inclusion; only base fee plus tip is actually charged, so the headroom is free.
func (a *canaryActor) feeCapWithHeadroom(t *testing.T) *big.Int {
	t.Helper()
	gasPrice, err := a.orderbook.Wallet.GetGasPrice(context.Background())
	require.NoError(t, err)
	return new(big.Int).Mul(gasPrice, big.NewInt(2))
}

// retryableBroadcastError reports whether a build or broadcast failure is a
// transient RPC-state race worth rebuilding for: a stale nonce from a lagging
// node, or a fee cap that fell below a freshly risen base fee
func retryableBroadcastError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "nonce") || strings.Contains(msg, "fee per gas")
}

// awaitBalanceDelta retries the balance read until the expected delta appears,
// tolerating RPC nodes that lag the fill block. A nil expectedDelta accepts any
// increase over the initial balance.
func (a *canaryActor) awaitBalanceDelta(t *testing.T, token geth_common.Address, initial, expectedDelta *big.Int, decreasing bool, label string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Minute)
	var current *big.Int
	for time.Now().Before(deadline) {
		current = a.balance(t, token)
		delta := new(big.Int).Sub(initial, current)
		if !decreasing {
			delta = new(big.Int).Sub(current, initial)
		}
		if expectedDelta == nil {
			if delta.Sign() > 0 {
				t.Logf("%s: delta %s", label, delta)
				return
			}
		} else if delta.Cmp(expectedDelta) == 0 {
			t.Logf("%s: delta %s", label, delta)
			return
		}
		time.Sleep(5 * time.Second)
	}
	t.Fatalf("%s: balance delta did not reach expectation within 2 minutes (initial %s, current %s)", label, initial, current)
}

// ensureErc20Allowance sends a one-time unlimited approval when the current
// allowance to the spender cannot cover the trade
func (a *canaryActor) ensureErc20Allowance(t *testing.T, token, spender geth_common.Address, required *big.Int) {
	t.Helper()
	if a.erc20Allowance(t, token, spender).Cmp(required) >= 0 {
		return
	}
	t.Logf("sending one-time ERC20 approval of %s to %s", token.Hex(), spender.Hex())
	approveData, err := a.erc20.Pack("approve", spender, constants.Uint256Max)
	require.NoError(t, err)
	a.sendTx(t, token, approveData)
}

// ensurePermit2RouterAllowance grants the router a standing allowance inside the
// Permit2 contract (used by the aggregation UsePermit2 flow, which relies on an
// existing Permit2 allowance rather than a per-trade signed permit)
func (a *canaryActor) ensurePermit2RouterAllowance(t *testing.T, token geth_common.Address, required *big.Int) {
	t.Helper()
	allowance, err := orderbook.GetPermit2Allowance(context.Background(), a.orderbook.Wallet, a.owner, token, a.router)
	require.NoError(t, err)
	nowUnix := big.NewInt(time.Now().Unix())
	if allowance.Amount.Cmp(required) >= 0 && allowance.Expiration.Cmp(nowUnix) > 0 {
		return
	}

	t.Logf("granting the router a standing Permit2 allowance for %s", token.Hex())
	permit2Abi, err := abi.JSON(strings.NewReader(permit2ApproveABI))
	require.NoError(t, err)
	maxUint160 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 160), big.NewInt(1))
	approveData, err := permit2Abi.Pack("approve", token, a.router, maxUint160, constants.Uint48Max)
	require.NoError(t, err)
	a.sendTx(t, a.permit2, approveData)
}

func (a *canaryActor) waitForReceipt(t *testing.T, hash geth_common.Hash, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		receipt, err := a.orderbook.Wallet.TransactionReceipt(context.Background(), hash)
		if err == nil {
			require.Equal(t, uint64(1), receipt.Status, "tx reverted: %s", hash.Hex())
			return
		}
		time.Sleep(3 * time.Second)
	}
	t.Fatalf("timed out waiting for receipt of %s", hash.Hex())
}

// buildEip2612Permit signs a classic EIP-2612 permit granting the router exactly the
// trade amount
func (a *canaryActor) buildEip2612Permit(t *testing.T, token geth_common.Address, amount *big.Int) string {
	t.Helper()
	permitData, err := a.orderbook.Wallet.GetContractDetailsForPermit(
		context.Background(), token, a.router, amount, time.Now().Add(30*time.Minute).Unix())
	require.NoError(t, err)
	permit, err := a.orderbook.Wallet.TokenPermit(*permitData)
	require.NoError(t, err)
	return permit
}

// buildPermit2OrderPermit signs a Permit2 PermitSingle granting the router exactly
// the trade amount, for embedding into a fusion order
func (a *canaryActor) buildPermit2OrderPermit(t *testing.T, token geth_common.Address, amount *big.Int) string {
	t.Helper()
	allowance, err := orderbook.GetPermit2Allowance(context.Background(), a.orderbook.Wallet, a.owner, token, a.router)
	require.NoError(t, err)

	expiration := big.NewInt(time.Now().Add(30 * time.Minute).Unix())
	permit, err := orderbook.BuildPermit2Calldata(a.orderbook.Wallet, orderbook.Permit2PermitParams{
		Token:       token,
		Amount:      amount,
		Expiration:  expiration,
		Nonce:       allowance.Nonce,
		Spender:     a.router,
		SigDeadline: expiration,
	})
	require.NoError(t, err)
	return permit
}

// pickDirection sells the side holding more trades worth of balance so the
// direction alternates once the wallet is roughly balanced
func pickDirection(t *testing.T, actor *canaryActor) (sellToken, buyToken geth_common.Address, sellAmount *big.Int) {
	t.Helper()
	chain := actor.chain
	weth := geth_common.HexToAddress(chain.weth)
	usdc := geth_common.HexToAddress(chain.usdc)
	wethBalance := actor.balance(t, weth)
	usdcBalance := actor.balance(t, usdc)

	sellWeth := new(big.Int).Mul(wethBalance, chain.usdcAmount).Cmp(new(big.Int).Mul(usdcBalance, chain.wethAmount)) >= 0
	sellToken, buyToken, sellAmount = weth, usdc, chain.wethAmount
	if !sellWeth {
		sellToken, buyToken, sellAmount = usdc, weth, chain.usdcAmount
	}
	require.True(t, actor.balance(t, sellToken).Cmp(sellAmount) >= 0,
		"canary wallet %s on %s cannot cover a %s trade of %s; fund it with dust WETH and USDC",
		actor.owner.Hex(), chain.name, sellToken.Hex(), sellAmount)
	t.Logf("canary direction on %s: sell %s of %s for %s", chain.name, sellAmount, sellToken.Hex(), buyToken.Hex())
	return sellToken, buyToken, sellAmount
}

// placeFusionOrderAndAwaitFill places one fusion order and waits for a resolver fill,
// asserting the sell amount left the wallet and the buy side arrived
func placeFusionOrderAndAwaitFill(t *testing.T, actor *canaryActor, fusionClient *fusion.Client, sellToken, buyToken geth_common.Address, sellAmount *big.Int, permit string, isPermit2 bool) {
	t.Helper()
	ctx := context.Background()

	initSellBalance := actor.balance(t, sellToken)
	initBuyBalance := actor.balance(t, buyToken)

	orderHash, err := fusionClient.PlaceOrderFromParams(ctx, fusion.OrderParams{
		WalletAddress:    actor.owner.Hex(),
		FromTokenAddress: strings.ToLower(sellToken.Hex()),
		ToTokenAddress:   strings.ToLower(buyToken.Hex()),
		Amount:           sellAmount.String(),
		Receiver:         constants.ZeroAddress,
		Preset:           fusion.Fast,
		Permit:           permit,
		IsPermit2:        isPermit2,
	})
	require.NoError(t, err, "the production API rejected the order")
	t.Logf("fusion order placed: %s", orderHash)

	deadline := time.Now().Add(5 * time.Minute)
	for {
		require.True(t, time.Now().Before(deadline), "order %s was not filled within 5 minutes", orderHash)
		time.Sleep(5 * time.Second)

		order, err := fusionClient.GetOrderStatus(ctx, orderHash)
		if err != nil {
			t.Logf("status poll failed, retrying: %v", err)
			continue
		}
		t.Logf("order status: %s", order.Status)
		switch order.Status {
		case "filled":
			actor.awaitBalanceDelta(t, sellToken, initSellBalance, sellAmount, true, "sell amount spent")
			actor.awaitBalanceDelta(t, buyToken, initBuyBalance, nil, false, "buy balance increased")
			return
		case "expired", "cancelled", "refunded", "false-predicate", "not-enough-balance-or-allowance", "wrong-permit":
			t.Fatalf("order %s ended without filling: %s", orderHash, order.Status)
		}
	}
}

// TestProductionCanaryFusion places dust-sized fusion orders on Base, one per
// allowance mechanism: direct ERC20 approval, EIP-2612 permit, and Permit2 permit
func TestProductionCanaryFusion(t *testing.T) {
	actor := newCanaryActor(t, canaryBase)
	usdc := geth_common.HexToAddress(actor.chain.usdc)
	weth := geth_common.HexToAddress(actor.chain.weth)

	fusionConfig, err := fusion.NewConfiguration(fusion.ConfigurationParams{
		ApiUrl:     canaryApiUrl,
		ApiKey:     actor.apiKey,
		ChainId:    actor.chain.chainId,
		PrivateKey: actor.walletKey,
	})
	require.NoError(t, err)
	fusionClient, err := fusion.NewClient(fusionConfig)
	require.NoError(t, err)

	t.Run("direct approval", func(t *testing.T) {
		sellToken, buyToken, sellAmount := pickDirection(t, actor)
		actor.ensureErc20Allowance(t, sellToken, actor.router, sellAmount)
		placeFusionOrderAndAwaitFill(t, actor, fusionClient, sellToken, buyToken, sellAmount, "", false)
	})

	// The EIP-2612 permit sells USDC (WETH has no permit function); the on-chain
	// permit overwrites any standing allowance with the exact amount, so the fill
	// consuming it to zero proves the permit executed
	t.Run("eip2612 permit", func(t *testing.T) {
		sellAmount := actor.chain.usdcAmount
		require.True(t, actor.balance(t, usdc).Cmp(sellAmount) >= 0, "canary wallet needs %s USDC on base", sellAmount)
		permit := actor.buildEip2612Permit(t, usdc, sellAmount)
		placeFusionOrderAndAwaitFill(t, actor, fusionClient, usdc, weth, sellAmount, permit, false)
		require.Eventually(t, func() bool {
			return actor.erc20Allowance(t, usdc, actor.router).Sign() == 0
		}, time.Minute, 5*time.Second, "the 2612 permit allowance must be fully consumed")
	})

	t.Run("permit2 permit", func(t *testing.T) {
		sellToken, buyToken, sellAmount := pickDirection(t, actor)
		actor.ensureErc20Allowance(t, sellToken, actor.permit2, sellAmount)
		permit := actor.buildPermit2OrderPermit(t, sellToken, sellAmount)
		placeFusionOrderAndAwaitFill(t, actor, fusionClient, sellToken, buyToken, sellAmount, permit, true)
		require.Eventually(t, func() bool {
			finalAllowance, err := orderbook.GetPermit2Allowance(context.Background(), actor.orderbook.Wallet, actor.owner, sellToken, actor.router)
			return err == nil && finalAllowance.Amount.Sign() == 0
		}, time.Minute, 5*time.Second, "the Permit2 allowance must be fully consumed")
	})
}

// TestProductionCanaryFusionPlus bridges USDC between Base and Arbitrum through a
// cross-chain fusion order, selling from whichever chain holds more USDC so the
// funds ping-pong between the two chains across runs
func TestProductionCanaryFusionPlus(t *testing.T) {
	baseActor := newCanaryActor(t, canaryBase)
	arbActor := newCanaryActor(t, canaryArbitrum)
	ctx := context.Background()

	baseUsdc := baseActor.balance(t, geth_common.HexToAddress(canaryBase.usdc))
	arbUsdc := arbActor.balance(t, geth_common.HexToAddress(canaryArbitrum.usdc))
	t.Logf("USDC balances: base=%s arbitrum=%s", baseUsdc, arbUsdc)

	src, dst := canaryBase, canaryArbitrum
	srcActor, dstActor := baseActor, arbActor
	if arbUsdc.Cmp(baseUsdc) > 0 {
		src, dst = canaryArbitrum, canaryBase
		srcActor, dstActor = arbActor, baseActor
	}
	srcToken := geth_common.HexToAddress(src.usdc)
	dstToken := geth_common.HexToAddress(dst.usdc)
	require.True(t, srcActor.balance(t, srcToken).Cmp(canaryFusionPlusAmount) >= 0,
		"canary wallet cannot cover a %s USDC bridge from %s; fund it on both chains", canaryFusionPlusAmount, src.name)
	t.Logf("fusion+ direction: %s USDC from %s to %s", canaryFusionPlusAmount, src.name, dst.name)

	srcActor.ensureErc20Allowance(t, srcToken, srcActor.router, canaryFusionPlusAmount)

	plusConfig, err := fusionplus.NewConfiguration(fusionplus.ConfigurationParams{
		ApiUrl:     canaryApiUrl,
		ApiKey:     srcActor.apiKey,
		PrivateKey: srcActor.walletKey,
	})
	require.NoError(t, err)
	plusClient, err := fusionplus.NewClient(plusConfig)
	require.NoError(t, err)

	quoteParams := fusionplus.QuoterControllerGetQuoteParamsFixed{
		SrcChain:        float32(src.chainId),
		DstChain:        float32(dst.chainId),
		SrcTokenAddress: strings.ToLower(srcToken.Hex()),
		DstTokenAddress: strings.ToLower(dstToken.Hex()),
		Amount:          canaryFusionPlusAmount.String(),
		WalletAddress:   srcActor.owner.Hex(),
		EnableEstimate:  true,
	}
	// The estimation inside the quote can race a base fee update; retry briefly
	var quote *fusionplus.GetQuoteOutputFixed
	for attempt := 0; ; attempt++ {
		quote, err = plusClient.GetQuote(ctx, quoteParams)
		if err == nil {
			break
		}
		require.True(t, attempt < 3, "the production API rejected the cross-chain quote: %v", err)
		t.Logf("cross-chain quote failed, retrying: %v", err)
		time.Sleep(10 * time.Second)
	}

	preset, err := fusionplus.GetPreset(quote.Presets, quote.RecommendedPreset)
	require.NoError(t, err)

	secrets := make([]string, int(preset.SecretsCount))
	secretHashes := make([]string, int(preset.SecretsCount))
	for i := range secrets {
		secrets[i], err = fusionplus.GetRandomBytes32()
		require.NoError(t, err)
		secretHashes[i], err = fusionplus.HashSecret(secrets[i])
		require.NoError(t, err)
	}
	var hashLock *fusionplus.HashLock
	if len(secrets) == 1 {
		hashLock, err = fusionplus.ForSingleFill(secrets[0])
	} else {
		hashLock, err = fusionplus.ForMultipleFills(secrets)
	}
	require.NoError(t, err)

	initDstBalance := dstActor.balance(t, dstToken)
	initSrcBalance := srcActor.balance(t, srcToken)

	orderHash, err := plusClient.PlaceOrder(ctx, quoteParams, quote, fusionplus.OrderParams{
		HashLock:     hashLock,
		SecretHashes: secretHashes,
		Receiver:     constants.ZeroAddress,
		Preset:       quote.RecommendedPreset,
	}, plusClient.Wallet)
	require.NoError(t, err, "the production API rejected the cross-chain order")
	t.Logf("fusion+ order placed: %s", orderHash)

	submitted := 0
	deadline := time.Now().Add(15 * time.Minute)
	for {
		require.True(t, time.Now().Before(deadline), "order %s did not complete within 15 minutes", orderHash)
		time.Sleep(5 * time.Second)

		order, err := plusClient.GetOrderByOrderHash(ctx, fusionplus.GetOrderByOrderHashParams{Hash: orderHash})
		if err != nil {
			t.Logf("status poll failed, retrying: %v", err)
			continue
		}
		t.Logf("order status: %s", order.Status)
		switch string(order.Status) {
		case "executed":
			srcActor.awaitBalanceDelta(t, srcToken, initSrcBalance, canaryFusionPlusAmount, true, "source USDC spent")
			// The destination transfer can land moments after the status flips
			arrivalDeadline := time.Now().Add(2 * time.Minute)
			for time.Now().Before(arrivalDeadline) {
				if dstActor.balance(t, dstToken).Cmp(initDstBalance) > 0 {
					t.Logf("fusion+ canary executed: received %s USDC on %s", new(big.Int).Sub(dstActor.balance(t, dstToken), initDstBalance), dst.name)
					return
				}
				time.Sleep(5 * time.Second)
			}
			t.Fatalf("order %s executed but destination USDC did not arrive within 2 minutes", orderHash)
		case "refunded", "cancelled", "expired":
			t.Fatalf("order %s ended without executing: %s", orderHash, order.Status)
		}

		fills, err := plusClient.GetReadyToAcceptFills(ctx, fusionplus.GetReadyToAcceptFillsParams{Hash: orderHash})
		if err != nil {
			t.Logf("fills poll failed, retrying: %v", err)
			continue
		}
		for ; submitted < len(fills.Fills) && submitted < len(secrets); submitted++ {
			require.NoError(t, plusClient.SubmitSecret(ctx, fusionplus.SecretInput{
				OrderHash: orderHash,
				Secret:    secrets[submitted],
			}))
			t.Logf("submitted secret %d", submitted)
		}
	}
}

// executeAggregationSwap requests swap calldata, broadcasts it, and asserts balances
func executeAggregationSwap(t *testing.T, actor *canaryActor, aggClient *aggregation.Client, sellToken, buyToken geth_common.Address, sellAmount *big.Int, permit string, usePermit2 bool) {
	t.Helper()
	ctx := context.Background()

	initSellBalance := actor.balance(t, sellToken)
	initBuyBalance := actor.balance(t, buyToken)

	swap, err := aggClient.GetSwap(ctx, aggregation.GetSwapParams{
		Src:        strings.ToLower(sellToken.Hex()),
		Dst:        strings.ToLower(buyToken.Hex()),
		Amount:     sellAmount.String(),
		From:       actor.owner.Hex(),
		Origin:     actor.owner.Hex(),
		Slippage:   1,
		Permit:     permit,
		UsePermit2: usePermit2,
	})
	require.NoError(t, err, "the production API rejected the swap request")

	var lastErr error
	for attempt := 0; attempt < 4; attempt++ {
		if attempt > 0 {
			time.Sleep(5 * time.Second)
		}
		tx, err := aggClient.TxBuilder.New().
			SetData(swap.TxNormalized.Data).
			SetTo(&swap.TxNormalized.To).
			SetGas(swap.TxNormalized.Gas).
			SetValue(swap.TxNormalized.Value).
			SetGasFeeCap(actor.feeCapWithHeadroom(t)).
			Build(ctx)
		if err == nil {
			var signedTx *types.Transaction
			signedTx, err = aggClient.Wallet.Sign(tx)
			require.NoError(t, err)
			err = aggClient.Wallet.BroadcastTransaction(ctx, signedTx)
			if err == nil {
				t.Logf("swap broadcast: %s", signedTx.Hash().Hex())
				actor.waitForReceipt(t, signedTx.Hash(), 3*time.Minute)
				lastErr = nil
				break
			}
		}
		lastErr = err
		require.True(t, retryableBroadcastError(lastErr), "failed to send swap: %v", lastErr)
		t.Logf("transient send failure, rebuilding and retrying: %v", lastErr)
	}
	require.NoError(t, lastErr, "failed to send swap after retries")

	actor.awaitBalanceDelta(t, sellToken, initSellBalance, sellAmount, true, "sell amount spent")
	actor.awaitBalanceDelta(t, buyToken, initBuyBalance, nil, false, "buy balance increased")
}

// TestProductionCanaryAggregation performs dust-sized classic swaps on Arbitrum, one
// per allowance mechanism: direct ERC20 approval, EIP-2612 permit, and a standing
// Permit2 allowance. Swaps are self-executed, so the wallet pays gas.
func TestProductionCanaryAggregation(t *testing.T) {
	actor := newCanaryActor(t, canaryArbitrum)
	usdc := geth_common.HexToAddress(actor.chain.usdc)
	weth := geth_common.HexToAddress(actor.chain.weth)

	aggConfig, err := aggregation.NewConfiguration(aggregation.ConfigurationParams{
		NodeUrl:    actor.rpcUrl,
		PrivateKey: actor.walletKey,
		ChainId:    actor.chain.chainId,
		ApiUrl:     canaryApiUrl,
		ApiKey:     actor.apiKey,
	})
	require.NoError(t, err)
	aggClient, err := aggregation.NewClient(aggConfig)
	require.NoError(t, err)

	t.Run("direct approval", func(t *testing.T) {
		sellToken, buyToken, sellAmount := pickDirection(t, actor)
		actor.ensureErc20Allowance(t, sellToken, actor.router, sellAmount)
		executeAggregationSwap(t, actor, aggClient, sellToken, buyToken, sellAmount, "", false)
	})

	t.Run("eip2612 permit", func(t *testing.T) {
		sellAmount := actor.chain.usdcAmount
		require.True(t, actor.balance(t, usdc).Cmp(sellAmount) >= 0, "canary wallet needs %s USDC on arbitrum", sellAmount)
		permit := actor.buildEip2612Permit(t, usdc, sellAmount)
		executeAggregationSwap(t, actor, aggClient, usdc, weth, sellAmount, permit, false)
	})

	t.Run("permit2 allowance", func(t *testing.T) {
		sellToken, buyToken, sellAmount := pickDirection(t, actor)
		actor.ensureErc20Allowance(t, sellToken, actor.permit2, sellAmount)
		actor.ensurePermit2RouterAllowance(t, sellToken, sellAmount)
		executeAggregationSwap(t, actor, aggClient, sellToken, buyToken, sellAmount, "", true)
	})
}
