package client

import (
	"context"
	"net/http"

	clienterrors "1inch-sdk-golang/client/errors"
	"1inch-sdk-golang/client/orderbook"
	"1inch-sdk-golang/client/swap"
	"1inch-sdk-golang/helpers"
)

type OrderbookService service

func (s *OrderbookService) CreateOrder(ctx context.Context, params orderbook.LimitOrderV3Request) (*orderbook.LimitOrderV3Data, *http.Response, error) {
	u := "/orderbook/v3.0/1/"

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	u, err = addQueryParameters(u, params)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("POST", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var allowanceResponse swap.AllowanceResponse
	res, err := s.client.Do(ctx, req, &allowanceResponse)
	if err != nil {
		return nil, nil, err
	}

	return nil, res, nil
}

// TODO Reusing the same request/response objects due to bad swagger spec
func (s *OrderbookService) GetOrdersByCreatorAddress(ctx context.Context, address string, params orderbook.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams) ([]*orderbook.OrderResponse, *http.Response, error) {
	u := "/orderbook/v3.0/1/address/{address}"

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
	u := "/orderbook/v3.0/1/all"

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
	u := "/orderbook/v3.0/1/count"

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
	u := "/orderbook/v3.0/1/events/{orderHash}"

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
	u := "/orderbook/v3.0/1/events"

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
	u := "/orderbook/v3.0/1/has-active-orders-with-permit/{walletAddress}/{token}"

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
