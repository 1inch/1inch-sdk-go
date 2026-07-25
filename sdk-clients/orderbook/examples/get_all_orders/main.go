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

	orders, err := client.GetAllOrders(ctx, orderbook.GetAllOrdersParams{
		LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams: orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{
			Page:     0,
			Limit:    2,
			Statuses: []float32{1},
		},
	})
	if err != nil {
		log.Fatalf("failed to get all orders: %v", err)
	}

	ordersIndented, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal response: %v", err)
	}

	fmt.Printf("Orders: %s\n", ordersIndented)
}
