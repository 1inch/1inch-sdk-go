package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/aggregation"
)

/*
This example demonstrates how to swap tokens on the Base network using the 1inch SDK.
The only thing you need to provide is your wallet address, wallet key, and dev portal token.
This can be done through your environment, or you can directly set them in the variables below
*/

var (
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const (
	UsdcBase   = "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"
	WethBase   = "0x4200000000000000000000000000000000000006"
	amountWeth = "10000000000000000"
	amountUsdc = "100000"
)

func main() {
	config, err := aggregation.NewConfiguration(aggregation.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: privateKey,
		ChainId:    constants.BaseChainId,
		ApiUrl:     "https://api.1inch.dev",
		ApiKey:     devPortalToken,
	})
	if err != nil {
		log.Fatalf("Failed to create configuration: %v", err)
	}
	client, err := aggregation.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	ctx := context.Background()

	swapData, err := client.GetSwap(ctx, aggregation.GetSwapParams{
		Src:      UsdcBase,
		Dst:      WethBase,
		Amount:   amountUsdc,
		From:     client.Wallet.Address().Hex(),
		Slippage: 1,
	})
	if err != nil {
		log.Fatalf("Failed to get swap data: %v", err)
	}

	tx, err := client.TxBuilder.New().SetData(swapData.TxNormalized.Data).SetTo(&swapData.TxNormalized.To).SetGas(swapData.TxNormalized.Gas).SetValue(swapData.TxNormalized.Value).Build(ctx)
	if err != nil {
		log.Fatalf("Failed to build transaction: %v", err)
	}
	signedTx, err := client.Wallet.Sign(tx)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	err = client.Wallet.BroadcastTransaction(ctx, signedTx)
	if err != nil {
		log.Fatalf("Failed to broadcast transaction: %v", err)
	}

	// Waiting for transaction, just an example of it
	fmt.Printf("Transaction has been broadcast. View it on Basescan here: https://basescan.org/tx/%s\n", signedTx.Hash().Hex())
	for {
		receipt, err := client.Wallet.TransactionReceipt(ctx, signedTx.Hash())
		if receipt != nil {
			fmt.Println("Transaction complete!")
			return
		}
		if err != nil {
			fmt.Println("Waiting for transaction to be mined")
		}
		select {
		case <-time.After(1000 * time.Millisecond): // check again after a delay
		case <-ctx.Done():
			fmt.Println("Context cancelled")
			return
		}
	}
}
