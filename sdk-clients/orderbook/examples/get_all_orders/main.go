package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const (
	chainId = 137
)

func main() {
	ctx := context.Background()

	config, err := orderbook.NewConfigurationAPI(chainId, "https://api.1inch.dev", devPortalToken)
	if err != nil {
		log.Fatal(err)
	}
	client, err := orderbook.NewClientOnlyAPI(config)
	if err != nil {
		log.Fatal(err)
	}

	orders, err := client.GetAllOrders(ctx, orderbook.GetAllOrdersParams{
		LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams: orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{
			Page:     0,
			Limit:    2,
			Statuses: []float32{1},
		},
	})
	if err != nil {
		log.Fatalf("Failed to get order count: %v", err)
	}

	ordersIndented, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal response: %v", err)
	}

	fmt.Printf("Orders: %s\n", ordersIndented)
}
