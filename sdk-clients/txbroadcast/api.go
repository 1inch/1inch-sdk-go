package txbroadcast

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/1inch/1inch-sdk-go/common"
)

// BroadcastPublicTransaction Do broadcast EVM transaction to the public mempool, takes the raw transaction data as input and returns the transaction hash
func (api *api) BroadcastPublicTransaction(ctx context.Context, params BroadcastRequest) (*BroadcastResponse, error) {
	u := fmt.Sprintf("tx-gateway/v1.1/%d/broadcast", api.chainId)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "POST",
		Params: nil,
		U:      u,
		Body:   body,
	}

	var response BroadcastResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// BroadcastPrivateTransaction Do broadcast EVM transaction to the private mempool, takes the raw transaction data as input and returns the transaction hash
func (api *api) BroadcastPrivateTransaction(ctx context.Context, params BroadcastRequest) (*BroadcastResponse, error) {
	u := fmt.Sprintf("tx-gateway/v1.1/%d//flashbots", api.chainId)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "POST",
		Params: nil,
		U:      u,
		Body:   body,
	}

	var response BroadcastResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
