package main

import (
	"context"
	"fmt"
	"os"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/aggregation"
	"github.com/1inch/1inch-sdk-go/sdk-clients/balances"
)

/*
This examples demonstrates how to swap tokens on the PolygonChainId network using the 1inch SDK.
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
	config, err := balances.NewConfiguration(constants.PolygonChainId, "https://api.1inch.dev", devPortalToken)
	if err != nil {
		return
	}
	client, err := balances.NewClient(config)
	if err != nil {
		return
	}
	ctx := context.Background()

	balancesAndAllowances, err := client.GetBalancesAndAllowances(ctx, balances.BalancesAndAllowancesParams{
		Wallets:     nil,
		FilterEmpty: false,
		Spender:     "",
	})
	if err != nil {
		return
	}

	fmt.Println(balancesAndAllowances)
}
