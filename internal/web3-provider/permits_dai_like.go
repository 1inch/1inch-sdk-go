package web3_provider

import (
	"context"
	"fmt"
	"math/big"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"github.com/1inch/1inch-sdk-go/common"
	"github.com/1inch/1inch-sdk-go/internal/web3-provider/multicall"
)

// TokenPermit Will return an erc2612 string struct if possible
func (w Wallet) TokenPermitDaiLike(cd common.ContractPermitDataDaiLike) (string, error) {
	ownerNoPrefix := remove0xPrefix(w.address.Hex())
	spenderNoPrefix := remove0xPrefix(cd.Spender)
	signature, err := w.createPermitSignatureDaiLike(&cd)
	if err != nil {
		return "", err
	}
	return "0x" + padStringWithZeroes(ownerNoPrefix) +
		padStringWithZeroes(spenderNoPrefix) +
		padStringWithZeroes(fmt.Sprintf("%x", cd.Nonce)) +
		padStringWithZeroes(fmt.Sprintf("%x", cd.Expiry)) +
		padStringWithZeroes(fmt.Sprintf("%x", boolToInt(cd.Allowed))) +
		convertSignatureToVRSString(signature), nil
}

func (w Wallet) createPermitSignatureDaiLike(cd *common.ContractPermitDataDaiLike) (string, error) {
	// Dynamically build the EIP712Domain types
	eip712DomainTypes := []apitypes.Type{
		{Name: "name", Type: "string"},
	}
	if !cd.IsDomainWithoutVersion {
		eip712DomainTypes = append(eip712DomainTypes, apitypes.Type{Name: "version", Type: "string"})
	}
	if !cd.IsSaltInsteadOfChainId {
		eip712DomainTypes = append(eip712DomainTypes, apitypes.Type{Name: "chainId", Type: "uint256"})
	} else {
		eip712DomainTypes = append(eip712DomainTypes, apitypes.Type{Name: "salt", Type: "bytes32"})
	}
	eip712DomainTypes = append(eip712DomainTypes, apitypes.Type{Name: "verifyingContract", Type: "address"})

	// Permit model fields
	permitFields := []apitypes.Type{
		{Name: "holder", Type: "address"},
		{Name: "spender", Type: "address"},
		{Name: "nonce", Type: "uint256"},
		{Name: "expiry", Type: "uint256"},
		{Name: "allowed", Type: "bool"},
	}

	domainData := apitypes.TypedDataDomain{
		Name:              cd.Name,
		VerifyingContract: cd.FromToken,
	}

	if cd.IsSaltInsteadOfChainId {
		domainData.Salt = cd.Salt
	} else {
		domainData.ChainId = math.NewHexOrDecimal256(int64(cd.ChainId))
	}
	if !cd.IsDomainWithoutVersion {
		domainData.Version = cd.Version
	}

	orderMessage := apitypes.TypedDataMessage{
		"holder":  cd.Holder,
		"spender": cd.Spender,
		"allowed": cd.Allowed,
		"nonce":   big.NewInt(cd.Nonce),
		"expiry":  big.NewInt(cd.Expiry),
	}

	typedData := apitypes.TypedData{
		Types: map[string][]apitypes.Type{
			"EIP712Domain": eip712DomainTypes,
			"Permit":       permitFields,
		},
		PrimaryType: "Permit",
		Domain:      domainData,
		Message:     orderMessage,
	}

	challengeHash, _, err := apitypes.TypedDataAndHash(typedData)
	if err != nil {
		return "", fmt.Errorf("error using TypedDataAndHash: %v", err)
	}

	signature, err := crypto.Sign(challengeHash, w.privateKey)
	if err != nil {
		return "", fmt.Errorf("error signing challenge hash: %v", err)
	}
	signature[64] += 27 // Adjust the `v` value

	// Convert createPermitSignature to hex string
	signatureHex := fmt.Sprintf("%x", signature)
	return signatureHex, nil
}

func (w Wallet) GetContractDetailsForPermitDaiLike(ctx context.Context, token gethCommon.Address, spender gethCommon.Address, deadline int64) (*common.ContractPermitData, error) {
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

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
