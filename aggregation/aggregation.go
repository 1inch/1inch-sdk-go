package aggregation

import (
	"context"
)

type RequestPayload struct {
	Method string
	Params interface{}
	U      string
	Body   []byte
}

type httpExecutor interface {
	ExecuteRequest(ctx context.Context, payload RequestPayload, v interface{}) error
}

type api struct {
	httpExecutor httpExecutor
}

type Client struct {
}
