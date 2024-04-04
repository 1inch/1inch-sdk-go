package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/1inch/1inch-sdk-go/aggregation"
	"github.com/1inch/1inch-sdk-go/aggregation/models"
)

/*
This example demonstrates how to swap tokens on the Polygon network using the 1inch SDK.
The only thing you need to provide is your wallet address, wallet key, and dev portal token.
This can be done through your environment, or you can directly set them in the variables below
*/

var (
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

func main() {
	config, err := aggregation.NewDefaultConfiguration(nodeUrl, privateKey, 137, "https://api.1inch.dev", devPortalToken)
	if err != nil {
		return
	}
	client, err := aggregation.NewClient(config)

	ctx := context.Background()

	swapData, err := client.GetSwap(ctx, models.AggregationControllerGetSwapParams{
		Src:      "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359",
		Dst:      "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270",
		Amount:   "1000",
		From:     client.Wallet.Address().Hex(),
		Slippage: 1,
	})
	if err != nil {
		fmt.Printf("Failed to get swap data: %v\n", err)
		return
	}

	nonce, err := client.Wallet.Nonce(ctx)
	if err != nil {
		fmt.Printf("Failed to get nonce: %v\n", err)
		return
	}

	gasTip, err := client.Wallet.GetGasTipCap(ctx)
	if err != nil {
		fmt.Printf("Failed to get gas tip: %v\n", err)
		return
	}

	gasFee, err := client.Wallet.GetGasFeeCap(ctx)
	if err != nil {
		fmt.Printf("Failed to get gas fee: %v\n", err)
		return
	}

	tx, err := client.BuildSwapTransaction(swapData, nonce, gasTip, gasFee)
	if err != nil {
		fmt.Printf("Failed to build transaction: %v\n", err)
		return
	}

	signedTx, err := client.Wallet.Sign(tx)
	if err != nil {
		fmt.Printf("Failed to sign transaction: %v\n", err)
		return
	}

	err = client.Wallet.BroadcastTransaction(ctx, signedTx)
	if err != nil {
		fmt.Printf("Failed to broadcast transaction: %v\n", err)
		return
	}

	// Waiting for transaction, just an example of it
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
