package main

import (
	"context"
	"fmt"
	"os"

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

func main() {
	config, err := balances.NewConfiguration(constants.EthereumChainId, "https://api.1inch.dev", devPortalToken)
	if err != nil {
		return
	}
	client, err := balances.NewClient(config)
	if err != nil {
		return
	}
	ctx := context.Background()

	//b1, err := client.GetBalancesOfCustomTokensByWalletAddress(ctx, balances.BalancesOfCustomTokensByWalletAddressParams{
	//	Wallets: "0x1C17622cfa9B6fD2043A76DfC39A5B5a109aa708",
	//	Tokens:  []string{"0x0d8775f648430679a709e98d2b0cb6250d2887ef", "0x58b6a8a3302369daec383334672404ee733ab239"},
	//})
	//if err != nil {
	//	return
	//}
	//fmt.Println(b1)

	//b2, err := client.GetBalancesOfCustomTokensByWalletAddressesList(ctx, balances.BalancesOfCustomTokensByWalletAddressesListParams{
	//	Wallets: []string{"0x1C17622cfa9B6fD2043A76DfC39A5B5a109aa708", "0x28C6c06298d514Db089934071355E5743bf21d60"},
	//	Tokens:  []string{"0x0d8775f648430679a709e98d2b0cb6250d2887ef", "0x58b6a8a3302369daec383334672404ee733ab239"},
	//})
	//if err != nil {
	//	return
	//}
	//fmt.Println(b2)

	balancesAndAllowances, err := client.GetBalancesAndAllowances(ctx, balances.BalancesAndAllowancesParams{
		Wallets:     []string{"0x1C17622cfa9B6fD2043A76DfC39A5B5a109aa708", "0x28C6c06298d514Db089934071355E5743bf21d60"},
		FilterEmpty: true,
		Spender:     "0x58b6a8a3302369daec383334672404ee733ab239",
	})
	if err != nil {
		return
	}

	fmt.Println(balancesAndAllowances)
	//
	//b, err := client.GetBalancesByWalletAddress(ctx, balances.BalancesByWalletAddressParams{WalletAddress: "0x1C17622cfa9B6fD2043A76DfC39A5B5a109aa708"})
	//if err != nil {
	//	return
	//}
	//fmt.Println(b)
	//

	//b3, err := client.GetAllowancesByWalletAddress(ctx, balances.AllowancesByWalletAddressParams{
	//	Wallet:  "0x1C17622cfa9B6fD2043A76DfC39A5B5a109aa708",
	//	Spender: "0x0d8775f648430679a709e98d2b0cb6250d2887ef",
	//})
	//if err != nil {
	//	return
	//}
	//fmt.Println(b3)

	//b4, err := client.GetAllowancesOfCustomTokensByWalletAddress(ctx, balances.AllowancesOfCustomTokensByWalletAddressParams{
	//	Wallet:  "0x1C17622cfa9B6fD2043A76DfC39A5B5a109aa708",
	//	Spender: "0x0d8775f648430679a709e98d2b0cb6250d2887ef",
	//	Tokens:  []string{"0x0d8775f648430679a709e98d2b0cb6250d2887ef", "0x58b6a8a3302369daec383334672404ee733ab239"},
	//})
	//if err != nil {
	//	return
	//}
	//fmt.Println(b4)

	//b5, err := client.GetBalancesAndAllowancesByWalletAddressList(ctx, balances.BalancesAndAllowancesByWalletAddressListParams{
	//	Wallet:  "0x083fc10cE7e97CaFBaE0fE332a9c4384c5f54E45",
	//	Spender: "0x111111125421ca6dc452d289314280a0f8842a65",
	//})
	//if err != nil {
	//	return
	//}
	//fmt.Println(b5)

	//b6, err := client.GetBalancesAndAllowancesOfCustomTokensByWalletAddressList(ctx, balances.BalancesAndAllowancesOfCustomTokensByWalletAddressParams{
	//	Wallet:  "0x1C17622cfa9B6fD2043A76DfC39A5B5a109aa708",
	//	Spender: "0x0d8775f648430679a709e98d2b0cb6250d2887ef",
	//	Tokens:  []string{"0x0d8775f648430679a709e98d2b0cb6250d2887ef", "0x58b6a8a3302369daec383334672404ee733ab239"},
	//})
	//if err != nil {
	//	return
	//}
	//fmt.Println(b6)
}
