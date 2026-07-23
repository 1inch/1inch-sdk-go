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
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/v4/constants"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/fusion"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/orderbook"
)

// Production canary configuration. The trade runs on Polygon with dust-sized
// amounts and alternates direction based on which side of the pair the canary
// wallet currently holds more trades of, so the same funds recycle indefinitely.
const (
	canaryChainId    = 137
	canaryWeth       = "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619"
	canaryUsdc       = "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"
	canaryWethAmount = "200000000000000" // 0.0002 WETH
	canaryUsdcAmount = "500000"          // 0.5 USDC
	canaryApiUrl     = "https://api.1inch.com"
)

// TestProductionCanary places one real permit2 fusion order against the
// production API and waits for a resolver to fill it. It runs only when the
// canary secrets are present:
//
//	DEV_PORTAL_TOKEN  1inch Developer Portal API key
//	CANARY_WALLET_KEY private key of a dedicated wallet holding dust amounts
//	CANARY_NODE_URL   Polygon RPC endpoint
//
// The wallet needs a small POL balance for the one-time ERC20 approvals to
// Permit2; fills themselves are gasless for the maker.
func TestProductionCanary(t *testing.T) {
	apiKey := os.Getenv("DEV_PORTAL_TOKEN")
	walletKey := os.Getenv("CANARY_WALLET_KEY")
	nodeUrl := os.Getenv("CANARY_NODE_URL")
	if apiKey == "" || walletKey == "" || nodeUrl == "" {
		t.Skip("set DEV_PORTAL_TOKEN, CANARY_WALLET_KEY, and CANARY_NODE_URL to run the production canary")
	}

	ctx := context.Background()

	orderbookConfig, err := orderbook.NewConfiguration(orderbook.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: walletKey,
		ChainId:    canaryChainId,
		ApiUrl:     canaryApiUrl,
		ApiKey:     apiKey,
	})
	require.NoError(t, err)
	orderbookClient, err := orderbook.NewClient(orderbookConfig)
	require.NoError(t, err)

	fusionConfig, err := fusion.NewConfiguration(fusion.ConfigurationParams{
		ApiUrl:     canaryApiUrl,
		ApiKey:     apiKey,
		ChainId:    canaryChainId,
		PrivateKey: walletKey,
	})
	require.NoError(t, err)
	fusionClient, err := fusion.NewClient(fusionConfig)
	require.NoError(t, err)

	erc20, err := abi.JSON(strings.NewReader(minimalErc20ABI))
	require.NoError(t, err)

	owner := orderbookClient.Wallet.Address()
	weth := geth_common.HexToAddress(canaryWeth)
	usdc := geth_common.HexToAddress(canaryUsdc)
	wethAmount, _ := new(big.Int).SetString(canaryWethAmount, 10)
	usdcAmount, _ := new(big.Int).SetString(canaryUsdcAmount, 10)

	liveBalance := func(token geth_common.Address) *big.Int {
		callData, err := erc20.Pack("balanceOf", owner)
		require.NoError(t, err)
		result, err := orderbookClient.Wallet.Call(ctx, token, callData)
		require.NoError(t, err)
		return new(big.Int).SetBytes(result)
	}

	wethBalance := liveBalance(weth)
	usdcBalance := liveBalance(usdc)
	t.Logf("canary wallet %s balances: %s WETH wei, %s USDC units", owner.Hex(), wethBalance, usdcBalance)

	// Sell the side holding more trades worth of balance so direction alternates
	// once the wallet is roughly balanced: compare wethBalance/wethAmount against
	// usdcBalance/usdcAmount via cross-multiplication
	sellWeth := new(big.Int).Mul(wethBalance, usdcAmount).Cmp(new(big.Int).Mul(usdcBalance, wethAmount)) >= 0

	sellToken, buyToken, sellAmount := weth, usdc, wethAmount
	if !sellWeth {
		sellToken, buyToken, sellAmount = usdc, weth, usdcAmount
	}
	require.True(t, liveBalance(sellToken).Cmp(sellAmount) >= 0,
		"canary wallet %s cannot cover a %s trade of %s; fund it with dust amounts of WETH and USDC on Polygon",
		owner.Hex(), sellToken.Hex(), sellAmount)
	t.Logf("canary direction: sell %s of %s for %s", sellAmount, sellToken.Hex(), buyToken.Hex())

	ensureCanaryPermit2Approval(t, ctx, orderbookClient, erc20, sellToken, sellAmount)

	allowance, err := orderbook.GetPermit2Allowance(ctx, orderbookClient.Wallet, owner, sellToken, geth_common.HexToAddress(constants.AggregationRouterV6))
	require.NoError(t, err)

	expiration := big.NewInt(time.Now().Add(30 * time.Minute).Unix())
	permit, err := orderbook.BuildPermit2Calldata(orderbookClient.Wallet, orderbook.Permit2PermitParams{
		Token:       sellToken,
		Amount:      sellAmount,
		Expiration:  expiration,
		Nonce:       allowance.Nonce,
		Spender:     geth_common.HexToAddress(constants.AggregationRouterV6),
		SigDeadline: expiration,
	})
	require.NoError(t, err)

	initSellBalance := liveBalance(sellToken)
	initBuyBalance := liveBalance(buyToken)

	orderHash, err := fusionClient.PlaceOrderFromParams(ctx, fusion.OrderParams{
		WalletAddress:    owner.Hex(),
		FromTokenAddress: strings.ToLower(sellToken.Hex()),
		ToTokenAddress:   strings.ToLower(buyToken.Hex()),
		Amount:           sellAmount.String(),
		Receiver:         constants.ZeroAddress,
		Preset:           fusion.Fast,
		Permit:           permit,
		IsPermit2:        true,
	})
	require.NoError(t, err, "the production API rejected the permit2 order")
	t.Logf("order placed: %s", orderHash)

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
			finalSellBalance := liveBalance(sellToken)
			finalBuyBalance := liveBalance(buyToken)
			require.Equal(t, sellAmount.String(), new(big.Int).Sub(initSellBalance, finalSellBalance).String(), "sell amount spent")
			require.Equal(t, 1, finalBuyBalance.Cmp(initBuyBalance), "buy balance increased")
			t.Logf("canary filled: received %s of %s", new(big.Int).Sub(finalBuyBalance, initBuyBalance), buyToken.Hex())
			return
		case "expired", "cancelled", "refunded", "false-predicate", "not-enough-balance-or-allowance", "wrong-permit":
			t.Fatalf("order %s ended without filling: %s", orderHash, order.Status)
		}
	}
}

