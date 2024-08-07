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
	config, err := spotprices.NewConfiguration(spotprices.ConfigurationParams{
		ChainId: constants.EthereumChainId,
		ApiUrl:  "https://api.1inch.dev",
		ApiKey:  devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}

	client, err := spotprices.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	customCurrencies, err := client.GetCustomCurrenciesList(ctx)
	if err != nil {
		log.Fatalf("failed to GetCustomCurrenciesList: %v", err)
	}

	customCurrenciesIndented, err := json.MarshalIndent(customCurrencies, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal customCurrencies: %v", err)
	}

	fmt.Printf("GetCustomCurrenciesList: %s\n", customCurrenciesIndented)
}
