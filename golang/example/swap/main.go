package main

import (
	"fmt"

	"dev-portal-sdk-go/client"
	"dev-portal-sdk-go/client/swap"
)

func main() {
	c := client.NewClient(nil)
	swapParams := swap.AggregationControllerGetQuoteParams{
		Src:    "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		Dst:    "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		Amount: "1000000",
	}
	message, _, err := c.GetQuote(swapParams)
	if err != nil {
		fmt.Printf("Failure: %v\n", err)
		return
	}

	fmt.Printf("Success: %v\n", message)
}
