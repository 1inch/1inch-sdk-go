package main

import (
	"context"
	"fmt"
	"os"

	"github.com/1inch/1inch-sdk-go/aggregation"
	"github.com/1inch/1inch-sdk-go/constants"
)

/*
This example demonstrates how to swap tokens on the PolygonChainId network using the 1inch SDK.
The only thing you need to provide is your wallet address, wallet key, and dev portal token.
This can be done through your environment, or you can directly set them in the variables below
*/

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const (
	PolygonDai  = "0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063"
	PolygonWeth = "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619"
)

func main() {
	config, err := aggregation.NewConfigurationAPI(constants.PolygonChainId, "https://api.1inch.dev", devPortalToken)
	if err != nil {
		return
	}
	client, err := aggregation.NewClientOnlyAPI(config)
	if err != nil {
		return
	}
	ctx := context.Background()

	quote, err := client.GetQuote(ctx, aggregation.GetQuoteParams{
		Src:    PolygonDai,
		Dst:    PolygonWeth,
		Amount: "10000000000000000",
	})

	if err != nil {
		return
	}

	fmt.Printf("Quote Amount: %+v\n", quote.ToAmount)
}
