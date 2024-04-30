package balance

import (
	"context"
	"fmt"

	"github.com/1inch/1inch-sdk-go/common"
)

// GetAggregatedBalancesAndAllowances Get balances and allowances by spender for list of EVM addresses
func (api *api) GetAggregatedBalancesAndAllowances(ctx context.Context, params AggregatedBalancesAndAllowancesParams) (*AggregatedBalancesAndAllowancesResponse, error) {
	u := fmt.Sprintf("/balance/v1.2/%d/aggregatedBalancesAndAllowances/%s", api.chainId, params.Spender)

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	payload := common.RequestPayload{
		Method: "GET",
		Params: nil,
		U:      u,
		Body:   nil,
	}

	var response AggregatedBalancesAndAllowancesResponse
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
