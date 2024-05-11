package gasprices

import (
	"context"
	"fmt"

	"github.com/1inch/1inch-sdk-go/common"
)

// GetGasPriceEIP1559 Get gas prices for specified network using EIP1559 spec if applicable
func (api *api) GetGasPriceEIP1559(ctx context.Context) (*Eip1559GasPriceResponse, error) {
	u := fmt.Sprintf("/gas-price/v1.5/%d", api.chainId)
	if !isEIP1559Applicable(api.chainId) {
		return nil, fmt.Errorf("chain id %d does not support eip1559", api.chainId)
	}
	payload := common.RequestPayload{
		Method: "GET",
		Params: nil,
		U:      u,
		Body:   nil,
	}

	var response Eip1559GasPriceResponse
	err := api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetGasPriceLegacy Get gas prices for specified network in legacy node (not EIP1559)
func (api *api) GetGasPriceLegacy(ctx context.Context) (*GetGasPriceLegacyResponse, error) {
	u := fmt.Sprintf("/gas-price/v1.5/%d", api.chainId)
	if isEIP1559Applicable(api.chainId) {
		return nil, fmt.Errorf("chain id %d does not support legacy gas price method", api.chainId)
	}
	payload := common.RequestPayload{
		Method: "GET",
		Params: nil,
		U:      u,
		Body:   nil,
	}

	var response GetGasPriceLegacyResponse
	err := api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
