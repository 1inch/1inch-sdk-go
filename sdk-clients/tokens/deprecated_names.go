package tokens

// Pre-v5 names for generated types. The generated names now come from the
// friendly operation-id mapping (codegen/mapping.json) and schema name
// overrides (codegen/overrides.go); these aliases keep existing integrations
// compiling and are identical types.

// Deprecated: Use GetCustomTokensParams instead.
type CustomTokensControllerGetTokensInfoParams = GetCustomTokensParams

// Deprecated: Use SearchAllChainsParams instead.
type SearchControllerSearchAllChainsParams = SearchAllChainsParams

// Deprecated: Use SearchSingleChainParams instead.
type SearchControllerSearchSingleChainParams = SearchSingleChainParams

// Deprecated: Use GetWhitelistedTokensListParams instead.
type TokenListControllerTokensListParams = GetWhitelistedTokensListParams

// Deprecated: Use GetWhitelistedTokensParams instead.
type TokenListControllerTokensParams = GetWhitelistedTokensParams
