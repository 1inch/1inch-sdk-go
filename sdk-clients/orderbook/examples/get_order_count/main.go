package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const (
	chainId = 8453
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

	orderCount, err := client.GetOrderCount(ctx, orderbook.GetOrderCountParams{
		Statuses:   []orderbook.OrderStatus{orderbook.ValidOrders},
		MakerAsset: "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
		TakerAsset: "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359",
	})
	if err != nil {
		log.Fatalf("Failed to get order count: %v", err)
	}

	fmt.Printf("Order count: %v\n", orderCount.Count)
}
