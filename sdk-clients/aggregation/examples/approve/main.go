package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/aggregation"
	"github.com/ethereum/go-ethereum/common"
)

var (
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const (
	PolygonDai  = "0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063"
	PolygonWeth = "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619"
)

func main() {
	config, err := aggregation.NewConfiguration(aggregation.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: privateKey,
		ChainId:    constants.PolygonChainId,
		ApiUrl:     "https://api.1inch.dev",
		ApiKey:     devPortalToken,
	})
	if err != nil {
		log.Fatalf("Failed to create configuration: %v\n", err)
	}
	client, err := aggregation.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}
	ctx := context.Background()

	amountToSwap := big.NewInt(1e18)

	allowanceData, err := client.GetApproveAllowance(ctx, aggregation.GetAllowanceParams{
		TokenAddress:  PolygonWeth,
		WalletAddress: client.Wallet.Address().Hex(),
	})

	allowance := new(big.Int)
	allowance.SetString(allowanceData.Allowance, 10)

	cmp := amountToSwap.Cmp(allowance)

	if cmp > 0 {
		approveData, err := client.GetApproveTransaction(ctx, aggregation.GetApproveParams{
			TokenAddress: PolygonWeth,
			Amount:       amountToSwap.String(),
		})
		if err != nil {
			log.Fatalf("Failed to get approve data: %v\n", err)
		}
		data, err := hex.DecodeString(approveData.Data[2:])
		if err != nil {
			log.Fatalf("Failed to decode approve data: %v\n", err)
		}

		to := common.HexToAddress(approveData.To)

		tx, err := client.TxBuilder.New().SetData(data).SetTo(&to).Build(ctx)
		if err != nil {
			log.Fatalf("Failed to build approve transaction: %v\n", err)
		}

		signedTx, err := client.Wallet.Sign(tx)
		if err != nil {
			log.Fatalf("Failed to sign approve transaction: %v\n", err)
		}

		err = client.Wallet.BroadcastTransaction(ctx, signedTx)
		if err != nil {
			log.Fatalf("Failed to broadcast approve transaction: %v\n", err)
		}

		fmt.Printf("Transaction has been broadcast. View it on Polygonscan here: %v\n", fmt.Sprintf("https://polygonscan.com/tx/%v", signedTx.Hash().Hex()))
		for {
			receipt, err := client.Wallet.TransactionReceipt(ctx, signedTx.Hash())
			if receipt != nil {
				fmt.Println("Transaction complete!")
				return
			}
			if err != nil {
				fmt.Println("Waiting for transaction to be mined")
			}
			select {
			case <-time.After(1000 * time.Millisecond): // check again after a delay
			case <-ctx.Done():
				fmt.Println("Context cancelled")
				return
			}
		}
	}
}
