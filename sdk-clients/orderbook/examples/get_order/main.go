package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/orderbook"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const (
	chainId = 137
)

func main() {
	if devPortalToken == "" {
		log.Fatal("set DEV_PORTAL_TOKEN to run this example")
	}

	ctx := context.Background()

	config, err := orderbook.NewConfigurationAPI(chainId, "https://api.1inch.com", devPortalToken)
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}
	client, err := orderbook.NewClientOnlyAPI(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	order, err := client.GetOrder(ctx, orderbook.GetOrderParams{
		OrderHash: "0x887b4e1b5ab0fd51884f25234fb725307056e0b2d33b881ea9013f9258fb444a",
	})
	if err != nil {
		log.Fatalf("failed to get order: %v", err)
	}

	orderIndented, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal response: %v", err)
	}

	fmt.Printf("Order: %s\n", orderIndented)
}
