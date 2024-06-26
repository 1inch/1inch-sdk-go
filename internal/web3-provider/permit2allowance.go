package web3_provider

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type AllowancePermitDetails struct {
	Token      string `json:"token"`
	Amount     string `json:"amount"`
	Expiration string `json:"expiration"`
	Nonce      string `json:"nonce"`
}

type AllowancePermitSingle struct {
	Details     AllowancePermitDetails `json:"details"`
	Spender     string                 `json:"spender"`
	SigDeadline string                 `json:"sigDeadline"`
}

var PERMIT_TYPES = map[string][]apitypes.Type{
	"EIP712Domain": {
		{Name: "name", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	},
	"PermitSingle": {
		{Name: "details", Type: "PermitDetails"},
		{Name: "spender", Type: "address"},
		{Name: "sigDeadline", Type: "uint256"},
	},
	"PermitDetails": {
		{Name: "token", Type: "address"},
		{Name: "amount", Type: "uint160"},
		{Name: "expiration", Type: "uint48"},
		{Name: "nonce", Type: "uint48"},
	},
}

var (
	MaxAllowanceTransferAmount = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 160), big.NewInt(1))
	MaxAllowanceExpiration     = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 48), big.NewInt(1))
	MaxOrderedNonce            = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 48), big.NewInt(1))
	MaxSigDeadline             = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
)

func GetTypedDataAllowancePermitSingle(permit AllowancePermitSingle, permit2Address common.Address, chainId int) (apitypes.TypedData, error) {
	err := validatePermit(permit)
	if err != nil {
		return apitypes.TypedData{}, err
	}
	values := apitypes.TypedDataMessage{
		"details": apitypes.TypedDataMessage{
			"token":      permit.Details.Token,
			"amount":     permit.Details.Amount,
			"expiration": permit.Details.Expiration,
			"nonce":      permit.Details.Nonce,
		},
		"spender":     permit.Spender,
		"sigDeadline": permit.SigDeadline,
	}

	return apitypes.TypedData{
		Domain: apitypes.TypedDataDomain{
			Name:              "Permit2",
			ChainId:           math.NewHexOrDecimal256(int64(chainId)),
			VerifyingContract: permit2Address.Hex(),
		},
		Types:       PERMIT_TYPES,
		Message:     values,
		PrimaryType: "PermitSingle",
	}, nil
}

func validatePermit(permit AllowancePermitSingle) error {
	nonce, ok := new(big.Int).SetString(permit.Details.Nonce, 10)
	if !ok {
		return fmt.Errorf("invalid nonce")
	}
	if nonce.Cmp(MaxOrderedNonce) > 0 {
		return fmt.Errorf("NONCE_OUT_OF_RANGE")
	}

	amount, ok := new(big.Int).SetString(permit.Details.Amount, 10)
	if !ok {
		return fmt.Errorf("invalid amount")
	}
	if amount.Cmp(MaxAllowanceTransferAmount) > 0 {
		return fmt.Errorf("AMOUNT_OUT_OF_RANGE")
	}

	expiration, ok := new(big.Int).SetString(permit.Details.Expiration, 10)
	if !ok {
		return fmt.Errorf("invalid expiration")
	}
	if expiration.Cmp(MaxAllowanceExpiration) > 0 {
		return fmt.Errorf("EXPIRATION_OUT_OF_RANGE")
	}

	sigDeadline, ok := new(big.Int).SetString(permit.SigDeadline, 10)
	if !ok {
		return fmt.Errorf("invalid sigDeadline")
	}
	if sigDeadline.Cmp(MaxSigDeadline) > 0 {
		return fmt.Errorf("SIG_DEADLINE_OUT_OF_RANGE")
	}

	return nil
}

func hashPermitData(permit AllowancePermitSingle, permit2Address common.Address, chainId int) (string, error) {
	typedData, err := GetTypedDataAllowancePermitSingle(permit, permit2Address, chainId)
	if err != nil {
		return "", err
	}

	challengeHash, _, err := apitypes.TypedDataAndHash(typedData)
	if err != nil {
		return "", fmt.Errorf("error using TypedDataAndHash: %v", err)
	}

	return "0x" + common.Bytes2Hex(challengeHash), nil
}
