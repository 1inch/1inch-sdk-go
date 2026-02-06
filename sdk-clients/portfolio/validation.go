package portfolio

import (
	"github.com/1inch/1inch-sdk-go/internal/validate"
)

func (params *GetCurrentValuePortfolioV4OverviewProtocolsCurrentValueGetParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Addresses, "Addresses", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "ChainId", validate.CheckChainIdInt, validationErrors)
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *GetProfitAndLossPortfolioV4OverviewProtocolsProfitAndLossGetParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Addresses, "Addresses", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "ChainId", validate.CheckChainIdInt, validationErrors)
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *GetCurrentValuePortfolioV4GeneralCurrentValueGetParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Addresses, "Addresses", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "ChainId", validate.CheckChainIdInt, validationErrors)
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *GetProfitAndLossPortfolioV4GeneralProfitAndLossGetParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Addresses, "Addresses", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "ChainId", validate.CheckChainIdInt, validationErrors)
	//validationErrors = validate.Parameter(params.Timerange, "Timerange", validate.CheckTimerange, validationErrors)  // TODO "x-go-type-skip-optional-pointer": true does not work as expected for parameters of type schema. Need to research this
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *GetDetailsPortfolioV4OverviewProtocolsDetailsGetParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Addresses, "Addresses", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "ChainId", validate.CheckChainIdInt, validationErrors)
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *GetCurrentValuePortfolioV4OverviewErc20CurrentValueGetParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Addresses, "Addresses", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "ChainId", validate.CheckChainIdInt, validationErrors)
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *GetProfitAndLossPortfolioV4OverviewErc20ProfitAndLossGetParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Addresses, "Addresses", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "ChainId", validate.CheckChainIdInt, validationErrors)
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *GetDetailsPortfolioV4OverviewErc20DetailsGetParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Addresses, "Addresses", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "ChainId", validate.CheckChainIdInt, validationErrors)
	//validationErrors = validate.Parameter(params.Timerange, "Timerange", validate.CheckTimerange, validationErrors)  // TODO "x-go-type-skip-optional-pointer": true does not work as expected for parameters of type schema. Need to research this
	return validate.ConsolidateValidationErrors(validationErrors)
}

func (params *GetValueChartPortfolioV4GeneralValueChartGetParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Addresses, "Addresses", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(params.ChainId, "ChainId", validate.CheckChainIdInt, validationErrors)
	//validationErrors = validate.Parameter(params.Timerange, "Timerange", validate.CheckTimerange, validationErrors)  // TODO "x-go-type-skip-optional-pointer": true does not work as expected for parameters of type schema. Need to research this
	return validate.ConsolidateValidationErrors(validationErrors)
}
