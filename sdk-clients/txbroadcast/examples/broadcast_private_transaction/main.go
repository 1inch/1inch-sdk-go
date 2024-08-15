package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/txbroadcast"
)

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

	broadcastPrivateResponse, err := client.BroadcastPrivateTransaction(ctx, txbroadcast.BroadcastRequest{
		RawTransaction: "<YOUR RAW TX here that you want to send to private mempool>",
	})
	if err != nil {
		log.Fatalf("failed to BroadcastPrivateTransaction: %v", err)
	}

	broadcastPrivateResponseIndented, err := json.MarshalIndent(broadcastPrivateResponse, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal broadcastPrivateResponse: %v", err)
	}

	fmt.Printf("BroadcastPrivateTransaction: %s\n", broadcastPrivateResponseIndented)
}
