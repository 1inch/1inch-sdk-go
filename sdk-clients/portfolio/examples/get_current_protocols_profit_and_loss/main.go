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

	timerange := portfolio.N1day

	response, err := client.GetProtocolsProfitAndLoss(ctx, portfolio.GetProfitAndLossPortfolioV4OverviewProtocolsProfitAndLossGetParams{
		Addresses: []string{"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"},
		ChainId:   1,
		Timerange: &timerange,
	})
	if err != nil {
		log.Fatalf("failed to get response: %v", err)
	}

	responseIndented, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("failed to MarshalIndent: %v", err)
	}

	fmt.Printf("Response: %s\n", responseIndented)
}
