package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/1inch/1inch-sdk-go/v4/constants"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/aggregation"
)

/*
This example swaps USDC for WETH on Base with a classic aggregation swap.

The wallet must already have granted the 1inch Aggregation Router an allowance
for the sell token (see the approve example), or use the swap_with_permit or
swap_with_permit2 examples for gasless approvals.

Requires the following environment variables:
  - DEV_PORTAL_TOKEN: 1inch Developer Portal API key
  - WALLET_KEY:       private key (64 hex chars, no 0x prefix)
  - NODE_URL:         RPC endpoint for Base
*/

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
)

const (
	UsdcBase   = "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"
	WethBase   = "0x4200000000000000000000000000000000000006"
	amountUsdc = "100000" // 0.1 USDC (6 decimals)
)

func main() {
	if devPortalToken == "" || privateKey == "" || nodeUrl == "" {
		log.Fatal("set DEV_PORTAL_TOKEN, WALLET_KEY, and NODE_URL to run this example")
	}

	config, err := aggregation.NewConfiguration(aggregation.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: privateKey,
		ChainId:    constants.BaseChainId,
		ApiUrl:     "https://api.1inch.com",
		ApiKey:     devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}
	client, err := aggregation.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()

	swapData, err := client.GetSwap(ctx, aggregation.GetSwapParams{
		Src:      UsdcBase,
		Dst:      WethBase,
		Amount:   amountUsdc,
		From:     client.Wallet.Address().Hex(),
		Slippage: 1, // 1% slippage
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
	fmt.Printf("Swap sent: https://basescan.org/tx/%s\n", signedTx.Hash().Hex())

	deadline := time.Now().Add(3 * time.Minute)
	for time.Now().Before(deadline) {
		receipt, err := client.Wallet.TransactionReceipt(ctx, signedTx.Hash())
		if err == nil {
			if receipt.Status != types.ReceiptStatusSuccessful {
				log.Fatalf("swap transaction reverted: %s", signedTx.Hash().Hex())
			}
			fmt.Println("Swap confirmed")
			return
		}
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("timed out waiting for receipt: %s", signedTx.Hash().Hex())
}
