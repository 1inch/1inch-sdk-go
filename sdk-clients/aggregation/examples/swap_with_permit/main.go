package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/1inch/1inch-sdk-go/v4/constants"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/aggregation"
)

/*
This example swaps FRAX for WETH on Polygon using a signed EIP-2612 permit
instead of a prior ERC20 approval, so no separate approve transaction is needed.
The permit grants the router exactly the trade amount and rides inside the swap
transaction. It is only signed when the current allowance cannot cover the trade.

Note that a 2612 permit overwrites any standing allowance the router holds for
the token.

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
	PolygonFrax = "0x45c32fa6df82ead1e2ef74d17b76547eddfaff89"
	PolygonWeth = "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619"
	amountFrax  = "10000000000000000" // 0.01 FRAX
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

	amountToSwap, ok := new(big.Int).SetString(amountFrax, 10)
	if !ok {
		log.Fatalf("invalid amount: %s", amountFrax)
	}

	allowanceData, err := client.GetApproveAllowance(ctx, aggregation.GetAllowanceParams{
		TokenAddress:  PolygonFrax,
		WalletAddress: client.Wallet.Address().Hex(),
	})
	if err != nil {
		log.Fatalf("failed to get allowance: %v", err)
	}
	allowance := new(big.Int)
	if _, ok := allowance.SetString(allowanceData.Allowance, 10); !ok {
		log.Fatalf("failed to parse allowance: %s", allowanceData.Allowance)
	}

	var permit string
	if amountToSwap.Cmp(allowance) > 0 {
		spender, err := client.GetApproveSpender(ctx)
		if err != nil {
			log.Fatalf("failed to get approve spender: %v", err)
		}

		permitDeadline := time.Now().Add(30 * time.Minute).Unix()
		permitData, err := client.Wallet.GetContractDetailsForPermit(ctx, common.HexToAddress(PolygonFrax), common.HexToAddress(spender.Address), amountToSwap, permitDeadline)
		if err != nil {
			log.Fatalf("failed to get permit details: %v", err)
		}
		permit, err = client.Wallet.TokenPermit(*permitData)
		if err != nil {
			log.Fatalf("failed to sign permit: %v", err)
		}
		fmt.Println("Permit signed")
	} else {
		fmt.Println("The router already has a sufficient allowance; swapping without a permit")
	}

	swapData, err := client.GetSwap(ctx, aggregation.GetSwapParams{
		Src:      PolygonFrax,
		Dst:      PolygonWeth,
		Amount:   amountToSwap.String(),
		From:     client.Wallet.Address().Hex(),
		Slippage: 1,
		Permit:   permit,
	})
	if err != nil {
		log.Fatalf("failed to get swap data: %v", err)
	}

	tx, err := client.TxBuilder.New().
		SetData(swapData.TxNormalized.Data).
		SetTo(&swapData.TxNormalized.To).
		SetGas(swapData.TxNormalized.Gas).
		SetValue(swapData.TxNormalized.Value).
		Build(ctx)
	if err != nil {
		log.Fatalf("failed to build swap transaction: %v", err)
	}
	signedTx, err := client.Wallet.Sign(tx)
	if err != nil {
		log.Fatalf("failed to sign swap transaction: %v", err)
	}
	if err := client.Wallet.BroadcastTransaction(ctx, signedTx); err != nil {
		log.Fatalf("failed to broadcast swap transaction: %v", err)
	}
	fmt.Printf("Swap sent: https://polygonscan.com/tx/%s\n", signedTx.Hash().Hex())

	deadline := time.Now().Add(3 * time.Minute)
	for time.Now().Before(deadline) {
		receipt, err := client.Wallet.TransactionReceipt(ctx, signedTx.Hash())
		if err == nil {
			if receipt.Status != types.ReceiptStatusSuccessful {
				log.Fatalf("swap transaction reverted: %s", signedTx.Hash().Hex())
			}
			fmt.Println("Swap confirmed")
			return
		}
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("timed out waiting for receipt: %s", signedTx.Hash().Hex())
}
