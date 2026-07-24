package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/1inch/1inch-sdk-go/v4/constants"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/aggregation"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/orderbook"
)

/*
This example swaps WETH for USDC on Arbitrum using a standing Permit2 allowance
instead of a direct ERC20 approval to the 1inch router.

Permit2 uses two layers of approval:

 1. A one-time on-chain ERC20 approval from the sell token to the canonical
    Permit2 contract (constants.Permit2Address).
 2. A standing allowance inside Permit2 granting the 1inch Aggregation Router
    spending rights, set here with an on-chain permit2.approve call.

Once both are in place, swaps only need the UsePermit2 flag; tokens flow through
Permit2 with its amount- and time-bounded allowance instead of an unlimited
router approval.

Requires the following environment variables:
  - DEV_PORTAL_TOKEN: 1inch Developer Portal API key
  - WALLET_KEY:       private key (64 hex chars, no 0x prefix)
  - NODE_URL:         RPC endpoint for Arbitrum
*/

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
)

const (
	arbitrumWeth = "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1"
	arbitrumUsdc = "0xaf88d065e77c8cC2239327C5EDb3A432268e5831"
	amountWeth   = "200000000000000" // 0.0002 WETH
	chainId      = 42161
)

const permit2ApproveABI = `[{"inputs":[{"internalType":"address","name":"token","type":"address"},{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint160","name":"amount","type":"uint160"},{"internalType":"uint48","name":"expiration","type":"uint48"}],"name":"approve","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

func main() {
	if devPortalToken == "" || privateKey == "" || nodeUrl == "" {
		log.Fatal("set DEV_PORTAL_TOKEN, WALLET_KEY, and NODE_URL to run this example")
	}
	ctx := context.Background()

	aggConfig, err := aggregation.NewConfiguration(aggregation.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: privateKey,
		ChainId:    chainId,
		ApiUrl:     "https://api.1inch.com",
		ApiKey:     devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}
	client, err := aggregation.NewClient(aggConfig)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	owner := client.Wallet.Address()
	sellToken := gethCommon.HexToAddress(arbitrumWeth)
	router := gethCommon.HexToAddress(constants.AggregationRouterV6)
	permit2 := gethCommon.HexToAddress(constants.Permit2Address)
	amount, ok := new(big.Int).SetString(amountWeth, 10)
	if !ok {
		log.Fatalf("invalid amount: %s", amountWeth)
	}

	// Step 1: one-time ERC20 approval of the sell token to the Permit2 contract
	erc20, err := abi.JSON(strings.NewReader(constants.Erc20ABI))
	if err != nil {
		log.Fatalf("failed to parse ERC20 ABI: %v", err)
	}
	allowanceData, err := erc20.Pack("allowance", owner, permit2)
	if err != nil {
		log.Fatalf("failed to pack allowance call: %v", err)
	}
	result, err := client.Wallet.Call(ctx, sellToken, allowanceData)
	if err != nil {
		log.Fatalf("failed to read ERC20 allowance: %v", err)
	}
	if new(big.Int).SetBytes(result).Cmp(amount) < 0 {
		fmt.Println("Sending one-time ERC20 approval to Permit2...")
		approveData, err := erc20.Pack("approve", permit2, constants.Uint256Max)
		if err != nil {
			log.Fatalf("failed to pack approve call: %v", err)
		}
		sendAndWait(ctx, client, sellToken, approveData)
	} else {
		fmt.Println("Permit2 already has a sufficient ERC20 approval")
	}

	// Step 2: standing Permit2 allowance for the router, bounded in amount and time
	allowance, err := orderbook.GetPermit2Allowance(ctx, client.Wallet, owner, sellToken, router)
	if err != nil {
		log.Fatalf("failed to read Permit2 allowance: %v", err)
	}
	nowUnix := big.NewInt(time.Now().Unix())
	if allowance.Amount.Cmp(amount) < 0 || allowance.Expiration.Cmp(nowUnix) <= 0 {
		fmt.Println("Granting the router a standing Permit2 allowance...")
		permit2Abi, err := abi.JSON(strings.NewReader(permit2ApproveABI))
		if err != nil {
			log.Fatalf("failed to parse Permit2 ABI: %v", err)
		}
		expiration := big.NewInt(time.Now().Add(30 * 24 * time.Hour).Unix())
		approveData, err := permit2Abi.Pack("approve", sellToken, router, amount, expiration)
		if err != nil {
			log.Fatalf("failed to pack permit2 approve call: %v", err)
		}
		sendAndWait(ctx, client, permit2, approveData)
	} else {
		fmt.Println("The router already has a sufficient Permit2 allowance")
	}

	// Step 3: swap with the UsePermit2 flag; tokens are pulled through Permit2
	swapData, err := client.GetSwap(ctx, aggregation.GetSwapParams{
		Src:        arbitrumWeth,
		Dst:        arbitrumUsdc,
		Amount:     amountWeth,
		From:       owner.Hex(),
		Slippage:   1,
		UsePermit2: true,
	})
	if err != nil {
		log.Fatalf("failed to get swap data: %v", err)
	}

	tx, err := client.TxBuilder.New().
		SetData(swapData.TxNormalized.Data).
		SetTo(&swapData.TxNormalized.To).
		SetGas(swapData.TxNormalized.Gas).
		SetValue(swapData.TxNormalized.Value).
		Build(ctx)
	if err != nil {
		log.Fatalf("failed to build swap transaction: %v", err)
	}
	signedTx, err := client.Wallet.Sign(tx)
	if err != nil {
		log.Fatalf("failed to sign swap transaction: %v", err)
	}
	if err := client.Wallet.BroadcastTransaction(ctx, signedTx); err != nil {
		log.Fatalf("failed to broadcast swap transaction: %v", err)
	}
	fmt.Printf("Swap sent: https://arbiscan.io/tx/%s\n", signedTx.Hash().Hex())
	waitForReceipt(ctx, client, signedTx.Hash())
	fmt.Println("Swap confirmed")
}

// sendAndWait builds, signs, broadcasts a transaction and waits for its receipt
func sendAndWait(ctx context.Context, client *aggregation.Client, to gethCommon.Address, data []byte) {
	tx, err := client.TxBuilder.New().SetData(data).SetTo(&to).Build(ctx)
	if err != nil {
		log.Fatalf("failed to build transaction: %v", err)
	}
	signedTx, err := client.Wallet.Sign(tx)
	if err != nil {
		log.Fatalf("failed to sign transaction: %v", err)
	}
	if err := client.Wallet.BroadcastTransaction(ctx, signedTx); err != nil {
		log.Fatalf("failed to broadcast transaction: %v", err)
	}
	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())
	waitForReceipt(ctx, client, signedTx.Hash())
	fmt.Println("Transaction confirmed")
}

// waitForReceipt polls for a transaction receipt until it lands or a deadline passes
func waitForReceipt(ctx context.Context, client *aggregation.Client, hash gethCommon.Hash) {
	deadline := time.Now().Add(3 * time.Minute)
	for time.Now().Before(deadline) {
		receipt, err := client.Wallet.TransactionReceipt(ctx, hash)
		if err == nil {
			if receipt.Status != types.ReceiptStatusSuccessful {
				log.Fatalf("transaction reverted: %s", hash.Hex())
			}
			return
		}
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("timed out waiting for receipt: %s", hash.Hex())
}
