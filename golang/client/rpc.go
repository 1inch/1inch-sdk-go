package client

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"1inch-sdk-golang/helpers/consts/amounts"
	"1inch-sdk-golang/helpers/consts/contracts"
)

func (c *Client) ExecuteSwap(fromToken string, swapData string) {
	client, err := ethclient.Dial(c.RpcUrlWithKey)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(c.WalletKey)
	if err != nil {
		log.Fatalf("Failed to convert private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Could not cast public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}

	// Assume spenderAddress is the address of the contract you are giving unlimited approval to spend your USDC
	spenderAddress := common.HexToAddress(contracts.AggregationRouterV5)

	// Parse the USDC contract ABI
	parsedABI, err := abi.JSON(strings.NewReader(contracts.UsdcAbi))
	if err != nil {
		log.Fatalf("Failed to parse USDC ABI: %v", err)
	}

	// Pack the transaction data with the method signature and parameters
	data, err := parsedABI.Pack("approve", spenderAddress, amounts.BigMaxUint256)
	if err != nil {
		log.Fatalf("Failed to pack data for approve: %v", err)
	}

	tx := getDynamicFeeTx(client, chainID, fromAddress, fromToken, data)

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	log.Printf("Transaction sent! Hash: %s\n", signedTx.Hash().Hex())

	receipt, err := waitForTransaction(client, signedTx.Hash())
	if err != nil {
		log.Fatalf("Failed to get transaction receipt: %v", err)
	}
	fmt.Printf("Transaction mined! Block hash: %v\n", receipt.BlockHash)

	hexData, err := hex.DecodeString(swapData[2:])
	if err != nil {
		log.Fatalf("Failed to decode swap data: %v", err)
	}
	tx2 := getDynamicFeeTx(client, chainID, fromAddress, contracts.AggregationRouterV5, hexData)

	// Sign the transaction
	signedTx2, err := types.SignTx(tx2, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), signedTx2)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	log.Printf("Transaction sent! Hash: %s\n", signedTx2.Hash().Hex())

	receipt2, err := waitForTransaction(client, signedTx2.Hash())
	if err != nil {
		log.Fatalf("Failed to get transaction receipt: %v", err)
	}
	fmt.Printf("Transaction mined! Block hash: %v\n", receipt2.BlockHash)
}

func (c *Client) RpcCallEmptyTransaction() {
	client, err := ethclient.Dial(c.RpcUrlWithKey) // TODO make one client and reuse it instead
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	privateKey, err := crypto.HexToECDSA(c.WalletKey)
	if err != nil {
		log.Fatalf("Failed to convert private key: %v", err)
	}

	publicKey := privateKey.Public() // TODO just use public key if possible
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Could not cast public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}

	tx := getDynamicFeeTx(client, chainID, fromAddress, contracts.AggregationRouterV5, []byte{})

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	fmt.Printf("Transaction sent! Hash: %s\n", signedTx.Hash().Hex())
	receipt, err := waitForTransaction(client, signedTx.Hash())
	if err != nil {
		log.Fatalf("Failed to get transaction receipt: %v", err)
	}
	fmt.Printf("Transaction receipt: %#v\n", receipt)
}

func getDynamicFeeTx(client *ethclient.Client, chainID *big.Int, fromAddress common.Address, to string, data []byte) *types.Transaction {
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	fmt.Printf("Current nonce: %v\n", nonce)

	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas tip cap: %v", err)
	}

	//fmt.Printf("Current gas tip cap: %v\n", gasTipCap)
	//
	//gasTipCap.Mul(gasTipCap, big.NewInt(100))

	gasFeeCap, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas fee cap: %v", err)
	}

	toAddress := common.HexToAddress(to)
	value := big.NewInt(0)     // in wei (0 eth)
	gasLimit := uint64(210000) // in units

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

func waitForTransaction(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
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
