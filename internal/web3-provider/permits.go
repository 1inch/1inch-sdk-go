package web3_provider

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"github.com/1inch/1inch-sdk-go/common"
	"github.com/1inch/1inch-sdk-go/internal/web3-provider/multicall"
)

// TokenPermit Will return an erc2612 string struct if possible
func (w Wallet) TokenPermit(cd common.ContractPermitData) (string, error) {
	ownerNoPrefix := remove0xPrefix(w.address.Hex())
	spenderNoPrefix := remove0xPrefix(cd.Spender)
	signature, err := w.createPermitSignature(&cd)
	if err != nil {
		return "", err
	}

	a := new(big.Int)
	a, ok := a.SetString(cd.Amount, 10)
	if !ok {
		return "", fmt.Errorf("bad amount")
	}

	return "0x" + padStringWithZeroes(ownerNoPrefix) +
		padStringWithZeroes(spenderNoPrefix) +
		padStringWithZeroes(fmt.Sprintf("%x", a)) +
		padStringWithZeroes(fmt.Sprintf("%x", cd.Deadline)) +
		convertSignatureToVRSString(signature), nil
}

func (w Wallet) createPermitSignature(cd *common.ContractPermitData) (string, error) {
	domainData := apitypes.TypedDataDomain{
		Name:              cd.Name,
		Version:           cd.Version,
		ChainId:           math.NewHexOrDecimal256(int64(cd.ChainId)),
		VerifyingContract: cd.FromToken,
	}

	amount, ok := new(big.Int).SetString(cd.Amount, 10)
	if !ok {
		return "", fmt.Errorf("failed to convert string (%v) to big.Int", cd.Amount)
	}

	// Order Message
	orderMessage := apitypes.TypedDataMessage{
		"owner":    cd.PublicAddress,
		"spender":  cd.Spender,
		"value":    amount,
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

	// Convert createPermitSignature to hex string
	signatureHex := fmt.Sprintf("%x", signature)
	return signatureHex, nil
}

func (w Wallet) GetContractDetailsForPermit(ctx context.Context, token gethCommon.Address, spender gethCommon.Address, deadline int64) (*common.ContractPermitData, error) {
	contractNameData, err := w.erc20ABI.Pack("name")
	if err != nil {
		return nil, err
	}

	contractVersionData, err := w.erc20ABI.Pack("version")
	if err != nil {
		return nil, err
	}

	contractNonceData, err := w.erc20ABI.Pack("nonce", []gethCommon.Address{token})
	if err != nil {
		return nil, err
	}

	callDataArray := []multicall.CallData{
		multicall.BuildCallData(token, contractNameData, 0),
		multicall.BuildCallData(token, contractVersionData, 0),
		multicall.BuildCallData(token, contractNonceData, 0),
	}

	mResult, err := w.multicall.Execute(ctx, callDataArray)
	if err != nil {
		return nil, err
	}

	var contractName string
	err = w.erc20ABI.UnpackIntoInterface(&contractName, "name", mResult[0])
	if err != nil {
		return nil, err
	}

	var contractVersion string
	err = w.erc20ABI.UnpackIntoInterface(&contractVersion, "version", mResult[1])
	if err != nil {
		return nil, err
	}

	var contractNonce int64
	err = w.erc20ABI.UnpackIntoInterface(&contractNonce, "nonce", mResult[2])
	if err != nil {
		return nil, err
	}

	return &common.ContractPermitData{
		FromToken:     token.Hex(),
		PublicAddress: w.address.Hex(),
		Spender:       spender.Hex(),
		ChainId:       int(w.ChainId()),
		Deadline:      deadline,
		Name:          contractName,
		Version:       contractVersion,
		Nonce:         contractNonce,
	}, nil
}

func padStringWithZeroes(s string) string {
	if len(s) >= 64 {
		return s
	}
	return strings.Repeat("0", 64-len(s)) + s
}

func remove0xPrefix(s string) string {
	if strings.HasPrefix(s, "0x") {
		return s[2:]
	}
	return s
}

// ConvertSignatureToVRSString converts a createPermitSignature from rsv to padded vrs format
func convertSignatureToVRSString(signature string) string {
	// explicit breakdown
	//r := createPermitSignature[:66]
	//s := createPermitSignature[66:128]
	//v := createPermitSignature[128:]
	return padStringWithZeroes(signature[128:]) + signature[:128]
}
