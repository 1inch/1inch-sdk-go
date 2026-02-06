package tokens

import (
	"github.com/1inch/1inch-sdk-go/internal/validate"
)

func (params *SearchControllerSearchAllChainsParams) Validate() error {
	var validationErrors []error
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *SearchControllerSearchSingleChainParams) Validate() error {
	var validationErrors []error
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *TokenListControllerTokensParams) Validate() error {
	var validationErrors []error
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *CustomTokensControllerGetTokensInfoParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Addresses, "addresses", validate.CheckEthereumAddressListRequired, validationErrors)
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *CustomTokensControllerGetTokenInfoParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Address, "address", validate.CheckEthereumAddressRequired, validationErrors)
	return validate.ConsolidateValidationErrors(validationErrors)
}
