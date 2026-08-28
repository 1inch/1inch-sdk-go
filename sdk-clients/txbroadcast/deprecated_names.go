package txbroadcast

// Pre-v5 names for generated types. The generated names now come from the
// friendly operation-id mapping (codegen/mapping.json) and schema name
// overrides (codegen/overrides.go); these aliases keep existing integrations
// compiling and are identical types.

// Deprecated: Use BroadcastPrivateTransactionJSONRequestBody instead.
type TxProcessorApiControllerBroadcastFlashbotsTransactionJSONRequestBody = BroadcastPrivateTransactionJSONRequestBody

// Deprecated: Use BroadcastPublicTransactionJSONRequestBody instead.
type TxProcessorApiControllerBroadcastTransactionJSONRequestBody = BroadcastPublicTransactionJSONRequestBody
