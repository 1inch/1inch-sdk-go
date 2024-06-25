package web3_provider

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
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

func getPermitData(permit PermitSingle, permit2Address common.Address, chainId *big.Int) (PermitSingleData, error) {
	err := validatePermitDetails(permit.Details)
	if err != nil {
		return PermitSingleData{}, err
	}

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
		"spender":           permit.Spender.Hex(),
		"sigDeadline":       permit.SigDeadline,
		"name":              domain.Name,
		"version":           domain.Version,
		"chainId":           domain.ChainId,
		"verifyingContract": domain.VerifyingContract,
	}

	return PermitSingleData{
		Domain: domain,
		Types:  PERMIT_TYPES,
		Values: values,
	}, nil
}

func validatePermitDetails(details PermitDetails) error {
	maxUint48 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 48), big.NewInt(1))
	maxUint160 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 160), big.NewInt(1))

	if details.Amount.Cmp(maxUint160) > 0 {
		return errors.New("AMOUNT_OUT_OF_RANGE")
	}
	if details.Expiration.Cmp(maxUint48) > 0 {
		return errors.New("EXPIRATION_OUT_OF_RANGE")
	}
	if details.Nonce.Cmp(maxUint48) > 0 {
		return errors.New("NONCE_OUT_OF_RANGE")
	}
	return nil
}

func hashPermitSingle(permit PermitSingle, permit2Address common.Address, chainId *big.Int) (string, error) {
	permitData, err := getPermitData(permit, permit2Address, chainId)
	if err != nil {
		return "", err
	}
	typedData := apitypes.TypedData{
		Types:       permitData.Types,
		PrimaryType: "PermitSingle",
		Domain:      permitData.Domain,
		Message:     permitData.Values,
	}

	challengeHash, _, err := apitypes.TypedDataAndHash(typedData)
	if err != nil {
		return "", fmt.Errorf("error using TypedDataAndHash: %v", err)
	}

	return fmt.Sprintf("%x", challengeHash), nil
}
