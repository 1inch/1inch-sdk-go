package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/sdk-clients/portfolio"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

func main() {
	config, err := portfolio.NewConfiguration(portfolio.ConfigurationParams{
		ApiUrl: "https://api.1inch.dev",
		ApiKey: devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}
	client, err := portfolio.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()

	isServiceAvailableResponse, err := client.IsServiceAvailable(ctx)
	if err != nil {
		log.Fatalf("failed to IsServiceAvailable: %v", err)
	}

	tokensProfitLossResponseIndented, err := json.MarshalIndent(isServiceAvailableResponse, "", "  ")
	if err != nil {
		log.Fatalf("failed to MarshalIndent: %v", err)
	}

	fmt.Printf("GetProtocolsDetails: %s\n", tokensProfitLossResponseIndented)
}
