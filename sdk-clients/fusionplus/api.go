package fusionplus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/1inch/1inch-sdk-go/common"
)

func (api *api) GetOrderByOrderHash(ctx context.Context, params GetOrderByOrderHashParams) (*GetOrderFillsByHashOutputFixed, error) {
	u := fmt.Sprintf("/fusion-plus/orders/v1.0/order/status/%s", params.Hash)

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		U:      u,
		Body:   nil,
	}

	var response GetOrderFillsByHashOutputFixed
	err := api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetReadyToAcceptFills(ctx context.Context, params GetOrderByOrderHashParams) (*ReadyToAcceptSecretFills, error) {
	u := fmt.Sprintf("/fusion-plus/orders/v1.0/order/ready-to-accept-secret-fills/%s", params.Hash)

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		U:      u,
		Body:   nil,
	}

	var response ReadyToAcceptSecretFills
	err := api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) SubmitSecret(ctx context.Context, params SecretInput) error {
	u := "/fusion-plus/relayer/v1.0/submit/secret"

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	payload := common.RequestPayload{
		Method: "POST",
		Params: params,
		U:      u,
		Body:   body,
	}

	err = api.httpExecutor.ExecuteRequest(ctx, payload, nil)
	if err != nil {
		return err
	}

	return nil
}

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

func (api *api) GetQuote(ctx context.Context, params QuoterControllerGetQuoteParamsFixed) (*GetQuoteOutputFixed, error) {
	u := "/fusion-plus/quoter/v1.0/quote/receive"

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

	// TODO must normalize response here

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

// PlaceOrder accepts a quote and submits it as a fusion plus order
func (api *api) PlaceOrder(ctx context.Context, quoteParams QuoterControllerGetQuoteParamsFixed, quote *GetQuoteOutputFixed, orderParams OrderParams, wallet common.Wallet) (string, error) {
	u := "/fusion-plus/relayer/v1.0/submit"

	err := orderParams.Validate()
	if err != nil {
		return "", err
	}

	preset, err := GetPreset(quote.Presets, orderParams.Preset)
	if err != nil {
		return "", fmt.Errorf("failed to get preset: %v", err)
	}

	// TODO orders will now be allowed with multiple secret hashes for now
	if len(orderParams.SecretHashes) > 1 {
		return "", fmt.Errorf("Multiple secret hashes are not supported at this time. Please ")
	}

	if !preset.AllowMultipleFills && len(orderParams.SecretHashes) > 1 {
		return "", fmt.Errorf("multiple secrets are required with multiple secret hashes")
	} else {
		// TODO support multiple secrets
	}

	fusionPlusOrder, err := CreateFusionPlusOrderData(quoteParams, quote, orderParams, wallet, int(quoteParams.SrcChain))
	if err != nil {
		return "", fmt.Errorf("failed to create order: %v", err)
	}

	signedOrder := SignedOrderInput{
		Extension: fusionPlusOrder.LimitOrder.Data.Extension,
		Order: OrderInput{
			Maker:        fusionPlusOrder.LimitOrder.Data.Maker,
			MakerAsset:   fusionPlusOrder.LimitOrder.Data.MakerAsset,
			MakerTraits:  fusionPlusOrder.LimitOrder.Data.MakerTraits,
			MakingAmount: fusionPlusOrder.LimitOrder.Data.MakingAmount,
			Receiver:     fusionPlusOrder.LimitOrder.Data.Receiver,
			Salt:         fusionPlusOrder.LimitOrder.Data.Salt,
			TakerAsset:   fusionPlusOrder.LimitOrder.Data.TakerAsset,
			TakingAmount: fusionPlusOrder.LimitOrder.Data.TakingAmount,
		},
		QuoteId: quote.QuoteId,
		//SecretHashes: orderParams.SecretHashes, // TODO this only should be submitted when there are multiple secrets
		Signature:  fusionPlusOrder.LimitOrder.Signature,
		SrcChainId: quoteParams.SrcChain,
	}

	body, err := json.Marshal(signedOrder)
	if err != nil {
		return "", err
	}

	payload := common.RequestPayload{
		Method: "POST",
		Params: nil,
		U:      u,
		Body:   body,
	}

	err = api.httpExecutor.ExecuteRequest(ctx, payload, nil)
	if err != nil {
		return "", err
	}

	return fusionPlusOrder.Hash, nil
}
