package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

var (
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const (
	limitOrderHash = "0x585e2e2488d9654926a80356ce135e0ef890cbb6ecec156cad347e091635029e"
	chainId        = 137
)

func main() {
	ctx := context.Background()

	config, err := orderbook.NewConfiguration(orderbook.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: privateKey,
		ChainId:    chainId,
		ApiUrl:     "https://api.1inch.dev",
		ApiKey:     devPortalToken,
	})
	if err != nil {
		log.Fatal(err)
	}
	client, err := orderbook.NewClient(config)

	params := orderbook.GetOrderParams{
		OrderHash:               limitOrderHash,
		SleepBetweenSubrequests: true,
	}
	getOrderResponse, err := client.GetOrderWithSignature(ctx, params)

	takerTraits := orderbook.NewTakerTraits(orderbook.TakerTraitsParams{
		Extension: getOrderResponse.Data.Extension,
	})

	fillOrderData, err := client.GetFillOrderCalldata(getOrderResponse, takerTraits)
	if err != nil {
		log.Fatalf("Failed to get fill order calldata: %v", err)
	}

	aggregationRouter, err := constants.Get1inchRouterFromChainId(chainId)
	if err != nil {
		log.Fatalf("Failed to get 1inch router address: %v", err)
	}
	aggregationRouterAddress := gethCommon.HexToAddress(aggregationRouter)

	tx, err := client.TxBuilder.New().SetData(fillOrderData).SetTo(&aggregationRouterAddress).SetGas(150000).Build(ctx)
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
