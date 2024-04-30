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
