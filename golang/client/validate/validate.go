package validate

import (
	"fmt"
	"math/big"
	"regexp"

	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
)

func CheckEthereumAddressPointer(parameter interface{}, variableName string) error {
	value, ok := parameter.(*string)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string pointer", variableName, "EthereumAddressPointer")
	}
	if value == nil {
		return nil
	}
	return CheckEthereumAddress(*value, variableName)
}

func CheckEthereumAddress(parameter interface{}, variableName string) error {
	value, ok := parameter.(string)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, "EthereumAddress")
	}

	if value == "" {
		return NewParameterMissingError(variableName)
	}

	re := regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
	if !re.MatchString(value) {
		return NewParameterValidationError(variableName, "not a valid Ethereum address")
	}
	return nil
}

func CheckBigIntPointer(parameter interface{}, variableName string) error {
	value, ok := parameter.(*string)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string pointer", variableName, "BigIntPointer")
	}
	if value == nil {
		return nil
	}
	return CheckBigInt(*value, variableName)
}

var bigIntMax, _ = helpers.BigIntFromString("115792089237316195423570985008687907853269984665640564039457584007913129639935")
var bigIntZero = big.NewInt(0)

func CheckBigInt(parameter interface{}, variableName string) error {
	value, ok := parameter.(string)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, "BigInt")
	}

	if value == "" {
		return NewParameterMissingError(variableName)
	}

	parsedValue, err := helpers.BigIntFromString(value)
	if err != nil {
		return NewParameterValidationError(variableName, "not a valid value")
	}
	if parsedValue.Cmp(bigIntMax) > 0 {
		return NewParameterValidationError(variableName, "too big to fit in uint256")
	}
	if parsedValue.Cmp(bigIntZero) < 0 {
		return NewParameterValidationError(variableName, "must be a positive value")
	}
	return nil
}

func CheckChainId(parameter interface{}, variableName string) error {
	value, ok := parameter.(int)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be an int", variableName, "ChainId")
	}

	if value == 0 {
		return NewParameterMissingError(variableName)
	}

	if !helpers.Contains(value, chains.ValidChainIds) {
		return NewParameterValidationError(variableName, fmt.Sprintf("invalid chain id, valid chain ids are: %v", chains.ValidChainIds))
	}
	return nil
}

func CheckPrivateKey(parameter interface{}, variableName string) error {
	address, ok := parameter.(string)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, "PrivateKey")
	}

	if address == "" {
		return NewParameterMissingError(variableName)
	}

	re := regexp.MustCompile(`^[a-fA-F0-9]{64}$`)
	if !re.MatchString(address) {
		return NewParameterValidationError(variableName, "not a valid private key")
	}
	return nil
}

func CheckFloat32NonZeroPointer(parameter interface{}, variableName string) error {
	value, ok := parameter.(*float32)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32 pointer", variableName, "Float32NonZeroPointer")
	}
	if value == nil {
		return nil
	}
	return CheckFloat32NonZero(*value, variableName)
}

func CheckFloat32NonZero(parameter interface{}, variableName string) error {
	value, ok := parameter.(float32)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32", variableName, "Float32NonZero")
	}
	if value < 1 {
		return NewParameterValidationError(variableName, "must be explicitly set to a value greater than 0")
	}
	return nil
}

func CheckFloat32Pointer(parameter interface{}, variableName string) error {
	value, ok := parameter.(*float32)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32 pointer", variableName, "Float32Pointer")
	}
	if value == nil {
		return nil
	}
	return CheckFloat32(*value, variableName)
}

func CheckFloat32(parameter interface{}, variableName string) error {
	value, ok := parameter.(float32)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32", variableName, "Float32")
	}
	if value < 0 {
		return NewParameterValidationError(variableName, "must be non-negative")
	}
	return nil
}

func CheckApprovalType(parameter interface{}, variableName string) error {
	value, ok := parameter.(int)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be an int", variableName, "ApprovalType")
	}

	if !helpers.Contains(value, []int{0, 1, 2}) {
		return NewParameterValidationError(variableName, "invalid approval type")
	}
	return nil
}

