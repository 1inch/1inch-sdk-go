package client

import (
	"fmt"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	Response     *http.Response `json:"-"`
	ErrorMessage string         `json:"error"`
	Description  string         `json:"description"`
	StatusCode   int            `json:"statusCode"`
	RequestId    string         `json:"requestId"`
	Meta         []struct {
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"meta"`
}

func (r *ErrorResponse) Error() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("ErrorMessage: %s\n", r.ErrorMessage))
	builder.WriteString(fmt.Sprintf("Description: %s\n", r.Description))
	builder.WriteString(fmt.Sprintf("StatusCode: %d\n", r.StatusCode))
	builder.WriteString(fmt.Sprintf("RequestId: %s\n", r.RequestId))

	if len(r.Meta) > 0 {
		builder.WriteString("Meta:\n")
		for _, meta := range r.Meta {
			builder.WriteString(fmt.Sprintf("  - Value: %s\n", meta.Value))
			builder.WriteString(fmt.Sprintf("    Type: %s\n", meta.Type))
		}
	} else {
		builder.WriteString("Meta: []\n")
	}

	return builder.String()
}
