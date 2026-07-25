package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/1inch/1inch-sdk-go/v4/constants"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/aggregation"
)

/*
This example grants the 1inch Aggregation Router an ERC20 allowance for WETH on
Polygon, using the approve transaction returned by the API. It skips the
transaction when the current allowance already covers the amount.

Requires the following environment variables:
  - DEV_PORTAL_TOKEN: 1inch Developer Portal API key
  - WALLET_KEY:       private key (64 hex chars, no 0x prefix)
  - NODE_URL:         RPC endpoint for Polygon
*/

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
)

const (
	PolygonWeth   = "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619"
	amountToAllow = "1000000000000000000" // 1 WETH
)

func main() {
	if devPortalToken == "" || privateKey == "" || nodeUrl == "" {
		log.Fatal("set DEV_PORTAL_TOKEN, WALLET_KEY, and NODE_URL to run this example")
	}

	config, err := aggregation.NewConfiguration(aggregation.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: privateKey,
		ChainId:    constants.PolygonChainId,
		ApiUrl:     "https://api.1inch.com",
		ApiKey:     devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}
	client, err := aggregation.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()

	amount, ok := new(big.Int).SetString(amountToAllow, 10)
	if !ok {
		log.Fatalf("invalid amount: %s", amountToAllow)
	}

	allowanceData, err := client.GetApproveAllowance(ctx, aggregation.GetAllowanceParams{
		TokenAddress:  PolygonWeth,
		WalletAddress: client.Wallet.Address().Hex(),
	})
	if err != nil {
		log.Fatalf("failed to get allowance: %v", err)
	}
	allowance := new(big.Int)
	if _, ok := allowance.SetString(allowanceData.Allowance, 10); !ok {
		log.Fatalf("failed to parse allowance: %s", allowanceData.Allowance)
	}

	if allowance.Cmp(amount) >= 0 {
		fmt.Printf("The router already has a sufficient allowance (%s)\n", allowance)
		return
	}

	approveData, err := client.GetApproveTransaction(ctx, aggregation.GetApproveParams{
		TokenAddress: PolygonWeth,
		Amount:       amount.String(),
	})
	if err != nil {
		log.Fatalf("failed to get approve transaction: %v", err)
	}
	data, err := hexutil.Decode(approveData.Data)
	if err != nil {
		log.Fatalf("failed to decode approve data: %v", err)
	}
	to := common.HexToAddress(approveData.To)

	tx, err := client.TxBuilder.New().SetData(data).SetTo(&to).Build(ctx)
	if err != nil {
		log.Fatalf("failed to build approve transaction: %v", err)
	}
	signedTx, err := client.Wallet.Sign(tx)
	if err != nil {
		log.Fatalf("failed to sign approve transaction: %v", err)
	}
	if err := client.Wallet.BroadcastTransaction(ctx, signedTx); err != nil {
		log.Fatalf("failed to broadcast approve transaction: %v", err)
	}
	fmt.Printf("Approve transaction sent: https://polygonscan.com/tx/%s\n", signedTx.Hash().Hex())

	deadline := time.Now().Add(3 * time.Minute)
	for time.Now().Before(deadline) {
		receipt, err := client.Wallet.TransactionReceipt(ctx, signedTx.Hash())
		if err == nil {
			if receipt.Status != types.ReceiptStatusSuccessful {
				log.Fatalf("approve transaction reverted: %s", signedTx.Hash().Hex())
			}
			fmt.Println("Approve transaction confirmed")
			return
		}
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("timed out waiting for receipt: %s", signedTx.Hash().Hex())
}
