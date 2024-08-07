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

const (
	tokenAddress1 = "0x0d8775f648430679a709e98d2b0cb6250d2887ef"
	tokenAddress2 = "0x58b6a8a3302369daec383334672404ee733ab239"
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

	// Get prices for requested tokens.
	requestedTokensPrices, err := client.GetPricesForRequestedTokens(ctx, spotprices.GetPricesRequestDto{
		Currency: spotprices.GetPricesRequestDtoCurrency(spotprices.USD),
		Tokens:   []string{tokenAddress1, tokenAddress2},
	})
	if err != nil {
		log.Fatalf("failed to GetPricesForRequestedTokens: %v", err)
	}

	// Marshal the response to a pretty-printed JSON format.
	requestedTokensPricesIndented, err := json.MarshalIndent(requestedTokensPrices, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal requestedTokensPrices: %v", err)
	}

	// Output the response.
	fmt.Printf("GetPricesForRequestedTokens: %s\n", requestedTokensPricesIndented)
}
