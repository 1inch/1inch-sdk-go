package orderbook

// GetSortByParameter is a helper function that returns the pointer of the currency type being used
func GetSortByParameter(currency LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortBy) *LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortBy {
	return &currency
}
