package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/1inch/1inch-sdk/golang/client/onchain"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/contracts"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/tokens"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"

	clienterrors "github.com/1inch/1inch-sdk/golang/client/errors"
	"github.com/1inch/1inch-sdk/golang/client/orderbook"
	"github.com/1inch/1inch-sdk/golang/helpers"
)

type OrderbookService service

// CreateOrder creates an order in the Limit Order Protocol
func (s *OrderbookService) CreateOrder(ctx context.Context, params orderbook.OrderRequest) (*orderbook.CreateOrderResponse, *http.Response, error) {
	u := fmt.Sprintf("/orderbook/v3.0/%d", s.client.ChainId)

	if s.client.WalletKey == "" {
		return nil, nil, fmt.Errorf("wallet key must be set in the client config")
	}

	validate := validator.New()
	err := validate.Struct(params)
	if err != nil {
		return nil, nil, err
	}

	aggregationRouter, err := contracts.Get1inchRouterFromChainId(s.client.ChainId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get 1inch router address: %v", err)
	}

	if params.FromToken == tokens.NativeToken || params.ToToken == tokens.NativeToken {
		return nil, nil, errors.New("native gas token is not supported")
	}

	fromTokenAddress := common.HexToAddress(params.FromToken)
	publicAddress := common.HexToAddress(params.SourceWallet)
	aggreateRouterAddress := common.HexToAddress(aggregationRouter)
	allowance, err := onchain.ReadContractAllowance(s.client.EthClient, fromTokenAddress, publicAddress, aggreateRouterAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read allowance: %v", err)
	}

	makingAmountBig, err := helpers.BigIntFromString(params.MakingAmount)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse making amount: %v", err)
	}
	if allowance.Cmp(makingAmountBig) <= 0 {
		if !params.SkipWarnings {
			ok, err := orderbook.ConfirmApprovalWithUser(s.client.EthClient, params.SourceWallet, params.FromToken)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to confirm approval: %v", err)
			}
			if !ok {
				return nil, nil, errors.New("user rejected approval")
			}
		}
		err := onchain.ApproveTokenForRouter(s.client.EthClient, s.client.ChainId, s.client.WalletKey, fromTokenAddress, publicAddress, aggreateRouterAddress)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to approve token for router: %v", err)
		}
	}

	order, err := orderbook.CreateLimitOrder(params, s.client.ChainId, s.client.WalletKey)
	if err != nil {
		return nil, nil, err
	}

	if !params.SkipWarnings {
		ok, err := orderbook.ConfirmLimitOrderWithUser(order, s.client.EthClient)
		if err != nil {
			return nil, nil, err
		}
		if !ok {
			return nil, nil, errors.New("user rejected trade")
		}
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

// GetOrdersByCreatorAddress returns all orders created by a given address in the Limit Order Protocol
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

// GetAllOrders returns all orders in the Limit Order Protocol
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

// GetCount returns the number of orders in the Limit Order Protocol
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

// GetEvent returns an event in the Limit Order Protocol by order hash
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

// GetEvents returns all events in the Limit Order Protocol
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

// TODO untested endpoint

// GetActiveOrdersWithPermit returns all orders in the Limit Order Protocol that are active and have a valid permit
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
