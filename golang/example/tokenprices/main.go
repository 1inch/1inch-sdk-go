package main

import (
	"fmt"

	"dev-portal-sdk-go/client"
	"dev-portal-sdk-go/client/spotprice"
)

func main() {
	c := client.NewClient()
	priceParameters := spotprice.PricesParameters{
		Currency: spotprice.CurrencyTypeUSD,
	}
	message, _, err := c.GetTokenPrices(priceParameters)
	if err != nil {
		fmt.Printf("Failure: %v\n", err)
		return
	}

	fmt.Printf("Success: %v\n", message)
}
