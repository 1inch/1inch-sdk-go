package orderbook

// Pre-v5 names for generated types. The generated names now come from the
// friendly operation-id mapping (codegen/mapping.json) and schema name
// overrides (codegen/overrides.go); these aliases keep existing integrations
// compiling and are identical types.

// Deprecated: Use CreateLimitOrderJSONRequestBody instead.
type LimitOrderV3SubscribedApiControllerCreateLimitOrderJSONRequestBody = CreateLimitOrderJSONRequestBody

// Deprecated: Use LimitOrdersQueryParams instead.
type LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams = LimitOrdersQueryParams

// Deprecated: Use LimitOrdersQueryParamsSortBy instead.
type LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortBy = LimitOrdersQueryParamsSortBy

// Deprecated: Use OrdersCountQueryParams instead.
type LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams = OrdersCountQueryParams

// Deprecated: Use EventsQueryParams instead.
type LimitOrderV3SubscribedApiControllerGetEventsParams = EventsQueryParams

// Deprecated: Use GetLimitOrderParams instead.
type LimitOrderV3SubscribedApiControllerGetLimitOrderParams = GetLimitOrderParams

// Deprecated: Use GetLimitOrderParamsSortBy instead.
type LimitOrderV3SubscribedApiControllerGetLimitOrderParamsSortBy = GetLimitOrderParamsSortBy

const (
	// Deprecated: Use LimitOrdersQueryParamsSortByCreateDateTime instead.
	LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortByCreateDateTime = LimitOrdersQueryParamsSortByCreateDateTime
	// Deprecated: Use LimitOrdersQueryParamsSortByMakerAmount instead.
	LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortByMakerAmount = LimitOrdersQueryParamsSortByMakerAmount
	// Deprecated: Use LimitOrdersQueryParamsSortByMakerRate instead.
	LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortByMakerRate = LimitOrdersQueryParamsSortByMakerRate
	// Deprecated: Use LimitOrdersQueryParamsSortByTakerAmount instead.
	LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortByTakerAmount = LimitOrdersQueryParamsSortByTakerAmount
	// Deprecated: Use LimitOrdersQueryParamsSortByTakerRate instead.
	LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortByTakerRate = LimitOrdersQueryParamsSortByTakerRate
	// Deprecated: Use GetLimitOrderParamsSortByCreateDateTime instead.
	LimitOrderV3SubscribedApiControllerGetLimitOrderParamsSortByCreateDateTime = GetLimitOrderParamsSortByCreateDateTime
	// Deprecated: Use GetLimitOrderParamsSortByMakerAmount instead.
	LimitOrderV3SubscribedApiControllerGetLimitOrderParamsSortByMakerAmount = GetLimitOrderParamsSortByMakerAmount
	// Deprecated: Use GetLimitOrderParamsSortByMakerRate instead.
	LimitOrderV3SubscribedApiControllerGetLimitOrderParamsSortByMakerRate = GetLimitOrderParamsSortByMakerRate
	// Deprecated: Use GetLimitOrderParamsSortByTakerAmount instead.
	LimitOrderV3SubscribedApiControllerGetLimitOrderParamsSortByTakerAmount = GetLimitOrderParamsSortByTakerAmount
	// Deprecated: Use GetLimitOrderParamsSortByTakerRate instead.
	LimitOrderV3SubscribedApiControllerGetLimitOrderParamsSortByTakerRate = GetLimitOrderParamsSortByTakerRate
)
