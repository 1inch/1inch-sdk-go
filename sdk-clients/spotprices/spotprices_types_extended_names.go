package spotprices

// GetCustomCurrenciesResponse is the clean name (dropping the Dto suffix) for
// the GetCustomCurrenciesList return; the generated CurrenciesResponseDto
// remains valid. (GetPricesRequestDto is intentionally not aliased here — the
// name GetPricesForRequestedTokensParams already exists as a generated type.)
type GetCustomCurrenciesResponse = CurrenciesResponseDto
