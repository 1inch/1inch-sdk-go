package main

import (
	"fmt"
	"log"
	"os"

	"dev-portal-sdk-go/client"
	"dev-portal-sdk-go/client/swap"
	"dev-portal-sdk-go/helpers"
	"dev-portal-sdk-go/helpers/consts/addresses"
	"dev-portal-sdk-go/helpers/consts/amounts"
	"dev-portal-sdk-go/helpers/consts/tokens"
)

func main() {

	// Build the config for the client
	config := client.Config{
		ApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
	}

	// Create the 1inch client
	c, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Build the config for the quote call
	quoteParams := swap.AggregationControllerGetQuoteParams{
		Src:    tokens.EthereumUsdc,
		Dst:    tokens.EthereumWeth,
		Amount: amounts.Ten6,
	}

	// Execute quote call
	quoteResponse, _, err := c.GetQuote(quoteParams)
	if err != nil {
		log.Fatalf("Failed to get quote: %v", err)
	}

	fmt.Printf("Quote return amount: %v\n", quoteResponse.ToAmount)

	helpers.Sleep()

	// Build the config for the swap call
	swapParams := swap.AggregationControllerGetSwapParams{
		Src:             tokens.EthereumUsdc,
		Dst:             tokens.EthereumWeth,
		From:            addresses.Vitalik,
		Amount:          amounts.Ten6,
		DisableEstimate: helpers.BoolPtr(true),
	}

	// Execute swap call
	// This will return the transaction data used by a wallet to execute the swap
	swapResponse, _, err := c.GetSwap(swapParams)
	if err != nil {
		log.Fatalf("Failed to get swap: %v", err)
	}

	prettyPrintSwapResponse(swapResponse)
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
