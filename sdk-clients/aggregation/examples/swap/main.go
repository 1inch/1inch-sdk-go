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
This example demonstrates how to swap tokens on the PolygonChainId network using the 1inch SDK.
The only thing you need to provide is your wallet address, wallet key, and dev portal token.
This can be done through your environment, or you can directly set them in the variables below
*/

var (
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const PolygonPol = "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"
const PolygonUsdc = "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"

func main() {
	config, err := aggregation.NewConfiguration(aggregation.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: privateKey,
		ChainId:    constants.PolygonChainId,
		ApiUrl:     "https://api.1inch.dev",
		ApiKey:     devPortalToken,
	})
	if err != nil {
		log.Fatalf("Failed to create configuration: %v\n", err)
	}
	client, err := aggregation.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}
	ctx := context.Background()

	swapData, err := client.GetSwap(ctx, aggregation.GetSwapParams{
		Src:      PolygonPol,
		Dst:      PolygonUsdc,
		Amount:   "1000000000000000000",
		From:     client.Wallet.Address().Hex(),
		Slippage: 1,
	})
	if err != nil {
		log.Fatalf("Failed to get swap data: %v\n", err)
	}

	tx, err := client.TxBuilder.New().SetData(swapData.TxNormalized.Data).SetTo(&swapData.TxNormalized.To).SetGas(swapData.TxNormalized.Gas).SetValue(swapData.TxNormalized.Value).Build(ctx)
	if err != nil {
		log.Fatalf("Failed to build transaction: %v\n", err)
	}
	signedTx, err := client.Wallet.Sign(tx)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v\n", err)
	}

	err = client.Wallet.BroadcastTransaction(ctx, signedTx)
	if err != nil {
		log.Fatalf("Failed to broadcast transaction: %v\n", err)
	}

	// Waiting for transaction, just an examples of it
	fmt.Printf("Transaction has been broadcast. View it on Polygonscan here: %v\n", fmt.Sprintf("https://polygonscan.com/tx/%v", signedTx.Hash().Hex()))
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
