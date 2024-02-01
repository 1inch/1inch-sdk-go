package validate

import (
	"fmt"
	"regexp"

	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
)

// CheckEthereumAddressPointer exits early if the pointer is nil because nil pointer parameters are optional
func CheckEthereumAddressPointer(address *string, variableName string) error {
	if address == nil {
		return nil
	}
	return CheckEthereumAddress(*address, variableName)
}

func CheckEthereumAddress(address string, variableName string) error {
	if address == "" {
		return NewParameterMissingError(variableName)
	}

	// Ethereum address starts with '0x' followed by 40 hexadecimal characters.
	re := regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
	isEthereumAddress := re.MatchString(address)
	if !isEthereumAddress {
		return NewParameterValidationError(variableName, "not a valid Ethereum address")
	}
	return nil
}

func CheckBigIntPointer(value *string, variableName string) error {
	if value == nil {
		return nil
	}
	return CheckBigInt(*value, variableName)
}

var maxBigInt, _ = helpers.BigIntFromString("115792089237316195423570985008687907853269984665640564039457584007913129639935")

func CheckBigInt(value string, variableName string) error {
	if value == "" {
		return NewParameterMissingError(variableName)
	}

	parsedValue, err := helpers.BigIntFromString(value)
	if err != nil {
		return NewParameterValidationError(variableName, "not a valid big integer")
	}
	if parsedValue.Cmp(maxBigInt) > 0 {
		return NewParameterValidationError(variableName, "too big to fit in uint256")
	}
	return nil
}

func CheckChainId(value int, variableName string) error {
	if value == 0 {
		return NewParameterMissingError(variableName)
	}
	if !helpers.Contains(value, chains.ValidChainIds) {
		return NewParameterValidationError(variableName, fmt.Sprintf("valid chain ids are: %v", chains.ValidChainIds))
	}
	return nil
}

func CheckPrivateKey(address string, variableName string) error {
	if address == "" {
		return NewParameterMissingError(variableName)
	}

	// Private keys are always 64 hexadecimal characters.
	re := regexp.MustCompile(`^[a-fA-F0-9]{64}$`)
	isPrivateKey := re.MatchString(address)
	if !isPrivateKey {
		return NewParameterValidationError(variableName, "not a valid private key")
	}
	return nil
}

func CheckApprovalType(value int, variableName string) error {
	if !helpers.Contains(value, []int{0, 1, 2}) {
		return NewParameterValidationError(variableName, "invalid approval type")
	}
	return nil
}

func CheckSlippage(value float32, variableName string) error {
	// Slippage of '0' is technically allowed, but it is much more likely the user forgot to set it in their request config, so it is disallowed for now
	if value == 0 {
		return NewParameterMissingError(variableName)
	}
	if value < 0.01 || value > 50 {
		return NewParameterValidationError(variableName, "invalid slippage value - only values 0.01-50 are allowed")
	}
	return nil
}

func CheckPagePointer(value *float32, variableName string) error {
	if value == nil {
		return nil
	}
	return CheckPage(*value, variableName)
}

func CheckPage(page float32, variableName string) error {
	if page < 1 {
		return NewParameterValidationError(variableName, "must be greater than 0")
	}
	return nil
}

func CheckLimitPointer(value *float32, variableName string) error {
	if value == nil {
		return nil
	}
	return CheckLimit(*value, variableName)
}

func CheckLimit(value float32, variableName string) error {
	if value == 0 {
		return NewParameterMissingError(variableName)
	}
	if value < 1 {
		return NewParameterValidationError(variableName, "must be greater than 0")
	}
	return nil
}

func CheckStatusesIntsPointer(value *[]float32, variableName string) error {
	if value == nil {
		return nil
	}
	return CheckStatusesInts(*value, variableName)
}

func CheckStatusesInts(statuses []float32, variableName string) error {
	if statuses == nil {
		return nil
	}
	if helpers.HasDuplicates(statuses) {
		return NewParameterValidationError(variableName, "must not contain duplicates")
	}
	validStatuses := []float32{1, 2, 3}
	if !helpers.IsSubset(statuses, validStatuses) {
		return NewParameterValidationError(variableName, "can only contain 1, 2, and/or 3")
	}
	return nil
}

func CheckStatusesStringsPointer(value *[]string, variableName string) error {
	if value == nil {
		return nil
	}
	return CheckStatusesStrings(*value, variableName)
}

func CheckStatusesStrings(statuses []string, variableName string) error {
	if statuses == nil {
		return nil
	}
	if helpers.HasDuplicates(statuses) {
		return NewParameterValidationError(variableName, "must not contain duplicates")
	}
	validStatuses := []string{"1", "2", "3"}
	if !helpers.IsSubset(statuses, validStatuses) {
		return NewParameterValidationError(variableName, "can only contain 1, 2, and/or 3")
	}
	return nil
}

func CheckSortByPointer(value *string, variableName string) error {
	if value == nil {
		return nil
	}
	return CheckSortBy(*value, variableName)
}

func CheckSortBy(sortBy string, variableName string) error {
	validSortBy := []string{"createDateTime", "takerRate", "makerRate", "makerAmount", "takerAmount"}
	if !helpers.Contains(sortBy, validSortBy) {
		return NewParameterValidationError(variableName, "can only contain createDateTime, takerRate, makerRate, makerAmount, or takerAmount")
	}
	return nil
}

func CheckOrderHash(value string, variableName string) error {
	if value == "" {
		return NewParameterMissingError(variableName)
	}
	// TODO add criteria that captures valid order hash strings here
	return nil
}
