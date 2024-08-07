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
	// Initialize a new configuration using the 1inch SDK for the Ethereum chain.
	config, err := txbroadcast.NewConfiguration(txbroadcast.ConfigurationParams{
		ChainId: constants.EthereumChainId,
		ApiUrl:  "https://api.1inch.dev",
		ApiKey:  devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}

	// Create a new client with the provided configuration.
	client, err := txbroadcast.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	// Create a new context for the API call.
	ctx := context.Background()

	// Broadcast a public transaction.
	broadcastPublicResponse, err := client.BroadcastPublicTransaction(ctx, txbroadcast.BroadcastRequest{
		RawTransaction: "<YOUR RAW TX here>",
	})
	if err != nil {
		log.Fatalf("failed to BroadcastPublicTransaction: %v", err)
	}

	// Marshal the response to a pretty-printed JSON format.
	broadcastPublicResponseIndented, err := json.MarshalIndent(broadcastPublicResponse, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal broadcastPublicResponse: %v", err)
	}

	// Output the response.
	fmt.Printf("BroadcastPublicTransaction: %s\n", broadcastPublicResponseIndented)
}
