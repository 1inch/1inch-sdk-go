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
	limitOrderHash = "0xddf0907b11b3dba4591372e201117f09f9cdc7ccfff175397bc57eeb7043a59f"
	chainId        = 137
)

func main() {
	ctx := context.Background()

	config, err := orderbook.NewDefaultConfiguration(nodeUrl, privateKey, uint64(chainId), "https://api.1inch.dev", devPortalToken)
	if err != nil {
		log.Fatal(err)
	}
	client, err := orderbook.NewClient(config)

	getOrderResponse, err := client.GetOrder(ctx, orderbook.GetOrderParams{
		OrderHash: limitOrderHash,
	})

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
