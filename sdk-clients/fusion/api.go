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

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		U:      u,
		Body:   nil,
	}

	var response GetQuoteOutputFixed
	err := api.httpExecutor.ExecuteRequest(ctx, payload, &response)
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

// TODO Evaluate how to properly accept the order data and the extension, signature, and quoteId

func (api *api) PlaceOrder(ctx context.Context, quote GetQuoteOutputFixed, orderParams OrderParams) (*GetQuoteOutputFixed, error) {
	u := fmt.Sprintf("/fusion/relayer/v2.0/%d/order/submit", api.chainId)

	// TODO some kind of input validation
	fusionOrderParamsData := FusionOrderParamsData{
		NetworkId: 1,
		Preset:    Fast, // TODO currently always choosing the fast preset
		Receiver:  orderParams.Receiver,
		//Nonce:                   nil,
		//Permit:                  "",
		//IsPermit2:               false,
		//AllowPartialFills:       false,
		//AllowMultipleFills:      false,
		//DelayAuctionStartTimeBy: nil,
		//OrderExpirationDelay:    nil,
	}

	additionalParams := AdditionalParams{
		FromAddress: "0x0",
	}

	fusionOrder := CreateOrder(orderParams, quote, fusionOrderParamsData, additionalParams)

	fusionOrderIndented, err := json.MarshalIndent(fusionOrder, "", "  ")
	if err != nil {
		return nil, err
	}

	fmt.Printf("Fusion Order: %s\n", string(fusionOrderIndented))
	return nil, err

	// TODO implement the rest of the function
	//signedOrder := SignedOrderInput{
	//	Extension: "",
	//	Order: OrderInput{
	//		Maker:        quote.Maker,
	//		MakerAsset:   "",
	//		MakerTraits:  "",
	//		MakingAmount: "",
	//		Receiver:     "",
	//		Salt:         "",
	//		TakerAsset:   "",
	//		TakingAmount: "",
	//	},
	//	QuoteId:   "",
	//	Signature: "",
	//}

	payload := common.RequestPayload{
		Method: "GET",
		Params: nil,
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
