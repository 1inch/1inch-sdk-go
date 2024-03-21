package aggregation

import (
	"context"
)

type actions interface {
	Swap(ctx context.Context, params SwapTokensParams) error
}

type RequestPayload struct {
	Method string
	Params string
	Body   []byte
}

type httpExecutor interface {
	ExecuteRequest(ctx context.Context, payload RequestPayload, v interface{}) error
}

type apiActions struct {
	httpExecutor httpExecutor
}

type Client struct {
}
