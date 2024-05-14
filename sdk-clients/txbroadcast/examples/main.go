package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/txbroadcast"
)

/*
This example demonstrates how to swap tokens on the EthereumChainId network using the 1inch SDK.
The only thing you need to provide is your wallet address, wallet key, and dev portal token.
This can be done through your environment, or you can directly set them in the variables below
*/

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

func main() {
	config, err := txbroadcast.NewConfiguration(txbroadcast.ConfigurationParams{
		ChainId: constants.EthereumChainId,
		ApiUrl:  "https://api.1inch.dev",
		ApiKey:  devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}
	client, err := txbroadcast.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()

	broadcastPublicResponse, err := client.BroadcastPublicTransaction(ctx, txbroadcast.BroadcastRequest{
		RawTransaction: "<YOOR RAW TX here>",
	})
	if err != nil {
		log.Fatalf("failed to BroadcastPublicTransaction: %v", err)
	}

	fmt.Println("BroadcastPublicTransaction:", broadcastPublicResponse)
	time.Sleep(time.Second)

	broadcastPrivateResponse, err := client.BroadcastPrivateTransaction(ctx, txbroadcast.BroadcastRequest{
		RawTransaction: "<YOOR RAW TX here that you want to send to private mempool>",
	})
	if err != nil {
		log.Fatalf("failed to BroadcastPrivateTransaction: %v", err)
	}

	fmt.Println("BroadcastPrivateTransaction:", broadcastPrivateResponse)
	time.Sleep(time.Second)
}
