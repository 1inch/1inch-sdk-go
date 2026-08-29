package tokens

import (
	"github.com/1inch/1inch-sdk-go/v5/internal/validate"
)

func (params *SearchAllChainsParams) Validate() error {
	var validationErrors []error
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *SearchSingleChainParams) Validate() error {
	var validationErrors []error
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *GetWhitelistedTokensParams) Validate() error {
	var validationErrors []error
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *GetCustomTokensParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Addresses, "addresses", validate.CheckEthereumAddressListRequired, validationErrors)
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *GetCustomTokenParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Address, "address", validate.CheckEthereumAddressRequired, validationErrors)
	return validate.ConsolidateValidationErrors(validationErrors)
}