// ensureCanaryPermit2Approval sends the one-time ERC20 approval of the sell token to
// the Permit2 contract when the existing allowance cannot cover the trade
func ensureCanaryPermit2Approval(t *testing.T, ctx context.Context, client *orderbook.Client, erc20 abi.ABI, token geth_common.Address, required *big.Int) {
	t.Helper()
	permit2 := geth_common.HexToAddress(constants.Permit2Address)

	allowanceData, err := erc20.Pack("allowance", client.Wallet.Address(), permit2)
	require.NoError(t, err)
	result, err := client.Wallet.Call(ctx, token, allowanceData)
	require.NoError(t, err)
	if new(big.Int).SetBytes(result).Cmp(required) >= 0 {
		return
	}

	t.Logf("sending one-time ERC20 approval of %s to Permit2", token.Hex())
	approveData, err := erc20.Pack("approve", permit2, constants.Uint256Max)
	require.NoError(t, err)
	toAddress := token
	tx, err := client.TxBuilder.New().SetData(approveData).SetTo(&toAddress).Build(ctx)
	require.NoError(t, err)
	signedTx, err := client.Wallet.Sign(tx)
	require.NoError(t, err)
	require.NoError(t, client.Wallet.BroadcastTransaction(ctx, signedTx))

	approvalDeadline := time.Now().Add(3 * time.Minute)
	for time.Now().Before(approvalDeadline) {
		receipt, err := client.Wallet.TransactionReceipt(ctx, signedTx.Hash())
		if err == nil {
			require.Equal(t, uint64(1), receipt.Status, "approval tx reverted: %s", signedTx.Hash().Hex())
			return
		}
		time.Sleep(3 * time.Second)
	}
	t.Fatalf("timed out waiting for approval receipt: %s", signedTx.Hash().Hex())
}
