package validate

import (
	"errors"
	"fmt"
)

type ParameterType int

const (
	EthereumAddress ParameterType = iota
	EthereumAddressPointer
	BigInt
	BigIntPointer
	ChainID
	PrivateKey
	ApprovalType
	Slippage
	Page
	PagePointer
	Limit
	LimitPointer
	StatusesInts
	StatusesIntsPointer
	StatusesStrings
	StatusesStringsPointer
	SortBy
	SortByPointer
	OrderHash
)

func GetParameterTypeName(parameterType ParameterType) string {
	switch parameterType {
	case EthereumAddress:
		return "EthereumAddress"
	case EthereumAddressPointer:
		return "EthereumAddressPointer"
	case BigInt:
		return "BigInt"
	case BigIntPointer:
		return "BigIntPointer"
	case ChainID:
		return "ChainID"
	case PrivateKey:
		return "PrivateKey"
	case ApprovalType:
		return "ApprovalType"
	case Slippage:
		return "Slippage"
	case Page:
		return "Page"
	case PagePointer:
		return "PagePointer"
	case Limit:
		return "Limit"
	case LimitPointer:
		return "LimitPointer"
	case StatusesInts:
		return "StatusesInts"
	case StatusesIntsPointer:
		return "StatusesIntsPointer"
	case StatusesStrings:
		return "StatusesStrings"
	case StatusesStringsPointer:
		return "StatusesStringsPointer"
	case SortBy:
		return "SortBy"
	case SortByPointer:
		return "SortByPointer"
	case OrderHash:
		return "OrderHash"
	default:
		return "Unknown"
	}
}

func Parameter(parameter interface{}, variableName string, parameterType ParameterType, validationErrors []error) []error {
	var err error

	switch parameterType {
	case EthereumAddress:
		if value, ok := parameter.(string); ok {
			err = CheckEthereumAddress(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, GetParameterTypeName(EthereumAddress))
		}
	case EthereumAddressPointer:
		if value, ok := parameter.(*string); ok {
			err = CheckEthereumAddressPointer(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string pointer", variableName, GetParameterTypeName(EthereumAddressPointer))
		}
	case BigInt:
		if value, ok := parameter.(string); ok {
			err = CheckBigInt(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, GetParameterTypeName(BigInt))
		}
	case BigIntPointer:
		if value, ok := parameter.(*string); ok {
			err = CheckBigIntPointer(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string pointer", variableName, GetParameterTypeName(BigIntPointer))
		}
	case ChainID:
		if value, ok := parameter.(int); ok {
			err = CheckChainId(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be an int", variableName, GetParameterTypeName(ChainID))
		}
	case PrivateKey:
		if value, ok := parameter.(string); ok {
			err = CheckPrivateKey(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, GetParameterTypeName(PrivateKey))
		}
	case ApprovalType:
		if value, ok := parameter.(int); ok {
			err = CheckApprovalType(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be an int", variableName, GetParameterTypeName(ApprovalType))
		}
	case Slippage:
		if value, ok := parameter.(float32); ok {
			err = CheckSlippage(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32", variableName, GetParameterTypeName(Slippage))
		}
	case Page:
		if value, ok := parameter.(float32); ok {
			err = CheckPage(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32", variableName, GetParameterTypeName(Page))
		}
	case PagePointer:
		if value, ok := parameter.(*float32); ok {
			err = CheckPagePointer(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32 pointer", variableName, GetParameterTypeName(PagePointer))
		}
	case Limit:
		if value, ok := parameter.(float32); ok {
			err = CheckLimit(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32", variableName, GetParameterTypeName(Limit))
		}
	case LimitPointer:
		if value, ok := parameter.(*float32); ok {
			err = CheckLimitPointer(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32 pointer", variableName, GetParameterTypeName(LimitPointer))
		}
	case StatusesInts:
		if value, ok := parameter.([]float32); ok {
			err = CheckStatusesInts(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a []float32", variableName, GetParameterTypeName(StatusesInts))
		}
	case StatusesIntsPointer:
		if value, ok := parameter.(*[]float32); ok {
			err = CheckStatusesIntsPointer(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a *[]float32 pointer", variableName, GetParameterTypeName(StatusesIntsPointer))
		}
	case StatusesStrings:
		if value, ok := parameter.([]string); ok {
			err = CheckStatusesStrings(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a []string", variableName, GetParameterTypeName(StatusesStrings))
		}
	case StatusesStringsPointer:
		if value, ok := parameter.(*[]string); ok {
			err = CheckStatusesStringsPointer(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a *[]string pointer", variableName, GetParameterTypeName(StatusesStringsPointer))
		}
	case SortBy:
		if value, ok := parameter.(string); ok {
			err = CheckSortBy(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, GetParameterTypeName(SortBy))
		}
	case SortByPointer:
		if value, ok := parameter.(*string); ok {
			err = CheckSortByPointer(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string pointer", variableName, GetParameterTypeName(SortByPointer))
		}
	case OrderHash:
		if value, ok := parameter.(string); ok {
			err = CheckOrderHash(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, GetParameterTypeName(OrderHash))
		}
	default:
		err = errors.New("unknown parameter type")
	}

	if err != nil {
		validationErrors = append(validationErrors, err)
	}
	return validationErrors
}
