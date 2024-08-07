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

	// Get the list of custom currencies.
	customCurrencies, err := client.GetCustomCurrenciesList(ctx)
	if err != nil {
		log.Fatalf("failed to GetCustomCurrenciesList: %v", err)
	}

	// Marshal the response to a pretty-printed JSON format.
	customCurrenciesIndented, err := json.MarshalIndent(customCurrencies, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal customCurrencies: %v", err)
	}

	// Output the response.
	fmt.Printf("GetCustomCurrenciesList: %s\n", customCurrenciesIndented)
}
