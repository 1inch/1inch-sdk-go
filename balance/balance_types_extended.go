package balance

// AggregatedBalancesAndAllowancesParams is used instead of codegen struct to right now as params for API handle
type AggregatedBalancesAndAllowancesParams struct {
	// List of EVM addresses
	Wallets []string `url:"wallets" json:"wallets"`

	// Will filter tokens with 0 balance from response
	FilterEmpty bool `url:"filterEmpty" json:"filterEmpty"`

	// EVM address of token spender (based on erc20 spec)
	Spender string
}

// BalancesByWalletAddressParams is used instead of codegen struct to right now as params for API handle
type BalancesByWalletAddressParams struct {
	WalletAddress string `url:"wallets" json:"walletAddress"`
}

// BalancesByWalletAddressParams is used instead of codegen struct to right now as params for API handle
type BalancesByWalletAddressResponse map[string]string

// BalancesOfCustomTokensByWalletAddressParams is used instead of codegen struct to right now as params for API handle
type BalancesOfCustomTokensByWalletAddressParams struct {
	WalletAddress string   `url:"wallets" json:"_"`
	Tokens        []string `url:"tokens" json:"tokens"`
}

// BalancesByWalletAddressParams is used instead of codegen struct to right now as params for API handle
type BalancesOfCustomTokensByWalletAddressResponse map[string]string

// BalancesOfCustomTokensByWalletAddressParams is used instead of codegen struct to right now as params for API handle
type BalancesOfCustomTokensByWalletAddressesListParams struct {
	WalletAddresses []string `url:"wallets" json:"wallets"`
	Tokens          []string `url:"tokens" json:"tokens"`
}

// BalancesByWalletAddressParams is used instead of codegen struct to right now as params for API handle
type BalancesOfCustomTokensByWalletAddressesListResponse map[string]string
