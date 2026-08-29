package portfolio

import (
	"github.com/1inch/1inch-sdk-go/v5/internal/validate"
)

func (params *GetProtocolsCurrentValueParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Addresses, "Addresses", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "ChainId", validate.CheckChainIdInt, validationErrors)
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *GetProtocolsProfitAndLossParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Addresses, "Addresses", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "ChainId", validate.CheckChainIdInt, validationErrors)
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *GetCurrentValueParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Addresses, "Addresses", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "ChainId", validate.CheckChainIdInt, validationErrors)
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *GetProfitAndLossParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Addresses, "Addresses", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "ChainId", validate.CheckChainIdInt, validationErrors)
	//validationErrors = validate.Parameter(params.Timerange, "Timerange", validate.CheckTimerange, validationErrors)  // TODO "x-go-type-skip-optional-pointer": true does not work as expected for parameters of type schema. Need to research this
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *GetProtocolsDetailsParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Addresses, "Addresses", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "ChainId", validate.CheckChainIdInt, validationErrors)
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *GetTokensCurrentValueParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Addresses, "Addresses", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "ChainId", validate.CheckChainIdInt, validationErrors)
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *GetTokensProfitAndLossParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Addresses, "Addresses", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "ChainId", validate.CheckChainIdInt, validationErrors)
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *GetTokensDetailsParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Addresses, "Addresses", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "ChainId", validate.CheckChainIdInt, validationErrors)
	//validationErrors = validate.Parameter(params.Timerange, "Timerange", validate.CheckTimerange, validationErrors)  // TODO "x-go-type-skip-optional-pointer": true does not work as expected for parameters of type schema. Need to research this
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *GetValueChartParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Addresses, "Addresses", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "ChainId", validate.CheckChainIdInt, validationErrors)
	//validationErrors = validate.Parameter(params.Timerange, "Timerange", validate.CheckTimerange, validationErrors)  // TODO "x-go-type-skip-optional-pointer": true does not work as expected for parameters of type schema. Need to research this
	return validate.ConsolidateValidationErrors(validationErrors)
}
