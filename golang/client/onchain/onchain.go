package onchain

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/abis"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/amounts"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
)

const gasLimit = uint64(21000000) // TODO make sure this value more dynamic

// TODO: this nonce value will compete with any pending transactions on the wallet. The user should be able to set this if they want

func GetNonce(ethClient *ethclient.Client, key string, publicAddress common.Address, nonceCache map[string]uint64) (uint64, error) {
	var err error
	nonce, ok := nonceCache[key]
	if !ok {
		nonce, err = ethClient.NonceAt(context.Background(), publicAddress, nil)
		if err != nil {
			return 0, fmt.Errorf("failed to get nonce: %v", err)
		}
		nonceCache[key] = nonce
	}
	return nonce, nil
}

func ExecuteTransaction(txConfig TxConfig, ethClient *ethclient.Client, nonceCache map[string]uint64) error {

	nonceCacheKey := fmt.Sprintf("%s+%d", txConfig.PublicAddress, txConfig.ChainId.Int64())
	nonce, err := GetNonce(ethClient, nonceCacheKey, txConfig.PublicAddress, nonceCache)

	swapTx, err := GetTx(ethClient, nonce, txConfig)

	signingKey, err := crypto.HexToECDSA(txConfig.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to convert private key: %v", err)
	}

	// Sign the transaction
	swapTxSigned, err := types.SignTx(swapTx, types.LatestSignerForChainID(txConfig.ChainId), signingKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send the transaction
	err = ethClient.SendTransaction(context.Background(), swapTxSigned)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}

	fmt.Printf("Transaction sent! (%s)\n", txConfig.Description)
	helpers.PrintBlockExplorerTxLink(int(txConfig.ChainId.Int64()), swapTxSigned.Hash().String())

	_, err = WaitForTransaction(ethClient, swapTxSigned.Hash())
	if err != nil {
		return fmt.Errorf("failed to get transaction receipt: %v", err)
	}

	// Update cache to avoid RPC nonce desync
	nonceCache[nonceCacheKey] = nonce + 1

	return nil
}

func GetTx(client *ethclient.Client, nonce uint64, config TxConfig) (*types.Transaction, error) {
	chainIdInt := int(config.ChainId.Int64())
	if chainIdInt == chains.Ethereum || chainIdInt == chains.Polygon {
		return GetDynamicFeeTx(client, nonce, config.ChainId, config.To, config.Value, config.Data)
	} else {
		return GetLegacyTx(client, nonce, config.To, config.Value, config.Data)
	}
}

func GetDynamicFeeTx(client *ethclient.Client, nonce uint64, chainID *big.Int, to string, value *big.Int, data []byte) (*types.Transaction, error) {

	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas tip cap: %v", err)
	}

	gasFeeCap, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas fee cap: %v", err)
	}

	//// Increase the gas fee cap by 25%
	gasFeeCap = gasFeeCap.Mul(gasFeeCap, big.NewInt(150))
	gasFeeCap = gasFeeCap.Div(gasFeeCap, big.NewInt(100))
	//// Increase the gas tip cap by 25%
	gasTipCap = gasTipCap.Mul(gasTipCap, big.NewInt(150))
	gasTipCap = gasTipCap.Div(gasTipCap, big.NewInt(100))

	toAddress := common.HexToAddress(to)

	return types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     value,
		Data:      data,
	}), nil
}

func GetLegacyTx(client *ethclient.Client, nonce uint64, to string, value *big.Int, data []byte) (*types.Transaction, error) {

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %v", err)
	}

	// Increase the gas fee cap by 25%
	gasPrice = gasPrice.Mul(gasPrice, big.NewInt(125))
	gasPrice = gasPrice.Div(gasPrice, big.NewInt(100))

	toAddress := common.HexToAddress(to)

	return types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &toAddress,
		Value:    value,
		Data:     data,
	}), nil
}

// ReadContractName reads the 'name' public variable from a contract.
func ReadContractName(client *ethclient.Client, contractAddress common.Address) (string, error) {
	parsedABI, err := abi.JSON(strings.NewReader(abis.Erc20)) // Make a generic version of this ABI
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

// ReadContractSymbol reads the 'symbol' public variable from a contract.
func ReadContractSymbol(client *ethclient.Client, contractAddress common.Address) (string, error) {
	parsedABI, err := abi.JSON(strings.NewReader(abis.Erc20)) // Make a generic version of this ABI
	if err != nil {
		return "", err
	}

	// Construct the call message
	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: parsedABI.Methods["symbol"].ID,
	}

	// Query the blockchain
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "", err
	}

	// Unpack the result
	var contractName string
	err = parsedABI.UnpackIntoInterface(&contractName, "symbol", result)
	if err != nil {
		return "", err
	}

	return contractName, nil
}

