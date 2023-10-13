package main

import (
	"fmt"

	"dev-portal-sdk-go/client"
	"dev-portal-sdk-go/client/tokenprices"
)

func main() {
	c := client.NewClient()
	priceParameters := tokenprices.PricesParameters{
		Currency: tokenprices.CurrencyTypeUSD,
	}
	message, _, err := c.GetTokenPrices(priceParameters)
	if err != nil {
		fmt.Printf("Failure: %v\n", err)
		return
	}

	fmt.Printf("Success: %v\n", message)
}
