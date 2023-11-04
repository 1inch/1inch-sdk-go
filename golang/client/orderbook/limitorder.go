package orderbook

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
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"1inch-sdk-golang/helpers/consts/contracts"
)

type Client struct {
	EthClient *ethclient.Client
}

func Trim0x(data string) string {
	if strings.HasPrefix(data, "0x") {
		return data[2:]
	}
	return data
}

func CumulativeSum(initial int) func(int) int {
	sum := initial
	return func(value int) int {
		sum += value
		return sum
	}
}

var GenerateSalt = func() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
}

const (
	nonceMethod = "nonce"
)

func (c *Client) FetchNonce(address string) (*big.Int, error) {

	parsedNonceManagerAbi, err := abi.JSON(strings.NewReader(nonceManagerAbi))
	if err != nil {
		log.Fatalf("Failed to parse nonceManager ABI: %v\n", err)
	}

	getNonceRequestData, err := parsedNonceManagerAbi.Pack(nonceMethod, big.NewInt(123), common.HexToAddress(address))
	if err != nil {
		log.Fatalf("Failed to pack ABI for %v: %v\n", nonceMethod, err)
	}

	nonceManagerAddressAsAddress := common.HexToAddress(nonceManagerAddress)

	getNonceRequestMessage := ethereum.CallMsg{
		To:   &nonceManagerAddressAsAddress,
		Data: getNonceRequestData,
	}
	nonceResponse, err := c.EthClient.CallContract(context.Background(), getNonceRequestMessage, nil)
	if err != nil {
		log.Fatalf("Failed to call contract: %v\n", err)
	}

	var nonce *big.Int
	err = parsedNonceManagerAbi.UnpackIntoInterface(&nonce, nonceMethod, nonceResponse)
	if err != nil {
		log.Fatalf("Failed to unpack data for %v: %v\n", nonceMethod, err)
	}

	return nonce, nil
}

func CreateLimitOrder(orderRequest OrderRequest, chainId int, key string) (*Order, error) {

	orderData := OrderData{
		MakerAsset:    orderRequest.FromToken,
		TakerAsset:    orderRequest.ToToken,
		MakingAmount:  fmt.Sprintf("%d", orderRequest.MakingAmount),
		TakingAmount:  fmt.Sprintf("%d", orderRequest.TakingAmount),
		Salt:          GenerateSalt(),
		Maker:         orderRequest.SourceWallet,
		AllowedSender: "0x0000000000000000000000000000000000000000", // TODO use this
		Receiver:      orderRequest.Receiver,
		Offsets:       "0",  // TODO use this
		Interactions:  "0x", // TODO use this
	}

	// Set up the domain data
	domainData := apitypes.TypedDataDomain{
		Name:              contracts.AggregationRouterV5Name,
		Version:           contracts.AggregationRouterV5VersionNumber,
		ChainId:           math.NewHexOrDecimal256(int64(chainId)),
		VerifyingContract: contracts.AggregationRouterV5,
	}

	orderMessage := apitypes.TypedDataMessage{
		"salt":          orderData.Salt,
		"makerAsset":    orderData.MakerAsset,
		"takerAsset":    orderData.TakerAsset,
		"maker":         orderData.Maker,
		"receiver":      orderData.Receiver,
		"allowedSender": orderData.AllowedSender,
		"makingAmount":  orderData.MakingAmount,
		"takingAmount":  orderData.TakingAmount,
		"offsets":       orderData.Offsets,
		"interactions":  common.FromHex(orderData.Interactions),
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
				{Name: "makerAsset", Type: "address"},
				{Name: "takerAsset", Type: "address"},
				{Name: "maker", Type: "address"},
				{Name: "receiver", Type: "address"},
				{Name: "allowedSender", Type: "address"},
				{Name: "makingAmount", Type: "uint256"},
				{Name: "takingAmount", Type: "uint256"},
				{Name: "offsets", Type: "uint256"},
				{Name: "interactions", Type: "bytes"},
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

	privateKey, err := crypto.HexToECDSA(key)
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
