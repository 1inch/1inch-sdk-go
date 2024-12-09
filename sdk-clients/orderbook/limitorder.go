package orderbook

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"github.com/1inch/1inch-sdk-go/constants"
)

func CreateLimitOrderMessage(orderRequest CreateOrderParams, chainId int) (*Order, error) {

	encodedExtension, err := orderRequest.Extension.Encode()
	if err != nil {
		return nil, fmt.Errorf("error encoding extension: %v", err)
	}

	orderData := OrderData{
		MakerAsset:    orderRequest.MakerAsset,
		TakerAsset:    orderRequest.TakerAsset,
		MakingAmount:  orderRequest.MakingAmount,
		TakingAmount:  orderRequest.TakingAmount,
		Salt:          GenerateSalt(encodedExtension),
		Maker:         orderRequest.Maker,
		AllowedSender: "0x0000000000000000000000000000000000000000",
		Receiver:      orderRequest.Taker,
		MakerTraits:   orderRequest.MakerTraits.Encode(),
		Extension:     encodedExtension,
	}

	aggregationRouter, err := constants.Get1inchRouterFromChainId(chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get 1inch router address: %v", err)
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
		return nil, fmt.Errorf("error hashing typed data: %v", err)
	}
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, fmt.Errorf("error hashing domain separator: %v", err)
	}

	// Add required prefix to the message
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))

	challengeHash := crypto.Keccak256Hash(rawData)
	challengeHashHex := challengeHash.Hex()

	// Sign the challenge hash
	signature, err := orderRequest.Wallet.SignBytes(challengeHash.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error signing challenge hash: %v", err)
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

func GenerateSalt(extension string) string {
	if extension == "0x" {
		return fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
	}

	byteConverted, err := stringToHexBytes(extension)
	if err != nil {
		panic(err)
	}

	keccakHash := crypto.Keccak256Hash(byteConverted)
	salt := new(big.Int).SetBytes(keccakHash.Bytes())
	// We need to keccak256 the extension and then bitwise & it with uint_160_max
	var uint160Max = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 160), big.NewInt(1))
	salt.And(salt, uint160Max)
	return fmt.Sprintf("0x%x", salt)
}

func stringToHexBytes(hexStr string) ([]byte, error) {
	// Strip the "0x" prefix if it exists
	cleanedStr := strings.TrimPrefix(hexStr, "0x")

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
