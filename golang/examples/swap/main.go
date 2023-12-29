package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"1inch-sdk-golang/client"
	"1inch-sdk-golang/client/swap"
	"1inch-sdk-golang/helpers"
	"1inch-sdk-golang/helpers/consts/amounts"
	"1inch-sdk-golang/helpers/consts/chains"
	"1inch-sdk-golang/helpers/consts/tokens"
)

func main() {

	// Build the config for the client
	config := client.Config{
		DevPortalApiKey:            os.Getenv("DEV_PORTAL_TOKEN"),
		WalletKey:                  os.Getenv("WALLET_KEY_POLY"),
		Web3HttpProviderUrlWithKey: os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY_POLYGON"),
		ChainId:                    chains.Polygon,
	}

	// Create the 1inch client
	c, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Build the config for the swap request
	swapParams := swap.AggregationControllerGetSwapParams{
		Src:             tokens.PolygonDai,
		Dst:             tokens.PolygonWeth,
		From:            os.Getenv("WALLET_ADDRESS_POLY"),
		Amount:          amounts.Ten16,
		DisableEstimate: helpers.GetPtr(true),
	}

	// Execute swap request
	// This will return the transaction data used by a wallet to execute the swap
	swapResponse, _, err := c.Swap.GetSwapData(context.Background(), swapParams)
	if err != nil {
		log.Fatalf("Failed to get swap: %v", err)
	}

	prettyPrintSwapResponse(swapResponse)

	err = c.Swap.ExecuteSwap(tokens.PolygonDai, swapResponse.Tx.Data)
	if err != nil {
		log.Fatalf("Failed to execute swap: %v", err)
	}
}

func prettyPrintSwapResponse(resp *swap.SwapResponse) {
	fmt.Println("Swap Response:")

	if resp.FromToken != nil {
		fmt.Println("FromToken:")
		PrettyPrintTokenInfo(*resp.FromToken)
	}
	if resp.Protocols != nil {
		fmt.Println("Protocols:")
		for _, protoGroup := range *resp.Protocols {
			for _, proto := range protoGroup {
				for _, p := range proto {
					fmt.Printf("\tFromTokenAddress: %s\n", p.FromTokenAddress)
					fmt.Printf("\tName: %s\n", p.Name)
					fmt.Printf("\tPart: %f\n", p.Part)
					fmt.Printf("\tToTokenAddress: %s\n", p.ToTokenAddress)
					fmt.Println()
				}
			}
		}
	}
	fmt.Printf("ToAmount: %s\n", resp.ToAmount)
	if resp.ToToken != nil {
		fmt.Println("ToToken:")
		PrettyPrintTokenInfo(*resp.ToToken)
	}
	fmt.Println("Transaction Data:")
	PrettyPrintTransactionData(resp.Tx)
}

func PrettyPrintTokenInfo(token swap.TokenInfo) {
	fmt.Printf("\tAddress: %s\n", token.Address)
	fmt.Printf("\tDecimals: %f\n", token.Decimals)
	if token.DomainVersion != nil {
		fmt.Printf("\tDomainVersion: %s\n", *token.DomainVersion)
	}
	if token.Eip2612 != nil {
		fmt.Printf("\tEip2612: %v\n", *token.Eip2612)
	}
	if token.IsFoT != nil {
		fmt.Printf("\tIsFoT: %v\n", *token.IsFoT)
	}
	fmt.Printf("\tLogoURI: %s\n", token.LogoURI)
	fmt.Printf("\tName: %s\n", token.Name)
	fmt.Printf("\tSymbol: %s\n", token.Symbol)
	if token.Tags != nil {
		fmt.Printf("\tTags: %v\n", *token.Tags)
	}
}

func PrettyPrintTransactionData(tx swap.TransactionData) {
	fmt.Printf("\tData: %s\n", tx.Data)
	fmt.Printf("\tFrom: %s\n", tx.From)
	fmt.Printf("\tGas: %f\n", tx.Gas)
	fmt.Printf("\tGasPrice: %s\n", tx.GasPrice)
	fmt.Printf("\tTo: %s\n", tx.To)
	fmt.Printf("\tValue: %s\n", tx.Value)
}
