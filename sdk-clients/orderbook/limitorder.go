package orderbook

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/1inch/1inch-sdk-go/internal/hexadecimal"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"github.com/1inch/1inch-sdk-go/constants"
)

func CreateLimitOrderMessage(orderRequest CreateOrderParams, chainId int) (*Order, error) {

	var makerTraitsEncoded string
	if orderRequest.MakerTraits == nil {
		makerTraitsEncoded = ""
	} else {
		makerTraitsEncoded = orderRequest.MakerTraits.Encode()
	}

	orderData := OrderData{
		MakerAsset:    orderRequest.MakerAsset,
		TakerAsset:    orderRequest.TakerAsset,
		MakingAmount:  orderRequest.MakingAmount,
		TakingAmount:  orderRequest.TakingAmount,
		Salt:          orderRequest.Salt,
		Maker:         orderRequest.Maker,
		AllowedSender: "0x0000000000000000000000000000000000000000",
		Receiver:      orderRequest.Taker,
		MakerTraits:   makerTraitsEncoded,
		Extension:     orderRequest.ExtensionEncoded,
	}

	aggregationRouter, err := constants.Get1inchRouterFromChainId(chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get 1inch router address: %w", err)
	}

	// Set up the domain data
	domainData := apitypes.TypedDataDomain{
		Name:              constants.AggregationRouterV6Name,
		Version:           constants.AggregationRouterV6VersionNumber,
		ChainId:           math.NewHexOrDecimal256(int64(chainId)),
		VerifyingContract: aggregationRouter,
	}

	orderMessage := apitypes.TypedDataMessage{
		"salt":         orderData.Salt,
		"makerAsset":   orderData.MakerAsset,
		"takerAsset":   orderData.TakerAsset,
		"maker":        orderData.Maker,
		"receiver":     orderData.Receiver,
		"makingAmount": orderData.MakingAmount,
		"takingAmount": orderData.TakingAmount,
		"makerTraits":  orderData.MakerTraits,
	}

	typedData := apitypes.TypedData{
		Types: map[string][]apitypes.Type{
			"EIP712Domain": {
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"Order": {
				{Name: "salt", Type: "uint256"},
				{Name: "maker", Type: "address"},
				{Name: "receiver", Type: "address"},
				{Name: "makerAsset", Type: "address"},
				{Name: "takerAsset", Type: "address"},
				{Name: "makingAmount", Type: "uint256"},
				{Name: "takingAmount", Type: "uint256"},
				{Name: "makerTraits", Type: "uint256"},
			},
		},
		PrimaryType: "Order",
		Domain: apitypes.TypedDataDomain{
			Name:              domainData.Name,
			Version:           domainData.Version,
			ChainId:           domainData.ChainId,
			VerifyingContract: domainData.VerifyingContract,
		},
		Message: orderMessage,
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to hash typed data: %w", err)
	}
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, fmt.Errorf("failed to hash domain separator: %w", err)
	}

	// Add required prefix to the message
	rawData := []byte{0x19, 0x01}
	rawData = append(rawData, domainSeparator...)
	rawData = append(rawData, typedDataHash...)

	challengeHash := crypto.Keccak256Hash(rawData)
	challengeHashHex := challengeHash.Hex()

	// Sign the challenge hash
	signature, err := orderRequest.Wallet.SignBytes(challengeHash.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to sign challenge hash: %w", err)
	}

	// add 27 to `v` value (last byte)
	signature[64] += 27

	// convert signature to hex string
	signatureHex := fmt.Sprintf("0x%x", signature)

	return &Order{
		OrderHash: challengeHashHex,
		Signature: signatureHex,
		Data:      orderData,
	}, err
}

var timeNow = func() int64 {
	return time.Now().UnixNano()
}

func GenerateSalt(extension string, customBaseSalt *big.Int) (string, error) {
	if extension == "0x" || extension == "" {
		return fmt.Sprintf("%d", timeNow()/int64(time.Millisecond)), nil
	}

	byteConverted, err := stringToHexBytes(extension)
	if err != nil {
		return "", err
	}

	keccakHash := crypto.Keccak256Hash(byteConverted)
	salt := new(big.Int).SetBytes(keccakHash.Bytes())
	// We need to keccak256 the extension and then bitwise & it with uint_160_max
	var uint160Max = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 160), big.NewInt(1))
	salt.And(salt, uint160Max)

	// Convert salt (20 bytes) to byte slice
	saltBytes := salt.Bytes()
	if len(saltBytes) < 20 {
		pad := make([]byte, 20-len(saltBytes))
		saltBytes = append(pad, saltBytes...) // pad to 20 bytes
	}

	var prefixBytes []byte
	if customBaseSalt != nil {
		prefixBytes = customBaseSalt.Bytes()
		if len(prefixBytes) > 12 {
			return "", fmt.Errorf("custom base salt exceeds 96 bits")
		}
		if len(prefixBytes) < 12 {
			pad := make([]byte, 12-len(prefixBytes))
			prefixBytes = append(pad, prefixBytes...)
		}
	} else {
		// Generate random 12-byte prefix
		prefixBytes = make([]byte, 12)
		_, err = rand.Read(prefixBytes)
		if err != nil {
			return "", err
		}
	}

	// Combine random prefix and salt
	full := append(prefixBytes, saltBytes...)

	return fmt.Sprintf("0x%x", full), nil
}

func stringToHexBytes(hexStr string) ([]byte, error) {
	// Strip the "0x" prefix if it exists
	cleanedStr := hexadecimal.Trim0x(hexStr)

	// Ensure the string has an even length by padding with a zero if it's odd
	if len(cleanedStr)%2 != 0 {
		cleanedStr = "0" + cleanedStr
	}

	// Decode the string into bytes
	bytes, err := hex.DecodeString(cleanedStr)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
