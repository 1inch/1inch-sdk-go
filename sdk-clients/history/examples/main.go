package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/history"
)

/*
This example demonstrates how to swap tokens on the EthereumChainId network using the 1inch SDK.
The only thing you need to provide is your wallet address, wallet key, and dev portal token.
This can be done through your environment, or you can directly set them in the variables below
*/

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

func main() {
	config, err := history.NewConfiguration(history.ConfigurationParams{
		ApiUrl: "https://api.1inch.dev",
		ApiKey: devPortalToken,
	})
	if err != nil {
		log.Fatal("failed to create configuration: %w", err)
	}
	client, err := history.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()

	historyEvents, err := client.GetHistoryEventsByAddress(ctx, history.EventsByAddressParams{
		Address: "0x266E77cE9034a023056ea2845CB6A20517F6FDB7",
		ChainId: constants.EthereumChainId,
	})
	if err != nil {
		log.Fatalf("failed to GetHistoryEventsByAddress: %v", err)
	}

	fmt.Println("GetHistoryEventsByAddress:", historyEvents)
	time.Sleep(time.Second)
}
