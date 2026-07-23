package orderbook

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	stdmath "math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	"github.com/1inch/1inch-sdk-go/common"
	"github.com/1inch/1inch-sdk-go/constants"
)

const permit2AllowanceABI = `[{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"allowance","outputs":[{"internalType":"uint160","name":"amount","type":"uint160"},{"internalType":"uint48","name":"expiration","type":"uint48"},{"internalType":"uint48","name":"nonce","type":"uint48"}],"stateMutability":"view","type":"function"}]`

var permit2AllowanceParsedABI, permit2AllowanceParsedABIErr = abi.JSON(strings.NewReader(permit2AllowanceABI))

var permit2CalldataArguments, permit2CalldataArgumentsErr = buildPermit2CalldataArguments()

func buildPermit2CalldataArguments() (abi.Arguments, error) {
	addressType, err := abi.NewType("address", "", nil)
	if err != nil {
		return nil, err
	}
	permitSingleType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "details", Type: "tuple", Components: []abi.ArgumentMarshaling{
			{Name: "token", Type: "address"},
			{Name: "amount", Type: "uint160"},
			{Name: "expiration", Type: "uint48"},
			{Name: "nonce", Type: "uint48"},
		}},
		{Name: "spender", Type: "address"},
		{Name: "sigDeadline", Type: "uint256"},
	})
	if err != nil {
		return nil, err
	}
	bytesType, err := abi.NewType("bytes", "", nil)
	if err != nil {
		return nil, err
	}
	return abi.Arguments{
		{Type: addressType},
		{Type: permitSingleType},
		{Type: bytesType},
	}, nil
}

// Permit2Allowance is the AllowanceTransfer state stored by the Permit2 contract
// for an (owner, token, spender) triple
type Permit2Allowance struct {
	Amount     *big.Int
	Expiration *big.Int
	Nonce      *big.Int
}

// Permit2PermitParams describes a Permit2 AllowanceTransfer PermitSingle message
type Permit2PermitParams struct {
	// Token is the ERC20 token whose allowance is being granted
	Token gethCommon.Address
	// Amount is the allowance granted to the spender (uint160)
	Amount *big.Int
	// Expiration is the allowance expiry timestamp (uint48)
	Expiration *big.Int
	// Nonce is the current Permit2 nonce for (owner, token, spender), see GetPermit2Allowance
	Nonce *big.Int
	// Spender receives the allowance, typically the 1inch Aggregation Router (constants.AggregationRouterV6)
	Spender gethCommon.Address
	// SigDeadline is the timestamp until which the signature itself is valid (uint256)
	SigDeadline *big.Int
}

// GetPermit2Allowance reads the AllowanceTransfer (amount, expiration, nonce) for
// (owner, token, spender) from the canonical Permit2 contract. The wallet must be
// RPC-connected (created with a node URL).
func GetPermit2Allowance(ctx context.Context, wallet common.Wallet, owner, token, spender gethCommon.Address) (*Permit2Allowance, error) {
	if permit2AllowanceParsedABIErr != nil {
		return nil, permit2AllowanceParsedABIErr
	}
	callData, err := permit2AllowanceParsedABI.Pack("allowance", owner, token, spender)
	if err != nil {
		return nil, err
	}
	result, err := wallet.Call(ctx, gethCommon.HexToAddress(constants.Permit2Address), callData)
	if err != nil {
		return nil, fmt.Errorf("failed to read Permit2 allowance: %w", err)
	}
	values, err := permit2AllowanceParsedABI.Unpack("allowance", result)
	if err != nil {
		return nil, err
	}
	return &Permit2Allowance{
		Amount:     values[0].(*big.Int),
		Expiration: values[1].(*big.Int),
		Nonce:      values[2].(*big.Int),
	}, nil
}

// BuildPermit2Calldata signs a Permit2 AllowanceTransfer PermitSingle message with the
// given wallet and returns the 352-byte permit calldata understood by the Limit Order
// Protocol: abi.encode(owner, PermitSingle, compactSignature). The wallet address is the
// permit owner and the wallet chain id is used for the EIP-712 domain. The signature is
// compacted to EIP-2098 form; a 65-byte signature would make the calldata unrecognizable
// to the protocol.
//
// The maker must have an ERC20 approval from Token to the Permit2 contract
// (constants.Permit2Address) for the permit to be executable. See
// BuildPermit2CalldataCompact for the smaller 96-byte encoding.
func BuildPermit2Calldata(wallet common.Wallet, params Permit2PermitParams) (string, error) {
	compactSig, err := signPermit2PermitSingle(wallet, params)
	if err != nil {
		return "", err
	}
	return encodePermit2Calldata(wallet.Address(), params, compactSig)
}

