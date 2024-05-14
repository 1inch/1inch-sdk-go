package traces

import (
	"context"
	"fmt"

	"github.com/1inch/1inch-sdk-go/common"
)

// GetSyncedInterval Get synced interval
func (api *api) GetSyncedInterval(ctx context.Context) (*ReadSyncedIntervalResponseDto, error) {
	u := fmt.Sprintf("traces/v1.0/chain/%d/synced-interval", api.chainId)

	payload := common.RequestPayload{
		Method: "GET",
		Params: nil,
		U:      u,
		Body:   nil,
	}

	var response ReadSyncedIntervalResponseDto
	err := api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetBlockTraceByNumber Get block trace by number
func (api *api) GetBlockTraceByNumber(ctx context.Context, param GetBlockTraceByNumberParam) (*CoreBuiltinBlockTracesDto, error) {
	u := fmt.Sprintf("traces/v1.0/chain/%d/block-trace/%d", api.chainId, param)

	payload := common.RequestPayload{
		Method: "GET",
		Params: nil,
		U:      u,
		Body:   nil,
	}

	var response CoreBuiltinBlockTracesDto
	err := api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetTxTraceByNumberAndHash Get transaction trace by block number and transaction hash
func (api *api) GetTxTraceByNumberAndHash(ctx context.Context, param GetTxTraceByNumberAndHashParams) (*TransactionTraceResponse, error) {
	u := fmt.Sprintf("traces/v1.0/chain/%d/block-trace/%d/tx-hash/%s", api.chainId, param.BlockNumber, param.TransactionHash)

	payload := common.RequestPayload{
		Method: "GET",
		Params: nil,
		U:      u,
		Body:   nil,
	}

	var response TransactionTraceResponse
	err := api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetTxTraceByNumberAndOffset Get transaction trace by block number and offset of transaction in block
func (api *api) GetTxTraceByNumberAndOffset(ctx context.Context, param GetTxTraceByNumberAndOffsetParams) (*TransactionTraceResponse, error) {
	u := fmt.Sprintf("traces/v1.0/chain/%d/block-trace/%d/offset/%d", api.chainId, param.BlockNumber, param.Offset)

	payload := common.RequestPayload{
		Method: "GET",
		Params: nil,
		U:      u,
		Body:   nil,
	}

	var response TransactionTraceResponse
	err := api.httpExecutor.ExecuteRequest(ctx, payload, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
