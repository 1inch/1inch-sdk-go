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

	broadcastPublicResponse, err := client.BroadcastPublicTransaction(ctx, txbroadcast.BroadcastRequest{
		RawTransaction: "<YOUR RAW TX here>",
	})
	if err != nil {
		log.Fatalf("failed to BroadcastPublicTransaction: %v", err)
	}

	broadcastPublicResponseIndented, err := json.MarshalIndent(broadcastPublicResponse, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal broadcastPublicResponse: %v", err)
	}

	fmt.Printf("BroadcastPublicTransaction: %s\n", broadcastPublicResponseIndented)
}
