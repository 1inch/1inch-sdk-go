package spotprices

import (
	"github.com/1inch/1inch-sdk-go/internal/validate"
)

func (params *GetWhitelistedTokensPricesParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(string(params.Currency), "currency", validate.CheckFiatCurrency, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *GetPricesRequestDto) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Tokens, "Tokens", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(string(params.Currency), "currency", validate.CheckFiatCurrency, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}
