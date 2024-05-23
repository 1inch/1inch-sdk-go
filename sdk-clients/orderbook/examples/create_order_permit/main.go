package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
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
	PolygonFRAX = "0x45c32fa6df82ead1e2ef74d17b76547eddfaff89"
	PolygonUsdc = "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"
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

	router, err := constants.Get1inchRouterFromChainId(chainId)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get 1inch router address: %v", err))
	}

	makingAmount := ten16
	makingAmountInt, err := strconv.ParseInt(makingAmount, 10, 64)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return
	}

	permitData, err := client.Wallet.GetContractDetailsForPermit(ctx, common.HexToAddress(PolygonFRAX), common.HexToAddress(router), big.NewInt(makingAmountInt), expireAfter)
	if err != nil {
		panic(err)
	}
	permit, err := client.Wallet.TokenPermit(*permitData)
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to get permit: %v\n", err))
	}

	extension, err := orderbook.NewExtension(orderbook.ExtensionParams{
		MakerAsset: PolygonFRAX,
		Permit:     permit,
	})
	if err != nil {
		log.Fatalf("Failed to create extension: %v\n", err)
	}

	makerTraits := orderbook.NewMakerTraits(orderbook.MakerTraitsParams{
		AllowedSender:      zeroAddress,
		ShouldCheckEpoch:   false,
		UsePermit2:         false,
		UnwrapWeth:         false,
		HasExtension:       true,
		HasPreInteraction:  false,
		HasPostInteraction: false,
		Expiry:             expireAfter,
		Nonce:              seriesNonce.Int64(),
		Series:             0, // TODO: Series 0 always?
	})

	createOrderResponse, err := client.CreateOrder(ctx, orderbook.CreateOrderParams{
		SeriesNonce:                    seriesNonce,
		MakerTraits:                    makerTraits,
		Extension:                      extension,
		PrivateKey:                     privateKey,
		ExpireAfter:                    expireAfter, // TODO update the field name to have "unix" suffix
		Maker:                          publicAddress.Hex(),
		MakerAsset:                     PolygonFRAX,
		TakerAsset:                     PolygonUsdc,
		MakingAmount:                   makingAmount,
		TakingAmount:                   ten6,
		Taker:                          zeroAddress,
		SkipWarnings:                   false,
		EnableOnchainApprovalsIfNeeded: false,
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
