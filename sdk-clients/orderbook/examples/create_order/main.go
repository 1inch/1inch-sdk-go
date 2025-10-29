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
	ten18       = "1000000000000000000"
	ten8        = "100000000"
	zeroAddress = "0x0000000000000000000000000000000000000000"
	chainId     = 137
)

var (
	makerAsset  = wmatic
	takerAsset  = usdc
	makerAmount = ten18
	takerAmount = ten8
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

	feeInfo, err := client.GetFeeInfo(ctx, orderbook.GetFeeInfoParams{
		MakerAsset:  makerAsset,
		TakerAsset:  takerAsset,
		MakerAmount: makerAmount,
		TakerAmount: takerAmount,
	})
	if err != nil {
		log.Fatalf("Failed to get fee info: %v", err)
	}

	buildOrderExtensionBytesParams := &orderbook.BuildOrderExtensionBytesParams{
		ExtensionTarget: feeInfo.ExtensionAddress,
		IntegratorFee: &orderbook.IntegratorFee{
			Integrator: zeroAddress,
			Protocol:   zeroAddress,
			Fee:        0,
			Share:      0,
		},
		ResolverFee: &orderbook.ResolverFee{
			Receiver:          feeInfo.ProtocolFeeReceiver,
			Fee:               feeInfo.FeeBps,
			WhitelistDiscount: feeInfo.WhitelistDiscountPercent,
		},
		Whitelist:      feeInfo.Whitelist,
		CustomReceiver: publicAddress.Hex(),
	}

	extensionEncoded, err := orderbook.BuildOrderExtensionBytes(buildOrderExtensionBytesParams)
	if err != nil {
		log.Fatalf("Failed to create extension: %v\n", err)
	}

	salt, err := orderbook.GenerateSaltNew(&orderbook.GetSaltParams{
		Extension: extensionEncoded,
	})
	if err != nil {
		log.Fatalf("Failed to generate salt: %v", err)
	}

	createOrderResponse, err := client.CreateOrder(ctx, orderbook.CreateOrderParams{
		Wallet:                         client.Wallet,
		Salt:                           fmt.Sprintf("%d", salt),
		Maker:                          publicAddress.Hex(),
		MakerAsset:                     makerAsset,
		TakerAsset:                     takerAsset,
		MakingAmount:                   makerAmount,
		TakingAmount:                   takerAmount,
		Taker:                          feeInfo.ExtensionAddress,
		SkipWarnings:                   false,
		EnableOnchainApprovalsIfNeeded: false,
		MakerTraits:                    orderbook.NewMakerTraitsDefault(),
		MakerTraitsEncoded:             orderbook.NewMakerTraitsDefault().Encode(),
		ExtensionEncoded:               extensionEncoded,
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
