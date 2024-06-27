package fusion

import (
	"fmt"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/internal/validate"
)

func (params *OrderApiControllerGetActiveOrdersParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Page, "Page", validate.CheckPage, validationErrors)
	validationErrors = validate.Parameter(params.Limit, "Limit", validate.CheckLimit, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *QuoterControllerGetQuoteParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.FromTokenAddress, "FromTokenAddress", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.ToTokenAddress, "ToTokenAddress", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Permit, "Permit", validate.CheckPermitHash, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *QuoterControllerGetQuoteWithCustomPresetsParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.FromTokenAddress, "FromTokenAddress", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.ToTokenAddress, "ToTokenAddress", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Amount, "Amount", validate.CheckBigIntRequired, validationErrors)
	validationErrors = validate.Parameter(params.WalletAddress, "WalletAddress", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.WalletAddress, "WalletAddress", validate.CheckEthereumAddressRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (body *PlaceOrderBody) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(body.Maker, "Maker", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(body.MakerAsset, "MakerAsset", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(body.MakingAmount, "MakingAmount", validate.CheckBigIntRequired, validationErrors)
	validationErrors = validate.Parameter(body.Receiver, "Receiver", validate.CheckEthereumAddressRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (body *OrderParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(body.Receiver, "Receiver", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(body.WalletAddress, "WalletAddress", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(body.FromTokenAddress, "FromTokenAddress", validate.CheckEthereumAddress, validationErrors)
	validationErrors = validate.Parameter(body.ToTokenAddress, "ToTokenAddress", validate.CheckEthereumAddress, validationErrors)
	validationErrors = validate.Parameter(body.Amount, "Amount", validate.CheckBigInt, validationErrors)
	validationErrors = validate.Parameter(body.Permit, "Permit", validate.CheckPermitHash, validationErrors)
	if body.Preset == "" {
		validationErrors = append(validationErrors, validate.NewParameterCustomError(fmt.Sprintf("Preset is required. Pass in one of the Fusion library constants: %v", constants.ValidFusionPresets)))
	}
	return validate.ConsolidateValidationErorrs(validationErrors)
}
