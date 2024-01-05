package onchain

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"1inch-sdk-golang/helpers/consts/contracts"
)

func GetDynamicFeeTx(client *ethclient.Client, chainID *big.Int, fromAddress common.Address, to string, data []byte) *types.Transaction {
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	fmt.Printf("Current nonce: %v\n", nonce)

	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas tip cap: %v", err)
	}

	gasFeeCap, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas fee cap: %v", err)
	}

	toAddress := common.HexToAddress(to)
	value := big.NewInt(0)     // in wei (0 eth)
	gasLimit := uint64(210000) // TODO make sure this value is always correct

	return types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     value,
		Data:      data,
	})
}

// ReadContractName reads the 'name' public variable from a contract.
func ReadContractName(client *ethclient.Client, contractAddress common.Address) (string, error) {
	parsedABI, err := abi.JSON(strings.NewReader(contracts.Erc20Abi)) // Make a generic version of this ABI
	if err != nil {
		return "", err
	}

	// Construct the call message
	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: parsedABI.Methods["name"].ID,
	}

	// Query the blockchain
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "", err
	}

	// Unpack the result
	var contractName string
	err = parsedABI.UnpackIntoInterface(&contractName, "name", result)
	if err != nil {
		return "", err
	}

	return contractName, nil
}

// ReadContractNonce reads the 'nonces' public variable from a contract.
func ReadContractNonce(client *ethclient.Client, publicAddress common.Address, contractAddress common.Address) (int64, error) {
	parsedABI, err := abi.JSON(strings.NewReader(contracts.Erc20Abi)) // Make a generic version of this ABI
	if err != nil {
		return -1, err
	}

	data, err := parsedABI.Pack("nonces", publicAddress)
	if err != nil {
		return -1, err
	}

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	// Query the blockchain
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return -1, err
	}

	// Unpack the result
	var nonce *big.Int
	err = parsedABI.UnpackIntoInterface(&nonce, "nonces", result)
	if err != nil {
		return -1, err
	}

	return nonce.Int64(), nil
}

func WaitForTransaction(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if receipt != nil {
			return receipt, nil
		}
		if err != nil {
			fmt.Println("Transaction not yet mined...")
		}
		select {
		case <-time.After(1 * time.Second): // check again after a delay
		case <-context.Background().Done():
			return nil, context.Background().Err()
		}
	}
}
