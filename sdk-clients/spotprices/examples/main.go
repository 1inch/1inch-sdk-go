package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/spotprices"
)

/*
This example demonstrates how to swap tokens on the EthereumChainId network using the 1inch SDK.
The only thing you need to provide is your wallet address, wallet key, and dev portal token.
This can be done through your environment, or you can directly set them in the variables below
*/

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const (
	mainWalletAddress      = "0x1C17622cfa9B6fD2043A76DfC39A5B5a109aa708"
	secondaryWalletAddress = "0x28C6c06298d514Db089934071355E5743bf21d60"

	tokenAddress1 = "0x0d8775f648430679a709e98d2b0cb6250d2887ef"
	tokenAddress2 = "0x58b6a8a3302369daec383334672404ee733ab239"

	spender     = "0x58b6a8a3302369daec383334672404ee733ab239"
	spenderInch = "0x111111125421ca6dc452d289314280a0f8842a65"
)

func main() {
	config, err := spotprices.NewConfiguration(constants.EthereumChainId, "https://api.1inch.dev", devPortalToken)
	if err != nil {
		return
	}
	client, err := spotprices.NewClient(config)
	if err != nil {
		return
	}
	ctx := context.Background()

	whitelistedTokensPrices, err := client.GetPricesForWhitelistedTokens(ctx, spotprices.GetWhitelistedTokensPricesParams{
		Currency: spotprices.GetWhitelistedTokensPricesParamsCurrency(spotprices.USD),
	})
	if err != nil {
		fmt.Println("failed to GetWhitelistedTokensPricesParams: %w", err)
		return
	}

	fmt.Println("GetWhitelistedTokensPricesParams:", whitelistedTokensPrices)
	time.Sleep(time.Second)
}
