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
	// Initialize a new configuration using the 1inch SDK for the Ethereum chain.
	config, err := gasprices.NewConfiguration(gasprices.ConfigurationParams{
		ChainId: constants.EthereumChainId,
		ApiUrl:  "https://api.1inch.dev",
		ApiKey:  devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}

	// Create a new client with the provided configuration.
	client, err := gasprices.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	// Create a new context for the API call.
	ctx := context.Background()

	// Get the EIP-1559 gas price.
	gasPriceEIP1559, err := client.GetGasPriceEIP1559(ctx)
	if err != nil {
		log.Fatalf("failed to GetGasPriceEIP1559: %v", err)
	}

	// Marshal the response to a pretty-printed JSON format.
	gasPriceEIP1559Indented, err := json.MarshalIndent(gasPriceEIP1559, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal gasPriceEIP1559: %v", err)
	}

	// Output the response.
	fmt.Printf("GetGasPriceEIP1559: %s\n", gasPriceEIP1559Indented)
}
