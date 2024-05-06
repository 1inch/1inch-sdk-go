package balance

import "github.com/1inch/1inch-sdk-go/internal/validate"

func (params *BalancesAndAllowancesParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.FilterEmpty, "filterEmpty", validate.CheckBoolean, validationErrors)
	validationErrors = validate.Parameter(params.Wallets, "wallets", validate.CheckAddressesList, validationErrors)
	validationErrors = validate.Parameter(params.Spender, "spender", validate.CheckEthereumAddressRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *BalancesByWalletAddressParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.WalletAddress, "walletAddress", validate.CheckEthereumAddressRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *BalancesOfCustomTokensByWalletAddressParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.WalletAddress, "walletAddress", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Tokens, "tokens", validate.CheckAddressesList, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *BalancesOfCustomTokensByWalletAddressesListParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.WalletAddresses, "wallets", validate.CheckAddressesList, validationErrors)
	validationErrors = validate.Parameter(params.Tokens, "tokens", validate.CheckAddressesList, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *BalancesAndAllowancesByWalletAddressListParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Wallet, "wallet", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Spender, "spender", validate.CheckEthereumAddressRequired, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}

func (params *BalancesAndAllowancesOfCustomTokensByWalletAddressListParams) Validate() error {
	var validationErrors []error
	validationErrors = validate.Parameter(params.Wallet, "wallet", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Spender, "spender", validate.CheckEthereumAddressRequired, validationErrors)
	validationErrors = validate.Parameter(params.Tokens, "spender", validate.CheckAddressesList, validationErrors)
	return validate.ConsolidateValidationErorrs(validationErrors)
}
