package fusion

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/1inch/1inch-sdk-go/common"
)

func (api *api) GetActiveOrders(ctx context.Context, params OrderApiControllerGetActiveOrdersParams) (*GetActiveOrdersOutput, error) {
	u := fmt.Sprintf("/fusion/orders/v2.0/%d/order/active", api.chainId)

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		U:      u,
		Body:   nil,
	}

	var response GetActiveOrdersOutput
	err := api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetSettlementContract(ctx context.Context) (*SettlementAddressOutput, error) {
	u := fmt.Sprintf("/fusion/orders/v2.0/%d/order/settlement", api.chainId)

	payload := common.RequestPayload{
		Method: "GET",
		Params: nil,
		U:      u,
		Body:   nil,
	}

	var response SettlementAddressOutput
	err := api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetQuote(ctx context.Context, params QuoterControllerGetQuoteParams) (*GetQuoteOutputFixed, error) {
	u := fmt.Sprintf("/fusion/quoter/v2.0/%d/quote/receive", api.chainId)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		U:      u,
		Body:   nil,
	}

	var response GetQuoteOutputFixed
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetQuoteWithCustomPreset(ctx context.Context, params QuoterControllerGetQuoteWithCustomPresetsParams, presetDetails QuoterControllerGetQuoteWithCustomPresetsJSONRequestBody) (*GetQuoteOutputFixed, error) {
	u := fmt.Sprintf("/fusion/quoter/v2.0/%d/quote/receive", api.chainId)

	body, err := json.Marshal(presetDetails)
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		U:      u,
		Body:   body,
	}

	var response GetQuoteOutputFixed
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// PlaceOrder accepts a quote and submits it as a fusion order
func (api *api) PlaceOrder(ctx context.Context, fusionQuote GetQuoteOutputFixed, orderParams OrderParams, wallet common.Wallet) error {
	u := fmt.Sprintf("/fusion/relayer/v2.0/%d/order/submit", api.chainId)

	err := orderParams.Validate()
	if err != nil {
		return err
	}

	// TODO This function can simply return the SignedOrderInput object
	_, limitOrder, err := CreateFusionOrderData(fusionQuote, orderParams, wallet, api.chainId)
	if err != nil {
		return fmt.Errorf("failed to create order: %v", err)
	}

	signedOrder := SignedOrderInput{
		Extension: limitOrder.Data.Extension,
		Order: OrderInput{
			Maker:        limitOrder.Data.Maker,
			MakerAsset:   limitOrder.Data.MakerAsset,
			MakerTraits:  limitOrder.Data.MakerTraits,
			MakingAmount: limitOrder.Data.MakingAmount,
			Receiver:     limitOrder.Data.Receiver,
			Salt:         limitOrder.Data.Salt,
			TakerAsset:   limitOrder.Data.TakerAsset,
			TakingAmount: limitOrder.Data.TakingAmount,
		},
		QuoteId:   fusionQuote.QuoteId,
		Signature: limitOrder.Signature,
	}

	body, err := json.Marshal(signedOrder)
	if err != nil {
		return err
	}

	payload := common.RequestPayload{
		Method: "POST",
		Params: nil,
		U:      u,
		Body:   body,
	}

	err = api.httpExecutor.ExecuteRequest(ctx, payload, nil)
	if err != nil {
		return err
	}

	return nil
}

func (api *api) PlaceOrders(ctx context.Context, body []PlaceOrderBody) (*GetQuoteOutput, error) {
	u := fmt.Sprintf("/fusion/relayer/v2.0/%d/order/submit/many", api.chainId)

	for _, order := range body {
		err := order.Validate()
		if err != nil {
			return nil, err
		}
	}

	bodyMarshaled, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: nil,
		U:      u,
		Body:   bodyMarshaled,
	}

	var response GetQuoteOutput
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
