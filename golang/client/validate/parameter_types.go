package validate

import (
	"errors"
	"fmt"
)

type ParameterType int

// Generic parameter types
const (
	EthereumAddress ParameterType = iota
	EthereumAddressPointer
	BigInt
	BigIntPointer
	ChainID
	PrivateKey
)

// Swap parameter types
const (
	ApprovalType ParameterType = iota + 50
	Slippage
)

// Orderbook parameter types
const (
	Page ParameterType = iota + 100
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

func Parameter(parameter interface{}, variableName string, parameterType ParameterType, validationErrors []error) []error {
	var err error

	switch parameterType {
	// Generic parameter types
	case EthereumAddress:
		if value, ok := parameter.(string); ok {
			err = CheckEthereumAddress(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, "EthereumAddress")
		}
	case EthereumAddressPointer:
		if value, ok := parameter.(*string); ok {
			err = CheckEthereumAddressPointer(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string pointer", variableName, "EthereumAddressPointer")
		}
	case BigInt:
		if value, ok := parameter.(string); ok {
			err = CheckBigInt(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, "BigInt")
		}
	case BigIntPointer:
		if value, ok := parameter.(*string); ok {
			err = CheckBigIntPointer(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string pointer", variableName, "BigIntPointer")
		}
	case ChainID:
		if value, ok := parameter.(int); ok {
			err = CheckChainId(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be an int", variableName, "ChainID")
		}
	case PrivateKey:
		if value, ok := parameter.(string); ok {
			err = CheckPrivateKey(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, "PrivateKey")
		}

	// Swap parameter types
	case ApprovalType:
		if value, ok := parameter.(int); ok {
			err = CheckApprovalType(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be an int", variableName, "ApprovalType")
		}
	case Slippage:
		if value, ok := parameter.(float32); ok {
			err = CheckSlippage(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32", variableName, "Slippage")
		}

	// Orderbook parameter types
	case Page:
		if value, ok := parameter.(float32); ok {
			err = CheckPage(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32", variableName, "Page")
		}
	case PagePointer:
		if value, ok := parameter.(*float32); ok {
			err = CheckPagePointer(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32 pointer", variableName, "PagePointer")
		}
	case Limit:
		if value, ok := parameter.(float32); ok {
			err = CheckLimit(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32", variableName, "Limit")
		}
	case LimitPointer:
		if value, ok := parameter.(*float32); ok {
			err = CheckLimitPointer(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32 pointer", variableName, "LimitPointer")
		}
	case StatusesInts:
		if value, ok := parameter.([]float32); ok {
			err = CheckStatusesInts(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a []float32", variableName, "StatusesInts")
		}
	case StatusesIntsPointer:
		if value, ok := parameter.(*[]float32); ok {
			err = CheckStatusesIntsPointer(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a *[]float32 pointer", variableName, "StatusesIntsPointer")
		}
	case StatusesStrings:
		if value, ok := parameter.([]string); ok {
			err = CheckStatusesStrings(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a []string", variableName, "StatusesStrings")
		}
	case StatusesStringsPointer:
		if value, ok := parameter.(*[]string); ok {
			err = CheckStatusesStringsPointer(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a *[]string pointer", variableName, "StatusesStringsPointer")
		}
	case SortBy:
		if value, ok := parameter.(string); ok {
			err = CheckSortBy(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, "SortBy")
		}
	case SortByPointer:
		if value, ok := parameter.(*string); ok {
			err = CheckSortByPointer(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string pointer", variableName, "SortByPointer")
		}
	case OrderHash:
		if value, ok := parameter.(string); ok {
			err = CheckOrderHash(value, variableName)
		} else {
			err = fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, "OrderHash")
		}
	default:
		err = errors.New("unknown parameter type")
	}

	if err != nil {
		validationErrors = append(validationErrors, err)
	}
	return validationErrors
}
