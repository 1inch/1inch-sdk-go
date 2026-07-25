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
This example runs the full classic swap flow on Base: check the router's USDC
allowance, approve the exact trade amount if it falls short, then swap USDC for
WETH.

Requires the following environment variables:
  - DEV_PORTAL_TOKEN: 1inch Developer Portal API key
  - WALLET_KEY:       private key (64 hex chars, no 0x prefix)
  - NODE_URL:         RPC endpoint for Base
*/

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
)

const (
	UsdcBase   = "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"
	WethBase   = "0x4200000000000000000000000000000000000006"
	amountUsdc = "100000" // 0.1 USDC (6 decimals)
)

func main() {
	if devPortalToken == "" || privateKey == "" || nodeUrl == "" {
		log.Fatal("set DEV_PORTAL_TOKEN, WALLET_KEY, and NODE_URL to run this example")
	}

	config, err := aggregation.NewConfiguration(aggregation.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: privateKey,
		ChainId:    constants.BaseChainId,
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

	walletAddr := client.Wallet.Address().Hex()
	amountToSwap, ok := new(big.Int).SetString(amountUsdc, 10)
	if !ok {
		log.Fatalf("invalid amount: %s", amountUsdc)
	}

	// Step 1: check the router's current allowance
	allowanceData, err := client.GetApproveAllowance(ctx, aggregation.GetAllowanceParams{
		TokenAddress:  UsdcBase,
		WalletAddress: walletAddr,
	})
	if err != nil {
		log.Fatalf("failed to get allowance: %v", err)
	}
	allowance := new(big.Int)
	if _, ok := allowance.SetString(allowanceData.Allowance, 10); !ok {
		log.Fatalf("failed to parse allowance: %s", allowanceData.Allowance)
	}

	// Step 2: approve if needed
	if allowance.Cmp(amountToSwap) < 0 {
		fmt.Println("Insufficient allowance; approving...")
		approveData, err := client.GetApproveTransaction(ctx, aggregation.GetApproveParams{
			TokenAddress: UsdcBase,
			Amount:       amountUsdc,
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
		fmt.Printf("Approve transaction sent: https://basescan.org/tx/%s\n", signedTx.Hash().Hex())
		waitForReceipt(ctx, client, signedTx.Hash())
		fmt.Println("Approve transaction confirmed")
	} else {
		fmt.Println("Sufficient allowance already present")
	}

	// Step 3: swap
	swapData, err := client.GetSwap(ctx, aggregation.GetSwapParams{
		Src:      UsdcBase,
		Dst:      WethBase,
		Amount:   amountUsdc,
		From:     walletAddr,
		Slippage: 1, // 1% slippage
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
	fmt.Printf("Swap sent: https://basescan.org/tx/%s\n", signedTx.Hash().Hex())
	waitForReceipt(ctx, client, signedTx.Hash())
	fmt.Println("Swap confirmed")
}

// waitForReceipt polls for a transaction receipt until it lands or a deadline passes
func waitForReceipt(ctx context.Context, client *aggregation.Client, hash common.Hash) {
	deadline := time.Now().Add(3 * time.Minute)
	for time.Now().Before(deadline) {
		receipt, err := client.Wallet.TransactionReceipt(ctx, hash)
		if err == nil {
			if receipt.Status != types.ReceiptStatusSuccessful {
				log.Fatalf("transaction reverted: %s", hash.Hex())
			}
			return
		}
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("timed out waiting for receipt: %s", hash.Hex())
}
