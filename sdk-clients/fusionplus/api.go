package fusionplus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/1inch/1inch-sdk-go/v5/common"
)

// Deprecated: Use GetOrderFillsByHash instead. The name now matches its return
// type (GetOrderFillsByHashOutput). This forwards to GetOrderFillsByHash.
func (api *api) GetOrderByOrderHash(ctx context.Context, params GetOrderByOrderHashParams) (*GetOrderFillsByHashOutput, error) {
	return api.GetOrderFillsByHash(ctx, params)
}

func (api *api) GetOrderFillsByHash(ctx context.Context, params GetOrderFillsByHashParams) (*GetOrderFillsByHashOutput, error) {
	u := fmt.Sprintf("/fusion-plus/orders/v1.1/order/status/%s", params.Hash)

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		Path:   u,
		Body:   nil,
	}

	var response GetOrderFillsByHashOutput
	err := api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (api *api) GetReadyToAcceptFills(ctx context.Context, params GetReadyToAcceptFillsParams) (*ReadyToAcceptSecretFills, error) {
	u := fmt.Sprintf("/fusion-plus/orders/v1.1/order/ready-to-accept-secret-fills/%s", params.Hash)

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		Path:   u,
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
	u := "/fusion-plus/relayer/v1.1/submit/secret"

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	payload := common.RequestPayload{
		Method: "POST",
		Params: params,
		Path:   u,
		Body:   body,
	}

	err = api.httpExecutor.ExecuteRequest(ctx, payload, nil)
	if err != nil {
		return err
	}

	return nil
}

func (api *api) GetQuote(ctx context.Context, params QuoteParams) (*Quote, error) {
	u := "/fusion-plus/quoter/v1.1/quote/receive"

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		Path:   u,
		Body:   nil,
	}

	var response Quote
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	// TODO must normalize response here

	return &response, nil
}

// PlaceOrder accepts a quote and submits it as a fusion plus order
func (api *api) PlaceOrder(ctx context.Context, quoteParams QuoteParams, quote *Quote, orderParams OrderParams, wallet common.Wallet) (string, error) {
	u := "/fusion-plus/relayer/v1.1/submit"

	err := orderParams.Validate()
	if err != nil {
		return "", err
	}

	preset, err := GetPreset(quote.Presets, orderParams.Preset)
	if err != nil {
		return "", fmt.Errorf("failed to get preset: %w", err)
	}

	// A multiple-fill order commits to one secret per fill part. The relayer
	// needs the secret hashes to rebuild the Merkle tree behind the hashlock.
	// Mirror the cross-chain-sdk contract: a preset that disallows multiple
	// fills cannot carry more than one hash, and a multiple-fill order must
	// carry exactly parts+1 hashes, where parts is the count packed into the
	// hashlock.
	if len(orderParams.SecretHashes) > 1 {
		if !preset.AllowMultipleFills {
			return "", fmt.Errorf("multiple secret hashes require a preset that allows multiple fills")
		}
		if orderParams.HashLock == nil {
			return "", fmt.Errorf("a multiple-fill order requires a hashlock")
		}
		if want := int(orderParams.HashLock.GetPartsCount()) + 1; len(orderParams.SecretHashes) != want {
			return "", fmt.Errorf("secret hashes length %d must equal the hashlock parts count plus one (%d)", len(orderParams.SecretHashes), want)
		}
	}

	fusionPlusOrder, err := CreateOrderData(quoteParams, quote, orderParams, wallet, quoteParams.SrcChain)
	if err != nil {
		return "", fmt.Errorf("failed to create order: %w", err)
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
		QuoteId:    quote.QuoteId,
		Signature:  fusionPlusOrder.LimitOrder.Signature,
		SrcChainId: quoteParams.SrcChain,
	}

	// The relayer needs the secret hashes only for a multiple-fill order. A
	// single-fill order omits them (its hashlock is the lone secret hash).
	if len(orderParams.SecretHashes) > 1 {
		signedOrder.SecretHashes = orderParams.SecretHashes
	}

	body, err := json.Marshal(signedOrder)
	if err != nil {
		return "", fmt.Errorf("failed to serialize order: %w", err)
	}

	payload := common.RequestPayload{
		Method: "POST",
		Params: nil,
		Path:   u,
		Body:   body,
	}

	err = api.httpExecutor.ExecuteRequest(ctx, payload, nil)
	if err != nil {
		return "", fmt.Errorf("failed to place order: %w", err)
	}

	return fusionPlusOrder.Hash, nil
}

// GetActiveOrders returns cross-chain orders that are currently open for filling
func (api *api) GetActiveOrders(ctx context.Context, params GetActiveOrdersParams) (*GetActiveOrdersOutput, error) {
	u := "/fusion-plus/orders/v1.1/order/active"

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		Path:   u,
		Body:   nil,
	}

	var response GetActiveOrdersOutput
	err := api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetSettlementContract returns the escrow factory contract address for a chain
func (api *api) GetSettlementContract(ctx context.Context, params GetSettlementContractParams) (*EscrowFactory, error) {
	u := "/fusion-plus/orders/v1.1/order/escrow"

	payload := common.RequestPayload{
		Method: "GET",
		Params: params,
		Path:   u,
		Body:   nil,
	}

	var response EscrowFactory
	err := api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
