package swap

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"github.com/1inch/1inch-sdk/golang/client/onchain"
	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/amounts"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/contracts"
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

func ConfirmExecuteSwapWithUser(config *ExecuteSwapConfig) (bool, error) {
	stdOut := helpers.StdOutPrinter{}
	return confirmExecuteSwapWithUser(config, os.Stdin, stdOut)
}

func confirmExecuteSwapWithUser(config *ExecuteSwapConfig, reader io.Reader, writer helpers.Printer) (bool, error) {
	var permitType string
	if config.IsPermitSwap {
		permitType = "Permit1"
	} else {
		permitType = "Contract approval"
	}

	writer.Printf("Swap summary:\n")
	writer.Printf("    %-30s %s %s\n", "Selling: ", helpers.SimplifyValue(config.Amount, int(config.FromToken.Decimals)), config.FromToken.Symbol)
	writer.Printf("    %-30s %s %s\n", "Buying (estimation):", helpers.SimplifyValue(config.EstimatedAmountOut, int(config.ToToken.Decimals)), config.ToToken.Symbol)
	writer.Printf("    %-30s %v%s\n", "Slippage:", config.Slippage, "%")
	writer.Printf("    %-30s %s\n", "Permision type:", permitType)
	writer.Printf("\n")
	writer.Printf("WARNING: This swap will be executed onchain next. The results are irreversible. Make sure the proposed trade looks correct before continuing!\n")
	writer.Printf("Would you like to execute this swap onchain now? [y/N]: ")

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

func ConfirmSwapDataWithUser(swapResponse *SwapResponse, fromAmount string, slippage float32) error {
	stdOut := helpers.StdOutPrinter{}
	return confirmSwapDataWithUser(swapResponse, fromAmount, slippage, stdOut)
}

func confirmSwapDataWithUser(swapResponse *SwapResponse, fromAmount string, slippage float32, writer helpers.Printer) error {
	writer.Printf("Swap summary:\n")
	writer.Printf("    %-30s %s %s\n", "Selling: ", helpers.SimplifyValue(fromAmount, int(swapResponse.FromToken.Decimals)), swapResponse.FromToken.Symbol)
	writer.Printf("    %-30s %s %s\n", "Buying (estimation):", helpers.SimplifyValue(swapResponse.ToAmount, int(swapResponse.ToToken.Decimals)), swapResponse.ToToken.Symbol)
	writer.Printf("    %-30s %v%s\n", "Slippage:", slippage, "%")
	writer.Printf("\n")
	writer.Printf("WARNING: Executing the transaction data generated by this function is irreversible. Make sure the proposed trade looks correct!\n")

	return nil
}

func ConfirmApprovalWithUser(ethClient *ethclient.Client, publicAddress string, tokenAddress string) (bool, error) {
	stdOut := helpers.StdOutPrinter{}
	return confirmApprovalWithUser(ethClient, publicAddress, tokenAddress, os.Stdin, stdOut)
}

func confirmApprovalWithUser(ethClient *ethclient.Client, publicAddress string, tokenAddress string, reader io.Reader, writer helpers.Printer) (bool, error) {
	tokenName, err := onchain.ReadContractSymbol(ethClient, common.HexToAddress(tokenAddress))
	if err != nil {
		return false, fmt.Errorf("failed to read name: %v", err)
	}

	writer.Printf("The aggregator contract does not have enough allowance to execute this swap! The SDK can give an " +
		"unlimited approval on your behalf. If you would like to use custom approval amount instead, do that manually " +
		"onchain, then run the SDK again\n")
	writer.Printf("Approval summary:\n")
	writer.Printf("    %-30s %s\n", "Wallet:", publicAddress)
	writer.Printf("    %-30s %s\n", "Swapping: ", tokenName)
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
