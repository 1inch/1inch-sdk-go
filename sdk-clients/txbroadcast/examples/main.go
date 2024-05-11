package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/txbroadcast"
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
	config, err := txbroadcast.NewConfiguration(constants.EthereumChainId, "https://api.1inch.dev", devPortalToken)
	if err != nil {
		return
	}
	client, err := txbroadcast.NewClient(config)
	if err != nil {
		return
	}
	ctx := context.Background()

	broadcastPublicResponse, err := client.BroadcastPublicTransaction(ctx, txbroadcast.BroadcastRequest{
		RawTransaction: "",
	})
	if err != nil {
		fmt.Println("failed to BroadcastPublicTransaction: %w", err)
		return
	}

	fmt.Println("BroadcastPublicTransaction:", broadcastPublicResponse)
	time.Sleep(time.Second)
}
