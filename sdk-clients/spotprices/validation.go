package spotprices

import (
	"github.com/1inch/1inch-sdk-go/internal/validate"
)

func (params *GetWhitelistedTokensPricesParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Currency, "currency", validate.CheckFiatCurrency, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *GetPricesRequestDto) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Tokens, "Tokens", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(params.Currency, "currency", validate.CheckFiatCurrency, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *BalancesOfCustomTokensByWalletAddressParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Wallet, "walletAddress", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Tokens, "tokens", validate.CheckEthereumAddressListRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *BalancesOfCustomTokensByWalletAddressesListParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Wallets, "wallets", validate.CheckEthereumAddressListRequired, validationErrors)
	validationErrors = validate.Parameter(params.Tokens, "tokens", validate.CheckEthereumAddressListRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *BalancesAndAllowancesByWalletAddressListParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Wallet, "wallet", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Spender, "spender", validate.CheckEthereumAddressRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *BalancesAndAllowancesOfCustomTokensByWalletAddressParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Wallet, "wallet", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Spender, "spender", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Tokens, "tokens", validate.CheckEthereumAddressListRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *AllowancesByWalletAddressParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Wallet, "wallet", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Spender, "spender", validate.CheckEthereumAddressRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *AllowancesOfCustomTokensByWalletAddressParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Wallet, "wallet", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Spender, "spender", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Tokens, "tokens", validate.CheckEthereumAddressListRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}
