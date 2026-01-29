package fusionorder

import (
	"errors"
	"fmt"

	"github.com/1inch/1inch-sdk-go/internal/hexadecimal"
)

// Prefix0x ensures a hex string has the 0x prefix
func Prefix0x(value string) string {
	if len(value) >= 2 && value[:2] == "0x" {
		return value
	}
	return "0x" + value
}

// ExtensionHexParams contains hex string parameters that need validation
type ExtensionHexParams struct {
	SettlementContract string
	MakerAssetSuffix   string
	TakerAssetSuffix   string
	Predicate          string
	CustomData         string
}

// ValidateExtensionHexParams validates common hex parameters used in extension creation
func ValidateExtensionHexParams(params ExtensionHexParams) error {
	if !hexadecimal.IsHexBytes(params.SettlementContract) {
		return fmt.Errorf("invalid settlement contract hex: %s", params.SettlementContract)
	}
	if !hexadecimal.IsHexBytes(params.MakerAssetSuffix) {
		return fmt.Errorf("invalid maker asset suffix hex: %s", params.MakerAssetSuffix)
	}
	if !hexadecimal.IsHexBytes(params.TakerAssetSuffix) {
		return fmt.Errorf("invalid taker asset suffix hex: %s", params.TakerAssetSuffix)
	}
	if !hexadecimal.IsHexBytes(params.Predicate) {
		return fmt.Errorf("invalid predicate hex: %s", params.Predicate)
	}
	if params.CustomData != "" {
		return errors.New("unsupported: custom data")
	}
	return nil
}
