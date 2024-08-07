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

	// Broadcast a private transaction.
	broadcastPrivateResponse, err := client.BroadcastPrivateTransaction(ctx, txbroadcast.BroadcastRequest{
		RawTransaction: "<YOUR RAW TX here that you want to send to private mempool>",
	})
	if err != nil {
		log.Fatalf("failed to BroadcastPrivateTransaction: %v", err)
	}

	// Marshal the response to a pretty-printed JSON format.
	broadcastPrivateResponseIndented, err := json.MarshalIndent(broadcastPrivateResponse, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal broadcastPrivateResponse: %v", err)
	}

	// Output the response.
	fmt.Printf("BroadcastPrivateTransaction: %s\n", broadcastPrivateResponseIndented)
}
