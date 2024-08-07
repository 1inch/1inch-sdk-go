package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/sdk-clients/nft"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

func main() {
	// Initialize a new configuration using the 1inch SDK.
	config, err := nft.NewConfiguration(nft.ConfigurationParams{
		ApiUrl: "https://api.1inch.dev",
		ApiKey: devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}

	// Create a new client with the provided configuration.
	client, err := nft.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	// Create a new context for the API call.
	ctx := context.Background()

	// Get the supported chains.
	chains, err := client.GetSupportedChains(ctx)
	if err != nil {
		log.Fatalf("failed to GetSupportedChains: %v", err)
	}

	// Marshal the response to a pretty-printed JSON format.
	chainsIndented, err := json.MarshalIndent(chains, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal chains: %v", err)
	}

	// Output the response.
	fmt.Printf("GetSupportedChains: %s\n", chainsIndented)
}
