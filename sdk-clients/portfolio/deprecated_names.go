package portfolio

// Pre-v5 names for generated types. The generated names now come from the
// friendly operation-id mapping (codegen/mapping.json) and schema name
// overrides (codegen/overrides.go); these aliases keep existing integrations
// compiling and are identical types.

// Deprecated: Use GetCurrentValueParams instead.
type GetCurrentValuePortfolioV4GeneralCurrentValueGetParams = GetCurrentValueParams

// Deprecated: Use GetTokensCurrentValueParams instead.
type GetCurrentValuePortfolioV4OverviewErc20CurrentValueGetParams = GetTokensCurrentValueParams

// Deprecated: Use GetProtocolsCurrentValueParams instead.
type GetCurrentValuePortfolioV4OverviewProtocolsCurrentValueGetParams = GetProtocolsCurrentValueParams

// Deprecated: Use GetTokensDetailsParams instead.
type GetDetailsPortfolioV4OverviewErc20DetailsGetParams = GetTokensDetailsParams

// Deprecated: Use GetProtocolsDetailsParams instead.
type GetDetailsPortfolioV4OverviewProtocolsDetailsGetParams = GetProtocolsDetailsParams

// Deprecated: Use GetProfitAndLossParams instead.
type GetProfitAndLossPortfolioV4GeneralProfitAndLossGetParams = GetProfitAndLossParams

// Deprecated: Use GetTokensProfitAndLossParams instead.
type GetProfitAndLossPortfolioV4OverviewErc20ProfitAndLossGetParams = GetTokensProfitAndLossParams

// Deprecated: Use GetProtocolsProfitAndLossParams instead.
type GetProfitAndLossPortfolioV4OverviewProtocolsProfitAndLossGetParams = GetProtocolsProfitAndLossParams

// Deprecated: Use GetValueChartParams instead.
type GetValueChartPortfolioV4GeneralValueChartGetParams = GetValueChartParams