// ReadContractDecimals reads the 'decimals' public variable from a contract.
func ReadContractDecimals(client *ethclient.Client, contractAddress common.Address) (uint8, error) {
	parsedABI, err := abi.JSON(strings.NewReader(abis.Erc20)) // Make a generic version of this ABI
	if err != nil {
		return 0, err
	}

	// Construct the call message
	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: parsedABI.Methods["decimals"].ID,
	}

	// Query the blockchain
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return 0, err
	}

	// Unpack the result
	var decimals uint8
	err = parsedABI.UnpackIntoInterface(&decimals, "decimals", result)
	if err != nil {
		return 0, err
	}

	return decimals, nil
}

// ReadContractNonce reads the 'nonces' public variable from a contract.
func ReadContractNonce(client *ethclient.Client, publicAddress common.Address, contractAddress common.Address) (int64, error) {
	parsedABI, err := abi.JSON(strings.NewReader(abis.Erc20)) // Make a generic version of this ABI
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

// ReadContractAllowance reads the allowance a given contract has for a wallet.
func ReadContractAllowance(client *ethclient.Client, erc20Address common.Address, publicAddress common.Address, spenderAddress common.Address) (*big.Int, error) {
	parsedABI, err := abi.JSON(strings.NewReader(abis.Erc20)) // Make a generic version of this ABI
	if err != nil {
		return nil, err
	}

	data, err := parsedABI.Pack("allowance", publicAddress, spenderAddress)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &erc20Address,
		Data: data,
	}

	// Query the blockchain
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	// Unpack the result
	var allowance *big.Int
	err = parsedABI.UnpackIntoInterface(&allowance, "allowance", result)
	if err != nil {
		return nil, err
	}

	return allowance, nil
}

func GetTypeHash(client *ethclient.Client, addressAsString string) (string, error) { // Pack the call to get the PERMIT_TYPEHASH constant

	// Parse the ABI
	parsedABI, err := abi.JSON(strings.NewReader(abis.Erc20))
	if err != nil {
		return "", fmt.Errorf("failed to parse contract ABI: %v", err)
	}

	// Create the contract address
	address := common.HexToAddress(addressAsString)

	data, err := parsedABI.Pack("PERMIT_TYPEHASH")
	if err != nil {
		return "", fmt.Errorf("failed to pack data for PERMIT_TYPEHASH: %v", err)
	}

	// Create the call message
	msg := ethereum.CallMsg{
		To:   &address,
		Data: data,
	}

	// Query the blockchain
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve the PERMIT_TYPEHASH: %v", err)
	}

	// Convert the result to bytes32
	var typeHash [32]byte
	copy(typeHash[:], result)

	// Convert the result to a string
	resultAsString := fmt.Sprintf("%x", typeHash)
	// If the varaible does not exist, it will be all zeros
	if string(resultAsString) == "0000000000000000000000000000000000000000000000000000000000000000" {
		return "", errors.New("PERMIT_TYPEHASH does not exist")
	}

	return resultAsString, nil
}

func ApproveTokenForRouter(client *ethclient.Client, nonceCache map[string]uint64, config Erc20ApprovalConfig) error {
	// Parse the USDC contract ABI to get the 'Approve' function signature
	parsedABI, err := abi.JSON(strings.NewReader(abis.Erc20))
	if err != nil {
		return fmt.Errorf("failed to parse USDC ABI: %v", err)
	}

	// Pack the transaction data with the method signature and parameters
	data, err := parsedABI.Pack("approve", config.SpenderAddress, amounts.BigMaxUint256)
	if err != nil {
		return fmt.Errorf("failed to pack data for approve: %v", err)
	}

	txConfig := TxConfig{
		Description:   "Approval",
		PublicAddress: config.PublicAddress,
		PrivateKey:    config.Key,
		ChainId:       big.NewInt(int64(config.ChainId)),
		Value:         big.NewInt(0),
		To:            config.Erc20Address.Hex(),
		Data:          data,
	}
	err = ExecuteTransaction(txConfig, client, nonceCache)
	if err != nil {
		return fmt.Errorf("failed to execute transaction: %v", err)
	}
	return nil
}

