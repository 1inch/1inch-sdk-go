package main

//
//import (
//	"context"
//	"log"
//	"os"
//
//	"github.com/1inch/1inch-sdk-go/internal/helpers/consts/addresses"
//	"github.com/1inch/1inch-sdk-go/internal/helpers/consts/chains"
//	"github.com/1inch/1inch-sdk-go/internal/helpers/consts/tokens"
//
//	"github.com/1inch/1inch-sdk-go/client"
//	"github.com/1inch/1inch-sdk-go/client/models"
//	"github.com/1inch/1inch-sdk-go/helpers"
//)
//
//func main() {
//
//	// Build the config for the client
//	config := models.ClientConfig{
//		DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
//		Web3HttpProviders: []models.Web3Provider{
//			{
//				ChainId: chains.Polygon,
//				Url:     os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY_POLYGON"),
//			},
//		},
//	}
//
//	// Create the 1inch client
//	c, err := client.NewClient(config)
//	if err != nil {
//		log.Fatalf("Failed to create client: %v", err)
//	}
//
//	// Build the config for the orders request
//	createOrderParams := models.CreateOrderParams{
//		ChainId:      chains.Polygon,
//		PrivateKey:   os.Getenv("WALLET_KEY"),
//		Maker:        os.Getenv("WALLET_ADDRESS"),
//		MakerAsset:   tokens.PolygonFrax,
//		TakerAsset:   tokens.PolygonUsdc,
//		MakingAmount: "100000000000000000",
//		TakingAmount: "100000",
//		Taker:        addresses.Zero,
//	}
//
//	// Execute orders request
//	createOrderResponse, _, err := c.OrderbookApi.CreateOrder(context.Background(), createOrderParams)
//	if err != nil {
//		log.Fatalf("Failed to create order: %v", err)
//	}
//
//	helpers.PrettyPrintStruct(createOrderResponse)
//}
