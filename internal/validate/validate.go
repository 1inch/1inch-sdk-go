package validate

import (
	"fmt"
	"math/big"
	"regexp"
	"slices"
	"strings"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/internal/bigint"
)

// Pre-compiled regexes to avoid recompilation on every validation call
var (
	ethAddressRegex      = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
	privateKeyRegex      = regexp.MustCompile(`^[a-fA-F0-9]{64}$`)
	protocolsRegex       = regexp.MustCompile(`^[a-zA-Z0-9_]+(,[a-zA-Z0-9_]+)*$`)
	connectorTokensRegex = regexp.MustCompile(`^0x[a-fA-F0-9]{40}(,0x[a-fA-F0-9]{40})*$`)
	permitHashRegex      = regexp.MustCompile(`^0x[a-fA-F0-9]*$`)
)

func CheckEthereumAddressRequired(value string, variableName string) error {
	if value == "" {
		return NewParameterMissingError(variableName)
	}

	return CheckEthereumAddress(value, variableName)
}

func CheckEthereumAddress(value string, variableName string) error {
	if value == "" {
		return nil
	}

	if !ethAddressRegex.MatchString(value) {
		return NewParameterValidationError(variableName, "not a valid Ethereum address")
	}
	return nil
}

func CheckEthereumAddressListRequired(addresses []string, variableName string) error {
	if len(addresses) == 0 {
		return NewParameterMissingError(variableName)
	}

	for _, address := range addresses {
		if address == "" {
			return NewParameterMissingError(variableName)
		}
		if !ethAddressRegex.MatchString(address) {
			return NewParameterValidationError(variableName, "not a valid Ethereum address")
		}
	}

	return nil
}

var bigIntMax, _ = bigint.FromString("115792089237316195423570985008687907853269984665640564039457584007913129639935")
var bigIntZero = big.NewInt(0)

