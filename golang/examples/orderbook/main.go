package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"1inch-sdk-golang/client"
	"1inch-sdk-golang/client/orderbook"
	"1inch-sdk-golang/helpers"
)

func main() {

	// Build the config for the client
	config := client.Config{
		DevPortalApiKey:            os.Getenv("DEV_PORTAL_TOKEN"),
		Web3HttpProviderUrlWithKey: os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY"),
		EtherscanApiKey:            os.Getenv("ETHERSCAN_TOKEN"),
		WalletAddress:              os.Getenv("WALLET_ADDRESS"),
		WalletKey:                  os.Getenv("WALLET_KEY"),
		LimitOrderContract:         "0x1111111254EEB25477B68fb85Ed929f73A960582",
		ChainId:                    1,
	}

	// Create the 1inch client
	c, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Build the config for the orders request
	limitOrdersParams := orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams{
		Page:  helpers.GetPtr(float32(1)),
		Limit: helpers.GetPtr(float32(2)),
	}

	// Execute orders request
	allOrdersResponse, _, err := c.Orderbook.GetAllOrders(context.Background(), limitOrdersParams)
	if err != nil {
		log.Fatalf("Failed to get quote: %v", err)
	}

	prettyPrintOrderResponse(allOrdersResponse)
}

func prettyPrintOrderResponse(orders []*orderbook.OrderResponse) {
	for _, order := range orders {
		jsonOrder, err := json.MarshalIndent(order, "", "  ")
		if err != nil {
			log.Fatalf("Error marshaling to JSON: %v", err)
		}
		fmt.Println(string(jsonOrder))
		fmt.Println("-------------------------------")
	}
}

func PrettyPrint(order *orderbook.Order) {
	fmt.Printf("OrderHash (hex): 0x%s\n", order.OrderHash)
	fmt.Printf("Signature (hex): 0x%s\n", order.Signature)

	// Marshal the struct into JSON
	jsonOrder, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling to JSON: %v", err)
	}
	fmt.Println(string(jsonOrder))
}
