package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/spotprices"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

func main() {
	// Initialize a new configuration using the 1inch SDK for the Ethereum chain.
	config, err := spotprices.NewConfiguration(spotprices.ConfigurationParams{
		ChainId: constants.EthereumChainId,
		ApiUrl:  "https://api.1inch.dev",
		ApiKey:  devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}

	// Create a new client with the provided configuration.
	client, err := spotprices.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	// Create a new context for the API call.
	ctx := context.Background()

	// Get prices for whitelisted tokens.
	whitelistedTokensPrices, err := client.GetPricesForWhitelistedTokens(ctx, spotprices.GetWhitelistedTokensPricesParams{
		Currency: spotprices.GetWhitelistedTokensPricesParamsCurrency(spotprices.USD),
	})
	if err != nil {
		log.Fatalf("failed to GetWhitelistedTokensPrices: %v", err)
	}

	// Marshal the response to a pretty-printed JSON format.
	whitelistedTokensPricesIndented, err := json.MarshalIndent(whitelistedTokensPrices, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal whitelistedTokensPrices: %v", err)
	}

	// Output the response.
	fmt.Printf("GetWhitelistedTokensPrices: %s\n", whitelistedTokensPricesIndented)
}
