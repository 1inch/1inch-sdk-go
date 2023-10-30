// Package orderbook provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package orderbook

// Defines values for LimitOrderV3SubscribedApiControllerGetLimitOrderParamsSortBy.
const (
	LimitOrderV3SubscribedApiControllerGetLimitOrderParamsSortByCreateDateTime LimitOrderV3SubscribedApiControllerGetLimitOrderParamsSortBy = "createDateTime"
	LimitOrderV3SubscribedApiControllerGetLimitOrderParamsSortByMakerAmount    LimitOrderV3SubscribedApiControllerGetLimitOrderParamsSortBy = "makerAmount"
	LimitOrderV3SubscribedApiControllerGetLimitOrderParamsSortByMakerRate      LimitOrderV3SubscribedApiControllerGetLimitOrderParamsSortBy = "makerRate"
	LimitOrderV3SubscribedApiControllerGetLimitOrderParamsSortByTakerAmount    LimitOrderV3SubscribedApiControllerGetLimitOrderParamsSortBy = "takerAmount"
	LimitOrderV3SubscribedApiControllerGetLimitOrderParamsSortByTakerRate      LimitOrderV3SubscribedApiControllerGetLimitOrderParamsSortBy = "takerRate"
)

// Defines values for LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortBy.
const (
	LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortByCreateDateTime LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortBy = "createDateTime"
	LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortByMakerAmount    LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortBy = "makerAmount"
	LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortByMakerRate      LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortBy = "makerRate"
	LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortByTakerAmount    LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortBy = "takerAmount"
	LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortByTakerRate      LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortBy = "takerRate"
)

// LimitOrderV3Data defines model for LimitOrderV3Data.
type LimitOrderV3Data struct {
	// AllowedSender An address of the wallet or contract who will be able to fill this order (equals to Zero address on public orders)
	AllowedSender *string `json:"allowedSender,omitempty"`

	// Interactions Merged offsets of each field in interactions
	Interactions *string `json:"interactions,omitempty"`

	// Maker An address of the maker (wallet or contract address)
	Maker string `json:"maker"`

	// MakerAsset Address of the maker asset
	MakerAsset string `json:"makerAsset"`

	// MakingAmount Order maker's token amount
	MakingAmount string `json:"makingAmount"`

	// Offsets An interaction call data. ABI encoded set of makerAssetData, takerAssetData, getMakingAmount, getTakingAmount, predicate, permit, preInteraction, postInteraction
	Offsets *string `json:"offsets,omitempty"`

	// Receiver An address of the wallet or contract who will receive filled amount (equals to Zero address for receiver == makerAddress)
	Receiver *string `json:"receiver,omitempty"`

	// Salt Some unique value. It is necessary to be able to create limit orders with the same parameters (so that they have a different hash)
	Salt string `json:"salt"`

	// TakerAsset Address of the taker asset
	TakerAsset string `json:"takerAsset"`

	// TakingAmount Order taker's token amount
	TakingAmount string `json:"takingAmount"`
}

// LimitOrderV3Request defines model for LimitOrderV3Request.
type LimitOrderV3Request struct {
	// Data Limit order data
	Data LimitOrderV3Data `json:"data"`

	// OrderHash Hash of the limit order typed data
	OrderHash string `json:"orderHash"`

	// Signature Signature of the limit order typed data (using signTypedData_v4)
	Signature string `json:"signature"`
}

// LimitOrderV3SubscribedApiControllerGetLimitOrderParams defines parameters for LimitOrderV3SubscribedApiControllerGetLimitOrder.
type LimitOrderV3SubscribedApiControllerGetLimitOrderParams struct {
	// Page Pagination step, default: 1 (page = offset / limit)
	Page *float32 `url:"page,omitempty" json:"page,omitempty"`

	// Limit Number of limit orders to receive (default: 100, max: 500)
	Limit *float32 `url:"limit,omitempty" json:"limit,omitempty"`

	// Statuses JSON an array of statuses by which limit orders will be filtered: 1 - valid limit orders, 2 - temporary invalid limit orders, 3 - invalid limit orders
	Statuses *[]float32                                                    `url:"statuses,omitempty" json:"statuses,omitempty"`
	SortBy   *LimitOrderV3SubscribedApiControllerGetLimitOrderParamsSortBy `url:"sortBy,omitempty" json:"sortBy,omitempty"`

	// TakerAsset Address of the taker asset
	TakerAsset *string `url:"takerAsset,omitempty" json:"takerAsset,omitempty"`

	// MakerAsset Address of the maker asset
	MakerAsset *string `url:"makerAsset,omitempty" json:"makerAsset,omitempty"`
}

// LimitOrderV3SubscribedApiControllerGetLimitOrderParamsSortBy defines parameters for LimitOrderV3SubscribedApiControllerGetLimitOrder.
type LimitOrderV3SubscribedApiControllerGetLimitOrderParamsSortBy string

// LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams defines parameters for LimitOrderV3SubscribedApiControllerGetAllLimitOrders.
type LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams struct {
	// Page Pagination step, default: 1 (page = offset / limit)
	Page *float32 `url:"page,omitempty" json:"page,omitempty"`

	// Limit Number of limit orders to receive (default: 100, max: 500)
	Limit *float32 `url:"limit,omitempty" json:"limit,omitempty"`

	// Statuses JSON an array of statuses by which limit orders will be filtered: 1 - valid limit orders, 2 - temporary invalid limit orders, 3 - invalid limit orders
	Statuses *[]float32                                                        `url:"statuses,omitempty" json:"statuses,omitempty"`
	SortBy   *LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortBy `url:"sortBy,omitempty" json:"sortBy,omitempty"`

	// TakerAsset Address of the maker asset
	TakerAsset *string `url:"takerAsset,omitempty" json:"takerAsset,omitempty"`

	// MakerAsset Address of the maker asset
	MakerAsset *string `url:"makerAsset,omitempty" json:"makerAsset,omitempty"`
}

// LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortBy defines parameters for LimitOrderV3SubscribedApiControllerGetAllLimitOrders.
type LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParamsSortBy string

// LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams defines parameters for LimitOrderV3SubscribedApiControllerGetAllOrdersCount.
type LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams struct {
	Statuses []string `url:"statuses" json:"statuses"`
}

// LimitOrderV3SubscribedApiControllerGetEventsParams defines parameters for LimitOrderV3SubscribedApiControllerGetEvents.
type LimitOrderV3SubscribedApiControllerGetEventsParams struct {
	Limit float32 `url:"limit" json:"limit"`
}

// LimitOrderV3SubscribedApiControllerCreateLimitOrderJSONRequestBody defines body for LimitOrderV3SubscribedApiControllerCreateLimitOrder for application/json ContentType.
type LimitOrderV3SubscribedApiControllerCreateLimitOrderJSONRequestBody = LimitOrderV3Request
