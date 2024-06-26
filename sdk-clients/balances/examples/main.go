package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/balances"
)

/*
This example demonstrates how to swap tokens on the EthereumChainId network using the 1inch SDK.
The only thing you need to provide is your wallet address, wallet key, and dev portal token.
This can be done through your environment, or you can directly set them in the variables below
*/

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const (
	mainWalletAddress      = "0x1C17622cfa9B6fD2043A76DfC39A5B5a109aa708"
	secondaryWalletAddress = "0x28C6c06298d514Db089934071355E5743bf21d60"

	tokenAddress1 = "0x0d8775f648430679a709e98d2b0cb6250d2887ef"
	tokenAddress2 = "0x58b6a8a3302369daec383334672404ee733ab239"

	spender     = "0x58b6a8a3302369daec383334672404ee733ab239"
	spenderInch = "0x111111125421ca6dc452d289314280a0f8842a65"
)

func main() {
	config, err := balances.NewConfiguration(balances.ConfigurationParams{
		ChainId: constants.EthereumChainId,
		ApiUrl:  "https://api.1inch.dev",
		ApiKey:  devPortalToken,
	})
	if err != nil {
		log.Fatalf("Failed to create configuration: %v\n", err)
	}
	client, err := balances.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}
	ctx := context.Background()

	balancesOfCustomTokensByWalletAddressResponse, err := client.GetBalancesOfCustomTokensByWalletAddress(ctx, balances.BalancesOfCustomTokensByWalletAddressParams{
		Wallet: mainWalletAddress,
		Tokens: []string{tokenAddress1, tokenAddress2},
	})
	if err != nil {
		log.Fatalf("failed to GetBalancesOfCustomTokensByWalletAddress: %v\n", err)
	}

	fmt.Println("GetBalancesOfCustomTokensByWalletAddress:", balancesOfCustomTokensByWalletAddressResponse)
	time.Sleep(time.Second)

	balancesOfCustomTokensByWalletAddressesListResponse, err := client.GetBalancesOfCustomTokensByWalletAddressesList(ctx, balances.BalancesOfCustomTokensByWalletAddressesListParams{
		Wallets: []string{mainWalletAddress, secondaryWalletAddress},
		Tokens:  []string{tokenAddress1, tokenAddress2},
	})
	if err != nil {
		log.Fatalf("failed to GetBalancesOfCustomTokensByWalletAddressesList: %v\n", err)
	}

	fmt.Println("GetBalancesOfCustomTokensByWalletAddressesList:", balancesOfCustomTokensByWalletAddressesListResponse)
	time.Sleep(time.Second)

	aggregatedBalancesAndAllowancesResponse, err := client.GetBalancesAndAllowances(ctx, balances.BalancesAndAllowancesParams{
		Wallets:     []string{mainWalletAddress, secondaryWalletAddress},
		FilterEmpty: true,
		Spender:     spender,
	})
	if err != nil {
		log.Fatalf("failed to GetBalancesAndAllowances: %v\n", err)
	}

	fmt.Println("aggregatedBalancesAndAllowancesResponse:", aggregatedBalancesAndAllowancesResponse)
	time.Sleep(time.Second)

	balancesByWalletAddressResponse, err := client.GetBalancesByWalletAddress(ctx, balances.BalancesByWalletAddressParams{Wallet: mainWalletAddress})
	if err != nil {
		log.Fatalf("failed to GetBalancesByWalletAddress: %v\n", err)
	}

	fmt.Println("GetBalancesByWalletAddress:", balancesByWalletAddressResponse)
	time.Sleep(time.Second)

	allowancesByWalletAddressResponse, err := client.GetAllowancesByWalletAddress(ctx, balances.AllowancesByWalletAddressParams{
		Wallet:  mainWalletAddress,
		Spender: spender,
	})
	if err != nil {
		log.Fatalf("failed to GetAllowancesByWalletAddress: %v\n", err)
	}

	fmt.Println("GetAllowancesByWalletAddress:", allowancesByWalletAddressResponse)
	time.Sleep(time.Second)

	allowancesOfCustomTokensByWalletAddressResponse, err := client.GetAllowancesOfCustomTokensByWalletAddress(ctx, balances.AllowancesOfCustomTokensByWalletAddressParams{
		Wallet:  mainWalletAddress,
		Spender: spender,
		Tokens:  []string{tokenAddress1, tokenAddress2},
	})
	if err != nil {
		log.Fatalf("failed to GetAllowancesOfCustomTokensByWalletAddress: %v\n", err)
	}

	fmt.Println("GetAllowancesOfCustomTokensByWalletAddress:", allowancesOfCustomTokensByWalletAddressResponse)
	time.Sleep(time.Second)

	balancesAndAllowancesByWalletAddressListResponse, err := client.GetBalancesAndAllowancesByWalletAddressList(ctx, balances.BalancesAndAllowancesByWalletAddressListParams{
		Wallet:  secondaryWalletAddress,
		Spender: spenderInch,
	})
	if err != nil {
		log.Fatalf("failed to GetBalancesAndAllowancesByWalletAddressList: %v\n", err)
	}

	fmt.Println("GetBalancesAndAllowancesByWalletAddressList:", balancesAndAllowancesByWalletAddressListResponse)
	time.Sleep(time.Second)

	balancesAndAllowancesOfCustomTokensByWalletAddressResponse, err := client.GetBalancesAndAllowancesOfCustomTokensByWalletAddressList(ctx, balances.BalancesAndAllowancesOfCustomTokensByWalletAddressParams{
		Wallet:  mainWalletAddress,
		Spender: spender,
		Tokens:  []string{tokenAddress1, tokenAddress2},
	})
	if err != nil {
		log.Fatalf("failed to GetBalancesAndAllowancesOfCustomTokensByWalletAddressList: %v\n", err)
	}

	fmt.Println("GetBalancesAndAllowancesOfCustomTokensByWalletAddressList:", balancesAndAllowancesOfCustomTokensByWalletAddressResponse)
}
