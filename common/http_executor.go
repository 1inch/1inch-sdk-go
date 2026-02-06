package common

import "context"

type RequestPayload struct {
	Method string
	Params any
	U      string
	Body   []byte
}

type HttpExecutor interface {
	ExecuteRequest(ctx context.Context, payload RequestPayload, v any) error
}
