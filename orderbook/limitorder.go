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

const (
	unwrapWethFlag          = 247
	allowMultipleFillsFlag  = 254
	needEpochCheckFlag      = 250
	usePermit2Flag          = 248
	hasExtensionFlag        = 249
	needPreinteractionFlag  = 252
	needPostinteractionFlag = 251
)

func BuildMakerTraits(params BuildMakerTraitsParams) string {
	// Convert allowedSender from hex string to big.Int
	allowedSenderInt := new(big.Int)
	allowedSenderInt.SetString(params.AllowedSender, 16)

	// Initialize tempPredicate as big.Int
	tempPredicate := new(big.Int)
	tempPredicate.Lsh(big.NewInt(params.Series), 160)
	tempPredicate.Or(tempPredicate, new(big.Int).Lsh(big.NewInt(params.Nonce), 120))
	tempPredicate.Or(tempPredicate, new(big.Int).Lsh(big.NewInt(params.Expiry), 80))
	tempPredicate.Or(tempPredicate, new(big.Int).And(allowedSenderInt, new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 80), big.NewInt(1))))

	if params.UnwrapWeth {
		tempPredicate.Or(tempPredicate, big.NewInt(1).Lsh(big.NewInt(1), unwrapWethFlag))
	}
	// This flag must be set
	tempPredicate.Or(tempPredicate, big.NewInt(1).Lsh(big.NewInt(1), allowMultipleFillsFlag))

	if params.ShouldCheckEpoch {
		tempPredicate.Or(tempPredicate, big.NewInt(1).Lsh(big.NewInt(1), needEpochCheckFlag))
	}
	if params.UsePermit2 {
		tempPredicate.Or(tempPredicate, big.NewInt(1).Lsh(big.NewInt(1), usePermit2Flag))
	}
	if params.HasExtension {
		tempPredicate.Or(tempPredicate, big.NewInt(1).Lsh(big.NewInt(1), hasExtensionFlag))
	}
	if params.HasPreInteraction {
		tempPredicate.Or(tempPredicate, big.NewInt(1).Lsh(big.NewInt(1), needPreinteractionFlag))
	}
	if params.HasPostInteraction {
		tempPredicate.Or(tempPredicate, big.NewInt(1).Lsh(big.NewInt(1), needPostinteractionFlag))
	}

	// Pad the predicate to 32 bytes with 0's on the left and convert to hex string
	paddedPredicate := fmt.Sprintf("%032x", tempPredicate)
	return "0x" + paddedPredicate
}

func CreateLimitOrderMessage(orderRequest CreateOrderParams, chainId int) (*Order, error) {

	orderData := OrderData{
		MakerAsset:    orderRequest.MakerAsset,
		TakerAsset:    orderRequest.TakerAsset,
		MakingAmount:  orderRequest.MakingAmount,
		TakingAmount:  orderRequest.TakingAmount,
		Salt:          GenerateSalt(orderRequest.Extension),
		Maker:         orderRequest.Maker,
		AllowedSender: "0x0000000000000000000000000000000000000000",
		Receiver:      orderRequest.Taker,
		MakerTraits:   orderRequest.MakerTraits,
		Extension:     orderRequest.Extension,
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

	privateKey, err := crypto.HexToECDSA(orderRequest.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error converting private key to ECDSA: %v", err)
	}

	// Sign the challenge hash
	signature, err := crypto.Sign(challengeHash.Bytes(), privateKey)
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

func GetInteractions(makerAsset string, permit string) ([]string, error) {

	makerAssetData := `0x`
	takerAssetData := `0x`
	getMakingAmount := `0x`
	getTakingAmount := `0x`
	predicate := `0x`
	preInteraction := `0x`
	postInteraction := `0x`

	// The maker token must be prepended to permit data for limit orders
	if permit != "0x" {
		permit = makerAsset + permit
	}

	return []string{makerAssetData, takerAssetData, getMakingAmount, getTakingAmount, predicate, permit, preInteraction, postInteraction}, nil
}

func GetOffsets(interactions []string) *big.Int {
	var lengthMap []int
	for _, interaction := range interactions {
		if interaction[:2] == "0x" {
			lengthMap = append(lengthMap, len(interaction)/2-1)
		} else {
			lengthMap = append(lengthMap, len(interaction)/2)
		}
	}

	cumulativeSum := 0
	bytesAccumulator := big.NewInt(0)
	var index uint64

	for _, length := range lengthMap {
		cumulativeSum += length
		shiftVal := big.NewInt(int64(cumulativeSum))
		shiftVal.Lsh(shiftVal, uint(32*index))           // Shift left
		bytesAccumulator.Add(bytesAccumulator, shiftVal) // Add to accumulator
		index++
	}

	return bytesAccumulator
}

func BuildExtension(interactionsConcatednated string, offsets *big.Int) string {
	if interactionsConcatednated == "0x" {
		return "0x"
	}
	offsetsBytes := offsets.Bytes()
	paddedOffsetHex := fmt.Sprintf("%064x", offsetsBytes)
	return "0x" + paddedOffsetHex + strings.TrimPrefix(interactionsConcatednated, "0x")
}

func ConcatenateInteractions(interactions []string) string {
	var builder strings.Builder

	for _, interaction := range interactions {
		// Remove "0x" prefix if present
		interaction = strings.TrimPrefix(interaction, "0x")
		builder.WriteString(interaction)
	}

	// Add "0x" prefix to the final result
	return fmt.Sprintf("%s", builder.String())
}
