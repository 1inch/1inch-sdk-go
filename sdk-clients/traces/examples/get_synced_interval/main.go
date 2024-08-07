package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/traces"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

func main() {
	// Initialize a new configuration using the 1inch SDK for the Ethereum chain.
	config, err := traces.NewConfiguration(constants.EthereumChainId, "https://api.1inch.dev", devPortalToken)
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}

	// Create a new client with the provided configuration.
	client, err := traces.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	// Create a new context for the API call.
	ctx := context.Background()

	// Get the synced interval.
	interval, err := client.GetSyncedInterval(ctx)
	if err != nil {
		log.Fatalf("failed to GetSyncedInterval: %v", err)
	}

	// Marshal the response to a pretty-printed JSON format.
	intervalIndented, err := json.MarshalIndent(interval, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal interval: %v", err)
	}

	// Output the response.
	fmt.Printf("GetSyncedInterval: %s\n", intervalIndented)
}
