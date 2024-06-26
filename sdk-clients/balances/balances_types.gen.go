// Package balances provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package balances

// AggregatedBalancesAndAllowancesResponse defines model for AggregatedBalancesAndAllowancesResponse.
type AggregatedBalancesAndAllowancesResponse = []struct {
	// Address Token contract address
	Address *string `json:"address,omitempty"`

	// Decimals Number of decimal places for the token
	Decimals *int `json:"decimals,omitempty"`

	// IsCustom Indicates if the token is custom or not
	IsCustom *bool `json:"isCustom,omitempty"`

	// LogoURI URL to the token logo image
	LogoURI *string `json:"logoURI,omitempty"`

	// Name Name of the token
	Name *string `json:"name,omitempty"`

	// Symbol Symbol of the token
	Symbol *string `json:"symbol,omitempty"`

	// Tags Additional tags for the token
	Tags *[]string `json:"tags,omitempty"`

	// Tracked Indicates if the token is tracked or not
	Tracked *bool `json:"tracked,omitempty"`

	// Type Type of token (e.g., 'ethereum', 'token')
	Type *string `json:"type,omitempty"`

	// Wallets Token balances and allowances for specific wallets
	Wallets *map[string]struct {
		// Allowance Allowance of the token for the wallet
		Allowance *string `json:"allowance,omitempty"`

		// Balance Balance of the token for the wallet
		Balance *string `json:"balance,omitempty"`
	} `json:"wallets,omitempty"`
}

// CustomTokensAndWalletsRequest defines model for CustomTokensAndWalletsRequest.
type CustomTokensAndWalletsRequest struct {
	// Tokens List of custom tokens
	Tokens []string `json:"tokens"`

	// Wallets List of wallets
	Wallets []string `json:"wallets"`
}

// CustomTokensRequest defines model for CustomTokensRequest.
type CustomTokensRequest struct {
	// Tokens List of custom tokens
	Tokens []string `json:"tokens"`
}

// GetAggregatedBalancesAndAllowancesParams defines parameters for GetAggregatedBalancesAndAllowances.
type GetAggregatedBalancesAndAllowancesParams struct {
	Wallets     []string `url:"wallets" json:"wallets"`
	FilterEmpty bool     `url:"filterEmpty" json:"filterEmpty"`
}

// ChainV12ControllerGetCustomAllowancesJSONRequestBody defines body for ChainV12ControllerGetCustomAllowances for application/json ContentType.
type ChainV12ControllerGetCustomAllowancesJSONRequestBody = CustomTokensRequest

// ChainV12ControllerGetCustomAllowancesAndBalancesJSONRequestBody defines body for ChainV12ControllerGetCustomAllowancesAndBalances for application/json ContentType.
type ChainV12ControllerGetCustomAllowancesAndBalancesJSONRequestBody = CustomTokensRequest

// ChainV12ControllerGetBalancesByMultipleWalletsJSONRequestBody defines body for ChainV12ControllerGetBalancesByMultipleWallets for application/json ContentType.
type ChainV12ControllerGetBalancesByMultipleWalletsJSONRequestBody = CustomTokensAndWalletsRequest

// ChainV12ControllerGetCustomBalancesJSONRequestBody defines body for ChainV12ControllerGetCustomBalances for application/json ContentType.
type ChainV12ControllerGetCustomBalancesJSONRequestBody = CustomTokensRequest
