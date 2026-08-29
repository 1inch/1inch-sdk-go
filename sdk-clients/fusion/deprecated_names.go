package fusion

// Pre-v5 names for generated types. The generated names now come from the
// friendly operation-id mapping (codegen/mapping.json) and schema name
// overrides (codegen/overrides.go); these aliases keep existing integrations
// compiling and are identical types.

// Deprecated: Use Quote instead.
type GetQuoteOutput = Quote

// Deprecated: Use GetActiveOrdersParams instead.
type OrderApiControllerGetActiveOrdersParams = GetActiveOrdersParams

// Deprecated: Use QuoteParams instead.
type QuoterControllerGetQuoteParams = QuoteParams

// Deprecated: Use CustomPresetQuoteJSONRequestBody instead.
type QuoterControllerGetQuoteWithCustomPresetsJSONRequestBody = CustomPresetQuoteJSONRequestBody

// Deprecated: Use CustomPresetQuoteParams instead.
type QuoterControllerGetQuoteWithCustomPresetsParams = CustomPresetQuoteParams

// Deprecated: Use SubmitOrder400JSONResponseBodyError instead.
type RelayerControllerSubmit400JSONResponseBodyError = SubmitOrder400JSONResponseBodyError

// Deprecated: Use SubmitOrder400JSONResponseBodyStatusCode instead.
type RelayerControllerSubmit400JSONResponseBodyStatusCode = SubmitOrder400JSONResponseBodyStatusCode

// Deprecated: Use SubmitOrderJSONRequestBody instead.
type RelayerControllerSubmitJSONRequestBody = SubmitOrderJSONRequestBody

// Deprecated: Use SubmitManyOrders400JSONResponseBodyError instead.
type RelayerControllerSubmitMany400JSONResponseBodyError = SubmitManyOrders400JSONResponseBodyError

// Deprecated: Use SubmitManyOrders400JSONResponseBodyStatusCode instead.
type RelayerControllerSubmitMany400JSONResponseBodyStatusCode = SubmitManyOrders400JSONResponseBodyStatusCode

// Deprecated: Use SubmitManyOrdersJSONBody instead.
type RelayerControllerSubmitManyJSONBody = SubmitManyOrdersJSONBody

// Deprecated: Use SubmitManyOrdersJSONRequestBody instead.
type RelayerControllerSubmitManyJSONRequestBody = SubmitManyOrdersJSONRequestBody

const (
	// Deprecated: Use SubmitOrder400JSONResponseBodyErrorBadRequest instead.
	RelayerControllerSubmit400JSONResponseBodyErrorBadRequest = SubmitOrder400JSONResponseBodyErrorBadRequest
	// Deprecated: Use SubmitOrder400JSONResponseBodyStatusCodeN400 instead.
	RelayerControllerSubmit400JSONResponseBodyStatusCodeN400 = SubmitOrder400JSONResponseBodyStatusCodeN400
	// Deprecated: Use SubmitManyOrders400JSONResponseBodyErrorBadRequest instead.
	RelayerControllerSubmitMany400JSONResponseBodyErrorBadRequest = SubmitManyOrders400JSONResponseBodyErrorBadRequest
	// Deprecated: Use SubmitManyOrders400JSONResponseBodyStatusCodeN400 instead.
	RelayerControllerSubmitMany400JSONResponseBodyStatusCodeN400 = SubmitManyOrders400JSONResponseBodyStatusCodeN400
)
