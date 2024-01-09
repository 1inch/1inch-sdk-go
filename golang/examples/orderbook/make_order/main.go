package main

import (
	"context"
	"log"
	"os"

	"github.com/1inch/1inch-sdk/golang/client"
	"github.com/1inch/1inch-sdk/golang/client/orderbook"
	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/addresses"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/contracts"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/tokens"
)

func main() {

	// Build the config for the client
	config := client.Config{
		DevPortalApiKey:    os.Getenv("DEV_PORTAL_TOKEN"),
		WalletKey:          os.Getenv("WALLET_KEY"),
		LimitOrderContract: contracts.AggregationRouterV5,
		ChainId:            chains.Polygon,
	}

	// Create the 1inch client
	c, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Build the config for the orders request
	createOrderParams := orderbook.OrderRequest{
		SourceWallet: os.Getenv("WALLET_ADDRESS"),
		FromToken:    tokens.PolygonDai,
		ToToken:      tokens.PolygonWeth,
		TakingAmount: 1000000000000000000,
		MakingAmount: 1000000000000000000,
		Receiver:     addresses.Zero,
	}

	// Execute orders request
	createOrderResponse, _, err := c.Orderbook.CreateOrder(context.Background(), createOrderParams)
	if err != nil {
		log.Fatalf("Failed to get quote: %v", err)
	}

	helpers.PrettyPrintStruct(createOrderResponse)
}
