package swap

type SwapTokensParams struct {
	ApprovalType  ApprovalType
	ChainId       int
	SkipWarnings  bool
	PublicAddress string
	WalletKey     string
	AggregationControllerGetSwapParams
}

type ApproveAllowanceParams struct {
	ChainId int
	ApproveControllerGetAllowanceParams
}

type ApproveSpenderParams struct {
	ChainId int
}

type ApproveTransactionParams struct {
	ChainId int
	ApproveControllerGetCallDataParams
}

type GetLiquiditySourcesParams struct {
	ChainId int
}

type GetQuoteParams struct {
	ChainId int
	AggregationControllerGetQuoteParams
}

type GetSwapDataParams struct {
	ChainId      int
	SkipWarnings bool
	AggregationControllerGetSwapParams
}

type GetTokensParams struct {
	ChainId int
}
