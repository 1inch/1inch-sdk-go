package balance

// BalancesAndAllowancesByWalletAddressListParams is used instead of codegen struct to right now as params for API handle
type BalancesAndAllowancesByWalletAddressListParams struct {
	Wallet  string `json:"-"`
	Spender string `json:"-"`
}

// BalancesByWalletAddressListParams is used instead of codegen struct to right now as params for API handle
type BalancesAndAllowancesByWalletAddressListResponse map[string]TokenDetails

// TokenDetails holds balance and allowance for an Ethereum address (token)
type TokenDetails struct {
	Balance   string `json:"balance"`
	Allowance string `json:"allowance"`
}

// BalancesAndAllowancesOfCustomTokensByWalletAddressListParams is used instead of codegen struct to right now as params for API handle
type BalancesAndAllowancesOfCustomTokensByWalletAddressListParams struct {
	Wallet  string   `json:"-"`
	Spender string   `json:"-"`
	Tokens  []string `json:"tokens"`
}

// BalancesAndAllowancesOfCustomTokensByWalletAddressListResponse is used instead of codegen struct to right now as params for API handle
type BalancesAndAllowancesOfCustomTokensByWalletAddressListResponse map[string]TokenDetails

// BalancesAndAllowancesParams is used instead of codegen struct to right now as params for API handle
type BalancesAndAllowancesParams struct {
	Wallets []string `url:"wallets" json:"wallets"`

	// Will filter tokens with 0 balance from response
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
