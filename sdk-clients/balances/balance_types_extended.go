package balances

// BalancesAndAllowancesByWalletAddressListParams is used instead of codegen struct to right now as params for API handle
type BalancesAndAllowancesByWalletAddressListParams struct {
	Wallet  string `json:"-"`
	Spender string `json:"-"`
}

// BalancesByWalletAddressListParams is used instead of codegen struct to right now as params for API handle
type BalancesAndAllowancesByWalletAddressListResponse map[string]TokenDetails

// TokenDetails holds balances and allowance for an Ethereum address (token)
type TokenDetails struct {
	Balance   string `json:"balance"`
	Allowance string `json:"allowance"`
}

// BalancesAndAllowancesOfCustomTokensByWalletAddressParams is used instead of codegen struct to right now as params for API handle
type BalancesAndAllowancesOfCustomTokensByWalletAddressParams struct {
	Wallet  string   `json:"-"`
	Spender string   `json:"-"`
	Tokens  []string `json:"tokens"`
}

// BalancesAndAllowancesOfCustomTokensByWalletAddressResponse is used instead of codegen struct to right now as params for API handle
type BalancesAndAllowancesOfCustomTokensByWalletAddressResponse map[string]TokenDetails

// BalancesAndAllowancesParams is used instead of codegen struct to right now as params for API handle
type BalancesAndAllowancesParams struct {
	Wallets []string `url:"wallets" json:"wallets"`

	// Will filter tokens with 0 balances from response
	FilterEmpty bool `url:"filterEmpty" json:"filterEmpty"`

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
	Wallets string   `url:"wallets" json:"_"`
	Tokens  []string `url:"tokens" json:"tokens"`
}

// BalancesByWalletAddressParams is used instead of codegen struct to right now as params for API handle
type BalancesOfCustomTokensByWalletAddressResponse map[string]string

// BalancesOfCustomTokensByWalletAddressParams is used instead of codegen struct to right now as params for API handle
type BalancesOfCustomTokensByWalletAddressesListParams struct {
	Wallets []string `url:"wallets" json:"wallets"`
	Tokens  []string `url:"tokens" json:"tokens"`
}

// BalancesByWalletAddressParams is used instead of codegen struct to right now as params for API handle
type BalancesOfCustomTokensByWalletAddressesListResponse map[string]map[string]string

// AllowancesByWalletAddressParams is used instead of codegen struct to right now as params for API handle
type AllowancesByWalletAddressParams struct {
	Wallet  string `json:"-"`
	Spender string `json:"-"`
}

// AllowancesByWalletAddressResponse is used instead of codegen struct to right now as params for API handle
type AllowancesByWalletAddressResponse map[string]string

// AllowancesByWalletAddressParams is used instead of codegen struct to right now as params for API handle
type AllowancesOfCustomTokensByWalletAddressParams struct {
	Wallet  string   `json:"-"`
	Spender string   `json:"-"`
	Tokens  []string `url:"tokens" json:"tokens"`
}

// AllowancesByWalletAddressResponse is used instead of codegen struct to right now as params for API handle
type AllowancesOfCustomTokensByWalletAddressResponse map[string]string