func CheckBigIntRequired(value string, variableName string) error {
	if value == "" {
		return NewParameterMissingError(variableName)
	}

	return CheckBigInt(value, variableName)
}
func CheckBigInt(value string, variableName string) error {
	if value == "" {
		return nil
	}

	parsedValue, err := bigint.FromString(value)
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

func CheckChainIdIntRequired(value int, variableName string) error {
	if value == 0 {
		return NewParameterMissingError(variableName)
	}

	return CheckChainIdInt(value, variableName)
}

func CheckChainIdInt(value int, variableName string) error {
	if value == 0 {
		return nil
	}

	if !slices.Contains(constants.ValidChainIds, value) {
		return NewParameterValidationError(variableName, fmt.Sprintf("is invalid, valid chain ids are: %v", constants.ValidChainIds))
	}
	return nil
}

func CheckChainIdFloat32Required(value float32, variableName string) error {
	if value == 0 {
		return NewParameterMissingError(variableName)
	}

	return CheckChainIdFloat32(value, variableName)
}

func CheckChainIdFloat32(value float32, variableName string) error {
	if value == 0 {
		return nil
	}

	if !slices.Contains(constants.ValidChainIds, int(value)) {
		return NewParameterValidationError(variableName, fmt.Sprintf("is invalid, valid chain ids are: %v", constants.ValidChainIds))
	}
	return nil
}

func CheckPrivateKeyRequired(value string, variableName string) error {
	if value == "" {
		return NewParameterMissingError(variableName)
	}

	return CheckPrivateKey(value, variableName)
}

func CheckPrivateKey(value string, variableName string) error {
	if value == "" {
		return nil
	}

	if !privateKeyRegex.MatchString(value) {
		return NewParameterValidationError(variableName, "not a valid private key")
	}
	return nil
}

func CheckApprovalType(value int, variableName string) error {
	if value == 0 {
		return nil
	}

	if !slices.Contains([]int{0, 1, 2}, value) {
		return NewParameterValidationError(variableName, "invalid approval type")
	}
	return nil
}

func CheckSlippageRequired(value float32, variableName string) error {
	if value == 0 {
		return NewParameterMissingError(variableName)
	}
	return CheckSlippage(value, variableName)
}

func CheckSlippage(value float32, variableName string) error {
	if value == 0 {
		return nil
	}
	if value < 0.01 || value > 50 {
		return NewParameterValidationError(variableName, fmt.Sprintf("invalid slippage value (%v) - only values 0.01-50 are allowed", value))
	}
	return nil
}

func CheckPage(value float32, variableName string) error {
	if value == 0 {
		return nil
	}

	if value < 1 {
		return NewParameterValidationError(variableName, "must be greater than 0")
	}
	return nil
}

func CheckLimit(value float32, variableName string) error {
	if value == 0 {
		return nil
	}

	if value < 1 {
		return NewParameterValidationError(variableName, "must be greater than 0")
	}
	return nil
}

func CheckStatusesInts(value []float32, variableName string) error {
	if value == nil {
		return nil
	}

	if HasDuplicates(value) {
		return NewParameterValidationError(variableName, "must not contain duplicates")
	}
	validStatuses := []float32{1, 2, 3}
	if !IsSubset(value, validStatuses) {
		return NewParameterValidationError(variableName, fmt.Sprintf("can only contain %v", validStatuses))
	}
	return nil
}

func CheckStatusesStrings(value []string, variableName string) error {
	if HasDuplicates(value) {
		return NewParameterValidationError(variableName, "must not contain duplicates")
	}
	validStatuses := []string{"1", "2", "3"}
	if !IsSubset(value, validStatuses) {
		return NewParameterValidationError(variableName, fmt.Sprintf("can only contain %v", validStatuses))
	}
	return nil
}

func CheckStatusesOrderStatus(value []int, variableName string) error {
	if value == nil {
		return nil
	}

	if HasDuplicates(value) {
		return NewParameterValidationError(variableName, "must not contain duplicates")
	}
	validStatuses := []int{1, 2, 3}
	if !IsSubset(value, validStatuses) {
		return NewParameterValidationError(variableName, fmt.Sprintf("can only contain %v", validStatuses))
	}
	return nil
}

func CheckSortBy(value string, variableName string) error {
	if value == "" {
		return nil
	}

	validSortBy := []string{"createDateTime", "takerRate", "makerRate", "makerAmount", "takerAmount"}
	if !slices.Contains(validSortBy, value) {
		return NewParameterValidationError(variableName, fmt.Sprintf("can only contain %v", validSortBy))
	}
	return nil
}

func CheckOrderHashRequired(value string, variableName string) error {
	if value == "" {
		return NewParameterMissingError(variableName)
	}
	return CheckOrderHash(value, variableName)
}

func CheckOrderHash(value string, variableName string) error {
	if value == "" {
		return nil
	}
	// TODO add criteria that captures valid order hash strings here
	return nil
}

func CheckProtocols(value string, variableName string) error {
	if value == "" {
		return nil
	}

	if !protocolsRegex.MatchString(value) {
		return NewParameterValidationError(variableName, "must be formatted as a single-string list exactly in the format 'Protocol1,Protocol2,Protocol3' without any "+
			"spaces between each protocol name. These names must match the exact protocol id used by the 1inch APIs "+
			"(use the Swap service's GetLiquiditySources function to see this list). Additionally, there cannot be a trailing comma at the end of the list.")
	}

	protocols := strings.Split(value, ",")
	protocolsMap := make(map[string]bool)
	for _, protocol := range protocols {
		if _, exists := protocolsMap[protocol]; exists {
			return NewParameterValidationError(variableName, "Duplicate protocol found in list")
		}
		protocolsMap[protocol] = true
	}

	return nil
}

func CheckFee(value float32, variableName string) error {
	if value < 0 {
		return NewParameterValidationError(variableName, "must be a positive value")
	}

	if value > 3 {
		return NewParameterValidationError(variableName, "must be a value between 0 and 3")
	}

	return nil
}

func CheckFloat32NonNegativeWhole(value float32, variableName string) error {
	if value < 0 {
		return NewParameterValidationError(variableName, "must be 0 or greater")
	}

	// Cast it to an int to truncate it, then cast back to float32 for the comparison
	if float32(int(value)) != value {
		return NewParameterValidationError(variableName, "must be an integer")
	}

	return nil
}

func CheckConnectorTokens(value string, variableName string) error {
	if value == "" {
		return nil
	}

	if !connectorTokensRegex.MatchString(value) {
		return NewParameterValidationError(variableName, "must be formatted as a single-string list exactly in the format '0x123,0x456,0x789' "+
			"without any spaces between each protocol name. Additionally, there cannot be a trailing comma at the end of the list.")
	}

	// Split the string by commas to get individual addresses
	addresses := strings.Split(value, ",")

	// Use a map to check for duplicates
	addressesMap := make(map[string]bool)

	for _, address := range addresses {
		if _, exists := addressesMap[address]; exists {
			return NewParameterValidationError(variableName, "Duplicate address found in list")
		}
		addressesMap[address] = true
	}

	return nil
}

func CheckPermitHash(value string, variableName string) error {
	if value == "" {
		return nil
	}

	if !permitHashRegex.MatchString(value) {
		return NewParameterValidationError(variableName, "not a valid permit hash")
	}
	return nil
}

func CheckFiatCurrency(value string, variableName string) error {
	if len(value) != 3 {
		return NewParameterValidationError(variableName, "must have len = 3 (like USD, EUR, etc)")
	}

	return nil
}

func CheckTimerange(value string, variableName string) error {
	validTimerangeValues := []string{"1day", "1week", "1month", "1year", "3years"}
	if !slices.Contains(validTimerangeValues, value) {
		return NewParameterValidationError(variableName, fmt.Sprintf("is invalid, valid timerange values are: %v", validTimerangeValues))
	}
	return nil
}

func CheckJsonRpcVersionRequired(value string, variableName string) error {
	if value == "" {
		return NewParameterMissingError(variableName)
	}

	return CheckJsonRpcVersion(value, variableName)
}

func CheckJsonRpcVersion(value string, variableName string) error {
	validJsonRpcValues := []string{"1.0", "2.0"}
	if !slices.Contains(validJsonRpcValues, value) {
		return NewParameterValidationError(variableName, fmt.Sprintf("is invalid, valid rpc version are: %v", validJsonRpcValues))
	}
	return nil
}

func CheckNodeType(value string, variableName string) error {
	validNodeTypes := []string{"archive", "full"}
	if !slices.Contains(validNodeTypes, value) {
		return NewParameterValidationError(variableName, fmt.Sprintf("is invalid, valid node types are: %v", validNodeTypes))
	}
	return nil
}
