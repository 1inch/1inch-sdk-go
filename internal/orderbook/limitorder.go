package orderbook

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"github.com/1inch/1inch-sdk-go/client/models"
	"github.com/1inch/1inch-sdk-go/helpers"
	"github.com/1inch/1inch-sdk-go/helpers/consts/contracts"
	"github.com/1inch/1inch-sdk-go/internal/onchain"
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

func BuildMakerTraits(params models.BuildMakerTraitsParams) string {
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

func CreateLimitOrderMessage(orderRequest models.CreateOrderParams, makerTraits string) (*models.Order, error) {

	orderData := models.OrderData{
		MakerAsset:    orderRequest.MakerAsset,
		TakerAsset:    orderRequest.TakerAsset,
		MakingAmount:  orderRequest.MakingAmount,
		TakingAmount:  orderRequest.TakingAmount,
		Salt:          GenerateSalt(),
		Maker:         orderRequest.Maker,
		AllowedSender: "0x0000000000000000000000000000000000000000",
		Receiver:      orderRequest.Taker,
		MakerTraits:   makerTraits,
		Extension:     "0x",
	}

	aggregationRouter, err := contracts.Get1inchRouterFromChainId(orderRequest.ChainId)
	if err != nil {
		return nil, fmt.Errorf("failed to get 1inch router address: %v", err)
	}

	// Set up the domain data
	domainData := apitypes.TypedDataDomain{
		Name:              contracts.AggregationRouterV6Name,
		Version:           contracts.AggregationRouterV6VersionNumber,
		ChainId:           math.NewHexOrDecimal256(int64(orderRequest.ChainId)),
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

	return &models.Order{
		OrderHash: challengeHashHex,
		Signature: signatureHex,
		Data:      orderData,
	}, err
}

func ConfirmLimitOrderWithUser(order *models.Order, ethClient *ethclient.Client) (bool, error) {
	stdOut := helpers.StdOutPrinter{}
	return confirmLimitOrderWithUser(order, ethClient, os.Stdin, stdOut)
}

func confirmLimitOrderWithUser(order *models.Order, ethClient *ethclient.Client, reader io.Reader, writer helpers.Printer) (bool, error) {
	makerTokenDecimals, err := onchain.ReadContractDecimals(ethClient, common.HexToAddress(order.Data.MakerAsset))
	if err != nil {
		return false, fmt.Errorf("failed to read decimals: %v", err)
	}

	makerTokenName, err := onchain.ReadContractSymbol(ethClient, common.HexToAddress(order.Data.MakerAsset))
	if err != nil {
		return false, fmt.Errorf("failed to read name: %v", err)
	}

	takerTokenDecimals, err := onchain.ReadContractDecimals(ethClient, common.HexToAddress(order.Data.TakerAsset))
	if err != nil {
		return false, fmt.Errorf("failed to read decimals: %v", err)
	}

	takerTokenName, err := onchain.ReadContractSymbol(ethClient, common.HexToAddress(order.Data.TakerAsset))
	if err != nil {
		return false, fmt.Errorf("failed to read name: %v", err)
	}

	writer.Printf("Order summary:\n")
	writer.Printf("    %-30s %s\n", "Wallet:", order.Data.Maker)
	writer.Printf("    %-30s %s %s\n", "Selling: ", helpers.SimplifyValue(order.Data.MakingAmount, int(makerTokenDecimals)), makerTokenName)
	writer.Printf("    %-30s %s %s\n", "Buying:", helpers.SimplifyValue(order.Data.TakingAmount, int(takerTokenDecimals)), takerTokenName)
	writer.Printf("\n")
	writer.Printf("WARNING: This order will be officially posted to the 1inch Limit Order protocol where anyone will be able to execute in onchain immediately. " +
		"Once executed, the results are irreversible. Make sure the proposed trade looks correct before continuing!\n")
	writer.Printf("Would you like to post this order to the 1inch API now? [y/N]: ")

	inputReader := bufio.NewReader(reader)
	input, _ := inputReader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	switch input {
	case "y":
		return true, nil
	default:
		return false, nil
	}
}

func ConfirmApprovalWithUser(ethClient *ethclient.Client, publicAddress string, tokenAddress string) (bool, error) {
	stdOut := helpers.StdOutPrinter{}
	return confirmApprovalWithUser(ethClient, publicAddress, tokenAddress, os.Stdin, stdOut)
}

func confirmApprovalWithUser(ethClient *ethclient.Client, publicAddress string, tokenAddress string, reader io.Reader, writer helpers.Printer) (bool, error) {
	tokenSymbol, err := onchain.ReadContractSymbol(ethClient, common.HexToAddress(tokenAddress))
	if err != nil {
		return false, fmt.Errorf("failed to read name: %v", err)
	}

	writer.Printf("The aggregator contract does not have enough allowance to execute the order! The SDK can give an " +
		"unlimited approval on your behalf. If you would like to use custom approval amount instead, do that manually " +
		"onchain, then run the SDK again\n")
	writer.Printf("Approval summary:\n")
	writer.Printf("    %-30s %s\n", "Wallet:", publicAddress)
	writer.Printf("    %-30s %s\n", "Selling: ", tokenSymbol)
	writer.Printf("    %-30s %s\n", "Approval amount: ", "unlimited")
	writer.Printf("\n")
	writer.Printf("Would you like post an onchain unlimited approval now? [y/N]: ")

	inputReader := bufio.NewReader(reader)
	input, _ := inputReader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	switch input {
	case "y":
		return true, nil
	default:
		return false, nil
	}
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

func concatenateInteractions(interactions []string) string {
	interactionsConcatenated := "0x"
	for _, interaction := range interactions {
		interactionsConcatenated += Trim0x(interaction)
	}
	return interactionsConcatenated
}

var GenerateSalt = func() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
}
