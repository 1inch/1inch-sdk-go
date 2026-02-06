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

var (
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const (
	PolygonFRAX = "0x45c32fa6df82ead1e2ef74d17b76547eddfaff89"
	PolygonUsdc = "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"
	ten16   = "10000000000000000"
	ten6    = "1000000"
	ten4    = "10000"
	chainId = 137
)

var (
	makerAsset  = PolygonUsdc
	takerAsset  = PolygonFRAX
	makerAmount = ten4
	takerAmount = ten16
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
		log.Fatalf("Failed to create configuration: %v", err)
	}
	client, err := orderbook.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ecdsaPrivateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatalf("error converting private key to ECDSA: %v", err)
	}
	publicKey := ecdsaPrivateKey.Public()
	publicAddress := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))

	expireAfter := time.Now().Add(time.Hour).Unix()

	router, err := constants.Get1inchRouterFromChainId(chainId)
	if err != nil {
		log.Fatalf("failed to get 1inch router address: %v", err)
	}

	makingAmountInt, err := strconv.ParseInt(makerAmount, 10, 64)
	if err != nil {
		log.Fatalf("Error converting string to int: %v", err)
	}

	permitData, err := client.Wallet.GetContractDetailsForPermit(ctx, common.HexToAddress(makerAsset), common.HexToAddress(router), big.NewInt(makingAmountInt), expireAfter)
	if err != nil {
		log.Fatalf("failed to get permit data: %v", err)
	}
	permit, err := client.Wallet.TokenPermit(*permitData)
	if err != nil {
		log.Fatalf("Failed to get permit: %v", err)
	}

	fmt.Printf("Permit: %v", permit)

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
			Integrator: constants.ZeroAddress,
			Protocol:   constants.ZeroAddress,
			Fee:        0,
			Share:      0,
		},
		MakerPermit: []byte(permit),
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
		log.Fatalf("Failed to create extension: %v", err)
	}

	salt, err := orderbook.GenerateSaltWithFees(&orderbook.GetSaltParams{
		Extension: extensionEncoded,
	})
	if err != nil {
		log.Fatalf("Failed to generate salt: %v", err)
	}

	createOrderResponse, err := client.CreateOrder(ctx, orderbook.CreateOrderParams{
		Wallet:                         client.Wallet,
		Salt:                           fmt.Sprintf("%d", salt),
		MakerTraits:                    orderbook.NewMakerTraitsDefault(), // Defaults to a 1 hour expiration
		ExtensionEncoded:               extensionEncoded,
		Maker:                          publicAddress.Hex(),
		MakerAsset:                     makerAsset,
		TakerAsset:                     takerAsset,
		MakingAmount:                   makerAmount,
		TakingAmount:                   takerAmount,
		Taker:                          feeInfo.ExtensionAddress,
		SkipWarnings:                   false,
		EnableOnchainApprovalsIfNeeded: false,
	})
	if err != nil {
		log.Fatalf("Failed to create order: %v", err)
	}
	if !createOrderResponse.Success {
		log.Fatalf("Request completed, but order creation status was a failure: %v", createOrderResponse)
	}

	// Sleep to accommodate free-tier API keys
	time.Sleep(time.Second)

	getOrderResponse, err := client.GetOrdersByCreatorAddress(ctx, orderbook.GetOrdersByCreatorAddressParams{
		CreatorAddress: publicAddress.Hex(),
	})
	if err != nil {
		log.Fatalf("Failed to get orders by creator address: %v", err)
	}

	orderIndented, err := json.MarshalIndent(getOrderResponse[0], "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal response: %v", err)
	}

	fmt.Printf("Order created: %s\n", orderIndented)
}
