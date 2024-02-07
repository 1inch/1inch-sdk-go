package main

import (
	"log"
	"os"

	"github.com/1inch/1inch-sdk/golang/client"
	"github.com/1inch/1inch-sdk/golang/client/swap"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/tokens"
)

func main() {

	// Build the config for the client
	config := client.Config{
		DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
		Web3HttpProviders: []client.Web3ProviderConfig{
			{
				ChainId: chains.Polygon,
				Url:     os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY_POLYGON"),
			},
		},
	}

	// Create the 1inch client
	c, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Build the config for the swap request
	swapParams := swap.SwapTokensParams{
		ApprovalType:  swap.PermitIfPossible,
		SkipWarnings:  false,
		PublicAddress: os.Getenv("WALLET_ADDRESS"),
		WalletKey:     os.Getenv("WALLET_KEY"),
		ChainId:       chains.Polygon,
		AggregationControllerGetSwapParams: swap.AggregationControllerGetSwapParams{
			Src:             tokens.PolygonFrax,
			Dst:             tokens.PolygonUsdc,
			From:            os.Getenv("WALLET_ADDRESS"),
			Amount:          "10000000000000000",
			Slippage:        0.5,
			DisableEstimate: true,
		},
	}

	err = c.Actions.SwapTokens(swapParams)
	if err != nil {
		log.Fatalf("Failed to swap tokens: %v", err)
	}
}