func CheckSlippage(parameter interface{}, variableName string) error {
	value, ok := parameter.(float32)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32", variableName, "Slippage")
	}
	if value == 0 {
		return NewParameterMissingError(variableName)
	}
	if value < 0.01 || value > 50 {
		return NewParameterValidationError(variableName, fmt.Sprintf("invalid slippage value (%v) - only values 0.01-50 are allowed", value))
	}
	return nil
}

func CheckPagePointer(parameter interface{}, variableName string) error {
	value, ok := parameter.(*float32)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32 pointer", variableName, "PagePointer")
	}
	if value == nil {
		return nil
	}
	return CheckPage(*value, variableName)
}

func CheckPage(parameter interface{}, variableName string) error {
	value, ok := parameter.(float32)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32", variableName, "Page")
	}

	if value < 1 {
		return NewParameterValidationError(variableName, "must be greater than 0")
	}
	return nil
}

func CheckLimitPointer(parameter interface{}, variableName string) error {
	value, ok := parameter.(*float32)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32 pointer", variableName, "LimitPointer")
	}
	if value == nil {
		return nil
	}
	return CheckLimit(*value, variableName)
}

func CheckLimit(parameter interface{}, variableName string) error {
	value, ok := parameter.(float32)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32", variableName, "Limit")
	}

	if value < 1 {
		return NewParameterValidationError(variableName, "must be greater than 0")
	}
	return nil
}

func CheckStatusesIntsPointer(parameter interface{}, variableName string) error {
	value, ok := parameter.(*[]float32)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a *[]float32 pointer", variableName, "StatusesIntsPointer")
	}
	if value == nil {
		return nil
	}
	return CheckStatusesInts(*value, variableName)
}

func CheckStatusesInts(parameter interface{}, variableName string) error {
	value, ok := parameter.([]float32)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a []float32", variableName, "StatusesInts")
	}

	if helpers.HasDuplicates(value) {
		return NewParameterValidationError(variableName, "must not contain duplicates")
	}
	validStatuses := []float32{1, 2, 3}
	if !helpers.IsSubset(value, validStatuses) {
		return NewParameterValidationError(variableName, fmt.Sprintf("can only contain %v", validStatuses))
	}
	return nil
}

func CheckStatusesStringsPointer(parameter interface{}, variableName string) error {
	value, ok := parameter.(*[]string)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a *[]string pointer", variableName, "StatusesStringsPointer")
	}
	if value == nil {
		return nil
	}
	return CheckStatusesStrings(*value, variableName)
}

func CheckStatusesStrings(parameter interface{}, variableName string) error {
	value, ok := parameter.([]string)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a []string", variableName, "StatusesStrings")
	}

	if helpers.HasDuplicates(value) {
		return NewParameterValidationError(variableName, "must not contain duplicates")
	}
	validStatuses := []string{"1", "2", "3"}
	if !helpers.IsSubset(value, validStatuses) {
		return NewParameterValidationError(variableName, fmt.Sprintf("can only contain %v", validStatuses))
	}
	return nil
}

func CheckSortByPointer(parameter interface{}, variableName string) error {
	value, ok := parameter.(*string)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string pointer", variableName, "SortByPointer")
	}
	if value == nil {
		return nil
	}
	return CheckSortBy(*value, variableName)
}

func CheckSortBy(parameter interface{}, variableName string) error {
	value, ok := parameter.(string)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, "SortBy")
	}

	validSortBy := []string{"createDateTime", "takerRate", "makerRate", "makerAmount", "takerAmount"}
	if !helpers.Contains(value, validSortBy) {
		return NewParameterValidationError(variableName, fmt.Sprintf("can only contain %v", validSortBy))
	}
	return nil
}

func CheckOrderHash(parameter interface{}, variableName string) error {
	value, ok := parameter.(string)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, "OrderHash")
	}

	if value == "" {
		return NewParameterMissingError(variableName)
	}
	// TODO add criteria that captures valid order hash strings here
	return nil
}
