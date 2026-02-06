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
)

func main() {
	config, err := aggregation.NewConfigurationAPI(
		constants.PolygonChainId,
		"https://api.1inch.dev",
		devPortalToken,
	)
	if err != nil {
		log.Fatalf("Failed to create configuration: %v", err)
	}
	client, err := aggregation.NewClientOnlyAPI(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	ctx := context.Background()

	response, err := client.GetApproveSpender(ctx)
	if err != nil {
		log.Fatalf("Failed to get approve allowance: %v", err)
	}

	log.Printf("1inch router: %v", response.Address)
}
