package main

import (
	"log"
	"os"

	"1inch-sdk-golang/actions"
	"1inch-sdk-golang/client"
	"1inch-sdk-golang/client/swap"
	"1inch-sdk-golang/helpers"
	"1inch-sdk-golang/helpers/consts/amounts"
	"1inch-sdk-golang/helpers/consts/chains"
	"1inch-sdk-golang/helpers/consts/tokens"
)

func main() {

	// Build the config for the client
	config := client.Config{
		DevPortalApiKey:            os.Getenv("DEV_PORTAL_TOKEN"),
		WalletKey:                  os.Getenv("WALLET_KEY"),
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

	err = actions.SwapTokens(c, swapParams)
	if err != nil {
		log.Fatalf("Failed to swap tokens: %v", err)
	}
}
