package fusionplus

// Pre-v5 names for generated types. The generated names now come from the
// friendly operation-id mapping (codegen/mapping.json) and schema name
// overrides (codegen/overrides.go); these aliases keep existing integrations
// compiling and are identical types.

// Deprecated: Use Quote instead.
type GetQuoteOutput = Quote

// Deprecated: Use GetActiveOrdersParams instead.
type OrderApiControllerGetActiveOrdersParams = GetActiveOrdersParams

// Deprecated: Use GetOrdersByMakerParams instead.
type OrderApiControllerGetOrdersByMakerParams = GetOrdersByMakerParams

// Deprecated: Use GetOrdersByHashesJSONRequestBody instead.
type OrderApiControllerGetOrdersByOrderHashesJSONRequestBody = GetOrdersByHashesJSONRequestBody

// Deprecated: Use GetSettlementContractParams instead.
type OrderApiControllerGetSettlementContractParams = GetSettlementContractParams

// Deprecated: Use BuildQuoteTypedDataJSONRequestBody instead.
type QuoterControllerBuildQuoteTypedDataJSONRequestBody = BuildQuoteTypedDataJSONRequestBody

// Deprecated: Use BuildQuoteTypedDataParams instead.
type QuoterControllerBuildQuoteTypedDataParams = BuildQuoteTypedDataParams

// Deprecated: Use QuoteParams instead.
type QuoterControllerGetQuoteParams = QuoteParams

// Deprecated: Use CustomPresetQuoteJSONRequestBody instead.
type QuoterControllerGetQuoteWithCustomPresetsJSONRequestBody = CustomPresetQuoteJSONRequestBody

// Deprecated: Use CustomPresetQuoteParams instead.
type QuoterControllerGetQuoteWithCustomPresetsParams = CustomPresetQuoteParams

// Deprecated: Use SubmitOrderJSONRequestBody instead.
type RelayerControllerSubmitJSONRequestBody = SubmitOrderJSONRequestBody

// Deprecated: Use SubmitManyOrdersJSONBody instead.
type RelayerControllerSubmitManyJSONBody = SubmitManyOrdersJSONBody

// Deprecated: Use SubmitManyOrdersJSONRequestBody instead.
type RelayerControllerSubmitManyJSONRequestBody = SubmitManyOrdersJSONRequestBody

// Deprecated: Use SubmitSecretsJSONRequestBody instead.
type RelayerControllerSubmitSecretsJSONRequestBody = SubmitSecretsJSONRequestBody
