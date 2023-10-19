package main

import (
	"fmt"
	"log"
	"os"

	"dev-portal-sdk-go/client"
	"dev-portal-sdk-go/client/spotprice"
)

func main() {

	config := client.Config{
		ApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
	}
	c, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	priceParameters := spotprice.ChainControllerByAddressesParams{
		Currency: spotprice.GetCurrencyType(spotprice.USD),
	}
	message, _, err := c.GetTokenPrices(priceParameters)
	if err != nil {
		log.Fatalf("Failed to get token prices: %v", err)
	}

	fmt.Println(message)
}
