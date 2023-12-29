package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"1inch-sdk-golang/client/fusion"
)

type FusionService service

func (s *FusionService) GetOrders(ctx context.Context, params fusion.OrderApiControllerGetActiveOrdersParams) (*fusion.GetActiveOrdersOutput, *http.Response, error) {
	u := fmt.Sprintf("/fusion/orders/v1.0/%d/order/active", s.client.ChainId)

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

	var getOrdersResponse fusion.GetActiveOrdersOutput
	res, err := s.client.Do(ctx, req, &getOrdersResponse)
	if err != nil {
		return nil, nil, err
	}

	return &getOrdersResponse, res, nil
}

func (s *FusionService) GetSettlementContract(ctx context.Context) (*fusion.SettlementAddressOutput, *http.Response, error) {
	u := fmt.Sprintf("/fusion/orders/v1.0/%d/order/settlement", s.client.ChainId)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var getSettlementContractResponse fusion.SettlementAddressOutput
	res, err := s.client.Do(ctx, req, &getSettlementContractResponse)
	if err != nil {
		return nil, nil, err
	}

	return &getSettlementContractResponse, res, nil
}

func (s *FusionService) GetQuote(ctx context.Context, params fusion.QuoterControllerGetQuoteParams) (*fusion.GetQuoteOutput, *http.Response, error) {
	u := fmt.Sprintf("/fusion/quoter/v1.0/%d/quote/receive", s.client.ChainId)

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

	var getQuoteResponse fusion.GetQuoteOutput
	res, err := s.client.Do(ctx, req, &getQuoteResponse)
	if err != nil {
		return nil, nil, err
	}

	return &getQuoteResponse, res, nil
}

// TODO CustomPresetInput struct does not properly represent Points information

func (s *FusionService) GetQuoteWithCustomPresets(ctx context.Context, params fusion.QuoterControllerGetQuoteParams, body fusion.CustomPresetInput) (*fusion.GetQuoteOutput, *http.Response, error) {
	u := fmt.Sprintf("/fusion/quoter/v1.0/%d/quote/receive", s.client.ChainId)

	err := params.Validate()
	if err != nil {
		return nil, nil, err
	}

	err = body.Validate()
	if err != nil {
		return nil, nil, err
	}

	u, err = addQueryParameters(u, params)
	if err != nil {
		return nil, nil, err
	}

	bodyMarshalled, err := json.Marshal(body)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("POST", u, bodyMarshalled)
	if err != nil {
		return nil, nil, err
	}

	var getQuoteResponse fusion.GetQuoteOutput
	res, err := s.client.Do(ctx, req, &getQuoteResponse)
	if err != nil {
		return nil, nil, err
	}

	return &getQuoteResponse, res, nil
}