func GetTimestampBelowCalldata(expiration int64) ([]byte, error) {

	expiration = time.Now().UnixMilli()

	parsedABI, err := abi.JSON(strings.NewReader(abis.SeriesNonceManager))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	// Pack the transaction data with the method signature and parameters
	data, err := parsedABI.Pack("timestampBelow", big.NewInt(expiration))
	if err != nil {
		return nil, fmt.Errorf("failed to pack data: %v", err)
	}

	return data, nil
}

func GetTimeSeriesManagerNonce(client *ethclient.Client, seriesNonceManager string, publicAddress string) (*big.Int, error) {

	function := "nonce"

	parsedABI, err := abi.JSON(strings.NewReader(abis.SeriesNonceManager))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	// Pack the transaction data with the method signature and parameters
	data, err := parsedABI.Pack(function, big.NewInt(0), common.HexToAddress(publicAddress))
	if err != nil {
		return nil, fmt.Errorf("failed to pack data: %v", err)
	}

	address := common.HexToAddress(seriesNonceManager)

	// Create the call message
	msg := ethereum.CallMsg{
		To:   &address,
		Data: data,
	}

	// Query the blockchain
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the PERMIT_TYPEHASH: %v", err)
	}

	// Unpack the result
	var nonce *big.Int
	err = parsedABI.UnpackIntoInterface(&nonce, function, result)
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

func GetTimestampBelowAndNonceEqualsCalldata(expiration int64, nonce *big.Int, publicAddress string) ([]byte, error) {
	var (
		timeInt    = new(big.Int).SetUint64(uint64(expiration))
		nonceInt   = new(big.Int).SetUint64(nonce.Uint64())
		seriesInt  = new(big.Int).SetUint64(uint64(0)) // Limit orders have a static series of 0
		accountInt = new(big.Int)
	)

	accountInt.SetString(publicAddress[2:], 16)

	timeInt.Lsh(timeInt, 216)
	nonceInt.Lsh(nonceInt, 176)
	seriesInt.Lsh(seriesInt, 160)

	result := new(big.Int)
	result.Or(result, timeInt)
	result.Or(result, nonceInt)
	result.Or(result, seriesInt)
	result.Or(result, accountInt)

	parsedABI, err := abi.JSON(strings.NewReader(abis.SeriesNonceManager))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	data, err := parsedABI.Pack("timestampBelowAndNonceEquals", result)
	if err != nil {
		return nil, fmt.Errorf("failed to pack data: %v", err)
	}

	return data, nil
}

func GetPredicateCalldata(seriesNonceManager string, getTimestampBelowAndNonceEqualsCalldata []byte) ([]byte, error) {
	parsedABI, err := abi.JSON(strings.NewReader(abis.AggregationRouterV5))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	// Pack the transaction data with the method signature and parameters
	data, err := parsedABI.Pack("arbitraryStaticCall", common.HexToAddress(seriesNonceManager), getTimestampBelowAndNonceEqualsCalldata)
	if err != nil {
		return nil, fmt.Errorf("failed to pack data: %v", err)
	}

	return data, nil
}

func RevokeApprovalForRouter(client *ethclient.Client, nonceCache map[string]uint64, config Erc20RevokeConfig) error {
	// Parse the USDC contract ABI to get the 'Approve' function signature
	parsedABI, err := abi.JSON(strings.NewReader(abis.Erc20))
	if err != nil {
		return fmt.Errorf("failed to parse ABI: %v", err)
	}

	// Pack the transaction data with the method signature and parameters
	data, err := parsedABI.Pack("decreaseAllowance", config.SpenderAddress, config.AllowanceDecreaseAmount)
	if err != nil {
		return fmt.Errorf("failed to pack data for approve: %v", err)
	}

	txConfig := TxConfig{
		Description:   "Revoke Approval",
		PublicAddress: config.PublicAddress,
		PrivateKey:    config.Key,
		ChainId:       big.NewInt(int64(config.ChainId)),
		Value:         big.NewInt(0),
		To:            config.Erc20Address.Hex(),
		Data:          data,
	}
	err = ExecuteTransaction(txConfig, client, nonceCache)
	if err != nil {
		return fmt.Errorf("failed to execute transaction: %v", err)
	}
	return nil
}

func WaitForTransaction(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if receipt != nil {
			fmt.Println("Transaction complete!")
			return receipt, nil
		}
		if err != nil {
			fmt.Println("Waiting for transaction to be mined")
		}
		select {
		case <-time.After(1000 * time.Millisecond): // check again after a delay
		case <-context.Background().Done():
			fmt.Println("Context cancelled")
			return nil, context.Background().Err()
		}
	}
}
