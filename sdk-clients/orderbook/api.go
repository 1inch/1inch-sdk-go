package orderbook

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/1inch/1inch-sdk-go/common"
	"github.com/1inch/1inch-sdk-go/internal/orderbook"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook/models"
)

const zeroAddress = "0x0000000000000000000000000000000000000000"

// CreateOrder creates an order in the Limit Order Protocol
func (api *api) CreateOrder(ctx context.Context, params models.CreateOrderParams) (*models.CreateOrderResponse, error) {
	u := fmt.Sprintf("/orderbook/v4.0/%d", api.chainId)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	// Orders only last one minute if not specified in the request
	if params.ExpireAfter == 0 {
		params.ExpireAfter = time.Now().Add(time.Hour).Unix()
	}

	// To post an order that is open to anyone, the taker address must be the zero address
	if params.Taker == "" {
		params.Taker = zeroAddress
	}

	buildMakerTraitsParams := models.BuildMakerTraitsParams{
		AllowedSender:      params.Taker,
		ShouldCheckEpoch:   false,
		UsePermit2:         false,
		UnwrapWeth:         false,
		HasExtension:       false,
		HasPreInteraction:  false,
		HasPostInteraction: false,
		Expiry:             params.ExpireAfter,
		Nonce:              params.SeriesNonce.Int64(),
		Series:             0, // TODO: Series 0 always?
	}
	makerTraits := orderbook.BuildMakerTraits(buildMakerTraitsParams)

	order, err := orderbook.CreateLimitOrderMessage(params, makerTraits)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "POST",
		Params: nil,
		U:      u,
		Body:   body,
	}

	var createOrderResponse models.CreateOrderResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &createOrderResponse)
	if err != nil {
		return nil, err
	}

	return &createOrderResponse, nil
}

// TODO Reusing the same request/response objects due to bad openapi spec

// GetOrdersByCreatorAddress returns all orders created by a given address in the Limit Order Protocol
func (api *api) GetOrdersByCreatorAddress(ctx context.Context, params models.GetOrdersByCreatorAddressParams) ([]models.OrderResponse, error) {
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

	var ordersResponse []models.OrderResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &ordersResponse)
	if err != nil {
		return nil, err
	}

	return ordersResponse, nil
}

// GetAllOrders returns all orders in the Limit Order Protocol
func (api *api) GetAllOrders(ctx context.Context, params models.GetAllOrdersParams) ([]models.OrderResponse, error) {
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

	var allOrdersResponse []models.OrderResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &allOrdersResponse)
	if err != nil {
		return nil, err
	}

	return allOrdersResponse, nil
}

// GetCount returns the number of orders in the Limit Order Protocol
func (api *api) GetCount(ctx context.Context, params models.GetCountParams) (*models.CountResponse, error) {
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

	var count models.CountResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &count)
	if err != nil {
		return nil, err
	}

	return &count, nil
}

// GetEvent returns an event in the Limit Order Protocol by order hash
func (api *api) GetEvent(ctx context.Context, params models.GetEventParams) (*models.EventResponse, error) {
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

	var event models.EventResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// GetEvents returns all events in the Limit Order Protocol
func (api *api) GetEvents(ctx context.Context, params models.GetEventsParams) ([]models.EventResponse, error) {
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

	var events []models.EventResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// TODO untested endpoint

// GetActiveOrdersWithPermit returns all orders in the Limit Order Protocol that are active and have a valid permit
func (api *api) GetActiveOrdersWithPermit(ctx context.Context, params models.GetActiveOrdersWithPermitParams) ([]models.OrderResponse, error) {
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

	var orders []models.OrderResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
