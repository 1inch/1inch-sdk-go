package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"dev-portal-sdk-go/client/swap"
)

func (c Client) GetSwap(params swap.AggregationControllerGetSwapParams) (*swap.SwapResponse, *http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/swap/v5.2/1/swap", c.BaseURL), nil)

	err = params.Validate()
	if err != nil {
		return nil, nil, fmt.Errorf("request validation error: %v", err)
	}

	query := getSwapAddQueryParameters(req.URL.Query(), params)
	req.URL.RawQuery = query.Encode()

	var swap swap.SwapResponse
	res, err := c.Do(context.Background(), req, &swap)
	if err != nil {
		return nil, nil, err
	}

	return &swap, res, nil
}

func getSwapAddQueryParameters(query url.Values, params swap.AggregationControllerGetSwapParams) url.Values {
	query.Add("src", params.Src)
	query.Add("dst", params.Dst)
	query.Add("amount", params.Amount)
	query.Add("from", params.From)
	query.Add("slippage", fmt.Sprintf("%f", params.Slippage))

	if params.Protocols != nil {
		query.Add("protocols", *params.Protocols)
	}
	if params.Fee != nil {
		query.Add("fee", fmt.Sprintf("%f", *params.Fee))
	}
	if params.GasPrice != nil {
		query.Add("gasPrice", *params.GasPrice)
	}
	if params.ComplexityLevel != nil {
		query.Add("complexityLevel", fmt.Sprintf("%f", *params.ComplexityLevel))
	}
	if params.Parts != nil {
		query.Add("parts", fmt.Sprintf("%f", *params.Parts))
	}
	if params.MainRouteParts != nil {
		query.Add("mainRouteParts", fmt.Sprintf("%f", *params.MainRouteParts))
	}
	if params.GasLimit != nil {
		query.Add("gasLimit", fmt.Sprintf("%f", *params.GasLimit))
	}
	if params.IncludeTokensInfo != nil && *params.IncludeTokensInfo {
		query.Add("includeTokensInfo", "true")
	}
	if params.IncludeProtocols != nil && *params.IncludeProtocols {
		query.Add("includeProtocols", "true")
	}
	if params.IncludeGas != nil && *params.IncludeGas {
		query.Add("includeGas", "true")
	}
	if params.ConnectorTokens != nil {
		query.Add("connectorTokens", *params.ConnectorTokens)
	}
	if params.Permit != nil {
		query.Add("permit", *params.Permit)
	}
	if params.Receiver != nil {
		query.Add("receiver", *params.Receiver)
	}
	if params.Referrer != nil {
		query.Add("referrer", *params.Referrer)
	}
	if params.AllowPartialFill != nil && *params.AllowPartialFill {
		query.Add("allowPartialFill", "true")
	}
	if params.DisableEstimate != nil && *params.DisableEstimate {
		query.Add("disableEstimate", "true")
	}
	return query
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
