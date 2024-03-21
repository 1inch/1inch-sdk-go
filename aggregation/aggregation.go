package aggregation

import (
	"context"
	"net/http"
)

type actions interface {
	Swap(ctx context.Context, params SwapTokensParams) error
}

type apiActions struct {
	httpClient http.Client
}

type Client struct {
}



