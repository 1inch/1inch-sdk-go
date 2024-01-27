package swap

type RequestParams struct {
	ChainId       int
	SkipWarnings  bool
	PublicAddress string
	WalletKey     string
}

type SwapTokensParams struct {
	ApprovalType ApprovalType
	RequestParams
	AggregationControllerGetSwapParams
}

type ApproveAllowanceParams struct {
	RequestParams
	ApproveControllerGetAllowanceParams
}

type ApproveSpenderParams struct {
	RequestParams
}

type ApproveTransactionParams struct {
	RequestParams
	ApproveControllerGetCallDataParams
}

type GetLiquiditySourcesParams struct {
	RequestParams
}

type GetQuoteParams struct {
	RequestParams
	AggregationControllerGetQuoteParams
}

type GetSwapDataParams struct {
	RequestParams
	AggregationControllerGetSwapParams
}

type GetTokensParams struct {
	RequestParams
}
