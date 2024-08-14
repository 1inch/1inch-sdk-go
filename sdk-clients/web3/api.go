package web3

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/1inch/1inch-sdk-go/common"
)

func (api *api) PerformRpcCallAgainstFullNode(ctx context.Context, params PerformRpcCallAgainstFullNodeParams) error {
	u := fmt.Sprintf("web3/%d", api.chainId)

	err := params.Validate()
	if err != nil {
		return err
	}

	body, err := json.Marshal(params.PostChainIdJSONBody)
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

func (api *api) PerformRpcCall(ctx context.Context, params PerformRpcCallParams) error {
	u := fmt.Sprintf("web3/%d/%s", api.chainId, params.PostChainIdNodeTypeParamsNodeType)

	err := params.Validate()
	if err != nil {
		return err
	}

	body, err := json.Marshal(params)
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
