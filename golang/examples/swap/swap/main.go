package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk/golang/client"
	"github.com/1inch/1inch-sdk/golang/client/swap"
	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/amounts"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/tokens"
)

func main() {

	// Build the config for the client
	config := client.Config{
		DevPortalApiKey:            os.Getenv("DEV_PORTAL_TOKEN"),
		Web3HttpProviderUrlWithKey: os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY_POLYGON"),
		ChainId:                    chains.Polygon,
	}

	// Create the 1inch client
	c, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Build the config for the swap request
	swapParams := swap.AggregationControllerGetSwapParams{
		Src:             tokens.PolygonFrax,
		Dst:             tokens.PolygonWeth,
		From:            os.Getenv("WALLET_ADDRESS"),
		Amount:          amounts.Ten16,
		DisableEstimate: helpers.GetPtr(true),
	}

	swapData, _, err := c.Swap.GetSwapData(context.Background(), swapParams, false)
	if err != nil {
		log.Fatalf("Failed to swap tokens: %v", err)
	}

	fmt.Printf("\nContract to send transaction to: %v\n", swapData.Tx.To)
	fmt.Printf("Transaction data: %v\n", swapData.Tx.Data)
}
