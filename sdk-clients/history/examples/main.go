package main

import (
	"context"
	"fmt"
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
	config, err := history.NewConfiguration("https://api.1inch.dev", devPortalToken)
	if err != nil {
		return
	}
	client, err := history.NewClient(config)
	if err != nil {
		return
	}
	ctx := context.Background()

	historyEvents, err := client.GetHistoryEventsByAddress(ctx, history.HistoryEventsByAddressParams{
		Address: "0x266E77cE9034a023056ea2845CB6A20517F6FDB7",
		ChainId: constants.EthereumChainId,
	})
	if err != nil {
		fmt.Println("failed to GetHistoryEventsByAddress: %w", err)
		return
	}

	fmt.Println("GetHistoryEventsByAddress:", historyEvents)
	time.Sleep(time.Second)
}
