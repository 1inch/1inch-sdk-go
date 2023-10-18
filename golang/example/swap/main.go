package main

import (
	"fmt"

	"dev-portal-sdk-go/client"
	"dev-portal-sdk-go/client/swap"
	"dev-portal-sdk-go/helpers"
	"dev-portal-sdk-go/helpers/consts/addresses"
	"dev-portal-sdk-go/helpers/consts/amounts"
	"dev-portal-sdk-go/helpers/consts/tokens"
)

func main() {
	c := client.NewClient(nil)
	quoteParams := swap.AggregationControllerGetQuoteParams{
		Src:    tokens.EthereumUsdc,
		Dst:    tokens.EthereumWeth,
		Amount: amounts.Ten6,
	}
	quoteResponse, _, err := c.GetQuote(quoteParams)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Quote: %v\n", quoteResponse.ToAmount)

	helpers.Sleep()

	swapParams := swap.AggregationControllerGetSwapParams{
		Src:             tokens.EthereumWeth,
		Dst:             tokens.EthereumUsdc,
		From:            addresses.Vitalik,
		Amount:          amounts.Ten18,
		DisableEstimate: helpers.BoolPtr(true),
	}
	swapResponse, _, err := c.GetSwap(swapParams)
	if err != nil {
		fmt.Println(err)
		return
	}

	client.PrettyPrintSwapResponse(swapResponse)
}
