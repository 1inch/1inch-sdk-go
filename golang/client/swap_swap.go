package client

import (
	"context"
	"fmt"
	"net/http"

	"dev-portal-sdk-go/client/swap"
)

func (c Client) GetSwap(params swap.AggregationControllerGetSwapParams) (*swap.SwapResponse, *http.Response, error) {
	u := "/swap/v5.2/1/swap"

	err := params.Validate()
	if err != nil {
		return nil, nil, fmt.Errorf("request validation error: %v", err)
	}

	u, err = addOptions(u, params)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var swap swap.SwapResponse
	res, err := c.Do(context.Background(), req, &swap)
	if err != nil {
		return nil, nil, err
	}

	return &swap, res, nil
}

func PrettyPrintSwapResponse(resp *swap.SwapResponse) {
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
