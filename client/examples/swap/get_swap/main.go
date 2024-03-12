package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/client"
	"github.com/1inch/1inch-sdk-go/client/models"
	"github.com/1inch/1inch-sdk-go/helpers/consts/amounts"
	"github.com/1inch/1inch-sdk-go/helpers/consts/chains"
	"github.com/1inch/1inch-sdk-go/helpers/consts/tokens"
	"github.com/1inch/1inch-sdk-go/helpers/consts/web3providers"
)

func main() {

	// Build the config for the client
	config := models.ClientConfig{
		DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
		Web3HttpProviders: []models.Web3Provider{
			{
				ChainId: chains.Polygon,
				Url:     web3providers.Polygon,
			},
		},
	}

	// Create the 1inch client
	c, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Build the config for the swap request
	swapParams := models.GetSwapParams{
		ChainId:      chains.Polygon,
		SkipWarnings: false,
		AggregationControllerGetSwapParams: models.AggregationControllerGetSwapParams{
			Src:             tokens.PolygonFrax,
			Dst:             tokens.PolygonWeth,
			From:            os.Getenv("WALLET_ADDRESS"),
			Amount:          amounts.Ten16,
			DisableEstimate: true,
			Slippage:        0.5,
		},
	}

	swapData, _, err := c.SwapApi.GetSwap(context.Background(), swapParams)
	if err != nil {
		log.Fatalf("Failed to swap tokens: %v", err)
	}

	swapDataRawIndented, err := json.MarshalIndent(swapData, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal swap data: %v", err)
	}

	fmt.Printf("%s\n", string(swapDataRawIndented))
}
