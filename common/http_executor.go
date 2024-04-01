package common

import "context"

type RequestPayload struct {
	Method string
	Params interface{}
	U      string
	Body   []byte
}

type HttpExecutor interface {
	ExecuteRequest(ctx context.Context, payload RequestPayload, v interface{}) error
}
