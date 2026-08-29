package balances

// Compile-time assertions that the deprecated "List" type names still resolve.
var (
	_ BalancesAndAllowancesByWalletAddressListParams    = BalancesAndAllowancesByWalletAddressParams{}
	_ BalancesOfCustomTokensByWalletAddressesListParams = BalancesOfCustomTokensByWalletAddressesParams{}
)
