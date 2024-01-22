package client

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
	"github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk/golang/client/onchain"
	"github.com/1inch/1inch-sdk/golang/client/swap"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/amounts"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/contracts"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/typehashes"
)

type ActionService service

// TODO temporarily adding a bool to the function call until config refactor

// SwapTokens executes a token swap operation using the 1inch Swap API.
//
// The helper function takes a client, swap parameters, and a flag to skip warnings. It executes the proposed swap onchain, using Permit if available.
//
// Parameters:
//   - c: A pointer to the client.Client instance. This client should be initialized and connected to the Ethereum network.
//   - swapParams: The parameters for the swap operation, of type swap.AggregationControllerGetSwapParams. It should contain details such as the source and destination tokens, the amount to swap, and the slippage tolerance.
//   - skipWarnings: A boolean flag indicating whether to skip warning prompts. If true, warning prompts will be suppressed; otherwise, they will be displayed.
//
// The function performs several key operations:
//   - Sets a 10-minute Permit1 deadline for the swap operation.
//   - Checks if the source token supports Permit1. If Permit1 is supported, it tries to use that instead of the traditional `Approve` swap.
//   - Executes the swap request onchain
//
// Note:
//   - The function currently has a hardcoded 10-minute deadline. Future versions will make this configurable.
//   - The Permit feature is used if the token typehash matches a known Permit typehash.
//
// Returns nil on successful execution of the swap. Any error during the process is returned as a non-nil error.
func (s *ActionService) SwapTokens(swapParams swap.AggregationControllerGetSwapParams, skipWarnings bool, approvalType swap.ApprovalType) error {

	// Always disable estimate so we can don onchain approvals for the swaps right before we execute
	swapParams.DisableEstimate = helpers.GetPtr(true)

	if s.client.WalletKey == "" {
		return fmt.Errorf("wallet key must be set in the client config")
	}

	deadline := time.Now().Add(1 * time.Minute).Unix() // TODO make this configurable

	executeSwapConfig := &swap.ExecuteSwapConfig{
		FromToken: swapParams.Src,
		ToToken:   swapParams.Dst,
		Amount:    swapParams.Amount,
		Slippage:  swapParams.Slippage,
	}

	if shouldTryPermit(s.client.ChainId, approvalType) {

	}

	var usePermit bool

	// Currently, Permit1 swaps are only tested on Ethereum and Polygon
	isPermitSupportedForCurrentChain := s.client.ChainId == chains.Ethereum || s.client.ChainId == chains.Polygon

	var typehash string
	var err error
	if isPermitSupportedForCurrentChain && approvalType != swap.ApprovalAlways {
		typehash, err = onchain.GetTypeHash(s.client.EthClient, swapParams.Src)
		if err == nil {
			// Typehash is present which means we can use Permit to save gas
			if typehash == typehashes.Permit1 {
				usePermit = true
			} else {
				log.Printf("Typehash exists, but it is not recognized: %v\n", typehash)
			}
		}
	}

	if usePermit || approvalType == swap.PermitAlways {
		name, err := onchain.ReadContractName(s.client.EthClient, common.HexToAddress(swapParams.Src))
		if err != nil {
			return fmt.Errorf("failed to read contract name: %v", err)
		}

		nonce, err := onchain.ReadContractNonce(s.client.EthClient, s.client.PublicAddress, common.HexToAddress(swapParams.Src))
		if err != nil {
			return fmt.Errorf("failed to read contract name: %v", err)
		}

		sig, err := swap.CreatePermitSignature(&swap.PermitSignatureConfig{
			FromToken:     swapParams.Src,
			Name:          name,
			PublicAddress: s.client.PublicAddress.Hex(),
			ChainId:       s.client.ChainId,
			Key:           s.client.WalletKey,
			Nonce:         nonce,
			Deadline:      deadline,
		})
		if err != nil {
			return fmt.Errorf("failed to create permit signature: %v", err)
		}

		aggregationRouter, err := contracts.Get1inchRouterFromChainId(s.client.ChainId)
		if err != nil {
			return fmt.Errorf("failed to get 1inch router address: %v", err)
		}

		permitParams := swap.CreatePermitParams(&swap.PermitParamsConfig{
			Owner:     strings.ToLower(s.client.PublicAddress.Hex()), // TODO remove ToLower and see if it still works
			Spender:   aggregationRouter,
			Value:     amounts.BigMaxUint256,
			Deadline:  deadline,
			Signature: sig,
		})

		executeSwapConfig.IsPermitSwap = true
		swapParams.Permit = &permitParams
		fmt.Println("Swapping using Permit1")
	}

	// Execute swap request
	// This will return the transaction data used by a wallet to execute the swap
	swapResponse, _, err := s.client.Swap.GetSwapData(context.Background(), swapParams, true)
	if err != nil {
		return fmt.Errorf("failed to get swap: %v", err)
	}

	executeSwapConfig.TransactionData = swapResponse.Tx.Data
	executeSwapConfig.EstimatedAmountOut = swapResponse.ToAmount
	executeSwapConfig.SkipWarnings = skipWarnings

	err = s.client.Swap.ExecuteSwap(executeSwapConfig)
	if err != nil {
		return fmt.Errorf("failed to execute swap: %v", err)
	}

	return nil
}

func shouldTryPermit(chainId int, approvalType swap.ApprovalType) bool {
	return approvalType == swap.PermitIfPossible || approvalType == swap.PermitAlways
}
