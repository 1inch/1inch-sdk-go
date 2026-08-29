package portfolio

var (
	_ GetPortfolioValueResponse         = GetProtocolsCurrentValueResponse{}
	_ GetPortfolioProfitAndLossResponse = GetProtocolsProfitAndLossResponse{}
	_ GetTokensProfitLossResponse       = GetTokensProfitAndLossResponse{}
	_ GetCurrentProfitLossResponse      = GetProfitAndLossResponse{}
)
