package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

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

	config, err := orderbook.NewConfiguration(orderbook.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: privateKey,
		ChainId:    chainId,
		ApiUrl:     "https://api.1inch.dev",
		ApiKey:     devPortalToken,
	})
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

	expireAfter := time.Now().Add(time.Hour).Unix()

	seriesNonce, err := client.GetSeriesNonce(ctx, publicAddress)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get series nonce: %v", err))
	}

	makerTraits := orderbook.NewMakerTraits(orderbook.MakerTraitsParams{
		AllowedSender:      zeroAddress,
		ShouldCheckEpoch:   false,
		UsePermit2:         false,
		UnwrapWeth:         false,
		HasExtension:       false,
		Expiry:             expireAfter,
		Nonce:              seriesNonce.Int64(),
		Series:             0, // TODO: Series 0 always?
		AllowMultipleFills: true,
		AllowPartialFills:  true,
	})

	createOrderResponse, err := client.CreateOrder(ctx, orderbook.CreateOrderParams{
		Wallet:                         client.Wallet,
		SeriesNonce:                    seriesNonce,
		ExpireAfter:                    expireAfter, // TODO update the field name to have "unix" suffix
		Maker:                          publicAddress.Hex(),
		MakerAsset:                     wmatic,
		TakerAsset:                     usdc,
		MakingAmount:                   ten16,
		TakingAmount:                   ten6,
		Taker:                          zeroAddress,
		SkipWarnings:                   false,
		EnableOnchainApprovalsIfNeeded: false,
		MakerTraits:                    makerTraits,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to create order: %v\n", err))
	}
	if !createOrderResponse.Success {
		log.Fatalf("Request completed, but order creation status was a failure: %v\n", createOrderResponse)
	}

	// Sleep to accommodate free-tier API keys
	time.Sleep(time.Second)

	getOrderResponse, err := client.GetOrdersByCreatorAddress(ctx, orderbook.GetOrdersByCreatorAddressParams{
		CreatorAddress: publicAddress.Hex(),
	})

	orderIndented, err := json.MarshalIndent(getOrderResponse[0], "", "  ")
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to marshal response: %v\n", err))
	}

	fmt.Printf("Order created: %s\n", orderIndented)
}
