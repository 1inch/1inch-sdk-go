package client

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/1inch/1inch-sdk/golang/client/tenderly"
	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/tokens"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/1inch/1inch-sdk/golang/client/onchain"
	"github.com/1inch/1inch-sdk/golang/client/swap"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/amounts"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/contracts"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/typehashes"
)

// This file provides helper functions that execute swaps onchain.

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
func (s *ActionService) SwapTokens(params swap.SwapTokensParams) error {

	// Always disable estimate so we can don onchain approvals for the swaps right before we execute
	params.DisableEstimate = helpers.GetPtr(true)

	// TODO find a better way of managing the matching between public and private keys
	privateKey, err := crypto.HexToECDSA(params.WalletKey)
	if err != nil {
		return fmt.Errorf("failed to convert private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("could not cast public key to ECDSA")
	}

	derivedPublicAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	if strings.ToLower(derivedPublicAddress.Hex()) != strings.ToLower(params.PublicAddress) {
		return fmt.Errorf("public address does not match with private key")
	}

	if params.WalletKey == "" {
		return fmt.Errorf("wallet key must be provided")
	}

	ethClient, err := s.client.GetEthClient(params.ChainId)
	if err != nil {
		return fmt.Errorf("failed to get eth client: %v", err)
	}

	deadline := time.Now().Add(1 * time.Minute).Unix() // TODO make this configurable

	executeSwapConfig := &swap.ExecuteSwapConfig{
		WalletKey:     params.WalletKey,
		ChainId:       params.ChainId,
		PublicAddress: params.PublicAddress,
		FromToken:     params.Src,
		ToToken:       params.Dst,
		Amount:        params.Amount,
		Slippage:      params.Slippage,
		SkipWarnings:  params.SkipWarnings,
	}

	var usePermit bool

	// Currently, Permit1 swaps are only tested on Ethereum and Polygon
	isPermitSupportedForCurrentChain := params.ChainId == chains.Ethereum || params.ChainId == chains.Polygon

	var typehash string
	if isPermitSupportedForCurrentChain && params.ApprovalType != swap.ApprovalAlways {
		typehash, err = onchain.GetTypeHash(ethClient, params.Src)
		if err == nil {
			// Typehash is present which means we can use Permit to save gas
			if typehash == typehashes.Permit1 {
				usePermit = true
			} else {
				fmt.Printf("WARN: Typehash exists, but it is not recognized: %v\n", typehash)
			}
		}
	}

	if usePermit || params.ApprovalType == swap.PermitAlways {
		name, err := onchain.ReadContractName(ethClient, common.HexToAddress(params.Src))
		if err != nil {
			return fmt.Errorf("failed to read contract name: %v", err)
		}

		nonce, err := onchain.ReadContractNonce(ethClient, derivedPublicAddress, common.HexToAddress(params.Src))
		if err != nil {
			return fmt.Errorf("failed to read contract name: %v", err)
		}

		sig, err := swap.CreatePermitSignature(&swap.PermitSignatureConfig{
			FromToken:     params.Src,
			Name:          name,
			PublicAddress: params.PublicAddress,
			ChainId:       params.ChainId,
			Key:           params.WalletKey,
			Nonce:         nonce,
			Deadline:      deadline,
		})
		if err != nil {
			return fmt.Errorf("failed to create permit signature: %v", err)
		}

		aggregationRouter, err := contracts.Get1inchRouterFromChainId(params.ChainId)
		if err != nil {
			return fmt.Errorf("failed to get 1inch router address: %v", err)
		}

		permitParams := swap.CreatePermitParams(&swap.PermitParamsConfig{
			Owner:     strings.ToLower(params.PublicAddress), // TODO remove ToLower and see if it still works
			Spender:   aggregationRouter,
			Value:     amounts.BigMaxUint256,
			Deadline:  deadline,
			Signature: sig,
		})

		executeSwapConfig.IsPermitSwap = true
		params.Permit = &permitParams
		fmt.Println("Swapping using Permit1")
	}

	// Execute swap request
	// This will return the transaction data used by a wallet to execute the swap
	swapResponse, _, err := s.client.Swap.GetSwapData(context.Background(), swap.GetSwapDataParams{
		RequestParams:                      params.RequestParams,
		AggregationControllerGetSwapParams: params.AggregationControllerGetSwapParams,
	})
	if err != nil {
		return fmt.Errorf("failed to get swap: %v", err)
	}

	executeSwapConfig.TransactionData = swapResponse.Tx.Data
	executeSwapConfig.EstimatedAmountOut = swapResponse.ToAmount

	err = s.client.Swap.ExecuteSwap(executeSwapConfig)
	if err != nil {
		return fmt.Errorf("failed to execute swap: %v", err)
	}

	return nil
}

// ExecuteSwap executes a swap on the Ethereum blockchain using swap data generated by GetSwapData
func (s *SwapService) ExecuteSwap(config *swap.ExecuteSwapConfig) error {

	if config.WalletKey == "" {
		return fmt.Errorf("wallet key must be set in the client config")
	}

	ethClient, err := s.client.GetEthClient(config.ChainId)
	if err != nil {
		return fmt.Errorf("failed to get eth client: %v", err)
	}

	if !config.SkipWarnings {
		ok, err := swap.ConfirmExecuteSwapWithUser(config, ethClient)
		if err != nil {
			return fmt.Errorf("failed to confirm swap: %v", err)
		}
		if !ok {
			return errors.New("user rejected trade")
		}
	}

	aggregationRouter, err := contracts.Get1inchRouterFromChainId(config.ChainId)
	if err != nil {
		return fmt.Errorf("failed to get 1inch router address: %v", err)
	}

	if !config.IsPermitSwap {
		err = s.executeSwapWithApproval(config, ethClient, aggregationRouter)
		if err != nil {
			return fmt.Errorf("failed to execute swap with approval: %v", err)
		}
	} else {
		err = s.executeSwapWithPermit(config, ethClient)
		if err != nil {
			return fmt.Errorf("failed to execute swap with permit: %v", err)
		}
	}

	return nil
}

func (s *SwapService) executeSwapWithApproval(config *swap.ExecuteSwapConfig, ethClient *ethclient.Client, aggregationRouter string) error {

	var value *big.Int
	var err error
	var approveFirst bool
	if config.FromToken != tokens.NativeToken {
		// When swapping erc20 tokens, the value set on the transaction will be 0
		value = big.NewInt(0)

		allowance, err := onchain.ReadContractAllowance(ethClient, common.HexToAddress(config.FromToken), common.HexToAddress(config.PublicAddress), common.HexToAddress(aggregationRouter))
		if err != nil {
			return fmt.Errorf("failed to read allowance: %v", err)
		}

		amountBig, err := helpers.BigIntFromString(config.Amount)
		if err != nil {
			return fmt.Errorf("failed to convert amount to big.Int: %v", err)
		}
		if allowance.Cmp(amountBig) <= 0 {
			if !config.SkipWarnings {
				ok, err := swap.ConfirmApprovalWithUser(ethClient, config.PublicAddress, config.FromToken)
				if err != nil {
					return fmt.Errorf("failed to confirm approval: %v", err)
				}
				if !ok {
					return errors.New("user rejected approval")
				}
			}

			approveFirst = true

			// Only run the approval if a tenderly key is not present
			if s.client.TenderlyKey == "" {
				erc20Config := onchain.Erc20ApprovalConfig{
					ChainId:        config.ChainId,
					Key:            config.WalletKey,
					Erc20Address:   common.HexToAddress(config.FromToken),
					PublicAddress:  common.HexToAddress(config.PublicAddress),
					SpenderAddress: common.HexToAddress(aggregationRouter),
				}
				err = onchain.ApproveTokenForRouter(ethClient, s.client.NonceCache, erc20Config)
				if err != nil {
					return fmt.Errorf("failed to approve token for router: %v", err)
				}
				helpers.Sleep()
			}
		}
	} else {
		// When swapping from the native token, there is no need for an approval and the amount passed in must be explicitly set
		value, err = helpers.BigIntFromString(config.Amount)
		if err != nil {
			return fmt.Errorf("failed to convert amount to big.Int: %v", err)
		}
	}

	hexData, err := hex.DecodeString(config.TransactionData[2:])
	if err != nil {
		return fmt.Errorf("failed to decode swap data: %v", err)
	}

	txConfig := onchain.TxConfig{
		Description:   "Swap",
		PublicAddress: common.HexToAddress(config.PublicAddress),
		PrivateKey:    config.WalletKey,
		ChainId:       big.NewInt(int64(config.ChainId)),
		Value:         value,
		To:            aggregationRouter,
		Data:          hexData,
	}

	if s.client.TenderlyKey != "" {
		_, err := tenderly.SimulateSwap(s.client.TenderlyKey, tenderly.SwapConfig{
			ChainId:         config.ChainId,
			PublicAddress:   config.PublicAddress,
			FromToken:       config.FromToken,
			TransactionData: config.TransactionData,
			ApproveFirst:    approveFirst,
			Value:           value.String(),
		})
		if err != nil {
			return fmt.Errorf("failed to execute tenderly simulation: %v", err)
		}
	} else {
		err = onchain.ExecuteTransaction(txConfig, ethClient, s.client.NonceCache)
		if err != nil {
			return fmt.Errorf("failed to execute transaction: %v", err)
		}
	}
	return nil
}

func (s *SwapService) executeSwapWithPermit(config *swap.ExecuteSwapConfig, ethClient *ethclient.Client) error {

	hexData, err := hex.DecodeString(config.TransactionData[2:])
	if err != nil {
		return fmt.Errorf("failed to decode swap data: %v", err)
	}

	aggregationRouter, err := contracts.Get1inchRouterFromChainId(config.ChainId)
	if err != nil {
		return fmt.Errorf("failed to get 1inch router address: %v", err)
	}

	txConfig := onchain.TxConfig{
		Description:   "Swap",
		PublicAddress: common.HexToAddress(config.PublicAddress),
		PrivateKey:    config.WalletKey,
		ChainId:       big.NewInt(int64(config.ChainId)),
		Value:         big.NewInt(0),
		To:            aggregationRouter,
		Data:          hexData,
	}
	if s.client.TenderlyKey != "" {
		_, err := tenderly.SimulateSwap(s.client.TenderlyKey, tenderly.SwapConfig{
			ChainId:         config.ChainId,
			PublicAddress:   config.PublicAddress,
			FromToken:       config.FromToken,
			TransactionData: config.TransactionData,
			Value:           "0",
		})
		if err != nil {
			return fmt.Errorf("failed to execute tenderly simulation: %v", err)
		}
	} else {
		err = onchain.ExecuteTransaction(txConfig, ethClient, s.client.NonceCache)
		if err != nil {
			return fmt.Errorf("failed to execute transaction: %v", err)
		}
	}
	return nil
}
