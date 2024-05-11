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
	tokenAddress1 = "0x0d8775f648430679a709e98d2b0cb6250d2887ef"
	tokenAddress2 = "0x58b6a8a3302369daec383334672404ee733ab239"
	tokenAddress3 = "0x320623b8e4ff03373931769a31fc52a4e78b5d70"
	tokenAddress4 = "0x71ab77b7dbb4fa7e017bc15090b2163221420282"
	tokenAddress5 = "0x256d1fce1b1221e8398f65f9b36033ce50b2d497"
	tokenAddress6 = "0x85f17cf997934a597031b2e18a9ab6ebd4b9f6a4"
	tokenAddress7 = "0x55c08ca52497e2f1534b59e2917bf524d4765257"
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

	requestedTokensPrices, err := client.GetPricesForRequestedTokens(ctx, spotprices.GetPricesRequestDto{
		Currency: spotprices.GetPricesRequestDtoCurrency(spotprices.USD),
		Tokens:   []string{tokenAddress1, tokenAddress2},
	})
	if err != nil {
		fmt.Println("failed to GetPricesForRequestedTokens: %w", err)
		return
	}

	fmt.Println("GetPricesForRequestedTokens:", requestedTokensPrices)
	time.Sleep(time.Second)

	requestedTokensPricesLarge, err := client.GetPricesForRequestedTokensLarge(ctx, spotprices.GetPricesRequestDto{
		Currency: spotprices.GetPricesRequestDtoCurrency(spotprices.USD),
		Tokens:   []string{tokenAddress1, tokenAddress2, tokenAddress3, tokenAddress4, tokenAddress5, tokenAddress6, tokenAddress7},
	})
	if err != nil {
		fmt.Println("failed to GetPricesForRequestedTokensLarge: %w", err)
		return
	}

	fmt.Println("GetPricesForRequestedTokensLarge:", requestedTokensPricesLarge)
	time.Sleep(time.Second)

	customCurrencies, err := client.GetCustomCurrenciesList(ctx)
	if err != nil {
		fmt.Println("failed to GetCustomCurrenciesList: %w", err)
		return
	}

	fmt.Println("GetCustomCurrenciesList:", customCurrencies)
	time.Sleep(time.Second)
}
