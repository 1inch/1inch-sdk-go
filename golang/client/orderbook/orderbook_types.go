package orderbook

type GetOrdersByCreatorAddressParams struct {
	RequestParams
	LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams
}

type GetAllOrdersParams struct {
	RequestParams
	LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams
}

type GetCountParams struct {
	RequestParams
	LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams
}

type GetEventParams struct {
	RequestParams
	OrderHash string
}

type GetEventsParams struct {
	RequestParams
	LimitOrderV3SubscribedApiControllerGetEventsParams
}

type GetActiveOrdersWithPermitParams struct {
	RequestParams
	Wallet string
	Token  string
}
