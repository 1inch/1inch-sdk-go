package web3_provider

import (
	"bytes"
	"math/big"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	PERMIT2_DOMAIN_NAME = "Permit2"
)

var (
	MaxUint256                 = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
	MaxSignatureTransferAmount = MaxUint256
	MaxUnorderedNonce          = MaxUint256
	MaxSigDeadline             = MaxUint256
)

type TypedDataDomain struct {
	Name              string
	ChainId           int
	VerifyingContract string
}

type TypedDataField struct {
	Name string
	Type string
}

type TokenPermissions struct {
	Token  string
	Amount *big.Int
}

type PermitTransferFrom struct {
	Permitted TokenPermissions
	Spender   string
	Nonce     *big.Int
	Deadline  *big.Int
}

type PermitBatchTransferFrom struct {
	Permitted []TokenPermissions
	Spender   string
	Nonce     *big.Int
	Deadline  *big.Int
}

type Witness struct {
	Witness         interface{}
	WitnessTypeName string
	WitnessType     map[string][]TypedDataField
}

type PermitTransferFromData struct {
	Domain TypedDataDomain
	Types  map[string][]TypedDataField
	Values PermitTransferFrom
}

var TOKEN_PERMISSIONS = []TypedDataField{
	{Name: "token", Type: "address"},
	{Name: "amount", Type: "uint256"},
}

var PERMIT_TRANSFER_FROM_TYPES = map[string][]TypedDataField{
	"PermitTransferFrom": {
		{Name: "permitted", Type: "TokenPermissions"},
		{Name: "spender", Type: "address"},
		{Name: "nonce", Type: "uint256"},
		{Name: "deadline", Type: "uint256"},
	},
	"TokenPermissions": TOKEN_PERMISSIONS,
}

func permit2Domain(permit2Address string, chainId int) TypedDataDomain {
	return TypedDataDomain{
		Name:              PERMIT2_DOMAIN_NAME,
		ChainId:           chainId,
		VerifyingContract: permit2Address,
	}
}

func permitTransferFromWithWitnessType(witness Witness) map[string][]TypedDataField {
	return map[string][]TypedDataField{
		"PermitWitnessTransferFrom": {
			{Name: "permitted", Type: "TokenPermissions"},
			{Name: "spender", Type: "address"},
			{Name: "nonce", Type: "uint256"},
			{Name: "deadline", Type: "uint256"},
			{Name: "witness", Type: witness.WitnessTypeName},
		},
		"TokenPermissions":      TOKEN_PERMISSIONS,
		witness.WitnessTypeName: witness.WitnessType[witness.WitnessTypeName],
	}
}

func validateTokenPermissions(permissions TokenPermissions) {
	if MaxSignatureTransferAmount.Cmp(permissions.Amount) < 0 {
		panic("AMOUNT_OUT_OF_RANGE")
	}
}

func getPermitTransferFromData(
	permit PermitTransferFrom,
	permit2Address string,
	chainId int,
	witness *Witness,
) PermitTransferFromData {
	if MaxSigDeadline.Cmp(permit.Deadline) < 0 {
		panic("SIG_DEADLINE_OUT_OF_RANGE")
	}
	if MaxUnorderedNonce.Cmp(permit.Nonce) < 0 {
		panic("NONCE_OUT_OF_RANGE")
	}

	validateTokenPermissions(permit.Permitted)

	domain := permit2Domain(permit2Address, chainId)
	var types map[string][]TypedDataField
	var values interface{}

	if witness != nil {
		types = permitTransferFromWithWitnessType(*witness)
		values = struct {
			PermitTransferFrom
			Witness interface{}
		}{permit, witness.Witness}
	} else {
		types = PERMIT_TRANSFER_FROM_TYPES
		values = permit
	}

	return PermitTransferFromData{
		Domain: domain,
		Types:  types,
		Values: values.(PermitTransferFrom),
	}
}

func hashPermitTransferFrom(
	permit PermitTransferFrom,
	permit2Address string,
	chainId int,
	witness *Witness,
) string {
	permitData := getPermitTransferFromData(permit, permit2Address, chainId, witness)
	domainHash := hashDomain(permitData.Domain)
	structHash := hashStruct(permitData.Types, permitData.Values)

	finalHash := crypto.Keccak256Hash(append([]byte{0x19, 0x01}, append(domainHash, structHash...)...)).Hex()
	return finalHash
}

func hashDomain(domain TypedDataDomain) []byte {
	var buffer bytes.Buffer

	buffer.Write(crypto.Keccak256([]byte("EIP712Domain(string name,uint256 chainId,address verifyingContract)")))

	nameHash := crypto.Keccak256([]byte(domain.Name))
	buffer.Write(nameHash)

	chainIdBytes := common.LeftPadBytes(new(big.Int).SetInt64(int64(domain.ChainId)).Bytes(), 32)
	buffer.Write(chainIdBytes)

	addressBytes := common.HexToAddress(domain.VerifyingContract).Bytes()
	buffer.Write(addressBytes)

	return crypto.Keccak256(buffer.Bytes())
}

func hashStruct(types map[string][]TypedDataField, values PermitTransferFrom) []byte {
	var buffer bytes.Buffer

	primaryType := ""
	for t := range types {
		primaryType = t
		break
	}

	buffer.Write(crypto.Keccak256([]byte(primaryType)))

	value := reflect.ValueOf(values)
	for _, field := range types[primaryType] {
		fieldName := field.Name
		fieldValue := value.FieldByName(strings.Title(fieldName))

		if field.Type == "address" {
			buffer.Write(common.HexToAddress(fieldValue.String()).Bytes())
		} else if field.Type == "uint256" {
			buffer.Write(common.LeftPadBytes(fieldValue.Interface().(*big.Int).Bytes(), 32))
		} else if field.Type == "TokenPermissions" {
			tokenPerm := fieldValue.Interface().(TokenPermissions)
			buffer.Write(common.HexToAddress(tokenPerm.Token).Bytes())
			buffer.Write(common.LeftPadBytes(tokenPerm.Amount.Bytes(), 32))
		}
	}

	return crypto.Keccak256(buffer.Bytes())
}
