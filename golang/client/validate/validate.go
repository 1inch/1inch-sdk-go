package validate

import (
	"fmt"
	"math/big"
	"regexp"

	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
)

func CheckEthereumAddressRequired(parameter interface{}, variableName string) error {
	value, ok := parameter.(string)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, "EthereumAddress")
	}

	if value == "" {
		return NewParameterMissingError(variableName)
	}

	return CheckEthereumAddress(value, variableName)
}
func CheckEthereumAddress(parameter interface{}, variableName string) error {
	value, ok := parameter.(string)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, "EthereumAddress")
	}
	if value == "" {
		return nil
	}

	re := regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
	if !re.MatchString(value) {
		return NewParameterValidationError(variableName, "not a valid Ethereum address")
	}
	return nil
}

var bigIntMax, _ = helpers.BigIntFromString("115792089237316195423570985008687907853269984665640564039457584007913129639935")
var bigIntZero = big.NewInt(0)

func CheckBigIntRequired(parameter interface{}, variableName string) error {
	value, ok := parameter.(string)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, "BigInt")
	}

	if value == "" {
		return NewParameterMissingError(variableName)
	}

	return CheckBigInt(value, variableName)
}
func CheckBigInt(parameter interface{}, variableName string) error {
	value, ok := parameter.(string)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, "BigInt")
	}

	if value == "" {
		return nil
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

func CheckChainIdRequired(parameter interface{}, variableName string) error {
	value, ok := parameter.(int)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be an int", variableName, "ChainId")
	}

	if value == 0 {
		return NewParameterMissingError(variableName)
	}

	return CheckChainId(value, variableName)
}

func CheckChainId(parameter interface{}, variableName string) error {
	value, ok := parameter.(int)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be an int", variableName, "ChainId")
	}
	if value == 0 {
		return nil
	}

	if !helpers.Contains(value, chains.ValidChainIds) {
		return NewParameterValidationError(variableName, fmt.Sprintf("is invalid, valid chain ids are: %v", chains.ValidChainIds))
	}
	return nil
}

func CheckPrivateKeyRequired(parameter interface{}, variableName string) error {
	address, ok := parameter.(string)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, "PrivateKey")
	}

	if address == "" {
		return NewParameterMissingError(variableName)
	}

	return CheckPrivateKey(address, variableName)
}
func CheckPrivateKey(parameter interface{}, variableName string) error {
	address, ok := parameter.(string)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, "PrivateKey")
	}

	if address == "" {
		return nil
	}

	re := regexp.MustCompile(`^[a-fA-F0-9]{64}$`)
	if !re.MatchString(address) {
		return NewParameterValidationError(variableName, "not a valid private key")
	}
	return nil
}

func CheckApprovalType(parameter interface{}, variableName string) error {
	value, ok := parameter.(int)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be an int", variableName, "ApprovalType")
	}

	if value == 0 {
		return nil
	}

	if !helpers.Contains(value, []int{0, 1, 2}) {
		return NewParameterValidationError(variableName, "invalid approval type")
	}
	return nil
}

func CheckSlippageRequired(parameter interface{}, variableName string) error {
	value, ok := parameter.(float32)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32", variableName, "Slippage")
	}
	if value == 0 {
		return NewParameterMissingError(variableName)
	}
	return CheckSlippage(value, variableName)
}
func CheckSlippage(parameter interface{}, variableName string) error {
	value, ok := parameter.(float32)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32", variableName, "Slippage")
	}
	if value == 0 {
		return nil
	}
	if value < 0.01 || value > 50 {
		return NewParameterValidationError(variableName, fmt.Sprintf("invalid slippage value (%v) - only values 0.01-50 are allowed", value))
	}
	return nil
}

func CheckPage(parameter interface{}, variableName string) error {
	value, ok := parameter.(float32)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32", variableName, "Page")
	}
	if value == 0 {
		return nil
	}

	if value < 1 {
		return NewParameterValidationError(variableName, "must be greater than 0")
	}
	return nil
}

func CheckLimit(parameter interface{}, variableName string) error {
	value, ok := parameter.(float32)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a float32", variableName, "Limit")
	}
	if value == 0 {
		return nil
	}

	if value < 1 {
		return NewParameterValidationError(variableName, "must be greater than 0")
	}
	return nil
}

func CheckStatusesInts(parameter interface{}, variableName string) error {
	value, ok := parameter.([]float32)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a []float32", variableName, "StatusesInts")
	}

	if value == nil {
		return nil
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

func CheckSortBy(parameter interface{}, variableName string) error {
	value, ok := parameter.(string)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, "SortBy")
	}

	if value == "" {
		return nil
	}

	validSortBy := []string{"createDateTime", "takerRate", "makerRate", "makerAmount", "takerAmount"}
	if !helpers.Contains(value, validSortBy) {
		return NewParameterValidationError(variableName, fmt.Sprintf("can only contain %v", validSortBy))
	}
	return nil
}

func CheckOrderHashRequired(parameter interface{}, variableName string) error {
	value, ok := parameter.(string)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, "OrderHash")
	}

	if value == "" {
		return NewParameterMissingError(variableName)
	}
	return CheckOrderHash(value, variableName)
}
func CheckOrderHash(parameter interface{}, variableName string) error {
	value, ok := parameter.(string)
	if !ok {
		return fmt.Errorf("for parameter '%v' to be validated as '%v', it must be a string", variableName, "OrderHash")
	}

	if value == "" {
		return nil
	}
	// TODO add criteria that captures valid order hash strings here
	return nil
}
