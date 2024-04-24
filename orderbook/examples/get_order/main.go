package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/orderbook"
	"github.com/1inch/1inch-sdk-go/orderbook/models"
)

var (
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const (
	limitOrderHash = "0x20b8c5bf0f381c4cfa2b320d7366954ec80472593092fc5a799fa11f17e52daa"
	chainId        = 137
)

func main() {
	// use the swap example program as a template for building the limit order fill flow

	ctx := context.Background()

	config, err := orderbook.NewDefaultConfiguration(nodeUrl, privateKey, uint64(chainId), "https://api.1inch.dev", devPortalToken)
	if err != nil {
		log.Fatal(err)
	}
	client, err := orderbook.NewClient(config)

	getOrderRresponse, err := client.GetOrder(ctx, models.GetOrderParams{
		ChainId:   chainId,
		OrderHash: limitOrderHash,
	})

	fmt.Printf("getOrderRresponse: %+v\n", getOrderRresponse)

}
