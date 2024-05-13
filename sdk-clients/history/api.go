package history

import (
	"context"
	"fmt"

	"github.com/1inch/1inch-sdk-go/common"
)

// GetHistoryEventsByAddress Returns history events for address
func (api *api) GetHistoryEventsByAddress(ctx context.Context, params HistoryEventsByAddressParams) (*HistoryResponseDto, error) {
	u := fmt.Sprintf("history/v2.0/history/%s/events", params.Address)

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

	var response HistoryResponseDto
	err = api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
