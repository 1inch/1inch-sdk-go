package web3_provider

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type PermitDetails struct {
	Token      common.Address
	Amount     *big.Int
	Expiration *big.Int
	Nonce      *big.Int
}

type PermitSingle struct {
	Details     PermitDetails
	Spender     common.Address
	SigDeadline *big.Int
}

type PermitSingleData struct {
	Domain apitypes.TypedDataDomain
	Types  apitypes.Types
	Values apitypes.TypedDataMessage
}

var PERMIT_DETAILS = []apitypes.Type{
	{Name: "token", Type: "address"},
	{Name: "amount", Type: "uint160"},
	{Name: "expiration", Type: "uint48"},
	{Name: "nonce", Type: "uint48"},
}

var PERMIT_TYPES = apitypes.Types{
	"PermitSingle": {
		{Name: "details", Type: "PermitDetails"},
		{Name: "spender", Type: "address"},
		{Name: "sigDeadline", Type: "uint256"},
	},
	"PermitDetails": PERMIT_DETAILS,
	"EIP712Domain": []apitypes.Type{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	},
}

func getPermitData(permit PermitSingle, permit2Address common.Address, chainId *big.Int) PermitSingleData {
	validatePermitDetails(permit.Details)

	domain := apitypes.TypedDataDomain{
		Name:              "Permit2",
		Version:           "1",
		ChainId:           math.NewHexOrDecimal256(chainId.Int64()),
		VerifyingContract: permit2Address.Hex(),
	}

	values := map[string]interface{}{
		"details": map[string]interface{}{
			"token":      permit.Details.Token.Hex(),
			"amount":     permit.Details.Amount,
			"expiration": permit.Details.Expiration,
			"nonce":      permit.Details.Nonce,
		},
		"spender":     permit.Spender.Hex(),
		"sigDeadline": permit.SigDeadline,
	}

	return PermitSingleData{
		Domain: domain,
		Types:  PERMIT_TYPES,
		Values: values,
	}
}

func validatePermitDetails(details PermitDetails) {
	maxUint48 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 48), big.NewInt(1))
	maxUint160 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 160), big.NewInt(1))

	if details.Amount.Cmp(maxUint160) > 0 {
		log.Fatal("AMOUNT_OUT_OF_RANGE")
	}
	if details.Expiration.Cmp(maxUint48) > 0 {
		log.Fatal("EXPIRATION_OUT_OF_RANGE")
	}
	if details.Nonce.Cmp(maxUint48) > 0 {
		log.Fatal("NONCE_OUT_OF_RANGE")
	}
}

func hashPermitSingle(permit PermitSingle, permit2Address common.Address, chainId *big.Int) (string, error) {
	permitData := getPermitData(permit, permit2Address, chainId)

	typedData := apitypes.TypedData{
		Types:       permitData.Types,
		PrimaryType: "PermitSingle",
		Domain:      permitData.Domain,
		Message:     permitData.Values,
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return "", err
	}

	domainSeparator, err := typedData.HashStruct("EIP712Domain", map[string]interface{}{
		"name":              typedData.Domain.Name,
		"version":           typedData.Domain.Version,
		"chainId":           typedData.Domain.ChainId,
		"verifyingContract": typedData.Domain.VerifyingContract,
	})
	if err != nil {
		return "", err
	}

	digest := crypto.Keccak256Hash(
		[]byte("\x19\x01"),
		domainSeparator,
		typedDataHash,
	)

	return digest.Hex(), nil
}
