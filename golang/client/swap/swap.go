package swap

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"1inch-sdk-golang/helpers/consts/amounts"
	"1inch-sdk-golang/helpers/consts/contracts"
)

type PermitSignatureConfig struct {
	FromToken     string
	Name          string
	PublicAddress string
	ChainId       int
	Key           string
	Nonce         int64
	Deadline      int64
}

func CreatePermitSignature(config *PermitSignatureConfig) (string, error) {
	// Domain Data
	domainData := apitypes.TypedDataDomain{
		Name:              config.Name,
		Version:           "1",
		ChainId:           math.NewHexOrDecimal256(int64(config.ChainId)),
		VerifyingContract: config.FromToken,
	}

	// Order Message
	orderMessage := apitypes.TypedDataMessage{
		"owner":    config.PublicAddress,
		"spender":  contracts.AggregationRouterV5,
		"value":    amounts.BigMaxUint256,
		"nonce":    big.NewInt(config.Nonce),
		"deadline": big.NewInt(config.Deadline),
	}

	// Typed Data
	typedData := apitypes.TypedData{
		Types: map[string][]apitypes.Type{
			"EIP712Domain": {
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"Permit": {
				{Name: "owner", Type: "address"},
				{Name: "spender", Type: "address"},
				{Name: "value", Type: "uint256"},
				{Name: "nonce", Type: "uint256"},
				{Name: "deadline", Type: "uint256"},
			},
		},
		PrimaryType: "Permit",
		Domain:      domainData,
		Message:     orderMessage,
	}

	// Hash the data
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return "", fmt.Errorf("error hashing typed data: %v", err)
	}
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return "", fmt.Errorf("error hashing domain separator: %v", err)
	}

	// Prepare the data for signing
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	challengeHash := crypto.Keccak256Hash(rawData)

	// Convert private key and sign
	privateKey, err := crypto.HexToECDSA(config.Key)
	if err != nil {
		return "", fmt.Errorf("error converting private key to ECDSA: %v", err)
	}
	signature, err := crypto.Sign(challengeHash.Bytes(), privateKey)
	if err != nil {
		return "", fmt.Errorf("error signing challenge hash: %v", err)
	}
	signature[64] += 27 // Adjust the `v` value

	// Convert signature to hex string
	signatureHex := fmt.Sprintf("0x%x", signature)
	return signatureHex, nil
}

type PermitParamsConfig struct {
	Owner     string
	Spender   string
	Value     *big.Int
	Deadline  int64
	Signature string
}

func CreatePermitParams(config *PermitParamsConfig) string {

	// example of a separated permit string
	// 000000000000000000000000xxxxxxxxxxx980858eb298a0843264cff21fd9c9 // owner
	// 0000000000000000000000001111111254eeb25477b68fb85ed929f73a960582 // spender
	// 0000000000000000000000000000000000c097ce7bc90715b34b9f1000000000 // value
	// 0000000000000000000000000000000000000000000000000000000063ada9c0 // deadline
	// 000000000000000000000000000000000000000000000000000000000000001b // v
	// 04dd10d79a8b12a5a93606f6872bb5b25ba3e41609be79409032f9dc6738792b // r
	// 08e0318c0dcd4ec8e3309ac0ff46d52d25e43369611402bc1ddd01fe0602ee56 // s

	ownerNoPrefix := Remove0xPrefix(config.Owner)
	spenderNoPrefix := Remove0xPrefix(config.Spender)
	signatureNoPrefix := Remove0xPrefix(config.Signature)

	return "0x" + padStringWithZeroes(ownerNoPrefix) +
		padStringWithZeroes(spenderNoPrefix) +
		padStringWithZeroes(fmt.Sprintf("%x", config.Value)) +
		padStringWithZeroes(fmt.Sprintf("%x", config.Deadline)) +
		ConvertSignatureToVRSString(signatureNoPrefix)
}

func padStringWithZeroes(s string) string {
	if len(s) >= 64 {
		return s
	}
	return strings.Repeat("0", 64-len(s)) + s
}

func Remove0xPrefix(s string) string {
	if strings.HasPrefix(s, "0x") {
		return s[2:]
	}
	return s
}

// ConvertSignatureToVRSString converts a signature from rsv to padded vrs format
func ConvertSignatureToVRSString(signature string) string {
	// explicit breakdown
	//r := signature[:66]
	//s := signature[66:128]
	//v := signature[128:]
	return padStringWithZeroes(signature[128:]) + signature[:128]
}

func GetTypeHash(client *ethclient.Client, addressAsString string) (string, error) { // Pack the call to get the PERMIT_TYPEHASH constant

	// Parse the ABI
	parsedABI, err := abi.JSON(strings.NewReader(contracts.Erc20Abi))
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
