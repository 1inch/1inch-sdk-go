package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook/models"
)

/*
This example demonstrates how to create an order on the Polygon network using the 1inch SDK.
The only thing you need to provide is your wallet address, wallet key, and dev portal token.
This can be done through your environment, or you can directly set them in the variables below
*/

var (
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const (
	wmatic      = "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"
	usdc        = "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"
	ten16       = "10000000000000000"
	ten6        = "1000000"
	zeroAddress = "0x0000000000000000000000000000000000000000"
	chainId     = 137
)

func main() {
	ctx := context.Background()

	config, err := orderbook.NewDefaultConfiguration(nodeUrl, privateKey, uint64(chainId), "https://api.1inch.dev", devPortalToken)
	if err != nil {
		log.Fatal(err)
	}
	client, err := orderbook.NewClient(config)

	ecdsaPrivateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatalf(fmt.Sprintf("error converting private key to ECDSA: %v", err))
	}
	publicKey := ecdsaPrivateKey.Public()
	publicAddress := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))

	seriesNonce, err := client.GetSeriesNonce(ctx, publicAddress)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get series nonce: %v", err))
	}

	createOrderResponse, err := client.CreateOrder(ctx, models.CreateOrderParams{
		ChainId:                        chainId,
		SeriesNonce:                    seriesNonce,
		PrivateKey:                     privateKey,
		ExpireAfter:                    time.Now().Add(time.Minute * 10).Unix(), // TODO update the field name to have "unix" suffix
		Maker:                          publicAddress.Hex(),
		MakerAsset:                     wmatic,
		TakerAsset:                     usdc,
		MakingAmount:                   ten16,
		TakingAmount:                   ten6,
		Taker:                          zeroAddress,
		SkipWarnings:                   false,
		EnableOnchainApprovalsIfNeeded: false,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to create order: %v\n", err))
	}

	fmt.Printf("Created order: %v\n", createOrderResponse)
}