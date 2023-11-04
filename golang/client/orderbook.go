package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	clienterrors "1inch-sdk-golang/client/errors"
	"1inch-sdk-golang/client/orderbook"
	"1inch-sdk-golang/helpers"
)

type OrderbookService service

func (s *OrderbookService) CreateOrder(ctx context.Context, params orderbook.OrderRequest) (*orderbook.CreateOrderResponse, *http.Response, error) {
	u := fmt.Sprintf("/orderbook/v3.0/%d", s.client.ChainId)

	validate := validator.New()
	err := validate.Struct(params)
	if err != nil {
		return nil, nil, err
	}

	order, err := orderbook.CreateLimitOrder(params, 137, s.client.WalletKey)
	if err != nil {
		return nil, nil, err
	}

	body, err := json.Marshal(order)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	var createOrderResponse orderbook.CreateOrderResponse
	res, err := s.client.Do(ctx, req, &createOrderResponse)
	if err != nil {
		return nil, nil, err
	}

	return &createOrderResponse, res, nil
}

// TODO Reusing the same request/response objects due to bad swagger spec
func (s *OrderbookService) GetOrdersByCreatorAddress(ctx context.Context, address string, params orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams) ([]*orderbook.OrderResponse, *http.Response, error) {
	u := fmt.Sprintf("/orderbook/v3.0/%d/address/{address}", s.client.ChainId)

	if !helpers.IsEthereumAddress(address) {
		return nil, nil, clienterrors.NewRequestValidationError("address must be a valid Ethereum address")
	}

	u, err := ReplacePathVariable(u, "address", address)
	if err != nil {
		return nil, nil, err
	}

	err = params.Validate()
	if err != nil {
		return nil, nil, err
	}

	u, err = addQueryParameters(u, params)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var ordersResponse []*orderbook.OrderResponse
	res, err := s.client.Do(ctx, req, &ordersResponse)
	if err != nil {
		return nil, nil, err
	}

	return ordersResponse, res, nil
}

func (s *OrderbookService) GetAllOrders(ctx context.Context, params orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams) ([]*orderbook.OrderResponse, *http.Response, error) {
	u := fmt.Sprintf("/orderbook/v3.0/%d/all", s.client.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	u, err = addQueryParameters(u, params)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var allOrdersResponse []*orderbook.OrderResponse
	res, err := s.client.Do(ctx, req, &allOrdersResponse)
	if err != nil {
		return nil, nil, err
	}

	return allOrdersResponse, res, nil
}

func (s *OrderbookService) GetCount(ctx context.Context, params orderbook.LimitOrderV3SubscribedApiControllerGetAllOrdersCountParams) (*orderbook.CountResponse, *http.Response, error) {
	u := fmt.Sprintf("/orderbook/v3.0/%d/count", s.client.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	u, err = addQueryParameters(u, params)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var count orderbook.CountResponse
	res, err := s.client.Do(ctx, req, &count)
	if err != nil {
		return nil, nil, err
	}

	return &count, res, nil
}

func (s *OrderbookService) GetEvent(ctx context.Context, orderHash string) (*orderbook.EventResponse, *http.Response, error) {
	u := fmt.Sprintf("/orderbook/v3.0/%d/events/{orderHash}", s.client.ChainId)

	u, err := ReplacePathVariable(u, "orderHash", orderHash)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var event orderbook.EventResponse
	res, err := s.client.Do(ctx, req, &event)
	if err != nil {
		return nil, nil, err
	}

	return &event, res, nil
}

func (s *OrderbookService) GetEvents(ctx context.Context, params orderbook.LimitOrderV3SubscribedApiControllerGetEventsParams) ([]*orderbook.EventResponse, *http.Response, error) {
	u := fmt.Sprintf("/orderbook/v3.0/%d/events", s.client.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	u, err = addQueryParameters(u, params)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []*orderbook.EventResponse
	res, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, nil, err
	}

	return events, res, nil
}

// TODO need docs
func (s *OrderbookService) GetActiveOrdersWithPermit(ctx context.Context, wallet string, token string) ([]*orderbook.OrderResponse, *http.Response, error) {
	u := fmt.Sprintf("/orderbook/v3.0/%d/has-active-orders-with-permit/{walletAddress}/{token}", s.client.ChainId)

	if !helpers.IsEthereumAddress(wallet) {
		return nil, nil, clienterrors.NewRequestValidationError("wallet must be a valid Ethereum address")
	}
	u, err := ReplacePathVariable(u, "walletAddress", wallet)
	if err != nil {
		return nil, nil, err
	}

	if !helpers.IsEthereumAddress(token) {
		return nil, nil, clienterrors.NewRequestValidationError("token must be a valid Ethereum address")
	}
	u, err = ReplacePathVariable(u, "token", token)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var orders []*orderbook.OrderResponse
	res, err := s.client.Do(ctx, req, &orders)
	if err != nil {
		return nil, nil, err
	}

	return orders, res, nil
}
