package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/orderbook"
	"github.com/1inch/1inch-sdk-go/orderbook/models"
)

/*
This example demonstrates how to fill an order on the Polygon network using the 1inch SDK.
You need to provide your wallet address, wallet key, dev portal token, and the order hash + chain ID of the order you would like to fill.
This can be done through your environment, or you can directly set them in the variables below
*/

var (
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const (
	limitOrderHash = "0x9fd4d0cccb92ae45ee8c86d579334695d525654d8bca2f63f1f85da339093932"
	chainId        = 137
)

func main() {
	ctx := context.Background()

	config, err := orderbook.NewDefaultConfiguration(nodeUrl, privateKey, uint64(chainId), "https://api.1inch.dev", devPortalToken)
	if err != nil {
		log.Fatal(err)
	}
	client, err := orderbook.NewClient(config)

	getOrderRresponse, err := client.GetOrder(ctx, models.GetOrderParams{
		OrderHash: limitOrderHash,
	})

	fillOrderData, err := client.GetFillOrderCalldata(getOrderRresponse)

	fmt.Printf("fillOrderData: %x\n", fillOrderData)

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
