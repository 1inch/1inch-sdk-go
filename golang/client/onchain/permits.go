package onchain

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"github.com/1inch/1inch-sdk/golang/helpers/consts/amounts"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/contracts"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/typehashes"
)

func CreatePermitSignature(config *PermitSignatureConfig) (string, error) {
	// Domain Data
	domainData := apitypes.TypedDataDomain{
		Name:              config.Name,
		Version:           "1",
		ChainId:           math.NewHexOrDecimal256(int64(config.ChainId)),
		VerifyingContract: config.FromToken,
	}

	aggregationRouter, err := contracts.Get1inchRouterFromChainId(config.ChainId)
	if err != nil {
		return "", fmt.Errorf("failed to get 1inch router address: %v", err)
	}

	// Order Message
	orderMessage := apitypes.TypedDataMessage{
		"owner":    config.PublicAddress,
		"spender":  aggregationRouter,
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

func ShouldUsePermit(ethClient *ethclient.Client, chainId int, srcToken string) bool {
	// Currently, Permit1 swaps are only tested on Ethereum and Polygon
	isPermitSupportedForCurrentChain := chainId == chains.Ethereum || chainId == chains.Polygon

	if isPermitSupportedForCurrentChain {
		typehash, err := GetTypeHash(ethClient, srcToken)
		if err == nil {
			// If a typehash for Permit1 is present, use that instead of Approve
			if typehash == typehashes.Permit1 {
				return true
			}
		}
	}
	return false
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
