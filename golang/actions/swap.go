package actions

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk/golang/client"
	"github.com/1inch/1inch-sdk/golang/client/onchain"
	"github.com/1inch/1inch-sdk/golang/client/swap"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/amounts"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/contracts"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/typehashes"
)

// TODO temporarily adding a bool to the function call until config refactor

func SwapTokens(c *client.Client, swapParams swap.AggregationControllerGetSwapParams, skipWarnings bool) error {

	deadline := time.Now().Add(10 * time.Minute).Unix() // TODO make this configurable

	executeSwapConfig := &swap.ExecuteSwapConfig{
		FromToken: swapParams.Src,
		ToToken:   swapParams.Dst,
		Amount:    swapParams.Amount,
		Slippage:  swapParams.Slippage,
	}

	typehash, err := swap.GetTypeHash(c.EthClient, swapParams.Src)
	if err == nil {
		// Typehash is present which means we can use Permit to save gas
		if typehash == typehashes.Permit1 {
			name, err := onchain.ReadContractName(c.EthClient, common.HexToAddress(swapParams.Src))
			if err != nil {
				return fmt.Errorf("failed to read contract name: %v", err)
			}

			nonce, err := onchain.ReadContractNonce(c.EthClient, c.PublicAddress, common.HexToAddress(swapParams.Src))
			if err != nil {
				return fmt.Errorf("failed to read contract name: %v", err)
			}

			sig, err := swap.CreatePermitSignature(&swap.PermitSignatureConfig{
				FromToken:     swapParams.Src,
				Name:          name,
				PublicAddress: c.PublicAddress.Hex(),
				ChainId:       c.ChainId,
				Key:           c.WalletKey,
				Nonce:         nonce,
				Deadline:      deadline,
			})
			if err != nil {
				return fmt.Errorf("failed to create permit signature: %v", err)
			}

			permitParams := swap.CreatePermitParams(&swap.PermitParamsConfig{
				Owner:     strings.ToLower(c.PublicAddress.Hex()), // TODO remove ToLower and see if it still works
				Spender:   contracts.AggregationRouterV5,
				Value:     amounts.BigMaxUint256,
				Deadline:  deadline,
				Signature: sig,
			})

			executeSwapConfig.IsPermitSwap = true
			swapParams.Permit = &permitParams
			fmt.Println("Permit supported by this token! Swapping using Permit1")
		} else {
			log.Printf("Typehash exists, but it is not recognized: %v\n", typehash)
		}
	}

	// Execute swap request
	// This will return the transaction data used by a wallet to execute the swap
	swapResponse, _, err := c.Swap.GetSwapData(context.Background(), swapParams, true)
	if err != nil {
		return fmt.Errorf("failed to get swap: %v", err)
	}

	//helpers.PrettyPrintStruct(swapResponse)

	executeSwapConfig.TransactionData = swapResponse.Tx.Data
	executeSwapConfig.EstimatedAmountOut = swapResponse.ToAmount
	executeSwapConfig.SkipWarnings = skipWarnings

	err = c.Swap.ExecuteSwap(executeSwapConfig)
	if err != nil {
		return fmt.Errorf("failed to execute swap: %v", err)
	}

	return nil
}
