package web3_provider

import (
	"bytes"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TypedDataField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type TypedDataDomain struct {
	Name              string `json:"name"`
	ChainId           int    `json:"chainId"`
	VerifyingContract string `json:"verifyingContract"`
}

type PermitDetails struct {
	Token      string `json:"token"`
	Amount     string `json:"amount"`
	Expiration string `json:"expiration"`
	Nonce      string `json:"nonce"`
}

type PermitSingle struct {
	Details     PermitDetails `json:"details"`
	Spender     string        `json:"spender"`
	SigDeadline string        `json:"sigDeadline"`
}

type TypedData struct {
	Domain apitypes.TypedDataDomain   `json:"domain"`
	Types  map[string][]apitypes.Type `json:"types"`
	Values apitypes.TypedDataMessage  `json:"values"`
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

func GetPermitData(permit PermitSingle, permit2Address common.Address, chainId int) (TypedData, error) {
	sigDeadline, ok := new(big.Int).SetString(permit.SigDeadline, 10)
	if !ok {
		return TypedData{}, fmt.Errorf("invalid sigDeadline")
	}
	if sigDeadline.Cmp(MaxSigDeadline) > 0 {
		return TypedData{}, fmt.Errorf("SIG_DEADLINE_OUT_OF_RANGE")
	}

	err := validatePermitDetails(permit.Details)
	if err != nil {
		return TypedData{}, err
	}
	values := map[string]interface{}{
		"details": map[string]interface{}{
			"token":      permit.Details.Token,
			"amount":     permit.Details.Amount,
			"expiration": permit.Details.Expiration,
			"nonce":      permit.Details.Nonce,
		},
		"spender":     permit.Spender,
		"sigDeadline": permit.SigDeadline,
	}

	return TypedData{
		Domain: apitypes.TypedDataDomain{
			Name:              "Permit2",
			ChainId:           math.NewHexOrDecimal256(int64(chainId)),
			VerifyingContract: permit2Address.Hex(),
		},
		Types:  PERMIT_TYPES,
		Values: values,
	}, nil
}

func validatePermitDetails(details PermitDetails) error {
	nonce, ok := new(big.Int).SetString(details.Nonce, 10)
	if !ok {
		return fmt.Errorf("invalid nonce")
	}
	if nonce.Cmp(MaxOrderedNonce) > 0 {
		return fmt.Errorf("NONCE_OUT_OF_RANGE")
	}

	amount, ok := new(big.Int).SetString(details.Amount, 10)
	if !ok {
		return fmt.Errorf("invalid amount")
	}
	if amount.Cmp(MaxAllowanceTransferAmount) > 0 {
		return fmt.Errorf("AMOUNT_OUT_OF_RANGE")
	}

	expiration, ok := new(big.Int).SetString(details.Expiration, 10)
	if !ok {
		return fmt.Errorf("invalid expiration")
	}
	if expiration.Cmp(MaxAllowanceExpiration) > 0 {
		return fmt.Errorf("EXPIRATION_OUT_OF_RANGE")
	}

	return nil
}

func hashPermitDetails(types map[string][]TypedDataField, details PermitDetails) []byte {
	var buffer bytes.Buffer

	buffer.Write(crypto.Keccak256([]byte("PermitDetails(address token,uint160 amount,uint48 expiration,uint48 nonce)")))

	value := reflect.ValueOf(details)
	for _, field := range types["PermitDetails"] {
		fieldName := field.Name
		fieldValue := value.FieldByName(cases.Title(language.English).String(fieldName))

		switch field.Type {
		case "address":
			buffer.Write(common.HexToAddress(fieldValue.String()).Bytes())
		case "uint160":
			intValue, ok := new(big.Int).SetString(fieldValue.String(), 10)
			if !ok {
				panic(fmt.Sprintf("Invalid uint160 value for field %s", fieldName))
			}
			buffer.Write(common.LeftPadBytes(intValue.Bytes(), 20))
		case "uint48":
			intValue, ok := new(big.Int).SetString(fieldValue.String(), 10)
			if !ok {
				panic(fmt.Sprintf("Invalid uint48 value for field %s", fieldName))
			}
			buffer.Write(common.LeftPadBytes(intValue.Bytes(), 6))
		}
	}

	return crypto.Keccak256(buffer.Bytes())
}

func hashStruct(types map[string][]TypedDataField, values PermitSingle) []byte {
	var buffer bytes.Buffer

	primaryType := "PermitSingle"
	buffer.Write(crypto.Keccak256([]byte("PermitSingle(PermitDetails details,address spender,uint256 sigDeadline)")))

	value := reflect.ValueOf(values)
	for _, field := range types[primaryType] {
		fieldName := field.Name
		// Manually handle the specific case for "sigDeadline" to "SigDeadline"
		var fieldValue reflect.Value
		if fieldName == "sigDeadline" {
			fieldValue = value.FieldByName("SigDeadline")
		} else {
			fieldValue = value.FieldByName(cases.Title(language.English).String(fieldName))
		}

		switch field.Type {
		case "address":
			buffer.Write(common.HexToAddress(fieldValue.String()).Bytes())
		case "uint160":
			intValue, ok := new(big.Int).SetString(fieldValue.String(), 10)
			if !ok {
				panic(fmt.Sprintf("Invalid uint160 value for field %s", fieldName))
			}
			buffer.Write(common.LeftPadBytes(intValue.Bytes(), 20))
		case "uint256":
			intValue, ok := new(big.Int).SetString(fieldValue.String(), 10)
			if !ok {
				panic(fmt.Sprintf("Invalid uint256 value for field %s", fieldName))
			}
			buffer.Write(common.LeftPadBytes(intValue.Bytes(), 32))
		case "uint48":
			intValue, ok := new(big.Int).SetString(fieldValue.String(), 10)
			if !ok {
				panic(fmt.Sprintf("Invalid uint48 value for field %s", fieldName))
			}
			buffer.Write(common.LeftPadBytes(intValue.Bytes(), 6))
		case "PermitDetails":
			buffer.Write(hashPermitDetails(types, fieldValue.Interface().(PermitDetails)))

		}
	}

	return crypto.Keccak256(buffer.Bytes())
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

func hashPermitData(permit PermitSingle, permit2Address common.Address, chainId int) (string, error) {
	permitData, err := GetPermitData(permit, permit2Address, chainId)
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

	return "0x" + common.Bytes2Hex(challengeHash), nil
}
