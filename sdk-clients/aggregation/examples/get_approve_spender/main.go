package main

import (
	"context"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/aggregation"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
	walletAddress  = os.Getenv("WALLET_ADDRESS")
)

const PolygonWeth = "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619"

func main() {
	config, err := aggregation.NewConfigurationAPI(
		constants.PolygonChainId,
		"https://api.1inch.dev",
		devPortalToken,
	)
	if err != nil {
		log.Fatalf("Failed to create configuration: %v\n", err)
	}
	client, err := aggregation.NewClientOnlyAPI(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}
	ctx := context.Background()

	response, err := client.GetApproveAllowance(ctx, aggregation.GetAllowanceParams{
		TokenAddress:  PolygonWeth,
		WalletAddress: walletAddress,
	})
	if err != nil {
		log.Fatalf("Failed to get tokens: %v\n", err)
	}

	log.Printf("Allowance: %v\n", response.Allowance)
}
