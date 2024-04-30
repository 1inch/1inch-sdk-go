package balance

import "github.com/1inch/1inch-sdk-go/internal/validate"

func (params *AggregatedBalancesAndAllowancesParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.FilterEmpty, "filterEmpty", validate.CheckBoolean, validationErrors)
	validationErrors = validate.Parameter(params.Wallets, "wallets", validate.CheckAddressesList, validationErrors)
	validationErrors = validate.Parameter(params.Spender, "spender", validate.CheckEthereumAddressRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}
