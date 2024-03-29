package web3_provider

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"github.com/1inch/1inch-sdk-go/internal/common"
	"github.com/1inch/1inch-sdk-go/internal/helpers/consts/abis"
)

// TokenPermit Will return an erc2612 string struct if possible
func (w Wallet) TokenPermit(cd common.ContractPermitData) (string, error) {
	domainData := apitypes.TypedDataDomain{
		Name:              cd.Name,
		Version:           cd.Version,
		ChainId:           math.NewHexOrDecimal256(int64(cd.ChainId)),
		VerifyingContract: cd.FromToken,
	}

	// Order Message
	orderMessage := apitypes.TypedDataMessage{
		"owner":    cd.PublicAddress,
		"spender":  cd.Spender,
		"value":    cd.Amount,
		"nonce":    big.NewInt(cd.Nonce),
		"deadline": big.NewInt(cd.Deadline),
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

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return "", fmt.Errorf("error hashing typed data: %v", err)
	}
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return "", fmt.Errorf("error hashing domain separator: %v", err)
	}

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	challengeHash := crypto.Keccak256Hash(rawData)

	signature, err := crypto.Sign(challengeHash.Bytes(), w.privateKey)
	if err != nil {
		return "", fmt.Errorf("error signing challenge hash: %v", err)
	}
	signature[64] += 27 // Adjust the `v` value

	// Convert signature to hex string
	signatureHex := fmt.Sprintf("%x", signature)
	ownerNoPrefix := Remove0xPrefix(w.address.Hex())
	spenderNoPrefix := Remove0xPrefix(cd.Spender)

	return "0x" + padStringWithZeroes(ownerNoPrefix) +
		padStringWithZeroes(spenderNoPrefix) +
		padStringWithZeroes(fmt.Sprintf("%x", cd.Amount)) +
		padStringWithZeroes(fmt.Sprintf("%x", cd.Deadline)) +
		ConvertSignatureToVRSString(signatureHex), nil
}

func (w Wallet) GetContractDetailsForPermit(ctx context.Context, token gethCommon.Address, spender gethCommon.Address, deadline int64) (*common.ContractPermitData, error) {
	parsedABI, err := abi.JSON(strings.NewReader(abis.Erc20)) // Make a generic version of this ABI
	if err != nil {
		return nil, err
	}

	contractName, err := callAndUnpackContractMethod(ctx, token, parsedABI, &w.ethClient, "name")
	if err != nil {
		return nil, err
	}

	contractVersion, err := callAndUnpackContractMethod(ctx, token, parsedABI, &w.ethClient, "version")
	if err != nil {
		return nil, err
	}

	contractNonceStr, err := callAndUnpackContractMethod(ctx, token, parsedABI, &w.ethClient, "nonce", []gethCommon.Address{token})
	if err != nil {
		return nil, err
	}

	contractNonce, err := strconv.ParseInt(contractNonceStr, 10, 64)
	if err != nil {
		return nil, err
	}

	return &common.ContractPermitData{
		FromToken:     token.Hex(),
		PublicAddress: w.address.Hex(),
		Spender:       spender.Hex(),
		ChainId:       int(w.chainID.Int64()),
		Deadline:      deadline,
		Name:          contractName,
		Version:       contractVersion,
		Nonce:         contractNonce,
	}, nil
}

func callAndUnpackContractMethod(ctx context.Context, token gethCommon.Address, parsedABI abi.ABI, client *ethclient.Client, methodName string, methodArgs ...interface{}) (string, error) {
	data, err := parsedABI.Pack(methodName, methodArgs...)
	if err != nil {
		return "", err
	}

	msg := ethereum.CallMsg{
		To:   &token,
		Data: data,
	}

	result, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		return "", err
	}

	var returnValue string
	err = parsedABI.UnpackIntoInterface(&returnValue, methodName, result)
	if err != nil {
		return "", err
	}

	return returnValue, nil
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
