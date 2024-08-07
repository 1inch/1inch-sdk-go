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
	tokenAddress3 = "0x320623b8e4ff03373931769a31fc52a4e78b5d70"
	tokenAddress4 = "0x71ab77b7dbb4fa7e017bc15090b2163221420282"
	tokenAddress5 = "0x256d1fce1b1221e8398f65f9b36033ce50b2d497"
	tokenAddress6 = "0x85f17cf997934a597031b2e18a9ab6ebd4b9f6a4"
	tokenAddress7 = "0x55c08ca52497e2f1534b59e2917bf524d4765257"
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

	// Get prices for a larger set of requested tokens.
	requestedTokensPricesLarge, err := client.GetPricesForRequestedTokensLarge(ctx, spotprices.GetPricesRequestDto{
		Currency: spotprices.GetPricesRequestDtoCurrency(spotprices.USD),
		Tokens:   []string{tokenAddress1, tokenAddress2, tokenAddress3, tokenAddress4, tokenAddress5, tokenAddress6, tokenAddress7},
	})
	if err != nil {
		log.Fatalf("failed to GetPricesForRequestedTokensLarge: %v", err)
	}

	// Marshal the response to a pretty-printed JSON format.
	requestedTokensPricesLargeIndented, err := json.MarshalIndent(requestedTokensPricesLarge, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal requestedTokensPricesLarge: %v", err)
	}

	// Output the response.
	fmt.Printf("GetPricesForRequestedTokensLarge: %s\n", requestedTokensPricesLargeIndented)
}
