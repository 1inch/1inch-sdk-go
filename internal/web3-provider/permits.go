package web3_provider

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
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

	return "0x" + padStringWithZeroes(ownerNoPrefix) +
		padStringWithZeroes(spenderNoPrefix) +
		padStringWithZeroes(fmt.Sprintf("%x", cd.Amount)) +
		padStringWithZeroes(fmt.Sprintf("%x", cd.Deadline)) +
		convertSignatureToVRSString(signature), nil
}

func ConvertHexToBytes32(hexStr string) ([32]byte, error) {
	var b32 [32]byte
	bytes := gethCommon.FromHex(hexStr)
	if len(bytes) != 32 {
		return b32, fmt.Errorf("hex must be exactly 32 bytes long")
	}
	copy(b32[:], bytes)
	return b32, nil
}

func (w Wallet) signatureWithDomainSeparator(cd *common.ContractPermitData) (string, error) {
	owner := gethCommon.HexToAddress(cd.PublicAddress)
	spender := gethCommon.HexToAddress(cd.Spender)
	value := cd.Amount
	nonce := big.NewInt(cd.Nonce)
	deadline := big.NewInt(cd.Deadline)

	// Convert hex strings to fixed-size byte arrays
	var domainSeparator [32]byte
	var typeHash [32]byte
	copy(domainSeparator[:], gethCommon.FromHex(cd.DomainSeparator))
	copy(typeHash[:], gethCommon.FromHex(cd.PermitTypeHash))

	// Create ABI type instances without incorrect nil usage
	bytes32Type, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create type for bytes32: %v", err)
	}
	addressType, err := abi.NewType("address", "", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create type for address: %v", err)
	}
	uint256Type, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create type for uint256: %v", err)
	}

	// Pack the data with ABI
	arguments := abi.Arguments{
		{Type: bytes32Type},
		{Type: bytes32Type},
		{Type: addressType},
		{Type: addressType},
		{Type: uint256Type},
		{Type: uint256Type},
		{Type: uint256Type},
	}

	encodedData, err := arguments.Pack(
		domainSeparator,
		typeHash,
		owner,
		spender,
		value,
		nonce,
		deadline,
	)
	if err != nil {
		return "", fmt.Errorf("failed to pack ABI arguments: %v", err)
	}

	// Hash the packed data
	hash := crypto.Keccak256Hash(encodedData)

	// Sign the hash
	signature, err := crypto.Sign(hash.Bytes(), w.privateKey)
	if err != nil {
		return "", fmt.Errorf("error signing the permit data: %v", err)
	}

	// Adjust the recovery byte (v value)
	if signature[64] < 27 {
		signature[64] += 27
	}

	signatureHex := fmt.Sprintf("%x", signature)
	return signatureHex, nil
}

func (w Wallet) createPermitSignature(cd *common.ContractPermitData) (string, error) {
	if cd.PermitTypeHash != "" && cd.DomainSeparator != "" {
		return w.signatureWithDomainSeparator(cd)
	}

	eip712DomainTypes := []apitypes.Type{
		{Name: "name", Type: "string"},
	}
	if !cd.IsDomainWithoutVersion {
		eip712DomainTypes = append(eip712DomainTypes, apitypes.Type{Name: "version", Type: "string"})
	}
	//if !cd.IsSaltInsteadOfChainId {
	eip712DomainTypes = append(eip712DomainTypes, apitypes.Type{Name: "chainId", Type: "uint256"})
	//}
	eip712DomainTypes = append(eip712DomainTypes, apitypes.Type{Name: "verifyingContract", Type: "address"})

	//if cd.IsSaltInsteadOfChainId {
	//	eip712DomainTypes = append(eip712DomainTypes, apitypes.Type{Name: "salt", Type: "bytes32"})
	//}
	// Permit model fields
	permitFields := []apitypes.Type{
		{Name: "owner", Type: "address"},
		{Name: "spender", Type: "address"},
		{Name: "value", Type: "uint256"},
		{Name: "nonce", Type: "uint256"},
		{Name: "deadline", Type: "uint256"},
	}

	domainData := apitypes.TypedDataDomain{
		Name:              cd.Name,
		VerifyingContract: cd.FromToken,
	}

	//if cd.IsSaltInsteadOfChainId {
	//	domainData.Salt = cd.Salt
	//} else {
	domainData.ChainId = math.NewHexOrDecimal256(int64(cd.ChainId))
	//}
	if !cd.IsDomainWithoutVersion {
		domainData.Version = cd.Version
	}

	//// Order Message
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

func (w Wallet) GetContractDetailsForPermit(ctx context.Context, token gethCommon.Address, spender gethCommon.Address, amount *big.Int, deadline int64) (*common.ContractPermitData, error) {
	contractNameData, err := w.erc20ABI.Pack("name")
	if err != nil {
		return nil, err
	}

	contractVersionData, err := w.erc20ABI.Pack("version")
	if err != nil {
		return nil, err
	}

	contractNonceData, err := w.erc20ABI.Pack("nonces", w.Address())
	if err != nil {
		return nil, err
	}

	contractPermitTypeHashData, err := w.erc20ABI.Pack("PERMIT_TYPEHASH")
	if err != nil {
		return nil, err
	}

	contractDomainSeparatorData, err := w.erc20ABI.Pack("DOMAIN_SEPARATOR")
	if err != nil {
		return nil, err
	}

	callDataArray := []multicall.CallData{
		multicall.BuildCallData(token, contractNameData, 0),
		multicall.BuildCallData(token, contractVersionData, 0),
		multicall.BuildCallData(token, contractNonceData, 0),
		multicall.BuildCallData(token, contractPermitTypeHashData, 0),
		multicall.BuildCallData(token, contractDomainSeparatorData, 0),
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
	if len(mResult[1]) == 0 {
		contractVersion = "1"
	} else {
		err = w.erc20ABI.UnpackIntoInterface(&contractVersion, "version", mResult[1])
		if err != nil {
			return nil, err
		}
	}

	var contractNonce *big.Int
	err = w.erc20ABI.UnpackIntoInterface(&contractNonce, "nonces", mResult[2])
	if err != nil {
		return nil, err
	}

	var permitTypeHashHex string
	var permitTypeHash [32]byte
	err = w.erc20ABI.UnpackIntoInterface(&permitTypeHash, "PERMIT_TYPEHASH", mResult[3])
	if err != nil {
		permitTypeHashHex = ""
	} else {
		permitTypeHashHex = "0x" + hex.EncodeToString(permitTypeHash[:])
	}

	var domainSeparatorHex string
	var domainSeparator [32]byte
	err = w.erc20ABI.UnpackIntoInterface(&domainSeparator, "DOMAIN_SEPARATOR", mResult[4])
	if err != nil {
		domainSeparatorHex = ""
	} else {
		domainSeparatorHex = "0x" + hex.EncodeToString(domainSeparator[:])
	}

	return &common.ContractPermitData{
		FromToken:       token.Hex(),
		PublicAddress:   w.address.Hex(),
		Spender:         spender.Hex(),
		ChainId:         int(w.ChainId()),
		Deadline:        deadline,
		Name:            contractName,
		Version:         contractVersion,
		Nonce:           contractNonce.Int64(),
		Amount:          amount,
		PermitTypeHash:  permitTypeHashHex,
		DomainSeparator: domainSeparatorHex,
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
