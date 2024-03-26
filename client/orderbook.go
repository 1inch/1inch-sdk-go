package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk-go/client/models"
	"github.com/1inch/1inch-sdk-go/helpers"
	"github.com/1inch/1inch-sdk-go/helpers/consts/addresses"
	"github.com/1inch/1inch-sdk-go/helpers/consts/contracts"
	"github.com/1inch/1inch-sdk-go/internal/onchain"
	"github.com/1inch/1inch-sdk-go/internal/orderbook"
	"github.com/1inch/1inch-sdk-go/internal/tenderly"
)

type OrderbookService service

// CreateOrder creates an order in the Limit Order Protocol
func (s *OrderbookService) CreateOrder(ctx context.Context, params models.CreateOrderParams) (*models.CreateOrderResponse, *http.Response, error) {
	u := fmt.Sprintf("/orderbook/v4.0/%d", params.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	// Orders only last one minute if not specified in the request
	if params.ExpireAfter == 0 {
		params.ExpireAfter = time.Now().Add(time.Hour).Unix()
	}

	// To post an order that is open to anyone, the taker address must be the zero address
	if params.Taker == "" {
		params.Taker = addresses.Zero
	}

	aggregationRouter, err := contracts.Get1inchRouterFromChainId(params.ChainId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get 1inch router address: %v", err)
	}

	fromTokenAddress := common.HexToAddress(params.MakerAsset)
	publicAddress := common.HexToAddress(params.Maker)
	aggregationRouterAddress := common.HexToAddress(aggregationRouter)
	ethClient, err := s.client.GetEthClient(params.ChainId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get eth client: %v", err)
	}

	allowance, err := onchain.ReadContractAllowance(ethClient, fromTokenAddress, publicAddress, aggregationRouterAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read allowance: %v", err)
	}

	makingAmountBig, err := helpers.BigIntFromString(params.MakingAmount)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse making amount: %v", err)
	}
	if allowance.Cmp(makingAmountBig) <= 0 {

		if !params.EnableOnchainApprovalsIfNeeded {
			return nil, nil, models.ErrorFailWhenApprovalIsNeeded
		}

		if params.ShowSummaryBeforeExecution {
			ok, err := orderbook.ConfirmApprovalWithUser(ethClient, params.Maker, params.MakerAsset)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to confirm approval: %v", err)
			}
			if !ok {
				return nil, nil, errors.New("user rejected approval")
			}
		}

		// Only run the approval if Tenderly data is not present
		if _, ok := ctx.Value(tenderly.SwapConfigKey).(tenderly.SimulationConfig); !ok {
			erc20Config := onchain.Erc20ApprovalConfig{
				ChainId:        params.ChainId,
				Key:            params.PrivateKey,
				Erc20Address:   fromTokenAddress,
				PublicAddress:  publicAddress,
				SpenderAddress: aggregationRouterAddress,
			}
			err := onchain.ApproveTokenForRouter(ctx, ethClient, s.client.NonceCache, erc20Config)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to approve token for router: %v", err)
			}
			helpers.Sleep()
		}
	}

	seriesNonceManager, err := contracts.GetSeriesNonceManagerFromChainId(params.ChainId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get series nonce manager address: %v", err)
	}

	currentNonce, err := onchain.GetTimeSeriesManagerNonce(ethClient, seriesNonceManager, params.Maker)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get series nonce: %v", err)
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
		Nonce:              currentNonce.Int64(),
		Series:             0, // TODO: Series 0 always?
	}
	makerTraits := orderbook.BuildMakerTraits(buildMakerTraitsParams)

	order, err := orderbook.CreateLimitOrderMessage(params, makerTraits)
	if err != nil {
		return nil, nil, err
	}

	if params.ShowSummaryBeforeExecution {
		ok, err := orderbook.ConfirmLimitOrderWithUser(order, ethClient)
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

	marshalledIndented, err := json.MarshalIndent(order, "", "  ")
	fmt.Println(string(marshalledIndented))

	req, err := s.client.NewRequest("POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	var createOrderResponse models.CreateOrderResponse
	res, err := s.client.Do(ctx, req, &createOrderResponse)
	if err != nil {
		return nil, nil, err
	}

	return &createOrderResponse, res, nil
}

// TODO Reusing the same request/response objects due to bad swagger spec

// GetOrdersByCreatorAddress returns all orders created by a given address in the Limit Order Protocol
func (s *OrderbookService) GetOrdersByCreatorAddress(ctx context.Context, params models.GetOrdersByCreatorAddressParams) ([]models.OrderResponse, *http.Response, error) {
	u := fmt.Sprintf("/orderbook/v4.0/%d/address/%s", params.ChainId, params.CreatorAddress)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	u, err = addQueryParameters(u, params.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var ordersResponse []models.OrderResponse
	res, err := s.client.Do(ctx, req, &ordersResponse)
	if err != nil {
		return nil, nil, err
	}

	return ordersResponse, res, nil
}

// GetAllOrders returns all orders in the Limit Order Protocol
func (s *OrderbookService) GetAllOrders(ctx context.Context, params models.GetAllOrdersParams) ([]models.OrderResponse, *http.Response, error) {
	u := fmt.Sprintf("/orderbook/v4.0/%d/all", params.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	u, err = addQueryParameters(u, params.LimitOrderV3SubscribedApiControllerGetAllLimitOrdersParams)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var allOrdersResponse []models.OrderResponse
	res, err := s.client.Do(ctx, req, &allOrdersResponse)
	if err != nil {
		return nil, nil, err
	}

	return allOrdersResponse, res, nil
}

// GetCount returns the number of orders in the Limit Order Protocol
func (s *OrderbookService) GetCount(ctx context.Context, params models.GetCountParams) (*models.CountResponse, *http.Response, error) {
	u := fmt.Sprintf("/orderbook/v4.0/%d/count", params.ChainId)

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

	var count models.CountResponse
	res, err := s.client.Do(ctx, req, &count)
	if err != nil {
		return nil, nil, err
	}

	return &count, res, nil
}

// GetEvent returns an event in the Limit Order Protocol by order hash
func (s *OrderbookService) GetEvent(ctx context.Context, params models.GetEventParams) (*models.EventResponse, *http.Response, error) {
	u := fmt.Sprintf("/orderbook/v4.0/%d/events/%s", params.ChainId, params.OrderHash)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var event models.EventResponse
	res, err := s.client.Do(ctx, req, &event)
	if err != nil {
		return nil, nil, err
	}

	return &event, res, nil
}

// GetEvents returns all events in the Limit Order Protocol
func (s *OrderbookService) GetEvents(ctx context.Context, params models.GetEventsParams) ([]models.EventResponse, *http.Response, error) {
	u := fmt.Sprintf("/orderbook/v4.0/%d/events", params.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	u, err = addQueryParameters(u, params.LimitOrderV3SubscribedApiControllerGetEventsParams)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var events []models.EventResponse
	res, err := s.client.Do(ctx, req, &events)
	if err != nil {
		return nil, nil, err
	}

	return events, res, nil
}

// TODO untested endpoint

// GetActiveOrdersWithPermit returns all orders in the Limit Order Protocol that are active and have a valid permit
func (s *OrderbookService) GetActiveOrdersWithPermit(ctx context.Context, params models.GetActiveOrdersWithPermitParams) ([]models.OrderResponse, *http.Response, error) {
	u := fmt.Sprintf("/orderbook/v4.0/%d/has-active-orders-with-permit/%s/%s", params.ChainId, params.Token, params.Wallet)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var orders []models.OrderResponse
	res, err := s.client.Do(ctx, req, &orders)
	if err != nil {
		return nil, nil, err
	}

	return orders, res, nil
}
