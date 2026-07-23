package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethCommon "github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk-go/v4/constants"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/fusion"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/orderbook"
)

/*
This example places a Fusion order whose maker funds are pulled through Permit2
instead of a direct ERC20 approval to the 1inch router.

Permit2 uses two layers of approval:

 1. A one-time on-chain ERC20 approval from the sell token to the canonical
    Permit2 contract (constants.Permit2Address). This is the only transaction
    the maker ever sends.
 2. A signed (gasless) PermitSingle message granting the 1inch Aggregation
    Router an allowance inside Permit2. This signature is embedded in the
    fusion order and executed on-chain by the protocol during the fill.

Requires the following environment variables:
  - WALLET_KEY:       private key (64 hex chars, no 0x prefix)
  - WALLET_ADDRESS:   address of the wallet
  - NODE_URL:         RPC endpoint for the chain (used for the one-time approval and nonce read)
  - DEV_PORTAL_TOKEN: 1inch Developer Portal API key
*/

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
	publicAddress  = os.Getenv("WALLET_ADDRESS")
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
)

const (
	usdc    = "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"
	weth    = "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619"
	amount  = "200000000000000"
	chainId = 137
)

func main() {
	ctx := context.Background()

	// The orderbook client is RPC-connected and handles the on-chain Permit2 steps
	orderbookConfig, err := orderbook.NewConfiguration(orderbook.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: privateKey,
		ChainId:    chainId,
		ApiUrl:     "https://api.1inch.com",
		ApiKey:     devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create orderbook configuration: %v", err)
	}
	orderbookClient, err := orderbook.NewClient(orderbookConfig)
	if err != nil {
		log.Fatalf("failed to create orderbook client: %v", err)
	}

	fusionConfig, err := fusion.NewConfiguration(fusion.ConfigurationParams{
		ApiUrl:     "https://api.1inch.com",
		ApiKey:     devPortalToken,
		ChainId:    chainId,
		PrivateKey: privateKey,
	})
	if err != nil {
		log.Fatalf("failed to create fusion configuration: %v", err)
	}
	fusionClient, err := fusion.NewClient(fusionConfig)
	if err != nil {
		log.Fatalf("failed to create fusion client: %v", err)
	}

	fromToken := weth
	toToken := usdc
	// The permit owner, order maker, and nonce lookups must all use the address
	// derived from the signing key; WALLET_ADDRESS is only sanity-checked against it
	owner := orderbookClient.Wallet.Address()
	if publicAddress != "" && !strings.EqualFold(publicAddress, owner.Hex()) {
		log.Fatalf("WALLET_ADDRESS %s does not match the address %s derived from WALLET_KEY", publicAddress, owner.Hex())
	}
	sellToken := gethCommon.HexToAddress(fromToken)
	router := gethCommon.HexToAddress(constants.AggregationRouterV6)
	makingAmount, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		log.Fatalf("invalid amount: %s", amount)
	}

	// Step 1: one-time ERC20 approval of the sell token to the Permit2 contract
	if err := ensurePermit2Approval(ctx, orderbookClient, sellToken, makingAmount); err != nil {
		log.Fatalf("failed to ensure Permit2 approval: %v", err)
	}

	// Step 2: read the current Permit2 nonce for (owner, token, router) and sign
	// a PermitSingle granting the router exactly the order amount
	allowance, err := orderbook.GetPermit2Allowance(ctx, orderbookClient.Wallet, owner, sellToken, router)
	if err != nil {
		log.Fatalf("failed to read Permit2 allowance: %v", err)
	}

	expiration := big.NewInt(time.Now().Add(30 * 24 * time.Hour).Unix())
	sigDeadline := big.NewInt(time.Now().Add(30 * time.Minute).Unix())
	permit, err := orderbook.BuildPermit2Calldata(orderbookClient.Wallet, orderbook.Permit2PermitParams{
		Token:       sellToken,
		Amount:      makingAmount,
		Expiration:  expiration,
		Nonce:       allowance.Nonce,
		Spender:     router,
		SigDeadline: sigDeadline,
	})
	if err != nil {
		log.Fatalf("failed to build Permit2 calldata: %v", err)
	}

	// Step 3: quote and place the order in one call. PlaceOrderFromParams propagates
	// Permit and IsPermit2 into both the quote request and the order.
	orderParams := fusion.OrderParams{
		WalletAddress:    owner.Hex(),
		FromTokenAddress: fromToken,
		ToTokenAddress:   toToken,
		Amount:           amount,
		Receiver:         constants.ZeroAddress,
		Preset:           fusion.Fast,
		Permit:           permit,
		IsPermit2:        true,
	}

	orderHash, err := fusionClient.PlaceOrderFromParams(ctx, orderParams)
	if err != nil {
		log.Fatalf("failed to place order: %v", err)
	}

	fmt.Printf("Order placed! Order hash: %s\n", orderHash)
	fmt.Println("Monitoring order until it completes...")

	for {
		<-time.After(1 * time.Second)
		order, err := fusionClient.GetOrderStatus(ctx, orderHash)
		if err != nil {
			fmt.Printf("failed to get order from order hash: %v", err)
			return
		}

		fmt.Printf("Order status: %s\n", order.Status)
		switch order.Status {
		case "filled":
			return
		case "expired", "cancelled", "refunded", "false-predicate", "not-enough-balance-or-allowance", "wrong-permit":
			fmt.Printf("Order ended without filling (status %s)\n", order.Status)
			return
		}
	}
}

