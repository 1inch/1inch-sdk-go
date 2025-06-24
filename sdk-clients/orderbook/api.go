package orderbook

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/1inch/1inch-sdk-go/common"
)

// CreateOrder creates an order in the Limit Order Protocol
func (api *api) CreateOrder(ctx context.Context, params CreateOrderParams) (*CreateOrderResponse, error) {
	u := fmt.Sprintf("/orderbook/v4.0/%d", api.chainId)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	if !params.MakerTraits.AllowMultipleFills || !params.MakerTraits.AllowPartialFills {
		return nil, errors.New("allowMultipleFills and allowPartialFills must be true")
	}

	//order, err := CreateLimitOrderMessage(params, int(api.chainId))
	//if err != nil {
	//	return nil, err
	//}

	body, err := json.Marshal(nil)
	//body, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "POST",
		Params: nil,
		U:      u,
		Body:   body,
	}

	var createOrderResponse CreateOrderResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &createOrderResponse)
	if err != nil {
		return nil, err
	}

	return &createOrderResponse, nil
}

// TODO Reusing the same request/response objects until the openapi spec is updated to include the correct object definitions

// GetOrdersByCreatorAddress returns all orders created by a given address in the Limit Order Protocol
func (api *api) GetOrdersByCreatorAddress(ctx context.Context, params GetOrdersByCreatorAddressParams) ([]*OrderResponse, error) {
	u := fmt.Sprintf("/orderbook/v4.0/%d/address/%s", api.chainId, params.CreatorAddress)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		U:      u,
	}

	var ordersResponse []*OrderResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &ordersResponse)
	if err != nil {
		return nil, err
	}

	return ordersResponse, nil
}

// GetOrder returns an order from Limit Order Protocol that matches a specific hash
func (api *api) GetOrder(ctx context.Context, params GetOrderParams) (*GetOrderByHashResponseExtended, error) {
	u := fmt.Sprintf("/orderbook/v4.0/%d/order/%s", api.chainId, params.OrderHash)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		U:      u,
	}

	var getOrderByHashResponse *GetOrderByHashResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &getOrderByHashResponse)
	if err != nil {
		return nil, err
	}

	return NormalizeGetOrderByHashResponse(getOrderByHashResponse)
}

// GetOrderWithSignature first looks up an order by hash, then does a second request to get the signature data
func (api *api) GetOrderWithSignature(ctx context.Context, params GetOrderParams) (*OrderExtendedWithSignature, error) {

	// First lookup the order by hash (no signature on this response)
	order, err := api.GetOrder(ctx, params)
	if err != nil {
		return nil, err
	}

	// For free accounts, this sleep is required to avoid 429 errors
	if params.SleepBetweenSubrequests {
		time.Sleep(time.Second)
	}

	// Second, lookup all orders by that creator (these orders will contain the signature data)
	allOrdersByCreator, err := api.GetOrdersByCreatorAddress(ctx, GetOrdersByCreatorAddressParams{
		CreatorAddress: order.OrderMaker,
	})
	if err != nil {
		return nil, err
	}

	// Filter through the second set of orders to find the signature
	for _, o := range allOrdersByCreator {
		if o.OrderHash == params.OrderHash {
			return &OrderExtendedWithSignature{
				GetOrderByHashResponse:   order.GetOrderByHashResponse,
				LimitOrderDataNormalized: order.LimitOrderDataNormalized,
				Signature:                o.Signature,
			}, nil
		}
	}

	return nil, errors.New("order not found")
}

// GetAllOrders returns all orders in the Limit Order Protocol
func (api *api) GetAllOrders(ctx context.Context, params GetAllOrdersParams) ([]OrderResponse, error) {
	u := fmt.Sprintf("/orderbook/v3.0/%d/all", api.chainId)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		U:      u,
	}

	var allOrdersResponse []OrderResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &allOrdersResponse)
	if err != nil {
		return nil, err
	}

	return allOrdersResponse, nil
}

// GetCount returns the number of orders in the Limit Order Protocol
func (api *api) GetCount(ctx context.Context, params GetCountParams) (*CountResponse, error) {
	u := fmt.Sprintf("/orderbook/v3.0/%d/count", api.chainId)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		U:      u,
	}

	var count CountResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &count)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

// GetEvent returns an event in the Limit Order Protocol by order hash
func (api *api) GetEvent(ctx context.Context, params GetEventParams) (*EventResponse, error) {
	u := fmt.Sprintf("/orderbook/v3.0/%d/events/%s", api.chainId, params.OrderHash)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		U:      u,
	}

	var event EventResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// GetEvents returns all events in the Limit Order Protocol
func (api *api) GetEvents(ctx context.Context, params GetEventsParams) ([]EventResponse, error) {
	u := fmt.Sprintf("/orderbook/v3.0/%d/events", api.chainId)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		U:      u,
	}

	var events []EventResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// TODO untested endpoint

// GetActiveOrdersWithPermit returns all orders in the Limit Order Protocol that are active and have a valid permit
func (api *api) GetActiveOrdersWithPermit(ctx context.Context, params GetActiveOrdersWithPermitParams) ([]OrderResponse, error) {
	u := fmt.Sprintf("/orderbook/v3.0/%d/has-active-orders-with-permit/%s/%s", api.chainId, params.Token, params.Wallet)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		U:      u,
	}

	var orders []OrderResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
