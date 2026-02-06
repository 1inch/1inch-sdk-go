package web3

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/1inch/1inch-sdk-go/common"
)

func (api *api) PerformRpcCallAgainstFullNode(ctx context.Context, params PerformRpcCallAgainstFullNodeParams) (map[string]any, error) {
	u := fmt.Sprintf("web3/%d", api.chainId)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(params.PostChainIdJSONBody)
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "POST",
		Params: nil,
		U:      u,
		Body:   body,
	}

	var response map[string]any
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (api *api) PerformRpcCall(ctx context.Context, params PerformRpcCallParams) (map[string]any, error) {
	u := fmt.Sprintf("web3/%d/%s", api.chainId, params.PostChainIdNodeTypeParamsNodeType)

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

	var response map[string]any
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
