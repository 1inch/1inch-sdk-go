package main

import (
	"context"
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
		DevPortalApiKey:     os.Getenv("DEV_PORTAL_TOKEN"),
		Web3HttpProviderUrl: os.Getenv("WEB_3_HTTP_PROVIDER_URL"),
		EtherscanApiKey:     os.Getenv("ETHERSCAN_TOKEN"),
		WalletAddress:       os.Getenv("WALLET_ADDRESS"),
		WalletKey:           os.Getenv("WALLET_KEY"),
		LimitOrderContract:  "0x1111111254EEB25477B68fb85Ed929f73A960582",
		ChainId:             "1",
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

	prettyPrint(allOrdersResponse)
}

func prettyPrint(orders []*orderbook.OrderResponse) {
	for _, order := range orders {
		fmt.Println("Signature:", order.Signature)
		fmt.Println("OrderHash:", order.OrderHash)
		fmt.Println("CreateDateTime:", order.CreateDateTime)
		fmt.Println("RemainingMakerAmount:", order.RemainingMakerAmount)
		fmt.Println("MakerBalance:", order.MakerBalance)
		fmt.Println("MakerAllowance:", order.MakerAllowance)
		fmt.Println("Data:")
		fmt.Println("\tMakerAsset:", order.Data.MakerAsset)
		fmt.Println("\tTakerAsset:", order.Data.TakerAsset)
		fmt.Println("\tSalt:", order.Data.Salt)
		fmt.Println("\tReceiver:", order.Data.Receiver)
		fmt.Println("\tAllowedSender:", order.Data.AllowedSender)
		fmt.Println("\tMakingAmount:", order.Data.MakingAmount)
		fmt.Println("\tTakingAmount:", order.Data.TakingAmount)
		fmt.Println("\tMaker:", order.Data.Maker)
		fmt.Println("\tInteractions:", order.Data.Interactions)
		fmt.Println("\tOffsets:", order.Data.Offsets)
		fmt.Println("MakerRate:", order.MakerRate)
		fmt.Println("TakerRate:", order.TakerRate)
		fmt.Println("IsMakerContract:", order.IsMakerContract)
		fmt.Println("OrderInvalidReason:", order.OrderInvalidReason)
		fmt.Println("-------------------------------")
	}
}
