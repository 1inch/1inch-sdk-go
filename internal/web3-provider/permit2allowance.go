package web3_provider

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type AllowancePermitParams struct {
	Token          string `json:"token"`
	Amount         string `json:"amount"`
	Expiration     string `json:"expiration"`
	Spender        string `json:"spender"`
	SigDeadline    string `json:"sigDeadline"`
	Permit2Address string
}

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

func (w Wallet) GetAllowancePermitSingle(ctx context.Context, params AllowancePermitParams) (apitypes.TypedData, error) {
	callData, err := w.erc20ABI.Pack("nonce", w.address.Hex(), params.Token, params.Spender)
	if err != nil {
		return apitypes.TypedData{}, fmt.Errorf("failed to pack allowance call data: %v", err)
	}

	permitAddress := common.HexToAddress(params.Permit2Address)

	msg := ethereum.CallMsg{
		To:   &permitAddress,
		Data: callData,
	}

	result, err := w.ethClient.CallContract(ctx, msg, nil)
	if err != nil {
		return apitypes.TypedData{}, fmt.Errorf("failed to call contract: %v", err)
	}

	var nonce *big.Int
	err = w.erc20ABI.UnpackIntoInterface(&nonce, "nonce", result)
	if err != nil {
		return apitypes.TypedData{}, fmt.Errorf("failed to unpack result: %v", err)
	}

	// Construct the permit details
	d := AllowancePermitSingle{
		Details: AllowancePermitDetails{
			Token:      params.Token,
			Amount:     params.Amount,
			Expiration: params.Expiration,
			Nonce:      nonce.String(),
		},
		Spender:     params.Spender,
		SigDeadline: params.SigDeadline,
	}

	permit, err := GetTypedDataAllowancePermitSingle(d, permitAddress, int(w.chainId.Int64()))
	if err != nil {
		return apitypes.TypedData{}, fmt.Errorf("failed to generate permit: %v", err)
	}

	return permit, nil
}

func GetTypedDataAllowancePermitSingle(permit AllowancePermitSingle, permit2Address common.Address, chainId int) (apitypes.TypedData, error) {
	err := validatePermit(permit)
	if err != nil {
		return apitypes.TypedData{}, err
	}

	return apitypes.TypedData{
		Domain: apitypes.TypedDataDomain{
			Name:              "Permit2",
			ChainId:           math.NewHexOrDecimal256(int64(chainId)),
			VerifyingContract: permit2Address.Hex(),
		},
		Types: map[string][]apitypes.Type{
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
		},
		Message: apitypes.TypedDataMessage{
			"details": apitypes.TypedDataMessage{
				"token":      permit.Details.Token,
				"amount":     permit.Details.Amount,
				"expiration": permit.Details.Expiration,
				"nonce":      permit.Details.Nonce,
			},
			"spender":     permit.Spender,
			"sigDeadline": permit.SigDeadline,
		},
		PrimaryType: "PermitSingle",
	}, nil
}

func validatePermit(permit AllowancePermitSingle) error {
	var (
		MaxAllowanceTransferAmount = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 160), big.NewInt(1))
		MaxAllowanceExpiration     = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 48), big.NewInt(1))
		MaxOrderedNonce            = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 48), big.NewInt(1))
		MaxSigDeadline             = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
	)

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

func AllowancePermitSingleTypedDataHash(permit AllowancePermitSingle, permit2Address common.Address, chainId int) (string, error) {
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

func (w Wallet) SignPermit2AllowanceAndPackToContract(permit AllowancePermitSingle) (string, error) {
	challengeHash, err := AllowancePermitSingleTypedDataHash(permit, *w.address, int(w.chainId.Int64()))
	if err != nil {
		return "", err
	}
	challengeHashWithoutPrefix := challengeHash[2:]
	challengeHashWithoutPrefixRaw := common.Hex2Bytes(challengeHashWithoutPrefix)

	signature, err := crypto.Sign(challengeHashWithoutPrefixRaw, w.privateKey)
	if err != nil {
		return "", fmt.Errorf("error signing challenge hash: %v", err)
	}
	signature[64] += 27 // Adjust the `v` value

	// Convert createPermitSignature to hex string
	signatureHex := fmt.Sprintf("%x", signature)

	// Step 5: Encode the permit data with the signature
	permitCall, err := w.permit2ABI.Pack("permit", w.address, permit, convertSignatureToVRSString(signatureHex))
	if err != nil {
		return "", fmt.Errorf("failed to encode function data: %v", err)
	}
	return padStringWithZeroes(common.Bytes2Hex(permitCall)), nil
}