// BuildPermit2CalldataCompact signs the same PermitSingle message but returns the
// 96-byte compact permit encoding:
// amount uint160 | expiration uint32 | nonce uint32 | sigDeadline uint32 | r | vs.
// The protocol expands the uint32 timestamps as (value - 1) truncated to uint48, so a
// stored 0 means max uint48 (unlimited). Expiration and SigDeadline must therefore be
// either max uint48 or at most 2^32 - 2, and Nonce must fit in uint32.
//
// Caution: fusion order fills through the currently deployed Aggregation Router v6
// revert on compact permits, because the router's expansion of the 20-byte amount
// leaves uncleaned upper bits that fail Permit2's uint160 calldata validation
// (verified on a mainnet fork). Use BuildPermit2Calldata for orders until a router
// deployment accepts the compact form.
func BuildPermit2CalldataCompact(wallet common.Wallet, params Permit2PermitParams) (string, error) {
	if params.Amount == nil || params.Expiration == nil || params.Nonce == nil || params.SigDeadline == nil {
		return "", errors.New("amount, expiration, nonce, and sig deadline are required")
	}

	expirationStored, err := compactPermit2Timestamp(params.Expiration)
	if err != nil {
		return "", fmt.Errorf("expiration: %w", err)
	}
	sigDeadlineStored, err := compactPermit2Timestamp(params.SigDeadline)
	if err != nil {
		return "", fmt.Errorf("sig deadline: %w", err)
	}
	if !params.Nonce.IsUint64() || params.Nonce.Uint64() > stdmath.MaxUint32 {
		return "", errors.New("nonce must fit in uint32 for the compact permit encoding")
	}
	if params.Amount.BitLen() > 160 {
		return "", errors.New("amount must fit in uint160")
	}

	compactSig, err := signPermit2PermitSingle(wallet, params)
	if err != nil {
		return "", err
	}

	out := make([]byte, 0, 96)
	amountBytes := make([]byte, 20)
	params.Amount.FillBytes(amountBytes)
	out = append(out, amountBytes...)
	out = binary.BigEndian.AppendUint32(out, expirationStored)
	out = binary.BigEndian.AppendUint32(out, uint32(params.Nonce.Uint64()))
	out = binary.BigEndian.AppendUint32(out, sigDeadlineStored)
	out = append(out, compactSig...)

	return fmt.Sprintf("0x%x", out), nil
}

// compactPermit2Timestamp encodes a uint48 timestamp for the compact permit form,
// which the protocol decodes as (stored - 1) truncated to uint48
func compactPermit2Timestamp(value *big.Int) (uint32, error) {
	if value.Cmp(constants.Uint48Max) == 0 {
		return 0, nil
	}
	if value.Sign() < 0 || value.BitLen() > 32 || value.Uint64()+1 > stdmath.MaxUint32 {
		return 0, errors.New("value must be max uint48 or at most 2^32 - 2")
	}
	return uint32(value.Uint64()) + 1, nil
}

// signPermit2PermitSingle signs the EIP-712 PermitSingle and returns the 64-byte
// EIP-2098 compact signature (r, vs)
func signPermit2PermitSingle(wallet common.Wallet, params Permit2PermitParams) ([]byte, error) {
	if params.Amount == nil || params.Expiration == nil || params.Nonce == nil || params.SigDeadline == nil {
		return nil, errors.New("amount, expiration, nonce, and sig deadline are required")
	}

	typedData := apitypes.TypedData{
		Types: map[string][]apitypes.Type{
			"EIP712Domain": {
				{Name: "name", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"PermitDetails": {
				{Name: "token", Type: "address"},
				{Name: "amount", Type: "uint160"},
				{Name: "expiration", Type: "uint48"},
				{Name: "nonce", Type: "uint48"},
			},
			"PermitSingle": {
				{Name: "details", Type: "PermitDetails"},
				{Name: "spender", Type: "address"},
				{Name: "sigDeadline", Type: "uint256"},
			},
		},
		PrimaryType: "PermitSingle",
		Domain: apitypes.TypedDataDomain{
			Name:              "Permit2",
			ChainId:           math.NewHexOrDecimal256(wallet.ChainId()),
			VerifyingContract: constants.Permit2Address,
		},
		Message: apitypes.TypedDataMessage{
			"details": map[string]any{
				"token":      params.Token.Hex(),
				"amount":     params.Amount.String(),
				"expiration": params.Expiration.String(),
				"nonce":      params.Nonce.String(),
			},
			"spender":     params.Spender.Hex(),
			"sigDeadline": params.SigDeadline.String(),
		},
	}

	digest, _, err := apitypes.TypedDataAndHash(typedData)
	if err != nil {
		return nil, fmt.Errorf("failed to hash PermitSingle typed data: %w", err)
	}

	signature, err := wallet.SignBytes(digest)
	if err != nil {
		return nil, fmt.Errorf("failed to sign PermitSingle digest: %w", err)
	}
	signature[64] += 27

	compact, err := CompressSignature(fmt.Sprintf("%x", signature))
	if err != nil {
		return nil, fmt.Errorf("failed to compress signature: %w", err)
	}
	return append(append([]byte{}, compact.R...), compact.VS...), nil
}

// encodePermit2Calldata encodes (address owner, PermitSingle permit, bytes signature)
func encodePermit2Calldata(owner gethCommon.Address, params Permit2PermitParams, signature []byte) (string, error) {
	if permit2CalldataArgumentsErr != nil {
		return "", permit2CalldataArgumentsErr
	}

	permitValue := struct {
		Details struct {
			Token      gethCommon.Address
			Amount     *big.Int
			Expiration *big.Int
			Nonce      *big.Int
		}
		Spender     gethCommon.Address
		SigDeadline *big.Int
	}{
		Details: struct {
			Token      gethCommon.Address
			Amount     *big.Int
			Expiration *big.Int
			Nonce      *big.Int
		}{
			Token:      params.Token,
			Amount:     params.Amount,
			Expiration: params.Expiration,
			Nonce:      params.Nonce,
		},
		Spender:     params.Spender,
		SigDeadline: params.SigDeadline,
	}

	packed, err := permit2CalldataArguments.Pack(owner, permitValue, signature)
	if err != nil {
		return "", fmt.Errorf("failed to encode permit2 calldata: %w", err)
	}
	return fmt.Sprintf("0x%x", packed), nil
}