// ensurePermit2Approval checks the ERC20 allowance from the sell token to the Permit2
// contract and sends an unlimited approval if it cannot cover the order amount
func ensurePermit2Approval(ctx context.Context, client *orderbook.Client, token gethCommon.Address, required *big.Int) error {
	erc20, err := abi.JSON(strings.NewReader(constants.Erc20ABI))
	if err != nil {
		return err
	}
	permit2 := gethCommon.HexToAddress(constants.Permit2Address)

	allowanceData, err := erc20.Pack("allowance", client.Wallet.Address(), permit2)
	if err != nil {
		return err
	}
	result, err := client.Wallet.Call(ctx, token, allowanceData)
	if err != nil {
		return fmt.Errorf("failed to read ERC20 allowance: %w", err)
	}
	if new(big.Int).SetBytes(result).Cmp(required) >= 0 {
		fmt.Println("Permit2 already has a sufficient ERC20 approval")
		return nil
	}

	// The unlimited approval is the common Permit2 pattern: per-trade limits are
	// enforced by the signed permits, which are amount-scoped and expiring. To keep
	// the ERC20 layer bounded too, replace constants.Uint256Max with the exact
	// required amount at the cost of one approval transaction per order.
	fmt.Println("Sending one-time ERC20 approval to Permit2...")
	approveData, err := erc20.Pack("approve", permit2, constants.Uint256Max)
	if err != nil {
		return err
	}
	tx, err := client.TxBuilder.New().SetData(approveData).SetTo(&token).Build(ctx)
	if err != nil {
		return fmt.Errorf("failed to build approval tx: %w", err)
	}
	signedTx, err := client.Wallet.Sign(tx)
	if err != nil {
		return fmt.Errorf("failed to sign approval tx: %w", err)
	}
	if err := client.Wallet.BroadcastTransaction(ctx, signedTx); err != nil {
		return fmt.Errorf("failed to broadcast approval tx: %w", err)
	}

	deadline := time.Now().Add(3 * time.Minute)
	for time.Now().Before(deadline) {
		receipt, err := client.Wallet.TransactionReceipt(ctx, signedTx.Hash())
		if err == nil {
			if receipt.Status != 1 {
				return fmt.Errorf("approval tx reverted: %s", signedTx.Hash().Hex())
			}
			fmt.Printf("Approval confirmed: %s\n", signedTx.Hash().Hex())
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("timed out waiting for approval receipt: %s", signedTx.Hash().Hex())
}
