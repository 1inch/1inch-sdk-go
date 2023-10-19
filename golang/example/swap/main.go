package main

import (
	"fmt"
	"log"
	"os"

	"dev-portal-sdk-go/client"
	"dev-portal-sdk-go/client/swap"
	"dev-portal-sdk-go/helpers"
	"dev-portal-sdk-go/helpers/consts/addresses"
	"dev-portal-sdk-go/helpers/consts/amounts"
	"dev-portal-sdk-go/helpers/consts/tokens"
)

func main() {

	config := client.Config{
		ApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
	}
	c, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	quoteParams := swap.AggregationControllerGetQuoteParams{
		Src:    tokens.EthereumUsdc,
		Dst:    tokens.EthereumWeth,
		Amount: amounts.Ten6,
	}
	quoteResponse, _, err := c.GetQuote(quoteParams)
	if err != nil {
		log.Fatalf("Failed to get quote: %v", err)
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
		log.Fatalf("Failed to get swap: %v", err)
	}

	client.PrettyPrintSwapResponse(swapResponse)
}
