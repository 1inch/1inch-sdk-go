package main

import (
	"context"
	"fmt"
	"os"

	"github.com/1inch/1inch-sdk-go/constants"
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
		Wallets:     []string{"0x1C17622cfa9B6fD2043A76DfC39A5B5a109aa708", "0x28C6c06298d514Db089934071355E5743bf21d60"},
		FilterEmpty: true,
		Spender:     "0x58b6a8a3302369daec383334672404ee733ab239",
	})
	if err != nil {
		return
	}

	fmt.Println(balancesAndAllowances)
}
