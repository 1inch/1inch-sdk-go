package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/gasprices"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

func main() {
	// Initialize a new configuration using the 1inch SDK for the Aurora chain.
	configLegacyChain, err := gasprices.NewConfiguration(gasprices.ConfigurationParams{
		ChainId: constants.AuroraChainId,
		ApiUrl:  "https://api.1inch.dev",
		ApiKey:  devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create configuration for legacy chain: %v", err)
	}

	// Create a new client with the provided configuration for the legacy chain.
	clientLegacyChain, err := gasprices.NewClient(configLegacyChain)
	if err != nil {
		log.Fatalf("failed to create legacy client: %v", err)
	}

	// Create a new context for the API call.
	ctx := context.Background()

	// Get the legacy gas price.
	gasPriceLegacy, err := clientLegacyChain.GetGasPriceLegacy(ctx)
	if err != nil {
		log.Fatalf("failed to GetGasPriceLegacy: %v", err)
	}

	// Marshal the response to a pretty-printed JSON format.
	gasPriceLegacyIndented, err := json.MarshalIndent(gasPriceLegacy, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal gasPriceLegacy: %v", err)
	}

	// Output the response.
	fmt.Printf("GetGasPriceLegacy: %s\n", gasPriceLegacyIndented)
}
