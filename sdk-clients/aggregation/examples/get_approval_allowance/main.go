package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/v4/constants"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/aggregation"
)

/*
This example fetches the current allowance a wallet has granted to the 1inch
Aggregation Router for a token.

Requires the following environment variables:
  - DEV_PORTAL_TOKEN: 1inch Developer Portal API key
  - WALLET_ADDRESS:   address of the wallet to check
*/

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
	walletAddress  = os.Getenv("WALLET_ADDRESS")
)

const PolygonWeth = "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619"

func main() {
	if devPortalToken == "" || walletAddress == "" {
		log.Fatal("set DEV_PORTAL_TOKEN and WALLET_ADDRESS to run this example")
	}

	config, err := aggregation.NewConfigurationAPI(
		constants.PolygonChainId,
		"https://api.1inch.com",
		devPortalToken,
	)
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}
	client, err := aggregation.NewClientOnlyAPI(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()

	response, err := client.GetApproveAllowance(ctx, aggregation.GetAllowanceParams{
		TokenAddress:  PolygonWeth,
		WalletAddress: walletAddress,
	})
	if err != nil {
		log.Fatalf("failed to get approve allowance: %v", err)
	}

	fmt.Printf("Router allowance for WETH: %s\n", response.Allowance)
}
